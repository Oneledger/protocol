/*

 */

package ons

import (
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

/*
	DomainSend info
*/
// DomainSend is a struct which encapsulates information required in a send to domain transaction.
// This is struct is serialized according to network strategy before sending over network.
type DomainSend struct {
	From       action.Address
	DomainName string
	Amount     action.Amount
}

// DomainSend satisfies the action.Msg interface
var _ action.Msg = DomainSend{}

var _ Ons = DomainSend{}

func (s DomainSend) OnsName() string {
	return s.DomainName
}

// Type method gives the transaction type of the
func (s DomainSend) Type() action.Type {
	return action.DOMAIN_SEND
}

// Signers gives the list of addresses of the signers (useful for verification_
func (s DomainSend) Signers() []action.Address {
	return []action.Address{s.From.Bytes()}
}

func (s DomainSend) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.DOMAIN_SEND.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: s.From.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(s.DomainName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

/*


	domainSendTx

*/
type domainSendTx struct {
}

var _ action.Tx = domainSendTx{}

func (domainSendTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee,
	memo string, signatures []action.Signature) (bool, error) {

	// validate basic signature
	ok, err := action.ValidateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	// validate transaction specific field
	send, ok := msg.(*DomainSend)
	if !ok {
		return false, action.ErrWrongTxType
	}

	// validate the amount
	if !send.Amount.IsValid(ctx) {
		return false, errors.Wrap(action.ErrInvalidAmount, send.Amount.String())
	}

	// validate the sender and receiver are not nil
	if send.From == nil || send.DomainName == "" {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainSendTx) ProcessCheck(ctx *action.Context, msg action.Msg,
	fee action.Fee) (bool, action.Response) {

	ctx.Logger.Debug("Processing Send Transaction for CheckTx", msg, fee)

	balances := ctx.Balances

	send, ok := msg.(*DomainSend)
	if !ok {
		return false, action.Response{Log: "failed to cast msg"}
	}

	// validate amount and get coin representation
	if !send.Amount.IsValid(ctx) {
		log := fmt.Sprint("amount is invalid", send.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}
	coin := send.Amount.ToCoin(ctx)

	// get sender balance
	b, _ := balances.Get(send.From.Bytes(), true)
	if b == nil {
		return false, action.Response{Log: "failed to get balance for sender"}
	}

	// check if balance is enough
	_, err := b.MinusCoin(coin)
	if err != nil {
		log := fmt.Sprint("error in minus coin", err)
		return false, action.Response{Log: log}
	}

	// check if domain of receiver exists
	receiverExists := ctx.Domains.Exists(send.DomainName)
	if !receiverExists {
		log := fmt.Sprintf("receiving domain does")
		return false, action.Response{Log: log}
	}

	return true, action.Response{Tags: send.Tags()}
}

func (domainSendTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	ctx.Logger.Debug("Processing Send to Domain Transaction for DeliverTx", msg, fee)

	balances := ctx.Balances

	send, ok := msg.(*DomainSend)
	if !ok {
		return false, action.Response{}
	}

	from, err := balances.Get(send.From.Bytes(), false)
	if err != nil {
		log := fmt.Sprint("Failed to get the balance of the owner ", send.From, "err", err)
		return false, action.Response{Log: log}
	}

	if !send.Amount.IsValid(ctx) {
		log := fmt.Sprint("amount is invalid", send.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	coin := send.Amount.ToCoin(ctx)

	//change owner balance
	from, err = from.MinusCoin(coin)
	if err != nil {
		log := fmt.Sprint("error in minus coin", err)
		return false, action.Response{Log: log}
	}
	err = balances.Set(send.From.Bytes(), *from)
	if err != nil {
		log := fmt.Sprint("error updating balance in send transaction", err)
		return false, action.Response{Log: log}
	}

	domain, err := ctx.Domains.Get(send.DomainName, true)
	if err != nil {
		log := fmt.Sprint("error getting domain:", err)
		return false, action.Response{Log: log}
	}
	if !domain.ActiveFlag {
		log := fmt.Sprint("domain inactive")
		return false, action.Response{Log: log}
	}
	if len(domain.AccountAddress) == 0 {
		log := fmt.Sprint("domain account address not set")
		return false, action.Response{Log: log}
	}

	to := domain.AccountAddress
	//change receiver balance
	toBalance, err := balances.Get(to, false)
	if err != nil {
		ctx.Logger.Error("failed to get the balance of the receipient", err)
	}
	if toBalance == nil {
		toBalance = balance.NewBalance()
	}
	toBalance.AddCoin(coin)
	err = balances.Set(to, *toBalance)
	if err != nil {
		return false, action.Response{Log: "balance set failed"}
	}
	return true, action.Response{Tags: send.Tags()}
}

func (domainSendTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
	// TODO: implement the fee charge for send
	return true, action.Response{GasWanted: 0, GasUsed: 0}
}
