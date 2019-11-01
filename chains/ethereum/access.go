package ethereum

import (
	"context"
	"crypto/ecdsa"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"math/big"
	"time"
)

const DefaultTimeout = 5 * time.Second

func NewAccess(rootDir string, cfg *Config, logger *log.Logger) (*Access, error) {
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
	privKey,err := crypto.HexToECDSA("2d3a5b8530d3024501885220b8ef645876bfc816a198a7f89e4b1458fc1644a8")
	if err != nil {
		logger.Error("Failed to create private key")
		return nil, err
	}

	return &Access{
		Contract:        ctrct,
		PrivateKey:      privKey,
		Client:          client,
		logger:          logger,
		ContractAddress: contractAddr,
	}, nil
}

// defaultContext returns the context.Access to be used in requests against the Ethereum client
func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

// Implements SContract interface
// var _ SContract = &Access{}

// Access provides the core fields required to interact with the Ethereum network. As of this moment (2019-08-21)
// it should only be used by validator nodes.
type Access struct {
	PrivateKey      *ecdsa.PrivateKey
	Client          *Client
	Contract        *Contract
	ContractAddress Address
	logger          *log.Logger
}

// Balance returns the current balance of the
func (acc Access) Balance() (*big.Int, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.BalanceAt(c, acc.Address(), nil)
}

// Nonce returns the nonce to use for our next transaction
func (acc Access) Nonce() (uint64, error) {
	c, cancel := defaultContext()
	defer cancel()
	return acc.Client.PendingNonceAt(c, acc.Address())
}



func (acc Access) IsContract() {

}

func (acc Access) PublicKey() ecdsa.PublicKey {
	return acc.PrivateKey.PublicKey
}

// VerifyContract returns true if we can verify that the current contract matches the
func (acc Access) VerifyContract(vs []Address) (bool, error) {
	// 1. Make sure good IsValidator
	// 2. Make sure bad !IsValidator
	return false, ErrNotImplemented
}

func (acc Access) isContract() (bool, error) {
	return true, nil
}

func (acc Access) Address() Address {
	return crypto.PubkeyToAddress(acc.PublicKey())
}
//addr keys.Address,
//func (acc *Access) Lock(wei *big.Int) (*Transaction, error) {
//	acc.logger.Info("locked wei")
//	opts := acc.TransactOpts()
//	opts.Value = wei
//	opts.From = acc.Address()
//	//var oltAddrAsEthAddr Address
//	//copy(oltAddrAsEthAddr[:], addr)
//	return acc.Contract.Lock(opts)
//}



func (acc *Access) Sign(wei *big.Int, recipient common.Address) (*Transaction,error){
	acc.logger.Info("Validator Signed Redeem")
	opts := acc.TransactOpts()
	opts.From = acc.Address()
	var redeemid = new(big.Int)
	redeemid.SetString("2", 10)
	return  acc.Contract.Sign(opts,redeemid,wei, recipient)
}


// LockFromSignedTx accepts a OneLedger address, and a signed transaction
// It handles the ethereum
func (acc *Access) LockFromSignedTx(oltAddr Address, rawTx []byte) (*big.Int, error) {
	// First accept the rawTx
	var tx *Transaction
	err := rlp.DecodeBytes(rawTx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode provided transaction bytes")
	}

	err = acc.Client.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, err
	}
	return tx.Value(), nil
}


func (acc *Access) CheckLock(amount int64,oltethadress string)(bool,error){
	acc.logger.Info("Checking Locked amount")
	opts := acc.CallOpts()
	opts.From = acc.Address()
	balance,err := acc.Contract.GetLockedBalance(opts,oltethadress)
	if err != nil {
		return false, errors.Wrap(err, "Unable to get locked balance for specified address")
	}
	currentBalance := balance.Int64()
	var previousBalance int64 //Save Balance of account on chain state ;
    return currentBalance-previousBalance == amount,nil
	//update chainstate if true else return false
	//Call mint function
}

/*
Methods for generating valid Options for calling contracts, creating transactions, etc.

These methods generally fill in the options with the next valid nonce and injects the private key inside
*/

// TransactOpts() returns a new transactor with the Access struct's private key
func (acc *Access) TransactOpts() *TransactOpts {
	return bind.NewKeyedTransactor(acc.PrivateKey)
}

func (acc *Access) CallOpts() *CallOpts {
	return &CallOpts{
		// Whyethere or not we want to operate on pending state or last known state
		Pending:     true,
		From:        acc.Address(),
		BlockNumber: nil,
		Context:     context.Background(),
	}
}
