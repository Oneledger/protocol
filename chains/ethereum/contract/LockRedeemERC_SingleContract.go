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
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = ethmath.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ERC20BasicABI is the input ABI used to generate the binding from.
const ERC20BasicABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numTokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numTokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"numTokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ERC20BasicFuncSigs maps the 4-byte function signature to its string representation.
var ERC20BasicFuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"313ce567": "decimals()",
	"06fdde03": "name()",
	"95d89b41": "symbol()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20BasicBin is the compiled bytecode used for deploying new contracts.
var ERC20BasicBin = "0x608060405234801561001057600080fd5b506040516107333803806107338339818101604052602081101561003357600080fd5b50516002819055336000908152602081905260409020556106da806100596000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce567146101a557806370a08231146101c357806395d89b41146101e9578063a9059cbb146101f1578063dd62ed3e1461021d57610093565b806306fdde0314610098578063095ea7b31461011557806318160ddd1461015557806323b872dd1461016f575b600080fd5b6100a061024b565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100da5781810151838201526020016100c2565b50505050905090810190601f1680156101075780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101416004803603604081101561012b57600080fd5b506001600160a01b038135169060200135610270565b604080519115158252519081900360200190f35b61015d6102d6565b60408051918252519081900360200190f35b6101416004803603606081101561018557600080fd5b506001600160a01b038135811691602081013590911690604001356102dc565b6101ad610437565b6040805160ff9092168252519081900360200190f35b61015d600480360360208110156101d957600080fd5b50356001600160a01b031661043c565b6100a0610457565b6101416004803603604081101561020757600080fd5b506001600160a01b038135169060200135610476565b61015d6004803603604081101561023357600080fd5b506001600160a01b0381358116916020013516610540565b604051806040016040528060098152602001682a32b9ba2a37b5b2b760b91b81525081565b3360008181526001602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60025490565b6001600160a01b03831660009081526020819052604081205482111561030157600080fd5b6001600160a01b038416600090815260016020908152604080832033845290915290205482111561033157600080fd5b6001600160a01b03841660009081526020819052604090205461035a908363ffffffff61056b16565b6001600160a01b038516600090815260208181526040808320939093556001815282822033835290522054610395908363ffffffff61056b16565b6001600160a01b03808616600090815260016020908152604080832033845282528083209490945591861681529081905220546103d8908363ffffffff6105b416565b6001600160a01b038085166000818152602081815260409182902094909455805186815290519193928816927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a35060019392505050565b601281565b6001600160a01b031660009081526020819052604090205490565b6040518060400160405280600381526020016254544360e81b81525081565b3360009081526020819052604081205482111561049257600080fd5b336000908152602081905260409020546104b2908363ffffffff61056b16565b33600090815260208190526040808220929092556001600160a01b038516815220546104e4908363ffffffff6105b416565b6001600160a01b038416600081815260208181526040918290209390935580518581529051919233927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a350600192915050565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60006105ad83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525061060e565b9392505050565b6000828201838110156105ad576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b6000818484111561069d5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561066257818101518382015260200161064a565b50505050905090810190601f16801561068f5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50505090039056fea265627a7a723158204d7bae91dc8e8e4a0102ca67656e053f0c94437a4825daeedf0ca258040475b364736f6c63430005100032"

// DeployERC20Basic deploys a new Ethereum contract, binding an instance of ERC20Basic to it.
func DeployERC20Basic(auth *bind.TransactOpts, backend bind.ContractBackend, total *big.Int) (common.Address, *types.Transaction, *ERC20Basic, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20BasicBin), backend, total)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// ERC20Basic is an auto generated Go binding around an Ethereum contract.
type ERC20Basic struct {
	ERC20BasicCaller     // Read-only binding to the contract
	ERC20BasicTransactor // Write-only binding to the contract
	ERC20BasicFilterer   // Log filterer for contract events
}

// ERC20BasicCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20BasicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20BasicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20BasicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20BasicSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20BasicSession struct {
	Contract     *ERC20Basic       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20BasicCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20BasicCallerSession struct {
	Contract *ERC20BasicCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ERC20BasicTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20BasicTransactorSession struct {
	Contract     *ERC20BasicTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ERC20BasicRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20BasicRaw struct {
	Contract *ERC20Basic // Generic contract binding to access the raw methods on
}

// ERC20BasicCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20BasicCallerRaw struct {
	Contract *ERC20BasicCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20BasicTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20BasicTransactorRaw struct {
	Contract *ERC20BasicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Basic creates a new instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20Basic(address common.Address, backend bind.ContractBackend) (*ERC20Basic, error) {
	contract, err := bindERC20Basic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Basic{ERC20BasicCaller: ERC20BasicCaller{contract: contract}, ERC20BasicTransactor: ERC20BasicTransactor{contract: contract}, ERC20BasicFilterer: ERC20BasicFilterer{contract: contract}}, nil
}

// NewERC20BasicCaller creates a new read-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicCaller(address common.Address, caller bind.ContractCaller) (*ERC20BasicCaller, error) {
	contract, err := bindERC20Basic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicCaller{contract: contract}, nil
}

// NewERC20BasicTransactor creates a new write-only instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20BasicTransactor, error) {
	contract, err := bindERC20Basic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransactor{contract: contract}, nil
}

// NewERC20BasicFilterer creates a new log filterer instance of ERC20Basic, bound to a specific deployed contract.
func NewERC20BasicFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20BasicFilterer, error) {
	contract, err := bindERC20Basic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicFilterer{contract: contract}, nil
}

// bindERC20Basic binds a generic wrapper to an already deployed contract.
func bindERC20Basic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20BasicABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _ERC20Basic.Contract.ERC20BasicCaller.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.ERC20BasicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Basic *ERC20BasicCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _ERC20Basic.Contract.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Basic *ERC20BasicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Basic.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address delegate) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) Allowance(opts *bind.CallOpts, owner common.Address, delegate common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "allowance", owner, delegate)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address delegate) constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) Allowance(owner common.Address, delegate common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.Allowance(&_ERC20Basic.CallOpts, owner, delegate)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address delegate) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) Allowance(owner common.Address, delegate common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.Allowance(&_ERC20Basic.CallOpts, owner, delegate)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) BalanceOf(opts *bind.CallOpts, tokenOwner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "balanceOf", tokenOwner)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, tokenOwner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address tokenOwner) constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	return _ERC20Basic.Contract.BalanceOf(&_ERC20Basic.CallOpts, tokenOwner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20Basic *ERC20BasicCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "decimals")
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20Basic *ERC20BasicSession) Decimals() (uint8, error) {
	return _ERC20Basic.Contract.Decimals(&_ERC20Basic.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint8)
func (_ERC20Basic *ERC20BasicCallerSession) Decimals() (uint8, error) {
	return _ERC20Basic.Contract.Decimals(&_ERC20Basic.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20Basic *ERC20BasicCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20Basic *ERC20BasicSession) Name() (string, error) {
	return _ERC20Basic.Contract.Name(&_ERC20Basic.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_ERC20Basic *ERC20BasicCallerSession) Name() (string, error) {
	return _ERC20Basic.Contract.Name(&_ERC20Basic.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20Basic *ERC20BasicCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "symbol")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20Basic *ERC20BasicSession) Symbol() (string, error) {
	return _ERC20Basic.Contract.Symbol(&_ERC20Basic.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_ERC20Basic *ERC20BasicCallerSession) Symbol() (string, error) {
	return _ERC20Basic.Contract.Symbol(&_ERC20Basic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _ERC20Basic.contract.Call(opts, &results, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20Basic *ERC20BasicCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Basic.Contract.TotalSupply(&_ERC20Basic.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address delegate, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactor) Approve(opts *bind.TransactOpts, delegate common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.contract.Transact(opts, "approve", delegate, numTokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address delegate, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicSession) Approve(delegate common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Approve(&_ERC20Basic.TransactOpts, delegate, numTokens)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address delegate, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactorSession) Approve(delegate common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Approve(&_ERC20Basic.TransactOpts, delegate, numTokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactor) Transfer(opts *bind.TransactOpts, receiver common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.contract.Transact(opts, "transfer", receiver, numTokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicSession) Transfer(receiver common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, receiver, numTokens)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address receiver, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactorSession) Transfer(receiver common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.Transfer(&_ERC20Basic.TransactOpts, receiver, numTokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address owner, address buyer, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactor) TransferFrom(opts *bind.TransactOpts, owner common.Address, buyer common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.contract.Transact(opts, "transferFrom", owner, buyer, numTokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address owner, address buyer, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicSession) TransferFrom(owner common.Address, buyer common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.TransferFrom(&_ERC20Basic.TransactOpts, owner, buyer, numTokens)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address owner, address buyer, uint256 numTokens) returns(bool)
func (_ERC20Basic *ERC20BasicTransactorSession) TransferFrom(owner common.Address, buyer common.Address, numTokens *big.Int) (*types.Transaction, error) {
	return _ERC20Basic.Contract.TransferFrom(&_ERC20Basic.TransactOpts, owner, buyer, numTokens)
}

// ERC20BasicApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Basic contract.
type ERC20BasicApprovalIterator struct {
	Event *ERC20BasicApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20BasicApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BasicApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20BasicApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20BasicApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BasicApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BasicApproval represents a Approval event raised by the ERC20Basic contract.
type ERC20BasicApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) FilterApproval(opts *bind.FilterOpts, tokenOwner []common.Address, spender []common.Address) (*ERC20BasicApprovalIterator, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Basic.contract.FilterLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicApprovalIterator{contract: _ERC20Basic.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20BasicApproval, tokenOwner []common.Address, spender []common.Address) (event.Subscription, error) {

	var tokenOwnerRule []interface{}
	for _, tokenOwnerItem := range tokenOwner {
		tokenOwnerRule = append(tokenOwnerRule, tokenOwnerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Basic.contract.WatchLogs(opts, "Approval", tokenOwnerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BasicApproval)
				if err := _ERC20Basic.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed tokenOwner, address indexed spender, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) ParseApproval(log types.Log) (*ERC20BasicApproval, error) {
	event := new(ERC20BasicApproval)
	if err := _ERC20Basic.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20BasicTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Basic contract.
type ERC20BasicTransferIterator struct {
	Event *ERC20BasicTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC20BasicTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20BasicTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC20BasicTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC20BasicTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20BasicTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20BasicTransfer represents a Transfer event raised by the ERC20Basic contract.
type ERC20BasicTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20BasicTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20BasicTransferIterator{contract: _ERC20Basic.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20BasicTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Basic.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20BasicTransfer)
				if err := _ERC20Basic.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 tokens)
func (_ERC20Basic *ERC20BasicFilterer) ParseTransfer(log types.Log) (*ERC20BasicTransfer, error) {
	event := new(ERC20BasicTransfer)
	if err := _ERC20Basic.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20ABI is the input ABI used to generate the binding from.
const IERC20ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IERC20FuncSigs maps the 4-byte function signature to its string representation.
var IERC20FuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// IERC20 is an auto generated Go binding around an Ethereum contract.
type IERC20 struct {
	IERC20Caller     // Read-only binding to the contract
	IERC20Transactor // Write-only binding to the contract
	IERC20Filterer   // Log filterer for contract events
}

// IERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type IERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type IERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IERC20Session struct {
	Contract     *IERC20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IERC20CallerSession struct {
	Contract *IERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IERC20TransactorSession struct {
	Contract     *IERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type IERC20Raw struct {
	Contract *IERC20 // Generic contract binding to access the raw methods on
}

// IERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IERC20CallerRaw struct {
	Contract *IERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IERC20TransactorRaw struct {
	Contract *IERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20 creates a new instance of IERC20, bound to a specific deployed contract.
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindIERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
}

// NewIERC20Caller creates a new read-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Caller(address common.Address, caller bind.ContractCaller) (*IERC20Caller, error) {
	contract, err := bindIERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Caller{contract: contract}, nil
}

// NewIERC20Transactor creates a new write-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC20Transactor, error) {
	contract, err := bindIERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Transactor{contract: contract}, nil
}

// NewIERC20Filterer creates a new log filterer instance of IERC20, bound to a specific deployed contract.
func NewIERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC20Filterer, error) {
	contract, err := bindIERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20Filterer{contract: contract}, nil
}

// bindIERC20 binds a generic wrapper to an already deployed contract.
func bindIERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _IERC20.Contract.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _IERC20.contract.Call(opts, &results, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _IERC20.contract.Call(opts, &results, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _IERC20.contract.Call(opts, &results, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, sender, recipient, amount)
}

// IERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20 contract.
type IERC20ApprovalIterator struct {
	Event *IERC20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IERC20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Approval represents a Approval event raised by the IERC20 contract.
type IERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IERC20ApprovalIterator{contract: _IERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Approval)
				if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) ParseApproval(log types.Log) (*IERC20Approval, error) {
	event := new(IERC20Approval)
	if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20 contract.
type IERC20TransferIterator struct {
	Event *IERC20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *IERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(IERC20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *IERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Transfer represents a Transfer event raised by the IERC20 contract.
type IERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IERC20TransferIterator{contract: _IERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Transfer)
				if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
	event := new(IERC20Transfer)
	if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCABI is the input ABI used to generate the binding from.
const LockRedeemERCABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_trafered\",\"type\":\"uint256\"}],\"name\":\"RedeemSuccessful\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress_\",\"type\":\"address\"}],\"name\":\"executeredeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTErcAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress_\",\"type\":\"address\"}],\"name\":\"getTotalErcBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"hasValidatorSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"tokenAddress_\",\"type\":\"address\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"verifyRedeem\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LockRedeemERCFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemERCFuncSigs = map[string]string{
	"bfb9e9f5": "addValidatorProposals(address)",
	"0d8f6b5b": "epochBlockHeight()",
	"311101b6": "executeredeem(address)",
	"231f97f1": "getOLTErcAddress()",
	"b0825584": "getTotalErcBalance(address)",
	"31b6a6d1": "hasValidatorSigned(address)",
	"facd743b": "isValidator(address)",
	"0e7d275d": "newThresholdProposals(uint256)",
	"5d593f8d": "numValidators()",
	"383ea59a": "proposeAddValidator(address)",
	"e0e887d0": "proposeNewThreshold(uint256)",
	"101a8538": "proposeRemoveValidator(address)",
	"7bde82f2": "redeem(uint256,address)",
	"0d00753a": "removeValidatorProposals(address)",
	"7cacde3f": "sign(uint256,address)",
	"fa52c7d8": "validators(address)",
	"91e39868": "verifyRedeem(address)",
	"62827733": "votingThreshold()",
}

// LockRedeemERCBin is the compiled bytecode used for deploying new contracts.
var LockRedeemERCBin = "0x608060405261708060025534801561001657600080fd5b50604051620010fc380380620010fc8339818101604052602081101561003b57600080fd5b810190808051604051939291908464010000000082111561005b57600080fd5b90830190602082018581111561007057600080fd5b825186602082028301116401000000008211171561008d57600080fd5b82525081516020918201928201910280838360005b838110156100ba5781810151838201526020016100a2565b505050509050016040525050506001815110156101095760405162461bcd60e51b815260040180806020018281038252602d815260200180620010a0602d913960400191505060405180910390fd5b60005b81518110156101ad57600082828151811061012357fe5b6020026020010151905060096000826001600160a01b03166001600160a01b03168152602001908152602001600020546000146101925760405162461bcd60e51b815260040180806020018281038252602f815260200180620010cd602f913960400191505060405180910390fd5b6101a4816001600160e01b036101cc16565b5060010161010c565b5060038151600202816101bc57fe5b0460010160018190555050610229565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a250565b610e6780620002396000396000f3fe608060405234801561001057600080fd5b50600436106101165760003560e01c806362827733116100a2578063b082558411610071578063b0825584146102d8578063bfb9e9f5146102fe578063e0e887d014610324578063fa52c7d814610341578063facd743b1461036757610116565b806362827733146102525780637bde82f21461025a5780637cacde3f1461028657806391e39868146102b257610116565b8063231f97f1116100e9578063231f97f1146101a0578063311101b6146101c457806331b6a6d1146101ea578063383ea59a146102245780635d593f8d1461024a57610116565b80630d00753a1461011b5780630d8f6b5b146101535780630e7d275d1461015b578063101a853814610178575b600080fd5b6101416004803603602081101561013157600080fd5b50356001600160a01b031661038d565b60408051918252519081900360200190f35b61014161039f565b6101416004803603602081101561017157600080fd5b50356103a5565b61019e6004803603602081101561018e57600080fd5b50356001600160a01b03166103b7565b005b6101a861043d565b604080516001600160a01b039092168252519081900360200190f35b61019e600480360360208110156101da57600080fd5b50356001600160a01b0316610441565b6102106004803603602081101561020057600080fd5b50356001600160a01b03166106de565b604080519115158252519081900360200190f35b61019e6004803603602081101561023a57600080fd5b50356001600160a01b031661070a565b6101416107b7565b6101416107bd565b61019e6004803603604081101561027057600080fd5b50803590602001356001600160a01b03166107c3565b61019e6004803603604081101561029c57600080fd5b50803590602001356001600160a01b0316610939565b610210600480360360208110156102c857600080fd5b50356001600160a01b0316610b2a565b610141600480360360208110156102ee57600080fd5b50356001600160a01b0316610b50565b6101416004803603602081101561031457600080fd5b50356001600160a01b0316610bce565b61019e6004803603602081101561033a57600080fd5b5035610be0565b6101416004803603602081101561035757600080fd5b50356001600160a01b0316610cbd565b6102106004803603602081101561037d57600080fd5b50356001600160a01b0316610ccf565b60056020526000908152604090205481565b60035481565b60066020526000908152604090205481565b33600090815260096020526040812054136103d157600080fd5b6001600160a01b0381166000908152600560209081526040808320338452600181019092529091205460ff16156104395760405162461bcd60e51b8152600401808060200182810382526030815260200180610dd76030913960400191505060405180910390fd5b5050565b3090565b336000818152600a60205260409020600101546001600160a01b0316146104995760405162461bcd60e51b8152600401808060200182810382526023815260200180610d8a6023913960400191505060405180910390fd5b600154336000908152600a60205260409020600401541015610502576040805162461bcd60e51b815260206004820152601a60248201527f4e6f7420656e6f7567682056616c696461746f7220766f746573000000000000604482015290519081900360640190fd5b336000908152600a60205260409020546001600160a01b03828116911614610568576040805162461bcd60e51b815260206004820152601460248201527315dc9bdb99c81d1bdad95b881cd95b1958dd195960621b604482015290519081900360640190fd5b336000908152600a602052604090206005015460ff161515600114156105d5576040805162461bcd60e51b815260206004820181905260248201527f52656465656d2068617320616c7265616479206265656e206578656375746564604482015290519081900360640190fd5b336000908152600a602090815260408083206001810154600390910154825163a9059cbb60e01b81526001600160a01b0392831660048201526024810191909152915185949185169363a9059cbb93604480820194929392918390030190829087803b15801561064457600080fd5b505af1158015610658573d6000803e3d6000fd5b505050506040513d602081101561066e57600080fd5b5050336000908152600a602090815260408083206003810184905560058101805460ff191660019081179091550154815193845290516001600160a01b0391909116927f80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad192908290030190a25050565b6001600160a01b03166000908152600a6020908152604080832033845260020190915290205460ff1690565b336000908152600960205260408120541361072457600080fd5b6001600160a01b0381166000908152600460209081526040808320338452600181019092529091205460ff161561078c5760405162461bcd60e51b815260040180806020018281038252602c815260200180610e07602c913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401815561043982610ceb565b60005481565b60015481565b336000908152600a6020526040902060030154156107e057600080fd5b60008211610835576040805162461bcd60e51b815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b336000908152600a6020526040902060060154431161089b576040805162461bcd60e51b815260206004820181905260248201527f72657175657374206973206c6f636b65642c206e6f7420617661696c61626c65604482015290519081900360640190fd5b336000818152600a6020908152604080832060058101805460ff1916905580546001600160a01b038781166001600160a01b03199283161783556004830195909555600182018054909116909517948590556003810187905560025443016006909101558051868152905193909216927f222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d29281900390910190a25050565b61094233610ccf565b610993576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e206c697374000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a602052604090206005015460ff1615610a04576040805162461bcd60e51b815260206004820152601b60248201527f72656465656d207265717565737420697320636f6d706c657465640000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a60205260409020600301548214610a73576040805162461bcd60e51b815260206004820152601960248201527f72656465656d20616d6f756e7420436f6d70726f6d6973656400000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a6020908152604080832033845260020190915290205460ff1615610aa757600080fd5b6001600160a01b0381166000818152600a6020818152604080842033808652600282018452828620805460ff191660019081179091559587905293835260040180549094019093558251918252810185905281517f3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db929181900390910190a25050565b6001546001600160a01b0382166000908152600a60205260409020600401541015919050565b604080516370a0823160e01b8152306004820152905160009183916001600160a01b038316916370a08231916024808301926020929190829003018186803b158015610b9b57600080fd5b505afa158015610baf573d6000803e3d6000fd5b505050506040513d6020811015610bc557600080fd5b50519392505050565b60046020526000908152604090205481565b3360009081526009602052604081205413610bfa57600080fd5b6000548110610c3a5760405162461bcd60e51b8152600401808060200182810382526041815260200180610d496041913960600191505060405180910390fd5b6000818152600660209081526040808320338452600181019092529091205460ff1615610c985760405162461bcd60e51b815260040180806020018281038252602a815260200180610dad602a913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401905550565b60096020526000908152604090205481565b6001600160a01b03166000908152600960205260408120541390565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a25056fe4e6577207468726573686f6c647320286d29206d757374206265206c657373207468616e20746865206e756d626572206f662076616c696461746f727320286e2952657175657374206e6f742063616c6c65642066726f6d20746f6b656e206f776e657273656e6465722068617320616c726561647920766f74656420666f7220746869732070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f20616464207468697320746f2070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f2061646420746869732061646472657373a265627a7a72315820bb50fc07b3175911cd673ee6f68d695014f569e601355cd8b46f33b946b4bb7064736f6c63430005100032696e73756666696369656e742076616c696461746f72732070617373656420746f20636f6e7374727563746f72666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

// DeployLockRedeemERC deploys a new Ethereum contract, binding an instance of LockRedeemERC to it.
func DeployLockRedeemERC(auth *bind.TransactOpts, backend bind.ContractBackend, initialValidators []common.Address) (common.Address, *types.Transaction, *LockRedeemERC, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemERCABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockRedeemERCBin), backend, initialValidators)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockRedeemERC{LockRedeemERCCaller: LockRedeemERCCaller{contract: contract}, LockRedeemERCTransactor: LockRedeemERCTransactor{contract: contract}, LockRedeemERCFilterer: LockRedeemERCFilterer{contract: contract}}, nil
}

// LockRedeemERC is an auto generated Go binding around an Ethereum contract.
type LockRedeemERC struct {
	LockRedeemERCCaller     // Read-only binding to the contract
	LockRedeemERCTransactor // Write-only binding to the contract
	LockRedeemERCFilterer   // Log filterer for contract events
}

// LockRedeemERCCaller is an auto generated read-only Go binding around an Ethereum contract.
type LockRedeemERCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemERCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LockRedeemERCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemERCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockRedeemERCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemERCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockRedeemERCSession struct {
	Contract     *LockRedeemERC    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockRedeemERCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockRedeemERCCallerSession struct {
	Contract *LockRedeemERCCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// LockRedeemERCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockRedeemERCTransactorSession struct {
	Contract     *LockRedeemERCTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LockRedeemERCRaw is an auto generated low-level Go binding around an Ethereum contract.
type LockRedeemERCRaw struct {
	Contract *LockRedeemERC // Generic contract binding to access the raw methods on
}

// LockRedeemERCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockRedeemERCCallerRaw struct {
	Contract *LockRedeemERCCaller // Generic read-only contract binding to access the raw methods on
}

// LockRedeemERCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockRedeemERCTransactorRaw struct {
	Contract *LockRedeemERCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLockRedeemERC creates a new instance of LockRedeemERC, bound to a specific deployed contract.
func NewLockRedeemERC(address common.Address, backend bind.ContractBackend) (*LockRedeemERC, error) {
	contract, err := bindLockRedeemERC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERC{LockRedeemERCCaller: LockRedeemERCCaller{contract: contract}, LockRedeemERCTransactor: LockRedeemERCTransactor{contract: contract}, LockRedeemERCFilterer: LockRedeemERCFilterer{contract: contract}}, nil
}

// NewLockRedeemERCCaller creates a new read-only instance of LockRedeemERC, bound to a specific deployed contract.
func NewLockRedeemERCCaller(address common.Address, caller bind.ContractCaller) (*LockRedeemERCCaller, error) {
	contract, err := bindLockRedeemERC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCCaller{contract: contract}, nil
}

// NewLockRedeemERCTransactor creates a new write-only instance of LockRedeemERC, bound to a specific deployed contract.
func NewLockRedeemERCTransactor(address common.Address, transactor bind.ContractTransactor) (*LockRedeemERCTransactor, error) {
	contract, err := bindLockRedeemERC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCTransactor{contract: contract}, nil
}

// NewLockRedeemERCFilterer creates a new log filterer instance of LockRedeemERC, bound to a specific deployed contract.
func NewLockRedeemERCFilterer(address common.Address, filterer bind.ContractFilterer) (*LockRedeemERCFilterer, error) {
	contract, err := bindLockRedeemERC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCFilterer{contract: contract}, nil
}

// bindLockRedeemERC binds a generic wrapper to an already deployed contract.
func bindLockRedeemERC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemERCABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemERC *LockRedeemERCRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _LockRedeemERC.Contract.LockRedeemERCCaller.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemERC *LockRedeemERCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.LockRedeemERCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemERC *LockRedeemERCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.LockRedeemERCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemERC *LockRedeemERCCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _LockRedeemERC.Contract.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemERC *LockRedeemERCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemERC *LockRedeemERCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.contract.Transact(opts, method, params...)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCaller) AddValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "addValidatorProposals", arg0)
	return *ret0, err
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.AddValidatorProposals(&_LockRedeemERC.CallOpts, arg0)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCallerSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.AddValidatorProposals(&_LockRedeemERC.CallOpts, arg0)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCaller) EpochBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "epochBlockHeight")
	return *ret0, err
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCSession) EpochBlockHeight() (*big.Int, error) {
	return _LockRedeemERC.Contract.EpochBlockHeight(&_LockRedeemERC.CallOpts)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCallerSession) EpochBlockHeight() (*big.Int, error) {
	return _LockRedeemERC.Contract.EpochBlockHeight(&_LockRedeemERC.CallOpts)
}

// GetOLTErcAddress is a free data retrieval call binding the contract method 0x231f97f1.
//
// Solidity: function getOLTErcAddress() constant returns(address)
func (_LockRedeemERC *LockRedeemERCCaller) GetOLTErcAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "getOLTErcAddress")
	return *ret0, err
}

// GetOLTErcAddress is a free data retrieval call binding the contract method 0x231f97f1.
//
// Solidity: function getOLTErcAddress() constant returns(address)
func (_LockRedeemERC *LockRedeemERCSession) GetOLTErcAddress() (common.Address, error) {
	return _LockRedeemERC.Contract.GetOLTErcAddress(&_LockRedeemERC.CallOpts)
}

// GetOLTErcAddress is a free data retrieval call binding the contract method 0x231f97f1.
//
// Solidity: function getOLTErcAddress() constant returns(address)
func (_LockRedeemERC *LockRedeemERCCallerSession) GetOLTErcAddress() (common.Address, error) {
	return _LockRedeemERC.Contract.GetOLTErcAddress(&_LockRedeemERC.CallOpts)
}

// GetTotalErcBalance is a free data retrieval call binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress_) constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCaller) GetTotalErcBalance(opts *bind.CallOpts, tokenAddress_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "getTotalErcBalance", tokenAddress_)
	return *ret0, err
}

// GetTotalErcBalance is a free data retrieval call binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress_) constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCSession) GetTotalErcBalance(tokenAddress_ common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.GetTotalErcBalance(&_LockRedeemERC.CallOpts, tokenAddress_)
}

// GetTotalErcBalance is a free data retrieval call binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress_) constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCallerSession) GetTotalErcBalance(tokenAddress_ common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.GetTotalErcBalance(&_LockRedeemERC.CallOpts, tokenAddress_)
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCaller) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "hasValidatorSigned", recipient_)
	return *ret0, err
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCSession) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeemERC.Contract.HasValidatorSigned(&_LockRedeemERC.CallOpts, recipient_)
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCallerSession) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeemERC.Contract.HasValidatorSigned(&_LockRedeemERC.CallOpts, recipient_)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCaller) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "isValidator", addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCSession) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeemERC.Contract.IsValidator(&_LockRedeemERC.CallOpts, addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCallerSession) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeemERC.Contract.IsValidator(&_LockRedeemERC.CallOpts, addr)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCaller) NewThresholdProposals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "newThresholdProposals", arg0)
	return *ret0, err
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _LockRedeemERC.Contract.NewThresholdProposals(&_LockRedeemERC.CallOpts, arg0)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCallerSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _LockRedeemERC.Contract.NewThresholdProposals(&_LockRedeemERC.CallOpts, arg0)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCaller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCSession) NumValidators() (*big.Int, error) {
	return _LockRedeemERC.Contract.NumValidators(&_LockRedeemERC.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCallerSession) NumValidators() (*big.Int, error) {
	return _LockRedeemERC.Contract.NumValidators(&_LockRedeemERC.CallOpts)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeemERC *LockRedeemERCCaller) ProposeRemoveValidator(opts *bind.CallOpts, v common.Address) error {
	var ()
	out := &[]interface{}{}
	err := _LockRedeemERC.contract.Call(opts, out, "proposeRemoveValidator", v)
	return err
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeemERC *LockRedeemERCSession) ProposeRemoveValidator(v common.Address) error {
	return _LockRedeemERC.Contract.ProposeRemoveValidator(&_LockRedeemERC.CallOpts, v)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeemERC *LockRedeemERCCallerSession) ProposeRemoveValidator(v common.Address) error {
	return _LockRedeemERC.Contract.ProposeRemoveValidator(&_LockRedeemERC.CallOpts, v)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCaller) RemoveValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "removeValidatorProposals", arg0)
	return *ret0, err
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.RemoveValidatorProposals(&_LockRedeemERC.CallOpts, arg0)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCallerSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.RemoveValidatorProposals(&_LockRedeemERC.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeemERC *LockRedeemERCCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeemERC *LockRedeemERCSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.Validators(&_LockRedeemERC.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeemERC *LockRedeemERCCallerSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeemERC.Contract.Validators(&_LockRedeemERC.CallOpts, arg0)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCaller) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "verifyRedeem", recipient_)
	return *ret0, err
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCSession) VerifyRedeem(recipient_ common.Address) (bool, error) {
	return _LockRedeemERC.Contract.VerifyRedeem(&_LockRedeemERC.CallOpts, recipient_)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCallerSession) VerifyRedeem(recipient_ common.Address) (bool, error) {
	return _LockRedeemERC.Contract.VerifyRedeem(&_LockRedeemERC.CallOpts, recipient_)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCaller) VotingThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeemERC.contract.Call(opts, &results, "votingThreshold")
	return *ret0, err
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCSession) VotingThreshold() (*big.Int, error) {
	return _LockRedeemERC.Contract.VotingThreshold(&_LockRedeemERC.CallOpts)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCallerSession) VotingThreshold() (*big.Int, error) {
	return _LockRedeemERC.Contract.VotingThreshold(&_LockRedeemERC.CallOpts)
}

// Executeredeem is a paid mutator transaction binding the contract method 0x311101b6.
//
// Solidity: function executeredeem(address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) Executeredeem(opts *bind.TransactOpts, tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "executeredeem", tokenAddress_)
}

// Executeredeem is a paid mutator transaction binding the contract method 0x311101b6.
//
// Solidity: function executeredeem(address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCSession) Executeredeem(tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Executeredeem(&_LockRedeemERC.TransactOpts, tokenAddress_)
}

// Executeredeem is a paid mutator transaction binding the contract method 0x311101b6.
//
// Solidity: function executeredeem(address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) Executeredeem(tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Executeredeem(&_LockRedeemERC.TransactOpts, tokenAddress_)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) ProposeAddValidator(opts *bind.TransactOpts, v common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "proposeAddValidator", v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeemERC *LockRedeemERCSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.ProposeAddValidator(&_LockRedeemERC.TransactOpts, v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.ProposeAddValidator(&_LockRedeemERC.TransactOpts, v)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) ProposeNewThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "proposeNewThreshold", threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeemERC *LockRedeemERCSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.ProposeNewThreshold(&_LockRedeemERC.TransactOpts, threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.ProposeNewThreshold(&_LockRedeemERC.TransactOpts, threshold)
}

// Redeem is a paid mutator transaction binding the contract method 0x7bde82f2.
//
// Solidity: function redeem(uint256 amount_, address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) Redeem(opts *bind.TransactOpts, amount_ *big.Int, tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "redeem", amount_, tokenAddress_)
}

// Redeem is a paid mutator transaction binding the contract method 0x7bde82f2.
//
// Solidity: function redeem(uint256 amount_, address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCSession) Redeem(amount_ *big.Int, tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Redeem(&_LockRedeemERC.TransactOpts, amount_, tokenAddress_)
}

// Redeem is a paid mutator transaction binding the contract method 0x7bde82f2.
//
// Solidity: function redeem(uint256 amount_, address tokenAddress_) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) Redeem(amount_ *big.Int, tokenAddress_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Redeem(&_LockRedeemERC.TransactOpts, amount_, tokenAddress_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "sign", amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemERC *LockRedeemERCSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Sign(&_LockRedeemERC.TransactOpts, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Sign(&_LockRedeemERC.TransactOpts, amount_, recipient_)
}

// LockRedeemERCAddValidatorIterator is returned from FilterAddValidator and is used to iterate over the raw logs and unpacked data for AddValidator events raised by the LockRedeemERC contract.
type LockRedeemERCAddValidatorIterator struct {
	Event *LockRedeemERCAddValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCAddValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCAddValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCAddValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCAddValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCAddValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCAddValidator represents a AddValidator event raised by the LockRedeemERC contract.
type LockRedeemERCAddValidator struct {
	Address common.Address
	Power   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterAddValidator(opts *bind.FilterOpts, _address []common.Address) (*LockRedeemERCAddValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCAddValidatorIterator{contract: _LockRedeemERC.contract, event: "AddValidator", logs: logs, sub: sub}, nil
}

// WatchAddValidator is a free log subscription operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchAddValidator(opts *bind.WatchOpts, sink chan<- *LockRedeemERCAddValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCAddValidator)
				if err := _LockRedeemERC.contract.UnpackLog(event, "AddValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddValidator is a log parse operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseAddValidator(log types.Log) (*LockRedeemERCAddValidator, error) {
	event := new(LockRedeemERCAddValidator)
	if err := _LockRedeemERC.contract.UnpackLog(event, "AddValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCDeleteValidatorIterator is returned from FilterDeleteValidator and is used to iterate over the raw logs and unpacked data for DeleteValidator events raised by the LockRedeemERC contract.
type LockRedeemERCDeleteValidatorIterator struct {
	Event *LockRedeemERCDeleteValidator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCDeleteValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCDeleteValidator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCDeleteValidator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCDeleteValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCDeleteValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCDeleteValidator represents a DeleteValidator event raised by the LockRedeemERC contract.
type LockRedeemERCDeleteValidator struct {
	Address common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDeleteValidator is a free log retrieval operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterDeleteValidator(opts *bind.FilterOpts, _address []common.Address) (*LockRedeemERCDeleteValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCDeleteValidatorIterator{contract: _LockRedeemERC.contract, event: "DeleteValidator", logs: logs, sub: sub}, nil
}

// WatchDeleteValidator is a free log subscription operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchDeleteValidator(opts *bind.WatchOpts, sink chan<- *LockRedeemERCDeleteValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCDeleteValidator)
				if err := _LockRedeemERC.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeleteValidator is a log parse operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseDeleteValidator(log types.Log) (*LockRedeemERCDeleteValidator, error) {
	event := new(LockRedeemERCDeleteValidator)
	if err := _LockRedeemERC.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the LockRedeemERC contract.
type LockRedeemERCNewEpochIterator struct {
	Event *LockRedeemERCNewEpoch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCNewEpoch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCNewEpoch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCNewEpoch represents a NewEpoch event raised by the LockRedeemERC contract.
type LockRedeemERCNewEpoch struct {
	EpochHeight *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*LockRedeemERCNewEpochIterator, error) {

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCNewEpochIterator{contract: _LockRedeemERC.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *LockRedeemERCNewEpoch) (event.Subscription, error) {

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCNewEpoch)
				if err := _LockRedeemERC.contract.UnpackLog(event, "NewEpoch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewEpoch is a log parse operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseNewEpoch(log types.Log) (*LockRedeemERCNewEpoch, error) {
	event := new(LockRedeemERCNewEpoch)
	if err := _LockRedeemERC.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCNewThresholdIterator is returned from FilterNewThreshold and is used to iterate over the raw logs and unpacked data for NewThreshold events raised by the LockRedeemERC contract.
type LockRedeemERCNewThresholdIterator struct {
	Event *LockRedeemERCNewThreshold // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCNewThresholdIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCNewThreshold)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCNewThreshold)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCNewThresholdIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCNewThresholdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCNewThreshold represents a NewThreshold event raised by the LockRedeemERC contract.
type LockRedeemERCNewThreshold struct {
	PrevThreshold *big.Int
	NewThreshold  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewThreshold is a free log retrieval operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterNewThreshold(opts *bind.FilterOpts) (*LockRedeemERCNewThresholdIterator, error) {

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCNewThresholdIterator{contract: _LockRedeemERC.contract, event: "NewThreshold", logs: logs, sub: sub}, nil
}

// WatchNewThreshold is a free log subscription operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchNewThreshold(opts *bind.WatchOpts, sink chan<- *LockRedeemERCNewThreshold) (event.Subscription, error) {

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCNewThreshold)
				if err := _LockRedeemERC.contract.UnpackLog(event, "NewThreshold", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewThreshold is a log parse operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseNewThreshold(log types.Log) (*LockRedeemERCNewThreshold, error) {
	event := new(LockRedeemERCNewThreshold)
	if err := _LockRedeemERC.contract.UnpackLog(event, "NewThreshold", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCRedeemRequestIterator is returned from FilterRedeemRequest and is used to iterate over the raw logs and unpacked data for RedeemRequest events raised by the LockRedeemERC contract.
type LockRedeemERCRedeemRequestIterator struct {
	Event *LockRedeemERCRedeemRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCRedeemRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCRedeemRequest)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCRedeemRequest)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCRedeemRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCRedeemRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCRedeemRequest represents a RedeemRequest event raised by the LockRedeemERC contract.
type LockRedeemERCRedeemRequest struct {
	Recepient       common.Address
	AmountRequested *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedeemRequest is a free log retrieval operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterRedeemRequest(opts *bind.FilterOpts, recepient []common.Address) (*LockRedeemERCRedeemRequestIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCRedeemRequestIterator{contract: _LockRedeemERC.contract, event: "RedeemRequest", logs: logs, sub: sub}, nil
}

// WatchRedeemRequest is a free log subscription operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchRedeemRequest(opts *bind.WatchOpts, sink chan<- *LockRedeemERCRedeemRequest, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCRedeemRequest)
				if err := _LockRedeemERC.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemRequest is a log parse operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseRedeemRequest(log types.Log) (*LockRedeemERCRedeemRequest, error) {
	event := new(LockRedeemERCRedeemRequest)
	if err := _LockRedeemERC.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCRedeemSuccessfulIterator is returned from FilterRedeemSuccessful and is used to iterate over the raw logs and unpacked data for RedeemSuccessful events raised by the LockRedeemERC contract.
type LockRedeemERCRedeemSuccessfulIterator struct {
	Event *LockRedeemERCRedeemSuccessful // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCRedeemSuccessfulIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCRedeemSuccessful)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCRedeemSuccessful)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCRedeemSuccessfulIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCRedeemSuccessfulIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCRedeemSuccessful represents a RedeemSuccessful event raised by the LockRedeemERC contract.
type LockRedeemERCRedeemSuccessful struct {
	Recepient      common.Address
	AmountTrafered *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRedeemSuccessful is a free log retrieval operation binding the contract event 0x80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad1.
//
// Solidity: event RedeemSuccessful(address indexed recepient, uint256 amount_trafered)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterRedeemSuccessful(opts *bind.FilterOpts, recepient []common.Address) (*LockRedeemERCRedeemSuccessfulIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "RedeemSuccessful", recepientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCRedeemSuccessfulIterator{contract: _LockRedeemERC.contract, event: "RedeemSuccessful", logs: logs, sub: sub}, nil
}

// WatchRedeemSuccessful is a free log subscription operation binding the contract event 0x80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad1.
//
// Solidity: event RedeemSuccessful(address indexed recepient, uint256 amount_trafered)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchRedeemSuccessful(opts *bind.WatchOpts, sink chan<- *LockRedeemERCRedeemSuccessful, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "RedeemSuccessful", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCRedeemSuccessful)
				if err := _LockRedeemERC.contract.UnpackLog(event, "RedeemSuccessful", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemSuccessful is a log parse operation binding the contract event 0x80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad1.
//
// Solidity: event RedeemSuccessful(address indexed recepient, uint256 amount_trafered)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseRedeemSuccessful(log types.Log) (*LockRedeemERCRedeemSuccessful, error) {
	event := new(LockRedeemERCRedeemSuccessful)
	if err := _LockRedeemERC.contract.UnpackLog(event, "RedeemSuccessful", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemERCValidatorSignedRedeemIterator is returned from FilterValidatorSignedRedeem and is used to iterate over the raw logs and unpacked data for ValidatorSignedRedeem events raised by the LockRedeemERC contract.
type LockRedeemERCValidatorSignedRedeemIterator struct {
	Event *LockRedeemERCValidatorSignedRedeem // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LockRedeemERCValidatorSignedRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemERCValidatorSignedRedeem)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LockRedeemERCValidatorSignedRedeem)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LockRedeemERCValidatorSignedRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemERCValidatorSignedRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemERCValidatorSignedRedeem represents a ValidatorSignedRedeem event raised by the LockRedeemERC contract.
type LockRedeemERCValidatorSignedRedeem struct {
	Recipient         common.Address
	ValidatorAddresss common.Address
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeemERC *LockRedeemERCFilterer) FilterValidatorSignedRedeem(opts *bind.FilterOpts, recipient []common.Address) (*LockRedeemERCValidatorSignedRedeemIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.FilterLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemERCValidatorSignedRedeemIterator{contract: _LockRedeemERC.contract, event: "ValidatorSignedRedeem", logs: logs, sub: sub}, nil
}

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeemERC *LockRedeemERCFilterer) WatchValidatorSignedRedeem(opts *bind.WatchOpts, sink chan<- *LockRedeemERCValidatorSignedRedeem, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeemERC.contract.WatchLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemERCValidatorSignedRedeem)
				if err := _LockRedeemERC.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorSignedRedeem is a log parse operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeemERC *LockRedeemERCFilterer) ParseValidatorSignedRedeem(log types.Log) (*LockRedeemERCValidatorSignedRedeem, error) {
	event := new(LockRedeemERCValidatorSignedRedeem)
	if err := _LockRedeemERC.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
var SafeMathBin = "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a72315820159fa05f1db86937744fcd83e9ceddfbcf4031ba01d8219fb13e794199fe8df564736f6c63430005100032"

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	results := make([]interface{}, 1)
	results[0] = result
	return _SafeMath.Contract.contract.Call(opts, &results, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
