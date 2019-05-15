package action

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

var _ Msg = Send{}

type Send struct {
	From   Address
	To     Address
	Amount balance.Coin
}

func (s Send) Signers() []Address {
	return []Address{s.From.Bytes()}
}

func (s Send) Type() Type {
	return SEND
}

func (s Send) Bytes() []byte {

	result, err := serialize.GetSerializer(serialize.NETWORK).Serialize(s)
	if err != nil {
		logger.Error("send bytes convert failed: ", err)
	}
	return result
}

var _ Tx = sendTx{}

type sendTx struct {
}

func (sendTx) Validate(ctx *Context, msg Msg, fee Fee, signatures []Signature) (bool, error) {
	send, ok := msg.(*Send)
	if !ok {
		return false, ErrWrongTxType
	}
	if !send.Amount.IsValid() {
		return false, errors.Wrap(ErrInvalidAmount, send.Amount.String())
	}
	if send.From == nil || send.To == nil {
		return false, ErrMissingData
	}

	base := BaseTx{
		send,
		fee,
		signatures,
		"",
	}

	return base.validateBasic()
}

func (sendTx) ProcessCheck(ctx *Context, msg Msg, fee Fee) (bool, Response) {
	logger.Debug("Processing Send Transaction for CheckTx", msg, fee)
	balances := ctx.Balances

	send, ok := msg.(*Send)
	if !ok {
		return false, Response{}
	}

	b, _ := balances.Get(send.From.Bytes(), true)
	if b == nil {
		return false, Response{}
	}

	coin := send.Amount
	if !enoughBalance(*b, coin) {
		return false, Response{}
	}

	return true, Response{Tags: send.Tags()}
}

func (sendTx) ProcessDeliver(ctx *Context, msg Msg, fee Fee) (bool, Response) {
	logger.Debug("Processing Send Transaction for DeliverTx", msg, fee)

	balances := ctx.Balances

	send, ok := msg.(*Send)
	if !ok {
		return false, Response{}
	}

	from, err := balances.Get(send.From.Bytes(), false)
	if err != nil {
		logger.Error("Failed to get the balance of the owner", send.From, "err", err)
		return false, Response{}
	}
	coin := send.Amount

	if !enoughBalance(*from, coin) {
		logger.Debug("Owner balance is not enough", from, send.Amount)
		return false, Response{}
	}

	//change owner balance
	from.MinusCoin(coin)
	err = balances.Set(send.From.Bytes(), *from)
	if err != nil {
		logger.Error("error updating balance in send transaction", err)
		return false, Response{}
	}

	//change receiver balance
	to, err := balances.Get(send.To.Bytes(), false)
	if err != nil {
		logger.Error("failed to get the balance of the receipient", err)
	}
	if to == nil {
		to = balance.NewBalance()
	}
	to.AddCoin(coin)
	err = balances.Set(send.To.Bytes(), *to)
	if err != nil {
		return false, Response{}
	}
	return true, Response{Tags: send.Tags()}
}

func (sendTx) ProcessFee(ctx *Context, fee Fee) (bool, Response) {
	panic("implement me")
	//todo: implement the fee charge for send
	return true, Response{GasWanted: 0, GasUsed: 0}
}

func enoughBalance(b Balance, value balance.Coin) bool {

	if !value.IsValid() {
		return false
	}

	total := balance.NewBalance()
	total.MinusCoin(value)
	if !b.IsEnoughBalance(*total) {
		return false
	}

	return true
}

func (s Send) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(SEND.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.from"),
		Value: s.From.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.to"),
		Value: s.To.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}
