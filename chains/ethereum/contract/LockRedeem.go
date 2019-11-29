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

// ContractABI is the input ABI used to generate the binding from.
const ContractABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"}]"

// ContractBin is the compiled bytecode used for deploying new contracts.
var ContractBin = "0x60806040526170806002553480156200001757600080fd5b5060405162001dda38038062001dda833981810160405260208110156200003d57600080fd5b81019080805160405193929190846401000000008211156200005e57600080fd5b838201915060208201858111156200007557600080fd5b82518660208202830111640100000000821117156200009357600080fd5b8083526020830192505050908051906020019060200280838360005b83811015620000cc578082015181840152602081019050620000af565b5050505090500160405250505060008151101562000136576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602d81526020018062001d7e602d913960400191505060405180910390fd5b60008090505b8151811015620002195760008282815181106200015557fe5b602002602001015190506000600960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414620001f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f81526020018062001dab602f913960400191505060405180910390fd5b6200020a816200023260201b60201c565b5080806001019150506200013c565b506200022b436200031760201b60201c565b5062000623565b6032600960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600160008082825401925050819055508073ffffffffffffffffffffffffffffffffffffffff167fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040518082815260200191505060405180910390a250565b6000600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054136200036457600080fd5b806003819055507febad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e3356003546040518082815260200191505060405180910390a160008090505b6007805490508110156200045f57600060078281548110620003c857fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905062000406816200023260201b60201c565b600460008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000808201600090555050508080600101915050620003aa565b5060076000620004709190620005d8565b60008090505b6008805490508110156200052b576000600882815481106200049457fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050620004d2816200053f60201b60201c565b600560008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600080820160009055505050808060010191505062000476565b50600860006200053c9190620005d8565b50565b600960008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009055600160008082825403925050819055508073ffffffffffffffffffffffffffffffffffffffff167f6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e560405160405180910390a250565b5080546000825590600052602060002090810190620005f89190620005fb565b50565b6200062091905b808211156200061c57600081600090555060010162000602565b5090565b90565b61174b80620006336000396000f3fe6080604052600436106100f35760003560e01c8063628277331161008a578063e0e887d011610059578063e0e887d01461044c578063f83d08ba14610487578063fa52c7d814610491578063facd743b146104f6576100f3565b806362827733146103265780637cacde3f14610351578063bfb9e9f5146103ac578063db006a7514610411576100f3565b8063287cc96b116100c6578063287cc96b14610228578063383ea59a1461025357806345dfa415146102a45780635d593f8d146102fb576100f3565b80630d00753a146100f85780630d8f6b5b1461015d5780630e7d275d14610188578063101a8538146101d7575b600080fd5b34801561010457600080fd5b506101476004803603602081101561011b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061055f565b6040518082815260200191505060405180910390f35b34801561016957600080fd5b5061017261057d565b6040518082815260200191505060405180910390f35b34801561019457600080fd5b506101c1600480360360208110156101ab57600080fd5b8101908080359060200190929190505050610583565b6040518082815260200191505060405180910390f35b3480156101e357600080fd5b50610226600480360360208110156101fa57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105a1565b005b34801561023457600080fd5b5061023d6106d9565b6040518082815260200191505060405180910390f35b34801561025f57600080fd5b506102a26004803603602081101561027657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506106f8565b005b3480156102b057600080fd5b506102b96108a6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561030757600080fd5b506103106108ae565b6040518082815260200191505060405180910390f35b34801561033257600080fd5b5061033b6108b4565b6040518082815260200191505060405180910390f35b34801561035d57600080fd5b506103aa6004803603604081101561037457600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506108ba565b005b3480156103b857600080fd5b506103fb600480360360208110156103cf57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610e13565b6040518082815260200191505060405180910390f35b34801561041d57600080fd5b5061044a6004803603602081101561043457600080fd5b8101908080359060200190929190505050610e31565b005b34801561045857600080fd5b506104856004803603602081101561046f57600080fd5b8101908080359060200190929190505050611250565b005b61048f611423565b005b34801561049d57600080fd5b506104e0600480360360208110156104b457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611507565b6040518082815260200191505060405180910390f35b34801561050257600080fd5b506105456004803603602081101561051957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061151f565b604051808215151515815260200191505060405180910390f35b60056020528060005260406000206000915090508060000154905081565b60035481565b60066020528060005260406000206000915090508060000154905081565b6000600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054136105ed57600080fd5b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002090508060010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156106d5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260308152602001806116bb6030913960400191505060405180910390fd5b5050565b60003073ffffffffffffffffffffffffffffffffffffffff1631905090565b6000600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541361074457600080fd5b6000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002090508060010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161561082c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602c8152602001806116eb602c913960400191505060405180910390fd5b60018160010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600181600001600082825401925050819055506108a28261156a565b5050565b600030905090565b60005481565b60015481565b6108c33361151f565b610935576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601d8152602001807f76616c696461746f72206e6f742070726573656e7420696e206c69737400000081525060200191505060405180910390fd5b81600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060020154146109ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260198152602001807f72656465656d20616d6f756e7420436f6d70726f6d697365640000000000000081525060200191505060405180910390fd5b600a60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615610a8357600080fd5b6001600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506001600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060030160008282540192505081905550600154600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003015410610d8d57600a60008273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc600a60008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201549081150290604051600060405180830381858888f19350505050158015610ca1573d6000803e3d6000fd5b506000600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201819055506001600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160006101000a81548160ff02191690831515021790555043600a60008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600501819055505b8073ffffffffffffffffffffffffffffffffffffffff167f3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db3384604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a25050565b60046020528060005260406000206000915090508060000154905081565b6000600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206002015414610e8057600080fd5b60008111610ef6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f616d6f756e742073686f756c6420626520626967676572207468616e2030000081525060200191505060405180910390fd5b43600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206005015410610fad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f72657175657374206973206c6f636b65642c206e6f7420617661696c61626c6581525060200191505060405180910390fd5b60001515600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060040160009054906101000a905050506000600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206003018190555033600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201819055506002544301600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060050181905550600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2600a60003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600201546040518082815260200191505060405180910390a250565b6000600960003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541361129c57600080fd5b60005481106112f6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260418152602001806116506041913960600191505060405180910390fd5b60006006600083815260200190815260200160002090508060010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16156113b2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602a815260200180611691602a913960400191505060405180910390fd5b60018160010160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600181600001600082825401925050819055505050565b600034101561149a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f4d7573742070617920612062616c616e6365206d6f7265207468616e2030000081525060200191505060405180910390fd5b7f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d4273334604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a1565b60096020528060005260406000206000915090505481565b600080600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054139050919050565b6032600960008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600160008082825401925050819055508073ffffffffffffffffffffffffffffffffffffffff167fb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b600960008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546040518082815260200191505060405180910390a25056fe4e6577207468726573686f6c647320286d29206d757374206265206c657373207468616e20746865206e756d626572206f662076616c696461746f727320286e2973656e6465722068617320616c726561647920766f74656420666f7220746869732070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f20616464207468697320746f2070726f706f73616c73656e6465722068617320616c726561647920766f74656420746f2061646420746869732061646472657373a265627a7a72315820c2d6331b1babb91556272aef46eca1897f964ffaac0d437df7709dd558b561ee64736f6c634300050b0032696e73756666696369656e742076616c696461746f72732070617373656420746f20636f6e7374727563746f72666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, initialValidators []common.Address) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ContractBin), backend, initialValidators)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractCaller) AddValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "addValidatorProposals", arg0)
	return *ret0, err
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.AddValidatorProposals(&_Contract.CallOpts, arg0)
}

// AddValidatorProposals is a free data retrieval call binding the contract method 0xbfb9e9f5.
//
// Solidity: function addValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractCallerSession) AddValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.AddValidatorProposals(&_Contract.CallOpts, arg0)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_Contract *ContractCaller) EpochBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "epochBlockHeight")
	return *ret0, err
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_Contract *ContractSession) EpochBlockHeight() (*big.Int, error) {
	return _Contract.Contract.EpochBlockHeight(&_Contract.CallOpts)
}

// EpochBlockHeight is a free data retrieval call binding the contract method 0x0d8f6b5b.
//
// Solidity: function epochBlockHeight() constant returns(uint256)
func (_Contract *ContractCallerSession) EpochBlockHeight() (*big.Int, error) {
	return _Contract.Contract.EpochBlockHeight(&_Contract.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_Contract *ContractCaller) GetOLTEthAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "getOLTEthAddress")
	return *ret0, err
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_Contract *ContractSession) GetOLTEthAddress() (common.Address, error) {
	return _Contract.Contract.GetOLTEthAddress(&_Contract.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_Contract *ContractCallerSession) GetOLTEthAddress() (common.Address, error) {
	return _Contract.Contract.GetOLTEthAddress(&_Contract.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_Contract *ContractCaller) GetTotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "getTotalEthBalance")
	return *ret0, err
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_Contract *ContractSession) GetTotalEthBalance() (*big.Int, error) {
	return _Contract.Contract.GetTotalEthBalance(&_Contract.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_Contract *ContractCallerSession) GetTotalEthBalance() (*big.Int, error) {
	return _Contract.Contract.GetTotalEthBalance(&_Contract.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_Contract *ContractCaller) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "isValidator", addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_Contract *ContractSession) IsValidator(addr common.Address) (bool, error) {
	return _Contract.Contract.IsValidator(&_Contract.CallOpts, addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_Contract *ContractCallerSession) IsValidator(addr common.Address) (bool, error) {
	return _Contract.Contract.IsValidator(&_Contract.CallOpts, addr)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_Contract *ContractCaller) NewThresholdProposals(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "newThresholdProposals", arg0)
	return *ret0, err
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_Contract *ContractSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _Contract.Contract.NewThresholdProposals(&_Contract.CallOpts, arg0)
}

// NewThresholdProposals is a free data retrieval call binding the contract method 0x0e7d275d.
//
// Solidity: function newThresholdProposals(uint256 ) constant returns(uint256 voteCount)
func (_Contract *ContractCallerSession) NewThresholdProposals(arg0 *big.Int) (*big.Int, error) {
	return _Contract.Contract.NewThresholdProposals(&_Contract.CallOpts, arg0)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_Contract *ContractCaller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_Contract *ContractSession) NumValidators() (*big.Int, error) {
	return _Contract.Contract.NumValidators(&_Contract.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_Contract *ContractCallerSession) NumValidators() (*big.Int, error) {
	return _Contract.Contract.NumValidators(&_Contract.CallOpts)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_Contract *ContractCaller) ProposeRemoveValidator(opts *bind.CallOpts, v common.Address) error {
	var ()
	out := &[]interface{}{}
	err := _Contract.contract.Call(opts, out, "proposeRemoveValidator", v)
	return err
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_Contract *ContractSession) ProposeRemoveValidator(v common.Address) error {
	return _Contract.Contract.ProposeRemoveValidator(&_Contract.CallOpts, v)
}

// ProposeRemoveValidator is a free data retrieval call binding the contract method 0x101a8538.
//
// Solidity: function proposeRemoveValidator(address v) constant returns()
func (_Contract *ContractCallerSession) ProposeRemoveValidator(v common.Address) error {
	return _Contract.Contract.ProposeRemoveValidator(&_Contract.CallOpts, v)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractCaller) RemoveValidatorProposals(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "removeValidatorProposals", arg0)
	return *ret0, err
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.RemoveValidatorProposals(&_Contract.CallOpts, arg0)
}

// RemoveValidatorProposals is a free data retrieval call binding the contract method 0x0d00753a.
//
// Solidity: function removeValidatorProposals(address ) constant returns(uint256 voteCount)
func (_Contract *ContractCallerSession) RemoveValidatorProposals(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.RemoveValidatorProposals(&_Contract.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_Contract *ContractCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_Contract *ContractSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.Validators(&_Contract.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(int256)
func (_Contract *ContractCallerSession) Validators(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.Validators(&_Contract.CallOpts, arg0)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_Contract *ContractCaller) VotingThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "votingThreshold")
	return *ret0, err
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_Contract *ContractSession) VotingThreshold() (*big.Int, error) {
	return _Contract.Contract.VotingThreshold(&_Contract.CallOpts)
}

// VotingThreshold is a free data retrieval call binding the contract method 0x62827733.
//
// Solidity: function votingThreshold() constant returns(uint256)
func (_Contract *ContractCallerSession) VotingThreshold() (*big.Int, error) {
	return _Contract.Contract.VotingThreshold(&_Contract.CallOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Contract *ContractTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Contract *ContractSession) Lock() (*types.Transaction, error) {
	return _Contract.Contract.Lock(&_Contract.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Contract *ContractTransactorSession) Lock() (*types.Transaction, error) {
	return _Contract.Contract.Lock(&_Contract.TransactOpts)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_Contract *ContractTransactor) ProposeAddValidator(opts *bind.TransactOpts, v common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "proposeAddValidator", v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_Contract *ContractSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _Contract.Contract.ProposeAddValidator(&_Contract.TransactOpts, v)
}

// ProposeAddValidator is a paid mutator transaction binding the contract method 0x383ea59a.
//
// Solidity: function proposeAddValidator(address v) returns()
func (_Contract *ContractTransactorSession) ProposeAddValidator(v common.Address) (*types.Transaction, error) {
	return _Contract.Contract.ProposeAddValidator(&_Contract.TransactOpts, v)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_Contract *ContractTransactor) ProposeNewThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "proposeNewThreshold", threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_Contract *ContractSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ProposeNewThreshold(&_Contract.TransactOpts, threshold)
}

// ProposeNewThreshold is a paid mutator transaction binding the contract method 0xe0e887d0.
//
// Solidity: function proposeNewThreshold(uint256 threshold) returns()
func (_Contract *ContractTransactorSession) ProposeNewThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ProposeNewThreshold(&_Contract.TransactOpts, threshold)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_Contract *ContractTransactor) Redeem(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "redeem", amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_Contract *ContractSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Redeem(&_Contract.TransactOpts, amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_Contract *ContractTransactorSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.Redeem(&_Contract.TransactOpts, amount_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_Contract *ContractTransactor) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "sign", amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_Contract *ContractSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Sign(&_Contract.TransactOpts, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_Contract *ContractTransactorSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Sign(&_Contract.TransactOpts, amount_, recipient_)
}

// ContractAddValidatorIterator is returned from FilterAddValidator and is used to iterate over the raw logs and unpacked data for AddValidator events raised by the Contract contract.
type ContractAddValidatorIterator struct {
	Event *ContractAddValidator // Event containing the contract specifics and raw log

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
func (it *ContractAddValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractAddValidator)
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
		it.Event = new(ContractAddValidator)
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
func (it *ContractAddValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractAddValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractAddValidator represents a AddValidator event raised by the Contract contract.
type ContractAddValidator struct {
	Address common.Address
	Power   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_Contract *ContractFilterer) FilterAddValidator(opts *bind.FilterOpts, _address []common.Address) (*ContractAddValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &ContractAddValidatorIterator{contract: _Contract.contract, event: "AddValidator", logs: logs, sub: sub}, nil
}

// WatchAddValidator is a free log subscription operation binding the contract event 0xb2076c69a79e1dfb01d613dcc63b7c42ae1962daf11d4f2151352135133f824b.
//
// Solidity: event AddValidator(address indexed _address, int256 _power)
func (_Contract *ContractFilterer) WatchAddValidator(opts *bind.WatchOpts, sink chan<- *ContractAddValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractAddValidator)
				if err := _Contract.contract.UnpackLog(event, "AddValidator", log); err != nil {
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
func (_Contract *ContractFilterer) ParseAddValidator(log types.Log) (*ContractAddValidator, error) {
	event := new(ContractAddValidator)
	if err := _Contract.contract.UnpackLog(event, "AddValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractDeleteValidatorIterator is returned from FilterDeleteValidator and is used to iterate over the raw logs and unpacked data for DeleteValidator events raised by the Contract contract.
type ContractDeleteValidatorIterator struct {
	Event *ContractDeleteValidator // Event containing the contract specifics and raw log

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
func (it *ContractDeleteValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDeleteValidator)
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
		it.Event = new(ContractDeleteValidator)
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
func (it *ContractDeleteValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDeleteValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDeleteValidator represents a DeleteValidator event raised by the Contract contract.
type ContractDeleteValidator struct {
	Address common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDeleteValidator is a free log retrieval operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_Contract *ContractFilterer) FilterDeleteValidator(opts *bind.FilterOpts, _address []common.Address) (*ContractDeleteValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &ContractDeleteValidatorIterator{contract: _Contract.contract, event: "DeleteValidator", logs: logs, sub: sub}, nil
}

// WatchDeleteValidator is a free log subscription operation binding the contract event 0x6d70afad774d81e8c32f930c6412789502b16ccf0a20f21679b249bdfac060e5.
//
// Solidity: event DeleteValidator(address indexed _address)
func (_Contract *ContractFilterer) WatchDeleteValidator(opts *bind.WatchOpts, sink chan<- *ContractDeleteValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "DeleteValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDeleteValidator)
				if err := _Contract.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
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
func (_Contract *ContractFilterer) ParseDeleteValidator(log types.Log) (*ContractDeleteValidator, error) {
	event := new(ContractDeleteValidator)
	if err := _Contract.contract.UnpackLog(event, "DeleteValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractLockIterator is returned from FilterLock and is used to iterate over the raw logs and unpacked data for Lock events raised by the Contract contract.
type ContractLockIterator struct {
	Event *ContractLock // Event containing the contract specifics and raw log

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
func (it *ContractLockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractLock)
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
		it.Event = new(ContractLock)
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
func (it *ContractLockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractLockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractLock represents a Lock event raised by the Contract contract.
type ContractLock struct {
	Sender         common.Address
	AmountReceived *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_Contract *ContractFilterer) FilterLock(opts *bind.FilterOpts) (*ContractLockIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return &ContractLockIterator{contract: _Contract.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_Contract *ContractFilterer) WatchLock(opts *bind.WatchOpts, sink chan<- *ContractLock) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractLock)
				if err := _Contract.contract.UnpackLog(event, "Lock", log); err != nil {
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
func (_Contract *ContractFilterer) ParseLock(log types.Log) (*ContractLock, error) {
	event := new(ContractLock)
	if err := _Contract.contract.UnpackLog(event, "Lock", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the Contract contract.
type ContractNewEpochIterator struct {
	Event *ContractNewEpoch // Event containing the contract specifics and raw log

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
func (it *ContractNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractNewEpoch)
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
		it.Event = new(ContractNewEpoch)
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
func (it *ContractNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractNewEpoch represents a NewEpoch event raised by the Contract contract.
type ContractNewEpoch struct {
	EpochHeight *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_Contract *ContractFilterer) FilterNewEpoch(opts *bind.FilterOpts) (*ContractNewEpochIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return &ContractNewEpochIterator{contract: _Contract.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xebad8099c467528a56c98b63c8d476d251cf1ffb4c75db94b4d23fa2b6a1e335.
//
// Solidity: event NewEpoch(uint256 epochHeight)
func (_Contract *ContractFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *ContractNewEpoch) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "NewEpoch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractNewEpoch)
				if err := _Contract.contract.UnpackLog(event, "NewEpoch", log); err != nil {
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
func (_Contract *ContractFilterer) ParseNewEpoch(log types.Log) (*ContractNewEpoch, error) {
	event := new(ContractNewEpoch)
	if err := _Contract.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractNewThresholdIterator is returned from FilterNewThreshold and is used to iterate over the raw logs and unpacked data for NewThreshold events raised by the Contract contract.
type ContractNewThresholdIterator struct {
	Event *ContractNewThreshold // Event containing the contract specifics and raw log

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
func (it *ContractNewThresholdIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractNewThreshold)
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
		it.Event = new(ContractNewThreshold)
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
func (it *ContractNewThresholdIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractNewThresholdIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractNewThreshold represents a NewThreshold event raised by the Contract contract.
type ContractNewThreshold struct {
	PrevThreshold *big.Int
	NewThreshold  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewThreshold is a free log retrieval operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_Contract *ContractFilterer) FilterNewThreshold(opts *bind.FilterOpts) (*ContractNewThresholdIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return &ContractNewThresholdIterator{contract: _Contract.contract, event: "NewThreshold", logs: logs, sub: sub}, nil
}

// WatchNewThreshold is a free log subscription operation binding the contract event 0x7a5c0f01d83576763cde96136a1c8a8c1c05ff95d8a184db483894a9b69b8b3a.
//
// Solidity: event NewThreshold(uint256 _prevThreshold, uint256 _newThreshold)
func (_Contract *ContractFilterer) WatchNewThreshold(opts *bind.WatchOpts, sink chan<- *ContractNewThreshold) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "NewThreshold")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractNewThreshold)
				if err := _Contract.contract.UnpackLog(event, "NewThreshold", log); err != nil {
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
func (_Contract *ContractFilterer) ParseNewThreshold(log types.Log) (*ContractNewThreshold, error) {
	event := new(ContractNewThreshold)
	if err := _Contract.contract.UnpackLog(event, "NewThreshold", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractRedeemRequestIterator is returned from FilterRedeemRequest and is used to iterate over the raw logs and unpacked data for RedeemRequest events raised by the Contract contract.
type ContractRedeemRequestIterator struct {
	Event *ContractRedeemRequest // Event containing the contract specifics and raw log

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
func (it *ContractRedeemRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractRedeemRequest)
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
		it.Event = new(ContractRedeemRequest)
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
func (it *ContractRedeemRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractRedeemRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractRedeemRequest represents a RedeemRequest event raised by the Contract contract.
type ContractRedeemRequest struct {
	Recepient       common.Address
	AmountRequested *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedeemRequest is a free log retrieval operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_Contract *ContractFilterer) FilterRedeemRequest(opts *bind.FilterOpts, recepient []common.Address) (*ContractRedeemRequestIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return &ContractRedeemRequestIterator{contract: _Contract.contract, event: "RedeemRequest", logs: logs, sub: sub}, nil
}

// WatchRedeemRequest is a free log subscription operation binding the contract event 0x222dc200773fe9b45015bf792e8fee37d651e3590c215806a5042404b6d741d2.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested)
func (_Contract *ContractFilterer) WatchRedeemRequest(opts *bind.WatchOpts, sink chan<- *ContractRedeemRequest, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractRedeemRequest)
				if err := _Contract.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
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
func (_Contract *ContractFilterer) ParseRedeemRequest(log types.Log) (*ContractRedeemRequest, error) {
	event := new(ContractRedeemRequest)
	if err := _Contract.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractValidatorSignedRedeemIterator is returned from FilterValidatorSignedRedeem and is used to iterate over the raw logs and unpacked data for ValidatorSignedRedeem events raised by the Contract contract.
type ContractValidatorSignedRedeemIterator struct {
	Event *ContractValidatorSignedRedeem // Event containing the contract specifics and raw log

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
func (it *ContractValidatorSignedRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractValidatorSignedRedeem)
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
		it.Event = new(ContractValidatorSignedRedeem)
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
func (it *ContractValidatorSignedRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractValidatorSignedRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractValidatorSignedRedeem represents a ValidatorSignedRedeem event raised by the Contract contract.
type ContractValidatorSignedRedeem struct {
	Recipient         common.Address
	ValidatorAddresss common.Address
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_Contract *ContractFilterer) FilterValidatorSignedRedeem(opts *bind.FilterOpts, recipient []common.Address) (*ContractValidatorSignedRedeemIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return &ContractValidatorSignedRedeemIterator{contract: _Contract.contract, event: "ValidatorSignedRedeem", logs: logs, sub: sub}, nil
}

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x3b76df4bf55914fbcbc8b02f6773984cc346db1e6aef40410dcee0f94c6a05db.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount)
func (_Contract *ContractFilterer) WatchValidatorSignedRedeem(opts *bind.WatchOpts, sink chan<- *ContractValidatorSignedRedeem, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractValidatorSignedRedeem)
				if err := _Contract.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
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
func (_Contract *ContractFilterer) ParseValidatorSignedRedeem(log types.Log) (*ContractValidatorSignedRedeem, error) {
	event := new(ContractValidatorSignedRedeem)
	if err := _Contract.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}
