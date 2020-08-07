package bidding

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
	numBids   = 10
	testDeadline = 1596740920
)

var (
	addrList    []keys.Address
	bidConvs     []*BidConv

	bidConvStore    *BidConvStore

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
		j := i / 2      //owner address list ranges from 0 - 4
		k := numPrivateKeys - 1 - j //bidder address list ranges from 9 - 5

		owner := addrList[j]
		asset := NewDomainAsset("test" + strconv.Itoa(i) + ".ol")
		bidder := addrList[k]

		bidConvs = append(bidConvs, NewBidConv(owner, asset, BidAssetOns,
			bidder, testDeadline))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create bid conversation store
	bidConvStore = NewBidConvStore("p_active", "p_succeed", "p_cancelled", "p_expired", "p_expiredFailed", cs)
}

func TestBidConvStore_Set(t *testing.T) {
	err := bidConvStore.Set(bidConvs[0])
	assert.Equal(t, nil, err)

	bidConv, err := bidConvStore.Get(bidConvs[0].BidId)
	assert.Equal(t, nil, err)

	assert.Equal(t, bidConv.BidId, bidConvs[0].BidId)
}

func TestBidConvStore_Exists(t *testing.T) {
	exists := bidConvStore.Exists(bidConvs[0].BidId)
	assert.Equal(t, true, exists)

	exists = bidConvStore.Exists(bidConvs[1].BidId)
	assert.Equal(t, false, exists)
}

func TestBidConvStore_Delete(t *testing.T) {
	_, err := bidConvStore.Get(bidConvs[0].BidId)
	assert.Equal(t, nil, err)

	res, err := bidConvStore.Delete(bidConvs[0].BidId)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = bidConvStore.Get(bidConvs[0].BidId)
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

//func TestBidConvStore_FilterBidConvs(t *testing.T) {
//	onsBidAddr0 := bidConvStore.FilterBidConvs(BidStateActive, addrList[0])
//	onsBidAddr1 := bidConvStore.FilterBidConvs(BidStateSucceed, addrList[1])
//	onsBidAddr2 := bidConvStore.FilterBidConvs(BidStateCancelled, addrList[2])
//	assert.Equal(t, 2, len(onsBidAddr0))
//	assert.Equal(t, 2, len(onsBidAddr1))
//	assert.Equal(t, 2, len(onsBidAddr2))
//}
