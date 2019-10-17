/*

 */

package internal

import (
	"fmt"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type BroadcastBTCData struct {
	tx          *wire.MsgTx
	trackerName string
}

func BroadcastBTC(data interface{}) {

	broadcastData := data.(BroadcastBTCData)

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18831",
		User:         "oltest01",
		Pass:         "olpass01",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	clt, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	cd := bitcoin.NewChainDriver("")
	txid, err := cd.BroadcastTx(broadcastData.tx, clt)
	if err != nil {
		return
	}

	// save to internal db
	fmt.Println(txid)
}
