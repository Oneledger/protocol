package bid_data

import (
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"testing"
)

var (
	bidOfferStore  *BidOfferStore
	cs             *storage.State
	address        keys.Address
	address2       keys.Address
	ID1            BidConvId
	ID2            BidConvId
	offerTimes     []int64
	offerTime      int64 = 1596701481
	offerTimeCount       = 7
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
		offerTimes = append(offerTimes, offerTime+int64(i*20))
	}
}

func generateIDs() {
	ID1 = generateBidConvID("Test", 123)
	ID2 = generateBidConvID("Test1", 234)
}

func TestBidOfferStore_AddActiveOffer(t *testing.T) {
	//fmt.Println("Adding New Proposer for funding")
	err := bidOfferStore.SetActiveOffer(*NewBidOffer(ID1, TypeBidOffer, offerTimes[0], *action.NewAmount("OLT", *balance.NewAmount(100)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetActiveOffer(*NewBidOffer(ID1, TypeCounterOffer, offerTimes[1], *action.NewAmount("OLT", *balance.NewAmount(200)), CounterOfferAmount))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetActiveOffer(*NewBidOffer(ID1, TypeBidOffer, offerTimes[2], *action.NewAmount("OLT", *balance.NewAmount(150)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetActiveOffer(*NewBidOffer(ID2, TypeBidOffer, offerTimes[3], *action.NewAmount("OLT", *balance.NewAmount(50)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetActiveOffer(*NewBidOffer(ID2, TypeCounterOffer, offerTimes[4], *action.NewAmount("OLT", *balance.NewAmount(300)), CounterOfferAmount))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetActiveOffer(*NewBidOffer(ID2, TypeBidOffer, offerTimes[5], *action.NewAmount("OLT", *balance.NewAmount(250)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
}

func TestBidOfferStore_SetInActiveOffer(t *testing.T) {
	//fmt.Println("Adding New Proposer for funding")
	err := bidOfferStore.SetInActiveOffer(*NewBidOffer(ID1, TypeBidOffer, offerTimes[0], *action.NewAmount("OLT", *balance.NewAmount(100)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetInActiveOffer(*NewBidOffer(ID1, TypeCounterOffer, offerTimes[1], *action.NewAmount("OLT", *balance.NewAmount(200)), CounterOfferAmount))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetInActiveOffer(*NewBidOffer(ID2, TypeBidOffer, offerTimes[3], *action.NewAmount("OLT", *balance.NewAmount(50)), BidAmountLocked))
	assert.NoError(t, err, "")
	cs.Commit()
	err = bidOfferStore.SetInActiveOffer(*NewBidOffer(ID2, TypeCounterOffer, offerTimes[4], *action.NewAmount("OLT", *balance.NewAmount(300)), CounterOfferAmount))
	assert.NoError(t, err, "")
	cs.Commit()

}

func TestBidOfferStore_Iterate(t *testing.T) {
	//fmt.Println("Iterating Stores")
	IDLIST := []BidConvId{ID2, ID1}
	ok := bidOfferStore.iterate(func(bidConvId BidConvId, offerType BidOfferType, offerTime int64, offer BidOffer) bool {
		fmt.Println("BidConvId : ", bidConvId, "BidOfferType :", offerType, "BidOfferTime :", offerTime)
		assert.Contains(t, IDLIST, bidConvId, "")
		return false
	})
	assert.True(t, ok, "")
}


func TestBidOfferStore_GetOffers(t *testing.T) {
	fmt.Println("Get offers for ID :  ", ID1)
	offer, err := bidOfferStore.GetActiveOffer(ID1, TypeInvalid)
	if err != nil {
		assert.Error(t, err,"")
	}
	fmt.Print(offer.Amount)
	offer, err = bidOfferStore.GetActiveOffer(ID2, TypeInvalid)
	if err != nil {
		assert.Error(t, err,"")
	}
	fmt.Print(offer.Amount)

	offers := bidOfferStore.GetInActiveOffers(ID1, TypeInvalid)
	for _, offer := range offers {
		fmt.Print(offer.Amount)
	}
	assert.EqualValues(t, 2, len(offers), "")

	offers = bidOfferStore.GetInActiveOffers(ID1, TypeInvalid)
	for _, offer := range offers {
		fmt.Print(offer.Amount)
	}
	assert.EqualValues(t, 2, len(offers), "")
}
