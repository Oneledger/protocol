package bidding

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/Oneledger/protocol/data/keys"
	"strconv"
)

type BidConv struct {
	BidConvId   BidConvId    `json:"bidId"`
	AssetOwner  keys.Address `json:"assetOwner"`
	Asset       BidAsset     `json:"asset"`
	AssetType   BidAssetType `json:"assetType"`
	Bidder      keys.Address `json:"bidder"`
	DeadlineUTC int64        `json:"deadlineUtc"`
}

func generateBidConvID(key string, blockHeight int64) BidConvId {
	uniqueKey := key + strconv.FormatInt(blockHeight, 10)
	hashHandler := sha256.New()
	_, err := hashHandler.Write([]byte(uniqueKey))
	if err != nil {
		return EmptyStr
	}
	return BidConvId(hex.EncodeToString(hashHandler.Sum(nil)))
}

func NewBidConv(owner keys.Address, asset BidAsset, assetType BidAssetType, bidder keys.Address, deadline int64, height int64) *BidConv {
	return &BidConv{
		BidConvId:   generateBidConvID(owner.String()+asset.ToString()+bidder.String(), height),
		AssetOwner:  owner,
		Asset:       asset,
		AssetType:   assetType,
		Bidder:      bidder,
		DeadlineUTC: deadline,
	}
}
