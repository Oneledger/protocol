package bidding

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bidding"
)

func GetBidMasterStore(ctx *action.Context) (*bidding.BidMasterStore, error) {
	store, err := ctx.ExtStores.Get("bidMaster")
	if err != nil {
		return nil, bidding.ErrGettingBidMasterStore
	}
	bidMasterStore, ok := store.(*bidding.BidMasterStore)
	if ok == false {
		return nil, bidding.ErrAssertingBidMasterStore
	}

	return bidMasterStore, nil
}
