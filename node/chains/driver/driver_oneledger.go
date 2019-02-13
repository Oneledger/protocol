package chaindriver

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
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

func (driver OneledgerDriver) GetChainAddress() []byte{
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
