/*

 */

package bitcoin

import (
	"errors"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type LockScriptStore struct {
	storage.SessionedStorage
	ser serialize.Serializer
}

func NewLockScriptStore(config config.Server, dbDir string) *LockScriptStore {

	store := storage.NewStorageDB(storage.KEYVALUE, "lockScriptStore", dbDir, config.Node.DB)

	return &LockScriptStore{
		store,
		serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (ls *LockScriptStore) SaveLockScript(lockScriptAddress, lockScript []byte) error {

	key := storage.StoreKey(lockScriptAddress)

	session := ls.BeginSession()

	err := session.Set(key, lockScript)
	if err != nil {
		return err
	}

	ok := session.Commit()
	if !ok {
		return errors.New("err committing to lockscript store")
	}

	return nil
}

func (ls *LockScriptStore) GetLockScript(lockScriptAddress []byte) ([]byte, error) {

	key := storage.StoreKey(lockScriptAddress)

	return ls.Get(key)
}
