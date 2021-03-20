package evm

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

var (
	_ AccountKeeper = (*MemoryKeeper)(nil)
	_ AccountKeeper = (*KeeperStore)(nil)
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(addr keys.Address) *keys.EthAccount
	GetAccount(addr keys.Address) *keys.EthAccount
	SetAccount(account keys.EthAccount)
	RemoveAccount(account keys.EthAccount)
}

type MemoryKeeper struct {
	data map[string]*keys.EthAccount
}

func NewMemoryKeeper() *MemoryKeeper {
	return &MemoryKeeper{
		data: make(map[string]*keys.EthAccount),
	}
}

func (mk *MemoryKeeper) NewAccountWithAddress(addr keys.Address) *keys.EthAccount {
	acc := keys.NewEthAccount(addr)
	mk.SetAccount(*acc)
	return mk.GetAccount(addr)
}

func (mk *MemoryKeeper) GetAccount(addr keys.Address) *keys.EthAccount {
	return mk.data[addr.Humanize()]
}

func (mk *MemoryKeeper) SetAccount(account keys.EthAccount) {
	mk.data[account.Address.Humanize()] = keys.NewEthAccount(account.Address)
}

func (mk *MemoryKeeper) RemoveAccount(account keys.EthAccount) {
	delete(mk.data, account.Address.Humanize())
}

type KeeperStore struct {
	state  *storage.State
	prefix []byte
}

func NewKeeperStore(state *storage.State) *KeeperStore {
	return &KeeperStore{
		state:  state,
		prefix: storage.Prefix("keeper"),
	}
}

func (ks *KeeperStore) NewAccountWithAddress(addr keys.Address) *keys.EthAccount {
	acc := keys.NewEthAccount(addr)
	ks.SetAccount(*acc)
	return ks.GetAccount(addr)
}

func (ks *KeeperStore) GetAccount(addr keys.Address) *keys.EthAccount {
	prefixKey := append(ks.prefix, addr.Bytes()...)

	dat, err := ks.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil
	}

	ea := &keys.EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil
	}
	return ea
}

func (ks *KeeperStore) SetAccount(account keys.EthAccount) {
	prefixKey := append(ks.prefix, account.Address.Bytes()...)
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(account)
	if err != nil {
		return
	}
	ks.state.Set(storage.StoreKey(prefixKey), dat)
}

func (ks *KeeperStore) RemoveAccount(account keys.EthAccount) {
	prefixed := append(ks.prefix, account.Address.Bytes()...)
	ks.state.Delete(prefixed)
}
