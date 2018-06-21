/*
	Copyright 2017 - 2018 OneLedger
*/
package bitcoin

import (
	"strconv"
	"testing"
	"time"

	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
	brpc "github.com/Oneledger/protocol/node/chains/bitcoin/rpc"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func XTestGeneration(t *testing.T) {
	log.Info("TESTING THE GENERATION")

	bitcoin := GetBtcClient(18831)

	if bitcoin == nil {
		log.Fatal("Can't Get Client")
	}
	log.Debug("Have a Bitcoin Client", "bitcoin", bitcoin)

	channel := ScheduleBlockGeneration(*bitcoin, 5)
	log.Debug("Gen", "channel", channel)

	time.Sleep(6 * time.Second)

	StopBlockGeneration(channel)
}

func Setup() *brpc.Bitcoind {
	bitcoin := GetBtcClient(18831)

	if bitcoin == nil {
		log.Fatal("Can't Get Client")
	}
	log.Debug("Have a Bitcoin Client", "bitcoin", bitcoin)

	return bitcoin
}

func Generate(bitcoin *brpc.Bitcoind) {
	log.Debug("About to Generate")
	text, err := bitcoin.Generate(5)
	if err != nil {
		log.Fatal("Generate", "err", err)
	}
	log.Debug("Generate", "text", text)
}

func Dump(bitcoin *brpc.Bitcoind) {
	// The last block hash on the longest chain...
	hash, err := bitcoin.GetBestBlockhash()
	if err != nil {
		log.Fatal("GetBestBlockhash", "tues", err, "xxx", bitcoin, "vvvv", bitcoin)
	}
	log.Debug("GetBestBlockhash", "hash", hash)

	// Number of blocks in the chain right now
	count, err := bitcoin.GetBlockCount()
	if err != nil {
		log.Fatal("GetBlockCount", "err", err)
	}
	log.Debug("GetBlockCount", "count", count)

	// All of the hashes
	for i := count - 10; i <= count; i++ {
		hash, err = bitcoin.GetBlockHash(i)
		if err != nil {
			log.Warn("GetBlockHash", "err", err)
		}
		log.Debug("GetBlockHash", "hash", hash)
	}

	results, err := bitcoin.ListAccounts(20)
	if err != nil {
		log.Fatal("ListAccounts", "err", err)
	}
	log.Debug("Accounts", "results", results)
}

func TestClient(t *testing.T) {
	log.Info("TESTING THE CLIENT")

	bitcoin := Setup()
	Dump(bitcoin)
}

// Do both Bob and Alice's side of the Hashed Timelock....
func TestHTLC(t *testing.T) {
	bitcoin := Setup()
	BobSetup(bitcoin)
}

func BobSetup(bitcoin *brpc.Bitcoind) {
	log.Debug("BobSetup", "bitcoin", bitcoin)
	chainParams := &chaincfg.TestNet3Params

	cp2Addr, err := btcutil.DecodeAddress("0x0001", chainParams)
	if err != nil {
		log.Error("failed to decode participant address", "err", err)
	}
	if !cp2Addr.IsForNet(chainParams) {
		log.Error("participant address is not intended for use on", "name", chainParams.Name)
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		log.Error("participant address is not P2PKH")
	}

	amountF64, err := strconv.ParseFloat("", 64)
	if err != nil {
		log.Error("failed to decode amount", "err", err)
	}
	amount, err := btcutil.NewAmount(amountF64)
	if err != nil {
		log.Error("failed to decode amount", "err", err)
	}

	initiate := htlc.NewInitiateCmd(cp2AddrP2PKH, amount, 1000)
	hash, err := initiate.RunCommand(bitcoin)

	log.Debug("Results", "initiateCmd", initiate, "hash", hash, "err", err)
}
