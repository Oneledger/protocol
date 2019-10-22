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
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

var pSzlr serialize.Serializer
var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("balance")

	serialize.RegisterConcrete(new(Balance), TagBalance)
	//serialize.RegisterConcrete(new(BalanceData), TagBalanceData)

	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)
}
