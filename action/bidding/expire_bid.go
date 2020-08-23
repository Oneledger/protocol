package bidding

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/bidding"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
)

var _ action.Msg = &ExpireBid{}

type ExpireBid struct {
	BidConvId        bidding.BidConvId `json:"bidConvId"`
	ValidatorAddress action.Address    `json:"validatorAddress"`
}

func (e ExpireBid) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	expireBid := ExpireBid{}
	err := expireBid.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), expireBid.Signers(), signedTx.Signatures)
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
	if expireBid.BidConvId.Err() != nil {
		return false, bidding.ErrInvalidBidConvId
	}

	//Check if validator address is valid oneLedger address

	err = expireBid.ValidatorAddress.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (e ExpireBid) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runExpireBid(ctx, tx)
}

func (e ExpireBid) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runExpireBid(ctx, tx)
}

func (e ExpireBid) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runExpireBid(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	expireBid := ExpireBid{}
	err := expireBid.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, expireBid.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, expireBid.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Exists(expireBid.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, expireBid.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Get(expireBid.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidConv, expireBid.Tags(), err)
	}

	//2. get the active offer(bid offer or counter offer)
	activeOffers := bidMasterStore.BidOffer.GetOffers(expireBid.BidConvId, bidding.BidOfferActive, bidding.TypeInvalid)
	// in this case there must be an offer
	if len(activeOffers) == 0 {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveOffer, expireBid.Tags(), err)
	} else if len(activeOffers) > 1 {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrTooManyActiveOffers, expireBid.Tags(), err)
	}
	activeOffer := activeOffers[0]

	//3. unlock amount and set offer to inactive(if active offer is bid offer from bidder)
	err = DeactivateOffer(false, bidConv.Bidder, ctx, &activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrDeactivateOffer, expireBid.Tags(), err)
	}

	err = CloseBidConv(bidConv, bidMasterStore, bidding.BidStateExpired)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrCloseBidConv, expireBid.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, expireBid.Tags(), "expire_bid_success")
}

func (e ExpireBid) Signers() []action.Address {
	return []action.Address{e.ValidatorAddress}
}

func (e ExpireBid) Type() action.Type {
	return action.BID_EXPIRE
}

func (e ExpireBid) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.bidConvId"),
		Value: []byte(e.BidConvId),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(e.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validatorAddress"),
		Value: e.ValidatorAddress.Bytes(),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (e ExpireBid) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ExpireBid) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, e)
}
