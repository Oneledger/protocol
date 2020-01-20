package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
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
	Amount *big.Int
	TokenAddress common.Address
}




type ContractType int8
const (
	ETH ContractType = 0x00
	ERC ContractType = 0x01
)

type Contract interface {
		IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error)
		VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (bool, error)
	    HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error)
		Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error)
	}


var _ Contract = ETHLRContract{}

type ETHLRContract struct {
	contract contract.LockRedeem
}

func (E ETHLRContract) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return E.contract.IsValidator(opts,addr)
}
func (E ETHLRContract) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.VerifyRedeem(opts,recipient_)
}
func (E ETHLRContract) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.HasValidatorSigned(opts,recipient_)
}
func (E ETHLRContract) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error){
	return E.contract.Sign(opts,amount_,recipient_)
}

func GetETHContract(contract contract.LockRedeem) *ETHLRContract {
	return &ETHLRContract{contract}
}


var _ Contract = &ERC20LRContract{}

type ERC20LRContract struct {
	contract contract.LockRedeemERC
}

func (E ERC20LRContract) IsValidator(opts *bind.CallOpts, addr common.Address) (bool, error) {
	return E.contract.IsValidator(opts,addr)
}
func (E ERC20LRContract) VerifyRedeem(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.VerifyRedeem(opts,recipient_)
}
func (E ERC20LRContract) HasValidatorSigned(opts *bind.CallOpts, recipient_ common.Address) (bool, error) {
	return E.contract.HasValidatorSigned(opts,recipient_)
}
func (E ERC20LRContract) Sign(opts *bind.TransactOpts, amount_ *big.Int, recipient_ common.Address) (*types.Transaction, error){
	return E.contract.Sign(opts,amount_,recipient_)
}

func GetERCContract(contract contract.LockRedeemERC) *ERC20LRContract {
	return &ERC20LRContract{contract}
}