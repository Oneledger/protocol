package bid_action

import (
	"encoding/json"
	"time"

	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CounterOffer{}

type CounterOffer struct {
	BidConvId  bid_data.BidConvId `json:"bidConvId"`
	AssetOwner keys.Address       `json:"assetOwner"`
	Amount     action.Amount      `json:"amount"`
}

var _ action.Tx = &CounterOfferTx{}

type CounterOfferTx struct {
}

func (c CounterOfferTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	counterOffer := CounterOffer{}
	err := counterOffer.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), counterOffer.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	// the currency should be OLT
	currency, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if currency.Name != counterOffer.Amount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, counterOffer.Amount.String())
	}

	//Check if bid ID is valid
	if counterOffer.BidConvId.Err() != nil {
		return false, bid_data.ErrInvalidBidConvId
	}

	//Check if owner address is valid oneLedger address
	err = counterOffer.AssetOwner.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (c CounterOfferTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CounterOffer Transaction for CheckTx", tx)
	return runCounterOffer(ctx, tx)
}

func (c CounterOfferTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CounterOffer Transaction for DeliverTx", tx)
	return runCounterOffer(ctx, tx)
}

func (c CounterOfferTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCounterOffer(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	counterOffer := CounterOffer{}
	err := counterOffer.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, counterOffer.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidMasterStore, counterOffer.Tags(), err)
	}

	if !bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Exists(counterOffer.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrBidConvNotFound, counterOffer.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Get(counterOffer.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidConv, counterOffer.Tags(), err)
	}

	//2. check owner's identity
	if !counterOffer.AssetOwner.Equal(bidConv.AssetOwner) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongAssetOwner, counterOffer.Tags(), err)
	}

	//3. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrExpiredBid, counterOffer.Tags(), err)
	}

	//4. check asset availability
	available, err := IsAssetAvailable(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner)
	if err != nil || available == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, counterOffer.Tags(), err)
	}

	//5. get active bid offer
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOffer(counterOffer.BidConvId, bid_data.TypeBidOffer)
	if err != nil || activeOffer == nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveBidOffer, counterOffer.Tags(), err)
	}
	//6. amount needs to be large than active bid offer from bidder
	offerCoin := counterOffer.Amount.ToCoin(ctx.Currencies)
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if offerCoin.LessThanEqualCoin(activeOfferCoin) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAmountLessThanActiveOffer, counterOffer.Tags(), err)
	}

	//7. unlock bidder's previous amount and deactivate the bidder's offer
	// this way we only lock amount from a bid offer from bidder
	// if the active offer is a counter offer from owner, no amount is locked from the bidder

	err = DeactivateOffer(false, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, counterOffer.Tags(), err)
	}

	//8. add new counter offer to offer store
	createCounterOffer := bid_data.NewBidOffer(
		counterOffer.BidConvId,
		bid_data.TypeCounterOffer,
		ctx.Header.Time.UTC().Unix(),
		counterOffer.Amount,
		bid_data.CounterOfferAmount,
	)

	err = bidMasterStore.BidOffer.SetActiveOffer(*createCounterOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAddingCounterOffer, counterOffer.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, counterOffer.Tags(), "create_counter_offer_success")
}

func (c CounterOffer) Signers() []action.Address {
	return []action.Address{c.AssetOwner}
}

func (c CounterOffer) Type() action.Type {
	return BID_CONTER_OFFER
}

func (c CounterOffer) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.bidConvId"),
		Value: []byte(c.BidConvId),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(c.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.assetOwner"),
		Value: c.AssetOwner.Bytes(),
	}
	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (c CounterOffer) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CounterOffer) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}
