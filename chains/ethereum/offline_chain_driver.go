package ethereum

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
)

type OfflineETHChainDriver struct {
}

// Implements ChainDriver interface
var _ ChainDriver = &OfflineETHChainDriver{}

// ParseErc20Lock parses the the ERC20 Lock transaction and returns LockErcRequest .
// LockErcRequest contains the parameters passed to the lock function function of LockRedeem contract.
func ParseErc20Lock(erc20list []ERC20Token, rawEthTx []byte) (*LockErcRequest, error) {
	ercParams := &LockErcRequest{}
	ethTx, err := DecodeTransaction(rawEthTx)
	if err != nil {
		return ercParams, err
	}
	token, err := GetToken(erc20list, *ethTx.To())
	if err != nil {
		return ercParams, err
	}
	contractAbi, err := abi.JSON(strings.NewReader(token.TokAbi))
	if err != nil {
		return ercParams, err
	}
	functionSignature, err := getSignFromName(&contractAbi, "transfer", contract.ERC20BasicFuncSigs)
	if err != nil {
		return ercParams, err
	}
	ercParams, err = parseERC20Lock(rawEthTx, functionSignature)
	if err != nil {
		return ercParams, err
	}
	return ercParams, nil
}

// VerifyLock verifies ETHER lock transaction to check whether the ETHER is being sent to our Smart Contract Address.
func VerifyLock(tx *types.Transaction, contractabi string) (bool, error) {

	contractAbi, err := abi.JSON(strings.NewReader(contractabi))
	if err != nil {
		return false, errors.Wrap(err, "Unable to get contract Abi from ChainDriver options")
	}
	bytesData, err := contractAbi.Pack("lock")
	if err != nil {
		return false, errors.Wrap(err, "Unable to to create Bytes data for Lock")
	}
	return bytes.Equal(bytesData, tx.Data()), nil

}
// VerifyLock verifies ERC lock transaction to check whether the ERC20 Token is being sent to our Smart Contract Address.
func VerfiyERC20Lock(rawTx []byte, tokenabi_ string, erc20contractaddr common.Address) (bool, error) {
	contractAbi, err := abi.JSON(strings.NewReader(tokenabi_))
	if err != nil {
		return false, errors.Wrap(err, "Unable to get contract Abi for Test Token from ChainDriver options")
	}
	methodSignature, err := getSignFromName(&contractAbi, "transfer", contract.ERC20BasicFuncSigs)
	if err != nil {
		return false, err
	}
	ercLockParams, err := parseERC20Lock(rawTx, methodSignature)
	if err != nil {
		return false, err
	}
	return bytes.Equal(ercLockParams.Receiver.Bytes(), erc20contractaddr.Bytes()), nil
}

// ParseERC20RedeemToken Parses the ERC Redeem Request and returns the Token which is being redeemed
func ParseERC20RedeemToken(rawTx []byte, tokenList []ERC20Token, lockredeemERCAbi string) (*ERC20Token, error) {
	ercRedeemParams, err := ParseERC20RedeemParams(rawTx, lockredeemERCAbi)
	if err != nil {
		return nil, err
	}
	token, err := GetToken(tokenList, ercRedeemParams.TokenAddress)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ParseERC20RedeemParams Parses the ERC Redeem Request and returns the RedeemERCRequest struct
// RedeemERCRequest contains the parameters that were passed when the redeem function in the LOCKREDEEMERC smart contract was called by the user
func ParseERC20RedeemParams(rawTx []byte, lockredeemERCAbi string) (*RedeemErcRequest, error) {

	contractAbi, err := StringTOABI(lockredeemERCAbi)
	if err != nil {
		return nil, err
	}
	methodSignature, err := getSignFromName(contractAbi, "redeem", contract.LockRedeemERCFuncSigs)
	if err != nil {
		return nil, err
	}
	ercRedeemParams, err := parseERC20Redeem(rawTx, methodSignature)
	if err != nil {
		return nil, err
	}
	return ercRedeemParams, nil
}
// ParseRedeem Parses the ERC Redeem Request and returns the RedeemRequest struct
// RedeemRequest contains the parameters that were passed when the redeem function in the LOCKREDEEM smart contract was called by the user
func ParseRedeem(data []byte, lockredeemAbi string) (req *RedeemRequest, err error) {
	contractAbi, err := StringTOABI(lockredeemAbi)
	if err != nil {
		return nil, err
	}
	methodSignature, err := getSignFromName(contractAbi, "redeem", contract.LockRedeemFuncSigs)
	if err != nil {
		return nil, err
	}
	ss := strings.Split(hex.EncodeToString(data), methodSignature)
	if len(ss) == 0 {
		return nil, errors.New("Transaction does not have the required input data")
	}
	if len(ss[1]) < 64 {
		return nil, errors.New("Transaction data is invalid")
	}
	d, err := hex.DecodeString(ss[1][:64])
	if err != nil {
		return nil, err
	}
	amt := big.NewInt(0).SetBytes(d)
	return &RedeemRequest{Amount: amt}, nil
}
// ParseLock takes a rawTX byteArray as input and returns a LockRequest struct
// LockRequest struct contains the amount which the user is trying to lock .
func ParseLock(data []byte) (req *LockRequest, err error) {

	tx, err := DecodeTransaction(data)
	if err != nil {
		return nil, err
	}

	return &LockRequest{Amount: tx.Value()}, nil
}

//DecodeTransaction takes a rawTx bytearray as input and returns a transaction
func DecodeTransaction(data []byte) (*types.Transaction, error) {
	tx := &types.Transaction{}

	err := rlp.DecodeBytes(data, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode Bytes")
	}

	return tx, nil
}


// GetToken takes a tokenlist and token address as input .
// It iterates over the Tokenlist and returns the token which matches the token address
func GetToken(erc20list []ERC20Token, tokAddr common.Address) (*ERC20Token, error) {
	for _, token := range erc20list {
		if token.TokAddr == tokAddr {
			return &token, nil
		}
	}
	return &ERC20Token{}, errors.New("Token not supported")
}

// StringTOABI takes string as an input and converts it into an ABI
func StringTOABI(contractAbi string) (*abi.ABI, error) {
	Abi, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get contract Abi for Test Token from ChainDriver options")
	}
	return &Abi, nil
}