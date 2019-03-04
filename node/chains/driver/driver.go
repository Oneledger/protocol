package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/common"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"math/big"
)

type ChainDriver interface {
	GetURL() string
	GetChainAddress(chainKey interface{}) []byte
	GetMethodsList() []string                   // @TODO return as strings now, but probably need to replace them with callback methods?
	ExecuteMethod(string, []byte) status.Code   // @TODO should the execute method return anything else?
	GetAddressFromByteArray([]byte) interface{}
	CreateSwapContract(receiver interface{}, account id.Account, value big.Int, timeout int64, hash [32]byte) common.Contract
	CreateSwapContractFromMessage(message []byte) common.Contract
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
