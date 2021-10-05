package olvm

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/vm"
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

type Transaction struct {
	From    action.Address  `json:"from"`
	To      *action.Address `json:"to"`
	Amount  action.Amount   `json:"amount"`
	Data    []byte          `json:"data"`
	Nonce   uint64          `json:"nonce"`
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
	tags := make([]kv.Pair, 0)

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
	tags = append(tags, action.UintTag("tx.nonce", tx.Nonce))
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

// ValidateEthTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (tx *Transaction) ValidateEthTx(keeper balance.AccountKeeper, ethTx *ethtypes.Transaction, minFee *big.Int, local bool) error {
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
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if keeper.GetBalance(tx.From).Cmp(ethTx.Cost()) < 0 {
		return ethcore.ErrInsufficientFunds
	}
	return nil
}

func (tx *Transaction) CheckIntrinsicGas(ethTx *ethtypes.Transaction) (gas uint64, err error) {
	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := ethcore.IntrinsicGas(ethTx.Data(), ethTx.AccessList(), ethTx.To() == nil, true, true)
	if err != nil {
		return intrGas, err
	}
	if ethTx.Gas() < intrGas {
		return intrGas, ethcore.ErrIntrinsicGas
	}
	return intrGas, nil
}

var _ action.Tx = olvmTx{}

type olvmTx struct {
}

func (otx olvmTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	if !ctx.StateDB.IsEnabled {
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
	return true, nil
}

func (otx olvmTx) ProcessCheck(ctx *action.Context, rawTx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing OLVM Transaction for CheckTx", rawTx)

	tx := &Transaction{}
	err := tx.Unmarshal(rawTx.Data)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, action.WrongFee)
	}

	ethTx := tx.tmToEthTx(ctx, rawTx)
	_, err = tx.CheckIntrinsicGas(ethTx)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, action.SkipFee)
	}
	err = tx.ValidateEthTx(
		ctx.StateDB.GetAccountKeeper(),
		tx.tmToEthTx(ctx, rawTx),
		ctx.FeePool.GetOpt().MinFee().Amount.BigInt(),
		true,
	)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, action.SkipFee)
	}
	return ResponseSuccess(tx.Tags(), action.SkipFee)
}

func (otx olvmTx) ProcessDeliver(ctx *action.Context, rawTx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing OLVM Transaction for DeliverTx", rawTx)
	ok, result = runOLVM(ctx, rawTx, false)
	return
}

func (otx olvmTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (ok bool, result action.Response) {
	ok, result = action.ContractFeeHandling(ctx, signedTx, gasUsed, start)
	ctx.Logger.Detailf("Processing OLVM Transaction for BasicFeeHandling: status - %t\n", ok)
	return ok, result
}

func runOLVM(ctx *action.Context, rawTx action.RawTx, isSimulation bool) (bool, action.Response) {
	tx := &Transaction{}
	err := tx.Unmarshal(rawTx.Data)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, action.WrongFee)
	}

	var (
		stateDB *vm.CommitStateDB
		gaspool = new(ethcore.GasPool)
	)

	// NOTE: Just left in case if process check if have a simulation (not recomended)
	if isSimulation {
		stateDB = ctx.StateDB.Copy()
		gaspool = gaspool.AddGas(vm.SimulationBlockGasLimit)
	} else {
		stateDB = ctx.StateDB
		gaspool = gaspool.AddGas(ctx.StateDB.GetAvailableGas())
	}

	ethTx := tx.tmToEthTx(ctx, rawTx)
	intrinsicGas, err := tx.CheckIntrinsicGas(ethTx)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, int64(intrinsicGas))
	}
	err = tx.ValidateEthTx(
		stateDB.GetAccountKeeper(),
		ethTx,
		ctx.FeePool.GetOpt().MinFee().Amount.BigInt(),
		isSimulation,
	)
	if err != nil {
		return ResponseFailed(tx.Tags(), err, int64(intrinsicGas))
	}

	evmTx := vm.NewEVMTransaction(stateDB, gaspool, ctx.Header, tx.From, tx.To, tx.Nonce, tx.Amount.Value.BigInt(), tx.Data, tx.AccessList, uint64(rawTx.Fee.Gas), rawTx.Fee.Price.Value.BigInt(), isSimulation)

	execResult, err := evmTx.Apply()

	if err != nil {
		ctx.Logger.Detailf("OLVM TX: Execution apply VM got err: %s\n", err.Error())
		return ResponseFailed(tx.Tags(), err, int64(intrinsicGas))
	} else if execResult.Failed() {
		ctx.Logger.Detailf("OLVM TX: Execution result VM got err: %s\n", execResult.Err.Error())
		return ResponseFailed(tx.Tags(), execResult.Err, int64(execResult.UsedGas))
	} else {
		tags, err := rlpLogsToEvents(tx.Tags(), stateDB.GetTxLogs())
		if err != nil {
			ctx.Logger.Detail("OLVM TX: Failed to process logs")
			return ResponseFailed(tx.Tags(), err, int64(execResult.UsedGas))
		}
		ctx.Logger.Detail("OLVM TX: Execution result VM got success")
		return ResponseSuccess(tags, int64(execResult.UsedGas))
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

func ResponseSuccess(tags kv.Pairs, gasUsed int64) (bool, action.Response) {
	result := action.Response{
		Events:  action.GetEvent(tags, HandlerName),
		GasUsed: gasUsed,
	}
	return true, result
}

func rlpLogsToEvents(tags kv.Pairs, logs []*ethtypes.Log) (kv.Pairs, error) {
	if len(logs) == 0 {
		return tags, nil
	}
	extTags := make(kv.Pairs, len(tags), len(tags)+len(logs))
	copy(extTags, tags)
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
