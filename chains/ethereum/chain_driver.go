package ethereum

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
)

type ChainDriver interface {
        //OnlineLock (rawTx []byte) (*big.Int, error)
       // OfflineLock (pubKey *ecdsa.PublicKey,lockAmount *big.Int) ([]byte,error)
	    PrepareUnsignedETHLock (pubKey *ecdsa.PublicKey,lockAmount *big.Int) ([]byte,error)
	    CheckFinality(txHash TransactionHash)
	    BroadcastTx (tx *types.Transaction) (TransactionHash,error)
        //ValidatorSignRedeem (wei *big.Int, recipient common.Address) (*Transaction,error)
}

func (acc *EthereumChainDriver) ValidatorSignRedeem(wei *big.Int, recipient common.Address) (*Transaction,error){
	acc.logger.Info("Validator Signed Redeem")
	opts := acc.TransactOpts()
	opts.From = acc.Address()
	var redeemid = new(big.Int)
	redeemid.SetString("2", 10)
	return  acc.Contract.Sign(opts,redeemid,wei, recipient)
}




func (acc *EthereumChainDriver) PrepareUnsignedETHLock (pubKey *ecdsa.PublicKey,lockAmount *big.Int) ([]byte,error){
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	c, cancel := defaultContext()
	defer cancel()
	nonce,err := acc.Client.PendingNonceAt(c,fromAddress)
	if err != nil {
		return nil,err
	}
	gasLimit := uint64(6721974)
	gasPrice, err := acc.Client.SuggestGasPrice(c)
	if err != nil {
		return nil,err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("lock")
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, lockAmount, gasLimit, gasPrice, bytesData)
	ts := types.Transactions{tx}
	rawTxBytes := ts.GetRlp(0)
	return rawTxBytes,nil
}

func (acc *EthereumChainDriver) CheckFinality (txHash TransactionHash) (*types.Receipt,error) {

		result, err := acc.Client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if result.Status == types.ReceiptStatusSuccessful {
		    acc.logger.Info("Received TX Receipt for : ",txHash)
			return result, nil
			}
			if result.Status == types.ReceiptStatusFailed {
				acc.logger.Warn("Receipt not found ")
				//err := Error("Transaction not added to blockchain yet / Failed to obtain receipt")
				return nil,nil
			}
		}
		acc.logger.Error("Unable to connect to Ethereum :" ,err)
		return nil,err
		}

func (acc *EthereumChainDriver) BroadcastTx (tx *types.Transaction) (TransactionHash,error){

	err := acc.Client.SendTransaction(context.Background(), tx)
	if err!= nil {
		acc.logger.Error("Error connecting to ETHEREUM :",err)
		return tx.Hash(),err
	}
	acc.logger.Info("Transaction Broadcasted to Ethereum ",tx.Hash())
	return tx.Hash(),nil

}


