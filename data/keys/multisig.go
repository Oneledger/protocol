package keys

import (
	"bytes"
	"encoding/json"

	"github.com/Oneledger/protocol/utils"

	"github.com/pkg/errors"
)

type MultiSigner interface {
	//Initialize the MultiSig
	// "msg" that need to be sign by this multisig
	// "threshold" used for threshold signature representing the minimal signatures for the multisig to be valid
	// "signers" are the expected signers that need to sign for the message, the signers will be sorted in the slice.
	Init(msg []byte, threshold int, signers []Address) error

	// Add a signature to the MultiSig. Return error if the signature not satisfy:
	//  - Valid s.PubKey in signature
	//  - PubKey's address doesn't match signer address of s.Index
	//  - Signed info is not the signature for message in the MultiSig
	AddSignature(s Signature) error

	// check MultiSig contains enough signature with threshold
	IsValid() bool

	//// Return the threshold for MultiSig
	//// threshold is between (0, len(Signer)]
	//// when threshold = len(Signer): need all signers to sign
	//Threshold() int
	//
	//// Return the expected signers for the MultiSig
	//Signers() []Address
	//
	//// Return the expected signing message for the MultiSig
	//Message() []byte
	//
	//// Return the address for the MultiSig, which act as unique identifier of MultiSig
	//// The address should only depend on (msg, m, signers) the signatures should not influence the address
	//Address() Address

	// Return the bytes for the MultiSig,
	// the Bytes should include all info
	Bytes() []byte

	// Set MultiSig from []byte
	FromBytes(b []byte) error
}

type Signature struct {
	Index  int       `json:"index"` //index in the MultiSig
	PubKey PublicKey `json:"pubkey"`
	Signed []byte    `json:"signed"`
}

var _ MultiSigner = &MultiSig{}

type MultiSig struct {

	// message that need to be sign by this multisig
	Msg []byte

	//used for threshold signature where M represent the minimal signatures for the multisig to be valid
	// M is between (0, len(Signer)]
	// when M = len(Signer): need all signers to sign
	M int

	//the expected signers that need to sign for the message, the signers will be sorted in the slice.
	Signers []Address

	//the collection of signatures for signers, the index of signatures should match the index of signers. this should
	// be empty when created and before add Signature
	Signatures []Signature
}

func (m *MultiSig) Init(msg []byte, threshold int, signers []Address) error {

	if msg == nil {
		return ErrMissMsg
	}

	if signers == nil {
		return ErrMissSigners
	}

	if threshold < 0 || threshold > len(signers) {
		return ErrInvalidThreshold
	}
	*m = MultiSig{
		Msg:        msg,
		M:          threshold,
		Signers:    signers,
		Signatures: make([]Signature, len(signers)),
	}
	return nil
}

func (m *MultiSig) AddSignature(s Signature) error {
	h, err := s.PubKey.GetHandler()
	if err != nil {
		return errors.Wrap(err, "failed to add")
	}
	if !bytes.Equal(m.Signers[s.Index], h.Address().Bytes()) {
		return ErrNotExpectedSigner
	}
	if !h.VerifyBytes(m.Msg, s.Signed) {
		return ErrInvalidSignedMsg
	}
	m.Signatures[s.Index] = s
	return nil
}

func (m MultiSig) IsValid() bool {
	cnt := 0
	for _, item := range m.Signatures {
		if item.Signed != nil {
			cnt++
		}
	}
	if cnt < m.M {
		return false
	}
	return true
}

func (m MultiSig) Address() Address {
	s := &MultiSig{m.Msg, m.M, m.Signers, nil}
	b := s.Bytes()
	return utils.Hash(b)
}

//func (m MultiSig) Threshold() int {
//	return m.M
//}
//
//func (m MultiSig) Signers() []Address {
//	return m.Signers
//}
//
//func (m MultiSig) Message() []byte {
//	return m.Msg
//}

func (m MultiSig) Bytes() []byte {
	signatures := m.Signatures
	m.Signatures = make([]Signature, 0)
	for _, item := range signatures {
		if item.Signed != nil {
			m.Signatures = append(m.Signatures, item)
		}
	}

	b, _ := json.Marshal(m)
	return b
}

func (m *MultiSig) FromBytes(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	signatures := m.Signatures
	m.Signatures = make([]Signature, len(m.Signers))
	for i, item := range signatures {
		m.Signatures[i] = item
	}
	return nil
}
