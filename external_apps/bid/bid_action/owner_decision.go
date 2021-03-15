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

var _ action.Msg = &OwnerDecision{}

type OwnerDecision struct {
	BidConvId bid_data.BidConvId   `json:"bidConvId"`
	Owner     keys.Address         `json:"owner"`
	Decision  bid_data.BidDecision `json:"decision"`
}

var _ action.Tx = &OwnerDecisionTx{}

type OwnerDecisionTx struct {
}

func (o OwnerDecisionTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	ownerDecision := OwnerDecision{}
	err := ownerDecision.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), ownerDecision.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if bid ID is valid
	if ownerDecision.BidConvId.Err() != nil {
		return false, bid_data.ErrInvalidBidConvId
	}

	//Check if owner address is valid oneLedger address
	err = ownerDecision.Owner.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (o OwnerDecisionTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing OwnerDecision Transaction for CheckTx", tx)
	return runOwnerDecision(ctx, tx)
}

func (o OwnerDecisionTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing OwnerDecision Transaction for DeliverTx", tx)
	return runOwnerDecision(ctx, tx)
}

func (o OwnerDecisionTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runOwnerDecision(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ownerDecision := OwnerDecision{}
	err := ownerDecision.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, ownerDecision.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidMasterStore, ownerDecision.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Exists(ownerDecision.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrBidConvNotFound, ownerDecision.Tags(), err)
	}
	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Get(ownerDecision.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidConv, ownerDecision.Tags(), err)
	}

	//2. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrExpiredBid, ownerDecision.Tags(), err)
	}

	//3. check asset availability
	available, err := IsAssetAvailable(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner)
	if err != nil || available == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, ownerDecision.Tags(), err)
	}

	//4. check owner's identity
	if !ownerDecision.Owner.Equal(bidConv.AssetOwner) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongAssetOwner, ownerDecision.Tags(), err)
	}

	//5. get active bid offer
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOffer(ownerDecision.BidConvId, bid_data.TypeBidOffer)
	// in this case, there must be an existing active offer
	if err != nil || activeOffer == nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveOffer, ownerDecision.Tags(), err)
	}

	//6. if reject
	if ownerDecision.Decision != bid_data.RejectBid && ownerDecision.Decision != bid_data.AcceptBid {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidOwnerDecision, ownerDecision.Tags(), err)
	} else if ownerDecision.Decision == bid_data.RejectBid {
		// deactivate offer and unlock amount
		err = DeactivateOffer(false, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, ownerDecision.Tags(), err)
		}
		// close bid conversation
		err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateRejected)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, ownerDecision.Tags(), err)
		}

		return helpers.LogAndReturnTrue(ctx.Logger, ownerDecision.Tags(), "owner_reject_bid_success")

	}

	//7. add the amount to owner, in this case offer amount is already being locked from bidder
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.AddToAddress(bidConv.AssetOwner.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAdddingAmountToOwner, ownerDecision.Tags(), err)
	}

	//8. change offer status to inactive and add it back to bid offer store
	err = DeactivateOffer(true, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, ownerDecision.Tags(), err)
	}

	//9. close the bid conversation
	err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateSucceed)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, ownerDecision.Tags(), err)
	}

	//10. exchange asset
	ok, err := ExchangeAsset(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner, bidConv.Bidder)
	if err != nil || ok == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrFailedToExchangeAsset, ownerDecision.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, ownerDecision.Tags(), "owner_accept_bid_success")
}

func (o OwnerDecision) Signers() []action.Address {
	return []action.Address{o.Owner}
}

func (o OwnerDecision) Type() action.Type {
	return BID_OWNER_DECISION
}

func (o OwnerDecision) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.bidConvId"),
		Value: []byte(o.BidConvId),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(o.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.assetOwner"),
		Value: o.Owner.Bytes(),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (o OwnerDecision) Marshal() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OwnerDecision) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, o)
}
