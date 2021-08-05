package legacy

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/Oneledger/protocol/action"
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
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	db "github.com/tendermint/tm-db"
)

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
		balance.NewAccountMapper(
			storage.NewState(storage.NewChainState("mapper", db)),
		),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)),
		),
		ctx.Logger,
	)
	ctx.StateDB.SetHeightHash(uint64(ctx.Header.Height), ethcmn.Hash{}, true)
	return ctx
}

func generateOLTKeyPair() (keys.Address, crypto.PubKey, *ed25519.PrivKeyEd25519) {
	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr.Bytes(), pubkey, &prikey
}

func generateETHKeyPair() (keys.Address, *ecdsa.PublicKey, *ecdsa.PrivateKey) {
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

func assemblyConnectData(from keys.Address, fromPubKey crypto.PubKey, fromPrivKey *ed25519.PrivKeyEd25519, to keys.Address, toPubKey *ecdsa.PublicKey, toPrivKey *ecdsa.PrivateKey, toggle bool, nonce []byte) action.SignedTx {
	legacy := &LegacyConnect{
		LegacyAddress: from,
		NewAddress:    to,
		Toggle:        toggle,
		Nonce:         ethcrypto.Keccak256(nonce),
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10000000000)},
		Gas:   int64(10),
	}
	data, _ := legacy.Marshal()
	rawTx := action.RawTx{
		Type: legacy.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}

	legacySignature, err := fromPrivKey.Sign(legacy.Nonce)
	if err != nil {
		panic(err)
	}

	signatures := make([]action.Signature, 0)
	signatures = append(signatures, action.Signature{
		Signer: keys.PublicKey{keys.ED25519, fromPubKey.Bytes()[5:]},
		Signed: legacySignature,
	})
	if toggle == true {
		newSignature, err := ethcrypto.Sign(legacy.Nonce, toPrivKey)
		if err != nil {
			panic(err)
		}
		compressedPub := ethcrypto.CompressPubkey(toPubKey)
		pubKey, err := keys.GetPublicKeyFromBytes(compressedPub, keys.ETHSECP)
		if err != nil {
			panic(err)
		}
		signatures = append(signatures, action.Signature{
			Signer: pubKey,
			Signed: newSignature,
		})
	}

	signed := action.SignedTx{
		RawTx:      rawTx,
		Signatures: signatures,
	}
	return signed
}

func TestRunner(t *testing.T) {

	t.Run("test link between 0lt account and eth 0x account and it is OK", func(t *testing.T) {
		from, fromPubKey, fromPrivKey := generateOLTKeyPair()
		to, toPubKey, toPrivKey := generateETHKeyPair()

		ctx := assemblyCtxData("OLT", 18, true, false, false, from.Bytes())
		mapper := ctx.StateDB.GetAccountMapper()

		// should not exists right now
		_, err := mapper.Get(from, keys.ED25519)
		assert.Error(t, err)
		_, err = mapper.Get(from, keys.ETHSECP)
		assert.Error(t, err)

		nonceUUID, _ := uuid.NewUUID()
		nonce := []byte(nonceUUID.String())

		tx := assemblyConnectData(from, fromPubKey, fromPrivKey, to, toPubKey, toPrivKey, true, nonce)

		ltx := &legacyConnectTx{}
		ok, err := ltx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, _ = ltx.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok)

		// should exists after execution
		fromAj, err := mapper.Get(from, keys.ED25519)
		assert.NoError(t, err)
		toAj, err := mapper.Get(from, keys.ETHSECP)
		assert.NoError(t, err)

		// verify correct algos and stored data
		assert.Equal(t, from, fromAj.Legacy.Address)
		assert.Equal(t, keys.ED25519, fromAj.Legacy.Algorithm)
		assert.Equal(t, to, fromAj.New.Address)
		assert.Equal(t, keys.ETHSECP, fromAj.New.Algorithm)

		assert.Equal(t, from, toAj.Legacy.Address)
		assert.Equal(t, keys.ED25519, toAj.Legacy.Algorithm)
		assert.Equal(t, to, toAj.New.Address)
		assert.Equal(t, keys.ETHSECP, toAj.New.Algorithm)
	})
}
