/*
	Copyright 2017 - 2018 OneLedger
*/
package bitcoin

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
	brpc "github.com/Oneledger/protocol/node/chains/bitcoin/rpc"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func TestSwap(t *testing.T) {
	log.Info("Successful Swap Test")
	testnet := Setup()
	AliceBobSuccessfulSwap(testnet)

	//TypeAddresses()
}

func XTestClient(t *testing.T) {
	log.Info("Client Test")

	testnet := Setup()
	Dump(testnet)
}

func XTestBlockGeneration(t *testing.T) {
	log.Info("TESTING THE GENERATION")

	testnet := GetBtcClient("127.0.0.1:18831")

	if testnet == nil {
		log.Fatal("Can't Get Client")
	}
	log.Debug("Have a Bitcoin Client", "testnet", testnet)

	channel := ScheduleBlockGeneration(*testnet, 5)
	log.Debug("Gen", "channel", channel)

	time.Sleep(6 * time.Second)

	StopBlockGeneration(channel)
}

func Setup() *brpc.Bitcoind {
	testnet := GetBtcClient("127.0.0.1:18831")

	if testnet == nil {
		log.Fatal("Can't Get Client")
	}
	log.Debug("Have a Bitcoin Client", "testnet", testnet)

	return testnet
}

func Generate(testnet *brpc.Bitcoind) {
	log.Debug("About to Generate")
	text, err := testnet.Generate(5)
	if err != nil {
		log.Fatal("Generate", "err", err)
	}
	log.Debug("Generate", "text", text)
}

func Dump(testnet *brpc.Bitcoind) {
	// The last block hash on the longest chain...
	hash, err := testnet.GetBestBlockhash()
	if err != nil {
		log.Fatal("GetBestBlockhash", "tues", err, "xxx", testnet, "vvvv", testnet)
	}
	log.Debug("GetBestBlockhash", "hash", hash)

	// Number of blocks in the chain right now
	count, err := testnet.GetBlockCount()
	if err != nil {
		log.Fatal("GetBlockCount", "err", err)
	}
	log.Debug("GetBlockCount", "count", count)

	// All of the hashes
	for i := count - 10; i <= count; i++ {
		hash, err = testnet.GetBlockHash(i)
		if err != nil {
			log.Warn("GetBlockHash", "err", err)
		}
		log.Debug("GetBlockHash", "hash", hash)
	}

	results, err := testnet.ListAccounts(20)
	if err != nil {
		log.Fatal("ListAccounts", "err", err)
	}
	log.Debug("Accounts", "results", results)
}

var addresses []string = []string{
	"2NAUVavVAkoFqWocNj7rSArVMV5awNB34Jc",
	"2NGYwMKYCMrNQP15c15iJG2VCgMYPi5uJXL",
	"a914ffa47e020c83e2feb6fe41a4c3178aeb29705e7187",
	"00142d85e5eb13acfd6f54d514a8d290174d5eb723b8",
	"031c7c3268741c5000a66b2f03a74a49204a6ed30b22255a508db463dbe294c05f",
	"2d85e5eb13acfd6f54d514a8d290174d5eb723b8",
	"031c7c3268741c5000a66b2f03a74a49204a6ed30b22255a508db463dbe294c05f",
	"bcrt1q9kz7t6cn4n7k74x4zj5d9yqhf40twgacvg3ev3",
	"00142d85e5eb13acfd6f54d514a8d290174d5eb723b8",
	"bcrt1q9kz7t6cn4n7k74x4zj5d9yqhf40twgacvg3ev3",
	"9d80ec70bf937f6f333279dda2ba89aec26214cf",
}

func TypeAddresses() {
	chainParams := &chaincfg.RegressionNetParams
	for i := 0; i < len(addresses); i++ {
		log.Debug("Testing", "address", addresses[i])
		pubkey, err := btcutil.DecodeAddress(addresses[i], chainParams)
		if err != nil {
			log.Warn("Bad value", "err", err)
		} else {
			log.Debug("PublicKey struct", "pubkey", pubkey, "type", reflect.TypeOf(pubkey))
		}
	}
}

// bitcoin-cli -regtest -rpcuser=oltest01 -rpcpassword=olpass01  -rpcport=18831 getrawchangeaddress
func GetRawAddress(testnet *brpc.Bitcoind) *btcutil.AddressPubKeyHash {
	addr, _ := testnet.GetRawChangeAddress()
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

func AliceBobSuccessfulSwap(testnet *brpc.Bitcoind) {
	log.Debug("AliceBobSuccessfulSwap", "testnet", testnet)

	timeout := int64(1000)
	secret := []byte("This is a secret")

	aliceAddress := GetRawAddress(testnet)
	bobAddress := GetRawAddress(testnet)

	/*
		aliceAddress := GetTestAddress()
		bobAddress := GetTestAddress()
	*/

	amount := GetAmount("1.3282384902")

	_, err := htlc.NewInitiateCmd(aliceAddress, amount, timeout).RunCommand(testnet)
	if err != nil {
		log.Fatal("Initiate", "err", err)
	}
	// Not threadsafe...
	contract := htlc.LastContract
	contractTx := htlc.LastContractTx

	hash, err := htlc.NewParticipateCmd(bobAddress, amount, secret, timeout).RunCommand(testnet)
	if err != nil {
		log.Fatal("Participate", "err", err)
	}

	hash, err = htlc.NewRedeemCmd(contract, contractTx, secret).RunCommand(testnet)
	if err != nil {
		log.Fatal("Redeem", "err", err)
	}

	log.Debug("Results", "hash", hash)
}

func GetContractTx(hash *chainhash.Hash) *wire.MsgTx {
	return nil
}

// TODO: Not working
func GetTestAddress() *btcutil.AddressPubKeyHash {

	chainParams := &chaincfg.RegressionNetParams

	// bitcoin-cli -regtest -rpcuser=oltest01 -rpcpassword=olpass01  -rpcport=18831 getnewaddress
	// bitcoin-cli -regtest -rpcuser=oltest01 -rpcpassword=olpass01  -rpcport=18831 validateaddress 2MvQKdE4pkgkuACicfoFznmuss6G4PVBhrP

	pubkey, _ := btcutil.DecodeAddress("2Mv5xcF9yaUKTZHjbJUJBV4ayUS7Xozfz6X", chainParams)

	log.Debug("PublicKey struct", "pubkey", pubkey, "type", reflect.TypeOf(pubkey))

	if !pubkey.IsForNet(chainParams) {
		log.Warn("participant address is not intended for use on", "name", chainParams.Name)
	}

	//stringAddress := pubkey.EncodeAddress()

	//cp2Addr, err := btcutil.NewAddressPubKeyHash([]byte(stringAddress), chainParams)
	//cp2Addr := &btcutil.AddressPubKeyHash{hash: []byte(stringAddress)}

	//cp2Addr, err := btcutil.DecodeAddress("0x0001", chainParams)
	//if err != nil {
	//	log.Warn("failed to decode participant address", "err", err, "addr", stringAddress)
	//}
	/*
		if cp2Addr == nil {
			log.Warn("Missing address")
		}
	*/

	//cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	//cp2AddrP2PKH := cp2Addr
	cp2AddrP2PKH := pubkey.(*btcutil.AddressPubKeyHash)

	/*
		ok := false
		if !ok {
			log.Warn("participant address is not P2PKH", "err", ok, "type",
				reflect.TypeOf(cp2Addr), "cp2Addr", cp2Addr)
		}
	*/

	return cp2AddrP2PKH
}
