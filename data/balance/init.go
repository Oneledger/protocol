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
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/serialize"
	"math/big"
	"os"
)

var pSzlr serialize.Serializer
var logger *log.Logger

var currencies map[string]Currency
var currenciesExtra map[string]Extra

func init() {

	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("balance")
	serial.Register(Balance{})
	serial.Register(Coin{})
	serial.Register(Currency{})

	serialize.RegisterConcrete(new(Balance), TagBalance)
	serialize.RegisterConcrete(new(BalanceData), TagBalanceData)

	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)

}

func RegisterCurrency(name string, ct chain.Type, id int,
	units *big.Float, decimal int, format uint8) {

	currencies[name] = Currency{name, ct, id}
	currenciesExtra[name] = Extra{units, decimal, format}
}
