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

func TestRunner_ecrecover(t *testing.T) {
	// pragma solidity ^0.5.0;

	// contract test {

	// 	function verify(bytes32 _message, uint8 _v, bytes32 _r, bytes32 _s) public view returns (address) {
	// 		address signer = ecrecover(_message, _v, _r, _s);
	// 		return signer;
	// 	}

	// }

	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address())
	code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610251806100206000396000f300608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063f1835db714610046575b600080fd5b34801561005257600080fd5b5061009e6004803603810190808035600019169060200190929190803560ff169060200190929190803560001916906020019092919080356000191690602001909291905050506100e0565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b600060606000806040805190810160405280601c81526020017f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250925082886040518083805190602001908083835b6020831015156101565780518252602082019150602081019050602083039250610131565b6001836020036101000a03801982511681845116808217855250505050505090500182600019166000191681526020019250505060405180910390209150600182888888604051600081526020016040526040518085600019166000191681526020018460ff1660ff1681526020018360001916600019168152602001826000191660001916815260200194505050505060206040516020810390808403906000865af115801561020b573d6000803e3d6000fd5b5050506020604051035190508093505050509493505050505600a165627a7a72305820d41e468fb32b8e5cd5ce468362016d770f132827d73bf66a87a5ea647eca80b00029")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test ecrecover func and it is OK", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 118765)

		// execute method with signed hello for address "0x7156526fbD7a3C72969B54f64e42c10fbb768C8a"
		// web3.sha3("hello")
		input := ethcmn.FromHex("0xf1835db71c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8000000000000000000000000000000000000000000000000000000000000001c9242685bf161793cc25603c231bc2f568eb630ea16aa137d2664ac80388256084f8ae3bd7535248d0bd448298cc2e2071e56992d0774dc340c368ae950852ada")
		data, leftOverGas, err := evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 6688)
		fmt.Println("Data: ", data)
		signer := ethcmn.BytesToAddress(data)
		assert.Equal(t, signer.Hex(), "0x7156526fbD7a3C72969B54f64e42c10fbb768C8a")
	})
}

func TestRunner_fallback(t *testing.T) {
	// pragma solidity ^0.4.0;

	// contract test {

	// 	address public addr;

	// 	function() external payable {
	// 		addr = msg.sender;
	// 	}

	// }

	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address())
	code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610126806100206000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063767800de146081575b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550005b348015608c57600080fd5b50609360d5565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820a2500fe275ffd84cf2dab4522f0a3c214a35a4ce81a9f71e1974b2e5f301b73e0029")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test fallback function and it is executed", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 58911)

		// check the addr and it must be zero
		input := ethcmn.FromHex("0x767800de")
		data, leftOverGas, err := evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 2342)
		fmt.Println("Data: ", data)
		addr := ethcmn.BytesToAddress(data)
		assert.Equal(t, addr.Hex(), "0x0000000000000000000000000000000000000000")

		// execute non existing method and apply fallback
		input = ethcmn.FromHex("0x40c10f190000000000000000000000007156526fbd7a3c72969b54f64e42c10fbb768c8a0000000000000000000000000000000000000000000000000000000000000001")
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 20251)

		// check the addr and it must be zero an sender address as fallback was applied
		input = ethcmn.FromHex("0x767800de")
		data, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 342)
		fmt.Println("Data: ", data)
		addr = ethcmn.BytesToAddress(data)
		assert.Equal(t, addr.Hex(), ethcmn.BytesToAddress(acc.Address()).Hex())
	})
}
