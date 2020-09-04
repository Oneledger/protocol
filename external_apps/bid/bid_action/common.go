package bid_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/external_apps/bid/bid_data"
)

func GetBidMasterStore(ctx *action.Context) (*bid_data.BidMasterStore, error) {
	store, err := ctx.ExtStores.Get("bidMaster")
	if err != nil {
		return nil, bid_data.ErrGettingBidMasterStore.Wrap(err)
	}
	bidMasterStore, ok := store.(*bid_data.BidMasterStore)
	if ok == false {
		return nil, bid_data.ErrAssertingBidMasterStore
	}

	return bidMasterStore, nil
}

func DeactivateOffer(deal bool, bidder action.Address, ctx *action.Context, activeOffer *bid_data.BidOffer, bidMasterStore *bid_data.BidMasterStore) error {
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if activeOffer.OfferType == bid_data.TypeOffer {
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
		// change offer status to inactive
		activeOffer.OfferStatus = bid_data.BidOfferInactive
		err := bidMasterStore.BidOffer.SetOffer(*activeOffer)
		if err != nil {
			return bid_data.ErrUpdateOffer
		}
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
