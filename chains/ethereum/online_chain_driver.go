package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

const DefaultTimeout = 5 * time.Second


type ETHChainDriver struct {
	cfg             *config.EthereumChainDriverConfig
	client          *Client
	contract        Contract
	logger          *log.Logger
	ContractAddress Address
	ContractABI     string
	ContractType    ContractType
}

// Implements ChainDriver interface
var _ ChainDriver = &ETHChainDriver{}
func NewChainDriver(cfg *config.EthereumChainDriverConfig, logger *log.Logger, contractAddress common.Address, contractAbi string, contractType ContractType) (*ETHChainDriver, error) {

	client, err := ethclient.Dial(cfg.Connection)
	if err != nil {
		return nil, err
	}

	_, err = contract.NewLockRedeem(contractAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to find contract at address ")
	}

	return &ETHChainDriver{
		cfg:             cfg,
		contract:        nil,
		client:          nil,
		logger:          logger,
		ContractAddress: contractAddress,
		ContractABI:     contractAbi,
		ContractType:    contractType,
	}, nil
}

// defaultContext returns the context.ETHChainDriver to be used in requests against the Ethereum client
func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}



// ETHChainDriver provides the core fields required to interact with the Ethereum network. As of this moment (2019-08-21)
// it should only be used by validator nodes.

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
// GetContract returns instance of an already deployed ETH/ERC LockRedeem Contract
func (acc *ETHChainDriver) GetContract() Contract {
	client := acc.GetClient()
	if acc.contract == nil {
		if acc.ContractType == ETH {
			ctrct, err := contract.NewLockRedeem(acc.ContractAddress, client)
			if err != nil {
				panic(err)
			}
			ctr := GetETHContract(*ctrct)
			acc.contract = ctr
		} else if acc.ContractType == ERC {
			ctrct, err := contract.NewLockRedeemERC(acc.ContractAddress, client)
			if err != nil {
				panic(err)
			}
			ctr := GetERCContract(*ctrct)
			acc.contract = ctr

		}
	}
	return acc.contract
}

// Balance returns the current balance of address
func (acc ETHChainDriver) Balance(addr Address) (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.GetClient().BalanceAt(c, addr, nil)
}

// Nonce returns the nonce of the address
func (acc ETHChainDriver) Nonce(addr Address) (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.GetClient().PendingNonceAt(c, addr)
}
// ChainID returns the ID for the current ethereum chain (Mainet/Ropsten/Ganache)
func (acc *ETHChainDriver) ChainId() (*big.Int, error) {
	return acc.GetClient().ChainID(context.Background())
}
// CallOpts creates a CallOpts object for contract calls (Only call no write )
func (acc *ETHChainDriver) CallOpts(addr Address) *CallOpts {
	return &CallOpts{
		Pending:     true,
		From:        addr,
		BlockNumber: nil,
		Context:     context.Background(),
	}
}

//SignRedeem creates an Ethereum transaction used by Validators to sign a redeem Request
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
	bytesData, err := contractAbi.Pack("sign", redeemAmount, recipient)
	if err != nil {
		return nil, err
	}
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gasLimit, gasPrice, bytesData)
	return tx, nil

}

// PrepareUnsignedETHLock creates a raw Transaction to lock ether.
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

// DecodeTransaction is an online wrapper for DecodeTransaction
func (acc *ETHChainDriver) DecodeTransaction(rawBytes []byte) (*types.Transaction, error) {
	return DecodeTransaction(rawBytes)
}

// GetTransactionMessage takes a trasaction as input and returns the message
func (acc *ETHChainDriver) GetTransactionMessage(tx *types.Transaction) (*types.Message, error) {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to convert Tx to message")
	}
	return &msg, nil
}

// CheckFinality verifies the finality of a transaction on the ethereum blockchain , waits for 12 block confirmations
func (acc *ETHChainDriver) CheckFinality(txHash TransactionHash) (*types.Receipt, error) {

	result, err := acc.GetClient().TransactionReceipt(context.Background(), txHash)
	fmt.Println("Client :", acc.GetClient())
	if err == nil {
		if result.Status == types.ReceiptStatusSuccessful {
			latestHeader, err := acc.client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return nil, errors.Wrap(err, "Unable to extract latest header")
			}
			diff := big.NewInt(0).Sub(latestHeader.Number, result.BlockNumber)
			if big.NewInt(12).Cmp(diff) < 0 {
				return nil, errors.New("Waiting for confirmation . Current Block Confirmations : " + diff.String())
			}
			return result, nil
		}
		if result.Status == types.ReceiptStatusFailed {
			acc.logger.Warn("Receipt not found ")
			b, _ := result.MarshalJSON()
			acc.logger.Error(string(b))
			return nil, nil
		}
	}
	acc.logger.Error("Transaction not added to Block yet :", err)
	return nil, err
}


// BroadcastTx takes a signed transaction as input and broadcasts it to the network
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
	acc.logger.Info("Transaction Broadcasted to Ethereum ", tx.Hash().Hex())
	return tx.Hash(), nil

}

// ParseERC20Redeem is the online wrapper for ParseERC20RedeemParams
func (Acc *ETHChainDriver) ParseERC20Redeem(rawTx []byte, lockredeemERCAbi string) (*RedeemErcRequest, error) {
	return ParseERC20RedeemParams(rawTx, lockredeemERCAbi)
}
// ParseERC20Redeem is the online wrapper for ParseRedeem
func (acc *ETHChainDriver) ParseRedeem(data []byte, abi string) (req *RedeemRequest, err error) {
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

	return ParseRedeem(data, abi)
}
// VerifyRedeem verifies if the Redeem request is Completed
func (acc *ETHChainDriver) VerifyRedeem(validatorAddress common.Address, recipient common.Address) (bool, error) {
	instance := acc.GetContract()
	ok, err := instance.VerifyRedeem(acc.CallOpts(validatorAddress), recipient)
	if err != nil {
		return false, errors.Wrap(err, "Unable to connect to ethereum smart contract")
	}
	return ok, nil
}
// HasValidatorSigned takes validator address and recipient address as input and verifies if the validator has already signed
func (acc *ETHChainDriver) HasValidatorSigned(validatorAddress common.Address, recipient common.Address) (bool, error) {
	instance := acc.GetContract()

	return instance.HasValidatorSigned(acc.CallOpts(validatorAddress), recipient)
}
