package balance

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type EthAccount struct {
	Address  keys.Address `json:"address"`
	CodeHash []byte       `json:"codeHash"`
	Amount   *Amount      `json:"amount"`
	Sequence uint64       `json:"sequence"`
}

func NewEthAccount(addr keys.Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		CodeHash: ethcrypto.Keccak256(nil),
		Amount:   NewAmountFromInt(0),
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

func (acc *EthAccount) IncrementNonce() uint64 {
	return atomic.AddUint64(&acc.Sequence, 1)
}

func (acc EthAccount) Balance() *big.Int {
	return acc.Amount.BigInt()
}

func (acc *EthAccount) AddBalance(amount *big.Int) {
	amt := NewAmountFromBigInt(amount)
	acc.Amount = acc.Amount.Plus(*amt)
}

func (acc *EthAccount) SubBalance(amount *big.Int) {
	amt := NewAmountFromBigInt(amount)
	newAmt, err := acc.Amount.Minus(*amt)
	if err != nil {
		panic(errors.Errorf("Failed to minus balance: %s", err))
	}
	acc.Amount = newAmt
}

func (acc *EthAccount) SetBalance(amount *big.Int) {
	acc.Amount = NewAmountFromBigInt(amount)
}

type AccountKeeper interface {
	NewAccountWithAddress(addr keys.Address) (*EthAccount, error)
	GetOrCreateAccount(addr keys.Address) (*EthAccount, error)
	GetAccount(addr keys.Address) (*EthAccount, error)
	GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error)
	SetAccount(account EthAccount) error
	RemoveAccount(account EthAccount)
	WithState(state *storage.State) AccountKeeper
}

var _ AccountKeeper = (*NesterAccountKeeper)(nil)

// NesterAccountKeeper is used to keep track balances and nonces
type NesterAccountKeeper struct {
	state  *storage.State
	prefix []byte
}

func NewNesterAccountKeeper(state *storage.State) AccountKeeper {
	return &NesterAccountKeeper{
		state:  state,
		prefix: storage.Prefix("keeper"),
	}
}

func (nak *NesterAccountKeeper) WithState(state *storage.State) AccountKeeper {
	nak.state = state
	return nak
}

func (nak *NesterAccountKeeper) NewAccountWithAddress(addr keys.Address) (*EthAccount, error) {
	acc := NewEthAccount(addr)

	err := nak.SetAccount(*acc)
	if err != nil {
		return nil, errors.Errorf("Failed to set account: %s", err)
	}
	acc, err = nak.GetAccount(addr)
	if err != nil {
		return nil, errors.Errorf("Failed to get account: %s", err)
	}
	return acc, nil
}

func (nak *NesterAccountKeeper) GetOrCreateAccount(addr keys.Address) (*EthAccount, error) {
	eoa, err := nak.GetAccount(addr)
	if err != nil {
		eoa, err = nak.NewAccountWithAddress(addr)
		if err != nil {
			return nil, err
		}
	}
	return eoa, nil
}

func (nak *NesterAccountKeeper) GetAccount(addr keys.Address) (*EthAccount, error) {
	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat, err := nak.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}

	if len(dat) == 0 {
		return nil, errors.Errorf("account does not exist")
	}

	eoa := &EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, eoa)
	if err != nil {
		return nil, err
	}
	return eoa, nil
}

func (nak *NesterAccountKeeper) GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error) {
	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat := nak.state.GetVersioned(height, storage.StoreKey(prefixKey))
	if len(dat) == 0 {
		return nil, errors.New(fmt.Sprintf("Previous state on height '%d' for addr '%s' not found", height, addr))
	}

	eoa := &EthAccount{}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, eoa)
	if err != nil {
		return nil, err
	}
	return eoa, nil
}

func (nak *NesterAccountKeeper) SetAccount(account EthAccount) error {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&account)
	if err != nil {
		return errors.Errorf("Failed to serialize: %s", err)
	}
	err = nak.state.Set(storage.StoreKey(prefixKey), dat)
	if err != nil {
		return errors.Errorf("Failed to update storage for account: %s", err)
	}
	return nil
}

func (nak *NesterAccountKeeper) RemoveAccount(account EthAccount) {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	nak.state.Delete(prefixKey)
}
