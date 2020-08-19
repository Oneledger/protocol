package bidding

import (
	"errors"
	"github.com/Oneledger/protocol/action"
)

type (
	BidConvId       string
	BidConvState    int
	BidAssetType	int
	BidConvStatus   bool
	BidOfferStatus  bool
	BidOfferType    int
	BidOfferAmountStatus int
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
	ToString() string
	ValidateAsset(ctx *action.Context) (bool, error)
}
