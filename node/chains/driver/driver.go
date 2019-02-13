package chaindriver

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

type ChainDriver interface {
	GetURL() string
	GetChainAddress() []byte
	GetMethodsList() []string                   // @TODO return as strings now, but probably need to replace them with callback methods?
	ExecuteMethod(string, []byte) status.Code   // @TODO should the execute method return anything else?
	GetAddressFromByteArray([]byte) interface{}
}

type ChainDriverBase struct {
}

var drivers map[data.ChainType]ChainDriver

func init() {
	serial.Register(ChainDriverBase{})

	// @TODO <temp> - Need to move initialization to a different place maybe
	drivers = make(map[data.ChainType]ChainDriver)

	drivers[data.BITCOIN] = ChainDriver(BitcoinDriver{})
	drivers[data.ETHEREUM] = ChainDriver(EthereumDriver{})
	drivers[data.ONELEDGER] = ChainDriver(OneledgerDriver{})
}

func GetDriverList() []data.ChainType {
	var list []data.ChainType

	for k, _ := range drivers {
		list = append(list, k)
	}

	return list
}

func GetDriver(id data.ChainType) ChainDriver {
	return drivers[id]
}
