package delegation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	memDb      db.DB
	store      *DelegationStore
	cs         *storage.State
	validator1 keys.Address
	validator2 keys.Address
	stakeAddr1 keys.Address
	stakeAddr2 keys.Address
	zero       *balance.Amount
	amt1       *balance.Amount
	amt2       *balance.Amount
	amt3       *balance.Amount
	amt4       *balance.Amount
	amt5       *balance.Amount
	withdraw1  *balance.Amount
	withdraw2  *balance.Amount
)

func setup() {
	fmt.Println("####### Testing delegation store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("delegation", memDb))

	store = NewDelegationStore("st", cs)
	generateAddresses()
}

func genAddress() keys.Address {
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub1.GetHandler()
	return h.Address()
}

func generateAddresses() {
	validator1 = genAddress()
	validator2 = genAddress()
	stakeAddr1 = genAddress()
	stakeAddr2 = genAddress()

	zero = balance.NewAmount(0)
	amt1 = balance.NewAmount(100)
	amt2 = balance.NewAmount(200)
	amt3 = balance.NewAmount(377)
	withdraw1 = balance.NewAmount(163)
	withdraw2 = balance.NewAmount(499)
}

func TestDelegationStore_DumpLoadState(t *testing.T) {
	setup()

	// validator1 stake/unstake/withdraw
	err := store.Stake(validator1, stakeAddr1, *balance.NewAmount(100))
	assert.Nil(t, err)
	err = store.Unstake(validator1, stakeAddr1, *balance.NewAmount(30), 11)
	assert.Nil(t, err)
	err = store.Withdraw(validator1, stakeAddr1, *balance.NewAmount(10))
	assert.Nil(t, err)

	// validator2 stake/unstake/withdraw
	err = store.Stake(validator2, stakeAddr2, *balance.NewAmount(200))
	assert.Nil(t, err)
	err = store.Unstake(validator2, stakeAddr2, *balance.NewAmount(70), 15)
	assert.Nil(t, err)
	err = store.Withdraw(validator2, stakeAddr2, *balance.NewAmount(60))
	assert.Nil(t, err)
	store.state.Commit()

	// prepare to dump
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "genesis.json")
	writer, err := os.Create(file)
	assert.Nil(t, err)
	defer func() { _ = os.Remove(file) }()
	state, err := store.dumpState()
	assert.Nil(t, err)

	// dump to Genesis
	str, err := json.MarshalIndent(state, "", " ")
	assert.Nil(t, err)
	_, err = writer.Write(str)
	assert.Nil(t, err)
	err = writer.Close()
	assert.Nil(t, err)

	// load from Genesis
	reader, err := os.Open(file)
	stateBytes, _ := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	stateLoaded := NewRewardCumuState()
	err = json.Unmarshal(stateBytes, stateLoaded)
	assert.Nil(t, err)
	assert.Equal(t, state, stateLoaded)
}
