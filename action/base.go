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
	Signers() []Address

	Type() string

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
	Data Msg    `json:"tx_data"`
	Fee Fee     `json:"fee"`
	Signatures  []Signature `json:"signature"`
	Memo        string  `json:"memo"`
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

		if !bytes.Equal(addr, address ) {
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

