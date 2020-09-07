package bid_action

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &BidderDecision{}

type BidderDecision struct {
	BidConvId bid_data.BidConvId   `json:"bidConvId"`
	Bidder    keys.Address        `json:"bidder"`
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
	//todo change all txs in bid, do not read chainstate here
	feeOpt, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, governance.ErrGetFeeOptions
	}
	err = action.ValidateFee(feeOpt, signedTx.Fee)
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
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runBidderDecision(ctx, tx)
}

func (b BidderDecisionTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runBidderDecision(ctx, tx)
}

func (b BidderDecisionTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
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

	//1. check asset availability
	assetOk, err := bidConv.Asset.ValidateAsset(ctx, bidConv.AssetOwner)
	if err != nil || assetOk == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, bidderDecision.Tags(), err)
	}

	//3. check bidder's identity
	if !bidderDecision.Bidder.Equal(bidConv.Bidder) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongBidder, bidderDecision.Tags(), err)
	}

	//4. get the active counter offer
	activeOffers := bidMasterStore.BidOffer.GetOffers(bidderDecision.BidConvId, bid_data.BidOfferActive, bid_data.TypeCounterOffer)
	// in this case there must be a counter offer from owner
	if len(activeOffers) == 0 {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveCounterOffer, bidderDecision.Tags(), err)
	} else if len(activeOffers) > 1 {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrTooManyActiveOffers, bidderDecision.Tags(), err)
	}
	activeOffer := activeOffers[0]

	//5. if reject
	if bidderDecision.Decision == bid_data.RejectBid {
		// deactivate offer and unlock amount depends on active offer type
		err = DeactivateOffer(false, bidConv.Bidder, ctx, &activeOffer, bidMasterStore)
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

	//6. deduct the amount from bidder, in this case no amount is currently being locked
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(bidderDecision.Bidder.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeductingAmountFromBidder, bidderDecision.Tags(), err)
	}

	//7. add the amount to owner
	err = ctx.Balances.AddToAddress(bidConv.AssetOwner.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAdddingAmountToOwner, bidderDecision.Tags(), err)
	}

	//8. change offer status to inactive and add it back to bid offer store
	err = DeactivateOffer(true, bidConv.Bidder, ctx, &activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, bidderDecision.Tags(), err)
	}

	//9. close the bid conversation
	err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateSucceed)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, bidderDecision.Tags(), err)
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
