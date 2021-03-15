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

var _ action.Msg = &BidderDecision{}

type BidderDecision struct {
	BidConvId bid_data.BidConvId   `json:"bidConvId"`
	Bidder    keys.Address         `json:"bidder"`
	Decision  bid_data.BidDecision `json:"decision"`
}

var _ action.Tx = &BidderDecisionTx{}

type BidderDecisionTx struct {
}

func (b BidderDecisionTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	bidderDecision := BidderDecision{}
	err := bidderDecision.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), bidderDecision.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if bid ID is valid
	if bidderDecision.BidConvId.Err() != nil {
		return false, bid_data.ErrInvalidBidConvId
	}

	//Check if bidder address is valid oneLedger address
	err = bidderDecision.Bidder.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (b BidderDecisionTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing BidderDecision Transaction for CheckTx", tx)
	return runBidderDecision(ctx, tx)
}

func (b BidderDecisionTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing BidderDecision Transaction for DeliverTx", tx)
	return runBidderDecision(ctx, tx)
}

func (b BidderDecisionTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runBidderDecision(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	bidderDecision := BidderDecision{}
	err := bidderDecision.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, bidderDecision.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidMasterStore, bidderDecision.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Exists(bidderDecision.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrBidConvNotFound, bidderDecision.Tags(), err)
	}
	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Get(bidderDecision.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidConv, bidderDecision.Tags(), err)
	}

	//2. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrExpiredBid, bidderDecision.Tags(), err)
	}

	//3. check asset availability
	available, err := IsAssetAvailable(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner)
	if err != nil || available == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, bidderDecision.Tags(), err)
	}

	//4. check bidder's identity
	if !bidderDecision.Bidder.Equal(bidConv.Bidder) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongBidder, bidderDecision.Tags(), err)
	}

	//5. get the active counter offer
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOffer(bidderDecision.BidConvId, bid_data.TypeCounterOffer)
	if err != nil || activeOffer == nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveCounterOffer, bidderDecision.Tags(), err)
	}

	//6. if reject
	if bidderDecision.Decision != bid_data.RejectBid && bidderDecision.Decision != bid_data.AcceptBid {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidBidderDecision, bidderDecision.Tags(), err)
	}
	if bidderDecision.Decision == bid_data.RejectBid {
		// deactivate offer
		err = DeactivateOffer(false, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, bidderDecision.Tags(), err)
		}
		// close bid conversation
		err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateRejected)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, bidderDecision.Tags(), err)
		}

		return helpers.LogAndReturnTrue(ctx.Logger, bidderDecision.Tags(), "bidder_reject_bid_success")

	}

	//7. deduct the amount from bidder, in this case no amount is currently being locked
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(bidderDecision.Bidder.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeductingAmountFromBidder, bidderDecision.Tags(), err)
	}

	//8. add the amount to owner
	err = ctx.Balances.AddToAddress(bidConv.AssetOwner.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAdddingAmountToOwner, bidderDecision.Tags(), err)
	}

	//9. deactivate offer
	err = DeactivateOffer(true, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, bidderDecision.Tags(), err)
	}

	//10. close the bid conversation
	err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateSucceed)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, bidderDecision.Tags(), err)
	}

	//11. exchange asset
	ok, err := ExchangeAsset(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner, bidConv.Bidder)
	if err != nil || ok == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrFailedToExchangeAsset, bidderDecision.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, bidderDecision.Tags(), "bidder_accept_bid_success")
}

func (b BidderDecision) Signers() []action.Address {
	return []action.Address{b.Bidder}
}

func (b BidderDecision) Type() action.Type {
	return BID_BIDDER_DECISION
}

func (b BidderDecision) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.bidConvId"),
		Value: []byte(b.BidConvId),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(b.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.assetOwner"),
		Value: b.Bidder.Bytes(),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (b BidderDecision) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BidderDecision) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, b)
}
