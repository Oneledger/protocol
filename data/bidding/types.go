package bidding

import (
	"errors"
)

type (
	BidConvId       string
	BidConvState    int
	BidAsset        interface{}
	BidAssetType	int
	BidConvStatus   bool
	BidOfferStatus  bool
	BidOfferType    bool
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
