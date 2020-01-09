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

// ContextABI is the input ABI used to generate the binding from.
const ContextABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// Context is an auto generated Go binding around an Ethereum contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// ERC20ABI is the input ABI used to generate the binding from.
const ERC20ABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// ERC20FuncSigs maps the 4-byte function signature to its string representation.
var ERC20FuncSigs = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"a457c2d7": "decreaseAllowance(address,uint256)",
	"39509351": "increaseAllowance(address,uint256)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
}

// ERC20Bin is the compiled bytecode used for deploying new contracts.
var ERC20Bin = "0x608060405261083b806100136000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806370a082311161005b57806370a0823114610149578063a457c2d71461016f578063a9059cbb1461019b578063dd62ed3e146101c757610088565b8063095ea7b31461008d57806318160ddd146100cd57806323b872dd146100e7578063395093511461011d575b600080fd5b6100b9600480360360408110156100a357600080fd5b506001600160a01b0381351690602001356101f5565b604080519115158252519081900360200190f35b6100d5610212565b60408051918252519081900360200190f35b6100b9600480360360608110156100fd57600080fd5b506001600160a01b03813581169160208101359091169060400135610218565b6100b96004803603604081101561013357600080fd5b506001600160a01b0381351690602001356102a5565b6100d56004803603602081101561015f57600080fd5b50356001600160a01b03166102f9565b6100b96004803603604081101561018557600080fd5b506001600160a01b038135169060200135610314565b6100b9600480360360408110156101b157600080fd5b506001600160a01b038135169060200135610382565b6100d5600480360360408110156101dd57600080fd5b506001600160a01b0381358116916020013516610396565b60006102096102026103c1565b84846103c5565b50600192915050565b60025490565b60006102258484846104b1565b61029b846102316103c1565b61029685604051806060016040528060288152602001610771602891396001600160a01b038a1660009081526001602052604081209061026f6103c1565b6001600160a01b03168152602081019190915260400160002054919063ffffffff61060d16565b6103c5565b5060019392505050565b60006102096102b26103c1565b8461029685600160006102c36103c1565b6001600160a01b03908116825260208083019390935260409182016000908120918c16815292529020549063ffffffff6106a416565b6001600160a01b031660009081526020819052604090205490565b60006102096103216103c1565b84610296856040518060600160405280602581526020016107e2602591396001600061034b6103c1565b6001600160a01b03908116825260208083019390935260409182016000908120918d1681529252902054919063ffffffff61060d16565b600061020961038f6103c1565b84846104b1565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b3390565b6001600160a01b03831661040a5760405162461bcd60e51b81526004018080602001828103825260248152602001806107be6024913960400191505060405180910390fd5b6001600160a01b03821661044f5760405162461bcd60e51b81526004018080602001828103825260228152602001806107296022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166104f65760405162461bcd60e51b81526004018080602001828103825260258152602001806107996025913960400191505060405180910390fd5b6001600160a01b03821661053b5760405162461bcd60e51b81526004018080602001828103825260238152602001806107066023913960400191505060405180910390fd5b61057e8160405180606001604052806026815260200161074b602691396001600160a01b038616600090815260208190526040902054919063ffffffff61060d16565b6001600160a01b0380851660009081526020819052604080822093909355908416815220546105b3908263ffffffff6106a416565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b6000818484111561069c5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610661578181015183820152602001610649565b50505050905090810190601f16801561068e5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b6000828201838110156106fe576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b939250505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737345524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa265627a7a72315820b17b47bed39a7fd0ea55074399062b51ddea6256e7d8b9e224f5436711d43ce364736f6c634300050b0032"

// DeployERC20 deploys a new Ethereum contract, binding an instance of ERC20 to it.
func DeployERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// ERC20 is an auto generated Go binding around an Ethereum contract.
type ERC20 struct {
	ERC20Caller     // Read-only binding to the contract
	ERC20Transactor // Write-only binding to the contract
	ERC20Filterer   // Log filterer for contract events
}

// ERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20Session struct {
	Contract     *ERC20            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20CallerSession struct {
	Contract *ERC20Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20TransactorSession struct {
	Contract     *ERC20Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20Raw struct {
	Contract *ERC20 // Generic contract binding to access the raw methods on
}

// ERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20CallerRaw struct {
	Contract *ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// ERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20TransactorRaw struct {
	Contract *ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20 creates a new instance of ERC20, bound to a specific deployed contract.
func NewERC20(address common.Address, backend bind.ContractBackend) (*ERC20, error) {
	contract, err := bindERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20{ERC20Caller: ERC20Caller{contract: contract}, ERC20Transactor: ERC20Transactor{contract: contract}, ERC20Filterer: ERC20Filterer{contract: contract}}, nil
}

// NewERC20Caller creates a new read-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Caller(address common.Address, caller bind.ContractCaller) (*ERC20Caller, error) {
	contract, err := bindERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Caller{contract: contract}, nil
}

// NewERC20Transactor creates a new write-only instance of ERC20, bound to a specific deployed contract.
func NewERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*ERC20Transactor, error) {
	contract, err := bindERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20Transactor{contract: contract}, nil
}

// NewERC20Filterer creates a new log filterer instance of ERC20, bound to a specific deployed contract.
func NewERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*ERC20Filterer, error) {
	contract, err := bindERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20Filterer{contract: contract}, nil
}

// bindERC20 binds a generic wrapper to an already deployed contract.
func bindERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20 *ERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20 *ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20 *ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_ERC20 *ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_ERC20 *ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20.Contract.Allowance(&_ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_ERC20 *ERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "balanceOf", account)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_ERC20 *ERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) constant returns(uint256)
func (_ERC20 *ERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20.Contract.BalanceOf(&_ERC20.CallOpts, account)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20Session) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_ERC20 *ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20.Contract.TotalSupply(&_ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Approve(&_ERC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ERC20 *ERC20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.DecreaseAllowance(&_ERC20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ERC20 *ERC20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.IncreaseAllowance(&_ERC20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.Transfer(&_ERC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_ERC20 *ERC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ERC20.Contract.TransferFrom(&_ERC20.TransactOpts, sender, recipient, amount)
}

// ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20 contract.
type ERC20ApprovalIterator struct {
	Event *ERC20Approval // Event containing the contract specifics and raw log

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
func (it *ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Approval)
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
		it.Event = new(ERC20Approval)
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
func (it *ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Approval represents a Approval event raised by the ERC20 contract.
type ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20ApprovalIterator{contract: _ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20 *ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Approval)
				if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ERC20 *ERC20Filterer) ParseApproval(log types.Log) (*ERC20Approval, error) {
	event := new(ERC20Approval)
	if err := _ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20 contract.
type ERC20TransferIterator struct {
	Event *ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20Transfer)
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
		it.Event = new(ERC20Transfer)
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
func (it *ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20Transfer represents a Transfer event raised by the ERC20 contract.
type ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20TransferIterator{contract: _ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20 *ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20Transfer)
				if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ERC20 *ERC20Filterer) ParseTransfer(log types.Log) (*ERC20Transfer, error) {
	event := new(ERC20Transfer)
	if err := _ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC20ABI is the input ABI used to generate the binding from.
const IERC20ABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

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
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, result, method, params...)
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
	return _IERC20.Contract.contract.Call(opts, result, method, params...)
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
	err := _IERC20.contract.Call(opts, out, "allowance", owner, spender)
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
	err := _IERC20.contract.Call(opts, out, "balanceOf", account)
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
	err := _IERC20.contract.Call(opts, out, "totalSupply")
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
const LockRedeemERCABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC20Interface\",\"outputs\":[{\"internalType\":\"contractERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTErcAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"executeredeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"getTotalErcBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_trafered\",\"type\":\"uint256\"}],\"name\":\"RedeemSuccessful\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"}]"

// LockRedeemERCFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemERCFuncSigs = map[string]string{
	"1bbc4b83": "ERC20Interface()",
	"bfb9e9f5": "addValidatorProposals(address)",
	"0d8f6b5b": "epochBlockHeight()",
	"311101b6": "executeredeem(address)",
	"231f97f1": "getOLTErcAddress()",
	"b0825584": "getTotalErcBalance(address)",
	"facd743b": "isValidator(address)",
	"0e7d275d": "newThresholdProposals(uint256)",
	"5d593f8d": "numValidators()",
	"383ea59a": "proposeAddValidator(address)",
	"e0e887d0": "proposeNewThreshold(uint256)",
	"101a8538": "proposeRemoveValidator(address)",
	"db006a75": "redeem(uint256)",
	"0d00753a": "removeValidatorProposals(address)",
	"7cacde3f": "sign(uint256,address)",
	"fa52c7d8": "validators(address)",
	"62827733": "votingThreshold()",
}

// LockRedeemERCBin is the compiled bytecode used for deploying new contracts.
var LockRedeemERCBin = "0x608060405261708060025534801561001657600080fd5b50604051610f8f380380610f8f8339818101604052602081101561003957600080fd5b810190808051604051939291908464010000000082111561005957600080fd5b90830190602082018581111561006e57600080fd5b825186602082028301116401000000008211171561008b57600080fd5b82525081516020918201928201910280838360005b838110156100b85781810151838201526020016100a0565b50505050905001604052505050600181511015610120576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602d815260200180610f33602d913960400191505060405180910390fd5b60005b81518110156101dd57600082828151811061013a57fe5b6020026020010151905060096000826001600160a01b03166001600160a01b03168152602001908152602001600020546000146101c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f815260200180610f60602f913960400191505060405180910390fd5b6101d4816001600160e01b036101e816565b50600101610123565b505060018055610245565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a250565b610cdf806102546000396000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c80635d593f8d116100a2578063bfb9e9f511610071578063bfb9e9f51461026f578063db006a7514610295578063e0e887d0146102b2578063fa52c7d8146102cf578063facd743b146102f55761010b565b80635d593f8d1461020d57806362827733146102155780637cacde3f1461021d578063b0825584146102495761010b565b80631bbc4b83116100de5780631bbc4b8314610195578063231f97f1146101b9578063311101b6146101c1578063383ea59a146101e75761010b565b80630d00753a146101105780630d8f6b5b146101485780630e7d275d14610150578063101a85381461016d575b600080fd5b6101366004803603602081101561012657600080fd5b50356001600160a01b031661032f565b60408051918252519081900360200190f35b610136610341565b6101366004803603602081101561016657600080fd5b5035610347565b6101936004803603602081101561018357600080fd5b50356001600160a01b0316610359565b005b61019d6103df565b604080516001600160a01b039092168252519081900360200190f35b61019d6103ee565b610193600480360360208110156101d757600080fd5b50356001600160a01b03166103f2565b610193600480360360208110156101fd57600080fd5b50356001600160a01b0316610549565b6101366105f6565b6101366105fc565b6101936004803603604081101561023357600080fd5b50803590602001356001600160a01b0316610602565b6101366004803603602081101561025f57600080fd5b50356001600160a01b03166107f0565b6101366004803603602081101561028557600080fd5b50356001600160a01b0316610886565b610193600480360360208110156102ab57600080fd5b5035610898565b610193600480360360208110156102c857600080fd5b5035610a53565b610136600480360360208110156102e557600080fd5b50356001600160a01b0316610b30565b61031b6004803603602081101561030b57600080fd5b50356001600160a01b0316610b42565b604080519115158252519081900360200190f35b60056020526000908152604090205481565b60035481565b60066020526000908152604090205481565b336000908152600960205260408120541361037357600080fd5b6001600160a01b0381166000908152600560209081526040808320338452600181019092529091205460ff16156103db5760405162461bcd60e51b8152600401808060200182810382526030815260200180610c276030913960400191505060405180910390fd5b5050565b600b546001600160a01b031681565b3090565b336000818152600a60205260409020546001600160a01b03161461041557600080fd5b600154336000908152600a6020526040902060030154101561043657600080fd5b600b80546001600160a01b0319166001600160a01b038381169190911791829055336000908152600a602090815260408083208054600290910154825163a9059cbb60e01b8152918616600483015260248201529051949093169363a9059cbb93604480820194918390030190829087803b1580156104b457600080fd5b505af11580156104c8573d6000803e3d6000fd5b505050506040513d60208110156104de57600080fd5b5050336000908152600a602090815260408083206002810184905560048101805460ff1916600117905554815193845290516001600160a01b0391909116927f80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad192908290030190a250565b336000908152600960205260408120541361056357600080fd5b6001600160a01b0381166000908152600460209081526040808320338452600181019092529091205460ff16156105cb5760405162461bcd60e51b815260040180806020018281038252602c815260200180610c7f602c913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff19168217905581540181556103db82610b5e565b60005481565b60015481565b61060b33610b42565b61065c576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e206c697374000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a602052604090206004015460ff16156106cd576040805162461bcd60e51b815260206004820152601b60248201527f72656465656d207265717565737420697320636f6d706c657465640000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a6020526040902060020154821461073c576040805162461bcd60e51b815260206004820152601960248201527f72656465656d20616d6f756e7420436f6d70726f6d6973656400000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a6020908152604080832033845260010190915290205460ff161561077057600080fd5b6001600160a01b0381166000818152600a602081815260408084203380865260018281018552838720805460ff1916821790559587905293835260030180549094019093558251918252810185905281517f3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db929181900390910190a25050565b600b80546001600160a01b0319166001600160a01b038381169190911791829055604080516370a0823160e01b81523060048201529051600093909216916370a0823191602480820192602092909190829003018186803b15801561085457600080fd5b505afa158015610868573d6000803e3d6000fd5b505050506040513d602081101561087e57600080fd5b505192915050565b60046020526000908152604090205481565b336000908152600a6020526040902060020154156108b557600080fd5b6000811161090a576040805162461bcd60e51b815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b336000908152600a602052604090206004015460ff1615156001146109605760405162461bcd60e51b8152600401808060200182810382526028815260200180610c576028913960400191505060405180910390fd5b336000908152600a602052604090206005015443116109c6576040805162461bcd60e51b815260206004820181905260248201527f72657175657374206973206c6f636b65642c206e6f7420617661696c61626c65604482015290519081900360640190fd5b336000818152600a6020908152604080832060048101805460ff19169055600381019390935582546001600160a01b0319169093178083556002808401869055544301600590930192909255825184815292516001600160a01b03909216927f222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2929081900390910190a250565b3360009081526009602052604081205413610a6d57600080fd5b6000548110610aad5760405162461bcd60e51b8152600401808060200182810382526041815260200180610bbc6041913960600191505060405180910390fd5b6000818152600660209081526040808320338452600181019092529091205460ff1615610b0b5760405162461bcd60e51b815260040180806020018281038252602a815260200180610bfd602a913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401905550565b60096020526000908152604090205481565b6001600160a01b03166000908152600960205260408120541390565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a25056fe4e6577207468726573686f6c647320286d29206d757374206265206c657373207468616e20746865206e756d626572206f662076616c696461746f727320286e2973656e6465722068617320616c726561647920766f74656420666f7220746869732070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f20616464207468697320746f2070726f706f73616c6561726c6965722072656465656d20686173206e6f74206265656e2065786563757465642079657473656e6465722068617320616c726561647920766f74656420746f2061646420746869732061646472657373a265627a7a723158205b5836c669084f6d59d19f9c3c08dd945993c87d90169167c3b87c4b3c1ba20164736f6c634300050b0032696e73756666696369656e742076616c696461746f72732070617373656420746f20636f6e7374727563746f72666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

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
	return _LockRedeemERC.Contract.LockRedeemERCCaller.contract.Call(opts, result, method, params...)
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
	return _LockRedeemERC.Contract.contract.Call(opts, result, method, params...)
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

// ERC20Interface is a free data retrieval call binding the contract method 0x1bbc4b83.
//
// Solidity: function ERC20Interface() constant returns(address)
func (_LockRedeemERC *LockRedeemERCCaller) ERC20Interface(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LockRedeemERC.contract.Call(opts, out, "ERC20Interface")
	return *ret0, err
}

// ERC20Interface is a free data retrieval call binding the contract method 0x1bbc4b83.
//
// Solidity: function ERC20Interface() constant returns(address)
func (_LockRedeemERC *LockRedeemERCSession) ERC20Interface() (common.Address, error) {
	return _LockRedeemERC.Contract.ERC20Interface(&_LockRedeemERC.CallOpts)
}

// ERC20Interface is a free data retrieval call binding the contract method 0x1bbc4b83.
//
// Solidity: function ERC20Interface() constant returns(address)
func (_LockRedeemERC *LockRedeemERCCallerSession) ERC20Interface() (common.Address, error) {
	return _LockRedeemERC.Contract.ERC20Interface(&_LockRedeemERC.CallOpts)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeemERC *LockRedeemERCCaller) AddValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemERC.contract.Call(opts, out, "addValidatorProposals", arg0)
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
	err := _LockRedeemERC.contract.Call(opts, out, "epochBlockHeight")
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
	err := _LockRedeemERC.contract.Call(opts, out, "getOLTErcAddress")
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

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemERC *LockRedeemERCCaller) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemERC.contract.Call(opts, out, "isValidator", addr)
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
	err := _LockRedeemERC.contract.Call(opts, out, "newThresholdProposals", arg0)
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
	err := _LockRedeemERC.contract.Call(opts, out, "numValidators")
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
	err := _LockRedeemERC.contract.Call(opts, out, "removeValidatorProposals", arg0)
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
	err := _LockRedeemERC.contract.Call(opts, out, "validators", arg0)
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

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeemERC *LockRedeemERCCaller) VotingThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemERC.contract.Call(opts, out, "votingThreshold")
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
// Solidity: function executeredeem(address tokenAddress) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) Executeredeem(opts *bind.TransactOpts, tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "executeredeem", tokenAddress)
}

// Executeredeem is a paid mutator transaction binding the contract method 0x311101b6.
//
// Solidity: function executeredeem(address tokenAddress) returns()
func (_LockRedeemERC *LockRedeemERCSession) Executeredeem(tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Executeredeem(&_LockRedeemERC.TransactOpts, tokenAddress)
}

// Executeredeem is a paid mutator transaction binding the contract method 0x311101b6.
//
// Solidity: function executeredeem(address tokenAddress) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) Executeredeem(tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Executeredeem(&_LockRedeemERC.TransactOpts, tokenAddress)
}

// GetTotalErcBalance is a paid mutator transaction binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress) returns(uint256)
func (_LockRedeemERC *LockRedeemERCTransactor) GetTotalErcBalance(opts *bind.TransactOpts, tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "getTotalErcBalance", tokenAddress)
}

// GetTotalErcBalance is a paid mutator transaction binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress) returns(uint256)
func (_LockRedeemERC *LockRedeemERCSession) GetTotalErcBalance(tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.GetTotalErcBalance(&_LockRedeemERC.TransactOpts, tokenAddress)
}

// GetTotalErcBalance is a paid mutator transaction binding the contract method 0xb0825584.
//
// Solidity: function getTotalErcBalance(address tokenAddress) returns(uint256)
func (_LockRedeemERC *LockRedeemERCTransactorSession) GetTotalErcBalance(tokenAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.GetTotalErcBalance(&_LockRedeemERC.TransactOpts, tokenAddress)
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

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemERC *LockRedeemERCTransactor) Redeem(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.contract.Transact(opts, "redeem", amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemERC *LockRedeemERCSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Redeem(&_LockRedeemERC.TransactOpts, amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemERC *LockRedeemERCTransactorSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemERC.Contract.Redeem(&_LockRedeemERC.TransactOpts, amount_)
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
var SafeMathBin = "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a723158202760e8f5a86558fd9a27d973797801f39e9e588776705a080e0e624bb0386c9764736f6c634300050b0032"

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
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
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
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
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
