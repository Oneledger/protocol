package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

var (
	ethCD                       *ETHChainDriver
	ethConfig                   = config.DefaultEthConfig("", "")
	LockRedeemABI               = contract.LockRedeemABI
	ERCLockRedeemABI            = contract.LockRedeemERCABI
	TestTokenABI                = contract.ERC20BasicABI
	LockRedeemContractAddr      = "0xFF51ABac8c8664e83AB0d94baac7312fD59ab873"
	TestTokenContractAddr       = "0xF055145EC2607feAcdDD732a8338a4311F961eD9"
	LockRedeemERC20ContractAddr = "0x34Ea04be3aC452BA23d8c56298178C8163674F0d"
	priv_key                    = "bdb082c7e42a946c477fa3efee4fb5bdece508b47592d8cb57f5e811cd840a40"
	addr                        = "0xa9258c306f392380E7A9aCcaD3C35230f7FC42F8"
	logger                      = log.NewLoggerWithPrefix(os.Stdout, "TestChainDRIVER")
	toAddress                   = common.HexToAddress(LockRedeemContractAddr)
	toAddressTestToken          = common.HexToAddress(TestTokenContractAddr)
	toAdddressLockRedeemERC     = common.HexToAddress(LockRedeemERC20ContractAddr)
	valuelockERC20              = big.NewInt(10)
	client                      *ethclient.Client
)

func init() {

	cdoptions := ChainDriverOption{
		ContractABI:     LockRedeemABI,
		ContractAddress: common.HexToAddress(LockRedeemContractAddr),
	}
	ethCD, _ = NewChainDriver(ethConfig, logger, cdoptions.ContractAddress, cdoptions.ContractABI, ETH)
	client, _ = ethclient.Dial(ethConfig.Connection)

}

func getPrivKey() *ecdsa.PrivateKey {
	UserprivKey, err := crypto.HexToECDSA(priv_key)
	if err != nil {
		logger.Error("Unable to create Private key")
	}
	return UserprivKey
}
func getAddress() common.Address {

	pubkey := getPrivKey().Public()
	ecdsapubkey, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}
	}
	addr := crypto.PubkeyToAddress(*ecdsapubkey)
	return addr
}

func GetSignedLockTX() (*Transaction, error) {
	rawtx, err := ethCD.PrepareUnsignedETHLock(getAddress(), big.NewInt(100000000000000000))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get rawLockTX")
	}
	chainid, err := ethCD.ChainId()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get ChainID")
	}
	tx, err := ethCD.DecodeTransaction(rawtx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to Decode rawTX to Transaction")
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainid), getPrivKey())
	if err != nil {
		return nil, errors.Wrap(err, "Unable to Sign Transaction")
	}
	return signedTx, nil
}

func CreateUnsignedERCLock() ([]byte, error) {
	tokenAbi, _ := abi.JSON(strings.NewReader(TestTokenABI))
	bytesData, err := tokenAbi.Pack("transfer", toAdddressLockRedeemERC, valuelockERC20)
	nonce, err := client.PendingNonceAt(context.Background(), getAddress())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974) // in units

	auth := bind.NewKeyedTransactor(getPrivKey())
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	tx := types.NewTransaction(nonce, toAddressTestToken, big.NewInt(0), gasLimit, gasPrice, bytesData)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), getPrivKey())

	if err != nil {
		return nil, err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	txNew := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes, txNew)
	return rawTxBytes, nil
}

func CreateERC20Redeem() ([]byte, error) {
	ERCLockRedeemAbi, _ := abi.JSON(strings.NewReader(ERCLockRedeemABI))
	bytesData, err := ERCLockRedeemAbi.Pack("redeem", valuelockERC20, toAddressTestToken)
	nonce, err := client.PendingNonceAt(context.Background(), getAddress())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974) // in units

	auth := bind.NewKeyedTransactor(getPrivKey())
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	tx := types.NewTransaction(nonce, toAdddressLockRedeemERC, big.NewInt(0), gasLimit, gasPrice, bytesData)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), getPrivKey())

	if err != nil {
		return nil, err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes, _ := rlp.EncodeToBytes(ts[0])
	txNew := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes, txNew)
	return rawTxBytes, nil

}
func BroadCastLock() (*TransactionHash, error) {

	signedTx, err := GetSignedLockTX()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to getSigned Transaction")
	}
	txHash, err := ethCD.BroadcastTx(signedTx)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get Broadcast")
	}
	return &txHash, nil
}

func TestETHChainDriver_CheckFinality(t *testing.T) {
	txHash, err := BroadCastLock()
	if err != nil {
		logger.Error(err)
		return
	}
	isFinalized := false
	for !isFinalized {
		status := ethCD.CheckFinality(*txHash, 2)
		if status == TXSuccess {
			logger.Info("Transaction confirmed")
			isFinalized = true
		}
	}
}

//func TestETHChainDriver_VerifyLock(t *testing.T) {
//	signedTx,err := GetSignedLockTX()
//	if err != nil {
//		fmt.Println(errors.Wrap(err,"Unable to getSigned Transaction"))
//		return
//	}
//	ok,err := VerifyLock(signedTx,LockRedeemABI)
//	if err != nil {
//		fmt.Println("Unable to verify lock transaction" ,err)
//		return
//		} else if !ok {
//		fmt.Println("Bytes data does not match (function name field is different)")
//		 return
//	}
//}

//func TestETHChainDriver_DecodeTransaction(t *testing.T) {
//
//	data := []byte{248, 119, 30, 139, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 131, 102, 145, 182, 148, 133, 169, 40, 34, 122, 76, 57, 145, 29, 183, 60, 107, 219, 177, 34, 247, 235, 75, 91, 41, 136, 13, 224, 182, 179, 167, 100, 0, 0, 132, 248, 61, 8, 186, 38, 160, 88, 197, 150, 4, 107, 150, 167, 134, 243, 218, 151, 78, 168, 77, 142, 208, 57, 238, 47, 174, 179, 216, 250, 88, 161, 95, 84, 224, 232, 164, 87, 223, 160, 32, 180, 100, 16, 32, 148, 1, 233, 140, 79, 133, 239, 61, 229, 214, 77, 17, 38, 236, 84, 206, 166, 250, 209, 255, 76, 62, 197, 84, 30, 2, 110}
//	tx, err := DecodeTransaction(data)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(tx.GasPrice())
//
//}

func TestVerfiyERC20Lock(t *testing.T) {
	rawERCLock, err := CreateUnsignedERCLock()
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err := VerfiyERC20Lock(rawERCLock, TestTokenABI, toAdddressLockRedeemERC)
	if err == nil {
		fmt.Println("VERIFY ERC LOCK :", ok)
	}

}

func TestParseERC20Redeem(t *testing.T) {
	rawERCRedeem, err := CreateERC20Redeem()
	if err != nil {
		fmt.Println(err)
		return
	}
	ERCLockRedeemAbi, err := abi.JSON(strings.NewReader(ERCLockRedeemABI))
	if err != nil {
		fmt.Println(err)
		return
	}
	sig, err := getSignFromName(&ERCLockRedeemAbi, "redeem", contract.LockRedeemERCFuncSigs)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := parseERC20Redeem(rawERCRedeem, sig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(req)
}
