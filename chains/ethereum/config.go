package ethereum

import (
	"crypto/ecdsa"
	"path/filepath"

	"github.com/ethereum/go-ethereum/ethclient"
)

const defaultKeyLocation = "eth/key.json"

type Config struct {
	// Path to the ethereum key
	ContractABI       string   `toml: contractAbi" desc:"AVI for the contract"`
	Connection        string   `toml:"connection" desc:"Connection string to the Ethereum node"`
	KeyLocation       string   `toml:"key" desc:"Relative path to the Ethereum key. Can be left blank"`
	ContractAddress   string   `toml:"contract_address" desc:"Address to the ethereum LockRedeem contract. This should not depend on  things."`
	InitialValidators []string `toml:"initial_validators" desc:"Addresses of the initial validators"`
}

func DefaultEthConfig() *Config {
	return &Config{
		Connection:  "https://ropsten.infura.io/v3/{API_KEY}",
		KeyLocation: defaultKeyLocation,
	}
}

func CreateEthConfig(connection string,keylocation string,address string,validators []string) *Config{
	return &Config{
		Connection:connection,
		KeyLocation:keylocation,
		ContractAddress:address,
        InitialValidators:validators,
	}
}

// Returns an instance of the private and public keypair with the given root directory
func (cfg *Config) Key(rootDir string) (*ecdsa.PrivateKey, error) {
	switch {
	case filepath.IsAbs(cfg.KeyLocation) && filepath.Ext(cfg.KeyLocation) != "":
		return ReadUnsafeKeyDump(cfg.KeyLocation)
	case filepath.IsAbs(cfg.KeyLocation):
		return ReadUnsafeKeyDump(cfg.KeyLocation)
	}

	return ReadUnsafeKeyDump(filepath.Join(rootDir, cfg.KeyLocation))
}

func (cfg *Config) Client() (*Client, error) {
	return ethclient.Dial(cfg.Connection)
}
