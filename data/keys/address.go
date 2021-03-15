package keys

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"fmt"
	"strings"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/utils"
)

//Address to be used as to reference a key-pair.
type Address []byte

func (a Address) String() string {
	return utils.PrefixAddress(hex.EncodeToString(a))
}

func (a Address) Bytes() []byte {
	return a
}

func (a Address) Equal(b Address) bool {
	return bytes.Equal(a, b)
}

func (a Address) Humanize() string {
	return utils.PrefixAddress(strings.ToLower(hex.EncodeToString(a)))
}

func (a Address) Err() error {
	switch {
	case len(a) == 0:
		return errors.New("address is empty")
	case len(a) != 20:
		return errors.New(fmt.Sprintf("address is the incorrect length: must be 20-bytes (40 hex characters after %s prefix)", utils.AddrPrefix))
	}
	return nil
}

var _ encoding.TextMarshaler = Address{}
var _ encoding.TextUnmarshaler = &Address{}

// MarshalText returns the text form for an Address. It returns a byteslice containing the hex encoding
// of the address including the 0lt prefix
func (a Address) MarshalText() ([]byte, error) {
	addrHex := hex.EncodeToString(a)
	return []byte(utils.PrefixAddress(addrHex)), nil
}

// UnmarshalText decodes the given text in a byteslice of characters,
// and works regardless of whether the 0lt prefix is present or not.
func (a *Address) UnmarshalText(text []byte) error {
	if a == nil {
		*a = Address{}
	}
	addrStr := string(text)

	// Cut off the hex prefix if it exists
	addrStr = utils.TrimAddress(addrStr)

	addrRaw, err := hex.DecodeString(addrStr)
	if err != nil {
		return errors.Wrap(err, "address text unmarshal failed, not a hex address")
	}

	*a = addrRaw

	return nil
}

// EthAccount implements the keys.Account interface and embeds with code hash for
// the contract
type EthAccount struct {
	Address  Address
	CodeHash []byte
	Coins    map[string]balance.Coin
	Sequence uint64
}

func NewEthAccount(addr Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		CodeHash: ethcrypto.Keccak256(nil),
		Coins:    make(map[string]balance.Coin),
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

func (acc EthAccount) Balance(currency balance.Currency) *balance.Amount {
	return acc.Coins[currency.Name].Amount
}

func (acc *EthAccount) AddBalance(coin balance.Coin) {
	balance := acc.Coins[coin.Currency.Name]
	acc.Coins[coin.Currency.Name] = balance.Plus(coin)
}

func (acc *EthAccount) SubBalance(coin balance.Coin) {
	balance := acc.Coins[coin.Currency.Name]
	newCoin, _ := balance.Minus(coin)
	acc.Coins[coin.Currency.Name] = newCoin
}

func (acc *EthAccount) SetBalance(coin balance.Coin) {
	acc.Coins[coin.Currency.Name] = coin
}
