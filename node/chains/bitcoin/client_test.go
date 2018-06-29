/*
	Copyright 2017 - 2018 OneLedger
*/
package bitcoin

import (
	"crypto/sha256"
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
	"encoding/hex"
	"crypto/rand"
	"bytes"
)

func SetChain() {
	testnode1 := Setup(1)
	testnode2 := Setup(2)
	//testnode3 := Setup(3)
	Generate(testnode1, 20)
	Generate(testnode2, 101)
	//Generate(testnode3, 100)
	hash, err := testnode1.GetBalance("", 0)
	log.Debug("Balance", "hash", hash, "err", err)
	hash, err = testnode2.GetBalance("", 0)
	log.Debug("Balance", "hash", hash, "err", err)
	//testnode3.GetBalance("",0)
}
func XTestSetChain(t *testing.T) {
	SetChain()
}

func TestSwap(t *testing.T) {
	log.Info("Setup Working Swap Test")
	testnode1 := Setup(1)
	testnode2 := Setup(2)
	testnode3 := Setup(3)

	var secret [32]byte
	_, err := rand.Read(secret[:])
	if err != nil {
		log.Error("failed to get random secret with 32 length", "err", err)
	}
	secretHash := sha256.Sum256(secret[:])

	log.Debug("secret pair", "secret",hex.EncodeToString(secret[:]), "secretHash", hex.EncodeToString(secretHash[:]))
	AliceBobSuccessfulSwap(testnode1, testnode2, testnode3, secret[:], secretHash)

	//TypeAddresses()
}

func XTestDump(t *testing.T) {
	testnode1 := Setup(1)
	testnode2 := Setup(2)
	Dump(testnode1)
	Dump(testnode2)
}

func XTestClient(t *testing.T) {
	log.Info("Client Test")

	testnode1 := Setup(1)
	Dump(testnode1)
}

func XTestBlockGeneration(t *testing.T) {
	log.Info("TESTING THE GENERATION")

	testnode1 := GetBtcClient("127.0.0.1:18831", &chaincfg.RegressionNetParams)

	if testnode1 == nil {
		log.Fatal("Can't Get Client")
	}
	log.Debug("Have a Bitcoin Client", "testnode1", testnode1)

	channel := ScheduleBlockGeneration(*testnode1, 5)
	log.Debug("Gen", "channel", channel)

	time.Sleep(6 * time.Second)

	StopBlockGeneration(channel)
}

func Setup(id int) *brpc.Bitcoind {
	var testnode *brpc.Bitcoind

	switch id {
	case 1:
		testnode = GetBtcClient("127.0.0.1:18831", &chaincfg.RegressionNetParams)
	case 2:
		testnode = GetBtcClient("127.0.0.1:18832", &chaincfg.RegressionNetParams)
	case 3:
		testnode = GetBtcClient("127.0.0.1:18833", &chaincfg.RegressionNetParams)
	default:
		log.Fatal("Invalid", "id", id)
	}

	if testnode == nil {
		log.Fatal("Can't Get Client", "config", chaincfg.RegressionNetParams)
	}
	log.Debug("Have a Bitcoin Client", "testnode", testnode)

	return testnode
}

func Generate(testnode *brpc.Bitcoind, count uint64) {
	log.Debug("About to Generate blocks", "cnt", count)
	text, err := testnode.Generate(count)
	if err != nil {
		log.Fatal("Generate", "err", err)
	}
	log.Debug("Generated", "text", text)
}

func Dump(testnode *brpc.Bitcoind) {
	// The last block hash on the longest chain...
	hash, err := testnode.GetBestBlockhash()
	if err != nil {
		log.Fatal("GetBestBlockhash", "err", err, "testnode", testnode)
	}
	log.Debug("GetBestBlockhash", "hash", hash)

	// Number of blocks in the chain right now
	count, err := testnode.GetBlockCount()
	if err != nil {
		log.Fatal("GetBlockCount", "err", err)
	}
	log.Debug("GetBlockCount", "count", count)

	// All of the hashes
	for i := count - 10; i <= count; i++ {
		hash, err = testnode.GetBlockHash(i)
		if err != nil {
			log.Warn("GetBlockHash", "err", err)
		}
		log.Debug("GetBlockHash", "hash", hash)
	}

	results, err := testnode.ListAccounts(20)
	if err != nil {
		log.Fatal("ListAccounts", "err", err)
	}
	log.Debug("Accounts", "results", results)
	for key, value := range results {
		log.Debug("Account", "key", key, "value", value)
	}
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
func GetRawAddress(testnode *brpc.Bitcoind) *btcutil.AddressPubKeyHash {
	addr, _ := testnode.GetRawChangeAddress()
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

func AliceBobSuccessfulSwap(testnode1 *brpc.Bitcoind, testnode2 *brpc.Bitcoind,
	testnode3 *brpc.Bitcoind, secret []byte, secretHash [32]byte) {

	log.Debug("AliceBobSuccessfulSwap", "testnode1", testnode1, "testnode2", testnode2)

	Generate(testnode3, 6)

	aliceTimeout := time.Now().Add( 2 * time.Minute).Unix()
	bobTimeout := time.Now().Add( 1 * time.Minute).Unix()
	aliceAddress := GetRawAddress(testnode1)
	bobAddress := GetRawAddress(testnode2)

	//testnode1.SendToAddress("oltest01", 234.0, "Fill Account", "Fill Account")
	//testnode1.SendToAddress("oltest02", 30000.0, "Fill Account", "Fill Account")
	//testnode2.SendToAddress("oltest02", 1000.0, "Fill an Account", "Fill the Account")
	//testnode2.SendToAddress("oltest01", 10210.0, "Fill an Account", "Fill the Account")

	Generate(testnode3, 20)

	log.Debug("Addresses", "alice", aliceAddress, "bob", bobAddress)

	amount := GetAmount("0.32822")

	log.Debug("==== ALICE INITIATE COMMAND")
	hash, err := htlc.NewInitiateCmd(bobAddress, amount, aliceTimeout, secretHash).RunCommand(testnode1)
	if err != nil {
		log.Warn("Initiate", "err", err)
	}

	// Not threadsafe...
	aliceContract := copyArray(htlc.LastContract)
	aliceContractTx := copyMsgTx(htlc.LastContractTx)

	time.Sleep(3 * time.Second)
	Generate(testnode3, 10)

	log.Debug("==== BOB AUDIT COMMAND")
	err = htlc.NewAuditContractCmd(aliceContract, aliceContractTx).RunCommand(testnode2)
	if err != nil {
		log.Warn("Audit", "err", err)
	}

	time.Sleep(3 * time.Second)
	Generate(testnode3, 10)

	log.Debug("==== BOB PARTICIPATE COMMAND")
	hash, err = htlc.NewParticipateCmd(aliceAddress, amount*2, secretHash, bobTimeout).RunCommand(testnode2)
	if err != nil {
		log.Warn("Participate", "err", err)
	}

	bobContract := copyArray(htlc.LastContract)
	bobContractTx := copyMsgTx(htlc.LastContractTx)

	log.Debug("==== ALICE AUDIT COMMAND")
	err = htlc.NewAuditContractCmd(bobContract, bobContractTx).RunCommand(testnode1)
	if err != nil {
		log.Warn("Audit", "err", err)
	}

	time.Sleep(5 * time.Second)
	Generate(testnode3, 6)

	log.Debug("==== ALICE REDEEM COMMAND")
	hash, err = htlc.NewRedeemCmd(bobContract, bobContractTx, secret).RunCommand(testnode1)
	if err != nil {
		log.Warn("Redeem", "err", err)
	}

	redemptionContractTx := copyMsgTx(htlc.LastContractTx)

	// TODO: Extract Secret
	log.Debug("==== BOB EXTRACT COMMAND")
	err = htlc.NewExtractSecretCmd(redemptionContractTx, secretHash).RunCommand(testnode2)
	if err != nil {
		log.Warn("Extract", "err", err)
	}
	extractSecret := copyArray(htlc.Secret)
	if !bytes.Equal(extractSecret, secret) {
		log.Warn("Extract Secret doesn't match", "extract", extractSecret, "original", secret)
	}
	log.Debug("Extracted Secret matches", "secret", hex.EncodeToString(extractSecret))

	log.Debug("==== BOB REDEEM COMMAND")
	hash, err = htlc.NewRedeemCmd(aliceContract, aliceContractTx, extractSecret).RunCommand(testnode2)
	if err != nil {
		log.Warn("Redeem", "err", err)
	}

	log.Debug("Results", "hash", hash)
	time.Sleep(3 * time.Second)
	Generate(testnode3, 10)
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
