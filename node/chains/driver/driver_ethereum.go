package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/ethereum/go-ethereum/common"
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

func (driver EthereumDriver) GetChainAddress() []byte{
	return ethereum.GetAddress().Bytes()
}

func (driver EthereumDriver) GetMethodsList() []string {
	list := []string{"getinfo"}

	return list
}

func (driver EthereumDriver) ExecuteMethod(method string, params []byte) status.Code {
	return status.NOT_IMPLEMENTED
}

func (driver EthereumDriver) GetAddressFromByteArray(address []byte) interface{} {
	result := common.BytesToAddress(address)
	log.Debug("ethereum address", "address", address, "resuslt", result)
	return result
}
