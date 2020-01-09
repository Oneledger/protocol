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
	ethCD         *ETHChainDriver
	ethConfig     = config.DefaultEthConfig()
	LockRedeemABI = contract.LockRedeemABI
	ERCLockRedeemABI =contract.LockRedeemERCABI
	TestTokenABI   =   contract.ERC20BasicABI
	contractAddr  = "0xdaF6850A7545705a80AB802f3b951dA0a635CE78"
	ERCcontractAddr ="0x4c6970C8d6748e21fAe1Bb5A4b7AE9bE2825bdFA"
	TestTokenContractAddr = "0x28075b7419bAA72420fe6539c95ADB8E8D039841"
	priv_key      = "0e342e8ad59b75c4dd8af7340b85efb3308a09a54eded1ee194b83132b6b1395"
	addr          = "0xa9258c306f392380E7A9aCcaD3C35230f7FC42F8"
	logger        = log.NewLoggerWithPrefix(os.Stdout, "TestChainDRIVER")
	toAddress = common.HexToAddress(contractAddr)
	toAddressTestToken = common.HexToAddress(TestTokenContractAddr)
	toAdddressLockRedeemERC = common.HexToAddress(ERCcontractAddr)
	valuelockERC20           = big.NewInt(10)
	client *ethclient.Client
)

func init() {

	cdoptions := ChainDriverOption{
		ContractABI:     LockRedeemABI,
		ContractAddress: common.HexToAddress(contractAddr),
	}
	ethCD, _ = NewChainDriver(ethConfig, logger, &cdoptions)
	client,_ = ethConfig.Client()

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

func CreateUnsignedERCLock() ([]byte,error){
	tokenAbi, _ := abi.JSON(strings.NewReader(TestTokenABI))
	bytesData, err := tokenAbi.Pack("transfer",toAdddressLockRedeemERC,valuelockERC20)
	nonce, err := client.PendingNonceAt(context.Background(), getAddress())
	if err != nil {
		return nil,err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil,err
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
		return nil,err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), getPrivKey())

	if err != nil {
		return nil,err
	}
	ts := types.Transactions{signedTx}
	rawTxBytes := ts.GetRlp(0)
	txNew := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes, txNew)
	return rawTxBytes,nil
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

//func TestETHChainDriver_CheckFinality(t *testing.T) {
//	 txHash,err := BroadCastLock()
//	 if err !=nil {
//	 	logger.Error(err)
//		 return
//	 }
//	 isFinalized := false
//     for !isFinalized {
//     	rec,err := ethCD.CheckFinality(*txHash)
//     	if err != nil {
//     		logger.Error(err)
//           // BroadCastLock()
//     		time.Sleep(5*time.Second)
//		}
//		if rec != nil {
//			logger.Info("Transaction confirmed , The transaction had been included at " ,rec.BlockNumber  )
//			isFinalized = true
//		}
//	 }
//}

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
	rawERCLock,err := CreateUnsignedERCLock()
	if err != nil {
		fmt.Println(err)
		return
	}
	ok,err:= VerfiyERC20Lock(rawERCLock,TestTokenABI,toAdddressLockRedeemERC)
	if err == nil {
		fmt.Println("VERIFY ERC LOCK :", ok)
	}

}