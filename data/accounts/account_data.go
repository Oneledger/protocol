/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package accounts

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"

	"github.com/pkg/errors"
)

var _ serialize.DataAdapter = &Account{}

type AccountData struct {
	Type chain.Type
	Name string

	PublicKeyData  []byte
	PrivateKeyData []byte
}

func (t *Account) NewDataInstance() serialize.Data {
	return &AccountData{}
}

func (t *Account) Data() serialize.Data {
	ad := &AccountData{
		Type: t.Type,
		Name: t.Name,
	}
	ad.PublicKeyData, _ = t.PublicKey.GobEncode()
	ad.PrivateKeyData, _ = t.PrivateKey.GobEncode()

	return ad
}

func (t *Account) SetData(a interface{}) error {
	ad, ok := a.(*AccountData)
	if !ok {
		return errors.New("Wrong data")
	}

	t.PublicKey = &keys.PublicKey{}
	t.PrivateKey = &keys.PrivateKey{}

	t.Type = ad.Type
	t.Name = ad.Name

	err := t.PublicKey.GobDecode(ad.PublicKeyData)
	if err != nil {
		return err
	}

	err = t.PrivateKey.GobDecode(ad.PrivateKeyData)
	if err != nil {
		return err
	}

	return nil
}

func (ad *AccountData) SerialTag() string {
	return ""
}
