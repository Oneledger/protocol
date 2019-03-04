package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/common"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"math/big"
)

type OneledgerDriver struct {
	Driver ChainDriverBase
}

func init() {
	serial.Register(OneledgerDriver{})
}

func (driver OneledgerDriver) GetURL() string {
	return global.Current.AppAddress
}

func (driver OneledgerDriver) GetChainAddress(chainKey interface{}) []byte{
	return nil
}

func (driver OneledgerDriver) GetMethodsList() []string {
	return nil
}

func (driver OneledgerDriver) ExecuteMethod(method string, params []byte) status.Code {
	return status.NOT_IMPLEMENTED
}

func (driver OneledgerDriver) GetAddressFromByteArray(address []byte) interface{} {
	return nil
}

func (driver OneledgerDriver) CreateSwapContract(receiver interface{}, account id.Account, value big.Int, timeout int64, hash [32]byte) common.Contract {
	return nil
}

func (driver OneledgerDriver) CreateSwapContractFromMessage(message []byte) common.Contract{
	return nil
}
