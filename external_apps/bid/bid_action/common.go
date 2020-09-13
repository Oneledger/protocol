package bid_action

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
)

func GetBidMasterStore(ctx *action.Context) (*bid_data.BidMasterStore, error) {
	store, err := ctx.ExtStores.Get("extBidMaster")
	if err != nil {
		return nil, bid_data.ErrGettingBidMasterStore.Wrap(err)
	}
	bidMasterStore, ok := store.(*bid_data.BidMasterStore)
	if ok == false {
		return nil, bid_data.ErrAssertingBidMasterStore
	}

	return bidMasterStore, nil
}

func IsAssetAvailable(ctx *action.Context, assetName string, assetType bid_data.BidAssetType, assetOwner keys.Address) (bool, error) {
	bidAssetTemplate := BidAssetMap[assetType]
	bidAsset := bidAssetTemplate.NewAssetWithName(assetName)
	//fmt.Printf("bidAsset: %p %v\n", bidAsset, bidAsset)
	assetOk, err := bidAsset.ValidateAsset(ctx, assetOwner)
	return assetOk, err
}

func ExchangeAsset(ctx *action.Context, assetName string, assetType bid_data.BidAssetType, assetOwner keys.Address, bidder keys.Address) (bool, error) {
	bidAssetTemplate := BidAssetMap[assetType]
	bidAsset := bidAssetTemplate.NewAssetWithName(assetName)
	//fmt.Printf("bidAsset: %p %v\n", bidAsset, bidAsset)
	exchangeOk, err := bidAsset.ExchangeAsset(ctx, bidder, assetOwner)
	return exchangeOk, err
}

func DeactivateOffer(deal bool, bidder action.Address, ctx *action.Context, activeOffer *bid_data.BidOffer, bidMasterStore *bid_data.BidMasterStore) error {
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if activeOffer.OfferType == bid_data.TypeBidOffer {
		// unlock the amount if no deal
		if deal == false {
			err := ctx.Balances.AddToAddress(bidder.Bytes(), activeOfferCoin)
			if err != nil {
				return bid_data.ErrUnlockAmount.Wrap(err)
			}
			// change amount status to unlocked
			activeOffer.AmountStatus = bid_data.BidAmountUnlocked
			// add reject time
			activeOffer.RejectTime = ctx.Header.Time.UTC().Unix()
		} else {
			// change amount status to transferred if there is a deal
			activeOffer.AmountStatus = bid_data.BidAmountTransferred
			// add accept time
			activeOffer.AcceptTime = ctx.Header.Time.UTC().Unix()
		}
	} else if activeOffer.OfferType == bid_data.TypeCounterOffer {
		if deal == false {
			// add reject time
			activeOffer.RejectTime = ctx.Header.Time.UTC().Unix()
		} else {
			// add accept time
			activeOffer.AcceptTime = ctx.Header.Time.UTC().Unix()
		}
	} else {
		return bid_data.ErrInvalidOfferType
	}
	// add updated offer as inactive offer
	err := bidMasterStore.BidOffer.SetInActiveOffer(*activeOffer)
	if err != nil {
		return bid_data.ErrSetOffer
	}
	//delete active offer, here only use bidConvId, same object is ok to be passed in
	err = bidMasterStore.BidOffer.DeleteActiveOffer(*activeOffer)
	if err != nil {
		fmt.Println("err in delete offer: ", err)
		return bid_data.ErrFailedToDeleteActiveOffer
	}
	return nil
}

func CloseBidConv(bidConv *bid_data.BidConv, bidMasterStore *bid_data.BidMasterStore, targetState bid_data.BidConvState) error {

	//add bid conversation to target store
	err := bidMasterStore.BidConv.WithPrefixType(targetState).Set(bidConv)
	if err != nil {
		return bid_data.ErrAddingBidConvToTargetStore.Wrap(err)
	}

	//delete it from ACTIVE store
	ok, err := bidMasterStore.BidConv.WithPrefixType(bid_data.BidStateActive).Delete(bidConv.BidConvId)
	if err != nil || !ok {
		return bid_data.ErrDeletingBidConvFromActiveStore.Wrap(err)
	}
	return nil
}
