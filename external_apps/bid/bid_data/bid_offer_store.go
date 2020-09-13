package bid_data

import (
	"fmt"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type BidOfferStore struct {
	State  *storage.State
	prefix []byte
}

func (bos *BidOfferStore) set(key storage.StoreKey, offer BidOffer) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(offer)
	if err != nil {
		return errors.Wrap(err, errorSerialization)
	}
	prefixed := append(bos.prefix, key...)
	err = bos.State.Set(prefixed, dat)
	return errors.Wrap(err, errorSettingRecord)
}

func (bos *BidOfferStore) get(key storage.StoreKey) (offer *BidOffer, err error) {
	prefixed := append(bos.prefix, key...)
	dat, err := bos.State.Get(prefixed)
	if err != nil {
		return nil, ErrGettingRecord.Wrap(err)
	}
	if len(dat) == 0 {
		return
	}
	offer = &BidOffer{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, offer)
	if err != nil {
		err = ErrFailedInDeserialization.Wrap(err)
	}
	return
}

func (bos *BidOfferStore) delete(key storage.StoreKey) (bool, error) {
	prefixed := append(bos.prefix, key...)
	res, err := bos.State.Delete(prefixed)
	if err != nil {
		return false, ErrDeletingRecord.Wrap(err)
	}
	return res, err
}

func assembleInactiveOfferPrefix(prefix []byte) []byte {
	prefixString := string(prefix) + InactiveOfferPrefix
	return storage.Prefix(prefixString)
}

func (bos *BidOfferStore) iterate(fn func(bidConvId BidConvId, offerType BidOfferType, offerTime int64, offer BidOffer) bool) bool {
	return bos.State.IterateRange(
		assembleInactiveOfferPrefix(bos.prefix),
		storage.Rangefix(string(assembleInactiveOfferPrefix(bos.prefix))),
		true,
		func(key, value []byte) bool {
			offer := &BidOffer{}
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, offer)
			if err != nil {
				return true
			}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			// key example: bidOffer_INACTIVE_bidConvId_offerType_offerTime
			bidConvId := arr[2]
			offerType, err := strconv.Atoi(arr[3])
			if err != nil {
				fmt.Println("Error Parsing Offer Type", err)
				return true
			}
			offerTime, err := strconv.ParseInt(arr[len(arr)-1], 10, 64)
			if err != nil {
				fmt.Println("Error Parsing Offer Time", err)
				return true
			}
			return fn(BidConvId(bidConvId), BidOfferType(offerType), int64(offerTime), *offer)
		},
	)
}

func (bos *BidOfferStore) WithState(state *storage.State) *BidOfferStore {
	bos.State = state
	return bos
}

func NewBidOfferStore(prefix string, state *storage.State) *BidOfferStore {
	return &BidOfferStore{
		State:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (bos *BidOfferStore) GetInActiveOffers(bId BidConvId, oType BidOfferType) []BidOffer {
	var bidOffers []BidOffer
	bos.iterate(func(bidConvId BidConvId, offerType BidOfferType, offerTime int64, offer BidOffer) bool {
		if len(bId) != 0 && bId != bidConvId {
			return false
		}
		if oType != TypeInvalid && offerType != oType {
			return false
		}
		bidOffers = append(bidOffers, offer)
		return false
	})
	return bidOffers
}

func (bos *BidOfferStore) GetActiveOffer(bId BidConvId, oType BidOfferType) (*BidOffer, error) {
	key := assembleActiveOfferKey(bId)
	offer, err := bos.get(key)
	if err != nil {
		return nil, err
	}
	if offer == nil {
		return nil, nil
	}
	if oType != TypeInvalid && offer.OfferType != oType {
		return nil, ErrWrongTypeForActiveOffer
	}
	return offer, nil
}

func assembleOfferKey(bidConvId BidConvId, offerType BidOfferType, offerTime int64) storage.StoreKey {
	key := storage.StoreKey(InactiveOfferPrefix + storage.DB_PREFIX + string(bidConvId) + storage.DB_PREFIX + strconv.Itoa(int(offerType)) + storage.DB_PREFIX + strconv.FormatInt(offerTime, 10))
	return key
}

func assembleActiveOfferKey(bidConvId BidConvId) storage.StoreKey {
	key := storage.StoreKey(ActiveOfferPrefix + storage.DB_PREFIX + string(bidConvId))
	return key
}

func (bos *BidOfferStore) SetInActiveOffer(offer BidOffer) error {
	key := assembleOfferKey(offer.BidConvId, offer.OfferType, offer.OfferTime)
	err := bos.set(key, offer)
	if err != nil {
		return err
	}
	return nil
}

func (bos *BidOfferStore) SetActiveOffer(offer BidOffer) error {
	key := assembleActiveOfferKey(offer.BidConvId)
	err := bos.set(key, offer)
	if err != nil {
		return err
	}
	return nil
}

func (bos *BidOfferStore) DeleteActiveOffer(offer BidOffer) error {
	key := assembleActiveOfferKey(offer.BidConvId)
	ok, err := bos.delete(key)
	if err != nil {
		return err
	}
	if ok == false {
		return errors.New("failed to delete active offer")
	}
	return nil
}
