package data

import (
	"testing"

	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	"github.com/tendermint/tm-db"
)

const (
	testType = "t"
)

var (
	chainstate *storage.ChainState
	test       *testStore
	stores     Router
)

var _ ExtStore = &testStore{}

type testStore struct {
	state  *storage.State
	prefix []byte //Current Store Prefix
}

func newTestStore(state *storage.State, prefix []byte) *testStore {
	return &testStore{state: state, prefix: prefix}
}

func (ts *testStore) Set(testKey string, testData string) error {
	prefixed := append(ts.prefix, testKey...)
	data := testData

	err := ts.state.Set(prefixed, []byte(data))

	return err
}

func (ts *testStore) WithState(state *storage.State) ExtStore {
	ts.state = state
	return ts
}

func (ts *testStore) Get(testKey string) (string, error) {
	var testData string
	prefixed := append(ts.prefix, testKey...)
	data, err := ts.state.Get(prefixed)
	if err != nil {
		return "", err
	}
	testData = string(data)
	return testData, nil
}

func init() {

	//Create a store
	db := db.NewDB("testDB", db.MemDBBackend, "")
	chainstate = storage.NewChainState("chainstate", db)

	test = newTestStore(storage.NewState(chainstate), []byte("test"))

	//Create data router
	stores = NewStorageRouter()

}

func TestStorageRouter_Add(t *testing.T) {

	//Add store to the data router
	_ = stores.Add(testType, test)

	db, _ := stores.Get(testType)
	testDb := db.(*testStore)

	testDb.Set("testKey", "testData")

}

func TestStorageRouter_Get(t *testing.T) {
	db, _ := stores.Get(testType)
	testDb := db.(*testStore)

	data, _ := testDb.Get("testKey")
	assert.Equal(t, data, "testData")

}
