package nexus

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/status_codes"
	"github.com/Oneledger/protocol/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

const (
	// txSlotSize is used to calculate how many data slots a single transaction
	// takes up based on its size. The slots are used as DoS protection, ensuring
	// that validating a new transaction remains a constant operation (in reality
	// O(maxslots), where max slots are 4 currently).
	txSlotSize = 32 * 1024

	// txMaxSize is the maximum size a single transaction can have. This field has
	// non-trivial consequences: larger transactions are significantly harder and
	// more expensive to propagate; larger transactions also take more resources
	// to validate whether they fit into the pool or not.
	txMaxSize = 4 * txSlotSize // 128KB
)

type Nexus struct {
	From   action.Address  `json:"from"`
	To     *action.Address `json:"to"`
	Amount action.Amount   `json:"amount"`
	Data   []byte          `json:"data"`
	Nonce  uint64          `json:"nonce"`
}

func (n Nexus) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *Nexus) Unmarshal(data []byte) error {
	return json.Unmarshal(data, n)
}

func (n Nexus) Signers() []action.Address {
	return []action.Address{n.From.Bytes()}
}

func (n Nexus) Type() action.Type {
	return action.NEXUS
}

func (n Nexus) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(n.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.from"),
		Value: n.From.Bytes(),
	}
	tags = append(tags, tag, tag2)
	if n.To != nil {
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.to"),
			Value: n.To.Bytes(),
		})
	}
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.value"),
		Value: n.Amount.Value.BigInt().Bytes(),
	})
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.data"),
		Value: n.Data,
	})
	tags = append(tags, action.UintTag("tx.nonce", n.Nonce))
	return tags
}

func (n Nexus) isSendTx() bool {
	if len(n.Data) > 0 {
		return false
	}
	return true
}

func (n *Nexus) tmToEthTx(ctx *action.Context, tx action.RawTx) *ethtypes.Transaction {
	var to *ethcmn.Address
	if n.To != nil {
		to = new(ethcmn.Address)
		*to = ethcmn.BytesToAddress(n.To.Bytes())
	}
	ethTx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    n.Nonce,
		To:       to,
		Value:    n.Amount.Value.BigInt(),
		Gas:      uint64(tx.Fee.Gas),
		GasPrice: tx.Fee.Price.Value.BigInt(),
		Data:     n.Data,
	})
	return ethTx
}

func (n *Nexus) getEthSigner(ctx *action.Context) ethtypes.Signer {
	return ethtypes.NewEIP155Signer(utils.HashToBigInt(ctx.Header.ChainID))
}

func (n *Nexus) validateSigner(ctx *action.Context, tx action.SignedTx) error {
	if len(tx.Signatures) != 1 {
		return action.ErrInvalidSignature
	}

	//validate basic signature
	signer := n.getEthSigner(ctx)
	ethTx := n.tmToEthTx(ctx, tx.RawTx)
	ethTx, err := ethTx.WithSignature(signer, tx.Signatures[0].Signed)
	if err != nil {
		return action.ErrInvalidSignature
	}
	// Make sure the transaction is signed properly.
	addr, err := signer.Sender(ethTx)
	if err != nil {
		return ethcore.ErrInvalidSender
	}
	if !n.From.Equal(addr.Bytes()) {
		return action.ErrInvalidAddress
	}
	return nil
}

// validateCheckTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (n *Nexus) validateTx(ctx *action.Context, tmTx action.RawTx, local bool) error {
	tx := n.tmToEthTx(ctx, tmTx)

	keeper := ctx.StateDB.GetAccountKeeper()
	currentState := keeper.GetState()
	// Accept only legacy transactions until EIP-2718/2930 activates.
	if tx.Type() != ethtypes.LegacyTxType {
		return ethtypes.ErrTxTypeNotSupported
	}
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > txMaxSize {
		return ethcore.ErrOversizedData
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return ethcore.ErrNegativeValue
	}
	currentMaxGas := uint64(currentState.GetCalculator().GetLimit())
	// Ensure the transaction doesn't exceed the current block limit gas.
	if currentMaxGas < tx.Gas() {
		return ethcore.ErrGasLimit
	}
	// Drop non-local transactions under our own minimal accepted gas price
	if !local && tx.GasPrice().Cmp(ctx.FeePool.GetOpt().MinFee().Amount.BigInt()) < 0 {
		return ethcore.ErrUnderpriced
	}
	// Ensure the transaction adheres to nonce ordering
	if keeper.GetNonce(n.From) > tx.Nonce() {
		return ethcore.ErrNonceTooLow
	}
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if keeper.GetBalance(n.From).Cmp(tx.Cost()) < 0 {
		return ethcore.ErrInsufficientFunds
	}
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := ethcore.IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, true, true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return ethcore.ErrIntrinsicGas
	}
	return nil
}

var _ action.Tx = nexusTx{}

type nexusTx struct {
}

func (n nexusTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	// zero means stateDB does not receive a block number so nexus update has not been applied
	if ctx.StateDB.GetCurrentHeight() == 0 {
		return false, errors.Errorf("not enabled")
	}

	nexus := &Nexus{}
	err := nexus.Unmarshal(tx.Data)
	if err != nil {
		return false, err
	}

	//validate basic signature
	err = nexus.validateSigner(ctx, tx)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !nexus.Amount.IsValid(ctx.Currencies) || nexus.Amount.Currency != action.DEFAULT_CURRENCY {
		return false, action.ErrInvalidCurrency
	}

	if nexus.From.Err() != nil || nexus.To != nil && nexus.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}
	return true, nil
}

func (n nexusTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Nexus Transaction for CheckTx", tx)
	nexus := &Nexus{}
	err := nexus.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, nexus.Tags(), err)
	}
	err = nexus.validateTx(ctx, tx, true)
	if err != nil {
		return false, action.Response{
			Events: action.GetEvent(nexus.Tags(), err.Error()),
			Log:    err.Error(),
		}
	}
	// NOTE: We do not need to run a logic on smart contract, but only calculate a gas
	// maybe it will be required a stateDB copy without deliver and clear it's dirties
	ok = true
	result = action.Response{GasUsed: -1}
	return
}

func (n nexusTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Nexus Transaction for DeliverTx", tx)
	ok, result = runNexus(ctx, tx)
	return
}

func (n nexusTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	// -1 if ProcessCheck
	if gasUsed == -1 {
		return true, action.Response{}
	}
	ok, result := action.ContractFeeHandling(ctx, signedTx, gasUsed)
	ctx.Logger.Detailf("Processing Nexus Transaction for BasicFeeHandling: %+v, status - %t\n", result, ok)
	return ok, result
}

func runNexus(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	nexus := &Nexus{}
	err := nexus.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	ctx.Logger.Detail("Nexus tx is send:", nexus.isSendTx())

	tags := nexus.Tags()

	err = nexus.validateTx(ctx, tx, false)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, status_codes.ProtocolError{
			Msg:  err.Error(),
			Code: status_codes.InvalidParams,
		}, tags, err)
	}

	evmTx := action.NewEVMTransaction(ctx.StateDB, ctx.Header, nexus.From, nexus.To, nexus.Nonce, nexus.Amount.Value.BigInt(), nexus.Data, false)

	vmenv := evmTx.NewEVM()
	execResult, err := evmTx.Apply(vmenv, tx)
	if err != nil {
		ctx.Logger.Debugf("Execution apply VM got err: %s\n", err.Error())
		tags = append(tags, action.UintTag("tx.status", ethtypes.ReceiptStatusFailed))
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.error"),
			Value: []byte(err.Error()),
		})
		return helpers.LogAndReturnFalse(ctx.Logger, status_codes.ProtocolError{
			Msg:  err.Error(),
			Code: status_codes.TxErrVMExecution,
		}, tags, err)
	}

	if execResult.Failed() {
		ctx.Logger.Debugf("Execution result got err: %s\n", execResult.Err.Error())
		tags = append(tags, action.UintTag("tx.status", ethtypes.ReceiptStatusFailed))
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.error"),
			Value: []byte(execResult.Err.Error()),
		})
	} else {
		tags = append(tags, action.UintTag("tx.status", ethtypes.ReceiptStatusSuccessful))
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.returnData"),
			Value: []byte(execResult.ReturnData),
		})
		if nexus.To == nil && len(execResult.ContractAddress.Bytes()) > 0 {
			ctx.Logger.Debugf("Contract created: %s\n", keys.Address(execResult.ContractAddress.Bytes()))
			tags = append(tags, kv.Pair{
				Key:   []byte("tx.contract"),
				Value: []byte(execResult.ContractAddress.Bytes()),
			})
		}
	}
	ctx.Logger.Debugf("Contract TX: status ok - %t, used gas - %d\n", !execResult.Failed(), execResult.UsedGas)
	return true, action.Response{
		Events:  action.GetEvent(tags, "nexus"),
		GasUsed: int64(execResult.UsedGas),
	}
}
