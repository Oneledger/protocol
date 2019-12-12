/*

 */

package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

func TestTrackerMarshal(t *testing.T) {

	addrs := []keys.Address{
		hexDump("aaaaaa"),
		hexDump("bbbbbb"),
		hexDump("cccccc"),
		hexDump("dddddd"),
	}

	tracker, err := NewTracker([]byte("lockscriptaddress"), 3, addrs)

	assert.NoError(t, err)

	db, err := storage.GetDatabase("chainstate", "/tmp/btc_test", "goleveldb")
	assert.NoError(t, err)

	cs := storage.NewChainState("testing", db)
	state := storage.NewState(cs)
	ts := NewTrackerStore("btct", state)

	err = ts.SetTracker("test", tracker)
	assert.NoError(t, err)

	trackerNew, err := ts.Get("test")
	assert.NoError(t, err)

	assert.Equal(t, trackerNew.Multisig.Signers[0].String(), tracker.Multisig.Signers[0].String())
	assert.Equal(t, trackerNew.Multisig.Signers[1].String(), tracker.Multisig.Signers[1].String())
	assert.Equal(t, trackerNew.Multisig.Signers[2].String(), tracker.Multisig.Signers[2].String())
	assert.Equal(t, trackerNew.Multisig.Signers[3].String(), tracker.Multisig.Signers[3].String())
}

func hexDump(a string) []byte {
	b, _ := hex.DecodeString(a)
	return b
}
