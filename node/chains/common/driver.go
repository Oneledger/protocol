package common

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/global"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/btcsuite/btcutil"
	"github.com/Oneledger/protocol/node/log"
	"github.com/ethereum/go-ethereum/common"
)

func GetSwapAddress(chain data.ChainType) []byte {
    switch chain {

    case data.BITCOIN:
        cli := bitcoin.GetBtcClient(global.Current.BTCAddress, &chaincfg.RegressionNetParams )
        return []byte(bitcoin.GetRawAddress(cli).String())
    case data.ETHEREUM:
        return ethereum.GetAddress().Bytes()

    default:
        return nil
    }
}



func GetBTCAddressFromByteArray(chain data.ChainType, address []byte) *btcutil.AddressPubKeyHash {

	if chain == data.BITCOIN {
		result, err := btcutil.DecodeAddress( string(address), &chaincfg.RegressionNetParams)
		if err != nil {
			log.Error("failed to get addressPubKeyHash")
			return nil
		}

		return result.(*btcutil.AddressPubKeyHash)
	} else {
		log.Fatal("not supported chain", "chain", chain)
	}
	return nil
}

func GetETHAddressFromByteArray(chain data.ChainType, address []byte) *common.Address {
    var result common.Address
    if chain == data.ETHEREUM {
        result = common.BytesToAddress(address)
        log.Debug("ethereum address","address", address, "resuslt", result)
        return &result
    } else {
        log.Fatal("not supported chain", "chain", chain)
    }

    return nil
}

type Contract interface {
    ToMessage() []byte
}
