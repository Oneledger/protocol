package bid_data

import (
	"errors"
	"github.com/Oneledger/protocol/action"
)

type (
	BidConvId            string
	BidConvState         int
	BidAssetType         int
	//todo delete this if not used
	BidConvStatus        bool
	BidOfferStatus       int
	BidOfferType         int
	BidOfferAmountStatus int
	BidDecision          int
)

func (id BidConvId) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("bid id is empty")
	case len(id) != SHA256LENGTH:
		return errors.New("bid id length is incorrect")
	}
	return nil
}

type BidAsset interface {
	SetName(name string)
	ToString() string
	ValidateAsset(ctx *action.Context, owner action.Address) (bool, error)
	ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error)
}
