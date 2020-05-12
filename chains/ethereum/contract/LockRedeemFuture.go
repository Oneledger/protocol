// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LockRedeemFutureABI is the input ABI used to generate the binding from.
const LockRedeemFutureABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_old_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"noofValidatorsinold\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":false,\"inputs\":[],\"name\":\"MigrateFromOld\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMigrationCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isActive\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LockRedeemFutureFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemFutureFuncSigs = map[string]string{
	"587ab37e": "MigrateFromOld()",
	"cdaf4028": "getMigrationCount()",
	"287cc96b": "getTotalEthBalance()",
	"22f3e2d4": "isActive()",
	"5d593f8d": "numValidators()",
	"fa52c7d8": "validators(address)",
}

// LockRedeemFutureBin is the compiled bytecode used for deploying new contracts.
var LockRedeemFutureBin = "0x60806040526000805460ff1916815560326001556004819055600581905560065534801561002c57600080fd5b5060405161027f38038061027f8339818101604052604081101561004f57600080fd5b508051602090910151600380546001600160a01b0319166001600160a01b0384161781556002820204600101600555506101ef9050806100906000396000f3fe6080604052600436106100555760003560e01c806322f3e2d414610074578063287cc96b1461009d578063587ab37e146100c45780635d593f8d146100db578063cdaf4028146100f0578063fa52c7d814610105575b6005546006541461006557600080fd5b6000805460ff19166001179055005b34801561008057600080fd5b50610089610138565b604080519115158252519081900360200190f35b3480156100a957600080fd5b506100b2610141565b60408051918252519081900360200190f35b3480156100d057600080fd5b506100d9610145565b005b3480156100e757600080fd5b506100b2610170565b3480156100fc57600080fd5b506100b2610176565b34801561011157600080fd5b506100b26004803603602081101561012857600080fd5b50356001600160a01b031661017c565b60005460ff1690565b4790565b6003546001600160a01b0316331461015c57600080fd5b60068054600101905561016e3261018e565b565b60045481565b60065490565b60026020526000908152604090205481565b600180546001600160a01b0390921660009081526002602052604090209190915560048054909101905556fea265627a7a7231582039b2d24c897e3c4d7ff269aa6e0bf74e6cec4eaa3c8cc8d3e4bf746a8dd93d3f64736f6c63430005100032"

// DeployLockRedeemFuture deploys a new Ethereum contract, binding an instance of LockRedeemFuture to it.
func DeployLockRedeemFuture(auth *bind.TransactOpts, backend bind.ContractBackend, _old_contract common.Address, noofValidatorsinold *big.Int) (common.Address, *types.Transaction, *LockRedeemFuture, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemFutureABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockRedeemFutureBin), backend, _old_contract, noofValidatorsinold)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockRedeemFuture{LockRedeemFutureCaller: LockRedeemFutureCaller{contract: contract}, LockRedeemFutureTransactor: LockRedeemFutureTransactor{contract: contract}, LockRedeemFutureFilterer: LockRedeemFutureFilterer{contract: contract}}, nil
}

// LockRedeemFuture is an auto generated Go binding around an Ethereum contract.
type LockRedeemFuture struct {
	LockRedeemFutureCaller     // Read-only binding to the contract
	LockRedeemFutureTransactor // Write-only binding to the contract
	LockRedeemFutureFilterer   // Log filterer for contract events
}

// LockRedeemFutureCaller is an auto generated read-only Go binding around an Ethereum contract.
type LockRedeemFutureCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemFutureTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LockRedeemFutureTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemFutureFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockRedeemFutureFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemFutureSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockRedeemFutureSession struct {
	Contract     *LockRedeemFuture // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockRedeemFutureCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockRedeemFutureCallerSession struct {
	Contract *LockRedeemFutureCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// LockRedeemFutureTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockRedeemFutureTransactorSession struct {
	Contract     *LockRedeemFutureTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// LockRedeemFutureRaw is an auto generated low-level Go binding around an Ethereum contract.
type LockRedeemFutureRaw struct {
	Contract *LockRedeemFuture // Generic contract binding to access the raw methods on
}

// LockRedeemFutureCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockRedeemFutureCallerRaw struct {
	Contract *LockRedeemFutureCaller // Generic read-only contract binding to access the raw methods on
}

// LockRedeemFutureTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockRedeemFutureTransactorRaw struct {
	Contract *LockRedeemFutureTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLockRedeemFuture creates a new instance of LockRedeemFuture, bound to a specific deployed contract.
func NewLockRedeemFuture(address common.Address, backend bind.ContractBackend) (*LockRedeemFuture, error) {
	contract, err := bindLockRedeemFuture(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockRedeemFuture{LockRedeemFutureCaller: LockRedeemFutureCaller{contract: contract}, LockRedeemFutureTransactor: LockRedeemFutureTransactor{contract: contract}, LockRedeemFutureFilterer: LockRedeemFutureFilterer{contract: contract}}, nil
}

// NewLockRedeemFutureCaller creates a new read-only instance of LockRedeemFuture, bound to a specific deployed contract.
func NewLockRedeemFutureCaller(address common.Address, caller bind.ContractCaller) (*LockRedeemFutureCaller, error) {
	contract, err := bindLockRedeemFuture(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemFutureCaller{contract: contract}, nil
}

// NewLockRedeemFutureTransactor creates a new write-only instance of LockRedeemFuture, bound to a specific deployed contract.
func NewLockRedeemFutureTransactor(address common.Address, transactor bind.ContractTransactor) (*LockRedeemFutureTransactor, error) {
	contract, err := bindLockRedeemFuture(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemFutureTransactor{contract: contract}, nil
}

// NewLockRedeemFutureFilterer creates a new log filterer instance of LockRedeemFuture, bound to a specific deployed contract.
func NewLockRedeemFutureFilterer(address common.Address, filterer bind.ContractFilterer) (*LockRedeemFutureFilterer, error) {
	contract, err := bindLockRedeemFuture(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockRedeemFutureFilterer{contract: contract}, nil
}

// bindLockRedeemFuture binds a generic wrapper to an already deployed contract.
func bindLockRedeemFuture(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemFutureABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemFuture *LockRedeemFutureRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeemFuture.Contract.LockRedeemFutureCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemFuture *LockRedeemFutureRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.LockRedeemFutureTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemFuture *LockRedeemFutureRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.LockRedeemFutureTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemFuture *LockRedeemFutureCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeemFuture.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemFuture *LockRedeemFutureTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemFuture *LockRedeemFutureTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.contract.Transact(opts, method, params...)
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCaller) GetMigrationCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemFuture.contract.Call(opts, out, "getMigrationCount")
	return *ret0, err
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureSession) GetMigrationCount() (*big.Int, error) {
	return _LockRedeemFuture.Contract.GetMigrationCount(&_LockRedeemFuture.CallOpts)
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCallerSession) GetMigrationCount() (*big.Int, error) {
	return _LockRedeemFuture.Contract.GetMigrationCount(&_LockRedeemFuture.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCaller) GetTotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemFuture.contract.Call(opts, out, "getTotalEthBalance")
	return *ret0, err
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureSession) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeemFuture.Contract.GetTotalEthBalance(&_LockRedeemFuture.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCallerSession) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeemFuture.Contract.GetTotalEthBalance(&_LockRedeemFuture.CallOpts)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() constant returns(bool)
func (_LockRedeemFuture *LockRedeemFutureCaller) IsActive(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemFuture.contract.Call(opts, out, "isActive")
	return *ret0, err
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() constant returns(bool)
func (_LockRedeemFuture *LockRedeemFutureSession) IsActive() (bool, error) {
	return _LockRedeemFuture.Contract.IsActive(&_LockRedeemFuture.CallOpts)
}

// IsActive is a free data retrieval call binding the contract method 0x22f3e2d4.
//
// Solidity: function isActive() constant returns(bool)
func (_LockRedeemFuture *LockRedeemFutureCallerSession) IsActive() (bool, error) {
	return _LockRedeemFuture.Contract.IsActive(&_LockRedeemFuture.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCaller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemFuture.contract.Call(opts, out, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureSession) NumValidators() (*big.Int, error) {
	return _LockRedeemFuture.Contract.NumValidators(&_LockRedeemFuture.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCallerSession) NumValidators() (*big.Int, error) {
	return _LockRedeemFuture.Contract.NumValidators(&_LockRedeemFuture.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemFuture.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemFuture.Contract.Validators(&_LockRedeemFuture.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint256)
func (_LockRedeemFuture *LockRedeemFutureCallerSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemFuture.Contract.Validators(&_LockRedeemFuture.CallOpts, arg0)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemFuture *LockRedeemFutureTransactor) MigrateFromOld(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemFuture.contract.Transact(opts, "MigrateFromOld")
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemFuture *LockRedeemFutureSession) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.MigrateFromOld(&_LockRedeemFuture.TransactOpts)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemFuture *LockRedeemFutureTransactorSession) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeemFuture.Contract.MigrateFromOld(&_LockRedeemFuture.TransactOpts)
}
