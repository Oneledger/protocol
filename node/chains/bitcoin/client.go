/*
	Copyright 2017 - 2018 OneLedger
*/

package bitcoin

import (
	"bytes"
	"github.com/Oneledger/protocol/node/data"
	"time"

	brpc "github.com/Oneledger/protocol/node/chains/bitcoin/rpc"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"encoding/base64"
	"net"
	"strings"

	"strconv"

	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func init() {
	var hash chainhash.Hash
	serial.Register(&hash)
	serial.Register(wire.OutPoint{})
	serial.Register(wire.TxIn{})
	serial.Register(wire.TxOut{})
	serial.Register(wire.TxWitness{})
	serial.Register(wire.MsgTx{})
	serial.Register(HTLContract{})
}

func GetChaincfg() *chaincfg.Params {

	return &chaincfg.RegressionNetParams
}

func GetBtcClient(address string) *brpc.Bitcoind {
	chainParams := GetChaincfg()

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
		log.Error("Can't get the btc rpc client at given address", "status", err)
		return nil
	}

	return cli
}

func getCredential(port int) (usr string, pass string) {

	var u, p string
	switch port {
	case 18831:
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
		log.Fatal("failed to decode amount", "status", err, "value", value)
	}

	amount, err := btcutil.NewAmount(number)
	if err != nil {
		log.Fatal("failed to create Bitcoin amount", "status", err, "number", number)
	}
	return amount
}

type HTLContract struct {
	Contract   []byte `json:"contract"`
	ContractTx []byte `json:"contractTx"`
}

func (h *HTLContract) Chain() data.ChainType {
	return data.BITCOIN
}

func (h *HTLContract) ToBytes() []byte {
	msg, err := serial.Serialize(h, serial.JSON)
	if err != nil {
		log.Error("Failed to serialize htlc", "status", err)
	}
	return msg
}

func (h *HTLContract) ToKey() []byte {
	key, err := btcutil.NewAddressScriptHash(h.Contract, GetChaincfg())
	if err != nil {
		log.Error("Failed to get the key for contract", "status", err)
		return nil
	}
	return key.ScriptAddress()
}

func (h HTLContract) GetMsgTx() *wire.MsgTx {
	if h.ContractTx == nil {
		log.Error("GetMsgTx contractTx nil", "contract", h)
		return nil
	}

	var output wire.MsgTx
	if err := output.Deserialize(bytes.NewReader(h.ContractTx)); err != nil {
		log.Error("GetMsgTx", "err", err)
		return nil
	}
	return &output
}

func (h *HTLContract) FromMsgTx(contract []byte, contractTx *wire.MsgTx) {
	h.Contract = contract
	var contractBuf bytes.Buffer
	contractBuf.Grow(contractTx.SerializeSize())
	contractTx.Serialize(&contractBuf)
	h.ContractTx = contractBuf.Bytes()
	return
}

func (h *HTLContract) FromBytes(message []byte) {
	_, err := serial.Deserialize(message, h, serial.JSON)
	if err != nil {
		log.Error("Failed to deserialize htlc", "err", err)
	}

	return
}
