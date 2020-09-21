package vpart_data

import (
	"fmt"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"strconv"
	"testing"
)

const (
	numParts = 9
	prefix   = "testPrefix"
)

var (
	parts []*VPart

	vPartStore *VPartStore
)

func init() {
	fmt.Println("####### TESTING PRODUCE STORE #######")

	for i := 0; i < numParts; i++ {

		parts = append(parts, NewVPart(Vin("1000000000000000"+strconv.Itoa(i)), "engine", "123", "aaa", "123456789", 2000))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	vPartStore = NewVPartStore(cs, prefix)
}

func TestVPartStore_Set(t *testing.T) {
	fmt.Println("parts: ", parts)
	fmt.Println("vPartStore: ", vPartStore)
	err := vPartStore.Set(parts[0])
	assert.Equal(t, nil, err)

	batch, err := vPartStore.Get(parts[0].VIN, parts[0].PartType)
	assert.Equal(t, nil, err)

	assert.Equal(t, batch.VIN, parts[0].VIN)
}

func TestVPartStore_Exists(t *testing.T) {
	exists := vPartStore.Exists(parts[0].VIN, parts[0].PartType)
	assert.Equal(t, true, exists)

	exists = vPartStore.Exists(parts[1].VIN, parts[1].PartType)
	assert.Equal(t, false, exists)
}

func TestVPartStore_Delete(t *testing.T) {
	_, err := vPartStore.Get(parts[0].VIN, parts[0].PartType)
	assert.Equal(t, nil, err)

	res, err := vPartStore.Delete(parts[0].VIN, parts[0].PartType)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = vPartStore.Get(parts[0].VIN, parts[0].PartType)
	assert.NotEqual(t, nil, err)
}

