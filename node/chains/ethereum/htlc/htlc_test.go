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
	"github.com/ethereum/go-ethereum/crypto"

)



func TestHtlc(t *testing.T) {

	cli, auth, opt := getClientAndAuth()
	auth.GasLimit = 2000000
	address, tx, htlcontract, err := DeployHtlc(auth, cli, auth.From)
	if err != nil {
		log.Error("Failed to deploy new contract: ", "err", err)
	}
	log.Debug("Contract pending deploy: ", "address", address)
	log.Debug("Transaction waiting to be mined: ", "transaction", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(10 * time.Second) // Allow it to be processed by the local node :P


	balance, err := htlcontract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "err", err)
	}
	log.Debug("balance:", "balance", balance)

	testHtlc_Funds(auth, opt, htlcontract)
	testHtlc_Setup(auth, htlcontract)
	testHtlc_Audit(auth, opt, htlcontract)
	testHtlc_Redeem(auth, opt, htlcontract)
	testHtlc_ExtractMsg(opt, address)
	testHtlc_Refund(auth, htlcontract)
	log.Debug("Contract: ", "address", address)
}

func getClientAndAuth() (*ethclient.Client, *bind.TransactOpts, *bind.CallOpts){
	cli, err := ethclient.Dial("/home/lan/go/test/ethereum/C/geth.ipc")
	if err != nil{
		log.Error("failed to get geth ipc ", "err", err)
	}
	key := `{"address":"8a309f95de0e47edb61de8fa0cf8bdd722271789","crypto":{"cipher":"aes-128-ctr","ciphertext":"81becb7ca37be737af147aa0552b1639b770d76ba98fa82069325fe1ce6e1aa1","cipherparams":{"iv":"5be20f263a46d6cca53cb0ae490245fd"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"12456c9a74778a06449596676cc90f2f046e306b5db74688600c04577529b9c2"},"mac":"6737c9dd93f0abc8e102590984790214c2b9dfc36ea6e2b769e80c19eb22e4e8"},"id":"fbaef12b-a667-4c4e-b4c7-7234ef37cbe9","version":3}`
	auth, err := bind.NewTransactor(strings.NewReader(key), "3456")
	if err != nil {
		log.Error("Failed to create authorized transactor: ", "err", err)
	}
	opt := &bind.CallOpts{Pending: true}
	return cli,auth,opt
}

func testHtlc_Funds(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc) {
	log.Debug("======test funds()======")
	auth.GasLimit = 2000000
	auth.Value = big.NewInt(1000000000000000000)

	log.Debug("auth",  "auth", auth)

	//funds
	tx, err := contract.Funds(auth)
	if err != nil {
		log.Error("fund contract failed")
	}
	log.Debug("Transaction waiting to be mined: ", "transaction", tx.Hash(),"value", tx.Value())

	time.Sleep(10 * time.Second)
	balance, err := contract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "err",  err)
	}
	log.Debug("balance:", "balance", balance)
	auth.Value = big.NewInt(0)
}

func getScrPair() (scrHash [32]byte, scr []byte) {
	s := "testswap"
	scrHash = [32]byte(crypto.Keccak256Hash([]byte(s)))
	scr = []byte(s)
	return
}

func testHtlc_Setup(auth *bind.TransactOpts, contract *Htlc) {

	log.Debug("======test setup()======")

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.Hex())



	scrHash, scr := getScrPair()

	log.Debug("secrets", "scr", scr, "scrHash", scrHash)


	log.Debug("auth",  "auth", auth)

	//setup
	tx, err := contract.Setup(auth, big.NewInt(25*3600), receiver, scrHash)
	if err != nil {
		log.Error("failed to setup: ","err", err)
	}
	log.Debug("Transaction waiting to be mined: ", "transaction", tx.Hash(),"value", tx.Value())

	time.Sleep(15 * time.Second)

	addr, err := contract.Receiver(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("failed to get receiver", "err", err)
	}
	log.Debug("receiver in contract: ", "address", addr)
}

func testHtlc_Redeem(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc) {

	log.Debug("======test redeem()======")

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.String())


	addr, err := contract.Receiver(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("failed to get receiver", "err", err)
	}
	log.Debug("receiver in contract: ", "address", addr)


	_, scr := getScrPair()

	tx, err := contract.Redeem(auth, scr)

	if err != nil {
		log.Error("redeem failed", "err", err)
	}
	time.Sleep(20 * time.Second)
	log.Debug("redeem transaction to be mined", "transaction", tx.Hash(), "value", tx.Value())

	balance, err := contract.Balance(opt)
	if err != nil {
		log.Error("Failed to retrieve balance: ", "err",  err)
	}
	log.Debug("balance:", "balance", balance)

}


func testHtlc_ExtractMsg(opt *bind.CallOpts, address common.Address) {

	log.Debug("======test extract()======")

	cli, err := ethclient.Dial("/home/lan/go/test/ethereum/A/geth.ipc")
	if err != nil {
		log.Error("failed to get geth ipc ", "err", err)
	}

	contract, err := NewHtlc(address, cli)
	if err != nil {
		log.Error("Failed to local the htlc contract at address", "err", err, "address", address)
	}

	balance, err := contract.Balance(opt)
	if err != nil {
		log.Error( "Failed to get balance from contract at address", "err", err, "address", address)
	}
	log.Debug("balance of contract","balance", balance)

	scr, err := contract.ExtractMsg(opt)
	if err != nil {
		log.Error("Failed to get the scr", "err", err)
	}
	log.Debug("secret is ", "scr", scr)

}

func testHtlc_Refund(auth *bind.TransactOpts, contract *Htlc) {

	log.Debug("======test refund()======")
	_, scr := getScrPair()

	tx, err := contract.Refund(auth,scr)
	if err != nil{
		log.Error("Failed to call the refund", )
	}
	log.Debug("refund transaction to be mined", "transaction", tx.Hash(), "value", tx.Value())

}

func testHtlc_Audit(auth *bind.TransactOpts, opt *bind.CallOpts, contract *Htlc) {

	log.Debug("======test audit()======")

	receiver := common.HexToAddress("0xd7858005867c3449f6673a91f6e4f719f10e12e5")
	log.Debug("receiver: ", "address", receiver.String())

	scrHash, _ := getScrPair()
	r, err :=contract.Audit(opt, receiver, big.NewInt(1000000000000000000), scrHash )
	if err != nil {
		log.Error("Failed to audit the contract", "err", err)
	}
	log.Debug("Correct audit result", "r", r)

	var falseScrHash [32]byte
	copy(falseScrHash[:], "falsetest")
	r, err =contract.Audit(opt, receiver,big.NewInt(10), falseScrHash )
	if err != nil {
		log.Error("Failed to audit the contract", "err", err)
	}
	log.Debug("Incorrect audit result", "r", r)
}
