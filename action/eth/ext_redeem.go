package eth

import (
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
)

var _ action.Msg = &Redeem{}

type Redeem struct {
	Owner  action.Address //User Oneledger address
	To     action.Address //User Ethereum address
	ETHTxn []byte
}

func (r Redeem) Signers() []action.Address {
	return []action.Address{r.Owner}
}

func (r Redeem) Type() action.Type {
	return action.ETH_REDEEM
}

func (r Redeem) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: r.Owner,
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (r Redeem) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Redeem) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

var _ action.Tx = ethRedeemTx{}

type ethRedeemTx struct {
}

func (ethRedeemTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	redeem := &Redeem{}
	err := redeem.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(err, action.ErrWrongTxType.Error())
	}
	err = action.ValidateBasic(signedTx.RawBytes(), redeem.Signers(), signedTx.Signatures)
	if err != nil {
		ctx.Logger.Error("validate basic failed", err)
		return false, err
	}

	// validate fee
	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		ctx.Logger.Error("validate fee failed", err)
		return false, err
	}

	if redeem.ETHTxn == nil {
		ctx.Logger.Error("eth txn is nil")
		return false, action.ErrMissingData
	}

	return true, nil

}

func (ethRedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return processCommon(ctx, tx)

}

func (ethRedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return processCommon(ctx, tx)
	// Create ethereum tracker

}

func processCommon(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	redeem := &Redeem{}
	err := redeem.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("")
		return false, action.Response{Log: action.ErrUnserializable.Error()}
	}

	req, err := ethereum.ParseRedeem(redeem.ETHTxn)
	if err != nil {
		return false, action.Response{Log: action.ErrInvalidExtTx.Error()}
	}

	c, ok := ctx.Currencies.GetCurrencyByName("ETH")
	if !ok {
		return false, action.Response{Log: "ETH not registered"}
	}

	coin := c.NewCoinFromAmount(*balance.NewAmountFromBigInt(req.Amount))
	err = ctx.Balances.MinusFromAddress(redeem.Owner, coin)
	if err != nil {
		return false, action.Response{Log: action.ErrNotEnoughFund.Error()}
	}

	validators, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}

	tracker := trackerlib.NewTracker(
		trackerlib.ProcessTypeRedeem,
		redeem.Owner,
		redeem.ETHTxn,
		ethcommon.BytesToHash(redeem.ETHTxn),
		validators,
	)

	tracker.State = trackerlib.New
	tracker.ProcessOwner = redeem.Owner
	tracker.SignedETHTx = redeem.ETHTxn
	tracker.To = redeem.To

	// Save eth Tracker
	err = ctx.ETHTrackers.Set(tracker)
	return true, action.Response{
		Data:      nil,
		Log:       "",
		Info:      "Transaction received ,Redeem in progress",
		GasWanted: 0,
		GasUsed:   0,
		Tags:      nil,
	}
}

func (ethRedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
