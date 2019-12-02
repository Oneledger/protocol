package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

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
		logger.Error("failed to dial the given ethereum connection", err)
		return nil, err
	}

	ctrct, err := contract.NewLockRedeem(option.ContractAddress, client)
	if err != nil {
		logger.Error("failed to create new contract")
		return nil, err
	}

	return &ETHChainDriver{
		Contract:        ctrct,
		Client:          client,
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
	Client          *Client
	Contract        *Contract
	ContractAddress Address
	logger          *log.Logger
	ContractABI     string
}

// Balance returns the current balance of the
func (acc ETHChainDriver) Balance(addr Address) (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.BalanceAt(c, addr, nil)
}

// Nonce returns the nonce to use for our next transaction
func (acc ETHChainDriver) Nonce(addr Address) (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.PendingNonceAt(c, addr)
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

func (acc *ETHChainDriver) SignRedeem(addr common.Address,redeemAmount *big.Int,recipient common.Address) (*Transaction, error) {


	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.Client.PendingNonceAt(c, addr)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.Client.SuggestGasPrice(c)
	if err != nil {
		return nil, err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("sign",redeemAmount,recipient)
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress,redeemAmount, gasLimit, gasPrice, bytesData)
	return tx,nil


}

func (acc *ETHChainDriver) PrepareUnsignedETHLock(addr common.Address, lockAmount *big.Int) ([]byte, error) {

	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.Client.PendingNonceAt(c, addr)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.Client.SuggestGasPrice(c)
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

func (acc *ETHChainDriver) PrepareUnsignedETHRedeem(addr common.Address, lockAmount *big.Int) ( *Transaction, error) {

	c, cancel := defaultContext()
	defer cancel()
	nonce, err := acc.Client.PendingNonceAt(c, addr)
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.Client.SuggestGasPrice(c)
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

	result, err := acc.Client.TransactionReceipt(context.Background(), txHash)
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

func (acc *ETHChainDriver) BroadcastTx(tx *types.Transaction) (TransactionHash, error) {

	_, _, err := acc.Client.TransactionByHash(context.Background(), tx.Hash())
	if err == nil {
		return tx.Hash(), nil
	}
	err = acc.Client.SendTransaction(context.Background(), tx)
	if err != nil {
		acc.logger.Error("Error connecting to ETHEREUM :", err)
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

	ss := strings.Split(hex.EncodeToString(data), "db006a75")
	d, _ := hex.DecodeString(ss[1][:64])
	amt := big.NewInt(0).SetBytes(d)

    return &RedeemRequest{Amount: amt},nil
}