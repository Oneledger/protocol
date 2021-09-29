package olvm

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/vm"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
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

	// enable VM debug mode to print opcodes, gas consumed etc.
	enableVMDebug = false
)

type Transaction struct {
	Nonce   uint64          `json:"nonce"`
	From    action.Address  `json:"from"`
	To      *action.Address `json:"to"`
	Amount  action.Amount   `json:"amount"`
	Data    []byte          `json:"data"`
	ChainID *big.Int        `json:"chainID"`
	// for future
	TxType     int64                `json:"type"`
	AccessList *ethtypes.AccessList `json:"accessList"`
}

func (tx Transaction) Marshal() ([]byte, error) {
	return json.Marshal(tx)
}

func (tx *Transaction) Unmarshal(data []byte) error {
	return json.Unmarshal(data, tx)
}

func (tx Transaction) Signers() []action.Address {
	return []action.Address{tx.From.Bytes()}
}

func (tx Transaction) Type() action.Type {
	return action.OLVM
}

func (tx Transaction) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0, 6)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(tx.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.from"),
		Value: tx.From.Bytes(),
	}
	tags = append(tags, tag, tag2)
	if tx.To != nil {
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.to"),
			Value: tx.To.Bytes(),
		})
	}
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.value"),
		Value: tx.Amount.Value.BigInt().Bytes(),
	})
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.data"),
		Value: tx.Data,
	})
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.nonce"),
		Value: []byte(strconv.FormatUint(tx.Nonce, 10)),
	})
	return tags
}

func (tx Transaction) isSendTx() bool {
	if len(tx.Data) > 0 {
		return false
	}
	return true
}

func (tx *Transaction) tmToEthTx(ctx *action.Context, raw action.RawTx) *ethtypes.Transaction {
	var to *ethcmn.Address
	if tx.To != nil {
		to = new(ethcmn.Address)
		*to = ethcmn.BytesToAddress(tx.To.Bytes())
	}
	ethTx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    tx.Nonce,
		To:       to,
		Value:    tx.Amount.Value.BigInt(),
		Gas:      uint64(raw.Fee.Gas),
		GasPrice: raw.Fee.Price.Value.BigInt(),
		Data:     tx.Data,
	})
	return ethTx
}

func (tx *Transaction) getEthSigner(ctx *action.Context) ethtypes.Signer {
	return ethtypes.NewEIP155Signer(utils.HashToBigInt(ctx.Header.ChainID))
}

func (tx *Transaction) validateSigner(ctx *action.Context, signedTx action.SignedTx) error {
	if len(signedTx.Signatures) != 1 {
		return errors.New("invalid signatures count")
	}

	//validate basic signature
	signer := tx.getEthSigner(ctx)
	ethTx := tx.tmToEthTx(ctx, signedTx.RawTx)
	ethTx, err := ethTx.WithSignature(signer, signedTx.Signatures[0].Signed)
	if err != nil {
		return err
	}
	// Accept only tx with matched chain ID
	if ethTx.ChainId().Cmp(tx.ChainID) != 0 {
		return ethtypes.ErrInvalidChainId
	}
	// Make sure the transaction is signed properly.
	addr, err := signer.Sender(ethTx)
	if err != nil {
		return ethcore.ErrInvalidSender
	}

	if !tx.From.Equal(addr.Bytes()) {
		return errors.New("mismatch sender")
	}
	return nil
}

// validateEthTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (tx *Transaction) validateEthTx(keeper balance.AccountKeeper, ethTx *ethtypes.Transaction, minFee *big.Int, local bool) error {
	// Accept only legacy transactions until EIP-2718/2930 activates.
	if ethTx.Type() != ethtypes.LegacyTxType {
		return ethtypes.ErrTxTypeNotSupported
	}
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(ethTx.Size()) > txMaxSize {
		return ethcore.ErrOversizedData
	}
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if ethTx.Value().Sign() < 0 {
		return ethcore.ErrNegativeValue
	}
	// Ensure the transaction doesn't exceed the current block limit gas.
	// NOTE: Change this in future when own tx pool will be implemented
	if vm.SimulationBlockGasLimit < ethTx.Gas() {
		return ethcore.ErrGasLimit
	}
	// Drop non-local transactions under our own minimal accepted gas price
	if !local && ethTx.GasPrice().Cmp(minFee) < 0 {
		return ethcore.ErrUnderpriced
	}
	// Ensure the transaction adheres to nonce ordering
	stNonce := keeper.GetNonce(tx.From)
	if msgNonce := ethTx.Nonce(); stNonce > msgNonce {
		return ethcore.ErrNonceTooLow
	}
	//  else if stNonce < msgNonce {
	// 	// TODO: Remove this from here and add sending to pool again
	// 	return ethcore.ErrNonceTooHigh
	// }
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if keeper.GetBalance(tx.From).Cmp(ethTx.Cost()) < 0 {
		return ethcore.ErrInsufficientFunds
	}
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := vm.IntrinsicGas(ethTx.Data(), ethTx.AccessList(), ethTx.To() == nil)
	if err != nil {
		return err
	}
	if ethTx.Gas() < intrGas {
		return ethcore.ErrIntrinsicGas
	}
	return nil
}

var _ action.Tx = olvmTx{}

type olvmTx struct {
}

func (otx olvmTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	if !ctx.StateDB.Enabled() {
		return false, errors.Errorf("not enabled")
	}

	tx := &Transaction{}
	err := tx.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}

	//validate basic signature
	err = tx.validateSigner(ctx, signedTx)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !tx.Amount.IsValid(ctx.Currencies) || tx.Amount.Currency != action.DEFAULT_CURRENCY {
		return false, action.ErrInvalidCurrency
	}

	if tx.From.Err() != nil || tx.To != nil && tx.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	err = tx.validateEthTx(
		ctx.StateDB.GetAccountKeeper(),
		tx.tmToEthTx(ctx, signedTx.RawTx),
		ctx.FeePool.GetOpt().MinFee().Amount.BigInt(),
		true,
	)
	if err != nil {
		return false, err
	}

	memoNonce, err := strconv.ParseUint(signedTx.Memo, 10, 0)
	if err != nil {
		return false, err
	}

	// double spend protection
	if memoNonce != tx.Nonce {
		return false, errors.New("wrong memo for nonce")
	}

	return true, nil
}

func (otx olvmTx) ProcessCheck(ctx *action.Context, rawTx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing OLVM Transaction for CheckTx", rawTx)
	return ResponseSuccess(action.GetEvent(make(kv.Pairs, 0), HandlerName), action.SkipFee)
}

func (otx olvmTx) ProcessDeliver(ctx *action.Context, rawTx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing OLVM Transaction for DeliverTx", rawTx)
	ok, result = runOLVM(ctx, rawTx)
	return
}

func (otx olvmTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (ok bool, result action.Response) {
	ok, result = action.ContractFeeHandling(ctx, signedTx, gasUsed, start)
	ctx.Logger.Detailf("Processing OLVM Transaction for BasicFeeHandling: status - %t\n", ok)
	return ok, result
}

func runOLVM(ctx *action.Context, rawTx action.RawTx) (bool, action.Response) {
	tx := &Transaction{}
	err := tx.Unmarshal(rawTx.Data)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, action.WrongFee)
	}

	evmTx := vm.NewEVMTransaction(
		ctx.StateDB,
		new(ethcore.GasPool).AddGas(ctx.StateDB.GetAvailableGas()),
		ctx.Header,
		tx.From,
		tx.To,
		tx.Nonce,
		tx.Amount.Value.BigInt(),
		tx.Data,
		tx.AccessList,
		uint64(rawTx.Fee.Gas),
		rawTx.Fee.Price.Value.BigInt(),
		false,
	)

	evmTx.SetVMDebug(enableVMDebug)

	tags := tx.Tags()

	if execResult, err := evmTx.Apply(); err != nil {
		ctx.Logger.Detailf("OLVM TX: Execution apply VM got err: %s\n", err.Error())

		tags = responseLogs(tags, ethtypes.ReceiptStatusFailed, err)
		return ResponseFailed(tags, err, action.WrongFee)
	} else if execResult.Failed() {
		ctx.Logger.Detailf("OLVM TX: Execution result VM got err: %s\n", execResult.Err.Error())

		tags = responseLogs(tags, ethtypes.ReceiptStatusFailed, execResult.Err)
		return ResponseSuccess(action.GetEvent(tags, HandlerName), int64(execResult.UsedGas))
	} else {
		ctx.Logger.Detail("OLVM TX: Execution result VM got success")

		logTags, err := rlpLogsToTags(ctx.StateDB.GetTxLogs())
		if err != nil {
			tags = responseLogs(tags, ethtypes.ReceiptStatusFailed, err)
			ctx.Logger.Detail("OLVM TX: Failed to process logs")
			return ResponseFailed(tags, err, int64(execResult.UsedGas))
		}

		tags = responseLogs(tags, ethtypes.ReceiptStatusSuccessful, nil)
		if len(execResult.ContractAddress.Bytes()) > 0 {
			ctx.Logger.Detailf("Contract created: %s\n", execResult.ContractAddress)
			tags = append(tags, kv.Pair{
				Key:   []byte("tx.contract"),
				Value: execResult.ContractAddress.Bytes(),
			})
		}

		events := action.GetEvent(tags, HandlerName)
		if len(logTags) > 0 {
			logEvent := types.Event{
				Type:       fmt.Sprintf("%s.logs", HandlerName),
				Attributes: logTags,
			}
			events = append(events, logEvent)
		}
		return ResponseSuccess(events, int64(execResult.UsedGas))
	}
}

func ResponseFailed(tags kv.Pairs, err error, gasUsed int64) (bool, action.Response) {
	result := action.Response{
		Events:  action.GetEvent(tags, HandlerName),
		Log:     action.ErrInvalidVmExecution.Wrap(err).Marshal(),
		GasUsed: gasUsed,
	}
	return false, result
}

func ResponseSuccess(events []types.Event, gasUsed int64) (bool, action.Response) {
	result := action.Response{
		Events:  events,
		GasUsed: gasUsed,
	}
	return true, result
}

func responseLogs(tags []kv.Pair, status uint64, err error) []kv.Pair {
	tags = append(tags, kv.Pair{
		Key:   []byte("tx.status"),
		Value: []byte(strconv.Itoa(int(status))),
	})
	if err != nil {
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.error"),
			Value: []byte(err.Error()),
		})
	}
	return tags
}

func rlpLogsToTags(logs []*ethtypes.Log) (kv.Pairs, error) {
	if len(logs) == 0 {
		return kv.Pairs{}, nil
	}
	extTags := make(kv.Pairs, 0, len(logs)+1)
	extTags = append(extTags, kv.Pair{
		Key:   []byte("tx.bloom"),
		Value: vm.LogsBloom(logs),
	})
	for _, log := range logs {
		rlpLog, err := vm.RLPLogConvert(*log).Encode()
		if err != nil {
			return kv.Pairs{}, err
		}
		extTags = append(extTags, kv.Pair{
			Key:   []byte(fmt.Sprintf("tx.logs.%d", log.Index)),
			Value: rlpLog,
		})
	}
	return extTags, nil
}
