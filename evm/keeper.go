package evm

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

var _ AccountKeeper = (*MemoryKeeper)(nil)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx *action.Context, addr keys.Address) *EthAccount
	GetAccount(ctx *action.Context, addr keys.Address) *EthAccount
	SetAccount(ctx *action.Context, account EthAccount)
	RemoveAccount(ctx *action.Context, account EthAccount)
}

type MemoryKeeper struct {
	data map[string]*EthAccount
}

func NewMemoryKeeper() *MemoryKeeper {
	return &MemoryKeeper{
		data: make(map[string]*EthAccount),
	}
}

func (mk *MemoryKeeper) NewAccountWithAddress(ctx *action.Context, addr keys.Address) *EthAccount {
	acc := NewEthAccount(addr)
	mk.SetAccount(ctx, *acc)
	return mk.GetAccount(ctx, addr)
}

func (mk *MemoryKeeper) GetAccount(ctx *action.Context, addr keys.Address) *EthAccount {
	return mk.data[addr.Humanize()]
}

func (mk *MemoryKeeper) SetAccount(ctx *action.Context, account EthAccount) {
	mk.data[account.Address.Humanize()] = NewEthAccount(account.Address)
}

func (mk *MemoryKeeper) RemoveAccount(ctx *action.Context, account EthAccount) {
	delete(mk.data, account.Address.Humanize())
}
