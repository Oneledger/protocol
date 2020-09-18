package farm_data

import (
	"fmt"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"strconv"
	"testing"
	"time"
)

const (
	numBatches = 9
	prefix     = "testPrefix"
)

var (
	produces []*Produce

	produceStore *ProduceStore
)

func init() {
	fmt.Println("####### TESTING PRODUCE STORE #######")

	for i := 0; i < numBatches; i++ {

		produces = append(produces, NewProduce(BatchID("10000000"+strconv.Itoa(i)), "apples", "F12345", "countryHome", "field A", time.Now().UTC().Unix(), "AAA", 100, "very good product"))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	produceStore = NewProduceStore(cs, prefix)
}

func TestProduceStore_Set(t *testing.T) {
	fmt.Println("products: ", produces)
	fmt.Println("produceStore: ", produceStore)
	err := produceStore.Set(produces[0])
	assert.Equal(t, nil, err)

	batch, err := produceStore.Get(produces[0].BatchID)
	assert.Equal(t, nil, err)

	assert.Equal(t, batch.BatchID, produces[0].BatchID)
}

func TestProduceStore_Exists(t *testing.T) {
	exists := produceStore.Exists(produces[0].BatchID)
	assert.Equal(t, true, exists)

	exists = produceStore.Exists(produces[1].BatchID)
	assert.Equal(t, false, exists)
}

func TestProduceStore_Delete(t *testing.T) {
	_, err := produceStore.Get(produces[0].BatchID)
	assert.Equal(t, nil, err)

	res, err := produceStore.Delete(produces[0].BatchID)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = produceStore.Get(produces[0].BatchID)
	assert.NotEqual(t, nil, err)
}
