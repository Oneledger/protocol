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
	"github.com/btcsuite/btcutil"
	"strconv"
	"github.com/btcsuite/btcd/wire"
	"github.com/Oneledger/protocol/node/comm"
)

func GetBtcClient(address string, chainParams *chaincfg.Params) *brpc.Bitcoind {
	addr := strings.Split(address, ":")
	if len(addr) < 2 {
		log.Error("address not in correct format", "fullAddress", address)
	}

	ip := net.ParseIP(addr[0])
	if ip == nil {
		log.Error("address can not be parsed", "addr", addr)
	}

	port := convert.GetInt(addr[1], 46688)

	usr, pass := getCredential(port)

	cli, err := brpc.New(ip.String(), port, usr, pass, false, chainParams)
	if err != nil {
		log.Error("Can't get the btc rpc client at given address", "err", err)
		return nil
	}

	return cli
}

func getCredential(port int) (usr string, pass string) {

	var u, p string
	switch port {
	case 18831 :
		u = "b2x0ZXN0MDE="
		p = "b2xwYXNzMDE="
	case 18832:
		u = "b2x0ZXN0MDI="
		p = "b2xwYXNzMDI="
	case 18833:
		u = "b2x0ZXN0MDM="
		p = "b2xwYXNzMDM="
	default:
		log.Fatal("Invalid", "port", port)
	}
	//todo: getCredential from database which should be randomly generated when register or import if user already has bitcoin node
	usrBytes, err := base64.StdEncoding.DecodeString(u)
	if err != nil {
		log.Error(err.Error())
		usr = ""
	}
	usr = string(usrBytes)
	passBytes, err := base64.StdEncoding.DecodeString(p)
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

func GetRawAddress(client *brpc.Bitcoind) *btcutil.AddressPubKeyHash {
	addr, _ := client.GetRawChangeAddress()
	if addr == nil {
		log.Fatal("Missing Address")
	}
	return addr.(*btcutil.AddressPubKeyHash)
}

func GetAmount(value string) btcutil.Amount {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal("failed to decode amount", "err", err, "value", value)
	}

	amount, err := btcutil.NewAmount(number)
	if err != nil {
		log.Fatal("failed to create Bitcoin amount", "err", err, "number", number)
	}
	return amount
}

type HTLContract struct {
	Contract 	[]byte		`json:"contract"`
	ContractTx 	wire.MsgTx	`json:"contractTx"`
}

func (h *HTLContract) ToMessage() []byte {
	msg, err := comm.Serialize(h)
	if err != nil {
	    log.Error("Failed to serialize htlc", "err", err)
    }
    return msg
}

func GetHTLCFromMessage(message []byte) *HTLContract{
	log.Debug("Parse message to BTC HTLC")
	register := &HTLContract{}

	result, err := comm.Deserialize(message, register)
	if err != nil {
		log.Error("Failed parse htlc contract", "err", err)
		return nil
	}
	return  result.(*HTLContract)
}