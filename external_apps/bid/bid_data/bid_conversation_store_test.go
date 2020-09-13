package bid_data

import (
	"fmt"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"strconv"
	"testing"
)

const (
	numPrivateKeys = 10
	numBids        = 10
	testDeadline   = 1596740920
	height         = 111
)

var (
	addrList []keys.Address
	bidConvs []*BidConv

	bidConvStore *BidConvStore

	assetNames []string
)

func init() {
	fmt.Println("####### TESTING BID CONV STORE #######")

	//Generate key pairs
	for i := 0; i < numPrivateKeys; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		addrList = append(addrList, h.Address())
	}

	//Create new bid conversations
	for i := 0; i < numBids; i++ {
		j := i / 2                  //owner address list ranges from 0 - 4
		k := numPrivateKeys - 1 - j //bidder address list ranges from 9 - 5

		owner := addrList[j]
		assetNames = append(assetNames, "test"+strconv.Itoa(i)+".ol")
		bidder := addrList[k]

		bidConvs = append(bidConvs, NewBidConv(owner, assetNames[i], BidAssetOns,
			bidder, testDeadline, height))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create bid conversation store
	bidConvStore = NewBidConvStore("p_active", "p_succeed", "p_cancelled", "p_expired", "p_rejected", cs)
}

func TestBidConvStore_Set(t *testing.T) {
	err := bidConvStore.Set(bidConvs[0])
	assert.Equal(t, nil, err)

	bidConv, err := bidConvStore.Get(bidConvs[0].BidConvId)
	assert.Equal(t, nil, err)

	assert.Equal(t, bidConv.BidConvId, bidConvs[0].BidConvId)
}

func TestBidConvStore_Exists(t *testing.T) {
	exists := bidConvStore.Exists(bidConvs[0].BidConvId)
	assert.Equal(t, true, exists)

	exists = bidConvStore.Exists(bidConvs[1].BidConvId)
	assert.Equal(t, false, exists)
}

func TestBidConvStore_Delete(t *testing.T) {
	_, err := bidConvStore.Get(bidConvs[0].BidConvId)
	assert.Equal(t, nil, err)

	res, err := bidConvStore.Delete(bidConvs[0].BidConvId)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = bidConvStore.Get(bidConvs[0].BidConvId)
	assert.NotEqual(t, nil, err)
}

func TestBidConvStore_Iterate(t *testing.T) {
	for _, val := range bidConvs {
		_ = bidConvStore.Set(val)
	}
	bidConvStore.state.Commit()

	bidConvCount := 0
	bidConvStore.Iterate(func(id BidConvId, bidConv *BidConv) bool {
		bidConvCount++
		return false
	})

	assert.Equal(t, numBids, bidConvCount)
}

func TestBidConvStore_FilterBidConvs(t *testing.T) {
	bidConvs := bidConvStore.FilterBidConvs(BidStateActive, addrList[0], assetNames[0], BidAssetOns,
		addrList[9])
	assert.Equal(t, 1, len(bidConvs))
}
