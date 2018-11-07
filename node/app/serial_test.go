package app

import (
	"flag"
	"math/big"
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/stretchr/testify/assert"
)

// Control the execution
func TestMain(m *testing.M) {
	flag.Parse()

	// Set the debug flags according to whether the -v flag is set in go test
	if testing.Verbose() {
		log.Debug("DEBUG TURNED ON")
		global.Current.Debug = true
	} else {
		log.Debug("DEBUG TURNED OFF")
		global.Current.Debug = false
	}

	// Run it all.
	code := m.Run()

	os.Exit(code)
}

func TestAccounts(t *testing.T) {
	global.Current.RootDir = "./"
	accounts := id.NewAccounts("MyAccounts")

	priv1, pub1 := id.GenerateKeys([]byte("testAccount1 password"))
	priv2, pub2 := id.GenerateKeys([]byte("testAccount1 password"))

	user1 := id.NewAccount(data.ONELEDGER, "testAccount1", pub1, priv1)
	user2 := id.NewAccount(data.ONELEDGER, "testAccount2", pub2, priv2)

	accounts.Add(user1)
	accounts.Add(user2)

	keys := accounts.FindAll()
	log.Dump("The accounts are:", keys)
}

func TestSwap(t *testing.T) {
	var swap *action.Swap

	party := action.Party{
		Key: id.AccountKey("2222222222222222222222"),
		Accounts: map[data.ChainType][]byte{
			0: []byte("01234567"),
			1: []byte("76543210"),
		},
	}

	currency := data.Currency{
		Name:  "Hey",
		Chain: 3,
		Id:    31212,
	}

	coin := data.Coin{
		Currency: currency,
		Amount:   big.NewInt(1000),
	}

	swap = &action.Swap{
		Base: nil,
		SwapMessage: action.SwapInit{
			Party:        party,
			CounterParty: party,
			Amount:       coin,
			Exchange:     coin,
			Fee:          coin,
			Gas:          coin,
		},
	}

	// Serialize the go data structure
	buffer, err := serial.Serialize(swap, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 *action.Swap

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, swap, result, "These should be equal")
}
