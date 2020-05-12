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

// LockRedeemV2ABI is the input ABI used to generate the binding from.
const LockRedeemV2ABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_lock_period\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_old_contract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"noofValidatorsinold\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_requested\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"redeemFeeCharged\",\"type\":\"uint256\"}],\"name\":\"RedeemRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"NewSmartContractAddress\",\"type\":\"address\"}],\"name\":\"ValidatorMigrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasReturned\",\"type\":\"uint256\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"constant\":false,\"inputs\":[],\"name\":\"MigrateFromOld\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"collectUserFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMigrationCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getRedeemBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"getSignatureCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"hasValidatorSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recepient_\",\"type\":\"address\"}],\"name\":\"isredeemAvailable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newSmartContractAddress\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"migrationSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"migrationSigners\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"verifyRedeem\",\"outputs\":[{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"verifyValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// LockRedeemV2FuncSigs maps the 4-byte function signature to its string representation.
var LockRedeemV2FuncSigs = map[string]string{
	"587ab37e": "MigrateFromOld()",
	"7edd7ccd": "collectUserFee()",
	"cdaf4028": "getMigrationCount()",
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
	"bbe4fb7d": "verifyValidator()",
}

// LockRedeemV2Bin is the compiled bytecode used for deploying new contracts.
var LockRedeemV2Bin = "0x60806040526000805460ff19168155615a9860065561933c600755662386f26fc100006008556001600955612710600a55601281905560135534801561004457600080fd5b506040516114423803806114428339818101604052606081101561006757600080fd5b5080516020820151604090920151600f91909155600160036002830281900482016004819055920401600555601180546001600160a01b039093166001600160a01b03199093169290921790915560125561137b806100c76000396000f3fe60806040526004361061011f5760003560e01c806391e39868116100a0578063db006a7511610064578063db006a75146103ae578063e75f7515146103cb578063f83d08ba146103fe578063fa52c7d814610406578063facd743b1461044f5761011f565b806391e39868146102d2578063a04d04981461031e578063bbe4fb7d14610351578063cdaf402814610366578063ce5494bb1461037b5761011f565b8063587ab37e116100e7578063587ab37e146102255780635d593f8d1461023c5780636c7d13df146102515780637cacde3f146102845780637edd7ccd146102bd5761011f565b80632138c6b91461013e57806327882c3a14610185578063287cc96b146101ac57806331b6a6d1146101c157806345dfa415146101f4575b6012546013541461012f57600080fd5b6000805460ff19166001179055005b34801561014a57600080fd5b506101716004803603602081101561016157600080fd5b50356001600160a01b0316610482565b604080519115158252519081900360200190f35b34801561019157600080fd5b5061019a6104b3565b60408051918252519081900360200190f35b3480156101b857600080fd5b5061019a6104b9565b3480156101cd57600080fd5b50610171600480360360208110156101e457600080fd5b50356001600160a01b03166104bd565b34801561020057600080fd5b50610209610508565b604080516001600160a01b039092168252519081900360200190f35b34801561023157600080fd5b5061023a61050c565b005b34801561024857600080fd5b5061019a610537565b34801561025d57600080fd5b5061019a6004803603602081101561027457600080fd5b50356001600160a01b031661053d565b34801561029057600080fd5b5061023a600480360360408110156102a757600080fd5b50803590602001356001600160a01b031661056c565b3480156102c957600080fd5b5061023a6109dd565b3480156102de57600080fd5b50610305600480360360208110156102f557600080fd5b50356001600160a01b0316610ad8565b60408051600092830b90920b8252519081900360200190f35b34801561032a57600080fd5b506101716004803603602081101561034157600080fd5b50356001600160a01b0316610b91565b34801561035d57600080fd5b50610171610ba6565b34801561037257600080fd5b5061019a610bb6565b34801561038757600080fd5b5061023a6004803603602081101561039e57600080fd5b50356001600160a01b0316610bbc565b61023a600480360360208110156103c457600080fd5b5035610fc5565b3480156103d757600080fd5b5061019a600480360360208110156103ee57600080fd5b50356001600160a01b031661116f565b61023a61119e565b34801561041257600080fd5b506104396004803603602081101561042957600080fd5b50356001600160a01b03166111e9565b6040805160ff9092168252519081900360200190f35b34801561045b57600080fd5b506101716004803603602081101561047257600080fd5b50356001600160a01b03166111fe565b6000805460ff1661049257600080fd5b506001600160a01b0316600090815260106020526040902060040154431190565b600b5481565b4790565b6000805460ff166104cd57600080fd5b336000908152600260208181526040808420546001600160a01b03871685526010909252909220600190810154909260ff161c061492915050565b3090565b6011546001600160a01b0316331461052357600080fd5b6013805460010190556105353261121e565b565b60015481565b6000805460ff1661054d57600080fd5b506001600160a01b031660009081526010602052604090206003015490565b60005460ff1661057b57600080fd5b610584336111fe565b6105d5576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b60005a90506105e3336111fe565b610634576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e206c697374000000604482015290519081900360640190fd5b6001600160a01b03821660009081526010602052604090206004015443106106a3576040805162461bcd60e51b815260206004820152601f60248201527f72656465656d2072657175657374206973206e6f7420617661696c61626c6500604482015290519081900360640190fd5b6001600160a01b0382166000908152601060205260409020600201548314610712576040805162461bcd60e51b815260206004820152601a60248201527f72656465656d20616d6f756e7420697320646966666572656e74000000000000604482015290519081900360640190fd5b336000908152600260208181526040808420546001600160a01b0387168552601090925283206001015460ff9091161c0614610795576040805162461bcd60e51b815260206004820152601b60248201527f76616c696461746f722068617320616c726561647920766f7465640000000000604482015290519081900360640190fd5b336000908152600260209081526040808320546001600160a01b038616845260109092529091206001808201805460ff90941682901b909301909255600301805490910190819055600454116108bc576001600160a01b0380831660009081526010602052604080822080546002909101549151929316918381818185875af1925050503d8060008114610845576040519150601f19603f3d011682016040523d82523d6000602084013e61084a565b606091505b5050905080610893576040805162461bcd60e51b815260206004820152601060248201526f2a3930b739b332b9103330b4b632b21760811b604482015290519081900360640190fd5b506001600160a01b03821660009081526010602052604081206002810191909155436004909101555b60006007546006545a8403010190506000600954600a54023a83020190506000336001600160a01b03168260405180600001905060006040518083038185875af1925050503d806000811461092d576040519150601f19603f3d011682016040523d82523d6000602084013e610932565b606091505b50509050806109725760405162461bcd60e51b81526004018080602001828103825260218152602001806113266021913960400191505060405180910390fd5b6001600160a01b03851660008181526010602090815260409182902060050180548690039055815133815290810189905280820185905290517f975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d158309181900360600190a2505050505050565b60005460ff166109ec57600080fd5b60006109f733610ad8565b60000b13610a365760405162461bcd60e51b81526004018080602001828103825260248152602001806112d76024913960400191505060405180910390fd5b336000818152601060205260408082206005015490519192918381818185875af1925050503d8060008114610a87576040519150601f19603f3d011682016040523d82523d6000602084013e610a8c565b606091505b5050905080610ad5576040805162461bcd60e51b815260206004820152601060248201526f2a3930b739b332b9103330b4b632b21760811b604482015290519081900360640190fd5b50565b6000805460ff16610ae857600080fd5b6001600160a01b03821660009081526010602052604090206004015415801590610b2c57506001600160a01b03821660009081526010602052604090206004015443115b610b61576001600160a01b03821660009081526010602052604090206002015415610b58576000610b5c565b6000195b610b8b565b6001600160a01b038216600090815260106020526040902060020154610b88576001610b8b565b60025b92915050565b600c6020526000908152604090205460ff1681565b6000610bb1336111fe565b905090565b60135490565b610bc5336111fe565b610c16576040805162461bcd60e51b815260206004820152601d60248201527f76616c696461746f72206e6f742070726573656e7420696e204c697374000000604482015290519081900360640190fd5b336000908152600c602052604090205460ff1615610c7b576040805162461bcd60e51b815260206004820152601860248201527f56616c696461746f72205369676e656420616c72656164790000000000000000604482015290519081900360640190fd5b336000908152600c60209081526040808320805460ff1916600117905580516f4d69677261746546726f6d4f6c64282960801b815281519081900360100181206001600160e01b031916818401528151808203600401815260249091019182905280516001600160a01b038616939192918291908401908083835b60208310610d155780518252601f199092019160209182019101610cf6565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610d77576040519150601f19603f3d011682016040523d82523d6000602084013e610d7c565b606091505b5050905080610dbc5760405162461bcd60e51b81526004018080602001828103825260248152602001806112b36024913960400191505060405180910390fd5b600b805460010190556001600160a01b0382166000908152600d6020526040902054610e2e57600e80546001810182556000919091527fbb7b4a454dc3493923482f07822329ed19e8244eff582cc204f8554c3620c3fd0180546001600160a01b0319166001600160a01b0384161790555b6001600160a01b0382166000908152600d6020526040902080546001019055600554600b541415610e64576000805460ff191690555b600454600b541415610fc157600080805b600e54811015610f245782600d6000600e8481548110610e9157fe5b60009182526020808320909101546001600160a01b031683528201929092526040019020541115610f1c57600d6000600e8381548110610ecd57fe5b60009182526020808320909101546001600160a01b03168352820192909252604001902054600e80549194509082908110610f0457fe5b6000918252602090912001546001600160a01b031691505b600101610e75565b506040516000906001600160a01b0383169047908381818185875af1925050503d8060008114610f70576040519150601f19603f3d011682016040523d82523d6000602084013e610f75565b606091505b5050905080610fbd576040805162461bcd60e51b815260206004820152600f60248201526e151c985b9cd9995c8819985a5b1959608a1b604482015290519081900360640190fd5b5050505b5050565b60005460ff16610fd457600080fd5b610fdd33610482565b6110185760405162461bcd60e51b815260040180806020018281038252602b8152602001806112fb602b913960400191505060405180910390fd5b6000811161106d576040805162461bcd60e51b815260206004820152601e60248201527f616d6f756e742073686f756c6420626520626967676572207468616e20300000604482015290519081900360640190fd5b60085433600090815260106020526040902060050154340110156110d8576040805162461bcd60e51b815260206004820152601760248201527f52656465656d20666565206e6f742070726f7669646564000000000000000000604482015290519081900360640190fd5b3360008181526010602090815260408083206003810184905580546001600160a01b03191690941780855560028501869055600f544301600486015560058501805434019081905560019095019390935580518581529182019390935282516001600160a01b03909216927feee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc92918290030190a250565b6000805460ff1661117f57600080fd5b506001600160a01b031660009081526010602052604090206005015490565b60005460ff166111ad57600080fd5b6040805133815234602082015281517f625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427929181900390910190a1565b60026020526000908152604090205460ff1681565b6001600160a01b031660009081526002602052604090205460ff16151590565b600380546001810182557fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b0180546001600160a01b0319166001600160a01b0384169081179091559054600082815260026020526040808220805460ff191660ff9094169390931790925590517f6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad369190a25056fe556e61626c6520746f204d696772617465206e657720536d61727420636f6e747261637472657175657374207369676e696e67206973207374696c6c20696e2070726f677265737372656465656d20746f20746869732061646472657373206973206e6f7420617661696c61626c65207965745472616e73666572206261636b20746f2076616c696461746f72206661696c6564a265627a7a72315820cb99766138b344bb5320f03c8bfd454dbcd5c056d7315f305f11ccd15dd53a3f64736f6c63430005100032"

// DeployLockRedeemV2 deploys a new Ethereum contract, binding an instance of LockRedeemV2 to it.
func DeployLockRedeemV2(auth *bind.TransactOpts, backend bind.ContractBackend, _lock_period *big.Int, _old_contract common.Address, noofValidatorsinold *big.Int) (common.Address, *types.Transaction, *LockRedeemV2, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemV2ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(LockRedeemV2Bin), backend, _lock_period, _old_contract, noofValidatorsinold)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LockRedeemV2{LockRedeemV2Caller: LockRedeemV2Caller{contract: contract}, LockRedeemV2Transactor: LockRedeemV2Transactor{contract: contract}, LockRedeemV2Filterer: LockRedeemV2Filterer{contract: contract}}, nil
}

// LockRedeemV2 is an auto generated Go binding around an Ethereum contract.
type LockRedeemV2 struct {
	LockRedeemV2Caller     // Read-only binding to the contract
	LockRedeemV2Transactor // Write-only binding to the contract
	LockRedeemV2Filterer   // Log filterer for contract events
}

// LockRedeemV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type LockRedeemV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type LockRedeemV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LockRedeemV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LockRedeemV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LockRedeemV2Session struct {
	Contract     *LockRedeemV2     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LockRedeemV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LockRedeemV2CallerSession struct {
	Contract *LockRedeemV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// LockRedeemV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LockRedeemV2TransactorSession struct {
	Contract     *LockRedeemV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LockRedeemV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type LockRedeemV2Raw struct {
	Contract *LockRedeemV2 // Generic contract binding to access the raw methods on
}

// LockRedeemV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LockRedeemV2CallerRaw struct {
	Contract *LockRedeemV2Caller // Generic read-only contract binding to access the raw methods on
}

// LockRedeemV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LockRedeemV2TransactorRaw struct {
	Contract *LockRedeemV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewLockRedeemV2 creates a new instance of LockRedeemV2, bound to a specific deployed contract.
func NewLockRedeemV2(address common.Address, backend bind.ContractBackend) (*LockRedeemV2, error) {
	contract, err := bindLockRedeemV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2{LockRedeemV2Caller: LockRedeemV2Caller{contract: contract}, LockRedeemV2Transactor: LockRedeemV2Transactor{contract: contract}, LockRedeemV2Filterer: LockRedeemV2Filterer{contract: contract}}, nil
}

// NewLockRedeemV2Caller creates a new read-only instance of LockRedeemV2, bound to a specific deployed contract.
func NewLockRedeemV2Caller(address common.Address, caller bind.ContractCaller) (*LockRedeemV2Caller, error) {
	contract, err := bindLockRedeemV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2Caller{contract: contract}, nil
}

// NewLockRedeemV2Transactor creates a new write-only instance of LockRedeemV2, bound to a specific deployed contract.
func NewLockRedeemV2Transactor(address common.Address, transactor bind.ContractTransactor) (*LockRedeemV2Transactor, error) {
	contract, err := bindLockRedeemV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2Transactor{contract: contract}, nil
}

// NewLockRedeemV2Filterer creates a new log filterer instance of LockRedeemV2, bound to a specific deployed contract.
func NewLockRedeemV2Filterer(address common.Address, filterer bind.ContractFilterer) (*LockRedeemV2Filterer, error) {
	contract, err := bindLockRedeemV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2Filterer{contract: contract}, nil
}

// bindLockRedeemV2 binds a generic wrapper to an already deployed contract.
func bindLockRedeemV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LockRedeemV2ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemV2 *LockRedeemV2Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeemV2.Contract.LockRedeemV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemV2 *LockRedeemV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.LockRedeemV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemV2 *LockRedeemV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.LockRedeemV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LockRedeemV2 *LockRedeemV2CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _LockRedeemV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LockRedeemV2 *LockRedeemV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LockRedeemV2 *LockRedeemV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.contract.Transact(opts, method, params...)
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) GetMigrationCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "getMigrationCount")
	return *ret0, err
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) GetMigrationCount() (*big.Int, error) {
	return _LockRedeemV2.Contract.GetMigrationCount(&_LockRedeemV2.CallOpts)
}

// GetMigrationCount is a free data retrieval call binding the contract method 0xcdaf4028.
//
// Solidity: function getMigrationCount() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) GetMigrationCount() (*big.Int, error) {
	return _LockRedeemV2.Contract.GetMigrationCount(&_LockRedeemV2.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeemV2 *LockRedeemV2Caller) GetOLTEthAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "getOLTEthAddress")
	return *ret0, err
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeemV2 *LockRedeemV2Session) GetOLTEthAddress() (common.Address, error) {
	return _LockRedeemV2.Contract.GetOLTEthAddress(&_LockRedeemV2.CallOpts)
}

// GetOLTEthAddress is a free data retrieval call binding the contract method 0x45dfa415.
//
// Solidity: function getOLTEthAddress() constant returns(address)
func (_LockRedeemV2 *LockRedeemV2CallerSession) GetOLTEthAddress() (common.Address, error) {
	return _LockRedeemV2.Contract.GetOLTEthAddress(&_LockRedeemV2.CallOpts)
}

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) GetRedeemBalance(opts *bind.CallOpts, recipient_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "getRedeemBalance", recipient_)
	return *ret0, err
}

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) GetRedeemBalance(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeemV2.Contract.GetRedeemBalance(&_LockRedeemV2.CallOpts, recipient_)
}

// GetRedeemBalance is a free data retrieval call binding the contract method 0xe75f7515.
//
// Solidity: function getRedeemBalance(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) GetRedeemBalance(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeemV2.Contract.GetRedeemBalance(&_LockRedeemV2.CallOpts, recipient_)
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) GetSignatureCount(opts *bind.CallOpts, recipient_ common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "getSignatureCount", recipient_)
	return *ret0, err
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) GetSignatureCount(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeemV2.Contract.GetSignatureCount(&_LockRedeemV2.CallOpts, recipient_)
}

// GetSignatureCount is a free data retrieval call binding the contract method 0x6c7d13df.
//
// Solidity: function getSignatureCount(address recipient_) constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) GetSignatureCount(recipient_ common.Address) (*big.Int, error) {
	return _LockRedeemV2.Contract.GetSignatureCount(&_LockRedeemV2.CallOpts, recipient_)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) GetTotalEthBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "getTotalEthBalance")
	return *ret0, err
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeemV2.Contract.GetTotalEthBalance(&_LockRedeemV2.CallOpts)
}

// GetTotalEthBalance is a free data retrieval call binding the contract method 0x287cc96b.
//
// Solidity: function getTotalEthBalance() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) GetTotalEthBalance() (*big.Int, error) {
	return _LockRedeemV2.Contract.GetTotalEthBalance(&_LockRedeemV2.CallOpts)
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Caller) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "hasValidatorSigned", recipient_)
	return *ret0, err
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Session) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeemV2.Contract.HasValidatorSigned(&_LockRedeemV2.CallOpts, recipient_)
}

// HasValidatorSigned is a free data retrieval call binding the contract method 0x31b6a6d1.
//
// Solidity: function hasValidatorSigned(address recipient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2CallerSession) HasValidatorSigned(recipient_ common.Address) (bool, error) {
	return _LockRedeemV2.Contract.HasValidatorSigned(&_LockRedeemV2.CallOpts, recipient_)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Caller) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "isValidator", addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Session) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeemV2.Contract.IsValidator(&_LockRedeemV2.CallOpts, addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address addr) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2CallerSession) IsValidator(addr common.Address) (bool, error) {
	return _LockRedeemV2.Contract.IsValidator(&_LockRedeemV2.CallOpts, addr)
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Caller) IsredeemAvailable(opts *bind.CallOpts, recepient_ common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "isredeemAvailable", recepient_)
	return *ret0, err
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Session) IsredeemAvailable(recepient_ common.Address) (bool, error) {
	return _LockRedeemV2.Contract.IsredeemAvailable(&_LockRedeemV2.CallOpts, recepient_)
}

// IsredeemAvailable is a free data retrieval call binding the contract method 0x2138c6b9.
//
// Solidity: function isredeemAvailable(address recepient_) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2CallerSession) IsredeemAvailable(recepient_ common.Address) (bool, error) {
	return _LockRedeemV2.Contract.IsredeemAvailable(&_LockRedeemV2.CallOpts, recepient_)
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) MigrationSignatures(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "migrationSignatures")
	return *ret0, err
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) MigrationSignatures() (*big.Int, error) {
	return _LockRedeemV2.Contract.MigrationSignatures(&_LockRedeemV2.CallOpts)
}

// MigrationSignatures is a free data retrieval call binding the contract method 0x27882c3a.
//
// Solidity: function migrationSignatures() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) MigrationSignatures() (*big.Int, error) {
	return _LockRedeemV2.Contract.MigrationSignatures(&_LockRedeemV2.CallOpts)
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Caller) MigrationSigners(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "migrationSigners", arg0)
	return *ret0, err
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Session) MigrationSigners(arg0 common.Address) (bool, error) {
	return _LockRedeemV2.Contract.MigrationSigners(&_LockRedeemV2.CallOpts, arg0)
}

// MigrationSigners is a free data retrieval call binding the contract method 0xa04d0498.
//
// Solidity: function migrationSigners(address ) constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2CallerSession) MigrationSigners(arg0 common.Address) (bool, error) {
	return _LockRedeemV2.Contract.MigrationSigners(&_LockRedeemV2.CallOpts, arg0)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Caller) NumValidators(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "numValidators")
	return *ret0, err
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2Session) NumValidators() (*big.Int, error) {
	return _LockRedeemV2.Contract.NumValidators(&_LockRedeemV2.CallOpts)
}

// NumValidators is a free data retrieval call binding the contract method 0x5d593f8d.
//
// Solidity: function numValidators() constant returns(uint256)
func (_LockRedeemV2 *LockRedeemV2CallerSession) NumValidators() (*big.Int, error) {
	return _LockRedeemV2.Contract.NumValidators(&_LockRedeemV2.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeemV2 *LockRedeemV2Caller) Validators(opts *bind.CallOpts, arg0 common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeemV2 *LockRedeemV2Session) Validators(arg0 common.Address) (uint8, error) {
	return _LockRedeemV2.Contract.Validators(&_LockRedeemV2.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) constant returns(uint8)
func (_LockRedeemV2 *LockRedeemV2CallerSession) Validators(arg0 common.Address) (uint8, error) {
	return _LockRedeemV2.Contract.Validators(&_LockRedeemV2.CallOpts, arg0)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeemV2 *LockRedeemV2Caller) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (int8, error) {
	var (
		ret0 = new(int8)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "verifyRedeem", recipient_)
	return *ret0, err
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeemV2 *LockRedeemV2Session) VerifyRedeem(recipient_ common.Address) (int8, error) {
	return _LockRedeemV2.Contract.VerifyRedeem(&_LockRedeemV2.CallOpts, recipient_)
}

// VerifyRedeem is a free data retrieval call binding the contract method 0x91e39868.
//
// Solidity: function verifyRedeem(address recipient_) constant returns(int8)
func (_LockRedeemV2 *LockRedeemV2CallerSession) VerifyRedeem(recipient_ common.Address) (int8, error) {
	return _LockRedeemV2.Contract.VerifyRedeem(&_LockRedeemV2.CallOpts, recipient_)
}

// VerifyValidator is a free data retrieval call binding the contract method 0xbbe4fb7d.
//
// Solidity: function verifyValidator() constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Caller) VerifyValidator(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _LockRedeemV2.contract.Call(opts, out, "verifyValidator")
	return *ret0, err
}

// VerifyValidator is a free data retrieval call binding the contract method 0xbbe4fb7d.
//
// Solidity: function verifyValidator() constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2Session) VerifyValidator() (bool, error) {
	return _LockRedeemV2.Contract.VerifyValidator(&_LockRedeemV2.CallOpts)
}

// VerifyValidator is a free data retrieval call binding the contract method 0xbbe4fb7d.
//
// Solidity: function verifyValidator() constant returns(bool)
func (_LockRedeemV2 *LockRedeemV2CallerSession) VerifyValidator() (bool, error) {
	return _LockRedeemV2.Contract.VerifyValidator(&_LockRedeemV2.CallOpts)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) MigrateFromOld(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "MigrateFromOld")
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemV2 *LockRedeemV2Session) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.MigrateFromOld(&_LockRedeemV2.TransactOpts)
}

// MigrateFromOld is a paid mutator transaction binding the contract method 0x587ab37e.
//
// Solidity: function MigrateFromOld() returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) MigrateFromOld() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.MigrateFromOld(&_LockRedeemV2.TransactOpts)
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) CollectUserFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "collectUserFee")
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeemV2 *LockRedeemV2Session) CollectUserFee() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.CollectUserFee(&_LockRedeemV2.TransactOpts)
}

// CollectUserFee is a paid mutator transaction binding the contract method 0x7edd7ccd.
//
// Solidity: function collectUserFee() returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) CollectUserFee() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.CollectUserFee(&_LockRedeemV2.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeemV2 *LockRedeemV2Session) Lock() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Lock(&_LockRedeemV2.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) Lock() (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Lock(&_LockRedeemV2.TransactOpts)
}

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) Migrate(opts *bind.TransactOpts, newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "migrate", newSmartContractAddress)
}

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeemV2 *LockRedeemV2Session) Migrate(newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Migrate(&_LockRedeemV2.TransactOpts, newSmartContractAddress)
}

// Migrate is a paid mutator transaction binding the contract method 0xce5494bb.
//
// Solidity: function migrate(address newSmartContractAddress) returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) Migrate(newSmartContractAddress common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Migrate(&_LockRedeemV2.TransactOpts, newSmartContractAddress)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) Redeem(opts *bind.TransactOpts, amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "redeem", amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemV2 *LockRedeemV2Session) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Redeem(&_LockRedeemV2.TransactOpts, amount_)
}

// Redeem is a paid mutator transaction binding the contract method 0xdb006a75.
//
// Solidity: function redeem(uint256 amount_) returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) Redeem(amount_ *big.Int) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Redeem(&_LockRedeemV2.TransactOpts, amount_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemV2 *LockRedeemV2Transactor) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.contract.Transact(opts, "sign", amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemV2 *LockRedeemV2Session) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Sign(&_LockRedeemV2.TransactOpts, amount_, recipient_)
}

// Sign is a paid mutator transaction binding the contract method 0x7cacde3f.
//
// Solidity: function sign(uint256 amount_, address recipient_) returns()
func (_LockRedeemV2 *LockRedeemV2TransactorSession) Sign(amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return _LockRedeemV2.Contract.Sign(&_LockRedeemV2.TransactOpts, amount_, recipient_)
}

// LockRedeemV2AddValidatorIterator is returned from FilterAddValidator and is used to iterate over the raw logs and unpacked data for AddValidator events raised by the LockRedeemV2 contract.
type LockRedeemV2AddValidatorIterator struct {
	Event *LockRedeemV2AddValidator // Event containing the contract specifics and raw log

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
func (it *LockRedeemV2AddValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemV2AddValidator)
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
		it.Event = new(LockRedeemV2AddValidator)
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
func (it *LockRedeemV2AddValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemV2AddValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemV2AddValidator represents a AddValidator event raised by the LockRedeemV2 contract.
type LockRedeemV2AddValidator struct {
	Address common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAddValidator is a free log retrieval operation binding the contract event 0x6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad36.
//
// Solidity: event AddValidator(address indexed _address)
func (_LockRedeemV2 *LockRedeemV2Filterer) FilterAddValidator(opts *bind.FilterOpts, _address []common.Address) (*LockRedeemV2AddValidatorIterator, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemV2.contract.FilterLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2AddValidatorIterator{contract: _LockRedeemV2.contract, event: "AddValidator", logs: logs, sub: sub}, nil
}

// WatchAddValidator is a free log subscription operation binding the contract event 0x6a7a7b9e5967ba1cf76c3d7d5a9b98e96f11754855b04564fada97b94741ad36.
//
// Solidity: event AddValidator(address indexed _address)
func (_LockRedeemV2 *LockRedeemV2Filterer) WatchAddValidator(opts *bind.WatchOpts, sink chan<- *LockRedeemV2AddValidator, _address []common.Address) (event.Subscription, error) {

	var _addressRule []interface{}
	for _, _addressItem := range _address {
		_addressRule = append(_addressRule, _addressItem)
	}

	logs, sub, err := _LockRedeemV2.contract.WatchLogs(opts, "AddValidator", _addressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemV2AddValidator)
				if err := _LockRedeemV2.contract.UnpackLog(event, "AddValidator", log); err != nil {
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
func (_LockRedeemV2 *LockRedeemV2Filterer) ParseAddValidator(log types.Log) (*LockRedeemV2AddValidator, error) {
	event := new(LockRedeemV2AddValidator)
	if err := _LockRedeemV2.contract.UnpackLog(event, "AddValidator", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemV2LockIterator is returned from FilterLock and is used to iterate over the raw logs and unpacked data for Lock events raised by the LockRedeemV2 contract.
type LockRedeemV2LockIterator struct {
	Event *LockRedeemV2Lock // Event containing the contract specifics and raw log

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
func (it *LockRedeemV2LockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemV2Lock)
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
		it.Event = new(LockRedeemV2Lock)
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
func (it *LockRedeemV2LockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemV2LockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemV2Lock represents a Lock event raised by the LockRedeemV2 contract.
type LockRedeemV2Lock struct {
	Sender         common.Address
	AmountReceived *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLock is a free log retrieval operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_LockRedeemV2 *LockRedeemV2Filterer) FilterLock(opts *bind.FilterOpts) (*LockRedeemV2LockIterator, error) {

	logs, sub, err := _LockRedeemV2.contract.FilterLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2LockIterator{contract: _LockRedeemV2.contract, event: "Lock", logs: logs, sub: sub}, nil
}

// WatchLock is a free log subscription operation binding the contract event 0x625fed9875dada8643f2418b838ae0bc78d9a148a18eee4ee1979ff0f3f5d427.
//
// Solidity: event Lock(address sender, uint256 amount_received)
func (_LockRedeemV2 *LockRedeemV2Filterer) WatchLock(opts *bind.WatchOpts, sink chan<- *LockRedeemV2Lock) (event.Subscription, error) {

	logs, sub, err := _LockRedeemV2.contract.WatchLogs(opts, "Lock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemV2Lock)
				if err := _LockRedeemV2.contract.UnpackLog(event, "Lock", log); err != nil {
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
func (_LockRedeemV2 *LockRedeemV2Filterer) ParseLock(log types.Log) (*LockRedeemV2Lock, error) {
	event := new(LockRedeemV2Lock)
	if err := _LockRedeemV2.contract.UnpackLog(event, "Lock", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemV2RedeemRequestIterator is returned from FilterRedeemRequest and is used to iterate over the raw logs and unpacked data for RedeemRequest events raised by the LockRedeemV2 contract.
type LockRedeemV2RedeemRequestIterator struct {
	Event *LockRedeemV2RedeemRequest // Event containing the contract specifics and raw log

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
func (it *LockRedeemV2RedeemRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemV2RedeemRequest)
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
		it.Event = new(LockRedeemV2RedeemRequest)
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
func (it *LockRedeemV2RedeemRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemV2RedeemRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemV2RedeemRequest represents a RedeemRequest event raised by the LockRedeemV2 contract.
type LockRedeemV2RedeemRequest struct {
	Recepient        common.Address
	AmountRequested  *big.Int
	RedeemFeeCharged *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRedeemRequest is a free log retrieval operation binding the contract event 0xeee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested, uint256 redeemFeeCharged)
func (_LockRedeemV2 *LockRedeemV2Filterer) FilterRedeemRequest(opts *bind.FilterOpts, recepient []common.Address) (*LockRedeemV2RedeemRequestIterator, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemV2.contract.FilterLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2RedeemRequestIterator{contract: _LockRedeemV2.contract, event: "RedeemRequest", logs: logs, sub: sub}, nil
}

// WatchRedeemRequest is a free log subscription operation binding the contract event 0xeee07ebdabc7ab1dc20be39b715e23aa8a85c6a8ae3c16f8334dace8d76683dc.
//
// Solidity: event RedeemRequest(address indexed recepient, uint256 amount_requested, uint256 redeemFeeCharged)
func (_LockRedeemV2 *LockRedeemV2Filterer) WatchRedeemRequest(opts *bind.WatchOpts, sink chan<- *LockRedeemV2RedeemRequest, recepient []common.Address) (event.Subscription, error) {

	var recepientRule []interface{}
	for _, recepientItem := range recepient {
		recepientRule = append(recepientRule, recepientItem)
	}

	logs, sub, err := _LockRedeemV2.contract.WatchLogs(opts, "RedeemRequest", recepientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemV2RedeemRequest)
				if err := _LockRedeemV2.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
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
func (_LockRedeemV2 *LockRedeemV2Filterer) ParseRedeemRequest(log types.Log) (*LockRedeemV2RedeemRequest, error) {
	event := new(LockRedeemV2RedeemRequest)
	if err := _LockRedeemV2.contract.UnpackLog(event, "RedeemRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemV2ValidatorMigratedIterator is returned from FilterValidatorMigrated and is used to iterate over the raw logs and unpacked data for ValidatorMigrated events raised by the LockRedeemV2 contract.
type LockRedeemV2ValidatorMigratedIterator struct {
	Event *LockRedeemV2ValidatorMigrated // Event containing the contract specifics and raw log

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
func (it *LockRedeemV2ValidatorMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemV2ValidatorMigrated)
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
		it.Event = new(LockRedeemV2ValidatorMigrated)
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
func (it *LockRedeemV2ValidatorMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemV2ValidatorMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemV2ValidatorMigrated represents a ValidatorMigrated event raised by the LockRedeemV2 contract.
type LockRedeemV2ValidatorMigrated struct {
	Validator               common.Address
	NewSmartContractAddress common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterValidatorMigrated is a free log retrieval operation binding the contract event 0x077478953a7559f9e01b2ceeb429ce87333fb7fc0ec16eb5eb9128463e30fa92.
//
// Solidity: event ValidatorMigrated(address validator, address NewSmartContractAddress)
func (_LockRedeemV2 *LockRedeemV2Filterer) FilterValidatorMigrated(opts *bind.FilterOpts) (*LockRedeemV2ValidatorMigratedIterator, error) {

	logs, sub, err := _LockRedeemV2.contract.FilterLogs(opts, "ValidatorMigrated")
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2ValidatorMigratedIterator{contract: _LockRedeemV2.contract, event: "ValidatorMigrated", logs: logs, sub: sub}, nil
}

// WatchValidatorMigrated is a free log subscription operation binding the contract event 0x077478953a7559f9e01b2ceeb429ce87333fb7fc0ec16eb5eb9128463e30fa92.
//
// Solidity: event ValidatorMigrated(address validator, address NewSmartContractAddress)
func (_LockRedeemV2 *LockRedeemV2Filterer) WatchValidatorMigrated(opts *bind.WatchOpts, sink chan<- *LockRedeemV2ValidatorMigrated) (event.Subscription, error) {

	logs, sub, err := _LockRedeemV2.contract.WatchLogs(opts, "ValidatorMigrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemV2ValidatorMigrated)
				if err := _LockRedeemV2.contract.UnpackLog(event, "ValidatorMigrated", log); err != nil {
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
func (_LockRedeemV2 *LockRedeemV2Filterer) ParseValidatorMigrated(log types.Log) (*LockRedeemV2ValidatorMigrated, error) {
	event := new(LockRedeemV2ValidatorMigrated)
	if err := _LockRedeemV2.contract.UnpackLog(event, "ValidatorMigrated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// LockRedeemV2ValidatorSignedRedeemIterator is returned from FilterValidatorSignedRedeem and is used to iterate over the raw logs and unpacked data for ValidatorSignedRedeem events raised by the LockRedeemV2 contract.
type LockRedeemV2ValidatorSignedRedeemIterator struct {
	Event *LockRedeemV2ValidatorSignedRedeem // Event containing the contract specifics and raw log

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
func (it *LockRedeemV2ValidatorSignedRedeemIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LockRedeemV2ValidatorSignedRedeem)
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
		it.Event = new(LockRedeemV2ValidatorSignedRedeem)
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
func (it *LockRedeemV2ValidatorSignedRedeemIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LockRedeemV2ValidatorSignedRedeemIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LockRedeemV2ValidatorSignedRedeem represents a ValidatorSignedRedeem event raised by the LockRedeemV2 contract.
type LockRedeemV2ValidatorSignedRedeem struct {
	Recipient         common.Address
	ValidatorAddresss common.Address
	Amount            *big.Int
	GasReturned       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterValidatorSignedRedeem is a free log retrieval operation binding the contract event 0x975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d15830.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount, uint256 gasReturned)
func (_LockRedeemV2 *LockRedeemV2Filterer) FilterValidatorSignedRedeem(opts *bind.FilterOpts, recipient []common.Address) (*LockRedeemV2ValidatorSignedRedeemIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeemV2.contract.FilterLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return &LockRedeemV2ValidatorSignedRedeemIterator{contract: _LockRedeemV2.contract, event: "ValidatorSignedRedeem", logs: logs, sub: sub}, nil
}

// WatchValidatorSignedRedeem is a free log subscription operation binding the contract event 0x975a8b0f36f1204c7939f566cea0503ea32284a2768a7f98ede91960b6d15830.
//
// Solidity: event ValidatorSignedRedeem(address indexed recipient, address validator_addresss, uint256 amount, uint256 gasReturned)
func (_LockRedeemV2 *LockRedeemV2Filterer) WatchValidatorSignedRedeem(opts *bind.WatchOpts, sink chan<- *LockRedeemV2ValidatorSignedRedeem, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _LockRedeemV2.contract.WatchLogs(opts, "ValidatorSignedRedeem", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LockRedeemV2ValidatorSignedRedeem)
				if err := _LockRedeemV2.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
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
func (_LockRedeemV2 *LockRedeemV2Filterer) ParseValidatorSignedRedeem(log types.Log) (*LockRedeemV2ValidatorSignedRedeem, error) {
	event := new(LockRedeemV2ValidatorSignedRedeem)
	if err := _LockRedeemV2.contract.UnpackLog(event, "ValidatorSignedRedeem", log); err != nil {
		return nil, err
	}
	return event, nil
}
