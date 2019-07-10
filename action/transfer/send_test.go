package transfer

import (
	"errors"
	"math/big"
	"os"
	"testing"

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

func assemblySendData(replaceFrom bool) (*action.BaseTx, crypto.Address) {

	from, fromPubkey, fromPrikey := generateKeyPair()
	to, _, _ := generateKeyPair()

	if replaceFrom {
		from = crypto.Address("")
	}
	amount := &action.Amount{
		Currency: "OLT",
		Value:    "100",
	}
	send := &Send{
		From:   from.Bytes(),
		To:     to.Bytes(),
		Amount: *amount,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", "10"},
		Gas:   int64(10),
	}
	tx := &action.BaseTx{
		Data: send,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := fromPrikey.Sign(tx.Bytes())
	signer := action.Signature{
		Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		Signed: signature,
	}
	tx.Signatures = []action.Signature{signer}

	return tx, from
}

func assemblyCtxData(currencyName string, currencyDecimal int64, setStore bool, setLogger bool, setCoin bool, setCoinAddr crypto.Address) *action.Context {

	ctx := &action.Context{}

	// store
	var store *balance.Store
	if setStore {
		store = balance.NewStore("test_balances", "test_dbpath", storage.CACHE, storage.PERSISTENT)
		ctx.Balances = store
	}
	// logger
	if setLogger {
		ctx.Logger = new(log.Logger)
	}
	// currencyList
	if currencyName != "" {
		// register new token OTT
		currencyList := balance.NewCurrencyList()
		currency := balance.Currency{
			Name:    currencyName,
			Chain:   chain.Type(1),
			Decimal: currencyDecimal,
		}
		err := currencyList.Register(currency)
		if err != nil {
			errors.New("register new token error")
		}
		ctx.Currencies = currencyList

		// set coin for account
		if setCoin {
			coin := balance.Coin{
				Currency: currency,
				Amount:   big.NewInt(100000000000000000),
			}
			ba := balance.NewBalance()
			ba = ba.AddCoin(coin)
			err = store.Set(setCoinAddr.Bytes(), *ba)
			if err != nil {
				errors.New("setup testing token balance error")
			}
			store.Commit()
			ctx.Balances = store
		}
	}
	return ctx
}

func TestSendTx_Validate(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	t.Run("validate basic send tx with invalid tx data, should return error", func(t *testing.T) {
		// invalid from address
		tx, _ := assemblySendData(true)
		ok, err := sendtx.Validate(ctx, tx.Data, tx.Fee, tx.Memo, tx.Signatures)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("cast tx to sendTx with invalid tx data, should return error", func(t *testing.T) {
		amount := action.Amount{
			Currency: "OLT",
			Value:    "100",
		}
		from, fromPubkey, fromPrikey := generateKeyPair()
		to, _, _ := generateKeyPair()
		// invalid send obj
		send := Send{
			From:   from.Bytes(),
			To:     to.Bytes(),
			Amount: amount,
		}
		fee := action.Fee{
			Price: action.Amount{"OLT", "10"},
			Gas:   int64(10),
		}
		tx := &action.BaseTx{
			Data: send,
			Fee:  fee,
			Memo: "test_memo",
		}
		signature, _ := fromPrikey.Sign(tx.Bytes())
		signer := action.Signature{
			Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
			Signed: signature,
		}
		ok, err := sendtx.Validate(ctx, send, tx.Fee, tx.Memo, []action.Signature{signer})
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("invalid amount tx data(invalid currency), should return error", func(t *testing.T) {
		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("OTT", int64(18), false, false, false, nil)
		ok, err := sendtx.Validate(ctx, tx.Data, tx.Fee, tx.Memo, tx.Signatures)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("valid tx data, should return ok", func(t *testing.T) {
		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("OLT", int64(18), false, false, false, nil)
		ok, err := sendtx.Validate(ctx, tx.Data, tx.Fee, tx.Memo, tx.Signatures)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestSendTx_ProcessCheck(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	t.Run("cast sendtx with invalid data, should return error", func(t *testing.T) {
		fee := action.Fee{
			Price: action.Amount{"OLT", "10"},
			Gas:   int64(10),
		}
		tx := &action.BaseTx{
			Data: nil,
		}

		ctx = &action.Context{
			Logger: new(log.Logger),
		}
		ok, _ := sendtx.ProcessCheck(ctx, tx.Data, fee)
		assert.False(t, ok)
	})
	t.Run("get from address balance with invalid data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, _ := assemblySendData(false)
		ctx = assemblyCtxData("", int64(0), true, true, false, nil)

		ok, _ := sendtx.ProcessCheck(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("invalid amount tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLL", int64(18), true, true, true, from)

		ok, _ := sendtx.ProcessCheck(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("check owner balance with invalid tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", int64(18), true, true, true, from)

		ok, _ := sendtx.ProcessCheck(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("check owner balance with valid test data, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", int64(10), true, true, true, from)

		ok, _ := sendtx.ProcessCheck(ctx, tx.Data, tx.Fee)
		assert.True(t, ok)
	})
}

func TestSendTx_ProcessDeliver(t *testing.T) {
	sendtx := &sendTx{}
	ctx := &action.Context{}
	t.Run("cast sendtx type with invalid tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("", int64(0), true, true, false, nil)
		tx, _ := assemblySendData(false)

		ok, _ := sendtx.ProcessDeliver(ctx, nil, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("get balance from balance store with invalid tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("", int64(0), true, true, false, nil)
		tx, _ := assemblySendData(false)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("invalid amount tx data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLL", int64(18), true, true, true, from)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("check owner balance with invalid test data, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", int64(18), true, true, true, from)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.Data, tx.Fee)
		assert.False(t, ok)
	})
	t.Run("change receiver balance with valid test data, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendData(false)
		ctx = assemblyCtxData("OLT", int64(10), true, true, true, from)

		ok, _ := sendtx.ProcessDeliver(ctx, tx.Data, tx.Fee)
		assert.True(t, ok)
	})
}
