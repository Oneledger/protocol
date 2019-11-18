package ethereum

import (
	"context"
	"crypto/ecdsa"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/data/keys"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

const DefaultTimeout = 5 * time.Second

func NewEthereumChainDriver(cfg *config.EthereumChainDriverConfig ,logger * log.Logger ,ethPrivKey *keys.PrivateKey) (*EthereumChainDriver, error) {

	// So first we need to grab the current address
	contractAddrSlice, err := hexutil.Decode(cfg.ContractAddress)
	if err != nil {
		logger.Error("failed to decode given contract address")
		return nil, err
	}

	var contractAddr Address
	copy(contractAddr[:], contractAddrSlice)

	client, err := ethclient.Dial(cfg.Connection)
	if err != nil {
		logger.Error("failed to dial the given ethereum connection")
		return nil, err
	}
    //TODO : Function to return AUTH object
    //TODO : Function to load contract from specific Address .
    //Todo : Deploy Smart contract and return Deployed addresss
	ctrct, err := contract.NewLockRedeem(contractAddr, client)
	if err != nil {
		logger.Error("failed to create new contract")
		return nil, err
	}

	//keyPath := filepath.Join(rootDir, cfg.KeyLocation)
	// READ FROM KEYSTORE ,ganache does not have a keystore
	//privKey, err := ReadUnsafeKeyDump(keyPath)
	//privKey := privateKey
	if err != nil {
		logger.Error("Failed to create private key")
		return nil, err
	}

	return &EthereumChainDriver{
		Contract:        ctrct,
		Client:          client,
		logger:          logger,
		ContractAddress: contractAddr,
		ContractABI:     cfg.ContractABI,
		PrivateKey:      ethPrivKey,
	}, nil
}

// defaultContext returns the context.EthereumChainDriver to be used in requests against the Ethereum client
func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// Implements SContract interface
// var _ SContract = &EthereumChainDriver{}

// EthereumChainDriver provides the core fields required to interact with the Ethereum network. As of this moment (2019-08-21)
// it should only be used by validator nodes.
type EthereumChainDriver struct {
	PrivateKey      *keys.PrivateKey
	Client          *Client
	Contract        *Contract
	ContractAddress Address
	logger          *log.Logger
	ContractABI     string
}

// Balance returns the current balance of the
func (acc EthereumChainDriver) Balance() (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.BalanceAt(c, acc.Address(), nil)
}

// Nonce returns the nonce to use for our next transaction
func (acc EthereumChainDriver) Nonce() (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.PendingNonceAt(c, acc.Address())
}



func (acc EthereumChainDriver) IsContract() {

}


// VerifyContract returns true if we can verify that the current contract matches the
func (acc EthereumChainDriver) VerifyContract(vs []Address) (bool, error) {
	// 1. Make sure good IsValidator
	// 2. Make sure bad !IsValidator
	return false, ErrNotImplemented
}

func (acc EthereumChainDriver) isContract() (bool, error) {
	return true, nil
}

func (acc EthereumChainDriver) Address() Address {
	return crypto.PubkeyToAddress(acc.PublicKey())
}


// TransactOpts() returns a new transactor with the EthereumChainDriver struct's private key
func (acc *EthereumChainDriver) TransactOpts() *TransactOpts {
	return bind.NewKeyedTransactor(acc.PrivateKey)
}

func (acc *EthereumChainDriver) CallOpts() *CallOpts {
	return &CallOpts{
		Pending:     true,
		From:        acc.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}
}
