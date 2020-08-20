package bidding

import (
	"fmt"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

type BidConvStore struct {
	state *storage.State
	szlr  serialize.Serializer

	prefix []byte //Current Store Prefix

	prefixActive         []byte
	prefixSucceed        []byte
	prefixRejected       []byte
	prefixCancelled      []byte
	prefixExpired        []byte
	prefixExpiredFailed  []byte
}

func (bcs *BidConvStore) Set(bid *BidConv) error {
	prefixed := append(bcs.prefix, bid.BidConvId...)
	data, err := bcs.szlr.Serialize(bid)
	if err != nil {
		return errors.Wrap(err, errorSerialization)
	}

	err = bcs.state.Set(prefixed, data)

	return errors.Wrap(err, errorSettingRecord)
}

func (bcs *BidConvStore) Get(bidId BidConvId) (*BidConv, error) {
	bid := &BidConv{}
	prefixed := append(bcs.prefix, bidId...)
	data, err := bcs.state.Get(prefixed)
	if err != nil {
		return nil, errors.Wrap(err, errorGettingRecord)
	}
	err = bcs.szlr.Deserialize(data, bid)
	if err != nil {
		return nil, errors.Wrap(err, errorDeSerialization)
	}

	return bid, nil
}

func (bcs *BidConvStore) Exists(key BidConvId) bool {
	active := append(bcs.prefixActive, key...)
	succeed := append(bcs.prefixSucceed, key...)
	rejected := append(bcs.prefixRejected, key...)
	cancelled := append(bcs.prefixCancelled, key...)
	expired := append(bcs.prefixExpired, key...)
	expiredFailed := append(bcs.prefixExpiredFailed, key...)
	return bcs.state.Exists(active) || bcs.state.Exists(succeed) || bcs.state.Exists(rejected) || bcs.state.Exists(cancelled) || bcs.state.Exists(expired) || bcs.state.Exists(expiredFailed)
}

func (bcs *BidConvStore) Delete(key BidConvId) (bool, error) {
	prefixed := append(bcs.prefix, key...)
	res, err := bcs.state.Delete(prefixed)
	if err != nil {
		return false, errors.Wrap(err, errorDeletingRecord)
	}
	return res, err
}

func (bcs *BidConvStore) GetIterable() storage.Iterable {
	return bcs.state.GetIterable()
}

func (bcs *BidConvStore) Iterate(fn func(id BidConvId, bid *BidConv) bool) (stopped bool) {
	return bcs.state.IterateRange(
		bcs.prefix,
		storage.Rangefix(string(bcs.prefix)),
		true,
		func(key, value []byte) bool {
			id := BidConvId(key)
			bid := &BidConv{}

			err := bcs.szlr.Deserialize(value, bid)
			if err != nil {
				return true
			}
			return fn(id, bid)
		},
	)
}

func (bcs *BidConvStore) GetState() *storage.State {
	return bcs.state
}

func (bcs *BidConvStore) WithState(state *storage.State) *BidConvStore {
	bcs.state = state
	return bcs
}

func (bcs *BidConvStore) WithPrefixType(prefixType BidConvState) *BidConvStore {
	switch prefixType {
	case BidStateActive:
		bcs.prefix = bcs.prefixActive
	case BidStateSucceed:
		bcs.prefix = bcs.prefixSucceed
	case BidStateRejected:
		bcs.prefix = bcs.prefixRejected
	case BidStateCancelled:
		bcs.prefix = bcs.prefixCancelled
	case BidStateExpired:
		bcs.prefix = bcs.prefixExpired
	//case BidStateExpireFailed:
	//	bcs.prefix = bcs.prefixExpiredFailed

	}
	return bcs
}

func (bcs *BidConvStore) QueryAllStores(key BidConvId) (*BidConv, BidConvState, error) {
	prefix := bcs.prefix
	defer func() { bcs.prefix = prefix }()

	bid, err := bcs.WithPrefixType(BidStateActive).Get(key)
	if err == nil {
		return bid, BidStateActive, nil
	}
	bid, err = bcs.WithPrefixType(BidStateSucceed).Get(key)
	if err == nil {
		return bid, BidStateSucceed, nil
	}
	bid, err = bcs.WithPrefixType(BidStateRejected).Get(key)
	if err == nil {
		return bid, BidStateRejected, nil
	}
	bid, err = bcs.WithPrefixType(BidStateCancelled).Get(key)
	if err == nil {
		return bid, BidStateCancelled, nil
	}
	bid, err = bcs.WithPrefixType(BidStateExpired).Get(key)
	if err == nil {
		return bid, BidStateExpired, nil
	}
	//bid, err = bcs.WithPrefixType(BidStateExpireFailed).Get(key)
	//if err == nil {
	//	return bid, BidStateExpireFailed, nil
	//}
	return nil, BidStateInvalid, errors.Wrap(err, errorGettingRecord)
}

func (bcs *BidConvStore) FilterBidConvs(bidState BidConvState, owner keys.Address, asset BidAsset, assetType BidAssetType, bidder keys.Address) []BidConv {
	prefix := bcs.prefix
	defer func() { bcs.prefix = prefix }()

	bidConvs := make([]BidConv, 0)
	bcs.WithPrefixType(bidState).Iterate(func(id BidConvId, bidConv *BidConv) bool {
		fmt.Println("owner", owner.String())
		fmt.Println("bidConv.AssetOwner", bidConv.AssetOwner.String())
		if len(owner) != 0 && !bidConv.AssetOwner.Equal(owner) {
			return false
		}
		fmt.Println("bidder", bidder.String())
		fmt.Println("bidConv.Bidder", bidConv.Bidder.String())
		if len(bidder) != 0 && !bidConv.Bidder.Equal(bidder) {
			return false
		}
		fmt.Println("assetType", assetType)
		fmt.Println("bidConv.AssetType", bidConv.AssetType)
		if bidConv.AssetType != assetType {
			return false
		}
		fmt.Println("asset", asset.ToString())
		fmt.Println("bidConv.Asset", bidConv.Asset.ToString())
		if !cmp.Equal(asset, bidConv.Asset) {
			return false
		}

		bidConvs = append(bidConvs, *bidConv)
		return false
	})
	return bidConvs
}

func (bcs *BidConvStore) GetIdForBidConv(bidState BidConvState, owner keys.Address, asset BidAsset, assetType BidAssetType, bidder keys.Address) (BidConvId, error) {
	prefix := bcs.prefix
	defer func() { bcs.prefix = prefix }()

	bidConvs := bcs.FilterBidConvs(bidState, owner, asset, assetType, bidder)
	if len(bidConvs) > 1 || len(bidConvs) == 0 {
		return "", errors.New("errTooManyActiveBidConvs")
	}
	return bidConvs[0].BidConvId, nil
}

func NewBidConvStore(prefixActive string, prefixSucceed string, prefixCancelled string, prefixExpired string, prefixExpiredFailed string, state *storage.State) *BidConvStore {
	return &BidConvStore{
		state:                state,
		szlr:                 serialize.GetSerializer(serialize.LOCAL),
		prefix:               []byte(prefixActive),
		prefixActive:         []byte(prefixActive),
		prefixSucceed:        []byte(prefixSucceed),
		prefixCancelled:      []byte(prefixCancelled),
		prefixExpired:        []byte(prefixExpired),
		prefixExpiredFailed:  []byte(prefixExpiredFailed),
	}
}
