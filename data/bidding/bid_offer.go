package bidding

import (
	"github.com/Oneledger/protocol/data/balance"
)

type BidOffer struct {
	BidId       BidConvId 		`json:"bidId"`
	OfferStatus BidOfferStatus  `json:"offerStatus"`
	OfferType   BidOfferType    `json:"offerType"`
	OfferTime   int64			`json:"offerTime"`
	AcceptTime  int64    		`json:"acceptTime"`
	RejectTime  int64			`json:"rejectTime"`
	Amount      balance.Amount  `json:"amount"`
}

func NewBidOffer(bidId BidConvId, offerType BidOfferType, offerTime int64, amount balance.Amount) *BidOffer {
	return &BidOffer{BidId: bidId, OfferStatus: BidOfferActive, OfferType: offerType, OfferTime: offerTime, AcceptTime: 0, RejectTime: 0, Amount: amount}
}


