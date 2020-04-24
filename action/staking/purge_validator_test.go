package staking

import (
	"encoding/hex"
	"testing"

	db2 "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/identity"

	"github.com/Oneledger/protocol/storage"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"

	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

// global setup_purgeValidatorStore
func setup_purgeValidatorStore(validator keys.Address) *identity.ValidatorStore {
	db := db2.NewDB("test", db2.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", db))
	vs := identity.NewValidatorStore("v", cs)

	apply := identity.Stake{
		ValidatorAddress: validator,
		StakeAddress:     validator,
		Pubkey: keys.PublicKey{
			KeyType: keys.ED25519,
			Data:    nil,
		},
		Name: "test_node",
	}
	vs.HandleStake(apply)

	return vs
}

func assemblyPurgeValidatorData(t *testing.T, pubkey keys.PublicKey, prikey keys.PrivateKey, admin keys.Address, validator keys.Address) action.SignedTx {
	pg := &Purge{
		AdminAddress:     admin,
		ValidatorAddress: validator,
	}
	fee := action.Fee{
		Price: action.Amount{
			Currency: "OLT",
			Value:    *balance.NewAmount(1000000000),
		},
		Gas: int64(80000),
	}
	data, _ := pg.Marshal()
	tx := action.RawTx{
		Type: pg.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	handler, _ := prikey.GetHandler()
	signature, _ := handler.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			{
				Signer: pubkey,
				Signed: signature,
			}},
	}
	return signed
}

func generatePurgeParam(t *testing.T) (keys.PublicKey, keys.PrivateKey, keys.Address, keys.Address) {
	pubkey, prikey, err := keys.NewKeyPairFromTendermint()
	assert.Nil(t, err)
	handler, err := pubkey.GetHandler()
	assert.Nil(t, err)
	admin := handler.Address()

	pub_validator, _, err := keys.NewKeyPairFromTendermint()
	assert.Nil(t, err)
	handler, err = pub_validator.GetHandler()
	assert.Nil(t, err)
	validator := handler.Address()

	return pubkey, prikey, admin, validator
}

func setupForPurgeValidator(t *testing.T) (*identity.ValidatorStore, action.SignedTx) {
	pubkey, prikey, admin, validator := generatePurgeParam(t)
	vs := setup_purgeValidatorStore(validator)
	tx := assemblyPurgeValidatorData(t, pubkey, prikey, admin, validator)
	return vs, tx
}

func TestPurgeValidator_Signers(t *testing.T) {
	t.Run("run empty Signer test case, shoule return an empty result", func(t *testing.T) {
		purge := Purge{}
		correctResult := []action.Address{nil}
		assert.Equal(t, correctResult, purge.Signers())
	})
	t.Run("run random Signer test case, should return the same result as the input data", func(t *testing.T) {
		hexData, err := hex.DecodeString("89507C7ABC6D1E9124FE94101A0AB38D5085E15A")
		assert.Nil(t, err)
		b := []byte(hexData)
		correctResult := []action.Address{b}
		purge := Purge{}
		purge.AdminAddress = keys.Address(b).Bytes()
		assert.Equal(t, correctResult, purge.Signers())
	})
}

func TestPurgeTx_Validate(t *testing.T) {
	t.Run("test apply validator with valid tx data, should be ok", func(t *testing.T) {
		_, tx := setupForPurgeValidator(t)
		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "OLT",
			Chain:   chain.ONELEDGER,
			Decimal: int64(18),
		}
		err := currencyList.Register(currency)
		assert.Nil(t, err, "register new currency should be ok")

		db := db2.NewDB("test", db2.MemDBBackend, "")
		cs := storage.NewChainState("chainstate", db)
		feePool := fees.NewStore("f", storage.NewState(cs))
		feePool.SetupOpt(&fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9})
		ctx := &action.Context{
			Currencies: currencyList,
			FeePool:    feePool,
		}
		handler := purgeTx{}
		ok, err := handler.Validate(ctx, tx)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestPurgeTx_ProcessCheck(t *testing.T) {
	t.Run("check with valid validator address, should return no error", func(t *testing.T) {
		vs, tx := setupForPurgeValidator(t)

		ctx := &action.Context{
			Validators: vs,
		}

		handler := purgeTx{}
		ok, _ := handler.ProcessCheck(ctx, tx.RawTx)
		assert.True(t, ok)
	})
}

func TestPurgeTx_ProcessDeliver(t *testing.T) {
	t.Run("process with valid data, should return no error", func(t *testing.T) {
		vs, tx := setupForPurgeValidator(t)

		currencyList := balance.NewCurrencySet()
		currency := balance.Currency{
			Name:    "OLT",
			Chain:   chain.ONELEDGER,
			Decimal: int64(18),
		}
		err := currencyList.Register(currency)
		assert.Nil(t, err, "register new currency should be ok")

		ctx := &action.Context{
			Validators: vs,
			Currencies: currencyList,
			FeeOpt:     &fees.FeeOption{FeeCurrency: currency, MinFeeDecimal: 9},
		}
		handler := purgeTx{}
		ok, _ := handler.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok)
	})
}
