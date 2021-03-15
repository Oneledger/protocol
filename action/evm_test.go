package action

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	db "github.com/tendermint/tm-db"
)

func setupServer() config.Server {
	networkConfig := &config.NetworkConfig{
		P2PAddress: "tcp://127.0.0.1:26601",
		RPCAddress: "tcp://127.0.0.1:26600",
		SDKAddress: "http://127.0.0.1:26603",
	}
	consensus := &config.ConsensusConfig{}
	p2pConfig := &config.P2PConfig{}
	mempool := &config.MempoolConfig{}
	nodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		FastSync: true,
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	config := config.Server{
		Node:      nodeConfig,
		Network:   networkConfig,
		Consensus: consensus,
		P2P:       p2pConfig,
		Mempool:   mempool,
	}
	return config
}

func assemblyCtxData(acc accounts.Account, currencyName string, currencyDecimal int, setStore bool, setLogger bool, setCoin bool, setCoinAddr crypto.Address) *Context {
	ctx := &Context{}
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	// store
	var store *balance.Store
	if setStore {
		store = balance.NewStore("tb", cs)
		ctx.Balances = store
	}
	// logger
	ctx.Logger = log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")
	// currencyList
	if currencyName != "" {
		// register new token OTT
		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    currencyName,
			Chain:   chain.Type(1),
			Decimal: int64(currencyDecimal),
		}
		err := currencyList.Register(currency)
		if err != nil {
			errors.New("register new token error")
		}
		ctx.Currencies = currencyList

		// set coin for account
		amt, err := balance.NewAmountFromString("100000000000000000000", 10)
		if setCoin {
			coin := balance.Coin{
				Currency: currency,
				Amount:   amt,
			}
			coin.MultiplyInt(currencyDecimal)
			err = store.AddToAddress(setCoinAddr.Bytes(), coin)
			if err != nil {
				errors.New("setup testing token balance error")
			}
			store.State.Commit()
			ctx.Balances = store
		}
	}
	ctx.StateDB = NewCommitStateDB(
		evm.NewContractStore(storage.NewState(storage.NewChainState("contracts", db))),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)),
			ctx.Balances,
			ctx.Currencies,
		),
		ctx.Logger,
	)
	ctx.FeeOpt = &fees.FeeOption{
		FeeCurrency: balance.Currency{
			Id:      0,
			Name:    "OLT",
			Chain:   0,
			Decimal: 18,
			Unit:    "nue",
		},
		MinFeeDecimal: 9,
	}

	ctx.GovernanceStore = governance.NewStore("tg", cs)
	ctx.FeePool = &fees.Store{}
	ctx.FeePool.SetupOpt(ctx.FeeOpt)
	proposalStore := governance.ProposalStore{}
	pOpt := governance.ProposalOptionSet{
		ConfigUpdate:      governance.ProposalOption{},
		CodeChange:        governance.ProposalOption{},
		General:           governance.ProposalOption{},
		BountyProgramAddr: "TestAddress",
	}
	proposalStore.SetOptions(&pOpt)
	ctx.ProposalMasterStore = &governance.ProposalMasterStore{
		Proposal:     &proposalStore,
		ProposalFund: nil,
		ProposalVote: nil,
	}

	rwz := rewards.NewRewardStore("r", "ri", "ra", cs)
	rwzc := rewards.NewRewardCumulativeStore("rc", cs)
	ctx.RewardMasterStore = rewards.NewRewardMasterStore(rwz, rwzc)
	rewardOptions := rewards.Options{
		RewardInterval:    150,
		RewardPoolAddress: "rewardspool",
	}
	ctx.RewardMasterStore.SetOptions(&rewardOptions)
	ctx.GovernanceStore.WithHeight(0).SetFeeOption(*ctx.FeeOpt)
	ctx.GovernanceStore.WithHeight(0).SetProposalOptions(pOpt)
	ctx.GovernanceStore.WithHeight(0).SetRewardOptions(rewardOptions)
	ctx.GovernanceStore.WithHeight(0).SetAllLUH()

	ctx.Header = &abci.Header{
		Height:  1,
		Time:    time.Now().AddDate(0, 0, 1),
		ChainID: "test-1",
	}
	return ctx
}

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
}

func assertExcError(t *testing.T, contractAddr ethcmn.Address, gas, leftOverGas, expected uint64) {
	msg := fmt.Sprintf("Wrong calculation for contract '0x%s' (gas used %d, but expected %d)", ethcmn.Bytes2Hex(contractAddr.Bytes()), gas-leftOverGas, expected)
	assert.Equal(t, uint64(expected), gas-leftOverGas, msg)
}

func getBool(res []byte) bool {
	trimRes := ethcmn.TrimLeftZeroes(res)
	if len(trimRes) == 1 {
		pr := int(trimRes[0])
		if pr == 1 {
			return true
		}
	}
	return false
}

func getNonce(ctx *Context, from keys.Address) (nonce uint64) {
	acc, _ := ctx.StateDB.GetAccountKeeper().GetAccount(from.Bytes())
	nonce = acc.Sequence
	return
}

func TestRunner(t *testing.T) {
	// pragma solidity >=0.7.0 <0.8.0;

	// contract Test {
	//     mapping(address => bool) private data;

	//     function set(bool res) public payable {
	//         data[msg.sender] = res;
	//     }

	//     function get() public view returns(bool) {
	//         return data[msg.sender];
	//     }
	// }

	// generating default data
	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address())
	code := ethcmn.FromHex("0x608060405234801561001057600080fd5b5061016d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80635f76f6ab1461003b5780636d4ce63c1461006b575b600080fd5b6100696004803603602081101561005157600080fd5b8101908080351515906020019092919050505061008b565b005b6100736100e4565b60405180821515815260200191505060405180910390f35b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690509056fea26469706673582212209bac4bf916f5d28c34ab5b6f59e791ea87337bed7abf384250e80a832d134f6364736f6c63430007060033")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test contract store and it is OK", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 73123)

		// execute method with true
		input := ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000001") // "set" method with true arg
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 22468)

		// get result true
		input = ethcmn.FromHex("0x6d4ce63c") // "get" method and get result true
		res, leftOverGas, err := evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		fmt.Printf("res: %+v\n", res)
		assertExcError(t, contractAddr, gas, leftOverGas, 444)
		assert.Equal(t, true, getBool(res), "Wrong data was stored")

		// execute method with false
		input = ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000000") // "set" method with false arg
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 568)

		// get result false
		input = ethcmn.FromHex("0x6d4ce63c") // "get" method and get result false
		res, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		fmt.Printf("res: %+v\n", res)
		assertExcError(t, contractAddr, gas, leftOverGas, 444)
		assert.Equal(t, false, getBool(res), "Wrong data was stored")
	})
}
