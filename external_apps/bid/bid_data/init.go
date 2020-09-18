package bid_data

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

func init() {
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

	EmptyStr = ""

	//Bid Offer Type
	TypeBidOffer     BidOfferType = 0x01
	TypeCounterOffer BidOfferType = 0x02
	TypeInvalid      BidOfferType = 0x03

	//Bid Offer Amount Lock Status
	BidAmountLocked      BidOfferAmountStatus = 0x01
	BidAmountUnlocked    BidOfferAmountStatus = 0x02
	CounterOfferAmount   BidOfferAmountStatus = 0x03
	BidAmountTransferred BidOfferAmountStatus = 0x04

	//Bid Decision
	AcceptBid BidDecision = 0x01
	RejectBid BidDecision = 0x02

	//Bid Asset Type
	BidAssetInvalid BidAssetType = 0xEE
	BidAssetOns     BidAssetType = 0x21
	BidAssetExample BidAssetType = 0x22

	//Bid Id length based on hash algorithm
	SHA256LENGTH int = 0x40

	ActiveOfferPrefix   string = "ACTIVE"
	InactiveOfferPrefix string = "INACTIVE"
)

type BidMasterStore struct {
	BidConv  *BidConvStore
	BidOffer *BidOfferStore
}

var _ data.ExtStore = &BidMasterStore{}

func (bm *BidMasterStore) WithState(state *storage.State) data.ExtStore {
	bm.BidConv.WithState(state)
	bm.BidOffer.WithState(state)
	return bm
}

func ConstructBidMasterStore(bc *BidConvStore, bo *BidOfferStore) *BidMasterStore {
	return &BidMasterStore{
		BidConv:  bc,
		BidOffer: bo,
	}
}

func NewBidMasterStore(chainstate *storage.ChainState) *BidMasterStore {
	bidConv := NewBidConvStore("extBidConvActive", "extBidConvSucceed", "extBidConvCancelled", "extBidConvExpired", "extBidConvRejected", storage.NewState(chainstate))
	bidOffer := NewBidOfferStore("extBidOffer", storage.NewState(chainstate))
	return ConstructBidMasterStore(bidConv, bidOffer)
}
