package smart_contract

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

func assemblyCtxData(currencyName string, currencyDecimal int, setStore bool, setLogger bool, setCoin bool, setCoinAddr crypto.Address) *action.Context {
	ctx := &action.Context{}
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	// store
	var store *balance.Store
	if setStore {
		store = balance.NewStore("tb", cs)
		ctx.Balances = store
	}
	// logger
	if setLogger {
		ctx.Logger = log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")
	}
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
	ctx.Logger = log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")
	ctx.StateDB = action.NewCommitStateDB(
		evm.NewContractStore(storage.NewState(storage.NewChainState("contracts", db))),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)),
			ctx.Balances,
			ctx.Currencies,
		),
		ctx.Logger,
	)
	return ctx
}

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
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

func assemblyExecuteData(from keys.Address, to *keys.Address, nonce uint64, fromPubKey crypto.PubKey, fromPrikey ed25519.PrivKeyEd25519, code []byte, gas int64) action.SignedTx {
	av := &Execute{
		From:   from,
		Amount: action.Amount{Currency: "OLT", Value: *balance.NewAmount(0)},
		Data:   code,
		Nonce:  nonce,
	}
	if to != nil {
		av.To = to
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10000000000)},
		Gas:   gas,
	}
	data, _ := av.Marshal()
	tx := action.RawTx{
		Type: av.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := fromPrikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			{
				Signer: keys.PublicKey{keys.ED25519, fromPubKey.Bytes()[5:]},
				Signed: signature,
			},
		},
	}
	return signed
}

func getReturnData(resp action.Response) (data []byte) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.data" {
				data = attr.Value
			}
		}
	}
	return
}

func getTxStatus(resp action.Response) (msgStatus uint64, msgError string) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.status" {
				msgStatus = binary.LittleEndian.Uint64(attr.Value)
			} else if string(attr.Key) == "tx.error" {
				msgError = string(attr.Value)
			}
		}
	}
	return
}

func getContractAddress(resp action.Response) (contractAddress ethcmn.Address) {
	for i := range resp.Events {
		evt := resp.Events[i]
		for j := range evt.Attributes {
			attr := evt.Attributes[j]
			if string(attr.Key) == "tx.contract" {
				contractAddress = ethcmn.BytesToAddress(attr.Value)
			}
		}
	}
	return
}

func blockCommit(ctx *action.Context) {
	// simulate block ender
	ctx.StateDB.Finalise(false)
	ctx.StateDB.UpdateAccounts()
	root, err := ctx.StateDB.Commit(false)
	if err != nil {
		panic(err)
	}
	ctx.StateDB.Reset(root)
}

func wrapProcessDeliver(stx *scExecuteTx, txHash ethcmn.Hash, ctx *action.Context, rawTx action.RawTx, f func(ctx *action.Context, tx action.RawTx) (bool, action.Response)) (bool, action.Response) {
	bhash := ethcmn.BytesToHash(utils.SHA2([]byte("block")))
	ctx.StateDB.SetHeightHash(2, bhash)
	ctx.StateDB.SetBlockHash(bhash)
	ctx.StateDB.Prepare(txHash)
	ok, resp := f(ctx, rawTx)
	blockCommit(ctx)
	return ok, resp
}

func getNonce(ctx *action.Context, from keys.Address) uint64 {
	return ctx.StateDB.GetNonce(ethcmn.BytesToAddress(from.Bytes()))
}

func TestRunner(t *testing.T) {
	// pragma solidity >=0.7.0 <0.8.0;

	// contract Test {

	// 	event TestEvent(address indexed owner);

	// 	mapping(address => bool) private data;

	// 	function set(bool res) public {
	// 		data[msg.sender] = res;
	// 		emit TestEvent(msg.sender);
	// 	}

	// 	function get() public view returns(bool) {
	// 		return data[msg.sender];
	// 	}

	// 	function checkRvt() public pure {
	// 		revert("hello");
	// 	}
	// }

	// generating default data
	ctx := assemblyCtxData("OLT", 18, true, false, false, nil)

	from, fromPubKey, fromPrikey := generateKeyPair()

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
	ctx.StateDB.GetAccountKeeper().SetAccount(*acc)

	t.Run("test contract store through the transaction and it is OK", func(t *testing.T) {
		txHash := ethcmn.BytesToHash(utils.SHA2([]byte("test")))

		stx := &scExecuteTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610233806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635f76f6ab146100465780636d4ce63c14610076578063cbed952214610096575b600080fd5b6100746004803603602081101561005c57600080fd5b810190808035151590602001909291905050506100a0565b005b61007e61013c565b60405180821515815260200191505060405180910390f35b61009e61018f565b005b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167fab77f9000c19702a713e62164a239e3764dde2ba5265c7551f9a49e0d304530d60405160405180910390a250565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905090565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f68656c6c6f00000000000000000000000000000000000000000000000000000081525060200191505060405180910390fdfea26469706673582212206872039b48bb16fb8cbf559a2e127d91b0af06f0d2d36b97faad6d0f9c335e7864736f6c63430007040033")
		fmt.Printf("code to deploy: %s\n", ethcmn.Bytes2Hex(code))

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, fromPubKey, fromPrikey, code, 232115)
		assert.Equal(t, int(nonce), 0)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		logs, _ := ctx.StateDB.GetLogs(txHash)
		assert.Equal(t, len(logs), 0, "Logs must be empty as tx not executed")

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx.RawTx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs, _ = ctx.StateDB.GetLogs(txHash)
		assert.Equal(t, len(logs), 0, "Logs must be empty as contract only deployed and no event in the constructor")

		status, errMsg := getTxStatus(resp)
		contractAddress := getContractAddress(resp)
		to := keys.Address(contractAddress.Bytes())
		assert.True(t, status == ethtypes.ReceiptStatusSuccessful, fmt.Sprintf("Got error: %s", errMsg))

		// going to set data
		input := ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000001")

		nonce = getNonce(ctx, from.Bytes())
		tx2 := assemblyExecuteData(from.Bytes(), &to, nonce, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 1)

		ok, err = stx.Validate(ctx, tx2)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx2.RawTx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs, _ = ctx.StateDB.GetLogs(txHash)
		assert.Equal(t, len(logs), 1, "Logs must not be empty as event was emited")

		status, errMsg = getTxStatus(resp)
		assert.True(t, status == ethtypes.ReceiptStatusSuccessful, fmt.Sprintf("Got error: %s", errMsg))

		// and after read it
		input = ethcmn.FromHex("0x6d4ce63c")

		nonce = getNonce(ctx, from.Bytes())
		tx3 := assemblyExecuteData(from.Bytes(), &to, nonce, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 2)

		ok, err = stx.Validate(ctx, tx3)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx3.RawTx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs, _ = ctx.StateDB.GetLogs(txHash)
		assert.Equal(t, len(logs), 1, "Logs must not be empty as not removed and taken from previous state")

		status, errMsg = getTxStatus(resp)
		data := getReturnData(resp)
		assert.True(t, getBool(data), "Data is not set as 'true'")
		assert.True(t, status == ethtypes.ReceiptStatusSuccessful, fmt.Sprintf("Got error: %s", errMsg))

		nonce = getNonce(ctx, from.Bytes())
		assert.Equal(t, int(nonce), 3)

		// TODO: Add some check of deployed code
		// storageCode := ctx.StateDB.GetCode(contractAddress)
		// fmt.Printf("deployed: %s\n", ethcmn.Bytes2Hex(storageCode))
		// res := bytes.Compare(code, storageCode)
		// fmt.Printf("resp: %d\n", res)
		// assert.True(t, res == 0, "Wrong code deployed")
	})

	t.Run("test contract store through the transaction with not enough gas and it is error", func(t *testing.T) {
		stx := &scExecuteTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610233806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635f76f6ab146100465780636d4ce63c14610076578063cbed952214610096575b600080fd5b6100746004803603602081101561005c57600080fd5b810190808035151590602001909291905050506100a0565b005b61007e61013c565b60405180821515815260200191505060405180910390f35b61009e61018f565b005b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167fab77f9000c19702a713e62164a239e3764dde2ba5265c7551f9a49e0d304530d60405160405180910390a250565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905090565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f68656c6c6f00000000000000000000000000000000000000000000000000000081525060200191505060405180910390fdfea26469706673582212206872039b48bb16fb8cbf559a2e127d91b0af06f0d2d36b97faad6d0f9c335e7864736f6c63430007040033")

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, fromPubKey, fromPrikey, code, 100)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := stx.ProcessDeliver(ctx, tx.RawTx)
		fmt.Printf("resp: %+v \n", resp)
		assert.False(t, ok)
		assert.True(t, strings.Contains(resp.Log, "300103")) // intrinsic gas too low code

		_, errMsg := getTxStatus(resp)
		assert.True(t, len(errMsg) == 0)
	})

	t.Run("test contract func exec on missed address and it is ok", func(t *testing.T) {
		stx := &scExecuteTx{}
		to_, _, _ := generateKeyPair()
		to := keys.Address(to_.Bytes())

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), &to, nonce, fromPubKey, fromPrikey, ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000001"), 100000)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := stx.ProcessDeliver(ctx, tx.RawTx)
		fmt.Printf("resp: %+v \n", resp)
		assert.True(t, ok)

		_, errMsg := getTxStatus(resp)
		assert.True(t, len(errMsg) == 0)
	})
}
