package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type BitcoinDriver struct {
	Driver ChainDriverBase
}

func init() {
	serial.Register(BitcoinDriver{})
}

func (driver BitcoinDriver) GetURL() string {
	return global.Current.BTCAddress
}

func (driver BitcoinDriver) GetChainAddress() []byte{
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress)
	return []byte(bitcoin.GetRawAddress(cli).String())
}

func (driver BitcoinDriver) GetMethodsList() []string {
	list := []string{"getinfo"}

	return list
}

func (driver BitcoinDriver) ExecuteMethod(method string, params []byte) status.Code {
	return status.NOT_IMPLEMENTED
}

func (driver BitcoinDriver) GetAddressFromByteArray(address []byte) interface{} {
	result, err := btcutil.DecodeAddress(string(address), &chaincfg.RegressionNetParams)

	if err != nil {
		log.Error("failed to get addressPubKeyHash")
		return nil
	}

	return result.(*btcutil.AddressPubKeyHash)
}
