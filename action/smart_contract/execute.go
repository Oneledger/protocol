package smart_contract

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Execute struct {
	From   action.Address  `json:"from"`
	To     *action.Address `json:"to"`
	Amount action.Amount   `json:"amount"`
	Data   []byte          `json:"data"`
}

func (e Execute) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Execute) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e Execute) Signers() []action.Address {
	return []action.Address{e.From.Bytes()}
}

func (e Execute) Type() action.Type {
	return action.SC_EXECUTE
}

func (e Execute) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(e.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.origin"),
		Value: e.From.Bytes(),
	}
	tags = append(tags, tag, tag2)
	if e.To != nil {
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.contract"),
			Value: e.To.Bytes(),
		})
	}
	return tags
}

var _ action.Tx = scExecuteTx{}

type scExecuteTx struct {
}

func (s scExecuteTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	execute := &Execute{}
	err := execute.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), execute.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !execute.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, execute.Amount.String())
	}

	if execute.From.Err() != nil || execute.To != nil && execute.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	if len(execute.Data) == 0 {
		return false, action.ErrMissingData
	}
	return true, nil
}

func (s scExecuteTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for CheckTx", tx)
	ok, result = runSCExecute(ctx, tx)
	// NOTE: clear internal logs and journal as it is only simulation
	ctx.StateDB.Reset(ethcmn.Hash{})
	return
}

func (s scExecuteTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for DeliverTx", tx)
	ok, result = runSCExecute(ctx, tx)
	return
}

func (s scExecuteTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	// checking the difference between nominal and used and get it's price
	// FIXME: Check gas calculation on the correctness
	return action.ContractFeeHandling(ctx, signedTx, storage.Gas(signedTx.Fee.Gas), gasUsed)
}

func runSCExecute(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	execute := &Execute{}
	err := execute.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !execute.Amount.IsValid(ctx.Currencies) {
		log := fmt.Sprint("amount is invalid", execute.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	evmTx := action.NewEVMTransaction(ctx.StateDB, ctx.Header, execute.From, execute.To, execute.Amount.Value.BigInt(), execute.Data)
	tags := execute.Tags()
	vmenv := evmTx.NewEVM()
	execResult, err := evmTx.Apply(vmenv, tx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, tags, err)
	}
	if execResult.Failed() {
		tags = append(tags, action.UintTag("tx.status", ethtypes.ReceiptStatusFailed))
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.error"),
			Value: []byte(execResult.Err.Error()),
		})
	} else {
		tags = append(tags, action.UintTag("tx.status", ethtypes.ReceiptStatusSuccessful))
		tags = append(tags, kv.Pair{
			Key:   []byte("tx.data"),
			Value: []byte(execResult.ReturnData),
		})
		if execute.To == nil {
			ctx.Logger.Debugf("Contract created: %s\n", execResult.ContractAddress)
			tags = append(tags, kv.Pair{
				Key:   []byte("tx.contract"),
				Value: []byte(execResult.ContractAddress.Bytes()),
			})
		}
	}
	ctx.Logger.Debugf("Contract TX: status ok - %t, used gas - %d, err - %s, data - %s\n", execResult.Failed(), execResult.UsedGas, execResult.Err, execResult.ReturnData)
	return true, action.Response{
		Events:  action.GetEvent(tags, "smart_contract_execute"),
		GasUsed: int64(execResult.UsedGas),
	}
}
