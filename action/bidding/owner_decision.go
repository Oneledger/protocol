package bidding

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/bidding"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &OwnerDecision{}

type OwnerDecision struct {
	BidConvId      	bidding.BidConvId		`json:"bidConvId"`
	Owner     		keys.Address 			`json:"owner"`
	Decision		bidding.BidDecision		`json:"decision"`
}

func (o OwnerDecision) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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
	feeOpt, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, governance.ErrGetFeeOptions
	}
	err = action.ValidateFee(feeOpt, signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if bid ID is valid
	if ownerDecision.BidConvId.Err() != nil {
		return false, bidding.ErrInvalidBidConvId
	}

	//Check if owner address is valid oneLedger address
	err = ownerDecision.Owner.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (o OwnerDecision) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runOwnerDecision(ctx, tx)
}

func (o OwnerDecision) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runOwnerDecision(ctx, tx)
}

func (o OwnerDecision) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
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
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, ownerDecision.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Exists(ownerDecision.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, ownerDecision.Tags(), err)
	}
	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Get(ownerDecision.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidConv, ownerDecision.Tags(), err)
	}

	//2. check owner's identity
	if !ownerDecision.Owner.Equal(bidConv.AssetOwner) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrWrongAssetOwner, ownerDecision.Tags(), err)
	}

	//3. check if there is active offer from bidder
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOfferForBidConvId(ownerDecision.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveOffers, ownerDecision.Tags(), err)
	}
	if activeOffer.OfferType == bidding.TypeCounterOffer {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveBidOffer, ownerDecision.Tags(), err)
	}

	//4. if reject
	if ownerDecision.Decision == bidding.RejectBid {
		// deactivate offer and unlock amount depends on active offer type
		err = DeactivateOffer(true, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrDeactivateOffer, ownerDecision.Tags(), err)
		}
		// close bid conversation
		err = CloseBidConv(bidConv, activeOffer, bidMasterStore, bidding.BidStateRejected)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrCloseBidConv, ownerDecision.Tags(), err)
		}

		return helpers.LogAndReturnTrue(ctx.Logger, ownerDecision.Tags(), "owner_reject_bid_success")

	}

	//5. add the amount to owner, in this case offer amount is already being locked from bidder
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.AddToAddress(bidConv.AssetOwner.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAdddingAmountToOwner, ownerDecision.Tags(), err)
	}

	//6. change offer status to inactive and add it back to bid offer store
	err = DeactivateOffer(false, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrDeactivateOffer, ownerDecision.Tags(), err)
	}

	//7. close the bid conversation
	err = CloseBidConv(bidConv, activeOffer, bidMasterStore, bidding.BidStateSucceed)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrCloseBidConv, ownerDecision.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, ownerDecision.Tags(), "owner_accept_bid_success")
}

func (o OwnerDecision) Signers() []action.Address {
	return []action.Address{o.Owner}
}

func (o OwnerDecision) Type() action.Type {
	return action.BID_CREATE
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
