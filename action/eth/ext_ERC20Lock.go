package eth

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	ethchaindriver "github.com/Oneledger/protocol/chains/ethereum"
)

type ERC20Lock struct {
	Locker action.Address
	ETHTxn []byte  // Raw Transaction for Locking Tokens
}

var _ action.Msg = &ERC20Lock{}

func (E ERC20Lock) Signers() []action.Address {
	return []action.Address{E.Locker}
}

func (E ERC20Lock) Type() action.Type {
	return action.ERC20_LOCK
}

func (E ERC20Lock) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(E.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: E.Locker.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (E ERC20Lock) Marshal() ([]byte, error) {
	return json.Marshal(E)
}

func (E ERC20Lock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, E)
}

type ethERC20LockTx struct {
}
var _ action.Tx = ethERC20LockTx{}

func (e ethERC20LockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	// unmarshal the tx message
	lock := &ERC20Lock{}
	err := lock.Unmarshal(signedTx.Data)
	if err != nil {
		ctx.Logger.Error("error wrong tx type")
		return false, errors.Wrap(err, action.ErrWrongTxType.Msg)
	}

	err = action.ValidateBasic(signedTx.RawBytes(), lock.Signers(), signedTx.Signatures)
	if err != nil {
		ctx.Logger.Error("validate basic failed", err)
		return false, err
	}

	// validate fee
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		ctx.Logger.Error("validate fee failed", err)
		return false, err
	}

	return true, nil
}

func (e ethERC20LockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Lock(ctx,tx)
}

func (e ethERC20LockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runERC20Lock(ctx,tx)
}

func (e ethERC20LockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	fmt.Println("processFee")
	return true,action.Response{Log:"processFee"}
}


func runERC20Lock (ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	erc20lock := &ERC20Lock{}
	err := erc20lock.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("wrong tx type", err)
		return false, action.Response{Log: "wrong tx type"}
	}
	ethTx, err := ethchaindriver.DecodeTransaction(erc20lock.ETHTxn)
	if err != nil {
		ctx.Logger.Error("decode eth txn err", err)
		return false, action.Response{
			Log: "decode eth txn error" + err.Error(),
		}
	}
    fmt.Println(ethTx)
	if err != nil {
		ctx.Logger.Error("Unable to Verify Data for Ethereum Lock")
		return false, action.Response{
			Log: "Unable to verify lock trasaction" + err.Error(),
		}
	}




	return true, action.Response{
		Log:"runlock",
	}
}

