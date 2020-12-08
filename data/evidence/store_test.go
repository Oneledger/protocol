package evidence

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
	store      *EvidenceStore
	cs         *storage.State
	validator1 keys.Address
	validator2 keys.Address
	validator3 keys.Address
	zero       *balance.Amount
	amt1       *balance.Amount
	amt2       *balance.Amount
	amt3       *balance.Amount
	amt4       *balance.Amount
	amt5       *balance.Amount
)

func setupStore() {
	fmt.Println("####### Testing  evidence store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("rewards", memDb))

	store = NewEvidenceStore("es", cs)
	generateAddresses()
}

func generateAddresses() {
	pub1, _, _ := keys.NewKeyPairFromTendermint()
	h1, _ := pub1.GetHandler()
	validator1 = h1.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	validator2 = h2.Address()

	pub3, _, _ := keys.NewKeyPairFromTendermint()
	h3, _ := pub2.GetHandler()
	validator3 = h3.Address()

	zero = balance.NewAmount(0)
	amt1 = balance.NewAmount(100)
	amt2 = balance.NewAmount(200)
	amt3 = balance.NewAmount(377)
}

func TestRewardsCumulativeStore_DumpLoadState(t *testing.T) {
	setupStore()

	// create allegations
	id1, _ := store.GenerateRequestID()
	id2, _ := store.GenerateRequestID()
	ar1 := NewAllegationRequest(id1, validator1, validator2, 10, "did something bad")
	ar2 := NewAllegationRequest(id2, validator1, validator3, 20, "too many votes")
	err := store.SetAllegationRequest(ar1)
	assert.Nil(t, err)
	err = store.SetAllegationRequest(ar2)
	assert.Nil(t, err)
	err = store.SetAllegationRequest(ar1)
	store.state.Commit()

	// dump chain state
	state, succeed := store.DumpState()
	assert.True(t, succeed)

	// create Genesis file
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "genesis.json")
	writer, err := os.Create(file)
	assert.Nil(t, err)
	defer func() { _ = os.Remove(file) }()

	// dump to Genesis file
	str, err := json.MarshalIndent(state, "", " ")
	assert.Nil(t, err)
	_, err = writer.Write(str)
	assert.Nil(t, err)
	err = writer.Close()
	assert.Nil(t, err)

	// load from Genesis file
	reader, err := os.Open(file)
	stateBytes, _ := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	stateLoaded := NewEvidenceState()
	err = json.Unmarshal(stateBytes, stateLoaded)
	assert.Nil(t, err)
	assert.Equal(t, state, stateLoaded)
}
