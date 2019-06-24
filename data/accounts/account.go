package accounts

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type Account struct {
	Type chain.Type `json:"type"`
	Name string     `json:"name"`

	PublicKey  *keys.PublicKey  `json:"publicKey"`
	PrivateKey *keys.PrivateKey `json:"privateKey"`
}

func NewAccount(t chain.Type, name string, privkey *keys.PrivateKey, pubkey *keys.PublicKey) (Account, error) {
	if pubkey == nil || privkey == nil {
		return Account{}, errors.New("empty keys")
	}

	acc := Account{
		Type:       t,
		Name:       name,
		PrivateKey: privkey,
		PublicKey:  pubkey,
	}

	if _, err := acc.PublicKey.GetHandler(); err != nil {
		return Account{}, err
	}
	if _, err := acc.PrivateKey.GetHandler(); err != nil {
		return Account{}, err
	}
	return acc, nil
}

func GenerateNewAccount(t chain.Type, name string) (Account, error) {
	pubKey, privKey, err := keys.NewKeyPairFromTendermint()
	if err != nil {
		return Account{}, err
	}
	acc := Account{
		Type:       t,
		Name:       name,
		PrivateKey: &privKey,
		PublicKey:  &pubKey,
	}
	return acc, nil
}

func (acc Account) Address() keys.Address {
	handler, err := acc.PublicKey.GetHandler()
	if err != nil {
		logger.Fatal("PublicKey format for account is wrong", err)
	}

	return keys.Address(handler.Address())
}

func (acc Account) Sign(msg []byte) ([]byte, error) {
	signer, err := acc.PrivateKey.GetHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to get the signer: %s ", err)
	}
	signed, err := signer.Sign(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the msg: %s", err)
	}
	return signed, nil
}

func (acc Account) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&acc)
	if err != nil {
		//logger.Fatal("account not serializable")
		return nil
	}
	return value
}

func (acc *Account) FromBytes(msg []byte) *Account {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, acc)
	if err != nil {
		logger.Fatal("failed to deserialize account from bytes", err)
		return nil
	}
	return acc
}

func (acc Account) String() string {
	//pkey := acc.PublicKey.GetABCIPubKey()
	result := ""
	result = result +
		fmt.Sprintf("name %s \n", acc.Name) +
		//fmt.Sprintf("type %s \n", acc.Type.String()) +
		//fmt.Sprintf("pubkey %s \n", pkey.String()) +
		fmt.Sprintf("address %s \n", acc.Address())
	return result
}
