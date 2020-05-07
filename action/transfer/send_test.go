package transfer

import (
	"errors"
	"os"
	"testing"

	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"

	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// global setup
func setup() []string {
	testDBs := []string{"test_dbpath"}
	return testDBs
}

// remove test db dir after
func teardown(dbPaths []string) {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
}

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
}

func assemblySendData(replaceFrom bool) (action.SignedTx, crypto.Address) {

	from, fromPubkey, fromPrikey := generateKeyPair()
	to, _, _ := generateKeyPair()

	if replaceFrom {
		from = crypto.Address("")
	}
	amt, _ := balance.NewAmountFromString("100000000000000000000", 10)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}
	send := &Send{
		From:   from.Bytes(),
		To:     to.Bytes(),
		Amount: *amount,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(1000000000)},
		Gas:   int64(10),
	}

	data, _ := send.Marshal()
	tx := action.RawTx{
		Type: send.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := fromPrikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			{
				Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
				Signed: signature,
			}},
	}

	return signed, from
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
		ctx.Logger = new(log.Logger)
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
		amt, err := balance.NewAmountFromString("100", 10)
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
	ctx.Govern = governance.NewStore("tg", cs)
	return ctx
}

func TestSendTx_Validate(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	t.Run("validate basic send tx with invalid tx data, should return error", func(t *testing.T) {
		// invalid from address
		tx, _ := assemblySendData(true)
		ok, err := sendtx.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	//t.Run("cast tx to sendTx with invalid tx data, should return error", func(t *testing.T) {
	//	amount := action.Amount{
	//		Currency: "OLT",
	//		Value:    *balance.NewAmount(100),
	//	}
	//	from, fromPubkey, fromPrikey := generateKeyPair()
	//	to, _, _ := generateKeyPair()
	//	// invalid send obj
	//	send := Send{
	//		From:   from.Bytes(),
	//		To:     to.Bytes(),
	//		Amount: amount,
	//	}
	//	fee := action.Fee{
	//		Price: action.Amount{"OLT", *balance.NewAmount(10)},
	//		Gas:   int64(10),
	//	}
	//	data, _ := send.Marshal()
	//	tx := action.RawTx{
	//		Type: send.Type(),
	//		Data: data,
	//		Fee:  fee,
	//		Memo: "test_memo",
	//	}
	//	signature, _ := fromPrikey.Sign(tx.RawBytes())
	//	signer := action.Signature{
	//		Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
	//		Signed: signature,
	//	}
	//	signed := action.SignedTx{
	//		RawTx:      tx,
	//		Signatures: []action.Signature{signer},
	//	}
	//	ok, err := sendtx.Validate(ctx, signed)
	//	assert.False(t, ok)
	//	assert.NotNil(t, err)
	//})
	t.Run("invalid amount tx data(invalid currency), should return error", func(t *testing.T) {
		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("OTT", 18, false, false, false, nil)
		ok, err := sendtx.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("valid tx data, should return ok", func(t *testing.T) {
		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("OLT", 18, false, false, false, nil)
		ok, err := sendtx.Validate(ctx, tx)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestSendTx_ProcessCheck(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	t.Run("cast sendtx with invalid data, should return error", func(t *testing.T) {
		tx := action.RawTx{
			Data: nil,
		}

		ctx = &action.Context{
			Logger: new(log.Logger),
		}
		ok, _ := sendtx.ProcessCheck(ctx, tx)
		assert.False(t, ok)
	})
	t.Run("get from address balance with invalid data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("", 0, true, true, false, nil)

		ok, _ := sendtx.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("invalid amount tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLL", 18, true, true, true, from)

		ok, _ := sendtx.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("check owner balance with invalid tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", 18, true, true, true, from)

		ok, _ := sendtx.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	//t.Run("check owner balance with valid test data, should return ok", func(t *testing.T) {
	//	testDB := setup()
	//	defer teardown(testDB)
	//
	//	tx, from := assemblySendData(false)
	//	ctx = assemblyCtxData("OLT", 10, true, true, true, from)
	//
	//	ok, _ := sendtx.ProcessCheck(ctx, tx.RawTx)
	//	assert.True(t, ok)
	//})
}

func TestSendTx_ProcessDeliver(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	//t.Run("cast sendtx type with invalid tx data, should return error", func(t *testing.T) {
	//	testDB := setup()
	//	defer teardown(testDB)
	//
	//	ctx = assemblyCtxData("", int64(0), true, true, false, nil)
	//	tx, _ := assemblySendData(false)
	//	tx.Data = []byte{}
	//	ok, _ := sendtx.ProcessDeliver(ctx, tx.RawTx)
	//	assert.False(t, ok)
	//})
	t.Run("get balance from balance store with invalid tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("", 0, true, true, false, nil)
		tx, _ := assemblySendData(false)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("invalid amount tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLL", 18, true, true, true, from)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("check owner balance with invalid test data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", 18, true, true, true, from)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	//t.Run("change receiver balance with valid test data, should return ok", func(t *testing.T) {
	//	testDB := setup()
	//	defer teardown(testDB)
	//
	//	tx, from := assemblySendData(false)
	//	ctx = assemblyCtxData("OLT", 21, true, true, true, from)
	//
	//	ok, _ := sendtx.ProcessDeliver(ctx, tx.RawTx)
	//	assert.True(t, ok)
	//})
}
