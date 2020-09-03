package bid_data

import (
	"github.com/Oneledger/protocol/action"
)

type BidOffer struct {
	BidConvId    BidConvId            `json:"bidConvId"`
	OfferStatus  BidOfferStatus       `json:"offerStatus"`
	OfferType    BidOfferType         `json:"offerType"`
	OfferTime    int64                `json:"offerTime"`
	AcceptTime   int64                `json:"acceptTime"`
	RejectTime   int64                `json:"rejectTime"`
	Amount       action.Amount        `json:"amount"`
	AmountStatus BidOfferAmountStatus `json:"amountStatus"`
}

func NewBidOffer(bidConvId BidConvId, offerType BidOfferType, offerTime int64, amount action.Amount, amountStatus BidOfferAmountStatus) *BidOffer {
	return &BidOffer{BidConvId: bidConvId, OfferStatus: BidOfferActive, OfferType: offerType, OfferTime: offerTime, AcceptTime: 0, RejectTime: 0, Amount: amount, AmountStatus: amountStatus}
}