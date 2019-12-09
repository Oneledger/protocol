package ethereum

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"
	"time"

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
	ethConfig     = config.DefaultEthConfig()
	LockRedeemABI = contract.LockRedeemABI
	contractAddr  = "0xAbc23065E977386EccE670f6ba143fe90e40F240"
	priv_key      = "6c24a44424c8182c1e3e995ad3ccfb2797e3f7ca845b99bea8dead7fc9dccd09"
	addr          = "0x3760136A6327e6E3f6e6284e2173dA2414e89894"
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

func BroadCastLock() (*TransactionHash,error) {

	rawtx,err := ethCD.PrepareUnsignedETHLock(getAddress(),big.NewInt(10000000000))
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
	txHash, err := ethCD.BroadcastTx(signedTx)
	if err !=nil {
		return nil,errors.Wrap(err,"Unable to get Broadcast")
	}
	return &txHash,nil
}

func TestETHChainDriver_CheckFinality(t *testing.T) {
	 txHash,err := BroadCastLock()
	 if err !=nil {
	 	logger.Error(err)
		 return
	 }
	 isFinalized := false
     for !isFinalized {
     	rec,err := ethCD.CheckFinality(txHash)
     	if err != nil {
     		logger.Error(err)
            BroadCastLock()
     		time.Sleep(5*time.Second)
		}
		if rec != nil {
			logger.Info("Transaction confirmed , The transaction had been included at " ,rec.BlockNumber  )
			isFinalized = true
		}
	 }
}
