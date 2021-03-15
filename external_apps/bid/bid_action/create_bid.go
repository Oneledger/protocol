package bid_action

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CreateBid{}

type CreateBid struct {
	BidConvId  bid_data.BidConvId    `json:"bidConvId"`
	AssetOwner keys.Address          `json:"assetOwner"`
	AssetName  string                `json:"assetName"`
	AssetType  bid_data.BidAssetType `json:"assetType"`
	Bidder     keys.Address          `json:"bidder"`
	Amount     action.Amount         `json:"amount"`
	Deadline   int64                 `json:"deadline"`
}

var _ action.Tx = &CreateBidTx{}

type CreateBidTx struct {
}

func (c CreateBidTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createBid := CreateBid{}
	err := createBid.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), createBid.Signers(), signedTx.Signatures)
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
	if currency.Name != createBid.Amount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, createBid.Amount.String())
	}

	//Check if bid ID is valid(if provided)
	if len(createBid.BidConvId) > 0 && createBid.BidConvId.Err() != nil {
		return false, bid_data.ErrInvalidBidConvId
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

func (c CreateBidTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateBid Transaction for CheckTx", tx)
	return runCreateBid(ctx, tx)
}

func (c CreateBidTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateBid Transaction for DeliverTx", tx)
	return runCreateBid(ctx, tx)
}

func (c CreateBidTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCreateBid(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	// if this is to create bid conversation, everything except bidConvId is needed
	// if this is just to add an offer from bidder, only needs bidConvId, bidder(to sign), amount
	createBid := CreateBid{}
	err := createBid.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, createBid.Tags(), err)
	}

	//1. check if this is to create a bid conversation or just add an offer
	bidConvId := createBid.BidConvId
	if len(createBid.BidConvId) == 0 {
		// check asset availability
		available, err := IsAssetAvailable(ctx, createBid.AssetName, createBid.AssetType, createBid.AssetOwner)
		if err != nil || available == false {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, createBid.Tags(), err)
		}
		bidConvId, err = createBid.createBidConv(ctx)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrFailedCreateBidConv, createBid.Tags(), err)
		}
	}

	//2. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidMasterStore, createBid.Tags(), err)
	}

	if !bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Exists(bidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrBidConvNotFound, createBid.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Get(bidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingBidConv, createBid.Tags(), err)
	}

	//3. check asset availability if this is just to add an offer
	if len(createBid.BidConvId) != 0 {
		available, err := IsAssetAvailable(ctx, bidConv.AssetName, bidConv.AssetType, bidConv.AssetOwner)
		if err != nil || available == false {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrInvalidAsset, createBid.Tags(), err)
		}
	}
	//4. check bidder's identity
	if !createBid.Bidder.Equal(bidConv.Bidder) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrWrongBidder, createBid.Tags(), err)
	}

	//5. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrExpiredBid, createBid.Tags(), err)
	}

	offerCoin := createBid.Amount.ToCoin(ctx.Currencies)

	//6. get the active counter offer
	activeCounterOffer, err := bidMasterStore.BidOffer.GetActiveOffer(bidConvId, bid_data.TypeCounterOffer)
	// in this case there can be no counter offer if this is the beginning of bid conversation
	if err != nil || (len(createBid.BidConvId) != 0 && activeCounterOffer == nil) {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrGettingActiveCounterOffer, createBid.Tags(), err)
	}
	if activeCounterOffer != nil {
		//7. amount needs to be less than active counter offer from owner
		activeOfferCoin := activeCounterOffer.Amount.ToCoin(ctx.Currencies)
		if activeOfferCoin.LessThanEqualCoin(offerCoin) {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAmountMoreThanActiveCounterOffer, createBid.Tags(), err)
		}
		//8. set active counter offer to inactive
		err = DeactivateOffer(false, bidConv.Bidder, ctx, activeCounterOffer, bidMasterStore)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrDeactivateOffer, createBid.Tags(), err)
		}
	}
	//9. lock amount
	err = ctx.Balances.MinusFromAddress(createBid.Bidder, offerCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrLockAmount, createBid.Tags(), err)
	}

	//10. add new offer to offer store
	createBidOffer := bid_data.NewBidOffer(
		bidConvId,
		bid_data.TypeBidOffer,
		ctx.Header.Time.UTC().Unix(),
		createBid.Amount,
		bid_data.BidAmountLocked,
	)

	err = bidMasterStore.BidOffer.SetActiveOffer(*createBidOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bid_data.ErrAddingOffer, createBid.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, createBid.Tags(), "create_bid_success")
}

func (c CreateBid) Signers() []action.Address {
	return []action.Address{c.Bidder}
}

func (c CreateBid) Type() action.Type {
	return BID_CREATE
}

func (c CreateBid) Tags() kv.Pairs {
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
		Value: []byte(c.AssetName),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.assetType"),
		Value: []byte(strconv.Itoa(int(c.AssetType))),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4)
	return tags
}

func (c CreateBid) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CreateBid) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}

func (c *CreateBid) createBidConv(ctx *action.Context) (bid_data.BidConvId, error) {
	createBidConv := bid_data.NewBidConv(
		c.AssetOwner,
		c.AssetName,
		c.AssetType,
		c.Bidder,
		c.Deadline,
		ctx.Header.Height,
	)
	//Validate bid deadline
	deadLine := time.Unix(createBidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return "", bid_data.ErrInvalidDeadline
	}

	//Check if any bid conversation with same asset, owner, bidder already exists in active store

	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return "", bid_data.ErrGettingBidMasterStore.Wrap(err)
	}
	filteredBidConvs := bidMasterStore.BidConv.FilterBidConvs(bid_data.BidStateActive, createBidConv.AssetOwner, createBidConv.AssetName, createBidConv.AssetType, createBidConv.Bidder)
	if len(filteredBidConvs) != 0 {
		return "", bid_data.ErrActiveBidConvExists
	}
	//Add bid conversation to DB
	activeBidConvs := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive)
	err = activeBidConvs.Set(createBidConv)
	if err != nil {
		return "", bid_data.ErrAddingBidConvToActiveStore.Wrap(err)
	}
	//pass the generated id
	return createBidConv.BidConvId, nil
}
