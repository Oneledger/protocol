package common

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/global"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/btcsuite/btcutil"
	"github.com/Oneledger/protocol/node/log"
	"reflect"
	"github.com/ethereum/go-ethereum/common"
	"github.com/Oneledger/protocol/node/action"
)

func GetSwapAddress(chain data.ChainType) []byte {

	if chain == data.BITCOIN {
		cli := bitcoin.GetBtcClient(global.Current.BTCAddress, &chaincfg.RegressionNetParams )

		return []byte(bitcoin.GetRawAddress(cli).String())
	} else if chain == data.ETHEREUM {

		return ethereum.GetAddress().Bytes()
	}
	return nil
}



func GetAddressFromByteArray(chain data.ChainType, address string, target interface{}) {

	if chain == data.BITCOIN {
		result, err := btcutil.DecodeAddress( address , &chaincfg.RegressionNetParams)
		if err != nil {
			log.Error("failed to get addressPubKeyHash")
			return
		}

		switch target.(type) {

		case *btcutil.AddressPubKeyHash:
			target = result.(*btcutil.AddressPubKeyHash)
		default:
			log.Fatal("not appropriate address type to convert yet", "address", address, "target", reflect.TypeOf(target))
		}
		return

	} else if chain == data.ETHEREUM {
		switch target.(type) {

		case common.Address:
			target = common.BytesToAddress([]byte(address))
		default:
			log.Fatal("not appropriate address type to convert yet", "address", address, "target", reflect.TypeOf(target))
		}
		return
	} else {
		log.Fatal("not supported chain", "chain", chain)
	}
}

type Contract interface {
	ToMessage() action.Message
}