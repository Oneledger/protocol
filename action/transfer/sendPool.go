package transfer

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
)

type SendPool struct {
	From     action.Address
	PoolName string
	Amount   action.Amount
}

func (s SendPool) Signers() []action.Address {
	return []action.Address{s.From.Bytes()}
}

func (s SendPool) Type() action.Type {
	return action.SENDPOOL
}

func (s SendPool) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(s.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.From"),
		Value: s.From.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.PoolName"),
		Value: []byte(s.PoolName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (s SendPool) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SendPool) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

var _ action.Msg = &SendPool{}

type sendPoolTx struct {
}

func (sendPoolTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	sendPool := &SendPool{}
	err := sendPool.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), sendPool.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}
	//validate transaction specific field
	if !sendPool.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, sendPool.Amount.String())
	}
	if sendPool.From.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	poolList := getPoolList(ctx)
	if poolList[sendPool.PoolName] == nil {
		return false, action.ErrPoolDoesNotExist
	}

	return true, nil
}

func (sendPoolTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runSendPool(ctx, tx)
}

func (sendPoolTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runSendPool(ctx, tx)
}

func (sendPoolTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

var _ action.Tx = &sendPoolTx{}

func runSendPool(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	sendPool := &SendPool{}
	err := sendPool.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, sendPool.Tags(), err)
	}

	// Get Coin
	coin := sendPool.Amount.ToCoin(ctx.Currencies)
	// Deduct from Sender
	err = ctx.Balances.MinusFromAddress(sendPool.From.Bytes(), coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, sendPool.Tags(), err)
	}
	// Get Pool Address
	toPool := getPoolList(ctx)[sendPool.PoolName]
	err = ctx.Balances.AddToAddress(toPool, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, sendPool.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, sendPool.Tags(), "Send to Pool Success")
}

func getPoolList(ctx *action.Context) map[string]action.Address {
	poolList := map[string]action.Address{}
	poolList["BountyProgram"] = action.Address(ctx.ProposalMasterStore.Proposal.GetOptions().BountyProgramAddr)
	poolList["FeePool"] = action.Address(fees.POOL_KEY)
	return poolList
}
