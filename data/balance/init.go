/*

   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

	Copyright 2017 - 2019 OneLedger

*/


package balance

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/serialize"
	"math/big"
)

var pSzlr serialize.Serializer


var currencies map[string]Currency
var currenciesExtra map[string]Extra

func init() {

	serial.Register(Balance{})
	serial.Register(Coin{})
	serial.Register(Currency{})

	serialize.RegisterConcrete(new(Balance), TagBalance)
	serialize.RegisterConcrete(new(BalanceData), TagBalanceData)

	/*
	currencies = map[string]Currency{
		OLT: {OLT, ONELEDGER, 0},
		BTC: {BTC, BITCOIN, 1},
		ETH: {ETH, ETHEREUM, 2},
		VT:  {VT, ONELEDGER, 3},
	}

	currenciesExtra = map[string]Extra{
		OLT: {big.NewFloat(1000000000000000000), 6, 'f'},
		BTC: {big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
		ETH: {big.NewFloat(1), 0, 'f'}, // TODO: This needs to be set correctly
		VT:  {big.NewFloat(1), 0, 'f'},
	}
	*/
	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)

}


func RegisterCurrency(name string, ct chain.ChainType, id int,
						units *big.Float, decimal int, format uint8) {

	currencies[name] = Currency{name, ct, id}
	currenciesExtra[name] = Extra{units, decimal, format}
}