package evidence

import (
	"testing"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func assemblyReleaseData() action.SignedTx {
	av := &Release{
		ValidatorAddress: from.Bytes(),
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10000000000)},
		Gas:   10,
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
				Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
				Signed: signature,
			},
		},
	}
	return signed
}

func TestReleaseTx_ProcessDeliver_OK(t *testing.T) {
	rtx := &releaseTx{}
	ctx := &action.Context{}

	t.Run("release non frozen validator and get error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("OLT", 1000)
		tx := assemblyReleaseData()

		// could not release as validator is not frozen
		ok, err := rtx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := rtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok, resp)
	})

	t.Run("release malicious validator and get ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("OLT", 1000)
		tx := assemblyReleaseData()

		createdAt := time.Now()

		lvh, err := ctx.EvidenceStore.CreateSuspiciousValidator(from.Bytes(), evidence.MISSED_REQUIRED_VOTES, 2, &createdAt)
		assert.NoError(t, err)
		assert.True(t, lvh.IsFrozen())

		ctx.Header = &abci.Header{
			Height: 10,
			Time:   time.Now().AddDate(0, 0, 1),
		}

		ok, err := rtx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := rtx.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		lvh, err = ctx.EvidenceStore.GetSuspiciousValidator(from.Bytes(), 0, 0)
		assert.NoError(t, err)
		assert.False(t, lvh.IsFrozen())
	})

	t.Run("release malicious validator with byzantine fault and get ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ctx = assemblyCtxData("OLT", 1000)
		tx := assemblyReleaseData()

		createdAt := time.Now()

		lvh, err := ctx.EvidenceStore.CreateSuspiciousValidator(from.Bytes(), evidence.BYZANTINE_FAULT, 2, &createdAt)
		assert.NoError(t, err)
		assert.True(t, lvh.IsFrozen())

		ok, err := rtx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := rtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok, resp)

		// update evidence options to test release
		opts, _ := ctx.GovernanceStore.GetEvidenceOptions()
		opts.ValidatorReleaseTime = 0
		ctx.GovernanceStore.SetEvidenceOptions(*opts)

		ctx.Header = &abci.Header{
			Height: 10,
			Time:   time.Now().AddDate(0, 0, 1),
		}

		ok, err = rtx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp = rtx.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		lvh, err = ctx.EvidenceStore.GetSuspiciousValidator(from.Bytes(), 0, 0)
		assert.NoError(t, err)
		assert.False(t, lvh.IsFrozen())
	})
}
