package bidding

import (
	"fmt"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"testing"
)

var (
	bidOfferStore    *BidOfferStore
	cs       *storage.State
	address  keys.Address
	address2 keys.Address
	ID1      BidConvId
	ID2      BidConvId
	offerTimes []int64
	offerTime1 int64 = 1596701481
	offerTimeCount = 7

)

func init() {
	fmt.Println("####### TESTING BID OFFER STORE #######")
	memDb := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("offer", memDb))
	bidOfferStore = NewBidOfferStore("test", cs)
	generateAddresses()
	generateIDs()
	generateOfferTimes()
}

func generateAddresses() {
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	address2 = h2.Address()
}

func generateOfferTimes() {
	for i := 0; i < offerTimeCount; i++ {
		offerTimes[i] = offerTime1 + int64(i * 20)
	}
}

func generateIDs() {
	ID1 = generateBidConvID("Test")
	ID2 = generateBidConvID("Test1")
}

func TestBidOfferStore_AddOffer(t *testing.T) {
	//fmt.Println("Adding New Proposer for funding")
	err := bidOfferStore.AddOffer(*NewBidOffer(ID1, TypeOffer, offerTimes[0], *balance.NewAmount(100)))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.AddOffer(*NewBidOffer(ID1, TypeCounterOffer, offerTimes[1], *balance.NewAmount(200)))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.AddOffer(*NewBidOffer(ID1, TypeOffer, offerTimes[2], *balance.NewAmount(150)))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.AddOffer(*NewBidOffer(ID2, TypeOffer, offerTimes[3], *balance.NewAmount(50)))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.AddOffer(*NewBidOffer(ID2, TypeCounterOffer, offerTimes[4], *balance.NewAmount(300)))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.AddOffer(*NewBidOffer(ID2, TypeOffer, offerTimes[5], *balance.NewAmount(250)))
	assert.NoError(t, err, "")
	cs.Commit()

}

func TestBidOfferStore_Iterate(t *testing.T) {
	//fmt.Println("Iterating Stores")
	IDLIST := []BidConvId{ID2, ID1}
	ok := bidOfferStore.iterate(func(bidConvId BidConvId, offerType BidOfferType, offerTime int64) bool {
		//fmt.Println("BidConvId : ", bidConvId, "BidOfferType :", offerType, "BidOfferTime :", offerTime)
		assert.Contains(t, IDLIST, bidConvId, "")
		return false
	})
	assert.True(t, ok, "")
}

//
func TestBidOfferStore_GetOffersForBidConvId(t *testing.T) {
	//fmt.Println("Get offers for ID :  ", ID1)
	offers := bidOfferStore.GetOffersForBidConvId(ID1, func(bidConvId BidConvId, offerType BidOfferType, offerTime int64) BidOffer {
		bidOffer, err := bidOfferStore.get(assembleBidOfferKey(bidConvId, offerType, offerTime))
		if err != nil {
			return BidOffer{}
		}
		return *bidOffer
	})
	//for _, offer := range offers {
	//	offer.Print()
	//}
	assert.EqualValues(t, 3, len(offers), "")
}

//
func TestBidOfferStore_GetActiveOfferForBidConvId(t *testing.T) {
	//fmt.Println("Get active offer for id :")

	offer, err := bidOfferStore.GetActiveOfferForBidConvId(ID1)
	assert.NoError(t, err, "")
	assert.NotEmpty(t, offer)
}