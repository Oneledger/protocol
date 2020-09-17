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
	numBatches     = 9
	prefix = "testPrefix"
)

var (
	products []*Product

	productStore *ProductStore
)

func init() {
	fmt.Println("####### TESTING PRODUCT STORE #######")

	//Create new bid conversations
	for i := 0; i < numBatches; i++ {

		products = append(products, NewProduct(BatchID("10000000" + strconv.Itoa(i)), "apples", "F12345", "countryHome", "field A", time.Now().UTC().Unix(), "AAA", 100, "very good product"))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create bid conversation store
	productStore = NewProductStore(cs, prefix)
}

func TestBidConvStore_Set(t *testing.T) {
	fmt.Println("products: ", products)
	fmt.Println("productStore: ", productStore)
	err := productStore.Set(products[0])
	assert.Equal(t, nil, err)

	batch, err := productStore.Get(products[0].BatchID)
	assert.Equal(t, nil, err)

	assert.Equal(t, batch.BatchID, products[0].BatchID)
}

func TestBidConvStore_Exists(t *testing.T) {
	exists := productStore.Exists(products[0].BatchID)
	assert.Equal(t, true, exists)

	exists = productStore.Exists(products[1].BatchID)
	assert.Equal(t, false, exists)
}

func TestBidConvStore_Delete(t *testing.T) {
	_, err := productStore.Get(products[0].BatchID)
	assert.Equal(t, nil, err)

	res, err := productStore.Delete(products[0].BatchID)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = productStore.Get(products[0].BatchID)
	assert.NotEqual(t, nil, err)
}

