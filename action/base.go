package action

import (
	"encoding/hex"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type MsgData []byte

type Msg interface {
	// return the necessary signers for the message, should have consistent order across the network
	Signers() []Address

	Type() Type

	Tags() common.KVPairs

	Marshal() ([]byte, error)

	Unmarshal([]byte) error
}

type Fee struct {
	Price Amount
	Gas   int64
}

type Signature struct {
	Signer keys.PublicKey
	Signed []byte
}

func (s Signature) Verify(msg []byte) bool {
	handler, err := s.Signer.GetHandler()
	if err != nil {
		return false
	}
	return handler.VerifyBytes(msg, s.Signed)
}

type RawTx struct {
	Type Type    `json:"tx_type"`
	Data MsgData `json:"tx_data"`
	Fee  Fee     `json:"fee"`
	Memo string  `json:"memo"`
}

func (t *RawTx) RawBytes() []byte {
	value, err := serialize.GetSerializer(serialize.NETWORK).Serialize(t)
	if err != nil {
		logger.Error("failed to serialize tx: ", t)
	}
	return value
}

type SignedTx struct {
	RawTx
	Signatures []Signature `json:"signatures"`
}

func (t *SignedTx) SignedBytes() []byte {
	value, err := serialize.GetSerializer(serialize.NETWORK).Serialize(t)
	if err != nil {
		logger.Error("failed to serialize tx: ", t)
	}
	return value
}

func ValidateBasic(data []byte, signerAddr []Address, signatures []Signature) (bool, error) {
	for i, s := range signerAddr {
		pkey := signatures[i].Signer
		h, err := pkey.GetHandler()
		if err != nil {
			return false, ErrInvalidPubkey
		}
		if !h.Address().Equal(s) {
			return false, errors.Wrap(ErrUnmatchSigner, hex.EncodeToString(h.Address())+","+hex.EncodeToString(s))
		}

		if !h.VerifyBytes(data, signatures[i].Signed) {
			return false, ErrInvalidSignature
		}
	}

	return true, nil
}
