package smart_contract

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Deploy struct {
	From   action.Address `json:"from"`
	Amount action.Amount  `json:"amount"`
	Data   []byte         `json:"data"`
}

func (d Deploy) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Deploy) Unmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d Deploy) Signers() []action.Address {
	return []action.Address{d.From.Bytes()}
}

func (d Deploy) Type() action.Type {
	return action.SC_DEPLOY
}

func (d Deploy) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(d.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: d.From.Bytes(),
	}
	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = scDeployTx{}

type scDeployTx struct {
}

func (s scDeployTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	deploy := &Deploy{}
	err := deploy.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), deploy.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !deploy.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, deploy.Amount.String())
	}

	if deploy.From.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	if len(deploy.Data) == 0 {
		return false, action.ErrMissingData
	}
	return true, nil
}

func (s scDeployTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for CheckTx", tx)
	ok, result = runSCDeploy(ctx, tx)
	return
}

func (s scDeployTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for DeliverTx", tx)
	ok, result = runSCDeploy(ctx, tx)
	return
}

func (s scDeployTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func execDeploy(ctx *action.Context, deploy *Deploy, tx action.RawTx) (*ethcore.ExecutionResult, error) {
	ecfg := action.NewEVMConfig(deploy.From, tx.Fee.Price.Value.BigInt(), uint64(tx.Fee.Gas), []int{})
	vmenv := action.NewEVM(ctx, ecfg)

	ethFrom := ethcmn.BytesToAddress(deploy.From)
	msg := ethtypes.NewMessage(ethFrom, nil, 0, deploy.Amount.Value.BigInt(), uint64(tx.Fee.Gas), tx.Fee.Price.Value.BigInt(), deploy.Data, make(ethtypes.AccessList, 0), false)

	ctx.CommitStateDB.Prepare(ethcmn.Hash{}, ethcmn.Hash{}, 0)
	er, err := ethcore.ApplyMessage(vmenv, msg, new(ethcore.GasPool).AddGas(uint64(uint64(tx.Fee.Gas))))
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	// Ensure any modifications are committed to the state
	// Only delete empty objects if EIP158/161 (a.k.a Spurious Dragon) is in effect
	ctx.CommitStateDB.Finalise(vmenv.ChainConfig().IsEIP158(big.NewInt(ctx.Header.Height)))
	return er, nil
}

func runSCDeploy(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	deploy := &Deploy{}
	err := deploy.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, deploy.Tags(), err)
	}

	if !deploy.Amount.IsValid(ctx.Currencies) {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAmount, deploy.Tags(), errors.New(fmt.Sprint("amount is invalid", deploy.Amount, ctx.Currencies)))
	}

	if _, err := execDeploy(ctx, deploy, tx); err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, deploy.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, deploy.Tags(), "smart_contract_deploy")
}
