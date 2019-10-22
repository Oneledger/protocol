/*

 */

package bitcoin

import (
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type LockScriptStore struct {
	storage.ChainState
	ser serialize.Serializer
}

func NewLockScriptStore(store *storage.ChainState) *LockScriptStore {
	return &LockScriptStore{
		*store,
		serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (ls *LockScriptStore) SaveLockScript(lockScriptAddress, lockScript []byte) error {

	key := storage.StoreKey(lockScriptAddress)
	return ls.Set(key, lockScript)
}

func (ls *LockScriptStore) GetLockScript(lockScriptAddress []byte) ([]byte, error) {

	key := storage.StoreKey(lockScriptAddress)

	return ls.Get(key)
}
