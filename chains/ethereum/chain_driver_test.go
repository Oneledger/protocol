package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

var (
	ethCD         *ETHChainDriver
	ethConfig     = config.DefaultEthConfigRoopsten()
	LockRedeemABI = contract.LockRedeemABI
	contractAddr  = "0xdaF6850A7545705a80AB802f3b951dA0a635CE78"
	priv_key      = "02038529C9AB706E9F4136F4A4EB51E866DBFE22D5E102FD3A22C14236E1C2EA"
	addr          = "0xa9258c306f392380E7A9aCcaD3C35230f7FC42F8"
	logger        = log.NewLoggerWithPrefix(os.Stdout, "TestChainDRIVER")
)

func init() {

	cdoptions := ChainDriverOption{
		ContractABI:     LockRedeemABI,
		ContractAddress: common.HexToAddress(contractAddr),
	}
	ethCD, _ = NewChainDriver(ethConfig, logger, &cdoptions)

}

func getPrivKey () *ecdsa.PrivateKey {
	UserprivKey, err := crypto.HexToECDSA(priv_key)
	if err != nil {
		logger.Error("Unable to create Private key")
	}
	return UserprivKey
}
func getAddress () common.Address {

	pubkey := getPrivKey().Public()
	ecdsapubkey, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}
	}
	addr := crypto.PubkeyToAddress(*ecdsapubkey)
	return addr
}

func GetSignedLockTX() (*Transaction,error){
	rawtx,err := ethCD.PrepareUnsignedETHLock(getAddress(),big.NewInt(100000000000000000))
	if err != nil {
		return nil,errors.Wrap(err,"Unable to get rawLockTX")
	}
	chainid,err := ethCD.ChainId()
	if err != nil {
		return nil,errors.Wrap(err,"Unable to get ChainID")
	}
	tx,err := ethCD.DecodeTransaction(rawtx)
	if err != nil {
		return nil,errors.Wrap(err,"Unable to Decode rawTX to Transaction")
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainid), getPrivKey())
	if err != nil {
		return nil,errors.Wrap(err,"Unable to Sign Transaction")
	}
	return signedTx,nil
}
func BroadCastLock() (*TransactionHash,error) {

	signedTx,err := GetSignedLockTX()
	if err != nil {
		return nil,errors.Wrap(err,"Unable to getSigned Transaction")
	}
	txHash, err := ethCD.BroadcastTx(signedTx)
	if err !=nil {
		return nil,errors.Wrap(err,"Unable to get Broadcast")
	}
	return &txHash,nil
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

func TestETHChainDriver_DecodeTransaction(t *testing.T) {
	buf := []byte{248,115,10,139,50,48,48,48,48,48,48,48,48,48,48,131,102,145,182,148,182,187,159,53,245,78,132,228,139,237,130,103,118,252,178,107,156,200,135,30,132,49,48,48,48,132,248,61,8,186,37,160,34,194,136,168,116,70,194,73,53,107,128,167,228,196,4,111,174,231,112,208,91,135,76,34,235,102,202,126,252,125,21,41,160,115,138,172,210,188,154,0,109,19,90,90,255,2,46,255,70,238,239,187,141,250,178,209,104,23,208,139,120,123,137,187,31}
	tx, err := DecodeTransaction(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx.GasPrice())
}
