/*
	Copyright 2017 - 2018 OneLedger
*/

package bitcoin

import (
	"time"

	brpc "github.com/Oneledger/protocol/node/chains/bitcoin/rpc"
	"github.com/btcsuite/btcd/chaincfg"

	"encoding/base64"
	"net"
	"strings"

	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/log"
)

func GetBtcClient(address string, id int, chainParams *chaincfg.Params) *brpc.Bitcoind {
	addr := strings.Split(address, ":")
	if len(addr) < 2 {
		log.Error("address not in correct format", "fullAddress", address)
	}

	ip := net.ParseIP(addr[0])
	if ip == nil {
		log.Error("address can not be parsed", "addr", addr)
	}

	port := convert.GetInt(addr[1], 46688)

	// TODO: Needs to be passed in as a param
	var usr, pass string
	switch id {
	case 1:
		usr, pass = getCredential()
	case 2:
		usr = "oltest02"
		pass = "olpass02"
	case 3:
		usr = "oltest03"
		pass = "olpass03"
	default:
		log.Fatal("Invalid", "id", id)
	}

	cli, err := brpc.New(ip.String(), port, usr, pass, false, chainParams)
	if err != nil {
		log.Error("Can't get the btc rpc client at given address", "err", err)
		return nil
	}

	return cli
}

func getCredential() (usr string, pass string) {
	//todo: getCredential from database which should be randomly generated when register or import if user already has bitcoin node
	usrBytes, err := base64.StdEncoding.DecodeString("b2x0ZXN0MDE=")
	if err != nil {
		log.Error(err.Error())
		usr = ""
	}
	usr = string(usrBytes)
	passBytes, err := base64.StdEncoding.DecodeString("b2xwYXNzMDE=")
	if err != nil {
		log.Error(err.Error())
		pass = ""
	}
	pass = string(passBytes)
	return
}

func ScheduleBlockGeneration(cli brpc.Bitcoind, interval time.Duration) chan bool {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				cli.Generate(1)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}

func StopBlockGeneration(stop chan bool) {
	close(stop)
}
