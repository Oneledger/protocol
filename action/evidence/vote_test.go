package evidence

import (
	"testing"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

func assemblyVoteData(requestID string, choice int8) action.SignedTx {
	av := &AllegationVote{
		RequestID: requestID,
		Address:   from.Bytes(),
		Choice:    choice,
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

func TestVoteTx_ProcessDeliver_OK(t *testing.T) {
	avtx := &allegationVoteTx{}
	ctx := &action.Context{}

	t.Run("vote for allegation request as yes and get ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ID := "test"
		ctx = assemblyCtxData("OLT", 1000)
		tx := assemblyVoteData(ID, evidence.YES)

		// setup base request
		ar := &evidence.AllegationRequest{
			ID:               ID,
			ReporterAddress:  from.Bytes(),
			MaliciousAddress: fromMal.Bytes(),
			BlockHeight:      2,
			ProofMsg:         "test",
			Status:           evidence.VOTING,
			Votes:            make([]*evidence.AllegationVote, 0),
		}
		err := ctx.EvidenceStore.SetAllegationRequest(ar)
		assert.NoError(t, err)
		at := &evidence.AllegationTracker{
			Requests: make(map[string]bool),
		}
		at.Requests[ID] = true
		err = ctx.EvidenceStore.SetAllegationTracker(at)
		assert.NoError(t, err)

		ok, err := avtx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := avtx.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		ar, err = ctx.EvidenceStore.GetAllegationRequest(ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ar.Votes))

		vote := ar.Votes[0]
		assert.Equal(t, from.Bytes(), vote.Address.Bytes())
		assert.Equal(t, evidence.YES, vote.Choice)

		// could not vote twice
		tx = assemblyVoteData(ID, evidence.NO)

		ok, resp = avtx.ProcessDeliver(ctx, tx.RawTx)
		assert.False(t, ok, resp)
	})
}
