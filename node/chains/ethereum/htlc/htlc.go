// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package htlc

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

// HtlcABI is the input ABI used to generate the binding from.
const HtlcABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"receiver_\",\"type\":\"address\"},{\"name\":\"balance_\",\"type\":\"uint256\"},{\"name\":\"scrHash_\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_lockPeriod\",\"type\":\"uint256\"},{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_scrHash\",\"type\":\"bytes32\"}],\"name\":\"setup\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lockPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"scrHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"scr_\",\"type\":\"bytes\"}],\"name\":\"refund\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"sender\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"scr_\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"extractMsg\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"scr_\",\"type\":\"bytes\"}],\"name\":\"validate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"funds\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"receiver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"startFromTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_sender\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Release\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Rollback\",\"type\":\"event\"}]"

// HtlcBin is the compiled bytecode used for deploying new contracts.
const HtlcBin = `608060405234801561001057600080fd5b50604051602080610dd083398101806040528101908080519060200190929190505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561006f57600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555042600481905550600060038190555050610d02806100ce6000396000f3006080604052600436106100c5576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806309ce7368146100ca5780633716c4f81461013d5780633fd8b02f146101b057806345d8b894146101db5780635eb7d9461461020e57806367e404ce1461028f5780639945e3d3146102e6578063b69ef8a814610367578063b77577cd14610392578063c16e50ef14610422578063c89f2ce41461048b578063f7260d3e14610495578063f85da5ca146104ec575b600080fd5b3480156100d657600080fd5b50610123600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035600019169060200190929190505050610517565b604051808215151515815260200191505060405180910390f35b34801561014957600080fd5b5061019660048036038101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080356000191690602001909291905050506105c0565b604051808215151515815260200191505060405180910390f35b3480156101bc57600080fd5b506101c56106f1565b6040518082815260200191505060405180910390f35b3480156101e757600080fd5b506101f06106f7565b60405180826000191660001916815260200191505060405180910390f35b34801561021a57600080fd5b50610275600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506106fd565b604051808215151515815260200191505060405180910390f35b34801561029b57600080fd5b506102a4610859565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102f257600080fd5b5061034d600480360381019080803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061087e565b604051808215151515815260200191505060405180910390f35b34801561037357600080fd5b5061037c6109df565b6040518082815260200191505060405180910390f35b34801561039e57600080fd5b506103a76109e5565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103e75780820151818401526020810190506103cc565b50505050905090810190601f1680156104145780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561042e57600080fd5b50610489600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610a87565b005b610493610b81565b005b3480156104a157600080fd5b506104aa610c05565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156104f857600080fd5b50610501610c2b565b6040518082815260200191505060405180910390f35b60008373ffffffffffffffffffffffffffffffffffffffff16600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151561057557600080fd5b8260025414151561058557600080fd5b61a8c042016004546003540111151561059d57600080fd5b8160001916600554600019161415156105b557600080fd5b600190509392505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561061d57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415151561065957600080fd5b62015180841015151561066b57600080fd5b600060025411151561067c57600080fd5b426003546004540110151561069057600080fd5b82600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550836003819055508160058160001916905550426004819055509392505050565b60035481565b60055481565b6000426003546004540111151561071357600080fd5b61071c82610a87565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc6002549081150290604051600060405180830381858888f19350505050158015610785573d6000803e3d6000fd5b5060006002819055507fbaf3b92e813efec2b7525399a930acf56a9ea74f17622f3f1080387356d1c711306000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600254604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a160019050919050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600061088982610a87565b816006908051906020019061089f929190610c31565b50600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc6002549081150290604051600060405180830381858888f1935050505015801561090a573d6000803e3d6000fd5b5060006002819055507fcb54aad3bd772fcfe1bc124e01bd1a91a91c9d80126d8b3014c4d9e687d5ca4830600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16600254604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001935050505060405180910390a160019050919050565b60025481565b606060068054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610a7d5780601f10610a5257610100808354040283529160200191610a7d565b820191906000526020600020905b815481529060010190602001808311610a6057829003601f168201915b5050505050905090565b60055460001916816040516020018082805190602001908083835b602083101515610ac75780518252602082019150602081019050602083039250610aa2565b6001836020036101000a0380198251168184511680821785525050505050509050019150506040516020818303038152906040526040518082805190602001908083835b602083101515610b305780518252602082019150602081019050602083039250610b0b565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902060001916141515610b6d57600080fd5b6000600254111515610b7e57600080fd5b50565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610bdc57600080fd5b600034111515610beb57600080fd5b6000600254141515610bfc57600080fd5b34600281905550565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60045481565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610c7257805160ff1916838001178555610ca0565b82800160010185558215610ca0579182015b82811115610c9f578251825591602001919060010190610c84565b5b509050610cad9190610cb1565b5090565b610cd391905b80821115610ccf576000816000905550600101610cb7565b5090565b905600a165627a7a72305820eb6eb7ea9128d317172ad3504a4ee10c3dba9189669d525438422c100d3320b80029`

// DeployHtlc deploys a new Ethereum contract, binding an instance of Htlc to it.
func DeployHtlc(auth *bind.TransactOpts, backend bind.ContractBackend, _sender common.Address) (common.Address, *types.Transaction, *Htlc, error) {
	parsed, err := abi.JSON(strings.NewReader(HtlcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HtlcBin), backend, _sender)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Htlc{HtlcCaller: HtlcCaller{contract: contract}, HtlcTransactor: HtlcTransactor{contract: contract}, HtlcFilterer: HtlcFilterer{contract: contract}}, nil
}

// Htlc is an auto generated Go binding around an Ethereum contract.
type Htlc struct {
	HtlcCaller     // Read-only binding to the contract
	HtlcTransactor // Write-only binding to the contract
	HtlcFilterer   // Log filterer for contract events
}

// HtlcCaller is an auto generated read-only Go binding around an Ethereum contract.
type HtlcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HtlcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HtlcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HtlcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HtlcSession struct {
	Contract     *Htlc             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HtlcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HtlcCallerSession struct {
	Contract *HtlcCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// HtlcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HtlcTransactorSession struct {
	Contract     *HtlcTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HtlcRaw is an auto generated low-level Go binding around an Ethereum contract.
type HtlcRaw struct {
	Contract *Htlc // Generic contract binding to access the raw methods on
}

// HtlcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HtlcCallerRaw struct {
	Contract *HtlcCaller // Generic read-only contract binding to access the raw methods on
}

// HtlcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HtlcTransactorRaw struct {
	Contract *HtlcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHtlc creates a new instance of Htlc, bound to a specific deployed contract.
func NewHtlc(address common.Address, backend bind.ContractBackend) (*Htlc, error) {
	contract, err := bindHtlc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Htlc{HtlcCaller: HtlcCaller{contract: contract}, HtlcTransactor: HtlcTransactor{contract: contract}, HtlcFilterer: HtlcFilterer{contract: contract}}, nil
}

// NewHtlcCaller creates a new read-only instance of Htlc, bound to a specific deployed contract.
func NewHtlcCaller(address common.Address, caller bind.ContractCaller) (*HtlcCaller, error) {
	contract, err := bindHtlc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HtlcCaller{contract: contract}, nil
}

// NewHtlcTransactor creates a new write-only instance of Htlc, bound to a specific deployed contract.
func NewHtlcTransactor(address common.Address, transactor bind.ContractTransactor) (*HtlcTransactor, error) {
	contract, err := bindHtlc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HtlcTransactor{contract: contract}, nil
}

// NewHtlcFilterer creates a new log filterer instance of Htlc, bound to a specific deployed contract.
func NewHtlcFilterer(address common.Address, filterer bind.ContractFilterer) (*HtlcFilterer, error) {
	contract, err := bindHtlc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HtlcFilterer{contract: contract}, nil
}

// bindHtlc binds a generic wrapper to an already deployed contract.
func bindHtlc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HtlcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Htlc *HtlcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Htlc.Contract.HtlcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Htlc *HtlcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Htlc.Contract.HtlcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Htlc *HtlcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Htlc.Contract.HtlcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Htlc *HtlcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Htlc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Htlc *HtlcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Htlc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Htlc *HtlcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Htlc.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0x09ce7368.
//
// Solidity: function audit(receiver_ address, balance_ uint256, scrHash_ bytes32) constant returns(bool)
func (_Htlc *HtlcCaller) Audit(opts *bind.CallOpts, receiver_ common.Address, balance_ *big.Int, scrHash_ [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "audit", receiver_, balance_, scrHash_)
	return *ret0, err
}

// Audit is a free data retrieval call binding the contract method 0x09ce7368.
//
// Solidity: function audit(receiver_ address, balance_ uint256, scrHash_ bytes32) constant returns(bool)
func (_Htlc *HtlcSession) Audit(receiver_ common.Address, balance_ *big.Int, scrHash_ [32]byte) (bool, error) {
	return _Htlc.Contract.Audit(&_Htlc.CallOpts, receiver_, balance_, scrHash_)
}

// Audit is a free data retrieval call binding the contract method 0x09ce7368.
//
// Solidity: function audit(receiver_ address, balance_ uint256, scrHash_ bytes32) constant returns(bool)
func (_Htlc *HtlcCallerSession) Audit(receiver_ common.Address, balance_ *big.Int, scrHash_ [32]byte) (bool, error) {
	return _Htlc.Contract.Audit(&_Htlc.CallOpts, receiver_, balance_, scrHash_)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_Htlc *HtlcCaller) Balance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "balance")
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_Htlc *HtlcSession) Balance() (*big.Int, error) {
	return _Htlc.Contract.Balance(&_Htlc.CallOpts)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_Htlc *HtlcCallerSession) Balance() (*big.Int, error) {
	return _Htlc.Contract.Balance(&_Htlc.CallOpts)
}

// ExtractMsg is a free data retrieval call binding the contract method 0xb77577cd.
//
// Solidity: function extractMsg() constant returns(bytes)
func (_Htlc *HtlcCaller) ExtractMsg(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "extractMsg")
	return *ret0, err
}

// ExtractMsg is a free data retrieval call binding the contract method 0xb77577cd.
//
// Solidity: function extractMsg() constant returns(bytes)
func (_Htlc *HtlcSession) ExtractMsg() ([]byte, error) {
	return _Htlc.Contract.ExtractMsg(&_Htlc.CallOpts)
}

// ExtractMsg is a free data retrieval call binding the contract method 0xb77577cd.
//
// Solidity: function extractMsg() constant returns(bytes)
func (_Htlc *HtlcCallerSession) ExtractMsg() ([]byte, error) {
	return _Htlc.Contract.ExtractMsg(&_Htlc.CallOpts)
}

// LockPeriod is a free data retrieval call binding the contract method 0x3fd8b02f.
//
// Solidity: function lockPeriod() constant returns(uint256)
func (_Htlc *HtlcCaller) LockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "lockPeriod")
	return *ret0, err
}

// LockPeriod is a free data retrieval call binding the contract method 0x3fd8b02f.
//
// Solidity: function lockPeriod() constant returns(uint256)
func (_Htlc *HtlcSession) LockPeriod() (*big.Int, error) {
	return _Htlc.Contract.LockPeriod(&_Htlc.CallOpts)
}

// LockPeriod is a free data retrieval call binding the contract method 0x3fd8b02f.
//
// Solidity: function lockPeriod() constant returns(uint256)
func (_Htlc *HtlcCallerSession) LockPeriod() (*big.Int, error) {
	return _Htlc.Contract.LockPeriod(&_Htlc.CallOpts)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() constant returns(address)
func (_Htlc *HtlcCaller) Receiver(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "receiver")
	return *ret0, err
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() constant returns(address)
func (_Htlc *HtlcSession) Receiver() (common.Address, error) {
	return _Htlc.Contract.Receiver(&_Htlc.CallOpts)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() constant returns(address)
func (_Htlc *HtlcCallerSession) Receiver() (common.Address, error) {
	return _Htlc.Contract.Receiver(&_Htlc.CallOpts)
}

// ScrHash is a free data retrieval call binding the contract method 0x45d8b894.
//
// Solidity: function scrHash() constant returns(bytes32)
func (_Htlc *HtlcCaller) ScrHash(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "scrHash")
	return *ret0, err
}

// ScrHash is a free data retrieval call binding the contract method 0x45d8b894.
//
// Solidity: function scrHash() constant returns(bytes32)
func (_Htlc *HtlcSession) ScrHash() ([32]byte, error) {
	return _Htlc.Contract.ScrHash(&_Htlc.CallOpts)
}

// ScrHash is a free data retrieval call binding the contract method 0x45d8b894.
//
// Solidity: function scrHash() constant returns(bytes32)
func (_Htlc *HtlcCallerSession) ScrHash() ([32]byte, error) {
	return _Htlc.Contract.ScrHash(&_Htlc.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() constant returns(address)
func (_Htlc *HtlcCaller) Sender(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "sender")
	return *ret0, err
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() constant returns(address)
func (_Htlc *HtlcSession) Sender() (common.Address, error) {
	return _Htlc.Contract.Sender(&_Htlc.CallOpts)
}

// Sender is a free data retrieval call binding the contract method 0x67e404ce.
//
// Solidity: function sender() constant returns(address)
func (_Htlc *HtlcCallerSession) Sender() (common.Address, error) {
	return _Htlc.Contract.Sender(&_Htlc.CallOpts)
}

// StartFromTime is a free data retrieval call binding the contract method 0xf85da5ca.
//
// Solidity: function startFromTime() constant returns(uint256)
func (_Htlc *HtlcCaller) StartFromTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Htlc.contract.Call(opts, out, "startFromTime")
	return *ret0, err
}

// StartFromTime is a free data retrieval call binding the contract method 0xf85da5ca.
//
// Solidity: function startFromTime() constant returns(uint256)
func (_Htlc *HtlcSession) StartFromTime() (*big.Int, error) {
	return _Htlc.Contract.StartFromTime(&_Htlc.CallOpts)
}

// StartFromTime is a free data retrieval call binding the contract method 0xf85da5ca.
//
// Solidity: function startFromTime() constant returns(uint256)
func (_Htlc *HtlcCallerSession) StartFromTime() (*big.Int, error) {
	return _Htlc.Contract.StartFromTime(&_Htlc.CallOpts)
}

// Validate is a free data retrieval call binding the contract method 0xc16e50ef.
//
// Solidity: function validate(scr_ bytes) constant returns()
func (_Htlc *HtlcCaller) Validate(opts *bind.CallOpts, scr_ []byte) error {
	var ()
	out := &[]interface{}{}
	err := _Htlc.contract.Call(opts, out, "validate", scr_)
	return err
}

// Validate is a free data retrieval call binding the contract method 0xc16e50ef.
//
// Solidity: function validate(scr_ bytes) constant returns()
func (_Htlc *HtlcSession) Validate(scr_ []byte) error {
	return _Htlc.Contract.Validate(&_Htlc.CallOpts, scr_)
}

// Validate is a free data retrieval call binding the contract method 0xc16e50ef.
//
// Solidity: function validate(scr_ bytes) constant returns()
func (_Htlc *HtlcCallerSession) Validate(scr_ []byte) error {
	return _Htlc.Contract.Validate(&_Htlc.CallOpts, scr_)
}

// Funds is a paid mutator transaction binding the contract method 0xc89f2ce4.
//
// Solidity: function funds() returns()
func (_Htlc *HtlcTransactor) Funds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "funds")
}

// Funds is a paid mutator transaction binding the contract method 0xc89f2ce4.
//
// Solidity: function funds() returns()
func (_Htlc *HtlcSession) Funds() (*types.Transaction, error) {
	return _Htlc.Contract.Funds(&_Htlc.TransactOpts)
}

// Funds is a paid mutator transaction binding the contract method 0xc89f2ce4.
//
// Solidity: function funds() returns()
func (_Htlc *HtlcTransactorSession) Funds() (*types.Transaction, error) {
	return _Htlc.Contract.Funds(&_Htlc.TransactOpts)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(scr_ bytes) returns(bool)
func (_Htlc *HtlcTransactor) Redeem(opts *bind.TransactOpts, scr_ []byte) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "redeem", scr_)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(scr_ bytes) returns(bool)
func (_Htlc *HtlcSession) Redeem(scr_ []byte) (*types.Transaction, error) {
	return _Htlc.Contract.Redeem(&_Htlc.TransactOpts, scr_)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(scr_ bytes) returns(bool)
func (_Htlc *HtlcTransactorSession) Redeem(scr_ []byte) (*types.Transaction, error) {
	return _Htlc.Contract.Redeem(&_Htlc.TransactOpts, scr_)
}

// Refund is a paid mutator transaction binding the contract method 0x5eb7d946.
//
// Solidity: function refund(scr_ bytes) returns(bool)
func (_Htlc *HtlcTransactor) Refund(opts *bind.TransactOpts, scr_ []byte) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "refund", scr_)
}

// Refund is a paid mutator transaction binding the contract method 0x5eb7d946.
//
// Solidity: function refund(scr_ bytes) returns(bool)
func (_Htlc *HtlcSession) Refund(scr_ []byte) (*types.Transaction, error) {
	return _Htlc.Contract.Refund(&_Htlc.TransactOpts, scr_)
}

// Refund is a paid mutator transaction binding the contract method 0x5eb7d946.
//
// Solidity: function refund(scr_ bytes) returns(bool)
func (_Htlc *HtlcTransactorSession) Refund(scr_ []byte) (*types.Transaction, error) {
	return _Htlc.Contract.Refund(&_Htlc.TransactOpts, scr_)
}

// Setup is a paid mutator transaction binding the contract method 0x3716c4f8.
//
// Solidity: function setup(_lockPeriod uint256, _receiver address, _scrHash bytes32) returns(bool)
func (_Htlc *HtlcTransactor) Setup(opts *bind.TransactOpts, _lockPeriod *big.Int, _receiver common.Address, _scrHash [32]byte) (*types.Transaction, error) {
	return _Htlc.contract.Transact(opts, "setup", _lockPeriod, _receiver, _scrHash)
}

// Setup is a paid mutator transaction binding the contract method 0x3716c4f8.
//
// Solidity: function setup(_lockPeriod uint256, _receiver address, _scrHash bytes32) returns(bool)
func (_Htlc *HtlcSession) Setup(_lockPeriod *big.Int, _receiver common.Address, _scrHash [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Setup(&_Htlc.TransactOpts, _lockPeriod, _receiver, _scrHash)
}

// Setup is a paid mutator transaction binding the contract method 0x3716c4f8.
//
// Solidity: function setup(_lockPeriod uint256, _receiver address, _scrHash bytes32) returns(bool)
func (_Htlc *HtlcTransactorSession) Setup(_lockPeriod *big.Int, _receiver common.Address, _scrHash [32]byte) (*types.Transaction, error) {
	return _Htlc.Contract.Setup(&_Htlc.TransactOpts, _lockPeriod, _receiver, _scrHash)
}

// HtlcReleaseIterator is returned from FilterRelease and is used to iterate over the raw logs and unpacked data for Release events raised by the Htlc contract.
type HtlcReleaseIterator struct {
	Event *HtlcRelease // Event containing the contract specifics and raw log

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
func (it *HtlcReleaseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HtlcRelease)
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
		it.Event = new(HtlcRelease)
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
func (it *HtlcReleaseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HtlcReleaseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HtlcRelease represents a Release event raised by the Htlc contract.
type HtlcRelease struct {
	Sender   common.Address
	Receiver common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRelease is a free log retrieval operation binding the contract event 0xcb54aad3bd772fcfe1bc124e01bd1a91a91c9d80126d8b3014c4d9e687d5ca48.
//
// Solidity: e Release(sender address, receiver address, value uint256)
func (_Htlc *HtlcFilterer) FilterRelease(opts *bind.FilterOpts) (*HtlcReleaseIterator, error) {

	logs, sub, err := _Htlc.contract.FilterLogs(opts, "Release")
	if err != nil {
		return nil, err
	}
	return &HtlcReleaseIterator{contract: _Htlc.contract, event: "Release", logs: logs, sub: sub}, nil
}

// WatchRelease is a free log subscription operation binding the contract event 0xcb54aad3bd772fcfe1bc124e01bd1a91a91c9d80126d8b3014c4d9e687d5ca48.
//
// Solidity: e Release(sender address, receiver address, value uint256)
func (_Htlc *HtlcFilterer) WatchRelease(opts *bind.WatchOpts, sink chan<- *HtlcRelease) (event.Subscription, error) {

	logs, sub, err := _Htlc.contract.WatchLogs(opts, "Release")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HtlcRelease)
				if err := _Htlc.contract.UnpackLog(event, "Release", log); err != nil {
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

// HtlcRollbackIterator is returned from FilterRollback and is used to iterate over the raw logs and unpacked data for Rollback events raised by the Htlc contract.
type HtlcRollbackIterator struct {
	Event *HtlcRollback // Event containing the contract specifics and raw log

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
func (it *HtlcRollbackIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HtlcRollback)
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
		it.Event = new(HtlcRollback)
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
func (it *HtlcRollbackIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HtlcRollbackIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HtlcRollback represents a Rollback event raised by the Htlc contract.
type HtlcRollback struct {
	Sender   common.Address
	Receiver common.Address
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRollback is a free log retrieval operation binding the contract event 0xbaf3b92e813efec2b7525399a930acf56a9ea74f17622f3f1080387356d1c711.
//
// Solidity: e Rollback(sender address, receiver address, value uint256)
func (_Htlc *HtlcFilterer) FilterRollback(opts *bind.FilterOpts) (*HtlcRollbackIterator, error) {

	logs, sub, err := _Htlc.contract.FilterLogs(opts, "Rollback")
	if err != nil {
		return nil, err
	}
	return &HtlcRollbackIterator{contract: _Htlc.contract, event: "Rollback", logs: logs, sub: sub}, nil
}

// WatchRollback is a free log subscription operation binding the contract event 0xbaf3b92e813efec2b7525399a930acf56a9ea74f17622f3f1080387356d1c711.
//
// Solidity: e Rollback(sender address, receiver address, value uint256)
func (_Htlc *HtlcFilterer) WatchRollback(opts *bind.WatchOpts, sink chan<- *HtlcRollback) (event.Subscription, error) {

	logs, sub, err := _Htlc.contract.WatchLogs(opts, "Rollback")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HtlcRollback)
				if err := _Htlc.contract.UnpackLog(event, "Rollback", log); err != nil {
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
