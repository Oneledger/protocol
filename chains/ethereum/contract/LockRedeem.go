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

// LockRedeemABI is the input ABI used to generate the binding from.
const LockRedeemABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_lock_period\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_old_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"noofValidatorsinold\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"redeemFeeCharged\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"NewSmartContractAddress\",\"type\":\"address\"}],\"name\":\"ValidatorMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasReturned\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":true,\"inputs\":[],\"name\":\"ActiveStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"MigrateFromOld\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"collectUserFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getRedeemBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getSignatureCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"hasValidatorSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recepient_\",\"type\":\"address\"}],\"name\":\"isredeemAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newSmartContractAddress\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"migrationSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"migrationSigners\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"verifyRedeem\",\"outputs\":[{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LockRedeemFuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemFuncSigs = map[string]string{
	"f95d2f57": "ActiveStatus()",
	"587ab37e": "MigrateFromOld()",
	"7edd7ccd": "collectUserFee()",
	"45dfa415": "getOLTEthAddress()",
	"e75f7515": "getRedeemBalance(address)",
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
var LockRedeemBin = "0x60806040526000805460ff19168155615a9860065561933c600755662386f26fc100006008556001600955612710600a5560118190556012553480156200004557600080fd5b506040516200159738038062001597833981810160405260808110156200006b57600080fd5b81019080805160405193929190846401000000008211156200008c57600080fd5b908301906020820185811115620000a257600080fd5b8251866020820283011164010000000082111715620000c057600080fd5b82525081516020918201928201910280838360005b83811015620000ef578181015183820152602001620000d5565b50505050919091016040908152602083015190830151606090930151909450919250620001199050565b60005b8451811015620001bd5760008582815181106200013557fe5b6020908102919091018101516001600160a01b0381166000908152600290925260409091205490915060ff16156200019f5760405162461bcd60e51b815260040180806020018281038252602f81526020018062001568602f913960400191505060405180910390fd5b620001b3816001600160e01b036200021316565b506001016200011c565b50600f839055835160039060020204600101600455835160039004600101600555601080546001600160a01b0319166001600160a01b0384161790556003600282020460010160118190555050505050620002a7565b600380546001810182557fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b0180546001600160a01b0319166001600160a01b0384169081179091559054600082815260026020526040808220805460ff191660ff9094169390931790925590517f6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad369190a250565b6112b180620002b76000396000f3fe6080604052600436106101145760003560e01c80637edd7ccd116100a0578063e75f751511610064578063e75f751514610396578063f83d08ba146103c9578063f95d2f57146103d1578063fa52c7d8146103e6578063facd743b1461042f57610114565b80637edd7ccd146102b257806391e39868146102c7578063a04d049814610313578063ce5494bb14610346578063db006a751461037957610114565b806345dfa415116100e757806345dfa415146101e9578063587ab37e1461021a5780635d593f8d146102315780636c7d13df146102465780637cacde3f1461027957610114565b80632138c6b91461013357806327882c3a1461017a578063287cc96b146101a157806331b6a6d1146101b6575b6011546012541461012457600080fd5b6000805460ff19166001179055005b34801561013f57600080fd5b506101666004803603602081101561015657600080fd5b50356001600160a01b0316610462565b604080519115158252519081900360200190f35b34801561018657600080fd5b5061018f610493565b60408051918252519081900360200190f35b3480156101ad57600080fd5b5061018f610499565b3480156101c257600080fd5b50610166600480360360208110156101d957600080fd5b50356001600160a01b031661049d565b3480156101f557600080fd5b506101fe6104e8565b604080516001600160a01b039092168252519081900360200190f35b34801561022657600080fd5b5061022f6104ec565b005b34801561023d57600080fd5b5061018f61050e565b34801561025257600080fd5b5061018f6004803603602081101561026957600080fd5b50356001600160a01b0316610514565b34801561028557600080fd5b5061022f6004803603604081101561029c57600080fd5b50803590602001356001600160a01b0316610543565b3480156102be57600080fd5b5061022f6109b4565b3480156102d357600080fd5b506102fa600480360360208110156102ea57600080fd5b50356001600160a01b0316610aaf565b60408051600092830b90920b8252519081900360200190f35b34801561031f57600080fd5b506101666004803603602081101561033657600080fd5b50356001600160a01b0316610b68565b34801561035257600080fd5b5061022f6004803603602081101561036957600080fd5b50356001600160a01b0316610b7d565b61022f6004803603602081101561038f57600080fd5b5035610f86565b3480156103a257600080fd5b5061018f600480360360208110156103b957600080fd5b50356001600160a01b0316611130565b61022f61115f565b3480156103dd57600080fd5b506101666111aa565b3480156103f257600080fd5b506104196004803603602081101561040957600080fd5b50356001600160a01b03166111b3565b6040805160ff9092168252519081900360200190f35b34801561043b57600080fd5b506101666004803603602081101561045257600080fd5b50356001600160a01b03166111c8565b6000805460ff1661047257600080fd5b506001600160a01b0316600090815260136020526040902060040154431190565b600b5481565b4790565b6000805460ff166104ad57600080fd5b336000908152600260208181526040808420546001600160a01b03871685526013909252909220600190810154909260ff161c061492915050565b3090565b6010546001600160a01b0316331461050357600080fd5b601280546001019055565b60015481565b6000805460ff1661052457600080fd5b506001600160a01b031660009081526013602052604090206003015490565b60005460ff1661055257600080fd5b61055b336111c8565b6105ac576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b60005a90506105ba336111c8565b61060b576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e206c697374000000604482015290519081900360640190fd5b6001600160a01b038216600090815260136020526040902060040154431061067a576040805162461bcd60e51b815260206004820152601f60248201527f72656465656d2072657175657374206973206e6f7420617661696c61626c6500604482015290519081900360640190fd5b6001600160a01b03821660009081526013602052604090206002015483146106e9576040805162461bcd60e51b815260206004820152601a60248201527f72656465656d20616d6f756e7420697320646966666572656e74000000000000604482015290519081900360640190fd5b336000908152600260208181526040808420546001600160a01b0387168552601390925283206001015460ff9091161c061461076c576040805162461bcd60e51b815260206004820152601b60248201527f76616c696461746f722068617320616c726561647920766f7465640000000000604482015290519081900360640190fd5b336000908152600260209081526040808320546001600160a01b038616845260139092529091206001808201805460ff90941682901b90930190925560030180549091019081905560045411610893576001600160a01b0380831660009081526013602052604080822080546002909101549151929316918381818185875af1925050503d806000811461081c576040519150601f19603f3d011682016040523d82523d6000602084013e610821565b606091505b505090508061086a576040805162461bcd60e51b815260206004820152601060248201526f2a3930b739b332b9103330b4b632b21760811b604482015290519081900360640190fd5b506001600160a01b03821660009081526013602052604081206002810191909155436004909101555b60006007546006545a8403010190506000600954600a54023a83020190506000336001600160a01b03168260405180600001905060006040518083038185875af1925050503d8060008114610904576040519150601f19603f3d011682016040523d82523d6000602084013e610909565b606091505b50509050806109495760405162461bcd60e51b815260040180806020018281038252602181526020018061125c6021913960400191505060405180910390fd5b6001600160a01b03851660008181526013602090815260409182902060050180548690039055815133815290810189905280820185905290517f975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d158309181900360600190a2505050505050565b60005460ff166109c357600080fd5b60006109ce33610aaf565b60000b13610a0d5760405162461bcd60e51b815260040180806020018281038252602481526020018061120d6024913960400191505060405180910390fd5b336000818152601360205260408082206005015490519192918381818185875af1925050503d8060008114610a5e576040519150601f19603f3d011682016040523d82523d6000602084013e610a63565b606091505b5050905080610aac576040805162461bcd60e51b815260206004820152601060248201526f2a3930b739b332b9103330b4b632b21760811b604482015290519081900360640190fd5b50565b6000805460ff16610abf57600080fd5b6001600160a01b03821660009081526013602052604090206004015415801590610b0357506001600160a01b03821660009081526013602052604090206004015443115b610b38576001600160a01b03821660009081526013602052604090206002015415610b2f576000610b33565b6000195b610b62565b6001600160a01b038216600090815260136020526040902060020154610b5f576001610b62565b60025b92915050565b600c6020526000908152604090205460ff1681565b610b86336111c8565b610bd7576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b336000908152600c602052604090205460ff1615610c3c576040805162461bcd60e51b815260206004820152601860248201527f56616c696461746f72205369676e656420616c72656164790000000000000000604482015290519081900360640190fd5b336000908152600c60209081526040808320805460ff1916600117905580516f4d69677261746546726f6d4f6c64282960801b815281519081900360100181206001600160e01b031916818401528151808203600401815260249091019182905280516001600160a01b038616939192918291908401908083835b60208310610cd65780518252601f199092019160209182019101610cb7565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610d38576040519150601f19603f3d011682016040523d82523d6000602084013e610d3d565b606091505b5050905080610d7d5760405162461bcd60e51b81526004018080602001828103825260248152602001806111e96024913960400191505060405180910390fd5b600b805460010190556001600160a01b0382166000908152600d6020526040902054610def57600e80546001810182556000919091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd0180546001600160a01b0319166001600160a01b0384161790555b6001600160a01b0382166000908152600d6020526040902080546001019055600554600b541415610e25576000805460ff191690555b600454600b541415610f8257600080805b600e54811015610ee55782600d6000600e8481548110610e5257fe5b60009182526020808320909101546001600160a01b031683528201929092526040019020541115610edd57600d6000600e8381548110610e8e57fe5b60009182526020808320909101546001600160a01b03168352820192909252604001902054600e80549194509082908110610ec557fe5b6000918252602090912001546001600160a01b031691505b600101610e36565b506040516000906001600160a01b0383169047908381818185875af1925050503d8060008114610f31576040519150601f19603f3d011682016040523d82523d6000602084013e610f36565b606091505b5050905080610f7e576040805162461bcd60e51b815260206004820152600f60248201526e151c985b9cd9995c8819985a5b1959608a1b604482015290519081900360640190fd5b5050505b5050565b60005460ff16610f9557600080fd5b610f9e33610462565b610fd95760405162461bcd60e51b815260040180806020018281038252602b815260200180611231602b913960400191505060405180910390fd5b6000811161102e576040805162461bcd60e51b815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b6008543360009081526013602052604090206005015434011015611099576040805162461bcd60e51b815260206004820152601760248201527f52656465656d20666565206e6f742070726f7669646564000000000000000000604482015290519081900360640190fd5b3360008181526013602090815260408083206003810184905580546001600160a01b03191690941780855560028501869055600f544301600486015560058501805434019081905560019095019390935580518581529182019390935282516001600160a01b03909216927feee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc92918290030190a250565b6000805460ff1661114057600080fd5b506001600160a01b031660009081526013602052604090206005015490565b60005460ff1661116e57600080fd5b6040805133815234602082015281517f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427929181900390910190a1565b60005460ff1690565b60026020526000908152604090205460ff1681565b6001600160a01b031660009081526002602052604090205460ff1615159056fe556e61626c6520746f204d696772617465206e657720536d61727420636f6e747261637472657175657374207369676e696e67206973207374696c6c20696e2070726f677265737372656465656d20746f20746869732061646472657373206973206e6f7420617661696c61626c65207965745472616e73666572206261636b20746f2076616c696461746f72206661696c6564a265627a7a7231582053ceb9b73a73a681943b0de1e427d4b6e8e46028087d4ff272f9e0f7f71db6ea64736f6c63430005100032666f756e64206e6f6e2d756e697175652076616c696461746f7220696e20696e697469616c56616c696461746f7273"

// DeployLockRedeem deploys a new Ethereum contract, binding an instance of LockRedeem to it.
func DeployLockRedeem(auth *bind.TransactOpts, backend bind.ContractBackend, initialValidators []common.Address, _lock_period *big.Int, _old_contract common.Address, noofValidatorsinold *big.Int) (common.Address, *types.Transaction, *LockRedeem, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockRedeemBin), backend, initialValidators, _lock_period, _old_contract, noofValidatorsinold)
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
	results := make([]interface{}, 1)
	results[0] = result
	return _LockRedeem.Contract.LockRedeemCaller.contract.Call(opts, &results, method, params...)
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
	results := make([]interface{}, 1)
	results[0] = result
	return _LockRedeem.Contract.contract.Call(opts, &results, method, params...)
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

// ActiveStatus is a free data retrieval call binding the contract method 0xf95d2f57.
//
// Solidity: function ActiveStatus() constant returns(bool)
func (_LockRedeem *LockRedeemCaller) ActiveStatus(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "ActiveStatus")
	return *ret0, err
}

// ActiveStatus is a free data retrieval call binding the contract method 0xf95d2f57.
//
// Solidity: function ActiveStatus() constant returns(bool)
func (_LockRedeem *LockRedeemSession) ActiveStatus() (bool, error) {
	return _LockRedeem.Contract.ActiveStatus(&_LockRedeem.CallOpts)
}

// ActiveStatus is a free data retrieval call binding the contract method 0xf95d2f57.
//
// Solidity: function ActiveStatus() constant returns(bool)
func (_LockRedeem *LockRedeemCallerSession) ActiveStatus() (bool, error) {
	return _LockRedeem.Contract.ActiveStatus(&_LockRedeem.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeem *LockRedeemCaller) GetOLTEthAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "getOLTEthAddress")
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

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) GetRedeemBalance(opts *bind.CallOpts, recipient_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "getRedeemBalance", recipient_)
	return *ret0, err
}

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemSession) GetRedeemBalance(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.GetRedeemBalance(&_LockRedeem.CallOpts, recipient_)
}

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemCallerSession) GetRedeemBalance(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeem.Contract.GetRedeemBalance(&_LockRedeem.CallOpts, recipient_)
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeem *LockRedeemCaller) GetSignatureCount(opts *bind.CallOpts, recipient_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "getSignatureCount", recipient_)
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "getTotalEthBalance")
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "hasValidatorSigned", recipient_)
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "isValidator", addr)
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

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeem *LockRedeemCaller) IsredeemAvailable(opts *bind.CallOpts, recepient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "isredeemAvailable", recepient_)
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "migrationSignatures")
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "migrationSigners", arg0)
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
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "numValidators")
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
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeem *LockRedeemCaller) Validators(opts *bind.CallOpts, arg0 common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeem *LockRedeemSession) Validators(arg0 common.Address) (uint8, error) {
	return _LockRedeem.Contract.Validators(&_LockRedeem.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeem *LockRedeemCallerSession) Validators(arg0 common.Address) (uint8, error) {
	return _LockRedeem.Contract.Validators(&_LockRedeem.CallOpts, arg0)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeem *LockRedeemCaller) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (int8, error) {
	var (
		ret0 = new(int8)
	)
	out := ret0
	results := make([]interface{}, 1)
	results[0] = out
	err := _LockRedeem.contract.Call(opts, &results, "verifyRedeem", recipient_)
	return *ret0, err
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeem *LockRedeemSession) VerifyRedeem(recipient_ common.Address) (int8, error) {
	return _LockRedeem.Contract.VerifyRedeem(&_LockRedeem.CallOpts, recipient_)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeem *LockRedeemCallerSession) VerifyRedeem(recipient_ common.Address) (int8, error) {
	return _LockRedeem.Contract.VerifyRedeem(&_LockRedeem.CallOpts, recipient_)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeem *LockRedeemTransactor) MigrateFromOld(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "MigrateFromOld")
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeem *LockRedeemSession) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeem.Contract.MigrateFromOld(&_LockRedeem.TransactOpts)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeem *LockRedeemTransactorSession) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeem.Contract.MigrateFromOld(&_LockRedeem.TransactOpts)
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeem *LockRedeemTransactor) CollectUserFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeem.contract.Transact(opts, "collectUserFee")
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeem *LockRedeemSession) CollectUserFee() (*types.Transaction, error) {
	return _LockRedeem.Contract.CollectUserFee(&_LockRedeem.TransactOpts)
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeem *LockRedeemTransactorSession) CollectUserFee() (*types.Transaction, error) {
	return _LockRedeem.Contract.CollectUserFee(&_LockRedeem.TransactOpts)
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
	Recepient        common.Address
	AmountRequested  *big.Int
	RedeemFeeCharged *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRedeemRequest is a free log retrieval operation binding the contract event 0xeee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested, uint256 redeemFeeCharged)
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

// WatchRedeemRequest is a free log subscription operation binding the contract event 0xeee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested, uint256 redeemFeeCharged)
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

// ParseRedeemRequest is a log parse operation binding the contract event 0xeee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested, uint256 redeemFeeCharged)
func (_LockRedeem *LockRedeemFilterer) ParseRedeemRequest(log types.Log) (*LockRedeemRedeemRequest, error) {
	event := new(LockRedeemRedeemRequest)
	if err := _LockRedeem.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemValidatorMigratedIterator is returned from FilterValidatorMigrated and is used to iterate over the raw logs and unpacked data for ValidatorMigrated events raised by the LockRedeem contract.
type LockRedeemValidatorMigratedIterator struct {
	Event *LockRedeemValidatorMigrated // Event containing the contract specifics and raw log

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
func (it *LockRedeemValidatorMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemValidatorMigrated)
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
		it.Event = new(LockRedeemValidatorMigrated)
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
func (it *LockRedeemValidatorMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemValidatorMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemValidatorMigrated represents a ValidatorMigrated event raised by the LockRedeem contract.
type LockRedeemValidatorMigrated struct {
	Validator               common.Address
	NewSmartContractAddress common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterValidatorMigrated is a free log retrieval operation binding the contract event 0x077478953a7559f9e01b2ceeb429ce87333fb7fc0ec16eb5eb9128463e30fa92.
//
// Solidity: event ValidatorMigrated(address validator, address NewSmartContractAddress)
func (_LockRedeem *LockRedeemFilterer) FilterValidatorMigrated(opts *bind.FilterOpts) (*LockRedeemValidatorMigratedIterator, error) {

	logs, sub, err := _LockRedeem.contract.FilterLogs(opts, "ValidatorMigrated")
	if err != nil {
		return nil, err
	}
	return &LockRedeemValidatorMigratedIterator{contract: _LockRedeem.contract, event: "ValidatorMigrated", logs: logs, sub: sub}, nil
}

// WatchValidatorMigrated is a free log subscription operation binding the contract event 0x077478953a7559f9e01b2ceeb429ce87333fb7fc0ec16eb5eb9128463e30fa92.
//
// Solidity: event ValidatorMigrated(address validator, address NewSmartContractAddress)
func (_LockRedeem *LockRedeemFilterer) WatchValidatorMigrated(opts *bind.WatchOpts, sink chan<- *LockRedeemValidatorMigrated) (event.Subscription, error) {

	logs, sub, err := _LockRedeem.contract.WatchLogs(opts, "ValidatorMigrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemValidatorMigrated)
				if err := _LockRedeem.contract.UnpackLog(event, "ValidatorMigrated", log); err != nil {
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

// ParseValidatorMigrated is a log parse operation binding the contract event 0x077478953a7559f9e01b2ceeb429ce87333fb7fc0ec16eb5eb9128463e30fa92.
//
// Solidity: event ValidatorMigrated(address validator, address NewSmartContractAddress)
func (_LockRedeem *LockRedeemFilterer) ParseValidatorMigrated(log types.Log) (*LockRedeemValidatorMigrated, error) {
	event := new(LockRedeemValidatorMigrated)
	if err := _LockRedeem.contract.UnpackLog(event, "ValidatorMigrated", log); err != nil {
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
	GasReturned       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d15830.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount, uint256 gasReturned)
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

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d15830.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount, uint256 gasReturned)
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

// ParseValidatorSignedRedeem is a log parse operation binding the contract event 0x975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d15830.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount, uint256 gasReturned)
func (_LockRedeem *LockRedeemFilterer) ParseValidatorSignedRedeem(log types.Log) (*LockRedeemValidatorSignedRedeem, error) {
	event := new(LockRedeemValidatorSignedRedeem)
	if err := _LockRedeem.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}
