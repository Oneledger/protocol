package bidding

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

func (bos *BidOfferStore) Set(key storage.StoreKey, offer BidOffer) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(offer)
	if err != nil {
		return errors.Wrap(err, errorSerialization)
	}
	prefixed := append(bos.prefix, key...)
	err = bos.State.Set(prefixed, dat)
	return errors.Wrap(err, errorSettingRecord)
}

func (bos *BidOfferStore) get(key storage.StoreKey) (offer *BidOffer, err error) {
	prefixed := append(bos.prefix, storage.StoreKey(key)...)
	dat, err := bos.State.Get(prefixed)
	//fmt.Println("dat :", dat, "err", err)
	if err != nil {
		return nil, errors.Wrap(err, errorGettingRecord)
	}
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, offer)
	if err != nil {
		err = errors.Wrap(err, errorDeSerialization)
	}
	return
}

func (bos *BidOfferStore) delete(key storage.StoreKey) (bool, error) {
	prefixed := append(bos.prefix, key...)
	res, err := bos.State.Delete(prefixed)
	if err != nil {
		return false, errors.Wrap(err, errorDeletingRecord)
	}
	return res, err
}

func (bos *BidOfferStore) iterate(fn func(bidConvId BidConvId, offStatus BidOfferStatus, offerType BidOfferType, offerTime int64, offer BidOffer) bool) bool {
	return bos.State.IterateRange(
		bos.prefix,
		storage.Rangefix(string(bos.prefix)),
		true,
		func(key, value []byte) bool {
			offer := &BidOffer{}
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, offer)
			if err != nil {
				return true
			}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			fmt.Println("key: ", key)
			fmt.Println("arr: ", arr)
			fmt.Println("arr[1]: ", arr[1])
			// key example: bidOffer_bidConvId_offerStatus_offerType_offerTime
			bidConvId := arr[1]
			fmt.Println("bidConvId: ", bidConvId)
			offerStatus, err := strconv.Atoi(arr[2])
			if err != nil {
				fmt.Println("Error Parsing Offer Status", err)
				return true
			}
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
			return fn(BidConvId(bidConvId), BidOfferStatus(offerStatus), BidOfferType(offerType), int64(offerTime), *offer)
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

func (bos *BidOfferStore) GetOffers(bId BidConvId, oStatus BidOfferStatus, oType BidOfferType) []BidOffer {
	var bidOffers []BidOffer

	bos.iterate(func(bidConvId BidConvId, offerStatus BidOfferStatus, offerType BidOfferType, offerTime int64, offer BidOffer) bool {
		if len(bId) != 0 && bId != bidConvId {
			return false
		}
		if oStatus != BidOfferInvalid && offerStatus != oStatus {
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

func (bos *BidOfferStore) GetActiveOfferForBidConvId(id BidConvId, offerType BidOfferType) (*BidOffer, error) {
	activeOffers := bos.GetOffers(id, BidOfferActive, offerType)
	if len(activeOffers) == 0 {
		return nil, errors.New("errNoActiveOffer")
	} else if len(activeOffers) > 1 {
		return nil, errors.New("errTooManyActiveOffers")
	}
	return &activeOffers[0], nil
}

func assembleBidOfferKey(bidConvId BidConvId, offerStatus BidOfferStatus, offerType BidOfferType, offerTime int64) storage.StoreKey {
	key := storage.StoreKey(string(bidConvId) + storage.DB_PREFIX + strconv.Itoa(int(offerStatus)) + storage.DB_PREFIX + strconv.Itoa(int(offerType)) + storage.DB_PREFIX + strconv.FormatInt(offerTime, 10))
	return key
}

func (bos *BidOfferStore) SetOffer(offer BidOffer) error {
	key := assembleBidOfferKey(offer.BidConvId, offer.OfferStatus, offer.OfferType, offer.OfferTime)
	err := bos.Set(key, offer)
	if err != nil {
		return err
	}
	return nil
}
