package governance

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/storage"
)

var govStore *Store

func init() {
	fmt.Println("####### TESTING GOVERNANCE STORE #######")

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create Governance store
	govStore = NewStore("g", cs)
}

func TestStore_InitialChain(t *testing.T) {

}

func TestStore_Initiated(t *testing.T) {

}

func TestStore_Exists(t *testing.T) {

}

func TestStore_Set(t *testing.T) {

}

func TestStore_Get(t *testing.T) {

}

func TestStore_SetBTCChainDriverOption(t *testing.T) {

}

func TestStore_GetBTCChainDriverOption(t *testing.T) {

}

func TestStore_SetCurrencies(t *testing.T) {

}

func TestStore_GetCurrencies(t *testing.T) {

}

func TestStore_SetEpoch(t *testing.T) {

}

func TestStore_GetEpoch(t *testing.T) {

}

func TestStore_SetETHChainDriverOption(t *testing.T) {

}

func TestStore_GetETHChainDriverOption(t *testing.T) {

}

func TestStore_SetFeeOption(t *testing.T) {

}

func TestStore_GetFeeOption(t *testing.T) {

}

func TestStore_SetONSOptions(t *testing.T) {

}

func TestStore_GetONSOptions(t *testing.T) {

}

func TestStore_SetProposalOptions(t *testing.T) {

}

func TestStore_GetProposalOptions(t *testing.T) {

}

func TestStore_SetLUH(t *testing.T) {
	err := govStore.WithHeight(1000).SetLUH()
	assert.NoError(t, err, "No error Expected")
}

func TestStore_GetLUH(t *testing.T) {
	height, err := govStore.GetLUH()
	assert.NoError(t, err, "No error Expected")
	assert.EqualValues(t, 1000, height, "Expected height is 1000 from Line 99")
}

func TestStoreGetAndSet(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 300000
	for i := 0; i < 30; i++ {
		h := rand.Intn(max-min+1) + min
		err := govStore.WithHeight(int64(h)).SetLUH()
		assert.NoError(t, err, "No error Expected")
		height, err := govStore.GetLUH()
		assert.NoError(t, err, "No error Expected")
		assert.EqualValues(t, h, height, "")
		fmt.Println("Testing passed for height : ", h)
	}

}
