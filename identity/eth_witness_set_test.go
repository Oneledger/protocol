package identity

import (
	"encoding/hex"
	"testing"

	"github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/storage"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
)

// test setup
func setupEthWitnessStore() *WitnessStore {
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewChainState("chainstate", db)
	ws := NewWitnessStore("etw", storage.NewState(cs))
	return ws
}

func setupInitialWitness(ws *WitnessStore) []keys.Address {
	addresses := make([]keys.Address, 4)
	addr0, _ := hex.DecodeString("72143ADE3D941025468792311A0AB38D5085E15A")
	addr1, _ := hex.DecodeString("A21437DF3C9410254A8792311A0A13255085E157")
	addr2, _ := hex.DecodeString("B2143CDE3D941025468792311A0AB38D5085E151")
	addr3, _ := hex.DecodeString("F2143AD5793B910D9410225ADC68B38D5085E11C")
	addresses[0] = addr0
	addresses[1] = addr1
	addresses[2] = addr2
	addresses[3] = addr3

	stakes := make([]Stake, 2)
	stakes[0] = Stake{
		ValidatorAddress: addr0,
		StakeAddress:     addr0,
		Pubkey: keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    nil,
		},
		Name: "test_node0",
	}
	stakes[1] = Stake{
		ValidatorAddress: addr1,
		StakeAddress:     addr1,
		Pubkey: keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    nil,
		},
		Name: "test_node1",
	}

	for _, stake := range stakes {
		ws.AddWitness(chain.ETHEREUM, stake)
	}

	ws.store.Commit()

	return addresses
}

func TestNewEthWitnessStore(t *testing.T) {
	ws := setupEthWitnessStore()
	assert.NotEmpty(t, ws)
}

func TestEthWitnessStore_Init_IsETHWitness(t *testing.T) {
	t.Run("run with witness node, should return true", func(t *testing.T) {
		ws := setupEthWitnessStore()
		addrs := setupInitialWitness(ws)
		ws.Init(chain.ETHEREUM, addrs[0])
		assert.Equal(t, true, ws.IsETHWitness())
	})
	t.Run("run with non-witness node, should return false", func(t *testing.T) {
		ws := setupEthWitnessStore()
		addrs := setupInitialWitness(ws)
		ws.Init(chain.ETHEREUM, addrs[2])
		assert.Equal(t, false, ws.IsETHWitness())
	})
}

func TestEthWitnessStore_Get_Exists(t *testing.T) {
	t.Run("get by valid witness address", func(t *testing.T) {
		ws := setupEthWitnessStore()
		addrs := setupInitialWitness(ws)
		exist := ws.Exists(chain.ETHEREUM, addrs[0])
		assert.True(t, exist)

		witness, err := ws.Get(chain.ETHEREUM, addrs[0])
		assert.Nil(t, err)
		assert.NotNil(t, witness)
		assert.Equal(t, witness.Address, addrs[0])
	})
	t.Run("get by invalid witness address", func(t *testing.T) {
		ws := setupEthWitnessStore()
		addrs := setupInitialWitness(ws)
		exist := ws.Exists(chain.ETHEREUM, addrs[2])
		assert.False(t, exist)

		witness, err := ws.Get(chain.ETHEREUM, addrs[2])
		assert.Nil(t, witness)
		assert.EqualError(t, err, "failed to get witness from store")
	})
}

func TestEthWitnessStore_Iterate(t *testing.T) {
	ws := setupEthWitnessStore()
	addrs := setupInitialWitness(ws)
	addrs_expected := addrs[:2]

	addrs_actual := []keys.Address{}
	ws.Iterate(chain.ETHEREUM, func(addr keys.Address, witness *Witness) bool {
		addrs_actual = append(addrs_actual, addr)
		return false
	})
	assert.Equal(t, addrs_expected, addrs_actual)
}

func TestEthWitnessStore_GetETHWitnessAddresses(t *testing.T) {
	ws := setupEthWitnessStore()
	addrs := setupInitialWitness(ws)
	addrs_expected := addrs[:2]
	addrs_actual, _ := ws.GetWitnessAddresses(chain.ETHEREUM)
	assert.EqualValues(t, addrs_expected, addrs_actual)
}

func TestEthWitnessStore_IsETHWitnessAddress(t *testing.T) {
	ws := setupEthWitnessStore()
	addrs := setupInitialWitness(ws)
	assert.True(t, ws.IsWitnessAddress(chain.ETHEREUM, addrs[0]))
	assert.True(t, ws.IsWitnessAddress(chain.ETHEREUM, addrs[1]))
	assert.False(t, ws.IsWitnessAddress(chain.ETHEREUM, addrs[2]))
	assert.False(t, ws.IsWitnessAddress(chain.ETHEREUM, addrs[3]))
}

func TestEthWitnessStore_AddWitness(t *testing.T) {
	ws := setupEthWitnessStore()
	addr, _ := hex.DecodeString("F2143AD5793BA10D9410225ADC68B38D5085E11C")
	stake := Stake{
		ValidatorAddress: addr,
		StakeAddress:     addr,
		Pubkey: keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    nil,
		},
		Name: "test_node_new",
	}
	err := ws.AddWitness(chain.ETHEREUM, stake)
	assert.Nil(t, err)
	witness, err = ws.Get(chain.ETHEREUM, addr)
	assert.Nil(t, err)
	assert.Equal(t, keys.Address(addr), witness.Address)
}
