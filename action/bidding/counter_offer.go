package bidding

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/bidding"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CounterOffer{}

type CounterOffer struct {
	BidConvId      	bidding.BidConvId		`json:"bidConvId"`
	AssetOwner 		keys.Address 			`json:"assetOwner"`
	Amount     		action.Amount           `json:"amount"`
	OfferTime		int64					`json:"offerTime"`
}

func (c CounterOffer) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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
	feeOpt, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, governance.ErrGetFeeOptions
	}
	err = action.ValidateFee(feeOpt, signedTx.Fee)
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
		return false, bidding.ErrInvalidBidConvId
	}

	//Check if owner address is valid oneLedger address
	err = counterOffer.AssetOwner.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (c CounterOffer) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runCounterOffer(ctx, tx)
}

func (c CounterOffer) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runCounterOffer(ctx, tx)
}

func (c CounterOffer) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
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
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, counterOffer.Tags(), err)
	}

	if !bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Exists(counterOffer.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, counterOffer.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Get(counterOffer.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidConv, counterOffer.Tags(), err)
	}

	//3. check owner's identity
	if !counterOffer.AssetOwner.Equal(bidConv.AssetOwner) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrWrongAssetOwner, counterOffer.Tags(), err)
	}

	//2. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrExpiredBid, counterOffer.Tags(), err)
	}

	//2. there should be no active counter offer from owner
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOfferForBidConvId(counterOffer.BidConvId)
	// in the counter offer case, there must be an active offer
	if err != nil || activeOffer == nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveOffer, counterOffer.Tags(), err)
	}
	if activeOffer.OfferType == bidding.TypeCounterOffer {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrActiveCounterOfferExists, counterOffer.Tags(), err)
	}

	//3. amount needs to be large than active bid offer from bidder
	offerCoin := counterOffer.Amount.ToCoin(ctx.Currencies)
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if offerCoin.LessThanEqualCoin(activeOfferCoin) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAmountLessThanActiveOffer, counterOffer.Tags(), err)
	}

	//4. unlock bidder's previous amount and deactivate the bidder's offer
	// this way we only lock amount from a bid offer from bidder
	// if the active offer is a counter offer from owner, no amount is locked from the bidder

	err = DeactivateOffer(true, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrDeactivateOffer, counterOffer.Tags(), err)
	}

	//5. add new counter offer to offer store
	createCounterOffer := bidding.NewBidOffer(
		counterOffer.BidConvId,
		bidding.TypeCounterOffer,
		counterOffer.OfferTime,
		counterOffer.Amount,
		bidding.CounterOfferAmount,
	)

	err = bidMasterStore.BidOffer.SetOffer(*createCounterOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAddingCounterOffer, counterOffer.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, counterOffer.Tags(), "create_counter_offer_success")
}

func (c CounterOffer) Signers() []action.Address {
	return []action.Address{c.AssetOwner}
}

func (c CounterOffer) Type() action.Type {
	return action.BID_CREATE
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
