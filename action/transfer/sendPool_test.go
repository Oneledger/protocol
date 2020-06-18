package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/tendermint/tendermint/crypto"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

func assemblySendPoolData(replaceFrom bool) (action.SignedTx, crypto.Address) {

	from, fromPubkey, fromPrikey := generateKeyPair()
	if replaceFrom {
		from = crypto.Address("")
	}
	amt, _ := balance.NewAmountFromString("100000000000000000000", 10)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}
	sendPool := &SendPool{
		From:     from.Bytes(),
		PoolName: "BountyProgram",
		Amount:   *amount,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(1000000000)},
		Gas:   int64(10),
	}

	data, _ := sendPool.Marshal()
	tx := action.RawTx{
		Type: sendPool.Type(),
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

func TestSendPoolTx_Validate(t *testing.T) {
	sendPooltx := &sendPoolTx{}
	ctx := &action.Context{}
	t.Run("validate basic send tx with invalid tx data, should return error", func(t *testing.T) {
		// invalid from address
		ctx = assemblyCtxData("OLT", 18, false, true, false, nil)
		tx, _ := assemblySendPoolData(true)
		ok, err := sendPooltx.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("invalid amount tx data(invalid currency), should return error", func(t *testing.T) {
		tx, _ := assemblySendPoolData(false)
		ctx = assemblyCtxData("OTT", 18, false, true, false, nil)
		ok, err := sendPooltx.Validate(ctx, tx)
		assert.False(t, ok)
		assert.NotNil(t, err)
	})
	t.Run("valid tx data, should return ok", func(t *testing.T) {
		tx, from := assemblySendPoolData(false)
		ctx = assemblyCtxData("OLT", 18, false, true, false, from)
		ok, err := sendPooltx.Validate(ctx, tx)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestSendPoolTx_ProcessCheckTx(t *testing.T) {
	sendPooltx := &sendPoolTx{}
	ctx := &action.Context{}
	t.Run("TX with invalid data, should return error", func(t *testing.T) {
		tx, _ := assemblySendPoolData(false)
		tx.Data = nil
		ctx = assemblyCtxData("", 0, true, true, false, nil)
		ok, _ := sendPooltx.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("From Account Does not have balance ,Should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendPoolData(false)
		ctx = assemblyCtxData("OLT", 18, true, true, false, from)
		ok, _ := sendPooltx.ProcessCheck(ctx, tx.RawTx)
		assert.False(t, ok)
	})
	t.Run("Successfull TX", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		tx, from := assemblySendPoolData(false)
		ctx = assemblyCtxData("OLT", 18, true, true, true, from)
		ok, _ := sendPooltx.ProcessCheck(ctx, tx.RawTx)
		assert.True(t, ok)
	})
}
