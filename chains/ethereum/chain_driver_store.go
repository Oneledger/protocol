package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

const DefaultTimeout = 5 * time.Second

func NewEthereumChainDriver(cfg *config.EthereumChainDriverConfig, logger *log.Logger, option *ChainDriverOption) (*EthereumChainDriver, error) {

	client, err := ethclient.Dial(cfg.Connection)
	if err != nil {
		logger.Error("failed to dial the given ethereum connection", err)
		return nil, err
	}
	//TODO : Function to return AUTH object
	//TODO : Function to load contract from specific Address .
	//Todo : Deploy Smart contract and return Deployed addresss
	ctrct, err := contract.NewLockRedeem(option.ContractAddress, client)
	if err != nil {
		logger.Error("failed to create new contract")
		return nil, err
	}

	//keyPath := filepath.Join(rootDir, cfg.KeyLocation)
	// READ FROM KEYSTORE ,ganache does not have a keystore
	//privKey, err := ReadUnsafeKeyDump(keyPath)
	//privKey := privateKey
	//if err != nil {
	//	logger.Error("Failed to create private key")
	//	return nil, err
	//}

	return &EthereumChainDriver{
		Contract:        ctrct,
		Client:          client,
		logger:          logger,
		ContractAddress: option.ContractAddress,
		ContractABI:     option.ContractABI,
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
	Client          *Client
	Contract        *Contract
	ContractAddress Address
	logger          *log.Logger
	ContractABI     string
}

// Balance returns the current balance of the
func (acc EthereumChainDriver) Balance(addr Address) (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.BalanceAt(c, addr, nil)
}

// Nonce returns the nonce to use for our next transaction
func (acc EthereumChainDriver) Nonce(addr Address) (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.PendingNonceAt(c, addr)
}

// VerifyContract returns true if we can verify that the current contract matches the
func (acc EthereumChainDriver) VerifyContract(vs []Address) (bool, error) {
	// 1. Make sure good IsValidator
	// 2. Make sure bad !IsValidator
	return false, ErrNotImplemented
}

func (acc EthereumChainDriver) isContract() (bool, error) {
	panic("implement me")
	return true, nil
}

func (acc *EthereumChainDriver) CallOpts(addr Address) *CallOpts {
	return &CallOpts{
		Pending:     true,
		From:        addr,
		BlockNumber: nil,
		Context:     context.Background(),
	}
}
