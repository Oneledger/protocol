package ethereum

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)
var (

	ethCD *ETHChainDriver
	ethConfig = config.DefaultEthConfig()
	LockRedeemABI = contract.LockRedeemABI
	contractAddr = "0xDae163dC84adD22c138CD7300AD0d22cfF2f2ec8"
	priv_key = "6c24a44424c8182c1e3e995ad3ccfb2797e3f7ca845b99bea8dead7fc9dccd09"
	addr = "0x"
	logger = log.NewLoggerWithPrefix(os.Stdout,"TestChainDRIVER")
)
func init() {

	cdoptions := ChainDriverOption{
		ContractABI: LockRedeemABI,
		ContractAddress: common.HexToAddress(contractAddr),
	}
	ethCD,_ = NewChainDriver(ethConfig,logger,&cdoptions)

}

func TestETHChainDriver_PrepareUnsignedETHRedeem(t *testing.T) {

}

func TestETHChainDriver_ParseRedeem(t *testing.T) {
	privkey, err := crypto.HexToECDSA(priv_key)
	assert.NoError(t, err, "wrong private key")
	pubkey := privkey.Public()
	publicKeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	assert.True(t, ok, "error casting public key to ECDSA")

	addr := crypto.PubkeyToAddress(*publicKeyECDSA)
	amount := big.NewInt(10000000000)
	tx, err := ethCD.PrepareUnsignedETHRedeem(addr, amount)
	chainID, err := ethCD.Client.NetworkID(context.Background())

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privkey)

	ts := types.Transactions{signedTx}
	data := ts.GetRlp(0)
 	req, err := ethCD.ParseRedeem(data)
 	assert.NoError(t, err, "failed to parse signed tx")
	assert.Equal(t, amount, req.Amount)
}