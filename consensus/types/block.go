package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	wire "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
	tps "github.com/tendermint/tendermint/types"
	"golang.org/x/crypto/ripemd160"
)

// Block defines the atomic unit of a oneledger blockchain.
// TODO: add Version byte
type Block struct {
	*Header    `json:"header"`
	*Data      `json:"data"`
	Evidence   EvidenceData `json:"evidence"`
	LastCommit *Commit      `json:"last_commit"`
}

// MakeBlock returns a new block with an empty header, except what can be computed from itself.
// It populates the same set of fields validated by ValidateBasic
func MakeBlock(height int64, txs []tps.Tx, commit *Commit) *Block {
	block := &Block{
		Header: &Header{
			Height: height,
			Time:   time.Now(),
			NumTxs: int64(len(txs)),
		},
		LastCommit: commit,
		Data: &Data{
			Txs: txs,
		},
	}
	block.FillHeader()
	return block
}

// AddEvidence appends the given evidence to the block
func (b *Block) AddEvidence(evidence []tps.Evidence) {
	b.Evidence.Evidence = append(b.Evidence.Evidence, evidence...)
}

// ValidateBasic performs basic validation that doesn't involve state data.
// It checks the internal consistency of the block.
func (b *Block) ValidateBasic() error {
	newTxs := int64(len(b.Data.Txs))
	if b.NumTxs != newTxs {
		return fmt.Errorf("Wrong Block.Header.NumTxs. Expected %v, got %v", newTxs, b.NumTxs)
	}
	if !bytes.Equal(b.LastCommitHash, b.LastCommit.Hash()) {
		return fmt.Errorf("Wrong Block.Header.LastCommitHash.  Expected %v, got %v", b.LastCommitHash, b.LastCommit.Hash())
	}
	if b.Header.Height != 1 {
		if err := b.LastCommit.ValidateBasic(); err != nil {
			return err
		}
	}
	if !bytes.Equal(b.DataHash, b.Data.Hash()) {
		return fmt.Errorf("Wrong Block.Header.DataHash.  Expected %v, got %v", b.DataHash, b.Data.Hash())
	}
	if !bytes.Equal(b.EvidenceHash, b.Evidence.Hash()) {
		return errors.New(cmn.Fmt("Wrong Block.Header.EvidenceHash.  Expected %v, got %v", b.EvidenceHash, b.Evidence.Hash()))
	}
	return nil
}

// FillHeader fills in any remaining header fields that are a function of the block data
func (b *Block) FillHeader() {
	if b.LastCommitHash == nil {
		b.LastCommitHash = b.LastCommit.Hash()
	}
	if b.DataHash == nil {
		b.DataHash = b.Data.Hash()
	}
	if b.EvidenceHash == nil {
		b.EvidenceHash = b.Evidence.Hash()
	}
}

// Hash computes and returns the block hash.
// If the block is incomplete, block hash is nil for safety.
func (b *Block) Hash() cmn.HexBytes {
	if b == nil || b.Header == nil || b.Data == nil || b.LastCommit == nil {
		return nil
	}
	b.FillHeader()
	return b.Header.Hash()
}

// MakePartSet returns a PartSet containing parts of a serialized block.
// This is the form in which the block is gossipped to peers.
func (b *Block) MakePartSet(partSize int) *tps.PartSet {
	return tps.NewPartSetFromData(wire.BinaryBytes(b), partSize)
}

// HashesTo is a convenience function that checks if a block hashes to the given argument.
// A nil block never hashes to anything, and nothing hashes to a nil hash.
func (b *Block) HashesTo(hash []byte) bool {
	if len(hash) == 0 {
		return false
	}
	if b == nil {
		return false
	}
	return bytes.Equal(b.Hash(), hash)
}

// String returns a string representation of the block
func (b *Block) String() string {
	return b.StringIndented("")
}

// StringIndented returns a string representation of the block
func (b *Block) StringIndented(indent string) string {
	if b == nil {
		return "nil-Block"
	}
	return fmt.Sprintf(`Block{
%s  %v
%s  %v
%s  %v
%s  %v
%s}#%v`,
		indent, b.Header.StringIndented(indent+"  "),
		indent, b.Data.StringIndented(indent+"  "),
		indent, b.Evidence.StringIndented(indent+"  "),
		indent, b.LastCommit.StringIndented(indent+"  "),
		indent, b.Hash())
}

// StringShort returns a shortened string representation of the block
func (b *Block) StringShort() string {
	if b == nil {
		return "nil-Block"
	} else {
		return fmt.Sprintf("Block#%v", b.Hash())
	}
}

//-----------------------------------------------------------------------------

// Header defines the structure of a oneledger block header
// TODO: limit header size
// NOTE: changes to the Header should be duplicated in the abci Header
type Header struct {
	// basic block info
	ChainID string    `json:"chain_id"`
	Height  int64     `json:"height"`
	Time    time.Time `json:"time"`
	NumTxs  int64     `json:"num_txs"`

	// prev block info
	LastBlockID tps.BlockID `json:"last_block_id"`
	TotalTxs    int64       `json:"total_txs"`

	// hashes of block data
	LastCommitHash cmn.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       cmn.HexBytes `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash  cmn.HexBytes `json:"validators_hash"`   // validators for the current block
	ConsensusHash   cmn.HexBytes `json:"consensus_hash"`    // consensus params for current block
	AppHash         cmn.HexBytes `json:"app_hash"`          // state after txs from the previous block
	LastResultsHash cmn.HexBytes `json:"last_results_hash"` // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash cmn.HexBytes `json:"evidence_hash"` // evidence included in the block
}

// Hash returns the hash of the header.
// Returns nil if ValidatorHash is missing.
func (h *Header) Hash() cmn.HexBytes {
	if len(h.ValidatorsHash) == 0 {
		return nil
	}
	return merkle.SimpleHashFromMap(map[string]merkle.Hasher{
		"ChainID":     wireHasher(h.ChainID),
		"Height":      wireHasher(h.Height),
		"Time":        wireHasher(h.Time),
		"NumTxs":      wireHasher(h.NumTxs),
		"TotalTxs":    wireHasher(h.TotalTxs),
		"LastBlockID": wireHasher(h.LastBlockID),
		"LastCommit":  wireHasher(h.LastCommitHash),
		"Data":        wireHasher(h.DataHash),
		"Validators":  wireHasher(h.ValidatorsHash),
		"App":         wireHasher(h.AppHash),
		"Consensus":   wireHasher(h.ConsensusHash),
		"Results":     wireHasher(h.LastResultsHash),
		"Evidence":    wireHasher(h.EvidenceHash),
	})
}

// StringIndented returns a string representation of the header
func (h *Header) StringIndented(indent string) string {
	if h == nil {
		return "nil-Header"
	}
	return fmt.Sprintf(`Header{
%s  ChainID:        %v
%s  Height:         %v
%s  Time:           %v
%s  NumTxs:         %v
%s  TotalTxs:       %v
%s  LastBlockID:    %v
%s  LastCommit:     %v
%s  Data:           %v
%s  Validators:     %v
%s  App:            %v
%s  Conensus:       %v
%s  Results:        %v
%s  Evidence:       %v
%s}#%v`,
		indent, h.ChainID,
		indent, h.Height,
		indent, h.Time,
		indent, h.NumTxs,
		indent, h.TotalTxs,
		indent, h.LastBlockID,
		indent, h.LastCommitHash,
		indent, h.DataHash,
		indent, h.ValidatorsHash,
		indent, h.AppHash,
		indent, h.ConsensusHash,
		indent, h.LastResultsHash,
		indent, h.EvidenceHash,
		indent, h.Hash())
}

//-------------------------------------

//-----------------------------------------------------------------------------

// SignedHeader is a header along with the commits that prove it
type SignedHeader struct {
	Header *Header `json:"header"`
	Commit *Commit `json:"commit"`
}

//-----------------------------------------------------------------------------

// Data contains the set of transactions included in the block
type Data struct {

	// Tx that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs tps.Txs `json:"txs"`

	// Volatile
	hash cmn.HexBytes
}

// Hash returns the hash of the data
func (data *Data) Hash() cmn.HexBytes {
	if data.hash == nil {
		data.hash = data.Txs.Hash() // NOTE: leaves of merkle tree are TxIDs
	}
	return data.hash
}

// StringIndented returns a string representation of the transactions
func (data *Data) StringIndented(indent string) string {
	if data == nil {
		return "nil-Data"
	}
	txStrings := make([]string, cmn.MinInt(len(data.Txs), 21))
	for i, tx := range data.Txs {
		if i == 20 {
			txStrings[i] = fmt.Sprintf("... (%v total)", len(data.Txs))
			break
		}
		txStrings[i] = fmt.Sprintf("Tx:%v", tx)
	}
	return fmt.Sprintf(`Data{
%s  %v
%s}#%v`,
		indent, strings.Join(txStrings, "\n"+indent+"  "),
		indent, data.hash)
}

//-----------------------------------------------------------------------------

// EvidenceData contains any evidence of malicious wrong-doing by validators
type EvidenceData struct {
	Evidence tps.EvidenceList `json:"evidence"`

	// Volatile
	hash cmn.HexBytes
}

// Hash returns the hash of the data.
func (data *EvidenceData) Hash() cmn.HexBytes {
	if data.hash == nil {
		data.hash = data.Evidence.Hash()
	}
	return data.hash
}

// StringIndented returns a string representation of the evidence.
func (data *EvidenceData) StringIndented(indent string) string {
	if data == nil {
		return "nil-Evidence"
	}
	evStrings := make([]string, cmn.MinInt(len(data.Evidence), 21))
	for i, ev := range data.Evidence {
		if i == 20 {
			evStrings[i] = fmt.Sprintf("... (%v total)", len(data.Evidence))
			break
		}
		evStrings[i] = fmt.Sprintf("Evidence:%v", ev)
	}
	return fmt.Sprintf(`Data{
%s  %v
%s}#%v`,
		indent, strings.Join(evStrings, "\n"+indent+"  "),
		indent, data.hash)
	return ""
}

//-------------------------------------------------------

type hasher struct {
	item interface{}
}

func (h hasher) Hash() []byte {
	hasher, n, err := ripemd160.New(), new(int), new(error)
	wire.WriteBinary(h.item, hasher, n, err)
	if *err != nil {
		panic(err)
	}
	return hasher.Sum(nil)

}

func wireHasher(item interface{}) merkle.Hasher {
	return hasher{item}
}
