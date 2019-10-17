/*

 */

package keys

import (
	"encoding/json"

	"github.com/Oneledger/protocol/utils"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

type BTCMultiSig struct {
	Msg []byte `json:"msg"`

	M int `json:"m"`

	Signers []btcutil.AddressPubKey `json:"signers"`

	Signatures []BTCSignature `json:"signatures"`
}

type BTCSignature struct {
	Index  int    `json:"index"`
	PubKey []byte `json:"pubkey"` // this should be a compressed public key
	Sign   []byte `json:"sign"`
}

func NewBTCMultiSig(msg []byte, m int, signers []btcutil.AddressPubKey) (*BTCMultiSig, error) {

	if msg == nil {
		return nil, ErrMissMsg
	}

	if signers == nil {
		return nil, ErrMissSigners
	}

	if m < 0 || m > len(signers) {
		return nil, ErrInvalidThreshold
	}

	return &BTCMultiSig{
		Msg:        msg,
		M:          m,
		Signers:    signers,
		Signatures: make([]BTCSignature, len(signers)),
	}, nil
}

func (m *BTCMultiSig) AddSignature(s *BTCSignature, params *chaincfg.Params) error {

	address, err := btcutil.NewAddressPubKey(s.PubKey, params)
	if err != nil {
		return errors.Wrap(err, "failed to add")
	}

	if address.String() != m.Signers[s.Index].String() {
		return ErrNotExpectedSigner
	}

	// TODO verify the signature using btc libs

	m.Signatures[s.Index] = *s
	return nil
}

func (m BTCMultiSig) IsValid() bool {
	cnt := 0
	for _, item := range m.Signatures {
		if item.Sign != nil {
			cnt++
		}
	}
	if cnt < m.M {
		return false
	}
	return true
}

func (m BTCMultiSig) Address() Address {
	s := &BTCMultiSig{m.Msg, m.M, m.Signers, nil}
	b := s.Bytes()
	return utils.Hash(b)
}

func (m BTCMultiSig) Bytes() []byte {

	signatures := m.Signatures
	m.Signatures = make([]BTCSignature, 0)

	for _, item := range signatures {
		if item.Sign != nil {
			m.Signatures = append(m.Signatures, item)
		}
	}

	b, _ := json.Marshal(m)
	return b
}

func (m BTCMultiSig) HasAddressSigned(addr btcutil.AddressPubKey) bool {
	index := len(m.Signers) + 100
	for i := range m.Signers {
		if m.Signers[i].String() == addr.String() {
			index = i
			break
		}
	}

	if index > len(m.Signers) {
		return false
	}

	if len(m.Signatures[index].Sign) == 0 {
		return false
	}

	return true
}

func (m *BTCMultiSig) FromBytes(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}

	signatures := m.Signatures
	m.Signatures = make([]BTCSignature, len(m.Signers))
	for i, item := range signatures {
		m.Signatures[i] = item
	}

	return nil
}

func (m *BTCMultiSig) GetSignerIndex(addr btcutil.AddressPubKey) (int, error) {

	for i := range m.Signers {
		if m.Signers[i].String() == addr.String() {
			return i, nil
		}
	}

	return 0, errors.New("address not found")
}

func (m *BTCMultiSig) GetSignatures() []BTCSignature {
	return m.Signatures
}
