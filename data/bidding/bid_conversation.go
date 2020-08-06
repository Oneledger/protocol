package bidding

import (
	"github.com/Oneledger/protocol/data/keys"
)

type BidConv struct {
	BidId      			BidConvId    	`json:"bidId"`
	AssetOwner 			keys.Address 	`json:"assetOwner"`
	Asset      			BidAsset 		`json:"asset"`
	AssetType 			BidAssetType 	`json:"assetType"`
	Bidder     			keys.Address 	`json:"bidder"`
	Deadline   			int64    		`json:"deadline"`
	Status     			BidConvStatus   `json:"status"`
	BidOffers   	 	[]BidOffer	 	`json:"bidOffers"`
	BidCounterOffers 	[]BidOffer	 	`json:"bidCounterOffers"`
}


func NewBidConv(bidId BidConvId, owner keys.Address, asset BidAsset, assetType BidAssetType, bidder keys.Address, deadline int64) *BidConv {
	return &BidConv{
		BidId:         		bidId,
		AssetOwner:    		owner,
		Asset:         		asset,
		AssetType:			assetType,
		Bidder:        		bidder,
		Deadline: 	   		deadline,
		Status:		   		BidConvOpen,
		BidOffers:        	[]BidOffer{},
		BidCounterOffers: 	[]BidOffer{},
	}
}
