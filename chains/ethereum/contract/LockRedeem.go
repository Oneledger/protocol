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

// LockRedeemABI is the input ABI used to generate the binding from.
const LockRedeemABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"}],\"name\":\"getRedeemAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_trafered\",\"type\":\"uint256\"}],\"name\":\"RedeemSuccessful\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"}]"

// LockRedeemFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemFuncSigs = map[string]string{
	"bfb9e9f5": "addValidatorProposals(address)",
	"0d8f6b5b": "epochBlockHeight()",
	"45dfa415": "getOLTEthAddress()",
	"2c3ee88c": "getRedeemAmount(uint256)",
	"287cc96b": "getTotalEthBalance()",
	"facd743b": "isValidator(address)",
	"f83d08ba": "lock()",
	"0e7d275d": "newThresholdProposals(uint256)",
	"5d593f8d": "numValidators()",
	"383ea59a": "proposeAddValidator(address)",
	"e0e887d0": "proposeNewThreshold(uint256)",
	"101a8538": "proposeRemoveValidator(address)",
	"db006a75": "redeem(uint256)",
	"0d00753a": "removeValidatorProposals(address)",
	"d5f37f95": "sign(uint256,uint256,address)",
	"fa52c7d8": "validators(address)",
	"62827733": "votingThreshold()",
}

// LockRedeemBin is the compiled bytecode used for deploying new contracts.
var LockRedeemBin = "0x60806040523480156200001157600080fd5b506040516200102b3803806200102b833981810160405260208110156200003757600080fd5b81019080805160405193929190846401000000008211156200005857600080fd5b9083019060208201858111156200006e57600080fd5b82518660208202830111640100000000821117156200008c57600080fd5b82525081516020918201928201910280838360005b83811015620000bb578181015183820152602001620000a1565b505050509190910160405250620000d192505050565b60005b815181101562000195576000828281518110620000ed57fe5b6020026020010151905060086000826001600160a01b03166001600160a01b031681526020019081526020016000205460001462000177576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f81526020018062000ffc602f913960400191505060405180910390fd5b6200018b816001600160e01b03620001b116565b50600101620000d4565b50620001aa436001600160e01b036200020e16565b50620003e0565b6001600160a01b03811660008181526008602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a250565b33600090815260086020526040812054136200022957600080fd5b60028190556040805182815290517febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359181900360200190a160005b600654811015620002cb576000600682815481106200027f57fe5b6000918252602090912001546001600160a01b03169050620002aa816001600160e01b03620001b116565b6001600160a01b031660009081526003602052604081205560010162000264565b50620002da60066000620003a2565b60005b6007548110156200034457600060078281548110620002f857fe5b6000918252602090912001546001600160a01b0316905062000323816001600160e01b036200035616565b6001600160a01b0316600090815260046020526040812055600101620002dd565b506200035360076000620003a2565b50565b6001600160a01b0381166000818152600860205260408082208290558154600019018255517f6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e59190a250565b5080546000825590600052602060002090810190620003539190620003dd91905b80821115620003d95760008155600101620003c3565b5090565b90565b610c0c80620003f06000396000f3fe6080604052600436106100fe5760003560e01c80635d593f8d11610095578063db006a7511610064578063db006a75146102fb578063e0e887d014610325578063f83d08ba1461034f578063fa52c7d814610357578063facd743b1461038a576100fe565b80635d593f8d1461025f5780636282773314610274578063bfb9e9f514610289578063d5f37f95146102bc576100fe565b8063287cc96b116100d1578063287cc96b146101bc5780632c3ee88c146101d1578063383ea59a146101fb57806345dfa4151461022e576100fe565b80630d00753a146101035780630d8f6b5b146101485780630e7d275d1461015d578063101a853814610187575b600080fd5b34801561010f57600080fd5b506101366004803603602081101561012657600080fd5b50356001600160a01b03166103d1565b60408051918252519081900360200190f35b34801561015457600080fd5b506101366103e3565b34801561016957600080fd5b506101366004803603602081101561018057600080fd5b50356103e9565b34801561019357600080fd5b506101ba600480360360208110156101aa57600080fd5b50356001600160a01b03166103fb565b005b3480156101c857600080fd5b50610136610481565b3480156101dd57600080fd5b50610136600480360360208110156101f457600080fd5b5035610486565b34801561020757600080fd5b506101ba6004803603602081101561021e57600080fd5b50356001600160a01b031661049b565b34801561023a57600080fd5b50610243610548565b604080516001600160a01b039092168252519081900360200190f35b34801561026b57600080fd5b5061013661054c565b34801561028057600080fd5b50610136610552565b34801561029557600080fd5b50610136600480360360208110156102ac57600080fd5b50356001600160a01b0316610558565b3480156102c857600080fd5b506101ba600480360360608110156102df57600080fd5b50803590602081013590604001356001600160a01b031661056a565b34801561030757600080fd5b506101ba6004803603602081101561031e57600080fd5b5035610729565b34801561033157600080fd5b506101ba6004803603602081101561034857600080fd5b50356108d5565b6101ba6109b2565b34801561036357600080fd5b506101366004803603602081101561037a57600080fd5b50356001600160a01b03166109ee565b34801561039657600080fd5b506103bd600480360360208110156103ad57600080fd5b50356001600160a01b0316610a00565b604080519115158252519081900360200190f35b60046020526000908152604090205481565b60025481565b60056020526000908152604090205481565b336000908152600860205260408120541361041557600080fd5b6001600160a01b0381166000908152600460209081526040808320338452600181019092529091205460ff161561047d5760405162461bcd60e51b8152600401808060200182810382526030815260200180610b7c6030913960400191505060405180910390fd5b5050565b303190565b60009081526009602052604090206001015490565b33600090815260086020526040812054136104b557600080fd5b6001600160a01b0381166000908152600360209081526040808320338452600181019092529091205460ff161561051d5760405162461bcd60e51b815260040180806020018281038252602c815260200180610bac602c913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401815561047d82610a1c565b3090565b60005481565b60015481565b60036020526000908152604090205481565b61057333610a00565b6105c4576040805162461bcd60e51b815260206004820152601e60248201527f76616c696461746f72206e6f74207072657373656e7420696e206c6973740000604482015290519081900360640190fd5b600083815260096020526040902060020154156106b9576000838152600960205260409020600101548214610637576040805162461bcd60e51b815260206004820152601460248201527315985b1a59185d1bdc90dbdb5c1c9bdb5a5cd95960621b604482015290519081900360640190fd5b6000838152600960205260409020546001600160a01b0382811691161461069c576040805162461bcd60e51b815260206004820152601460248201527315985b1a59185d1bdc90dbdb5c1c9bdb5a5cd95960621b604482015290519081900360640190fd5b6000838152600960205260409020600201805460010190556106f9565b6000838152600960205260409020600180820184905581546001600160a01b0319166001600160a01b0384161782556002820155600301805460ff191690555b60405133907f27a7fc932a73cc79830d493313d5b19cb1828092eeb13a3c33f5309932ee8c8490600090a2505050565b6000818152600960205260409020546001600160a01b0316331461077e5760405162461bcd60e51b8152600401808060200182810382526043815260200180610b0f6043913960600191505060405180910390fd5b60008181526009602052604090206003015460ff16156107cf5760405162461bcd60e51b8152600401808060200182810382526028815260200180610abb6028913960400191505060405180910390fd5b60015460008281526009602052604090206002015410156108215760405162461bcd60e51b815260040180806020018281038252602c815260200180610ae3602c913960400191505060405180910390fd5b600081815260096020526040808220805460019091015491516001600160a01b039091169282156108fc02929190818181858888f1935050505015801561086c573d6000803e3d6000fd5b5060008181526009602090815260409182902060038101805460ff191660019081179091558154910154835190815292516001600160a01b03909116927f80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad192908290030190a250565b33600090815260086020526040812054136108ef57600080fd5b600054811061092f5760405162461bcd60e51b8152600401808060200182810382526041815260200180610a7a6041913960600191505060405180910390fd5b6000818152600560209081526040808320338452600181019092529091205460ff161561098d5760405162461bcd60e51b815260040180806020018281038252602a815260200180610b52602a913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401905550565b6040805133815234602082015281517f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427929181900390910190a1565b60086020526000908152604090205481565b6001600160a01b03166000908152600860205260408120541390565b6001600160a01b03811660008181526008602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a25056fe4e6577207468726573686f6c647320286d29206d757374206265206c657373207468616e20746865206e756d626572206f662076616c696461746f727320286e2952656465656d20616c7265616479206578656375746564206f6e20746869732072656465656d49444e6f7420656e6f7567682056616c696461746f7220766f74657320746f20657865637574652052656465656d52656465656d2063616e206f6e6c7920626520696e746974617465642062792074686520726563657069656e74206f662072656465656d207472616e73616374696f6e73656e6465722068617320616c726561647920766f74656420666f7220746869732070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f20616464207468697320746f2070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f2061646420746869732061646472657373a265627a7a723158209efbb8829706652cddab81d1bf4d5b2d26b36002d489d5d7304832d07feaf52a64736f6c634300050b0032666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

// DeployLockRedeem deploys a new Ethereum contract, binding an instance of LockRedeem to it.
func DeployLockRedeem(auth *bind.TransactOpts, backend bind.ContractBackend, initialValidators []common.Address) (common.Address, *types.Transaction, *LockRedeem, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockRedeemBin), backend, initialValidators)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockRedeem{LockRedeemCaller: LockRedeemCaller{contract: contract}, LockRedeemTransactor: LockRedeemTransactor{contract: contract}, LockRedeemFilterer: LockRedeemFilterer{contract: contract}}, nil
}

// LockRedeem is an auto generated Go binding around an Ethereum contract.
type LockRedeem struct {
	LockRedeemCaller     // Read-only binding to the contract
	LockRedeemTransactor // Write-only binding to the contract
	LockRedeemFilterer   // Log filterer for contract events
}

// LockRedeemCaller is an auto generated read-only Go binding around an Ethereum contract.
type LockRedeemCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LockRedeemTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockRedeemFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockRedeemSession struct {
	Contract     *LockRedeem       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockRedeemCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockRedeemCallerSession struct {
	Contract *LockRedeemCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// LockRedeemTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockRedeemTransactorSession struct {
	Contract     *LockRedeemTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// LockRedeemRaw is an auto generated low-level Go binding around an Ethereum contract.
type LockRedeemRaw struct {
	Contract *LockRedeem // Generic contract binding to access the raw methods on
}

// LockRedeemCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockRedeemCallerRaw struct {
	Contract *LockRedeemCaller // Generic read-only contract binding to access the raw methods on
}

// LockRedeemTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockRedeemTransactorRaw struct {
	Contract *LockRedeemTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLockRedeem creates a new instance of LockRedeem, bound to a specific deployed contract.
func NewLockRedeem(address common.Address, backend bind.ContractBackend) (*LockRedeem, error) {
	contract, err := bindLockRedeem(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockRedeem{LockRedeemCaller: LockRedeemCaller{contract: contract}, LockRedeemTransactor: LockRedeemTransactor{contract: contract}, LockRedeemFilterer: LockRedeemFilterer{contract: contract}}, nil
}

// NewLockRedeemCaller creates a new read-only instance of LockRedeem, bound to a specific deployed contract.
func NewLockRedeemCaller(address common.Address, caller bind.ContractCaller) (*LockRedeemCaller, error) {
	contract, err := bindLockRedeem(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemCaller{contract: contract}, nil
}

// NewLockRedeemTransactor creates a new write-only instance of LockRedeem, bound to a specific deployed contract.
func NewLockRedeemTransactor(address common.Address, transactor bind.ContractTransactor) (*LockRedeemTransactor, error) {
	contract, err := bindLockRedeem(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemTransactor{contract: contract}, nil
}

// NewLockRedeemFilterer creates a new log filterer instance of LockRedeem, bound to a specific deployed contract.
func NewLockRedeemFilterer(address common.Address, filterer bind.ContractFilterer) (*LockRedeemFilterer, error) {
	contract, err := bindLockRedeem(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockRedeemFilterer{contract: contract}, nil
}

// bindLockRedeem binds a generic wrapper to an already deployed contract.
func bindLockRedeem(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeem *LockRedeemRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeem.Contract.LockRedeemCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeem *LockRedeemRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeem.Contract.LockRedeemTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeem *LockRedeemRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeem.Contract.LockRedeemTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeem *LockRedeemCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeem.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeem *LockRedeemTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeem.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeem *LockRedeemTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeem.Contract.contract.Transact(opts, method, params...)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCaller) AddValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "addValidatorProposals", arg0)
	return *ret0, err
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.AddValidatorProposals(&_LockRedeem.CallOpts, arg0)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCallerSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.AddValidatorProposals(&_LockRedeem.CallOpts, arg0)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) EpochBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "epochBlockHeight")
	return *ret0, err
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeem *LockRedeemSession) EpochBlockHeight() (*big.Int, error) {
	return _LockRedeem.Contract.EpochBlockHeight(&_LockRedeem.CallOpts)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) EpochBlockHeight() (*big.Int, error) {
	return _LockRedeem.Contract.EpochBlockHeight(&_LockRedeem.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeem *LockRedeemCaller) GetOLTEthAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "getOLTEthAddress")
	return *ret0, err
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeem *LockRedeemSession) GetOLTEthAddress() (common.Address, error) {
	return _LockRedeem.Contract.GetOLTEthAddress(&_LockRedeem.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeem *LockRedeemCallerSession) GetOLTEthAddress() (common.Address, error) {
	return _LockRedeem.Contract.GetOLTEthAddress(&_LockRedeem.CallOpts)
}

// GetRedeemAmount is a free data retrieval call binding the contract method 0x2c3ee88c.
//
// Solidity: function getRedeemAmount(uint256 redeemID_) constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) GetRedeemAmount(opts *bind.CallOpts, redeemID_ *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "getRedeemAmount", redeemID_)
	return *ret0, err
}

// GetRedeemAmount is a free data retrieval call binding the contract method 0x2c3ee88c.
//
// Solidity: function getRedeemAmount(uint256 redeemID_) constant returns(uint256)
func (_LockRedeem *LockRedeemSession) GetRedeemAmount(redeemID_ *big.Int) (*big.Int, error) {
	return _LockRedeem.Contract.GetRedeemAmount(&_LockRedeem.CallOpts, redeemID_)
}

// GetRedeemAmount is a free data retrieval call binding the contract method 0x2c3ee88c.
//
// Solidity: function getRedeemAmount(uint256 redeemID_) constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) GetRedeemAmount(redeemID_ *big.Int) (*big.Int, error) {
	return _LockRedeem.Contract.GetRedeemAmount(&_LockRedeem.CallOpts, redeemID_)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) GetTotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "getTotalEthBalance")
	return *ret0, err
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeem *LockRedeemSession) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeem.Contract.GetTotalEthBalance(&_LockRedeem.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeem.Contract.GetTotalEthBalance(&_LockRedeem.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "isValidator", addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeem *LockRedeemSession) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeem.Contract.IsValidator(&_LockRedeem.CallOpts, addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeem.Contract.IsValidator(&_LockRedeem.CallOpts, addr)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCaller) NewThresholdProposals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "newThresholdProposals", arg0)
	return *ret0, err
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _LockRedeem.Contract.NewThresholdProposals(&_LockRedeem.CallOpts, arg0)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCallerSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _LockRedeem.Contract.NewThresholdProposals(&_LockRedeem.CallOpts, arg0)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeem *LockRedeemSession) NumValidators() (*big.Int, error) {
	return _LockRedeem.Contract.NumValidators(&_LockRedeem.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) NumValidators() (*big.Int, error) {
	return _LockRedeem.Contract.NumValidators(&_LockRedeem.CallOpts)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeem *LockRedeemCaller) ProposeRemoveValidator(opts *bind.CallOpts, v common.Address) error {
	var ()
	out := &[]interface{}{}
	err := _LockRedeem.contract.Call(opts, out, "proposeRemoveValidator", v)
	return err
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeem *LockRedeemSession) ProposeRemoveValidator(v common.Address) error {
	return _LockRedeem.Contract.ProposeRemoveValidator(&_LockRedeem.CallOpts, v)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_LockRedeem *LockRedeemCallerSession) ProposeRemoveValidator(v common.Address) error {
	return _LockRedeem.Contract.ProposeRemoveValidator(&_LockRedeem.CallOpts, v)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCaller) RemoveValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "removeValidatorProposals", arg0)
	return *ret0, err
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.RemoveValidatorProposals(&_LockRedeem.CallOpts, arg0)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_LockRedeem *LockRedeemCallerSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.RemoveValidatorProposals(&_LockRedeem.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeem *LockRedeemCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeem *LockRedeemSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.Validators(&_LockRedeem.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_LockRedeem *LockRedeemCallerSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.Validators(&_LockRedeem.CallOpts, arg0)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) VotingThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "votingThreshold")
	return *ret0, err
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeem *LockRedeemSession) VotingThreshold() (*big.Int, error) {
	return _LockRedeem.Contract.VotingThreshold(&_LockRedeem.CallOpts)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) VotingThreshold() (*big.Int, error) {
	return _LockRedeem.Contract.VotingThreshold(&_LockRedeem.CallOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeem *LockRedeemTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeem *LockRedeemSession) Lock() (*types.Transaction, error) {
	return _LockRedeem.Contract.Lock(&_LockRedeem.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeem *LockRedeemTransactorSession) Lock() (*types.Transaction, error) {
	return _LockRedeem.Contract.Lock(&_LockRedeem.TransactOpts)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeem *LockRedeemTransactor) ProposeAddValidator(opts *bind.TransactOpts, v common.Address) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "proposeAddValidator", v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeem *LockRedeemSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.ProposeAddValidator(&_LockRedeem.TransactOpts, v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_LockRedeem *LockRedeemTransactorSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.ProposeAddValidator(&_LockRedeem.TransactOpts, v)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeem *LockRedeemTransactor) ProposeNewThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "proposeNewThreshold", threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeem *LockRedeemSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.ProposeNewThreshold(&_LockRedeem.TransactOpts, threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_LockRedeem *LockRedeemTransactorSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.ProposeNewThreshold(&_LockRedeem.TransactOpts, threshold)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 redeemID_) returns()
func (_LockRedeem *LockRedeemTransactor) Redeem(opts *bind.TransactOpts, redeemID_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "redeem", redeemID_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 redeemID_) returns()
func (_LockRedeem *LockRedeemSession) Redeem(redeemID_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.Redeem(&_LockRedeem.TransactOpts, redeemID_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 redeemID_) returns()
func (_LockRedeem *LockRedeemTransactorSession) Redeem(redeemID_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.Redeem(&_LockRedeem.TransactOpts, redeemID_)
}

// Sign is a paid mutator transaction binding the contract method 0xd5f37f95.
//
// Solidity: function sign(uint256 redeemID_, uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemTransactor) Sign(opts *bind.TransactOpts, redeemID_ *big.Int, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "sign", redeemID_, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0xd5f37f95.
//
// Solidity: function sign(uint256 redeemID_, uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemSession) Sign(redeemID_ *big.Int, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Sign(&_LockRedeem.TransactOpts, redeemID_, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0xd5f37f95.
//
// Solidity: function sign(uint256 redeemID_, uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemTransactorSession) Sign(redeemID_ *big.Int, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Sign(&_LockRedeem.TransactOpts, redeemID_, amount_, recipient_)
}

// LockRedeemAddValidatorIterator is returned from FilterAddValidator and is used to iterate over the raw logs and unpacked data for AddValidator events raised by the LockRedeem contract.
type LockRedeemAddValidatorIterator struct {
	Event *LockRedeemAddValidator // Event containing the contract specifics and raw log

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
func (it *LockRedeemAddValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemAddValidator)
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
		it.Event = new(LockRedeemAddValidator)
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
func (it *LockRedeemAddValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemAddValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemAddValidator represents a AddValidator event raised by the LockRedeem contract.
type LockRedeemAddValidator struct {
	Address common.Address
	Power   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_LockRedeem *LockRedeemFilterer) FilterAddValidator(opts *bind.FilterOpts, _address []common.Address) (*LockRedeemAddValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemAddValidatorIterator{contract: _LockRedeem.contract, event: "AddValidator", logs: logs, sub: sub}, nil
}

// WatchAddValidator is a free log subscription operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_LockRedeem *LockRedeemFilterer) WatchAddValidator(opts *bind.WatchOpts, sink chan<- *LockRedeemAddValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemAddValidator)
				if err := _LockRedeem.contract.UnpackLog(event, "AddValidator", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseAddValidator(log types.Log) (*LockRedeemAddValidator, error) {
	event := new(LockRedeemAddValidator)
	if err := _LockRedeem.contract.UnpackLog(event, "AddValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemDeleteValidatorIterator is returned from FilterDeleteValidator and is used to iterate over the raw logs and unpacked data for DeleteValidator events raised by the LockRedeem contract.
type LockRedeemDeleteValidatorIterator struct {
	Event *LockRedeemDeleteValidator // Event containing the contract specifics and raw log

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
func (it *LockRedeemDeleteValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemDeleteValidator)
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
		it.Event = new(LockRedeemDeleteValidator)
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
func (it *LockRedeemDeleteValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemDeleteValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemDeleteValidator represents a DeleteValidator event raised by the LockRedeem contract.
type LockRedeemDeleteValidator struct {
	Address common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDeleteValidator is a free log retrieval operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_LockRedeem *LockRedeemFilterer) FilterDeleteValidator(opts *bind.FilterOpts, _address []common.Address) (*LockRedeemDeleteValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemDeleteValidatorIterator{contract: _LockRedeem.contract, event: "DeleteValidator", logs: logs, sub: sub}, nil
}

// WatchDeleteValidator is a free log subscription operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_LockRedeem *LockRedeemFilterer) WatchDeleteValidator(opts *bind.WatchOpts, sink chan<- *LockRedeemDeleteValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemDeleteValidator)
				if err := _LockRedeem.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseDeleteValidator(log types.Log) (*LockRedeemDeleteValidator, error) {
	event := new(LockRedeemDeleteValidator)
	if err := _LockRedeem.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemLockIterator is returned from FilterLock and is used to iterate over the raw logs and unpacked data for Lock events raised by the LockRedeem contract.
type LockRedeemLockIterator struct {
	Event *LockRedeemLock // Event containing the contract specifics and raw log

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
func (it *LockRedeemLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemLock)
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
		it.Event = new(LockRedeemLock)
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
func (it *LockRedeemLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemLock represents a Lock event raised by the LockRedeem contract.
type LockRedeemLock struct {
	Sender         common.Address
	AmountReceived *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_LockRedeem *LockRedeemFilterer) FilterLock(opts *bind.FilterOpts) (*LockRedeemLockIterator, error) {

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return &LockRedeemLockIterator{contract: _LockRedeem.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_LockRedeem *LockRedeemFilterer) WatchLock(opts *bind.WatchOpts, sink chan<- *LockRedeemLock) (event.Subscription, error) {

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemLock)
				if err := _LockRedeem.contract.UnpackLog(event, "Lock", log); err != nil {
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

// ParseLock is a log parse operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_LockRedeem *LockRedeemFilterer) ParseLock(log types.Log) (*LockRedeemLock, error) {
	event := new(LockRedeemLock)
	if err := _LockRedeem.contract.UnpackLog(event, "Lock", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the LockRedeem contract.
type LockRedeemNewEpochIterator struct {
	Event *LockRedeemNewEpoch // Event containing the contract specifics and raw log

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
func (it *LockRedeemNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemNewEpoch)
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
		it.Event = new(LockRedeemNewEpoch)
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
func (it *LockRedeemNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemNewEpoch represents a NewEpoch event raised by the LockRedeem contract.
type LockRedeemNewEpoch struct {
	EpochHeight *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_LockRedeem *LockRedeemFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*LockRedeemNewEpochIterator, error) {

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &LockRedeemNewEpochIterator{contract: _LockRedeem.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_LockRedeem *LockRedeemFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *LockRedeemNewEpoch) (event.Subscription, error) {

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemNewEpoch)
				if err := _LockRedeem.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseNewEpoch(log types.Log) (*LockRedeemNewEpoch, error) {
	event := new(LockRedeemNewEpoch)
	if err := _LockRedeem.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemNewThresholdIterator is returned from FilterNewThreshold and is used to iterate over the raw logs and unpacked data for NewThreshold events raised by the LockRedeem contract.
type LockRedeemNewThresholdIterator struct {
	Event *LockRedeemNewThreshold // Event containing the contract specifics and raw log

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
func (it *LockRedeemNewThresholdIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemNewThreshold)
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
		it.Event = new(LockRedeemNewThreshold)
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
func (it *LockRedeemNewThresholdIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemNewThresholdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemNewThreshold represents a NewThreshold event raised by the LockRedeem contract.
type LockRedeemNewThreshold struct {
	PrevThreshold *big.Int
	NewThreshold  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewThreshold is a free log retrieval operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_LockRedeem *LockRedeemFilterer) FilterNewThreshold(opts *bind.FilterOpts) (*LockRedeemNewThresholdIterator, error) {

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return &LockRedeemNewThresholdIterator{contract: _LockRedeem.contract, event: "NewThreshold", logs: logs, sub: sub}, nil
}

// WatchNewThreshold is a free log subscription operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_LockRedeem *LockRedeemFilterer) WatchNewThreshold(opts *bind.WatchOpts, sink chan<- *LockRedeemNewThreshold) (event.Subscription, error) {

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemNewThreshold)
				if err := _LockRedeem.contract.UnpackLog(event, "NewThreshold", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseNewThreshold(log types.Log) (*LockRedeemNewThreshold, error) {
	event := new(LockRedeemNewThreshold)
	if err := _LockRedeem.contract.UnpackLog(event, "NewThreshold", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemRedeemSuccessfulIterator is returned from FilterRedeemSuccessful and is used to iterate over the raw logs and unpacked data for RedeemSuccessful events raised by the LockRedeem contract.
type LockRedeemRedeemSuccessfulIterator struct {
	Event *LockRedeemRedeemSuccessful // Event containing the contract specifics and raw log

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
func (it *LockRedeemRedeemSuccessfulIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemRedeemSuccessful)
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
		it.Event = new(LockRedeemRedeemSuccessful)
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
func (it *LockRedeemRedeemSuccessfulIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemRedeemSuccessfulIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemRedeemSuccessful represents a RedeemSuccessful event raised by the LockRedeem contract.
type LockRedeemRedeemSuccessful struct {
	Recepient      common.Address
	AmountTrafered *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRedeemSuccessful is a free log retrieval operation binding the contract event 0x80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad1.
//
// Solidity: event RedeemSuccessful(address indexed recepient, uint256 amount_trafered)
func (_LockRedeem *LockRedeemFilterer) FilterRedeemSuccessful(opts *bind.FilterOpts, recepient []common.Address) (*LockRedeemRedeemSuccessfulIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "RedeemSuccessful", recepientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemRedeemSuccessfulIterator{contract: _LockRedeem.contract, event: "RedeemSuccessful", logs: logs, sub: sub}, nil
}

// WatchRedeemSuccessful is a free log subscription operation binding the contract event 0x80cfc930fa1029f5fdb639588b474e55c8051b1a9b635f90fe3af3508cfd8ad1.
//
// Solidity: event RedeemSuccessful(address indexed recepient, uint256 amount_trafered)
func (_LockRedeem *LockRedeemFilterer) WatchRedeemSuccessful(opts *bind.WatchOpts, sink chan<- *LockRedeemRedeemSuccessful, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "RedeemSuccessful", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemRedeemSuccessful)
				if err := _LockRedeem.contract.UnpackLog(event, "RedeemSuccessful", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseRedeemSuccessful(log types.Log) (*LockRedeemRedeemSuccessful, error) {
	event := new(LockRedeemRedeemSuccessful)
	if err := _LockRedeem.contract.UnpackLog(event, "RedeemSuccessful", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemValidatorSignedRedeemIterator is returned from FilterValidatorSignedRedeem and is used to iterate over the raw logs and unpacked data for ValidatorSignedRedeem events raised by the LockRedeem contract.
type LockRedeemValidatorSignedRedeemIterator struct {
	Event *LockRedeemValidatorSignedRedeem // Event containing the contract specifics and raw log

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
func (it *LockRedeemValidatorSignedRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemValidatorSignedRedeem)
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
		it.Event = new(LockRedeemValidatorSignedRedeem)
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
func (it *LockRedeemValidatorSignedRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemValidatorSignedRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemValidatorSignedRedeem represents a ValidatorSignedRedeem event raised by the LockRedeem contract.
type LockRedeemValidatorSignedRedeem struct {
	ValidatorAddresss common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x27a7fc932a73cc79830d493313d5b19cb1828092eeb13a3c33f5309932ee8c84.
//
// Solidity: event ValidatorSignedRedeem(address indexed validator_addresss)
func (_LockRedeem *LockRedeemFilterer) FilterValidatorSignedRedeem(opts *bind.FilterOpts, validator_addresss []common.Address) (*LockRedeemValidatorSignedRedeemIterator, error) {

	var validator_addresssRule []interface{}
	for _, validator_addresssItem := range validator_addresss {
		validator_addresssRule = append(validator_addresssRule, validator_addresssItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "ValidatorSignedRedeem", validator_addresssRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemValidatorSignedRedeemIterator{contract: _LockRedeem.contract, event: "ValidatorSignedRedeem", logs: logs, sub: sub}, nil
}

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x27a7fc932a73cc79830d493313d5b19cb1828092eeb13a3c33f5309932ee8c84.
//
// Solidity: event ValidatorSignedRedeem(address indexed validator_addresss)
func (_LockRedeem *LockRedeemFilterer) WatchValidatorSignedRedeem(opts *bind.WatchOpts, sink chan<- *LockRedeemValidatorSignedRedeem, validator_addresss []common.Address) (event.Subscription, error) {

	var validator_addresssRule []interface{}
	for _, validator_addresssItem := range validator_addresss {
		validator_addresssRule = append(validator_addresssRule, validator_addresssItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "ValidatorSignedRedeem", validator_addresssRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemValidatorSignedRedeem)
				if err := _LockRedeem.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
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

// ParseValidatorSignedRedeem is a log parse operation binding the contract event 0x27a7fc932a73cc79830d493313d5b19cb1828092eeb13a3c33f5309932ee8c84.
//
// Solidity: event ValidatorSignedRedeem(address indexed validator_addresss)
func (_LockRedeem *LockRedeemFilterer) ParseValidatorSignedRedeem(log types.Log) (*LockRedeemValidatorSignedRedeem, error) {
	event := new(LockRedeemValidatorSignedRedeem)
	if err := _LockRedeem.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}
