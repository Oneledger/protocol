package htlc

import (
	"testing"
	"github.com/Oneledger/protocol/node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"strings"

	"math/big"
	"github.com/ethereum/go-ethereum/common"

	"crypto/sha256"
)



func TestHtlc(t *testing.T) {

	cli, auth, opt := getClientAndAuth()
	auth.GasLimit = 2000000
	address, tx, htlcontract, err := DeployHtlc(auth, cli, auth.From)
	if err != nil {
		log.Error("Failed to deploy new contract: ", "status", err)
	}
	log.Debug("Contract pending deploy: ", "address", address)
	log.Debug("Transaction waiting to be mined: ", "transaction", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(10 * time.Second) // Allow it to be processed by the local node :P


	balance, err := htlcontract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "status", err)
	}
	log.Debug("balance:", "balance", balance)
	value := new(big.Int)
	value.SetString("100000000000000000000", 10)
	testHtlc_Funds(auth, opt, htlcontract, value)
	//testHtlc_Setup(auth, htlcontract)
	testHtlc_Audit(auth, opt, htlcontract, value)
	//testHtlc_Redeem(auth, opt, htlcontract)
	testHtlc_ExtractMsg(opt, address)
	testHtlc_Refund(auth, htlcontract)
	log.Debug("Contract: ", "address", address)
}

func getClientAndAuth() (*ethclient.Client, *bind.TransactOpts, *bind.CallOpts){
	cli, err := ethclient.Dial("/home/lan/go/test/ethereum/B/geth.ipc")
	if err != nil{
		log.Error("failed to get geth ipc ", "status", err)
	}
	key := `{"address":"aafa2d8980a730b02195f9c8dfeafeb3e69a69ca","crypto":{"cipher":"aes-128-ctr","ciphertext":"cfc36b7deb503116482371b7d2596aa936758b8247279efce461cf0344ae4b31","cipherparams":{"iv":"fc200b937116258856dd0e5a085e011d"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7354c4523dfc70372c8c34616c15dc21448ac40617ffc3a7b3a9af7ee32c37e6"},"mac":"1b322fa3c5789cede87144783f2bd8c4588e5094e68f7880640ca9a5458b8aab"},"id":"87363a39-0171-4640-ba12-b5aacad7aed2","version":3}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "2345")
	if err != nil {
		log.Error("Failed to create authorized transactor: ", "status", err)
	}
	opt := &bind.CallOpts{Pending: true}
	return cli,auth,opt
}

func testHtlc_Funds(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc, v *big.Int) {
	log.Debug("======test funds()======")
	auth.GasLimit = 2000000
	auth.Value = v

	log.Debug("auth",  "auth", auth)

	//funds

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.Hex())

	scrHash, scr := getScrPair()

	log.Debug("secrets", "scr", scr, "scrHash", scrHash)

	tx, err := contract.Funds(auth, big.NewInt(25*3600), receiver, scrHash)
	if err != nil {
		log.Error("fund contract failed")
	}
	log.Debug("Transaction waiting to be mined: ", "transaction", tx.Hash(),"value", tx.Value())

	time.Sleep(10 * time.Second)
	balance, err := contract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "status",  err)
	}
	log.Debug("balance:", "balance", balance)
	auth.Value = big.NewInt(0)

	addr, err := contract.Receiver(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("failed to get receiver", "status", err)
	}
	log.Debug("receiver in contract: ", "address", addr)

	sh, err := contract.ScrHash(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("can't get the secret hash", "status", err)
	}
	log.Debug("scrHash in contract", "scrhash",sh)
}

func getScrPair() (scrHash [32]byte, scr []byte) {
	s := "testswap"
	scrHash = [32]byte(sha256.Sum256([]byte(s)))
	scr = []byte(s)
	return
}

func testHtlc_Redeem(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc) {

	log.Debug("======test redeem()======")

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.String())

    balance, err := contract.Balance(opt)
    if err != nil {
        log.Error("Failed to retrieve balance: ", "status",  err)
    }
    log.Debug("balance before redeem", "balance", balance)

	addr, err := contract.Receiver(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("failed to get receiver", "status", err)
	}
	log.Debug("receiver in contract: ", "address", addr)


	_, scr := getScrPair()

	tx, err := contract.Redeem(auth, scr)

	if err != nil {
		log.Error("redeem failed", "status", err)
	}
	time.Sleep(20 * time.Second)
	log.Debug("redeem transaction to be mined", "transaction", tx.Hash(), "value", tx.Value())

	balance, err = contract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "status",  err)
	}
	log.Debug("balance after redeem", "balance", balance)

}


func testHtlc_ExtractMsg(opt *bind.CallOpts, address common.Address) {

	log.Debug("======test extract()======")

	cli, err := ethclient.Dial("/home/lan/go/test/ethereum/A/geth.ipc")
	if err != nil {
		log.Error("failed to get geth ipc ", "status", err)
	}

	contract, err := NewHtlc(address, cli)
	if err != nil {
		log.Error("Failed to local the htlc contract at address", "status", err, "address", address)
	}

	balance, err := contract.Balance(opt)
	if err != nil {
		log.Error( "Failed to get balance from contract at address", "status", err, "address", address)
	}
	log.Debug("balance of contract","balance", balance)

	scr, err := contract.ExtractMsg(opt)
	if err != nil {
		log.Error("Failed to get the scr", "status", err)
	}
	log.Debug("secret is ", "scr", scr)

}

func testHtlc_Refund(auth *bind.TransactOpts, contract *Htlc) {

	log.Debug("======test refund()======")

	balance, err := contract.Balance(&bind.CallOpts{Pending: false})

	log.Debug("balance before refund", "balance", balance)

	tx, err := contract.Refund(auth)
	if err != nil{
		log.Error("Failed to call the refund", )
	}
	log.Debug("refund transaction to be mined", "transaction", tx.Hash(), "value", tx.Value())
	time.Sleep(20 * time.Second)
    balance, err = contract.Balance(&bind.CallOpts{Pending: false})
    log.Debug("balance after refund", "balance", balance)
}

func testHtlc_Audit(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc, v *big.Int) {

	log.Debug("======test audit()======")

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.String())

	scrHash, _ := getScrPair()
	r, err :=contract.Audit(opt, receiver, v, scrHash )
	if err != nil {
		log.Error("Failed to audit the contract", "status", err)
	}
	log.Debug("Correct audit result", "r", r)

	var falseScrHash [32]byte
	copy(falseScrHash[:], "falsetest")
	r, err =contract.Audit(opt, receiver, big.NewInt(10), falseScrHash )
	if err != nil && r == false {
		log.Debug("Failed to audit the contract as expected because used false secret", "r", r)
	} else {
		log.Error("something wrong with audit", "status", err, "r", r)
	}

}
