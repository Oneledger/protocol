package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
)

type RedeemStatus int8

const (
	NewRedeem       RedeemStatus = -1
	Ongoing         RedeemStatus = 0
	Success         RedeemStatus = 1
	Expired         RedeemStatus = 2
	ErrorConnecting RedeemStatus = -2
)

func (r RedeemStatus) String() string {
	switch r {
	case -1:
		return "NewRedeem"
	case 0:
		return "OnGoing"
	case 1:
		return "Success"
	case 2:
		return "Expired"
	case 3:
		return "Error connecting to ethereum"
	}
	return "Unknown Type / Check smart contract implementation"
}

type CheckFinalityStatus int8

const (
	TxBlockNotFound        CheckFinalityStatus = 0x01
	BlockHashFailed        CheckFinalityStatus = 0x02
	UnabletoGetHeader      CheckFinalityStatus = 0x03
	NotEnoughConfirmations CheckFinalityStatus = 0x04
	ReciptNotFound         CheckFinalityStatus = 0x05
	TransactionNotMined    CheckFinalityStatus = 0x06
	TXSuccess              CheckFinalityStatus = 0x07
)

type RedeemRequest struct {
	Amount *big.Int
}

type LockErcRequest struct {
	Receiver    common.Address
	TokenAmount *big.Int
}

type LockRequest struct {
	Amount *big.Int
}

type RedeemErcRequest struct {
	Amount       *big.Int
	TokenAddress common.Address
}

type ContractType int8

const (
	ETH ContractType = 0x00
	ERC ContractType = 0x01
)

type Contract interface {
	IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error)
	VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (int8, error)
	HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error)
	Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error)
	IsRedeemAvailable(opts *bind.CallOpts, recipient common.Address) (bool, error)
}

var _ Contract = ETHLRContract{}

type ETHLRContract struct {
	contract contract.LockRedeem
}

func (E ETHLRContract) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return E.contract.IsValidator(opts, addr)
}
func (E ETHLRContract) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (int8, error) {
	return E.contract.VerifyRedeem(opts, recipient_)
}

func (E ETHLRContract) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.HasValidatorSigned(opts, recipient_)
}

func (E ETHLRContract) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return E.contract.Sign(opts, amount_, recipient_)
}

func (E ETHLRContract) IsRedeemAvailable(opts *bind.CallOpts, recipient common.Address) (bool, error) {
	return E.contract.IsredeemAvailable(opts, recipient)
}

func GetETHContract(contract contract.LockRedeem) *ETHLRContract {
	return &ETHLRContract{contract}
}

var _ Contract = &ERC20LRContract{}

type ERC20LRContract struct {
	contract contract.LockRedeemERC
}

func (E ERC20LRContract) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return E.contract.IsValidator(opts, addr)
}

func (E ERC20LRContract) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (int8, error) {
	//return E.contract.VerifyRedeem(opts, recipient_)
	panic("Implement verify redeem for status updates")
}

func (E ERC20LRContract) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.HasValidatorSigned(opts, recipient_)
}

func (E ERC20LRContract) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error) {
	return E.contract.Sign(opts, amount_, recipient_)
}

func (E ERC20LRContract) IsRedeemAvailable(opts *bind.CallOpts, recipient common.Address) (bool, error) {
	panic("IsRedeemAvailable not implemented")
}
func GetERCContract(contract contract.LockRedeemERC) *ERC20LRContract {
	return &ERC20LRContract{contract}
}
