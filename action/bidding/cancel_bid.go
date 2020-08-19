package bidding

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/bidding"
	"strconv"

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

	//Check if bid ID is valid(if provided)
	if len(createBid.BidConvId) > 0 && createBid.BidConvId.Err() != nil {
		return false, bidding.ErrInvalidBidConvId
	}

	//Check if bidder and owner address is valid oneLedger address(if bid id is not provided)
	if len(createBid.BidConvId) == 0 {
		err = createBid.Bidder.Err()
		if err != nil {
			return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
		}

		err = createBid.AssetOwner.Err()
		if err != nil {
			return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
		}
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
	createBid := CreateBid{}
	err := createBid.Unmarshal(tx.Data)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(createBid.Tags(), "create_bid_offer_failed_deserialize"),
			Log:    action.ErrWrongTxType.Wrap(err).Marshal(),
		}
		return false, result
	}

	//1. check asset availability
	assetOk, err := createBid.Asset.ValidateAsset(ctx)
	if err != nil || assetOk == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrInvalidAsset, createBid.Tags(), err)
	}

	//2. check if this is to create a bid conversation or just add an offer
	if len(createBid.BidConvId) == 0 {
		err := createBid.createBidConv(ctx)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrFailedCreateBidConv, createBid.Tags(), err)
		}
	}

	//3. verify bidConvId exists
	store, err := ctx.ExtStores.Get("bidMaster")
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, createBid.Tags(), err)
	}
	bidMasterStore, ok := store.(*bidding.BidMasterStore)
	if ok == false {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAssertingBidMasterStore, createBid.Tags(), err)
	}
	if !bidMasterStore.BidConv.Exists(createBid.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, createBid.Tags(), err)
	}

	//4. add the offer to offer store
	createBidOffer := bidding.NewBidOffer(
		createBid.BidConvId,
		createBid.OfferTime,
		createBid.Amount,
	)

	err = bidMasterStore.BidOffer.AddOffer(*createBidOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAddingOffer, createBid.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, createBid.Tags(), "create_bid_offer_success")
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
		Value: c.AssetOwner.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.asset"),
		Value: []byte(c.Asset.ToString()),
	}
	tag4 := kv.Pair{
		Key: []byte("tx.assetType"),
		Value: []byte(strconv.Itoa(int(c.AssetType))),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4)
	return tags
}

func (c CancelBid) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CancelBid) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}
