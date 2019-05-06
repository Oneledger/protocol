package action

import (
	"bytes"
	"fmt"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)


type Address = keys.Address

type Coin = balance.Coin

type Context struct {
	Accounts data.Store
	Balances data.Store
}

type Msg interface {
	// return the necessary signers for the message, should have consistent order across the network
	GetSigners() []Address

	Type() string

	GetBytes() []byte


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
	Data Msg
	Fee Fee
	Signatures  []Signature
	Memo        string
}


func (t *BaseTx) Sign(ctx Context) error {
	addrs := t.Data.GetSigners()

	if t.Signatures == nil {
		t.Signatures = make([]Signature, len(addrs))
	}

	for i, addr := range addrs {
		value, err := ctx.Accounts.Get(addr.Bytes())
		if err != nil {
			return fmt.Errorf("failed to get account for sign: %s", err)
		}
		account := (&accounts.Account{}).FromBytes(value)
		signed, err := account.Sign(t.Data.GetBytes())
		if err != nil {
			return fmt.Errorf("failed to sign: %s", err)
		}
		t.Signatures[i] = Signature{account.PublicKey, signed}
	}
	return nil
}


func (t *BaseTx) SignWithAddress(ctx Context, address Address) error {
	addrs := t.Data.GetSigners()

	if t.Signatures == nil {
		t.Signatures = make([]Signature, len(addrs))
	}

	for i, addr := range addrs {

		if !bytes.Equal(addr, address ) {
			continue
		}

		value, err := ctx.Accounts.Get(addr.Bytes())
		if err != nil {
			return fmt.Errorf("failed to get account for sign: %s", err)
		}

		account := (&accounts.Account{}).FromBytes(value)
		signed, err := account.Sign(t.Data.GetBytes())
		if err != nil {
			return fmt.Errorf("failed to sign: %s", err)
		}
		t.Signatures[i] = Signature{account.PublicKey, signed}
	}
	return nil
}

