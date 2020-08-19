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

var _ action.Msg = &CancelBid{}

type CancelBid struct {
	BidConvId      	bidding.BidConvId		`json:"bidConvId"`
	Bidder     		keys.Address 			`json:"bidder"`
}

func (c CancelBid) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	cancelBid := CancelBid{}
	err := cancelBid.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), cancelBid.Signers(), signedTx.Signatures)
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
	if cancelBid.BidConvId.Err() != nil {
		return false, bidding.ErrInvalidBidConvId
	}

	//Check if bidder address is valid oneLedger address

	err = cancelBid.Bidder.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (c CancelBid) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runCancelBid(ctx, tx)
}

func (c CancelBid) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runCancelBid(ctx, tx)
}

func (c CancelBid) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCancelBid(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	cancelBid := CancelBid{}
	err := cancelBid.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, cancelBid.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, cancelBid.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Exists(cancelBid.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, cancelBid.Tags(), err)
	}

	//2. check if there is active counter offer from owner
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOfferForBidConvId(cancelBid.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveOffers, cancelBid.Tags(), err)
	}
	if activeOffer.OfferType == bidding.TypeCounterOffer {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrActiveCounterOfferExists, cancelBid.Tags(), err)
	}

	//3. unlock amount
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.AddToAddress(cancelBid.Bidder.Bytes(), activeOfferCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrUnlockAmount, cancelBid.Tags(), err)
	}

	//4. change amount status to unlocked and add it back to bid offer store
	activeOffer.AmountStatus = bidding.BidAmountUnlocked
	err = bidMasterStore.BidOffer.SetOffer(*activeOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrUpdateBidOffer, cancelBid.Tags(), err)
	}

	//5. add bid conversation to CANCELLED store

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Get(activeOffer.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidConv, cancelBid.Tags(), err)
	}
	err = bidMasterStore.BidConv.WithPrefixType(bidding.BidStateCancelled).Set(bidConv)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAddingBidConvToCancelledStore, cancelBid.Tags(), err)
	}

	// 5. delete it from ACTIVE store
	ok, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Delete(activeOffer.BidConvId)
	if err != nil || !ok {
		return false, action.Response{Log: bidding.ErrDeletingBidConvFromActiveStore.Marshal()}
	}

	return helpers.LogAndReturnTrue(ctx.Logger, cancelBid.Tags(), "cancel_bid_success")
}

func (c CancelBid) Signers() []action.Address {
	return []action.Address{c.Bidder}
}

func (c CancelBid) Type() action.Type {
	return action.BID_CREATE
}

func (c CancelBid) Tags() kv.Pairs {
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
		Value: c.Bidder.Bytes(),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (c CancelBid) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CancelBid) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}
