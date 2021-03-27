package action

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
)

func TestStateDBRunner(t *testing.T) {
	db := db.NewDB("test", db.MemDBBackend, "")

	balances := balance.NewStore("tb", storage.NewState(storage.NewChainState("balance", db)))

	currencies := balance.NewCurrencySet()
	currency := balance.Currency{
		Name:    "OLT",
		Chain:   chain.Type(1),
		Decimal: int64(18),
	}
	currencies.Register(currency)

	logger = log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")

	stateDB := NewCommitStateDB(
		evm.NewContractStore(storage.NewState(storage.NewChainState("contracts", db))),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)),
			balances,
			currencies,
		),
		logger,
	)

	from, _, _ := generateKeyPair()

	acc := &balance.EthAccount{
		Address: from.Bytes(),
		Coins: balance.Coin{
			Currency: balance.Currency{
				Id:      0,
				Name:    "OLT",
				Chain:   0,
				Decimal: 18,
				Unit:    "nue",
			},
			Amount: balance.NewAmountFromInt(10000),
		},
	}
	stateDB.GetAccountKeeper().SetAccount(*acc)

	t.Run("test contract code store and it is ok", func(t *testing.T) {
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b5061016d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80635f76f6ab1461003b5780636d4ce63c1461006b575b600080fd5b6100696004803603602081101561005157600080fd5b8101908080351515906020019092919050505061008b565b005b6100736100e4565b60405180821515815260200191505060405180910390f35b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690509056fea2646970667358221220ef09e2f46f4d83d3c8af213cd936666dbb273e3f612b70d008a1d8bbf6d14a1d64736f6c63430007040033")
		stateObject := newStateObject(stateDB, acc)
		stateObject.SetCode(acc.EthAddress().Hash(), code)
		fmt.Printf("code: %s\n", code)
		fmt.Printf("stateObject.CodeHash(): %s\n", stateObject.CodeHash())
		assert.True(t, bytes.Equal(code, stateObject.CodeHash()), "Wrong code in cache")
	})
}
