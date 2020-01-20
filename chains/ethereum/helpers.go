package ethereum

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)


func mapkey(m map[string]string, value string) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}
func parseERC20Lock(data []byte, functionSig string) (req *LockErcRequest, err error) {
	ss := strings.Split(hex.EncodeToString(data), functionSig)

	tokenAmount, err := hex.DecodeString(ss[1][64:128])
	if err != nil {
		return nil, err
	}
	receiver := ss[1][24:64]
	amt := big.NewInt(0).SetBytes(tokenAmount)
	return &LockErcRequest{
		Receiver:    common.HexToAddress(receiver),
		TokenAmount: amt,
	}, nil
}

func parseERC20Redeem(data []byte, functionSig string) (req *RedeemErcRequest, err error) {
	ss := strings.Split(hex.EncodeToString(data), functionSig)
	tokenAddress := ss[1][88:128]
	amount,err := hex.DecodeString(ss[1][:64])
	if err != nil {
		return nil, err
	}
	amt := big.NewInt(0).SetBytes(amount)
	fmt.Println(tokenAddress)
	return &RedeemErcRequest{
		Amount:    amt,
		TokenAddress: common.HexToAddress(tokenAddress),
	}, nil
}

func getSignfromName(contractAbi *abi.ABI, methodName string,funcSigs map[string]string) (string, error) {
	method, exists := contractAbi.Methods[methodName]
	if !exists {
		return "", errors.New("Function not found in abi ")
	}
	signature, ok := mapkey(funcSigs, method.Sig())
	if !ok {
		return "", errors.New("Method Signature does not exist")
	}
	return signature, nil
}