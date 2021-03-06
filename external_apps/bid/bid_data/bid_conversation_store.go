package bid_data

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/google/go-cmp/cmp"
)

type BidConvStore struct {
	state *storage.State
	szlr  serialize.Serializer

	prefix []byte //Current Store Prefix

	prefixActive        []byte
	prefixSucceed       []byte
	prefixRejected      []byte
	prefixCancelled     []byte
	prefixExpired       []byte
	prefixExpiredFailed []byte
}

func (bcs *BidConvStore) Set(bid *BidConv) error {
	prefixed := append(bcs.prefix, bid.BidConvId...)
	data, err := bcs.szlr.Serialize(bid)
	if err != nil {
		return ErrFailedInSerialization.Wrap(err)
	}

	err = bcs.state.Set(prefixed, data)
	if err != nil {
		return ErrSettingRecord.Wrap(err)
	}

	return nil
}

func (bcs *BidConvStore) Get(bidId BidConvId) (*BidConv, error) {
	bid := &BidConv{}
	prefixed := append(bcs.prefix, bidId...)
	data, err := bcs.state.Get(prefixed)
	if err != nil {
		return nil, ErrGettingRecord.Wrap(err)
	}
	err = bcs.szlr.Deserialize(data, bid)
	if err != nil {
		return nil, ErrFailedInDeserialization.Wrap(err)
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
		return false, ErrDeletingRecord.Wrap(err)
	}
	return res, err
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
	return nil, BidStateInvalid, ErrGettingRecord.Wrap(err)
}

func (bcs *BidConvStore) FilterBidConvs(bidState BidConvState, owner keys.Address, assetName string, assetType BidAssetType, bidder keys.Address) []BidConv {
	prefix := bcs.prefix
	defer func() { bcs.prefix = prefix }()

	bidConvs := make([]BidConv, 0)
	bcs.WithPrefixType(bidState).Iterate(func(id BidConvId, bidConv *BidConv) bool {
		if len(owner) != 0 && !bidConv.AssetOwner.Equal(owner) {
			return false
		}
		if len(bidder) != 0 && !bidConv.Bidder.Equal(bidder) {
			return false
		}
		if bidConv.AssetType != assetType {
			return false
		}
		if len(assetName) != 0 && !cmp.Equal(assetName, bidConv.AssetName) {
			return false
		}

		bidConvs = append(bidConvs, *bidConv)
		return false
	})
	return bidConvs
}

func NewBidConvStore(prefixActive string, prefixSucceed string, prefixCancelled string, prefixExpired string, prefixRejected string, state *storage.State) *BidConvStore {
	return &BidConvStore{
		state:           state,
		szlr:            serialize.GetSerializer(serialize.LOCAL),
		prefix:          []byte(prefixActive),
		prefixActive:    []byte(prefixActive),
		prefixSucceed:   []byte(prefixSucceed),
		prefixCancelled: []byte(prefixCancelled),
		prefixExpired:   []byte(prefixExpired),
		prefixRejected:  []byte(prefixRejected),
	}
}
