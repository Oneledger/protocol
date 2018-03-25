package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"

	tps "github.com/tendermint/tendermint/types"
	"time"
)
func randVoteSet(height int64, round int, type_ byte, numValidators int, votingPower int64) (*tps.VoteSet, *tps.ValidatorSet, []*tps.PrivValidatorFS) {
	valSet, privValidators := tps.RandValidatorSet(numValidators, votingPower)
	return tps.NewVoteSet("test_chain_id", height, round, type_, valSet), valSet, privValidators
}


func TestValidateBlock(t *testing.T) {
	txs := []tps.Tx{tps.Tx("foo"), tps.Tx("bar")}
	lastID := makeBlockIDRandom()
	h := int64(3)

	voteSet, _, vals := randVoteSet(h-1, 1, VoteTypePrecommit,
		10, 1)
	commit, err := MakeCommit(lastID, h-1, 1, voteSet, vals)
	require.NoError(t, err)

	block := MakeBlock(h, txs, commit)
	require.NotNil(t, block)

	// proper block must pass
	err = block.ValidateBasic()
	require.NoError(t, err)

	// tamper with NumTxs
	block = MakeBlock(h, txs, commit)
	block.NumTxs += 1
	err = block.ValidateBasic()
	require.Error(t, err)

	// remove 1/2 the commits
	block = MakeBlock(h, txs, commit)
	block.LastCommit.Precommits = commit.Precommits[:commit.Size()/2]
	block.LastCommit.hash = nil // clear hash or change wont be noticed
	err = block.ValidateBasic()
	require.Error(t, err)

	// tamper with LastCommitHash
	block = MakeBlock(h, txs, commit)
	block.LastCommitHash = []byte("something else")
	err = block.ValidateBasic()
	require.Error(t, err)

	// tamper with data
	block = MakeBlock(h, txs, commit)
	block.Data.Txs[0] = tps.Tx("something else")
	block.Data.hash = nil // clear hash or change wont be noticed
	err = block.ValidateBasic()
	require.Error(t, err)

	// tamper with DataHash
	block = MakeBlock(h, txs, commit)
	block.DataHash = cmn.RandBytes(len(block.DataHash))
	err = block.ValidateBasic()
	require.Error(t, err)
}

func makeBlockIDRandom() BlockID {
	blockHash, blockPartsHeader := crypto.CRandBytes(32), tps.PartSetHeader{123, crypto.CRandBytes(32)}
	return BlockID{blockHash, blockPartsHeader}
}

func makeBlockID(hash string, partSetSize int, partSetHash string) BlockID {
	return BlockID{
		Hash: []byte(hash),
		PartsHeader: tps.PartSetHeader{
			Total: partSetSize,
			Hash:  []byte(partSetHash),
		},
	}

}

func MakeCommit(blockID tps.BlockID, height int64, round int,
	voteSet *tps.VoteSet,
	validators []*tps.PrivValidatorFS, ) (*Commit, error) {
	return &Commit{ BlockID: blockID,height}
}

func signAddVote(privVal *tps.PrivValidatorFS, vote *Vote, voteSet *tps.VoteSet) (signed bool, err error) {
	vote.Signature, err = privVal.Signer.Sign(SignBytes(voteSet.ChainID(), vote))
	if err != nil {
		return false, err
	}
	return voteSet.AddVote(vote)
}
