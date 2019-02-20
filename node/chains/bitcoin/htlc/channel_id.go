/*
	Copyright 2017-2018 OneLedger

	OneLedger Wallet
*/
package htlc

import (
	"encoding/binary"
	"encoding/hex"
	"math"

	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const (
	// MaxFundingTxOutputs is the maximum number of allowed outputs on a
	// funding transaction within the protocol. This is due to the fact
	// that we use 2-bytes to encode the index within the funding output
	// during the funding workflow. Funding transaction with more outputs
	// than this are considered invalid within the protocol.
	MaxFundingTxOutputs = math.MaxUint16
)

// ChannelID is a series of 32-bytes that uniquely identifies all channels
// within the network. The ChannelID is computed using the outpoint of the
// funding transaction (the txid, and output index). Given a funding output the
// ChannelID can be calculated by XOR'ing the big-endian serialization of the
// txid and the big-endian serialization of the output index, truncated to
// 2 bytes.
type ChannelID [32]byte

// ConnectionWideID is an all-zero ChannelID, which is used to represent a
// message intended for all channels to specific peer.
var ConnectionWideID = ChannelID{}

// String returns the string representation of the ChannelID. This is just the
// hex string encoding of the ChannelID itself.
func (c ChannelID) String() string {
	return hex.EncodeToString(c[:])
}

// NewChanIDFromOutPoint converts a target OutPoint into a ChannelID that is
// usable within the network. In order to convert the OutPoint into a ChannelID,
// we XOR the lower 2-bytes of the txid within the OutPoint with the big-endian
// serialization of the Index of the OutPoint, truncated to 2-bytes.
func NewChanIDFromOutPoint(op *wire.OutPoint) ChannelID {
	// First we'll copy the txid of the outpoint into our channel ID slice.
	var cid ChannelID
	copy(cid[:], op.Hash[:])

	// With the txid copied over, we'll now XOR the lower 2-bytes of the
	// partial channelID with big-endian serialization of output index.
	xorTxid(&cid, uint16(op.Index))

	return cid
}

// xorTxid performs the transformation needed to transform an OutPoint into a
// ChannelID. To do this, we expect the cid parameter to contain the txid
// unaltered and the outputIndex to be the output index
func xorTxid(cid *ChannelID, outputIndex uint16) {
	var buf [32]byte
	binary.BigEndian.PutUint16(buf[30:], outputIndex)

	cid[30] = cid[30] ^ buf[30]
	cid[31] = cid[31] ^ buf[31]
}

// GenPossibleOutPoints generates all the possible outputs given a channel ID.
// In order to generate these possible outpoints, we perform a brute-force
// search through the candidate output index space, performing a reverse
// mapping from channelID back to OutPoint.
func (c *ChannelID) GenPossibleOutPoints() [MaxFundingTxOutputs]wire.OutPoint {
	var possiblePoints [MaxFundingTxOutputs]wire.OutPoint
	for i := uint32(0); i < MaxFundingTxOutputs; i++ {
		cidCopy := *c
		xorTxid(&cidCopy, uint16(i))

		possiblePoints[i] = wire.OutPoint{
			Hash:  chainhash.Hash(cidCopy),
			Index: i,
		}
	}

	return possiblePoints
}

// IsChanPoint returns true if the OutPoint passed corresponds to the target
// ChannelID.
func (c ChannelID) IsChanPoint(op *wire.OutPoint) bool {
	candidateCid := NewChanIDFromOutPoint(op)

	return candidateCid == c
}

// ShortChannelID represents the set of data which is needed to retrieve all
// necessary data to validate the channel existence.
type ShortChannelID struct {
	// BlockHeight is the height of the block where funding transaction
	// located.
	//
	// NOTE: This field is limited to 3 bytes.
	BlockHeight uint32

	// TxIndex is a position of funding transaction within a block.
	//
	// NOTE: This field is limited to 3 bytes.
	TxIndex uint32

	// TxPosition indicating transaction output which pays to the channel.
	TxPosition uint16
}

// NewShortChanIDFromInt returns a new ShortChannelID which is the decoded
// version of the compact channel ID encoded within the uint64. The format of
// the compact channel ID is as follows: 3 bytes for the block height, 3 bytes
// for the transaction index, and 2 bytes for the output index.
func NewShortChanIDFromInt(chanID uint64) ShortChannelID {
	return ShortChannelID{
		BlockHeight: uint32(chanID >> 40),
		TxIndex:     uint32(chanID>>16) & 0xFFFFFF,
		TxPosition:  uint16(chanID),
	}
}

// ToUint64 converts the ShortChannelID into a compact format encoded within a uint64 (8 bytes).
func (c ShortChannelID) ToUint64() uint64 {
	return ((uint64(c.BlockHeight) << 40) | (uint64(c.TxIndex) << 16) | (uint64(c.TxPosition)))
}

// String generates a human-readable representation of the channel ID.
func (c ShortChannelID) String() string {
	return fmt.Sprintf("%d:%d:%d", c.BlockHeight, c.TxIndex, c.TxPosition)
}
