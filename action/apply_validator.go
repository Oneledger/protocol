package action

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/serialize"
	"github.com/tendermint/tendermint/libs/common"
	"google.golang.org/grpc/status"
)

type ApplyValidator struct {
	Address   Address
	Stake     Amount
	NodeName  string
	TmAddress Address
	TmPubKey  keys.PublicKey
}

func (apply ApplyValidator) Signers() []Address {
	return []Address{apply.Address.Bytes()}
}

func (apply ApplyValidator) Type() Type {
	return APPLYVALIDATOR
}

func (apply ApplyValidator) Bytes() []byte {
	result, err := serialize.GetSerializer(serialize.NETWORK).Serialize(apply)
	if err != nil {
		logger.Error("send bytes convert failed: ", err)
	}
	return result
}

var _ Tx = applyTx{}

type applyTx struct {
}

func (applyTx) Validate(ctx *Context, msg Msg, fee Fee, memo string, signatures []Signature) (bool, error) {
	apply := msg.(ApplyValidator)

	ok, err := validateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	if msg == nil || len(apply.Address) == 0 {
		return false, ErrMissingData
	}

	if apply.TmAddress == nil {
		return false, ErrMissingData
	}
	_, err = apply.TmPubKey.GetHandler()
	if err != nil {
		return false, ErrInvalidPubkey
	}

	coin := apply.Stake.ToCoin(ctx)
	if coin.LessThanCoin(coin.Currency.NewCoinFromInt(0)) {
		return false, ErrInvalidAmount
	}

	if coin.Currency.Name != "VT" {
		return false, ErrInvalidAmount
	}

	return true, nil
}

func (a applyTx) ProcessCheck(ctx *Context, msg Msg, fee Fee) (bool, Response) {
	apply := msg.(ApplyValidator)

	result, err := checkBalances(ctx, apply.Address, apply.Stake)
	if err != nil {
		return false, Response{Log: err.Error()}
	}
	return result, Response{Tags: apply.Tags()}
}

func checkBalances(ctx *Context, address Address, stake Amount) (bool, error) {

	balances := ctx.Balances

	//check identity's VT is equal to the stake
	balance, err := balances.Get(address, false)
	if err != nil {
		return false, ErrNotEnoughFund
	}
	c, ok := ctx.Currencies.GetCurrencyByName("VT")
	if !ok {
		return false, ErrInvalidAmount
	}
	if balance.GetCoin(c).LessThanCoin(stake.ToCoin(ctx)) {
		return false, ErrNotEnoughFund
	}
	return true, nil
}

func (applyTx) ProcessDeliver(ctx *Context, msg Msg, fee Fee) (bool, Response) {
	apply := msg.(ApplyValidator)

	_, err := checkBalances(ctx, apply.Address, apply.Stake)
	if err != nil {
		return false, Response{Log: err.Error()}
	}

	validators := ctx.Validators
	validator := validators.GetValidator(apply.TmAddress)

	stake := identity.Stake{
		Address: validator.Address,
		Pubkey:  apply.TmPubKey,
	}
	validators.HandleStake(stake)
	return false, Response{}
}

func (applyTx) ProcProcessFee(ctx *Context, fee Fee) (bool, Response) {
	panic("implement me")
	//todo: implement fee charge for apply
	return true, Response{}
}

func (apply ApplyValidator) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(APPLYVALIDATOR.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.from"),
		Value: apply.Address.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.to"),
		Value: s.To.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}
