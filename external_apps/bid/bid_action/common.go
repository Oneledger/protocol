package bid_action

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bidding"
)

func GetBidMasterStore(ctx *action.Context) (*bidding.BidMasterStore, error) {
	store, err := ctx.ExtStores.Get("bidMaster")
	if err != nil {
		return nil, bidding.ErrGettingBidMasterStore.Wrap(err)
	}
	bidMasterStore, ok := store.(*bidding.BidMasterStore)
	if ok == false {
		return nil, bidding.ErrAssertingBidMasterStore
	}

	return bidMasterStore, nil
}

func DeactivateOffer(deal bool, bidder action.Address, ctx *action.Context, activeOffer *bidding.BidOffer, bidMasterStore *bidding.BidMasterStore) error {
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if activeOffer.OfferType == bidding.TypeOffer {
		// unlock the amount if no deal
		if deal == false {
			err := ctx.Balances.AddToAddress(bidder.Bytes(), activeOfferCoin)
			if err != nil {
				return bidding.ErrUnlockAmount.Wrap(err)
			}
			// change amount status to unlocked
			activeOffer.AmountStatus = bidding.BidAmountUnlocked
			// add reject time
			activeOffer.RejectTime = ctx.Header.Time.UTC().Unix()
		} else {
			// change amount status to transferred if there is a deal
			activeOffer.AmountStatus = bidding.BidAmountTransferred
			// add accept time
			activeOffer.AcceptTime = ctx.Header.Time.UTC().Unix()
		}
	} else if activeOffer.OfferType == bidding.TypeCounterOffer {
		// change offer status to inactive
		activeOffer.OfferStatus = bidding.BidOfferInactive
		err := bidMasterStore.BidOffer.SetOffer(*activeOffer)
		if err != nil {
			return bidding.ErrUpdateOffer
		}
		if deal == false {
			// add reject time
			activeOffer.RejectTime = ctx.Header.Time.UTC().Unix()
		} else {
			// add accept time
			activeOffer.AcceptTime = ctx.Header.Time.UTC().Unix()
		}
	} else {
		return bidding.ErrInvalidOfferType
	}
	return nil
}

func CloseBidConv(bidConv *bidding.BidConv, bidMasterStore *bidding.BidMasterStore, targetState bidding.BidConvState) error {

	//add bid conversation to target store
	err := bidMasterStore.BidConv.WithPrefixType(targetState).Set(bidConv)
	if err != nil {
		return bidding.ErrAddingBidConvToTargetStore.Wrap(err)
	}

	//delete it from ACTIVE store
	ok, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Delete(bidConv.BidConvId)
	if err != nil || !ok {
		return bidding.ErrDeletingBidConvFromActiveStore.Wrap(err)
	}
	return nil
}