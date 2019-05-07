/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package config

import (
	"math/big"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
)

func init() {
	// reblws: This should be moved up to app.setupState
	chain.RegisterChainType("UNKNOWN", 0)
	chain.RegisterChainType("OneLedger", 1)
	chain.RegisterChainType("Bitcoin", 2)
	chain.RegisterChainType("Ethereum", 3)

	balance.RegisterCurrency("OLT", chain.Type(1), big.NewFloat(1000000000000000000), 6, 'f')
	balance.RegisterCurrency("BTC", chain.Type(2), big.NewFloat(1), 0, 'f')
	balance.RegisterCurrency("ETH", chain.Type(3), big.NewFloat(1), 0, 'f')
	balance.RegisterCurrency("VT", chain.Type(1), big.NewFloat(1), 0, 'f')
}
