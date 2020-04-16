package storage

import (
	"math/rand"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/tendermint/tm-db"
)

var testcase = make(map[string][]byte)
var keys = make([]string, 0, 100)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func init() {

	for i := 0; i < 100; i++ {
		key := String(rand.Int() % 10)
		testcase[key] = []byte(String(rand.Int() % 10))
		keys = append(keys, key)
	}

}

func TestState_BeginTxSession(t *testing.T) {
	state := NewState(NewChainState("test", getCacheDB()))

	state.BeginTxSession()

	for _, k := range keys {
		v, _ := testcase[k]
		state.Set([]byte(k), v)
	}

	state.CommitTxSession()

	hash, version := state.Commit()

	state2 := NewState(NewChainState("test1", getCacheDB()))

	state2.BeginTxSession()
	for _, k := range keys {
		v, _ := testcase[k]
		state2.Set([]byte(k), v)
	}
	state2.CommitTxSession()
	hash2, version2 := state2.Commit()

	assert.Equal(t, hash2, hash, "hash should match")
	assert.Equal(t, version2, version, "version should match")
}

func TestState_Commit(t *testing.T) {
	state := NewState(NewChainState("test", getCacheDB()))

	for _, k := range keys {
		v, _ := testcase[k]
		state.Set([]byte(k), v)
	}

	hash, version := state.Commit()

	state2 := NewState(NewChainState("test1", getCacheDB()))

	for _, k := range keys {
		v, _ := testcase[k]
		state2.Set([]byte(k), v)
	}

	hash2, version2 := state2.Commit()

	assert.Equal(t, hash2, hash, "hash should match")
	assert.Equal(t, version2, version, "version should match")
}

func TestState_CommitTxSession(t *testing.T) {
	state := NewState(NewChainState("test", getCacheDB()))

	for _, k := range keys {
		v, _ := testcase[k]
		state.Set([]byte(k), v)
	}

	hash, version := state.Commit()

	state2 := NewState(NewChainState("test1", getCacheDB()))

	for _, k := range keys {
		state2.BeginTxSession()
		v, _ := testcase[k]
		state2.Set([]byte(k), v)
		state2.CommitTxSession()
	}

	hash2, version2 := state2.Commit()

	assert.Equal(t, hash2, hash, "hash should match")
	assert.Equal(t, version2, version, "version should match")
}
func getCacheDB() db.DB {
	return db.NewDB("test", db.MemDBBackend, "")
}
