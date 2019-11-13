/*

 */

package bitcoin

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

type KeyDB struct {
	data map[btcutil.Address]*btcec.PrivateKey
}

func NewKeyDB() *KeyDB {
	return &KeyDB{
		map[btcutil.Address]*btcec.PrivateKey{},
	}
}

func (k *KeyDB) Add(key btcutil.Address, privateKey *btcec.PrivateKey) {
	k.data[key] = privateKey
}

func (k KeyDB) GetKey(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
	var err error
	key, ok := k.data[address]
	if !ok {
		err = errors.New("address not found")
	}
	return key, ok, err
}
