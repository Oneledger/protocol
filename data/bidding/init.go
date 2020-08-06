package bidding

import (
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("ons_bid")
}

const (
	//Bid States
	BidStateInvalid       BidConvState = 0xEE
	BidStateActive        BidConvState = 0x01
	BidStateSucceed       BidConvState = 0x02
	BidStateCancelled     BidConvState = 0x03
	BidStateExpired       BidConvState = 0x04
	BidStateExpiredFailed BidConvState = 0x05

	//Error Codes
	errorSerialization   = "321"
	errorDeSerialization = "322"
	errorSettingRecord   = "323"
	errorGettingRecord   = "324"
	errorDeletingRecord  = "325"

	EmptyStr = ""

	//Bid Conversation Status
	BidConvOpen		BidConvStatus = true
	BidConvClosed   BidConvStatus = false

	//Bid Offer Status
	BidOfferActive  	BidOfferStatus = true
	BidOfferRejected    BidOfferStatus = false

	//Bid Offer Type
	TypeOffer         BidOfferType = true
	TypeCounterOffer  BidOfferType= false

	//Bid Asset Type
	BidAssetOns BidAssetType = 0x11

	//Bid Id length based on hash algorithm
	SHA256LENGTH int = 0x40


)

type BidMasterStore struct {
	BidConv            *BidConvStore
	BidOffer    	   *BidOfferStore
}

func (bm *BidMasterStore) WithState(state *storage.State) *BidMasterStore {
	bm.BidConv.WithState(state)
	bm.BidOffer.WithState(state)
	return bm
}

func NewBidMasterStore(bc *BidConvStore, bo *BidOfferStore) *BidMasterStore {
	return &BidMasterStore{
		BidConv:     	bc,
		BidOffer: 		bo,
	}
}