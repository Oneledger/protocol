package ethereum

import (
	"context"
	"math/big"
	"strings"
	"time"

	ethereum2 "github.com/ethereum/go-ethereum"
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
const gasLimit = 700000

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
	gasLimit := uint64(gasLimit)
	gasPrice, err := acc.GetClient().SuggestGasPrice(c)
	gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(0).Div(gasPrice, big.NewInt(2)))
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
	gasLimit := uint64(gasLimit)
	gasPrice, err := acc.GetClient().SuggestGasPrice(c)
	if err != nil {
		return nil, err
	}
	contractAbi, _ := abi.JSON(strings.NewReader(acc.ContractABI))
	bytesData, _ := contractAbi.Pack("lock")
	toAddress := acc.ContractAddress
	tx := types.NewTransaction(nonce, toAddress, lockAmount, gasLimit, gasPrice, bytesData)
	ts := types.Transactions{tx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	return rawTxBytes, nil
}

// DecodeTransaction is an online wrapper for DecodeTransaction
func (acc *ETHChainDriver) DecodeTransaction(rawBytes []byte) (*types.Transaction, error) {
	return DecodeTransaction(rawBytes)
}

// GetTransactionMessage takes a trasaction as input and returns the message
func (acc *ETHChainDriver) GetTransactionMessage(tx *types.Transaction) (*types.Message, error) {
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to convert Tx to message")
	}
	return &msg, nil
}

// CheckFinalityStatus verifies the finality of a transaction on the ethereum blockchain , waits for 12 block confirmations
func (acc *ETHChainDriver) CheckFinality(txHash TransactionHash, blockConfirmation int64) CheckFinalityStatus {
	result, err := acc.GetClient().TransactionReceipt(context.Background(), txHash)
	if err != nil {
		acc.logger.Debug("Transaction not added to Block yet :", err)
		return TransactionNotMined
	}
	if result.Status == types.ReceiptStatusSuccessful {
		latestHeader, err := acc.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			acc.logger.Debug("Original Receipt Successful , But unable to get Latest Header (Connection Problem)", err)
			return UnabletoGetHeader
		}
		diff := big.NewInt(0).Sub(latestHeader.Number, result.BlockNumber)
		if big.NewInt(blockConfirmation).Cmp(diff) > 0 {
			acc.logger.Debug("Waiting for confirmation . Current Block Confirmations : " + diff.String())
			return NotEnoughConfirmations
		}
		txHeaderCalculated, err := acc.client.HeaderByNumber(context.Background(), result.BlockNumber)
		if err != nil {
			acc.logger.Debug("Block Confirmed ,Error in getting Header of Original block at T - "+diff.String(), "  ", err)
			return TxBlockNotFound
		}
		if txHeaderCalculated.Hash() != result.BlockHash {
			acc.logger.Debug("BlockHash Check Failed ,Uncle block Detected")
			return BlockHashFailed
		}
		return TXSuccess
	}

	acc.logger.Debug("Receipt status failed ")
	b, _ := result.MarshalJSON()
	acc.logger.Debug(string(b))
	return ReciptNotFound

}

// Bool value is for recipt status
// Err is to handle error , we need to ignore NotFound and wait for the tx
func (acc *ETHChainDriver) VerifyReceipt(txHash TransactionHash) (VerifyReceiptStatus, error) {

	result, err := acc.GetClient().TransactionReceipt(context.Background(), txHash)
	if err == ethereum2.NotFound && result == nil {
		return NotFound, nil
	}
	if err != nil {
		return Other, err
	}
	if result.Status == types.ReceiptStatusFailed {
		return Failed, nil
	}
	if result.Status == types.ReceiptStatusSuccessful {
		return Found, nil
	}
	// Returning generic result ( no err )
	return Other, nil
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
	//acc.logger.Info("Transaction Broadcasted to Ethereum ", tx.Hash().Hex())
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
//BeforeRedeem : -1  (amount == 0 and until = 0)
//Ongoing : 0  (amount > 0 and until > block.number)
//Success : 1  (amount = 0 and until < block.number)
//Expired : 2  (amount > 0 and until < block.number)
//ErrConnecting  : - 2  (Ethereum connection cannot be established)

func (acc *ETHChainDriver) VerifyRedeem(validatorAddress common.Address, recipient common.Address) RedeemStatus {
	instance := acc.GetContract()
	redeemStatus, err := instance.VerifyRedeem(acc.CallOpts(validatorAddress), recipient)
	if err != nil {
		return ErrorConnecting
	}
	return RedeemStatus(redeemStatus)
}

// HasValidatorSigned takes validator address and recipient address as input and verifies if the validator has already signed
func (acc *ETHChainDriver) HasValidatorSigned(validatorAddress common.Address, recipient common.Address) (bool, error) {
	instance := acc.GetContract()
	return instance.HasValidatorSigned(acc.CallOpts(validatorAddress), recipient)
}
