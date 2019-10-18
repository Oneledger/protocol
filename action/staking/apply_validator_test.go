package staking

import (
	"encoding/hex"
	"errors"
	"github.com/Oneledger/protocol/data/fees"
	"os"
	"testing"

	db2 "github.com/tendermint/tendermint/libs/db"

	"github.com/Oneledger/protocol/serialize"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/identity"

	"github.com/Oneledger/protocol/storage"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"

	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

// global setup
func setup() (*identity.ValidatorStore, *balance.Store) {
	db := db2.NewDB("test", db2.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	b := balance.NewStore("tb", cs)
	v := identity.NewValidatorStore("tv", config.Server{}, cs)
	return v, b
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

func assemblyApplyValidatorData(addr crypto.Address, pubkey crypto.PubKey, prikey ed25519.PrivKeyEd25519, stake action.Amount, purge bool) action.SignedTx {

	av := &ApplyValidator{
		StakeAddress:     addr.Bytes(),
		Stake:            stake,
		NodeName:         "test_node",
		ValidatorAddress: addr.Bytes(),
		ValidatorPubKey:  keys.PublicKey{keys.ED25519, pubkey.Bytes()[5:]},
		Purge:            purge,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10)},
		Gas:   int64(10),
	}
	data, _ := av.Marshal()
	tx := action.RawTx{
		Type: av.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := prikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			action.Signature{
				Signer: keys.PublicKey{keys.ED25519, pubkey.Bytes()[5:]},
				Signed: signature,
			}},
	}
	return signed
}

func setupForApplyValidatorWithPurge() action.SignedTx {

	addr, pubkey, prikey := generateKeyPair()

	stake := action.Amount{"VT", *balance.NewAmount(1)}
	tx := assemblyApplyValidatorData(addr, pubkey, prikey, stake, true)

	return tx
}

func setupForApplyValidator() action.SignedTx {

	addr, pubkey, prikey := generateKeyPair()

	stake := action.Amount{"VT", *balance.NewAmount(1)}
	tx := assemblyApplyValidatorData(addr, pubkey, prikey, stake, false)

	return tx
}

func setupForInvalidStakeAmount() action.SignedTx {

	addr, pubkey, prikey := generateKeyPair()

	stake := action.Amount{"OLT", *balance.NewAmount(-10)}
	tx := assemblyApplyValidatorData(addr, pubkey, prikey, stake, false)

	return tx
}

func setupForTypeCastApplyValidator() action.SignedTx {

	addr, pubkey, prikey := generateKeyPair()

	// test cast for type casting error
	av := ApplyValidator{
		StakeAddress:     addr.Bytes(),
		Stake:            action.Amount{"OLT", *balance.NewAmount(10)},
		NodeName:         "test_node",
		ValidatorAddress: addr.Bytes(),
		ValidatorPubKey:  keys.PublicKey{keys.ED25519, pubkey.Bytes()[5:]},
		Purge:            false,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(1)},
		Gas:   int64(10),
	}

	data, _ := av.Marshal()
	tx := action.RawTx{
		Type: av.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := prikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			action.Signature{
				Signer: keys.PublicKey{keys.ED25519, pubkey.Bytes()[5:]},
				Signed: signature,
			}},
	}
	return signed
}

func TestApplyValidator_Signers(t *testing.T) {
	t.Run("run empty Signer test case, shoule return an empty result", func(t *testing.T) {
		apply := ApplyValidator{}
		correctResult := []action.Address{nil}
		assert.Equal(t, correctResult, apply.Signers())
	})
	t.Run("run random Signer test case, should return the same result as the input data", func(t *testing.T) {
		hexData, err := hex.DecodeString("89507C7ABC6D1E9124FE94101A0AB38D5085E15A")
		assert.Nil(t, err)
		b := []byte(hexData)
		correctResult := []action.Address{b}
		apply := ApplyValidator{}
		apply.StakeAddress = keys.Address(b).Bytes()
		assert.Equal(t, correctResult, apply.Signers())
	})
}

func TestApplyTx_Validate(t *testing.T) {
	t.Run("test validator basic part, should return no error but validator basic check should be false", func(t *testing.T) {
		tx := setupForApplyValidator()

		signer := keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    []byte(""),
		}
		signature := action.Signature{
			Signer: signer,
			Signed: []byte(""),
		}
		signatures := []action.Signature{signature}
		txNoSig := tx
		txNoSig.Signatures = signatures
		ctx := &action.Context{}
		handler := applyTx{}
		ok, err := handler.Validate(ctx, txNoSig)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})

	//t.Run("test apply validator type cast part, should be not ok", func(t *testing.T) {
	//	tx := setupForTypeCastApplyValidator()
	//	ctx := &action.Context{}
	//	handler := applyTx{}
	//	ok, err := handler.Validate(ctx, tx)
	//	assert.False(t, ok)
	//	assert.NotNil(t, err)
	//})

	t.Run("test apply validator with an invalid currency type, should be not ok", func(t *testing.T) {
		tx := setupForInvalidStakeAmount()

		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "OLT",
			Chain:   chain.Type(1),
			Decimal: int64(18),
		}
		err := currencyList.Register(currency)
		assert.Nil(t, err, "register new OLT token should be ok")
		ctx := &action.Context{
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := applyTx{}
		ok, err := handler.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})

	t.Run("test apply validator with an invalid currency type, should be not ok", func(t *testing.T) {
		tx := setupForApplyValidator()

		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "OLT",
			Chain:   chain.Type(1),
			Decimal: int64(18),
		}
		err := currencyList.Register(currency)
		assert.Nil(t, err, "register new OLT token should be ok")
		ctx := &action.Context{
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := applyTx{}
		ok, err := handler.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})

	t.Run("test apply validator with valid tx data, should be ok", func(t *testing.T) {
		tx := setupForApplyValidator()
		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "VT",
			Chain:   chain.Type(1),
			Decimal: int64(18),
		}
		err := currencyList.Register(currency)
		assert.Nil(t, err, "register new currency should be ok")
		ctx := &action.Context{
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := applyTx{}
		ok, err := handler.Validate(ctx, tx)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestApplyTx_ProcessCheck(t *testing.T) {
	//t.Run("cast applyvalidator tx with invalid data, should return error", func(t *testing.T) {
	//	tx := setupForTypeCastApplyValidator()
	//	ctx := &action.Context{}
	//	handler := applyTx{}
	//	ok, _ := handler.ProcessCheck(ctx, tx.RawTx)
	//	assert.False(t, ok)
	//})
	t.Run("check balance with valid data, should return no error", func(t *testing.T) {
		vs, bs := setup()
		tx := setupForApplyValidator()
		_ = vs

		ctx := &action.Context{
			Balances: bs,
		}
		handler := applyTx{}

		ok, _ := handler.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
}

func TestApplyTx_ProcessDeliver(t *testing.T) {
	t.Run(" process with valid data, should return no error", func(t *testing.T) {
		vs, bs := setup()

		tx := setupForApplyValidator()

		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "VT",
			Chain:   chain.Type(1),
			Decimal: int64(18),
		}
		coin := balance.Coin{
			Currency: currency,
			Amount:   balance.NewAmount(100000000000000000),
		}

		apply := &ApplyValidator{}
		err := serialize.GetSerializer(serialize.JSON).Deserialize(tx.Data, apply)

		err = bs.AddToAddress(apply.StakeAddress, coin)

		assert.Nil(t, err, "set some VT token for test, should be ok")
		err = currencyList.Register(currency)
		assert.Nil(t, err, "register new currency should be ok")
		ctx := &action.Context{
			Balances:   bs,
			Validators: vs,
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := applyTx{}
		ok, _ := handler.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok)
	})
	t.Run("process with valid data but using purge flag, should return no error", func(t *testing.T) {
		vs, bs := setup()
		tx := setupForApplyValidatorWithPurge()
		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "VT",
			Chain:   chain.Type(1),
			Decimal: int64(18),
		}
		coin := balance.Coin{
			Currency: currency,
			Amount:   balance.NewAmount(100000000000000000),
		}

		apply := &ApplyValidator{}
		err := serialize.GetSerializer(serialize.JSON).Deserialize(tx.Data, apply)

		err = bs.AddToAddress(apply.StakeAddress, coin)
		assert.Nil(t, err, "set some VT token for test, should be ok")
		err = currencyList.Register(currency)
		assert.Nil(t, err, "register new currency should be ok")
		ctx := &action.Context{
			Balances:   bs,
			Validators: vs,
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := applyTx{}
		ok, _ := handler.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok)
	})
}
