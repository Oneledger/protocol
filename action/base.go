package action

import (
	"bytes"
	"fmt"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type Context struct {
	Accounts data.Store
	Balances storage.ChainState
}

type Msg interface {
	// return the necessary signers for the message, should have consistent order across the network
	Signers() []Address

	Type() Type

	Bytes() []byte
}

type Fee struct {
	Price balance.Coin
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

func (t *BaseTx) Sign(ctx Context) error {
	addrs := t.Data.Signers()

	if t.Signatures == nil {
		t.Signatures = make([]Signature, len(addrs))
	}

	for i, addr := range addrs {
		signed, err := sign(ctx, addr, t.Data.Bytes())
		if err != nil {
			return err
		}
		t.Signatures[i] = signed
	}
	return nil
}

func (t *BaseTx) SignWithAddress(ctx Context, address Address) error {
	addrs := t.Data.Signers()

	if t.Signatures == nil {
		t.Signatures = make([]Signature, len(addrs))
	}

	for i, addr := range addrs {

		if !bytes.Equal(addr, address) {
			continue
		}

		signed, err := sign(ctx, addr, t.Data.Bytes())
		if err != nil {
			return err
		}
		t.Signatures[i] = signed
	}
	return nil
}

func (t *BaseTx) valideBasic() (bool, error) {
	msg := t.Data
	signatures := t.Signatures
	fee := t.Fee
	for i, s := range msg.Signers() {
		pkey := signatures[i].Signer
		h, err := pkey.GetHandler()
		if err != nil {
			return false, ErrInvalidPubkey
		}
		if !h.Address().Equal(s) {
			return false, ErrUnmatchSigner
		}
		if !h.VerifyBytes(msg.Bytes(), signatures[i].Signed) {
			return false, ErrInvalidSignature
		}
	}

	if !verifyMinimumFee(fee) {
		return false, ErrInvalidFee
	}

	return true, nil
}

func sign(ctx Context, address Address, msg []byte) (Signature, error) {
	value, err := ctx.Accounts.Get(address.Bytes())
	if err != nil {
		return Signature{}, fmt.Errorf("failed to get account for sign: %s", err)
	}

	account := (&accounts.Account{}).FromBytes(value)
	signed, err := account.Sign(msg)
	if err != nil {
		return Signature{}, fmt.Errorf("failed to sign: %s", err)
	}
	return Signature{account.PublicKey, signed}, nil
}

func verifyMinimumFee(fee Fee) bool {
	//todo: implement minimum fee check
	return true
}
