package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

const DefaultTimeout = 5 * time.Second

type ChainDriver interface {
}

func NewChainDriver(cfg *config.EthereumChainDriverConfig, logger *log.Logger, option *ChainDriverOption) (*ETHChainDriver, error) {

	client, err := ethclient.Dial(cfg.Connection)
	if err != nil {
		return nil, err
	}

	_, err = contract.NewLockRedeem(option.ContractAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to find contract at address ")
	}

	return &ETHChainDriver{
		cfg:             cfg,
		contract:        nil,
		client:          nil,
		logger:          logger,
		ContractAddress: option.ContractAddress,
		ContractABI:     option.ContractABI,
	}, nil
}

// defaultContext returns the context.ETHChainDriver to be used in requests against the Ethereum client
func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// Implements SContract interface
var _ ChainDriver = &ETHChainDriver{}

// ETHChainDriver provides the core fields required to interact with the Ethereum network. As of this moment (2019-08-21)
// it should only be used by validator nodes.
type ETHChainDriver struct {
	cfg             *config.EthereumChainDriverConfig
	client          *Client
	contract        *Contract
	logger          *log.Logger
	ContractAddress Address
	ContractABI     string
}

func (acc *ETHChainDriver) GetClient() *Client {
	if acc.client == nil {
		client, err := ethclient.Dial(acc.cfg.Connection)
		if err != nil {
			panic(err)
		}
		acc.client = client
	}
	return acc.client
}

func (acc *ETHChainDriver) GetContract() *Contract {
	client := acc.GetClient()

	if acc.contract == nil {
		ctrct, err := contract.NewLockRedeem(acc.ContractAddress, client)
		if err != nil {
			panic(err)
		}
		acc.contract = ctrct
	}
	return acc.contract
}

// Balance returns the current balance of the
func (acc ETHChainDriver) Balance(addr Address) (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.GetClient().BalanceAt(c, addr, nil)
}

// Nonce returns the nonce to use for our next transaction
func (acc ETHChainDriver) Nonce(addr Address) (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.GetClient().PendingNonceAt(c, addr)
}

// VerifyContract returns true if we can verify that the current contract matches the
func (acc ETHChainDriver) VerifyContract(vs []Address) (bool, error) {
	// 1. Make sure good IsValidator
	// 2. Make sure bad !IsValidator
	return false, ErrNotImplemented
}

func (acc *ETHChainDriver) CallOpts(addr Address) *CallOpts {
	return &CallOpts{
		Pending:     true,
		From:        addr,
		BlockNumber: nil,
		Context:     context.Background(),
	}
}

func (acc *ETHChainDriver) SignRedeem(fromaddr common.Address, redeemAmount *big.Int, recipient common.Address) (*Transaction, error) {

	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.GetClient().PendingNonceAt(c, fromaddr)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.GetClient().SuggestGasPrice(c)
	if err != nil {
		return nil, err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("sign", redeemAmount, recipient)
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, redeemAmount, gasLimit, gasPrice, bytesData)
	return tx, nil

}

func (acc *ETHChainDriver) PrepareUnsignedETHLock(addr common.Address, lockAmount *big.Int) ([]byte, error) {

	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.GetClient().PendingNonceAt(c, addr)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.GetClient().SuggestGasPrice(c)
	if err != nil {
		return nil, err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("lock")
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, lockAmount, gasLimit, gasPrice, bytesData)
	ts := types.Transactions{tx}
	rawTxBytes := ts.GetRlp(0)
	return rawTxBytes, nil
}

func (acc *ETHChainDriver) DecodeTransaction(rawBytes []byte) (*types.Transaction, error) {
	return DecodeTransaction(rawBytes)
}

func (acc *ETHChainDriver) GetTransactionMessage(tx *types.Transaction) (*types.Message, error) {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to convert Tx to message")
	}
	return &msg, nil
}

func (acc *ETHChainDriver) PrepareUnsignedETHRedeem(addr common.Address, lockAmount *big.Int) (*Transaction, error) {

	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.GetClient().PendingNonceAt(c, addr)
	if err != nil {
		return nil, err
	}
	fmt.Println(nonce, "for validator:", addr)
	gasLimit := uint64(6721974)
	gasPrice, err := acc.GetClient().SuggestGasPrice(c)
	if err != nil {
		return nil, err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("redeem", lockAmount)
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gasLimit, gasPrice, bytesData)

	return tx, nil
}

func (acc *ETHChainDriver) CheckFinality(txHash TransactionHash) (*types.Receipt, error) {

	result, err := acc.GetClient().TransactionReceipt(context.Background(), txHash)
	if err == nil {
		if result.Status == types.ReceiptStatusSuccessful {
			acc.logger.Info("Received TX Receipt for : ", txHash)
			return result, nil
		}
		if result.Status == types.ReceiptStatusFailed {
			acc.logger.Warn("Receipt not found ")
			//err := Error("Transaction not added to blockchain yet / Failed to obtain receipt")
			return nil, nil
		}
	}
	acc.logger.Error("Unable to connect to Ethereum :", err)
	return nil, err
}

func (acc *ETHChainDriver) ChainId() (*big.Int, error) {
	return acc.GetClient().ChainID(context.Background())
}

func (acc *ETHChainDriver) BroadcastTx(tx *types.Transaction) (TransactionHash, error) {

	_, _, err := acc.GetClient().TransactionByHash(context.Background(), tx.Hash())
	if err == nil {
		return tx.Hash(), nil
	}
	err = acc.GetClient().SendTransaction(context.Background(), tx)
	if err != nil {
		acc.logger.Error("Error connecting to Ethereum :", err)
		return tx.Hash(), err
	}
	acc.logger.Info("Transaction Broadcasted to Ethereum ", tx.Hash())
	return tx.Hash(), nil

}

func (acc *ETHChainDriver) ParseRedeem(data []byte) (req *RedeemRequest, err error) {
	//ethTx := &types.Transaction{}
	//err = rlp.DecodeBytes(data, ethTx)
	//if err != nil {
	//	return nil, errors.Wrap(err,"Unable to decode transaction")
	//}
	//
	//encodedData := ethTx.Data()
	//contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	//
	//r := make(map[string]interface{})
	//fmt.Println(contractAbi.Methods)
	//err = contractAbi.Unpack(&r, "db006a75", encodedData)
	//if err != nil {
	//	fmt.Println(r)
	//	return nil, errors.Wrap(err,"Unable to create Redeem Request")
	//}
	//TODO : Refactor to use UNPACK

	return ParseRedeem(data)
}

func (acc *ETHChainDriver) VerifyRedeem(validatorAddress common.Address, recepient common.Address) (bool, error) {
	instance := acc.GetContract()
	ok, err := instance.VerifyRedeem(acc.CallOpts(validatorAddress), recepient)
	if err != nil {
		return false, errors.Wrap(err, "Unable to connect to ethereum smart contract")
	}
	return ok, nil
}

func ParseRedeem(data []byte) (req *RedeemRequest, err error) {
	ss := strings.Split(hex.EncodeToString(data), "db006a75")
	if len(ss) == 0 {
		return nil, errors.New("Transaction does not have the required input data")
	}
	if len(ss[1]) < 64 {
		return nil, errors.New("Transaction data is invalid")
	}
	d, err := hex.DecodeString(ss[1][:64])
	if err != nil {
		return nil, err
	}
	amt := big.NewInt(0).SetBytes(d)
	return &RedeemRequest{Amount: amt}, nil
}

func ParseLock(data []byte) (req *LockRequest, err error) {
	tx, err := DecodeTransaction(data)
	if err != nil {
		return nil, err
	}
	return &LockRequest{Amount: tx.Value()}, nil
}

func DecodeTransaction(data []byte) (*types.Transaction, error) {
	tx := &types.Transaction{}
	err := rlp.DecodeBytes(data, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode Bytes")
	}
	return tx, nil
}
