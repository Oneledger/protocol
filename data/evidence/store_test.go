package evidence

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	memDb       db.DB
	store       *EvidenceStore
	cs          *storage.State
	validators  []keys.Address
	allegations []*AllegationRequest
	suspicious  []*LastValidatorHistory
)

func setupStore() {
	fmt.Println("####### Testing  evidence store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("state", memDb))

	store = NewEvidenceStore("es", cs)
	generateEvidences()
}

func generateEvidences() {
	for i := 0; i < 3; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		validators = append(validators, h.Address())
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].String() < validators[j].String()
	})

	// create allegations
	id1, _ := store.GenerateRequestID()
	id2, _ := store.GenerateRequestID()
	ar1 := NewAllegationRequest(id1, validators[0], validators[1], 10, "did something bad")
	ar2 := NewAllegationRequest(id2, validators[0], validators[2], 20, "too many votes")

	// sort allegations
	allegations = []*AllegationRequest{ar1, ar2}
	sort.Slice(allegations, func(i, j int) bool {
		return allegations[i].ID < allegations[j].ID
	})
}

func TestRewardsCumulativeStore_DumpLoadState(t *testing.T) {
	setupStore()

	// create allegation1 and vote
	err := store.SetAllegationRequest(allegations[0])
	assert.Nil(t, err)
	vote1 := &AllegationVote{
		Address: validators[0],
		Choice:  YES,
	}
	vote2 := &AllegationVote{
		Address: validators[1],
		Choice:  NO,
	}
	allegations[0].Votes = []*AllegationVote{vote1, vote2}
	err = store.Vote(allegations[0].ID, vote1.Address, vote1.Choice)
	assert.Nil(t, err)
	err = store.Vote(allegations[0].ID, vote2.Address, vote2.Choice)
	assert.Nil(t, err)

	// create allegation2 and vote
	err = store.SetAllegationRequest(allegations[1])
	assert.Nil(t, err)
	vote3 := &AllegationVote{
		Address: validators[0],
		Choice:  YES,
	}
	vote4 := &AllegationVote{
		Address: validators[1],
		Choice:  YES,
	}
	allegations[1].Votes = []*AllegationVote{vote3, vote4}
	err = store.Vote(allegations[1].ID, vote3.Address, vote3.Choice)
	assert.Nil(t, err)
	err = store.Vote(allegations[1].ID, vote4.Address, vote4.Choice)
	assert.Nil(t, err)

	// create suspicious validators
	now := time.Now()
	tm1 := now.AddDate(0, 0, 1)
	tm2 := tm1.AddDate(0, 0, 7)
	lvh1, _ := store.CreateSuspiciousValidator(validators[1], MISSED_REQUIRED_VOTES, 1100, &tm2)
	lvh2, _ := store.CreateSuspiciousValidator(validators[2], BYZANTINE_FAULT, 5600, &tm2)
	assert.Nil(t, err)
	suspicious = []*LastValidatorHistory{lvh1, lvh2}

	// dump chain state
	store.state.Commit()
	state, succeed := store.DumpState()
	assert.True(t, succeed)
	//Assert Allegations
	for i, allegation := range allegations {
		assert.EqualValues(t, allegation, state.Allegations[i])
	}
	//Assert Suspicious Validators
	for i, validator := range suspicious {
		assert.EqualValues(t, validator, state.SuspiciousValidators[i])
	}

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
	assert.EqualValues(t, state, stateLoaded)
}
