package ethereum

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
)

var (
	ethCD         *ETHChainDriver
	ethConfig     = config.DefaultEthConfig()
	LockRedeemABI = contract.LockRedeemABI
	contractAddr  = "0xDae163dC84adD22c138CD7300AD0d22cfF2f2ec8"
	priv_key      = "6c24a44424c8182c1e3e995ad3ccfb2797e3f7ca845b99bea8dead7fc9dccd09"
	addr          = "0x"
	logger        = log.NewLoggerWithPrefix(os.Stdout, "TestChainDRIVER")
)

func init() {

	cdoptions := ChainDriverOption{
		ContractABI:     LockRedeemABI,
		ContractAddress: common.HexToAddress(contractAddr),
	}
	ethCD, _ = NewChainDriver(ethConfig, logger, &cdoptions)

}

func TestETHChainDriver_PrepareUnsignedETHRedeem(t *testing.T) {

}
