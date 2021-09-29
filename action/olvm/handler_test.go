package olvm

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
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
	"github.com/Oneledger/protocol/vm"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmstore "github.com/tendermint/tendermint/store"
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
	gc := storage.NewGasCalculator(100_100_500)
	cs := storage.NewState(storage.NewChainState("balance", db)).WithGas(gc)

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
	ctx.FeePool = fees.NewStore("f", cs)
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
	bs := tmstore.NewBlockStore(db)
	ctx.Logger = log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")
	ctx.StateDB = vm.NewCommitStateDB(
		evm.NewContractStore(storage.NewState(storage.NewChainState("contracts", db)).WithGas(gc)),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)).WithGas(gc),
			ctx.Balances,
			ctx.Currencies,
		),
		ctx.Logger,
	)
	ctx.StateDB.SetBlockStore(bs)
	bhash := ethcmn.BytesToHash(utils.SHA2([]byte("block")))
	ctx.StateDB.SetBlockHash(bhash)
	ctx.State = store.State
	return ctx
}

func generateKeyPair() (keys.Address, *ecdsa.PublicKey, *ecdsa.PrivateKey) {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)
	prikey, err := ecdsa.GenerateKey(ethcrypto.S256(), reader)
	if err != nil {
		panic("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}
	pubkey := prikey.PublicKey
	data := ethcrypto.PubkeyToAddress(pubkey)
	addr := make([]byte, 20)
	copy(addr[:], data.Bytes())

	return addr, &pubkey, prikey
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

func assemblyExecuteData(from keys.Address, to *keys.Address, nonce uint64, value *big.Int, chainID *big.Int, fromPubKey *ecdsa.PublicKey, fromPrikey *ecdsa.PrivateKey, code []byte, gas uint64) action.SignedTx {
	av := &Transaction{
		From:    from,
		Amount:  action.Amount{Currency: "OLT", Value: *balance.NewAmountFromBigInt(value)},
		Data:    code,
		Nonce:   nonce,
		ChainID: chainID,
	}
	if to != nil {
		av.To = to
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(vm.DefaultGasPrice.Int64())},
		Gas:   int64(gas),
	}
	data, _ := av.Marshal()
	rawTx := action.RawTx{
		Type: av.Type(),
		Data: data,
		Fee:  fee,
		Memo: strconv.FormatUint(av.Nonce, 10),
	}

	var big8 = big.NewInt(8)
	chainIdMul := new(big.Int).Mul(chainID, big.NewInt(2))

	var ethTo *ethcmn.Address
	if to != nil {
		ethTo = new(ethcmn.Address)
		*ethTo = ethcmn.BytesToAddress(to.Bytes())
	}
	legacyTx := &ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: fee.Price.Value.BigInt(),
		Gas:      uint64(gas),
		To:       ethTo,
		Value:    value,
		Data:     code,
	}
	tx := ethtypes.NewTx(legacyTx)
	signer := ethtypes.NewEIP155Signer(chainID)
	tx, err := ethtypes.SignTx(tx, signer, fromPrikey)
	if err != nil {
		panic(err)
	}

	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, chainIdMul)
	V.Sub(V, big8)

	pub, err := utils.RecoverPlain(signer.Hash(tx), R, S, V, true)
	if err != nil {
		panic(err)
	}

	compressedPub := ethcrypto.CompressPubkey(pub)
	pubKey, err := keys.GetPublicKeyFromBytes(compressedPub, keys.ETHSECP)
	if err != nil {
		panic(err)
	}

	signature := utils.ToUncompressedSig(R, S, V)
	signed := action.SignedTx{
		RawTx: rawTx,
		Signatures: []action.Signature{
			{
				Signer: pubKey,
				Signed: signature,
			},
		},
	}
	return signed
}

func wrapProcessDeliver(stx *olvmTx, txHash ethcmn.Hash, ctx *action.Context, signedTx action.SignedTx, f func(ctx *action.Context, tx action.RawTx) (bool, action.Response)) (bool, action.Response) {
	bhash := ethcmn.BytesToHash(utils.SHA2([]byte(fmt.Sprintf("block"))))
	ctx.StateDB.Prepare(txHash)
	ok, resp := f(ctx, signedTx.RawTx)
	ctx.StateDB.Finality(resp.Events)
	ctx.StateDB.Reset()
	ctx.StateDB.SetBlockHash(bhash)
	// bs := ctx.StateDB.GetBlockStore()
	// bs.SaveBlock(&types.Block{
	// 	Header: types.Header{
	// 		Version: version.Consensus{
	// 			Block: version.Protocol(ctx.Header.Version.Block),
	// 			App:   version.Protocol(ctx.Header.Version.App),
	// 		},
	// 		ChainID: ctx.Header.ChainID,
	// 		Height:  globalHeight,
	// 		Time:    time.Now(),
	// 	},
	// 	LastCommit: types.NewCommit(globalHeight, 0, types.BlockID{}, make([]types.CommitSig, 0)),
	// }, &types.PartSet{}, &types.Commit{})
	// globalHeight++
	// ctx.Header = &abci.Header{
	// 	Height:  globalHeight,
	// 	Time:    time.Now().AddDate(0, 0, 1),
	// 	ChainID: "test-1",
	// }
	return ok, resp
}

func getNonce(ctx *action.Context, from keys.Address) uint64 {
	return ctx.StateDB.GetNonce(ethcmn.BytesToAddress(from.Bytes()))
}

var etherDecimals = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

func newAcc(ctx *action.Context, from keys.Address, amount int64) *balance.EthAccount {
	currency, _ := ctx.Currencies.GetCurrencyByName("OLT")
	acc := balance.NewEthAccount(from.Bytes(), balance.Coin{
		Currency: currency,
		Amount:   balance.NewAmountFromBigInt(new(big.Int).Mul(big.NewInt(amount), etherDecimals)),
	})
	err := ctx.StateDB.GetAccountKeeper().SetAccount(*acc)
	if err != nil {
		panic(err)
	}
	return acc
}

func legacyAcc(ctx *action.Context, from keys.Address, amount int64) {
	currency, _ := ctx.Currencies.GetCurrencyByName("OLT")
	coin := balance.Coin{
		Currency: currency,
		Amount:   balance.NewAmountFromBigInt(new(big.Int).Mul(big.NewInt(amount), etherDecimals)),
	}
	err := ctx.Balances.SetBalance(from, coin)
	if err != nil {
		panic(err)
	}
}

func getTxStatus(resp action.Response) uint64 {
	for _, evt := range resp.Events {
		for _, attr := range evt.Attributes {
			if bytes.Equal(attr.Key, []byte("tx.status")) {
				status, _ := strconv.Atoi(string(attr.Value))
				return uint64(status)
			}
		}
	}
	return ethtypes.ReceiptStatusSuccessful
}

func getErr(resp action.Response) string {
	for _, evt := range resp.Events {
		for _, attr := range evt.Attributes {
			if bytes.Equal(attr.Key, []byte("tx.error")) {
				return string(attr.Value)
			}
		}
	}
	return ""
}

func getTxLogs(resp action.Response) []*ethtypes.Log {
	var logs []*ethtypes.Log
	for _, evt := range resp.Events {
		for _, attr := range evt.Attributes {
			if bytes.Contains(attr.Key, []byte("tx.logs")) {
				log, err := new(vm.RLPLog).Decode(attr.Value)
				if err == nil {
					logs = append(logs, log)
				}
			}
		}
	}
	return logs
}

func TestRunner_Send(t *testing.T) {
	// generating default data
	ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
	currency, _ := ctx.Currencies.GetCurrencyByName("OLT")

	chainID := utils.HashToBigInt(ctx.Header.ChainID)
	sendTxFee := new(big.Int).Mul(big.NewInt(int64(vm.TxGas)), big.NewInt(vm.DefaultGasPrice.Int64()))

	txHash := ethcmn.BytesToHash(utils.SHA2([]byte("test")))

	t.Run("test send amount from one eoa to another and it is OK", func(t *testing.T) {
		from, fromPubKey, fromPrikey := generateKeyPair()
		to, _, _ := generateKeyPair()

		sender := newAcc(ctx, from, 10000)
		err := ctx.StateDB.GetAccountKeeper().SetAccount(*sender)
		if err != nil {
			panic(err)
		}

		stx := &olvmTx{}

		value := big.NewInt(100)
		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from, &to, nonce, value, chainID, fromPubKey, fromPrikey, make([]byte, 0), vm.TxGas)
		assert.Equal(t, 0, int(nonce))

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, _ = wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)

		keeper := ctx.StateDB.GetAccountKeeper()
		fromAcc, err := keeper.GetAccount(from)
		assert.NoError(t, err)
		toAcc, err := keeper.GetAccount(to)
		assert.NoError(t, err)

		resAmt, _ := sender.Coins.Minus(balance.Coin{
			Currency: currency,
			Amount:   (*balance.Amount)(new(big.Int).Add(value, sendTxFee)),
		})
		assert.Equal(t, resAmt.Amount.BigInt(), fromAcc.Coins.Amount.BigInt())
		assert.Equal(t, value, toAcc.Coins.Amount.BigInt())

		fromNonce := getNonce(ctx, from.Bytes())
		assert.Equal(t, 1, int(fromNonce))

		toNonce := getNonce(ctx, to.Bytes())
		assert.Equal(t, 0, int(toNonce))
	})

	t.Run("test send amount to some eoa and it is OK", func(t *testing.T) {
		from, fromPubKey, fromPrikey := generateKeyPair()
		to, _, _ := generateKeyPair()

		acc := newAcc(ctx, from, 10000)
		err := ctx.StateDB.GetAccountKeeper().SetAccount(*acc)
		if err != nil {
			panic(err)
		}

		stx := &olvmTx{}

		value := big.NewInt(100)
		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from, &to, nonce, value, chainID, fromPubKey, fromPrikey, make([]byte, 0), vm.TxGas)
		assert.Equal(t, 0, int(nonce))

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, _ = wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)

		keeper := ctx.StateDB.GetAccountKeeper()
		fromAcc, _ := keeper.GetAccount(from)
		toAcc, _ := keeper.GetAccount(to)

		resAmt, _ := acc.Coins.Minus(balance.Coin{
			Currency: currency,
			Amount:   (*balance.Amount)(new(big.Int).Add(value, sendTxFee)),
		})
		assert.Equal(t, resAmt.Amount.BigInt(), fromAcc.Coins.Amount.BigInt())
		assert.Equal(t, value, toAcc.Coins.Amount.BigInt())

		fromNonce := getNonce(ctx, from.Bytes())
		assert.Equal(t, 1, int(fromNonce))

		toNonce := getNonce(ctx, to.Bytes())
		assert.Equal(t, 0, int(toNonce))
	})

	// t.Run("test send amount with high nonce to some eoa and it is error", func(t *testing.T) {
	// 	from, fromPubKey, fromPrikey := generateKeyPair()
	// 	to, _, _ := generateKeyPair()
	// 	newAcc(ctx, from, 10000)

	// 	stx := &olvmTx{}

	// 	value := big.NewInt(100)
	// 	nonce := uint64(100_500)
	// 	tx := assemblyExecuteData(from, &to, nonce, value, chainID, fromPubKey, fromPrikey, make([]byte, 0), vm.TxGas)

	// 	ok, err := stx.Validate(ctx, tx)
	// 	assert.Error(t, err)
	// 	assert.False(t, ok)
	// })

	t.Run("test send minus amount to some eoa and it is error", func(t *testing.T) {
		from, fromPubKey, fromPrikey := generateKeyPair()
		to, _, _ := generateKeyPair()
		newAcc(ctx, from, 10000)

		stx := &olvmTx{}

		value := big.NewInt(-100)
		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from, &to, nonce, value, chainID, fromPubKey, fromPrikey, make([]byte, 0), vm.TxGas)
		assert.Equal(t, 0, int(nonce))

		ok, err := stx.Validate(ctx, tx)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("test send overflow amount to some eoa and it is error", func(t *testing.T) {
		from, fromPubKey, fromPrikey := generateKeyPair()
		to, _, _ := generateKeyPair()
		newAcc(ctx, from, 10000)

		stx := &olvmTx{}

		value := new(big.Int).Mul(big.NewInt(10001), etherDecimals)
		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from, &to, nonce, value, chainID, fromPubKey, fromPrikey, make([]byte, 0), vm.TxGas)
		assert.Equal(t, 0, int(nonce))

		ok, err := stx.Validate(ctx, tx)
		assert.True(t, strings.Contains(err.Error(), "insufficient funds for gas * price + value"))
		assert.False(t, ok)

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.False(t, ok)

		assert.Equal(t, uint64(0), getTxStatus(resp))
		assert.Equal(t, ctx.StateDB.GetBalance(ethcmn.BytesToAddress(to.Bytes())), big.NewInt(0))
	})
}

func TestRunner_ContractWithSuicide(t *testing.T) {
	// // SPDX-License-Identifier: MIT
	// pragma solidity ^0.8.4;

	// contract Test {

	//     address public constant dead = 0x000000000000000000000000000000000000dEaD;

	//     function kill() public payable {
	//         selfdestruct(payable(dead));
	//     }
	// }
	txHash := ethcmn.BytesToHash(utils.SHA2([]byte("test")))

	t.Run("test contract deployment and after suicide and it is OK", func(t *testing.T) {
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		// legacy acc set of balance
		initBalance := int64(10000)
		legacyAcc(ctx, from, initBalance)

		assert.Equal(t, new(big.Int).Mul(big.NewInt(initBalance), etherDecimals), ctx.StateDB.GetBalance(ethcmn.BytesToAddress(from.Bytes())))

		stx := &olvmTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610106806100206000396000f3fe60806040526004361060265760003560e01c806336cf7c8714602b57806341c0e1b5146051575b600080fd5b348015603657600080fd5b50603d6059565b6040516048919060b7565b60405180910390f35b6057605f565b005b61dead81565b61dead73ffffffffffffffffffffffffffffffffffffffff16ff5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600060a382607a565b9050919050565b60b181609a565b82525050565b600060208201905060ca600083018460aa565b9291505056fea2646970667358221220ad88262232f5428f4bbc9bede4bc2b66388ae8cf3e2b5682cd3cbd8af2ee83cc64736f6c63430008090033")
		fmt.Printf("code to deploy: %s\n", ethcmn.Bytes2Hex(code))

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, code, 232115)
		assert.Equal(t, int(nonce), 0)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)
		totalFeeCost := new(big.Int).Mul(big.NewInt(resp.GasUsed), vm.DefaultGasPrice)

		logs := getTxLogs(resp)
		assert.Equal(t, len(logs), 0, "Logs must be empty as contract only deployed and no event in the constructor")

		senderBalance := ctx.StateDB.GetBalance(ethcmn.BytesToAddress(from.Bytes()))

		to := keys.Address(ethcrypto.CreateAddress(ethcmn.BytesToAddress(from.Bytes()), nonce).Bytes())
		assert.True(t, len(ctx.StateDB.GetCode(ethcmn.BytesToAddress(to.Bytes()))) > 0)
		assert.Equal(t, new(big.Int).Sub(new(big.Int).Mul(big.NewInt(initBalance), etherDecimals), totalFeeCost), senderBalance)

		// going to suicide
		input := ethcmn.FromHex("0x41c0e1b5")
		deadAddress := ethcmn.HexToAddress("0x000000000000000000000000000000000000dEaD")
		value := new(big.Int).Mul(big.NewInt(20), etherDecimals)

		nonce = getNonce(ctx, from.Bytes())
		tx2 := assemblyExecuteData(from.Bytes(), &to, nonce, value, chainID, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 1)

		ok, err = stx.Validate(ctx, tx2)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx2, stx.ProcessDeliver)
		assert.True(t, ok)
		totalFeeCost = new(big.Int).Add(totalFeeCost, new(big.Int).Mul(big.NewInt(resp.GasUsed), vm.DefaultGasPrice))

		deadBalance := ctx.StateDB.GetBalance(deadAddress)
		assert.Equal(t, value, deadBalance)

		senderBalance = ctx.StateDB.GetBalance(ethcmn.BytesToAddress(from.Bytes()))
		assert.Equal(t, new(big.Int).Sub(new(big.Int).Mul(big.NewInt(initBalance), etherDecimals), new(big.Int).Add(totalFeeCost, value)), senderBalance)

		assert.True(t, len(ctx.StateDB.GetCode(ethcmn.BytesToAddress(to.Bytes()))) == 0)
	})
}

func TestRunner_BlockHashContract(t *testing.T) {
	// // SPDX-License-Identifier: MIT
	// pragma solidity ^0.8.7;

	// contract Test {

	//     event BlockHash(bytes32 hash);

	//     function printBlockHashes(uint64 blockNum) public {
	//         emit BlockHash(blockhash(blockNum));
	//     }
	// }

	txHash := ethcmn.BytesToHash(utils.SHA2([]byte("test")))

	t.Run("test contract blockhash and it is OK", func(t *testing.T) {
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		// legacy acc set of balance
		legacyAcc(ctx, from, 10000)

		stx := &olvmTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610159806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063f647ad4114610030575b600080fd5b61004a600480360381019061004591906100c2565b61004c565b005b7f57a862a8a40fa048d03f6073ebf87b1ad72de9f7d6e9c163ee04a9d175ff6cd5814060405161007c9190610108565b60405180910390a150565b600080fd5b6000819050919050565b61009f8161008c565b81146100aa57600080fd5b50565b6000813590506100bc81610096565b92915050565b6000602082840312156100d8576100d7610087565b5b60006100e6848285016100ad565b91505092915050565b6000819050919050565b610102816100ef565b82525050565b600060208201905061011d60008301846100f9565b9291505056fea26469706673582212206d975049cbc4bb739c3463471c74b6dd98692e2b03d22a539a1f0e532f9ecbfe64736f6c63430008090033")
		fmt.Printf("code to deploy: %s\n", ethcmn.Bytes2Hex(code))

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, code, 232115)
		assert.Equal(t, int(nonce), 0)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs := getTxLogs(resp)
		assert.Equal(t, len(logs), 0, "Logs must be empty as contract only deployed and no event in the constructor")

		to := keys.Address(ethcrypto.CreateAddress(ethcmn.BytesToAddress(from.Bytes()), nonce).Bytes())
		assert.True(t, len(ctx.StateDB.GetCode(ethcmn.BytesToAddress(to.Bytes()))) > 0)

		input := ethcmn.FromHex("0xf647ad4100000000000000000000000000000000000000000000000000000000000003e7")

		nonce = getNonce(ctx, from.Bytes())
		tx2 := assemblyExecuteData(from.Bytes(), &to, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 1)

		ok, err = stx.Validate(ctx, tx2)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx2, stx.ProcessDeliver)
		assert.True(t, ok)

		logs = getTxLogs(resp)
		assert.Equal(t, len(logs), 1, "Logs must not be empty as contract function emit BlockHash event")

		assert.Equal(t, ethcmn.BytesToHash(logs[0].Data), ethcmn.Hash{})
	})
}

func TestRunner_BaseSmartContract(t *testing.T) {
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

	txHash := ethcmn.BytesToHash(utils.SHA2([]byte("test")))

	t.Run("test contract deployment when user has legacy balance without ethereum balance and it is OK", func(t *testing.T) {
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		// legacy acc set of balance
		legacyAcc(ctx, from, 10000)

		stx := &olvmTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610233806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635f76f6ab146100465780636d4ce63c14610076578063cbed952214610096575b600080fd5b6100746004803603602081101561005c57600080fd5b810190808035151590602001909291905050506100a0565b005b61007e61013c565b60405180821515815260200191505060405180910390f35b61009e61018f565b005b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167fab77f9000c19702a713e62164a239e3764dde2ba5265c7551f9a49e0d304530d60405160405180910390a250565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905090565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f68656c6c6f00000000000000000000000000000000000000000000000000000081525060200191505060405180910390fdfea26469706673582212206872039b48bb16fb8cbf559a2e127d91b0af06f0d2d36b97faad6d0f9c335e7864736f6c63430007040033")
		fmt.Printf("code to deploy: %s\n", ethcmn.Bytes2Hex(code))

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, code, 232115)
		assert.Equal(t, int(nonce), 0)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs := getTxLogs(resp)
		assert.Equal(t, len(logs), 0, "Logs must be empty as contract only deployed and no event in the constructor")
	})

	t.Run("test contract store through the transaction and it is OK", func(t *testing.T) {
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		newAcc(ctx, from, 10000)

		stx := &olvmTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610233806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635f76f6ab146100465780636d4ce63c14610076578063cbed952214610096575b600080fd5b6100746004803603602081101561005c57600080fd5b810190808035151590602001909291905050506100a0565b005b61007e61013c565b60405180821515815260200191505060405180910390f35b61009e61018f565b005b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167fab77f9000c19702a713e62164a239e3764dde2ba5265c7551f9a49e0d304530d60405160405180910390a250565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905090565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f68656c6c6f00000000000000000000000000000000000000000000000000000081525060200191505060405180910390fdfea26469706673582212206872039b48bb16fb8cbf559a2e127d91b0af06f0d2d36b97faad6d0f9c335e7864736f6c63430007040033")
		fmt.Printf("code to deploy: %s\n", ethcmn.Bytes2Hex(code))

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, code, 232115)
		assert.Equal(t, int(nonce), 0)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := wrapProcessDeliver(stx, txHash, ctx, tx, stx.ProcessDeliver)
		assert.True(t, ok)

		logs := getTxLogs(resp)
		assert.Equal(t, len(logs), 0, "Logs must be empty as contract only deployed and no event in the constructor")

		to := keys.Address(ethcrypto.CreateAddress(ethcmn.BytesToAddress(from.Bytes()), nonce).Bytes())

		// going to set data
		input := ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000001")

		nonce = getNonce(ctx, from.Bytes())
		tx2 := assemblyExecuteData(from.Bytes(), &to, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 1)

		ok, err = stx.Validate(ctx, tx2)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx2, stx.ProcessDeliver)
		assert.True(t, ok)

		logs = getTxLogs(resp)
		assert.Equal(t, len(logs), 1, "Logs must not be empty as event was emited")

		// and after read it
		input = ethcmn.FromHex("0x6d4ce63c")

		nonce = getNonce(ctx, from.Bytes())
		tx3 := assemblyExecuteData(from.Bytes(), &to, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, input, 132115)
		assert.Equal(t, int(nonce), 2)

		ok, err = stx.Validate(ctx, tx3)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = wrapProcessDeliver(stx, txHash, ctx, tx3, stx.ProcessDeliver)
		assert.True(t, ok)

		logs = getTxLogs(resp)
		assert.Equal(t, len(logs), 0, "Logs must be empty as logs state were cleared")

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
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		newAcc(ctx, from, 10000)

		stx := &olvmTx{}
		code := ethcmn.FromHex("0x608060405234801561001057600080fd5b50610233806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635f76f6ab146100465780636d4ce63c14610076578063cbed952214610096575b600080fd5b6100746004803603602081101561005c57600080fd5b810190808035151590602001909291905050506100a0565b005b61007e61013c565b60405180821515815260200191505060405180910390f35b61009e61018f565b005b806000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167fab77f9000c19702a713e62164a239e3764dde2ba5265c7551f9a49e0d304530d60405160405180910390a250565b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905090565b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260058152602001807f68656c6c6f00000000000000000000000000000000000000000000000000000081525060200191505060405180910390fdfea26469706673582212206872039b48bb16fb8cbf559a2e127d91b0af06f0d2d36b97faad6d0f9c335e7864736f6c63430007040033")

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), nil, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, code, 100)

		ok, err := stx.Validate(ctx, tx)
		assert.True(t, strings.Contains(err.Error(), "intrinsic gas too low"))
		assert.False(t, ok)
	})

	t.Run("test contract func exec on missed address and it is ok", func(t *testing.T) {
		// generating default data
		ctx := assemblyCtxData("OLT", 18, true, false, false, nil)
		chainID := utils.HashToBigInt(ctx.Header.ChainID)
		from, fromPubKey, fromPrikey := generateKeyPair()

		newAcc(ctx, from, 10000)

		stx := &olvmTx{}
		to_, _, _ := generateKeyPair()
		to := keys.Address(to_.Bytes())

		nonce := getNonce(ctx, from.Bytes())
		tx := assemblyExecuteData(from.Bytes(), &to, nonce, big.NewInt(0), chainID, fromPubKey, fromPrikey, ethcmn.FromHex("0x5f76f6ab0000000000000000000000000000000000000000000000000000000000000001"), 100000)

		ok, err := stx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := stx.ProcessDeliver(ctx, tx.RawTx)
		fmt.Printf("resp: %+v \n", resp)
		assert.True(t, ok)
	})
}
