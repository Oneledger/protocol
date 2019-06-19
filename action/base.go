package action

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/tendermint/tendermint/libs/common"
)

type Msg interface {
	// return the necessary signers for the message, should have consistent order across the network
	Signers() []Address

	Type() Type

	Tags() common.KVPairs
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

type BaseTx struct {
	Data       Msg         `json:"tx_data"`
	Fee        Fee         `json:"fee"`
	Signatures []Signature `json:"signatures"`
	Memo       string      `json:"memo"`
}

func (t *BaseTx) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.NETWORK).Serialize(t)
	if err != nil {
		logger.Error("failed to serialize tx: ", t)
	}
	return value
}

func ValidateBasic(msg Msg, fee Fee, memo string, signatures []Signature) (bool, error) {
	toVerify := (&BaseTx{msg, fee, nil, memo}).Bytes()
	for i, s := range msg.Signers() {
		pkey := signatures[i].Signer
		h, err := pkey.GetHandler()
		if err != nil {
			return false, ErrInvalidPubkey
		}
		if !h.Address().Equal(s) {
			return false, ErrUnmatchSigner
		}

		if !h.VerifyBytes(toVerify, signatures[i].Signed) {
			return false, ErrInvalidSignature
		}
	}

	if !verifyMinimumFee(fee) {
		return false, ErrInvalidFee
	}

	return true, nil
}

func verifyMinimumFee(fee Fee) bool {
	//TODO: implement minimum fee check
	return true
}
