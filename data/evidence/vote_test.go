/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package evidence

import (
	"errors"
	"os"
	"testing"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/abci/types"
	db "github.com/tendermint/tm-db"
)

// global setup
func setup() []string {
	testDBs := []string{"test_dbpath"}
	return testDBs
}

// remove test db dir after
func teardown(dbPaths []string) {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
}

func TestCumulativeVote(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)

	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("evidence", db))

	es := NewEvidenceStore("tes", cs)

	Y := int64(3)
	vM := make(map[string]bool)

	baddr1 := keys.Address{}
	_ = baddr1.UnmarshalText([]byte("0lte952e380a48d5237630fae75a79d7f7616ff35a9"))

	baddr2 := keys.Address{}
	_ = baddr2.UnmarshalText([]byte("0lt6d2b781f89132ccfefdcdd9453bbc34651cf5c67"))

	baddr3 := keys.Address{}
	_ = baddr3.UnmarshalText([]byte("0ltea0dbe4afc96eb320533a149f42675eecd0b8987"))

	votes := []types.VoteInfo{
		{
			Validator: types.Validator{
				Address: baddr1,
				Power:   1,
			},
			SignedLastBlock: false,
		},
		{
			Validator: types.Validator{
				Address: baddr2,
				Power:   2,
			},
			SignedLastBlock: true,
		},
		{
			Validator: types.Validator{
				Address: baddr3,
				Power:   5,
			},
			SignedLastBlock: true,
		},
	}

	for _, vote := range votes {
		hAddr := keys.Address(vote.Validator.Address).Humanize()
		vM[hAddr] = vote.SignedLastBlock
	}

	for i := int64(1); i <= Y; i++ {
		err := es.SetVoteBlock(i, votes)
		assert.NoError(t, err)
		cv, err := es.GetCumulativeVote()
		assert.NoError(t, err)
		err = es.SetCumulativeVote(cv, i, Y)
		assert.NoError(t, err)
	}

	for i := int64(1); i <= Y; i++ {
		vb, err := es.GetVoteBlock(i)
		assert.NoError(t, err)
		assert.Equal(t, i, vb.Height)

		assert.Equal(t, 2, len(vb.Addresses))

		for _, address := range vb.Addresses {
			res, _ := vM[address.Humanize()]
			assert.True(t, res)
		}
	}

	// check cumulative vote
	cv, err := es.GetCumulativeVote()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(cv.Addresses))

	for i := range cv.Addresses {
		count := cv.Addresses[i]
		assert.Equal(t, int64(3), count)
	}

	baddr4 := keys.Address{}
	_ = baddr4.UnmarshalText([]byte("0lt581891d8411faaa15bfd3020b8bd78942cb74ecb"))
	// adding H > Y
	votes = []types.VoteInfo{
		{
			Validator: types.Validator{
				Address: baddr4,
				Power:   5,
			},
			SignedLastBlock: true,
		},
	}
	err = es.SetVoteBlock(Y+1, votes)
	assert.NoError(t, err)
	cv, err = es.GetCumulativeVote()
	assert.NoError(t, err)
	err = es.SetCumulativeVote(cv, Y+1, Y)
	assert.NoError(t, err)

	// check cumulative vote
	cv, err = es.GetCumulativeVote()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(cv.Addresses))

	data := make(map[string]int64)
	data["0lt6d2b781f89132ccfefdcdd9453bbc34651cf5c67"] = 2
	data["0ltea0dbe4afc96eb320533a149f42675eecd0b8987"] = 2
	data["0lt581891d8411faaa15bfd3020b8bd78942cb74ecb"] = 1
	for addr := range cv.Addresses {
		count := cv.Addresses[addr]
		assert.Equal(t, data[addr], count)
	}
}
