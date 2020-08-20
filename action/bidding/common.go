package bidding

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

func DeactivateOffer(unlock bool, bidder action.Address, ctx *action.Context, activeOffer *bidding.BidOffer, bidMasterStore *bidding.BidMasterStore) error {
	activeOfferCoin := activeOffer.Amount.ToCoin(ctx.Currencies)
	if activeOffer.OfferType == bidding.TypeOffer {
		// unlock the amount if no deal
		if unlock == true {
			err := ctx.Balances.AddToAddress(bidder.Bytes(), activeOfferCoin)
			if err != nil {
				return bidding.ErrUnlockAmount.Wrap(err)
			}
			// change amount status to unlocked
			activeOffer.AmountStatus = bidding.BidAmountUnlocked
		} else {
			// change amount status to transferred if there is a deal
			activeOffer.AmountStatus = bidding.BidAmountTransferred
		}
	}
	// change offer status to inactive
	activeOffer.OfferStatus = bidding.BidOfferInactive
	err := bidMasterStore.BidOffer.SetOffer(*activeOffer)
	if err != nil {
		return bidding.ErrUpdateBidOffer
	}
	return nil
}

func CloseBidConv(bidConv *bidding.BidConv, activeOffer *bidding.BidOffer, bidMasterStore *bidding.BidMasterStore, targetState bidding.BidConvState) error {

	//add bid conversation to target store
	err := bidMasterStore.BidConv.WithPrefixType(targetState).Set(bidConv)
	if err != nil {
		return bidding.ErrAddingBidConvToTargetStore.Wrap(err)
	}

	//delete it from ACTIVE store
	ok, err := bidMasterStore.BidConv.WithPrefixType(bidding.BidStateActive).Delete(activeOffer.BidConvId)
	if err != nil || !ok {
		return bidding.ErrDeletingBidConvFromActiveStore.Wrap(err)
	}
	return nil
}
