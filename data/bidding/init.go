package bidding

import (
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("bidding")
	serialize.RegisterConcrete(new(DomainAsset), "domain_asset")
	serialize.RegisterConcrete(new(ExampleAsset), "example_asset")
}

const (
	//Bid States
	BidStateInvalid   BidConvState = 0xEE
	BidStateActive    BidConvState = 0x01
	BidStateSucceed   BidConvState = 0x02
	BidStateCancelled BidConvState = 0x03
	BidStateExpired   BidConvState = 0x04
	BidStateRejected  BidConvState = 0x05
	//BidStateExpireFailed BidConvState = 0x05

	//Error Codes
	errorSerialization   = "321"
	errorDeSerialization = "322"
	errorSettingRecord   = "323"
	errorGettingRecord   = "324"
	errorDeletingRecord  = "325"

	EmptyStr = ""

	//Bid Offer Status
	BidOfferActive   BidOfferStatus = 0x01
	BidOfferInactive BidOfferStatus = 0x02
	BidOfferInvalid  BidOfferStatus = 0x03

	//Bid Offer Type
	TypeOffer        BidOfferType = 0x01
	TypeCounterOffer BidOfferType = 0x02
	TypeInvalid      BidOfferType = 0x03

	//Bid Offer Amount Lock Status
	BidAmountLocked      BidOfferAmountStatus = 0x01
	BidAmountUnlocked    BidOfferAmountStatus = 0x02
	CounterOfferAmount   BidOfferAmountStatus = 0x03
	BidAmountTransferred BidOfferAmountStatus = 0x04

	//Bid Decision
	AcceptBid BidDecision = true
	RejectBid BidDecision = false

	//Bid Asset Type
	BidAssetInvalid BidAssetType = 0xEE
	BidAssetOns     BidAssetType = 0x21

	//Bid Id length based on hash algorithm
	SHA256LENGTH int = 0x40
)

type BidMasterStore struct {
	BidConv  *BidConvStore
	BidOffer *BidOfferStore
}

func (bm *BidMasterStore) WithState(state *storage.State) *BidMasterStore {
	bm.BidConv.WithState(state)
	bm.BidOffer.WithState(state)
	return bm
}

func NewBidMasterStore(bc *BidConvStore, bo *BidOfferStore) *BidMasterStore {
	return &BidMasterStore{
		BidConv:  bc,
		BidOffer: bo,
	}
}
