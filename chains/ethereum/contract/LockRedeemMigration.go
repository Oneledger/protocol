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
const LockRedeemABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getSignatureCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"hasValidatorSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recepient_\",\"type\":\"address\"}],\"name\":\"isredeemAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newSmartContractAddress\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"migrationSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"migrationSigners\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"verifyRedeem\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LockRedeemFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemFuncSigs = map[string]string{
	"45dfa415": "getOLTEthAddress()",
	"6c7d13df": "getSignatureCount(address)",
	"287cc96b": "getTotalEthBalance()",
	"31b6a6d1": "hasValidatorSigned(address)",
	"facd743b": "isValidator(address)",
	"2138c6b9": "isredeemAvailable(address)",
	"f83d08ba": "lock()",
	"ce5494bb": "migrate(address)",
	"27882c3a": "migrationSignatures()",
	"a04d0498": "migrationSigners(address)",
	"5d593f8d": "numValidators()",
	"db006a75": "redeem(uint256)",
	"7cacde3f": "sign(uint256,address)",
	"fa52c7d8": "validators(address)",
	"91e39868": "verifyRedeem(address)",
}

// LockRedeemBin is the compiled bytecode used for deploying new contracts.
var LockRedeemBin = "0x60806040526000805460ff1916905561708060075534801561002057600080fd5b50604051610ea8380380610ea88339818101604052602081101561004357600080fd5b810190808051604051939291908464010000000082111561006357600080fd5b90830190602082018581111561007857600080fd5b825186602082028301116401000000008211171561009557600080fd5b82525081516020918201928201910280838360005b838110156100c25781810151838201526020016100aa565b5050505091909101604052506100d792505050565b60005b815181101561017a5760008282815181106100f157fe5b6020026020010151905060026000826001600160a01b03166001600160a01b031681526020019081526020016000205460001461015f5760405162461bcd60e51b815260040180806020018281038252602f815260200180610e79602f913960400191505060405180910390fd5b610171816001600160e01b036101a816565b506001016100da565b506000805460ff191660011790558051600390600202046001016003908155815104600101600455506101f5565b6001600160a01b038116600081815260026020526040808220603290556001805481019055517f6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad369190a250565b610c75806102046000396000f3fe6080604052600436106100e85760003560e01c80637cacde3f1161008a578063db006a7511610059578063db006a75146102f0578063f83d08ba1461031a578063fa52c7d814610322578063facd743b14610355576100e8565b80637cacde3f1461021c57806391e3986814610257578063a04d04981461028a578063ce5494bb146102bd576100e8565b806331b6a6d1116100c657806331b6a6d11461017057806345dfa415146101a35780635d593f8d146101d45780636c7d13df146101e9576100e8565b80632138c6b9146100ed57806327882c3a14610134578063287cc96b1461015b575b600080fd5b3480156100f957600080fd5b506101206004803603602081101561011057600080fd5b50356001600160a01b0316610388565b604080519115158252519081900360200190f35b34801561014057600080fd5b506101496103d0565b60408051918252519081900360200190f35b34801561016757600080fd5b506101496103d6565b34801561017c57600080fd5b506101206004803603602081101561019357600080fd5b50356001600160a01b03166103da565b3480156101af57600080fd5b506101b8610406565b604080516001600160a01b039092168252519081900360200190f35b3480156101e057600080fd5b5061014961040a565b3480156101f557600080fd5b506101496004803603602081101561020c57600080fd5b50356001600160a01b0316610410565b34801561022857600080fd5b506102556004803603604081101561023f57600080fd5b50803590602001356001600160a01b031661042e565b005b34801561026357600080fd5b506101206004803603602081101561027a57600080fd5b50356001600160a01b0316610723565b34801561029657600080fd5b50610120600480360360208110156102ad57600080fd5b50356001600160a01b0316610744565b3480156102c957600080fd5b50610255600480360360208110156102e057600080fd5b50356001600160a01b0316610759565b3480156102fc57600080fd5b506102556004803603602081101561031357600080fd5b5035610a2d565b610255610ba2565b34801561032e57600080fd5b506101496004803603602081101561034557600080fd5b50356001600160a01b0316610bed565b34801561036157600080fd5b506101206004803603602081101561037857600080fd5b50356001600160a01b0316610bff565b6001600160a01b03811660009081526008602052604081206005015415806103ca57506001600160a01b03821660009081526008602052604090206005015443115b92915050565b60055481565b4790565b6001600160a01b0316600090815260086020908152604080832033845260010190915290205460ff1690565b3090565b60015481565b6001600160a01b031660009081526008602052604090206003015490565b60005460ff1661043d57600080fd5b3360009081526002602052604090205461049e576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b6001600160a01b03811660009081526008602052604090206004015460ff161561050f576040805162461bcd60e51b815260206004820152601b60248201527f72656465656d207265717565737420697320636f6d706c657465640000000000604482015290519081900360640190fd5b6001600160a01b038116600090815260086020526040902060020154821461057e576040805162461bcd60e51b815260206004820152601960248201527f72656465656d20616d6f756e7420436f6d70726f6d6973656400000000000000604482015290519081900360640190fd5b6001600160a01b038116600090815260086020908152604080832033845260010190915290205460ff16156105b257600080fd5b6001600160a01b038116600081815260086020818152604080842033855260018181018452918520805460ff191683179055949093525260039182018054909101908190559054116106dc576001600160a01b0380821660009081526008602052604080822080546002909101549151929316918381818185875af1925050503d806000811461065e576040519150601f19603f3d011682016040523d82523d6000602084013e610663565b606091505b50509050806106ac576040805162461bcd60e51b815260206004820152601060248201526f2a3930b739b332b9103330b4b632b21760811b604482015290519081900360640190fd5b506001600160a01b03811660009081526008602052604081206002810191909155600401805460ff191660011790555b604080513381526020810184905281516001600160a01b038416927f3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db928290030190a25050565b6001600160a01b031660009081526008602052604090206004015460ff1690565b60066020526000908152604090205460ff1681565b336000908152600260205260409020546107ba576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b3360009081526006602052604090205460ff161561081f576040805162461bcd60e51b815260206004820152601860248201527f56616c696461746f72205369676e656420616c72656164790000000000000000604482015290519081900360640190fd5b336000908152600660209081526040808320805460ff1916600190811790915560058054909101905580516f4d69677261746546726f6d4f6c64282960801b815281519081900360100181206001600160e01b031916818401528151808203600401815260249091019182905280516001600160a01b038616939192918291908401908083835b602083106108c55780518252601f1990920191602091820191016108a6565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610927576040519150601f19603f3d011682016040523d82523d6000602084013e61092c565b606091505b505090508061096c5760405162461bcd60e51b8152600401808060200182810382526024815260200180610c1d6024913960400191505060405180910390fd5b6004546005541415610983576000805460ff191690555b6003546005541415610a29576040516000906001600160a01b0384169047908381818185875af1925050503d80600081146109da576040519150601f19603f3d011682016040523d82523d6000602084013e6109df565b606091505b5050905080610a27576040805162461bcd60e51b815260206004820152600f60248201526e151c985b9cd9995c8819985a5b1959608a1b604482015290519081900360640190fd5b505b5050565b60005460ff16610a3c57600080fd5b3360009081526008602052604090206002015415610a5957600080fd5b60008111610aae576040805162461bcd60e51b815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b336000908152600860205260409020600501544311610b14576040805162461bcd60e51b815260206004820181905260248201527f72657175657374206973206c6f636b65642c206e6f7420617661696c61626c65604482015290519081900360640190fd5b33600081815260086020908152604080832060048101805460ff19169055600381019390935582546001600160a01b031916909317808355600283018590556007544301600590930192909255825184815292516001600160a01b03909216927f222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2929081900390910190a250565b60005460ff16610bb157600080fd5b6040805133815234602082015281517f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427929181900390910190a1565b60026020526000908152604090205481565b6001600160a01b031660009081526002602052604090205415159056fe556e61626c6520746f204d696772617465206e657720536d61727420636f6e7472616374a265627a7a7231582016efac2c97d498deea8150377c74dac62afa930d54f95b0b44526da32739b1dc64736f6c63430005100032666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

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
// Solidity: function isValidator(address v) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) IsValidator(opts *bind.CallOpts, v common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "isValidator", v)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address v) constant returns(bool)
func (_LockRedeem *LockRedeemSession) IsValidator(v common.Address) (bool, error) {
	return _LockRedeem.Contract.IsValidator(&_LockRedeem.CallOpts, v)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address v) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) IsValidator(v common.Address) (bool, error) {
	return _LockRedeem.Contract.IsValidator(&_LockRedeem.CallOpts, v)
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) IsredeemAvailable(opts *bind.CallOpts, recepient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "isredeemAvailable", recepient_)
	return *ret0, err
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeem *LockRedeemSession) IsredeemAvailable(recepient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.IsredeemAvailable(&_LockRedeem.CallOpts, recepient_)
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) IsredeemAvailable(recepient_ common.Address) (bool, error) {
	return _LockRedeem.Contract.IsredeemAvailable(&_LockRedeem.CallOpts, recepient_)
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) MigrationSignatures(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "migrationSignatures")
	return *ret0, err
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeem *LockRedeemSession) MigrationSignatures() (*big.Int, error) {
	return _LockRedeem.Contract.MigrationSignatures(&_LockRedeem.CallOpts)
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) MigrationSignatures() (*big.Int, error) {
	return _LockRedeem.Contract.MigrationSignatures(&_LockRedeem.CallOpts)
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) MigrationSigners(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeem.contract.Call(opts, out, "migrationSigners", arg0)
	return *ret0, err
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeem *LockRedeemSession) MigrationSigners(arg0 common.Address) (bool, error) {
	return _LockRedeem.Contract.MigrationSigners(&_LockRedeem.CallOpts, arg0)
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) MigrationSigners(arg0 common.Address) (bool, error) {
	return _LockRedeem.Contract.MigrationSigners(&_LockRedeem.CallOpts, arg0)
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

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint256)
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
// Solidity: function validators(address ) constant returns(uint256)
func (_LockRedeem *LockRedeemSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.Validators(&_LockRedeem.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint256)
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

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeem *LockRedeemTransactor) Migrate(opts *bind.TransactOpts, newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "migrate", newSmartContractAddress)
}

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeem *LockRedeemSession) Migrate(newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Migrate(&_LockRedeem.TransactOpts, newSmartContractAddress)
}

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeem *LockRedeemTransactorSession) Migrate(newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeem.Contract.Migrate(&_LockRedeem.TransactOpts, newSmartContractAddress)
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
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0x6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad36.
//
// Solidity: event AddValidator(address indexed _address)
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

// WatchAddValidator is a free log subscription operation binding the contract event 0x6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad36.
//
// Solidity: event AddValidator(address indexed _address)
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

// ParseAddValidator is a log parse operation binding the contract event 0x6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad36.
//
// Solidity: event AddValidator(address indexed _address)
func (_LockRedeem *LockRedeemFilterer) ParseAddValidator(log types.Log) (*LockRedeemAddValidator, error) {
	event := new(LockRedeemAddValidator)
	if err := _LockRedeem.contract.UnpackLog(event, "AddValidator", log); err != nil {
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
