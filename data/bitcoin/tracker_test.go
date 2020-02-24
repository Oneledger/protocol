/*

 */

package bitcoin

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/config"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	store   *TrackerStore
	cs      *storage.State
	tracker *Tracker
)

func init() {
	addrs := []keys.Address{
		hexDump("aaaaaa"),
		hexDump("bbbbbb"),
		hexDump("cccccc"),
		hexDump("dddddd"),
	}

	var err error
	tracker, err = NewTracker([]byte("lockscriptaddress"), 3, addrs)
	if err != nil {
		return
	}

	db, err := storage.GetDatabase("chainstate", "/tmp/btc_test", "goleveldb")
	if err != nil {
		return
	}

	cs := storage.NewChainState("testing", db)
	state := storage.NewState(cs)
	store = NewTrackerStore("btct", state)
}

func TestTrackerMarshal(t *testing.T) {

	err := store.SetTracker("test", tracker)
	assert.NoError(t, err)

	trackerNew, err := store.Get("test")
	assert.NoError(t, err)

	assert.Equal(t, trackerNew.Multisig.Signers[0].String(), tracker.Multisig.Signers[0].String())
	assert.Equal(t, trackerNew.Multisig.Signers[1].String(), tracker.Multisig.Signers[1].String())
	assert.Equal(t, trackerNew.Multisig.Signers[2].String(), tracker.Multisig.Signers[2].String())
	assert.Equal(t, trackerNew.Multisig.Signers[3].String(), tracker.Multisig.Signers[3].String())
}

func TestNewBTCConfig(t *testing.T) {

	cdConfig := config.DefaultChainDriverConfig()

	btcConfig := NewBTCConfig(cdConfig, "testnet3")
	store.SetConfig(btcConfig)
	storedConfig := store.GetConfig()

	assert.Equal(t, storedConfig.BTCParams.Name, "testnet3")
}

func hexDump(a string) []byte {
	b, _ := hex.DecodeString(a)
	return b
}
