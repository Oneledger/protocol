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

var _ action.Msg = &CancelBid{}

type CancelBid struct {
	BidConvId bid_data.BidConvId `json:"bidConvId"`
	Bidder    keys.Address       `json:"bidder"`
}

var _ action.Tx = &CancelBidTx{}

type CancelBidTx struct {
}

func (c CancelBidTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if bid ID is valid
	if cancelBid.BidConvId.Err() != nil {
		return false, bid_data.ErrInvalidBidConvId
	}

	//Check if bidder address is valid oneLedger address

	err = cancelBid.Bidder.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}

	return true, nil
}

func (c CancelBidTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CancelBid Transaction for CheckTx", tx)
	return runCancelBid(ctx, tx)
}

func (c CancelBidTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CancelBid Transaction for DeliverTx", tx)
	return runCancelBid(ctx, tx)
}

func (c CancelBidTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCancelBid(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	// bidder can cancel a bid as long as bid is in ACTIVE store

	cancelBid := CancelBid{}
	err := cancelBid.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, cancelBid.Tags(), err)
	}

	//1. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidMasterStore, cancelBid.Tags(), err)
	}
	if !bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Exists(cancelBid.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrBidConvNotFound, cancelBid.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Get(cancelBid.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidConv, cancelBid.Tags(), err)
	}

	//2. check bidder's identity
	if !cancelBid.Bidder.Equal(bidConv.Bidder) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongBidder, cancelBid.Tags(), err)
	}

	//3. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrExpiredBid, cancelBid.Tags(), err)
	}

	//4. get the active offer/counter offer
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOffer(cancelBid.BidConvId, bid_data.TypeInvalid)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveOffer, cancelBid.Tags(), err)
	}
	//5. unlock amount if needed and deactivate it
	err = DeactivateOffer(false, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, cancelBid.Tags(), err)
	}

	//6. close bid and put to CANCELLED store
	err = CloseBidConv(bidConv, bidMasterStore, bid_data.BidStateCancelled)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrCloseBidConv, cancelBid.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, cancelBid.Tags(), "cancel_bid_success")
}

func (c CancelBid) Signers() []action.Address {
	return []action.Address{c.Bidder}
}

func (c CancelBid) Type() action.Type {
	return BID_CANCEL
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
