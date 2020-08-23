package client

import (
	"github.com/Oneledger/protocol/data/bidding"
	"github.com/Oneledger/protocol/data/keys"
)

type ListBidConvRequest struct {
	BidConvId bidding.BidConvId `json:"bidConvId"`
}

type BidConvStat struct {
	BidConv bidding.BidConv    `json:"bidConv"`
	Offers  []bidding.BidOffer `json:"bidOffers"`
}

type ActiveOfferStat struct {
	ActiveOffer bidding.BidOffer `json:"activeOffer"`
	BidConv     bidding.BidConv  `json:"bidConv"`
}

type ListBidConvsReply struct {
	BidConvStats []BidConvStat `json:"bidConvStats"`
	Height       int64         `json:"height"`
}

type ListBidConvsRequest struct {
	State     bidding.BidConvState `json:"state"`
	Owner     keys.Address         `json:"owner"`
	AssetName string               `json:"assetName"`
	AssetType bidding.BidAssetType `json:"assetType"`
	Bidder    keys.Address         `json:"bidder"`
}

type ListActiveOffersRequest struct {
	Owner     keys.Address         `json:"owner"`
	AssetName string               `json:"assetName"`
	AssetType bidding.BidAssetType `json:"assetType"`
	Bidder    keys.Address         `json:"bidder"`
	OfferType bidding.BidOfferType `json:"offerType"`
}

type ListActiveOffersReply struct {
	ActiveOffers []ActiveOfferStat `json:"activeOffers"`
	Height       int64             `json:"height"`
}
