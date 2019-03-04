package chaindriver

import (
	oneledger_common "github.com/Oneledger/protocol/node/chains/common"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type EtheriumAddress common.Address

type EthereumDriver struct {
	Driver ChainDriverBase
}

func init() {
	serial.Register(EthereumDriver{})
}

func (driver EthereumDriver) GetURL() string {
	return global.Current.ETHAddress
}

func (driver EthereumDriver) GetChainAddress(chainKey interface{}) []byte{
	return ethereum.GetAddress(chainKey).Bytes()
}

func (driver EthereumDriver) GetMethodsList() []string {
	list := []string{"getinfo"}

	return list
}

func (driver EthereumDriver) ExecuteMethod(method string, params []byte) status.Code {
	return status.NOT_IMPLEMENTED
}

func (driver EthereumDriver) CreateSwapContract(receiver []byte, account id.Account, value big.Int, timeout int64, hash [32]byte) oneledger_common.Contract {
	address := common.BytesToAddress(receiver)

	contract := ethereum.CreateHtlContract(account.GetChainKey())

	if contract == nil {
		return nil
	}

	log.Debug("Create ETH HTLC", "value", value, "receiver", receiver, "hash", hash)

	err := contract.Funds(account.GetChainKey(), &value, big.NewInt(timeout), address, hash)
	if err != nil {
		return nil
	}

	return contract
}

func (driver EthereumDriver) CreateSwapContractFromMessage(message []byte) oneledger_common.Contract{
	contract := &ethereum.HTLContract{}

	contract.FromBytes(message)

	return contract
}
