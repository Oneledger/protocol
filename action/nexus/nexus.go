package nexus

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
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

func (n Nexus) validateSigner(ctx *action.Context, tx action.SignedTx) error {
	if len(tx.Signatures) != 1 {
		return action.ErrInvalidSignature
	}

	//validate basic signature
	to := ethcmn.Address{}
	if n.To != nil {
		to = ethcmn.BytesToAddress(n.To.Bytes())
	}
	chainId := utils.HashToBigInt(ctx.Header.ChainID)
	signer := ethtypes.NewEIP155Signer(chainId)
	ethTx := ethtypes.NewTransaction(n.Nonce, to, n.Amount.Value.BigInt(), uint64(tx.Fee.Gas), tx.Fee.Price.Value.BigInt(), n.Data)
	ethTx, err := ethTx.WithSignature(signer, tx.Signatures[0].Signed)
	if err != nil {
		return errors.Wrap(action.ErrInvalidSignature, err.Error())
	}
	addr, err := signer.Sender(ethTx)
	if err != nil {
		return action.ErrInvalidSignature
	}
	if !n.From.Equal(addr.Bytes()) {
		return action.ErrInvalidAddress
	}
	return nil
}

var _ action.Tx = nexusTx{}

type nexusTx struct {
}

func (n nexusTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	// zero means stateDB does not receive a block number so nexus update has not been applied
	if ctx.StateDB.GetCurrentHeight() == 0 {
		return false, errors.Wrap(action.ErrWrongTxType, "not enabled")
	}

	nexus := &Nexus{}
	err := nexus.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
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
	if !nexus.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, nexus.Amount.String())
	}

	if nexus.From.Err() != nil || nexus.To != nil && nexus.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	if len(nexus.Data) == 0 {
		return false, action.ErrMissingData
	}
	return true, nil
}

func (n nexusTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Nexus Transaction for CheckTx", tx)
	// NOTE: We do not need to run a logic, but only calculate a gas
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

	if nexus.isSendTx() {
		return runSend(ctx, tx, nexus)
	}
	return runSmartContract(ctx, tx, nexus)
}

func runSend(ctx *action.Context, tx action.RawTx, nexus *Nexus) (bool, action.Response) {
	coin := nexus.Amount.ToCoin(ctx.Currencies)

	err := ctx.Balances.MinusFromAddress(nexus.From.Bytes(), coin)
	if err != nil {
		log := fmt.Sprint("error debiting balance in send transaction ", nexus.From, "err", err)
		return false, action.Response{Log: log}
	}

	err = ctx.Balances.AddToAddress(nexus.To.Bytes(), coin)
	if err != nil {
		log := fmt.Sprint("error crediting balance in send transaction ", nexus.From, "err", err)
		return false, action.Response{Log: log}
	}

	return true, action.Response{
		Events:  action.GetEvent(nexus.Tags(), "nexus_send"),
		GasUsed: int64(ethparams.TxGas),
	}
}

func runSmartContract(ctx *action.Context, tx action.RawTx, nexus *Nexus) (bool, action.Response) {
	evmTx := action.NewEVMTransaction(ctx.StateDB, ctx.Header, nexus.From, nexus.To, nexus.Nonce, nexus.Amount.Value.BigInt(), nexus.Data, false)
	tags := nexus.Tags()
	vmenv := evmTx.NewEVM()
	execResult, err := evmTx.Apply(vmenv, tx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, tags, err)
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
			Key:   []byte("tx.data"),
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
		Events:  action.GetEvent(tags, "nexus_contract"),
		GasUsed: int64(execResult.UsedGas),
	}
}
