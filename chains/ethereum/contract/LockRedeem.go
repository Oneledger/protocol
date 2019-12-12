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
const LockRedeemABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"hasValidatorSigned\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getSignatureCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint256\"},{\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"verifyRedeem\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount_requested\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"}]"

// LockRedeemFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemFuncSigs = map[string]string{
	"bfb9e9f5": "addValidatorProposals(address)",
	"0d8f6b5b": "epochBlockHeight()",
	"45dfa415": "getOLTEthAddress()",
	"6c7d13df": "getSignatureCount(address)",
	"287cc96b": "getTotalEthBalance()",
	"31b6a6d1": "hasValidatorSigned(address)",
	"facd743b": "isValidator(address)",
	"f83d08ba": "lock()",
	"0e7d275d": "newThresholdProposals(uint256)",
	"5d593f8d": "numValidators()",
	"383ea59a": "proposeAddValidator(address)",
	"e0e887d0": "proposeNewThreshold(uint256)",
	"101a8538": "proposeRemoveValidator(address)",
	"db006a75": "redeem(uint256)",
	"0d00753a": "removeValidatorProposals(address)",
	"7cacde3f": "sign(uint256,address)",
	"fa52c7d8": "validators(address)",
	"91e39868": "verifyRedeem(address)",
	"62827733": "votingThreshold()",
}

// LockRedeemBin is the compiled bytecode used for deploying new contracts.
var LockRedeemBin = "0x60806040526170806002553480156200001757600080fd5b50604051620010a9380380620010a9833981018060405260208110156200003d57600080fd5b8101908080516401000000008111156200005657600080fd5b820160208101848111156200006a57600080fd5b81518560208202830111640100000000821117156200008857600080fd5b509093506200009692505050565b60005b815181101562000157576000828281518110620000b257fe5b6020026020010151905060096000826001600160a01b03166001600160a01b03168152602001908152602001600020546000146200013c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f8152602001806200107a602f913960400191505060405180910390fd5b6200014d816200018960201b60201c565b5060010162000099565b5060038151600202816200016757fe5b046001016001819055506200018243620001e660201b60201c565b5062000399565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a250565b60038190556040805182815290517febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3359181900360200190a160005b60075481101562000286576000600782815481106200023c57fe5b600091825260209182902001546001600160a01b031691506200026590829062000189811b901c565b6001600160a01b031660009081526004602052604081205560010162000221565b5062000295600760006200035b565b60005b600854811015620002fd57600060088281548110620002b357fe5b600091825260209182902001546001600160a01b03169150620002dc9082906200030f811b901c565b6001600160a01b031660009081526005602052604081205560010162000298565b506200030c600860006200035b565b50565b6001600160a01b0381166000818152600960205260408082208290558154600019018255517f6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e59190a250565b50805460008255906000526020600020908101906200030c91906200039691905b808211156200039257600081556001016200037c565b5090565b90565b610cd180620003a96000396000f3fe6080604052600436106101145760003560e01c806362827733116100a0578063db006a7511610064578063db006a751461038e578063e0e887d0146103b8578063f83d08ba146103e2578063fa52c7d8146103ea578063facd743b1461041d57610114565b806362827733146102a75780636c7d13df146102bc5780637cacde3f146102ef57806391e3986814610328578063bfb9e9f51461035b57610114565b8063287cc96b116100e7578063287cc96b146101d257806331b6a6d1146101e7578063383ea59a1461022e57806345dfa415146102615780635d593f8d1461029257610114565b80630d00753a146101195780630d8f6b5b1461015e5780630e7d275d14610173578063101a85381461019d575b600080fd5b34801561012557600080fd5b5061014c6004803603602081101561013c57600080fd5b50356001600160a01b0316610450565b60408051918252519081900360200190f35b34801561016a57600080fd5b5061014c610462565b34801561017f57600080fd5b5061014c6004803603602081101561019657600080fd5b5035610468565b3480156101a957600080fd5b506101d0600480360360208110156101c057600080fd5b50356001600160a01b031661047a565b005b3480156101de57600080fd5b5061014c610503565b3480156101f357600080fd5b5061021a6004803603602081101561020a57600080fd5b50356001600160a01b0316610508565b604080519115158252519081900360200190f35b34801561023a57600080fd5b506101d06004803603602081101561025157600080fd5b50356001600160a01b0316610534565b34801561026d57600080fd5b506102766105e4565b604080516001600160a01b039092168252519081900360200190f35b34801561029e57600080fd5b5061014c6105e8565b3480156102b357600080fd5b5061014c6105ee565b3480156102c857600080fd5b5061014c600480360360208110156102df57600080fd5b50356001600160a01b03166105f4565b3480156102fb57600080fd5b506101d06004803603604081101561031257600080fd5b50803590602001356001600160a01b0316610612565b34801561033457600080fd5b5061021a6004803603602081101561034b57600080fd5b50356001600160a01b0316610896565b34801561036757600080fd5b5061014c6004803603602081101561037e57600080fd5b50356001600160a01b03166108b7565b34801561039a57600080fd5b506101d0600480360360208110156103b157600080fd5b50356108c9565b3480156103c457600080fd5b506101d0600480360360208110156103db57600080fd5b5035610a34565b6101d0610b17565b3480156103f657600080fd5b5061014c6004803603602081101561040d57600080fd5b50356001600160a01b0316610b53565b34801561042957600080fd5b5061021a6004803603602081101561044057600080fd5b50356001600160a01b0316610b65565b60056020526000908152604090205481565b60035481565b60066020526000908152604090205481565b336000908152600960205260408120541361049457600080fd5b6001600160a01b0381166000908152600560209081526040808320338452600181019092529091205460ff16156104ff57604051600160e51b62461bcd028152600401808060200182810382526030815260200180610c4a6030913960400191505060405180910390fd5b5050565b303190565b6001600160a01b03166000908152600a6020908152604080832033845260010190915290205460ff1690565b336000908152600960205260408120541361054e57600080fd5b6001600160a01b0381166000908152600460209081526040808320338452600181019092529091205460ff16156105b957604051600160e51b62461bcd02815260040180806020018281038252602c815260200180610c7a602c913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff19168217905581540181556104ff82610b81565b3090565b60005481565b60015481565b6001600160a01b03166000908152600a602052604090206003015490565b61061b33610b65565b61066f5760408051600160e51b62461bcd02815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e206c697374000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a602052604090206004015460ff16156106e35760408051600160e51b62461bcd02815260206004820152601b60248201527f72656465656d207265717565737420697320636f6d706c657465640000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a602052604090206002015482146107555760408051600160e51b62461bcd02815260206004820152601960248201527f72656465656d20616d6f756e7420436f6d70726f6d6973656400000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152600a6020908152604080832033845260010190915290205460ff161561078957600080fd5b6001600160a01b0381166000818152600a6020818152604080842033855260018181018452918520805460ff19168317905594909352526003909101805482019081905590541161084f576001600160a01b038082166000908152600a60205260408082208054600290910154915193169281156108fc0292818181858888f1935050505015801561081f573d6000803e3d6000fd5b506001600160a01b0381166000908152600a602052604081206002810191909155600401805460ff191660011790555b604080513381526020810184905281516001600160a01b038416927f3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db928290030190a25050565b6001600160a01b03166000908152600a602052604090206004015460ff1690565b60046020526000908152604090205481565b336000908152600a6020526040902060020154156108e657600080fd5b6000811161093e5760408051600160e51b62461bcd02815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b336000908152600a602052604090206005015443116109a75760408051600160e51b62461bcd02815260206004820181905260248201527f72657175657374206973206c6f636b65642c206e6f7420617661696c61626c65604482015290519081900360640190fd5b336000818152600a6020908152604080832060048101805460ff19169055600381019390935582546001600160a01b0319169093178083556002808401869055544301600590930192909255825184815292516001600160a01b03909216927f222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2929081900390910190a250565b3360009081526009602052604081205413610a4e57600080fd5b6000548110610a9157604051600160e51b62461bcd028152600401808060200182810382526041815260200180610bdf6041913960600191505060405180910390fd5b6000818152600660209081526040808320338452600181019092529091205460ff1615610af257604051600160e51b62461bcd02815260040180806020018281038252602a815260200180610c20602a913960400191505060405180910390fd5b33600090815260018281016020526040909120805460ff191682179055815401905550565b6040805133815234602082015281517f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427929181900390910190a1565b60096020526000908152604090205481565b6001600160a01b03166000908152600960205260408120541390565b6001600160a01b03811660008181526009602090815260408083206032815583546001019093559154825190815291517fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b9281900390910190a25056fe4e6577207468726573686f6c647320286d29206d757374206265206c657373207468616e20746865206e756d626572206f662076616c696461746f727320286e2973656e6465722068617320616c726561647920766f74656420666f7220746869732070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f20616464207468697320746f2070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f2061646420746869732061646472657373a165627a7a72305820803c9dc1b85152d261790cd36571f13965b82859ea97b1f9573198981d7c51690029666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

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

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) GetSignatureCount(opts *bind.CallOpts, recipient_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "getSignatureCount", recipient_)
	return *ret0, err
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemSession) GetSignatureCount(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.GetSignatureCount(&_LockRedeem.CallOpts, recipient_)
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) GetSignatureCount(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.GetSignatureCount(&_LockRedeem.CallOpts, recipient_)
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

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "hasValidatorSigned", recipient_)
	return *ret0, err
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemSession) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.HasValidatorSigned(&_LockRedeem.CallOpts, recipient_)
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.HasValidatorSigned(&_LockRedeem.CallOpts, recipient_)
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

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "verifyRedeem", recipient_)
	return *ret0, err
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemSession) VerifyRedeem(recipient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.VerifyRedeem(&_LockRedeem.CallOpts, recipient_)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) VerifyRedeem(recipient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.VerifyRedeem(&_LockRedeem.CallOpts, recipient_)
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
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeem *LockRedeemTransactor) Redeem(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "redeem", amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeem *LockRedeemSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.Redeem(&_LockRedeem.TransactOpts, amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeem *LockRedeemTransactorSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeem.Contract.Redeem(&_LockRedeem.TransactOpts, amount_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemTransactor) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "sign", amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Sign(&_LockRedeem.TransactOpts, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeem *LockRedeemTransactorSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Sign(&_LockRedeem.TransactOpts, amount_, recipient_)
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

// LockRedeemRedeemRequestIterator is returned from FilterRedeemRequest and is used to iterate over the raw logs and unpacked data for RedeemRequest events raised by the LockRedeem contract.
type LockRedeemRedeemRequestIterator struct {
	Event *LockRedeemRedeemRequest // Event containing the contract specifics and raw log

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
func (it *LockRedeemRedeemRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemRedeemRequest)
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
		it.Event = new(LockRedeemRedeemRequest)
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
func (it *LockRedeemRedeemRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemRedeemRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemRedeemRequest represents a RedeemRequest event raised by the LockRedeem contract.
type LockRedeemRedeemRequest struct {
	Recepient       common.Address
	AmountRequested *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedeemRequest is a free log retrieval operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_LockRedeem *LockRedeemFilterer) FilterRedeemRequest(opts *bind.FilterOpts, recepient []common.Address) (*LockRedeemRedeemRequestIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemRedeemRequestIterator{contract: _LockRedeem.contract, event: "RedeemRequest", logs: logs, sub: sub}, nil
}

// WatchRedeemRequest is a free log subscription operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_LockRedeem *LockRedeemFilterer) WatchRedeemRequest(opts *bind.WatchOpts, sink chan<- *LockRedeemRedeemRequest, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemRedeemRequest)
				if err := _LockRedeem.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
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
func (_LockRedeem *LockRedeemFilterer) ParseRedeemRequest(log types.Log) (*LockRedeemRedeemRequest, error) {
	event := new(LockRedeemRedeemRequest)
	if err := _LockRedeem.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
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
	Recipient         common.Address
	ValidatorAddresss common.Address
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeem *LockRedeemFilterer) FilterValidatorSignedRedeem(opts *bind.FilterOpts, recipient []common.Address) (*LockRedeemValidatorSignedRedeemIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemValidatorSignedRedeemIterator{contract: _LockRedeem.contract, event: "ValidatorSignedRedeem", logs: logs, sub: sub}, nil
}

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeem *LockRedeemFilterer) WatchValidatorSignedRedeem(opts *bind.WatchOpts, sink chan<- *LockRedeemValidatorSignedRedeem, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "ValidatorSignedRedeem", recipientRule)
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

// ParseValidatorSignedRedeem is a log parse operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_LockRedeem *LockRedeemFilterer) ParseValidatorSignedRedeem(log types.Log) (*LockRedeemValidatorSignedRedeem, error) {
	event := new(LockRedeemValidatorSignedRedeem)
	if err := _LockRedeem.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}
