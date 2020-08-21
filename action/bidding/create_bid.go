package bidding

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/bidding"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CreateBid{}

type CreateBid struct {
	BidConvId      	bidding.BidConvId		`json:"bidConvId"`
	AssetOwner 		keys.Address 			`json:"assetOwner"`
	Asset      		bidding.BidAsset 		`json:"asset"`
	AssetType 		bidding.BidAssetType 	`json:"assetType"`
	Bidder     		keys.Address 			`json:"bidder"`
	Amount     		action.Amount           `json:"amount"`
	OfferTime		int64					`json:"offerTime"`
	Deadline	 	int64                   `json:"deadline"`
}

func (c CreateBid) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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
	if currency.Name != createBid.Amount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, createBid.Amount.String())
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

func (c CreateBid) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for CheckTx", tx)
	return runCreateBid(ctx, tx)
}

func (c CreateBid) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing CreateProposal Transaction for DeliverTx", tx)
	return runCreateBid(ctx, tx)
}

func (c CreateBid) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCreateBid(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createBid := CreateBid{}
	err := createBid.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, createBid.Tags(), err)
	}

	//1. check asset availability
	assetOk, err := createBid.Asset.ValidateAsset(ctx, createBid.AssetOwner)
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

	//3. verify bidConvId exists in ACTIVE store
	bidMasterStore, err := GetBidMasterStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidMasterStore, createBid.Tags(), err)
	}

	if !bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Exists(createBid.BidConvId) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrBidConvNotFound, createBid.Tags(), err)
	}

	bidConv, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Get(createBid.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingBidConv, createBid.Tags(), err)
	}

	//3. check bidder's identity
	if !createBid.Bidder.Equal(bidConv.Bidder) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrWrongBidder, createBid.Tags(), err)
	}

	//4. check expiry
	deadLine := time.Unix(bidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrExpiredBid, createBid.Tags(), err)
	}

	//5. there should be no active bid offer from bidder
	activeOffer, err := bidMasterStore.BidOffer.GetActiveOfferForBidConvId(createBid.BidConvId)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrGettingActiveOffer, createBid.Tags(), err)
	}
	offerCoin := createBid.Amount.ToCoin(ctx.Currencies)

	if activeOffer != nil {
		if activeOffer.OfferType == bidding.TypeOffer {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrActiveBidOfferExists, createBid.Tags(), err)
		}
		//5. amount needs to be less than active counter offer from owner
		activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
		if activeOfferCoin.LessThanEqualCoin(offerCoin) {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAmountMoreThanActiveCounterOffer, createBid.Tags(), err)
		}
		//6. set active counter offer to inactive
		err = DeactivateOffer(true, bidConv.Bidder, ctx, activeOffer, bidMasterStore)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrDeactivateOffer, createBid.Tags(), err)
		}
	}

	//7. lock amount
	err = ctx.Balances.MinusFromAddress(createBid.Bidder.Bytes(), offerCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrLockAmount, createBid.Tags(), err)
	}

	//8. add new offer to offer store
	createBidOffer := bidding.NewBidOffer(
		createBid.BidConvId,
		bidding.TypeOffer,
		createBid.OfferTime,
		createBid.Amount,
		bidding.BidAmountLocked,
	)

	err = bidMasterStore.BidOffer.SetOffer(*createBidOffer)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, bidding.ErrAddingOffer, createBid.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, createBid.Tags(), "create_bid_success")
}

func (c CreateBid) Signers() []action.Address {
	return []action.Address{c.Bidder}
}

func (c CreateBid) Type() action.Type {
	return action.BID_CREATE
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
		Value: []byte(c.Asset.ToString()),
	}
	tag4 := kv.Pair{
		Key: []byte("tx.assetType"),
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

func (c *CreateBid) createBidConv(ctx *action.Context) error {
	createBidConv := bidding.NewBidConv(
		c.AssetOwner,
		c.Asset,
		c.AssetType,
		c.Bidder,
		c.Deadline,
	)
	//Validate bid deadline
	deadLine := time.Unix(createBidConv.DeadlineUTC, 0)

	if deadLine.Before(ctx.Header.Time.UTC()) {
		return bidding.ErrInvalidDeadline
	}

	//Check if any bid conversation with same asset, owner, bidder already exists in active store
	store, err := ctx.ExtStores.Get("bidMaster")
	if err != nil {
		return bidding.ErrGettingBidMasterStore.Wrap(err)
	}
	bidMasterStore := store.(*bidding.BidMasterStore)
	filteredBidConvs := bidMasterStore.BidConv.FilterBidConvs(bidding.BidStateActive, createBidConv.AssetOwner, createBidConv.Asset, createBidConv.AssetType, createBidConv.Bidder)
	if len(filteredBidConvs) != 0 {
		return bidding.ErrActiveBidConvExists
	}
	//Add bid conversation to DB
	activeBidConvs := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive)
	err = activeBidConvs.Set(createBidConv)
	if err != nil {
		return bidding.ErrAddingBidConvToActiveStore.Wrap(err)
	}

	return nil
}
