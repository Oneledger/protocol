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

func TestRunner_TokenMapWithDeleteFunctionality(t *testing.T) {
	// // SPDX-License-Identifier: MIT
	// pragma solidity ^0.8.4;

	// import "@openzeppelin/contracts/access/Ownable.sol";

	// contract TokenManager is Ownable {
	//     address private _WETH;
	//     address private _rootToken;

	//     struct Token {
	//         address token;
	//         bool isLocked;
	//     }

	//     mapping(address => Token) private _store;

	//     event TokenAdded(address indexed exitToken, address indexed receiveToken);
	//     event TokenRemoved(address indexed exitToken, address indexed receiveToken);
	//     event RootUpdated(address indexed prev, address indexed current);

	//     constructor(address rootToken_, address WETH_) {
	//         _rootToken = rootToken_;
	//         _WETH = WETH_;
	//     }

	//     /**
	//      * @dev This should be responsible to connect tokens between chains
	//      * @param exitToken token address on exit calldata
	//      * @param receiveToken token address on chain receive
	//      */
	//     function issue(address exitToken, address receiveToken) external onlyOwner {
	//         require(!_store[exitToken].isLocked, "TokenManager: EXIST");

	//         Token memory token = Token(receiveToken, true);
	//         _store[exitToken] = token;
	//         emit TokenAdded(exitToken, receiveToken);
	//     }

	//     /**
	//      * @dev This should be responsible to remove tokens connection between chains
	//      * @param exitToken token address on exit calldata
	//      */
	//     function revoke(address exitToken) external onlyOwner {
	//         require(_store[exitToken].isLocked, "TokenManager: NOT_EXIST");

	//         address receiveToken = _store[exitToken].token;
	//         delete _store[exitToken];
	//         emit TokenRemoved(exitToken, receiveToken);
	//     }

	//     /**
	//      * @dev This should be responsible to get token mapping on current chain
	//      * @param exitToken token address on exit calldata
	//      */
	//     function get(address exitToken)
	//         external
	//         view
	//         returns (bool success, address token)
	//     {
	//         Token memory tokenData = _store[exitToken];
	//         return (tokenData.isLocked, tokenData.token);
	//     }

	//     function WETH() external view returns (address) {
	//         return _WETH;
	//     }

	//     function isRootToken(address token) external view returns (bool) {
	//         return _rootToken == token;
	//     }

	//     function isWETH(address token) external view returns (bool) {
	//         require(_WETH != address(0), "TokenManager: ZERO_WETH");
	//         return _WETH == token;
	//     }
	// }

	// generating default data
	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
	code := ethcmn.FromHex("0x60806040523480156200001157600080fd5b5060405162001161380380620011618339818101604052810190620000379190620001c4565b620000576200004b620000e160201b60201c565b620000e960201b60201c565b81600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050506200025e565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600081519050620001be8162000244565b92915050565b60008060408385031215620001de57620001dd6200023f565b5b6000620001ee85828601620001ad565b92505060206200020185828601620001ad565b9150509250929050565b600062000218826200021f565b9050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600080fd5b6200024f816200020b565b81146200025b57600080fd5b50565b610ef3806200026e6000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c80638da5cb5b116100665780638da5cb5b1461011e578063ad5c46481461013c578063b40464b91461015a578063c2bc2efc14610176578063f2fde38b146101a757610093565b806342056b9e14610098578063715018a6146100c857806374a8f103146100d2578063882a1cfd146100ee575b600080fd5b6100b260048036038101906100ad9190610b26565b6101c3565b6040516100bf9190610c7b565b60405180910390f35b6100d06102ae565b005b6100ec60048036038101906100e79190610b26565b610336565b005b61010860048036038101906101039190610b26565b610583565b6040516101159190610c7b565b60405180910390f35b6101266105dd565b6040516101339190610c60565b60405180910390f35b610144610606565b6040516101519190610c60565b60405180910390f35b610174600480360381019061016f9190610b53565b610630565b005b610190600480360381019061018b9190610b26565b610877565b60405161019e929190610c96565b60405180910390f35b6101c160048036038101906101bc9190610b26565b61094d565b005b60008073ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415610256576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161024d90610cdf565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16149050919050565b6102b6610a45565b73ffffffffffffffffffffffffffffffffffffffff166102d46105dd565b73ffffffffffffffffffffffffffffffffffffffff161461032a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161032190610d1f565b60405180910390fd5b6103346000610a4d565b565b61033e610a45565b73ffffffffffffffffffffffffffffffffffffffff1661035c6105dd565b73ffffffffffffffffffffffffffffffffffffffff16146103b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103a990610d1f565b60405180910390fd5b600360008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900460ff16610441576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161043890610d3f565b60405180910390fd5b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600080820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556000820160146101000a81549060ff021916905550508073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fbbe55b1ff108e23e5ff1a6f5d36946eec15ec0ca0ded2bfed4cdcf697ca9046060405160405180910390a35050565b60008173ffffffffffffffffffffffffffffffffffffffff16600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16149050919050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610638610a45565b73ffffffffffffffffffffffffffffffffffffffff166106566105dd565b73ffffffffffffffffffffffffffffffffffffffff16146106ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106a390610d1f565b60405180910390fd5b600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160149054906101000a900460ff161561073c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161073390610cff565b60405180910390fd5b600060405180604001604052808373ffffffffffffffffffffffffffffffffffffffff16815260200160011515815250905080600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160000160146101000a81548160ff0219169083151502179055509050508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fdffbd9ded1c09446f09377de547142dcce7dc541c8b0b028142b1eba7026b9e760405160405180910390a3505050565b6000806000600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206040518060400160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016000820160149054906101000a900460ff1615151515815250509050806020015181600001519250925050915091565b610955610a45565b73ffffffffffffffffffffffffffffffffffffffff166109736105dd565b73ffffffffffffffffffffffffffffffffffffffff16146109c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109c090610d1f565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610a39576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a3090610cbf565b60405180910390fd5b610a4281610a4d565b50565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b600081359050610b2081610ea6565b92915050565b600060208284031215610b3c57610b3b610dae565b5b6000610b4a84828501610b11565b91505092915050565b60008060408385031215610b6a57610b69610dae565b5b6000610b7885828601610b11565b9250506020610b8985828601610b11565b9150509250929050565b610b9c81610d70565b82525050565b610bab81610d82565b82525050565b6000610bbe602683610d5f565b9150610bc982610db3565b604082019050919050565b6000610be1601783610d5f565b9150610bec82610e02565b602082019050919050565b6000610c04601383610d5f565b9150610c0f82610e2b565b602082019050919050565b6000610c27602083610d5f565b9150610c3282610e54565b602082019050919050565b6000610c4a601783610d5f565b9150610c5582610e7d565b602082019050919050565b6000602082019050610c756000830184610b93565b92915050565b6000602082019050610c906000830184610ba2565b92915050565b6000604082019050610cab6000830185610ba2565b610cb86020830184610b93565b9392505050565b60006020820190508181036000830152610cd881610bb1565b9050919050565b60006020820190508181036000830152610cf881610bd4565b9050919050565b60006020820190508181036000830152610d1881610bf7565b9050919050565b60006020820190508181036000830152610d3881610c1a565b9050919050565b60006020820190508181036000830152610d5881610c3d565b9050919050565b600082825260208201905092915050565b6000610d7b82610d8e565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600080fd5b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b7f546f6b656e4d616e616765723a205a45524f5f57455448000000000000000000600082015250565b7f546f6b656e4d616e616765723a20455849535400000000000000000000000000600082015250565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b7f546f6b656e4d616e616765723a204e4f545f4558495354000000000000000000600082015250565b610eaf81610d70565b8114610eba57600080fd5b5056fea26469706673582212206f936a06b093c9c4eab8145df2547df00468efc15cadcdf624f372831c648a3664736f6c6343000806003300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test contract issue and revoke func and it is OK", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 795262)

		// issuing a new token
		input := ethcmn.FromHex("0xb40464b900000000000000000000000097a4571c88563a0e934cae884f4b1ff0aedee8240000000000000000000000000000000000000000000000000000000000000000")
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 25506)

		// revoking a new token
		input = ethcmn.FromHex("0x74a8f10300000000000000000000000097a4571c88563a0e934cae884f4b1ff0aedee824")
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 3475)

		// must fail as already deleted
		input = ethcmn.FromHex("0x74a8f10300000000000000000000000097a4571c88563a0e934cae884f4b1ff0aedee824")
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.Error(t, err, "Execute contract not failed")
	})
}

func TestRunner_IncrementFromSafeMath(t *testing.T) {
	// // SPDX-License-Identifier: MIT
	// pragma solidity ^0.8.4;

	// import "@openzeppelin/contracts/utils/math/SafeMath.sol";

	// contract Incrementor {
	//     using SafeMath for uint256;

	//     uint256 public incr;

	//     function increment() public {
	//         incr = incr.add(1);
	//     }
	// }

	// generating default data
	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
	code := ethcmn.FromHex("0x608060405234801561001057600080fd5b5061018c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063119fbbd41461003b578063d09de08a14610059575b600080fd5b610043610063565b60405161005091906100ac565b60405180910390f35b610061610069565b005b60005481565b61007f600160005461008790919063ffffffff16565b600081905550565b6000818361009591906100c7565b905092915050565b6100a68161011d565b82525050565b60006020820190506100c1600083018461009d565b92915050565b60006100d28261011d565b91506100dd8361011d565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561011257610111610127565b5b828201905092915050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea2646970667358221220206bd15026183641523bfcf445666088e5c5960b19f01272714e3ac8211f4f2264736f6c63430008060033")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test contract increment func and it is OK", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 79329)

		// execute method with true
		input := ethcmn.FromHex("0xd09de08a") // "increment()"
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 22515)
	})
}

func TestRunner_ComplexIncrementFromSafeMath(t *testing.T) {
	// // SPDX-License-Identifier: MIT
	// pragma solidity ^0.8.4;

	// import "@openzeppelin/contracts/utils/math/SafeMath.sol";

	// abstract contract AIncrementorStorage {
	//     mapping(address => uint256) internal _nonces;
	//     uint256 internal _chainId;

	//     function increment() public virtual;
	// }

	// contract IncrementorAddress is AIncrementorStorage {
	//     using SafeMath for uint256;

	//     constructor() {
	//         assembly {
	//             sstore(_chainId.slot, chainid())
	//         }
	//     }

	//     event Enter(
	//         address indexed token,
	//         address indexed exitor,
	//         address indexed enterer,
	//         uint256 amount,
	//         uint256 nonce,
	//         uint256 chainId
	//     );

	//     function emitEnter(
	//         address token,
	//         address from,
	//         address to,
	//         uint256 amount
	//     ) internal {
	//         _nonces[to] = _nonces[to].add(1);
	//         emit Enter(token, from, to, amount, _nonces[to], _chainId);
	//     }

	//     function increment() public override {
	//         emitEnter(address(0), address(0), msg.sender, 1);
	//     }
	// }

	// generating default data
	acc, _ := accounts.GenerateNewAccount(1, "test")
	ctx := assemblyCtxData(acc, "OLT", 18, true, false, false, nil)
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
	code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50466001556102c5806100246000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063d09de08a14610030575b600080fd5b61003861003a565b005b61004860008033600161004a565b565b61009c60016000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546101a490919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167ffc56b948fab2c51225cc7a8775e3e87c97b66a507945dc738ee7a0b3c7425d4b846000808873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054600154604051610196939291906101c9565b60405180910390a450505050565b600081836101b29190610200565b905092915050565b6101c381610256565b82525050565b60006060820190506101de60008301866101ba565b6101eb60208301856101ba565b6101f860408301846101ba565b949350505050565b600061020b82610256565b915061021683610256565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561024b5761024a610260565b5b828201905092915050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea264697066735822122035229669a40316eee6fc1cafa42d308a9a000bfc1dfc12def31b96e792853ddf64736f6c63430008060033")
	gas := uint64(3000000)
	value := big.NewInt(0)

	accRef := ethvm.AccountRef(ethcmn.BytesToAddress(acc.Address()))
	nonce := getNonce(ctx, acc.Address())
	evm := NewEVMTransaction(ctx.StateDB, ctx.Header, acc.Address(), nil, nonce, value, code, true).NewEVM()

	t.Run("test contract increment func and it is OK", func(t *testing.T) {
		// deploy a contract
		_, contractAddr, leftOverGas, err := evm.Create(accRef, code, gas, value)
		assert.NoError(t, err, "Set contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 164095)

		// // execute method with true
		input := ethcmn.FromHex("0xd09de08a") // "increment()"
		_, leftOverGas, err = evm.Call(accRef, contractAddr, input, gas, value)
		assert.NoError(t, err, "Execute contract failed")
		assertExcError(t, contractAddr, gas, leftOverGas, 26067)
	})
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
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
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
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
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
	ctx.StateDB.GetAccountKeeper().NewAccountWithAddress(acc.Address(), true)
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
