package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"

	"strings"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	se "github.com/Oneledger/protocol/service/ethereum"
	oclient "github.com/Oneledger/protocol/client"
	//"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"os"
	//"github.com/Oneledger/protocol/log"
	//"os"
)

func main() {
	//var initialValidators = []string{"0xdE11f49F87A0eF71A805Bb47Ba0473432AA5E07a"}
	var keypath = "/scripts/testeth/keydump"
	var contractaddress = "0x68943A1d6Bab608f47258B8A9bdd9c384c1b56A1"
	const LockRedeemABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"removeValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"epochBlockHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newThresholdProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeRemoveValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalEthBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"}],\"name\":\"getRedeemAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"v\",\"type\":\"address\"}],\"name\":\"proposeAddValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getOLTEthAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"numValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"votingThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addValidatorProposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"voteCount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"recipient_\",\"type\":\"address\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"redeemID_\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"proposeNewThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"initialValidators\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"AddValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recepient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_trafered\",\"type\":\"uint256\"}],\"name\":\"RedeemSuccessful\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator_addresss\",\"type\":\"address\"}],\"name\":\"ValidatorSignedRedeem\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"DeleteValidator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"epochHeight\",\"type\":\"uint256\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount_received\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prevThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newThreshold\",\"type\":\"uint256\"}],\"name\":\"NewThreshold\",\"type\":\"event\"}]"

	config := config.CreateEthConfig("http://127.0.0.1:7545", keypath, contractaddress,LockRedeemABI)

	var log = log.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	UserprivKey, err := crypto.HexToECDSA("ee7af353ce3bc37c01187abae7a2d7d1ca22c2a4b88850fc171b988cad924be0")
	// Account 10 from Ganache
	//access,err := ethereum.NewEthereumChainDriver("/home/tanmay/Codebase/protocol",config,log,privKey)
	// todo Remove private Key parameter
	//if(err!=nil){
	//	fmt.Println(err)
	//}

	contractAbi, _ := abi.JSON(strings.NewReader(LockRedeemABI))
	bytesData, _ := contractAbi.Pack("lock")
	client, err := config.Client()
	if err != nil {

		return
	}

	publicKey := UserprivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")

	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(6721974)              // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {

		log.Fatal(err)
	}
	toAddress := common.HexToAddress(contractaddress)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, bytesData)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), UserprivKey)
	if err != nil {

		log.Fatal(err)
	}
	ts := types.Transactions{signedTx}

	rawTxBytes := ts.GetRlp(0)
	txNew := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes, txNew)

	if err != nil {
		fmt.Println(err)
		//fmt.Println("2")
		return
	}
	fmt.Println(txNew.Value())

	fullnode, err := oclient.NewServiceClient("localhost:26602")
	if err != nil {

	}
	result, err := fullnode.ListCurrencies()
	if err != nil {

	}

	pubKey, privKey, err := keys.NewKeyPairFromTendermint()

	accReply, err := fullnode.AddAccount(oclient.AddAccountRequest{
		Type:       chain.Type(0),
		Name:       "test",
		PublicKey:  &pubKey,
		PrivateKey: &privKey,
	})
	if err != nil {

	}
	acc := accReply.Account

	olt, ok := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")
	req := se.OLTLockRequest{
		RawTx:   rawTxBytes,
		Address: acc.Address(),
		Fee:     action.Amount{olt.Name, *balance.NewAmountFromInt(1000000000)},
		Gas:     40000,
	}
	reply := &se.OLTLockReply{}
	err = fullnode.Client.Call("eth.CreateRawExtLock", req, reply)

	signed, err := acc.Sign(reply.RawTX)

	bresult, err := fullnode.TxCommit(oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signed,
		PublicKey: pubKey,
	})

	fmt.Println("broadcast result: ", bresult)
	//d := time.Now().Add(1000 * time.Millisecond)
	//ctx, cancel := context.WithDeadline(context.Background(), d)
	//defer cancel()
	//err = client.SendTransaction(ctx, txNew)
	//locktx,err := access.LockFromSignedTx(signedTx)
	//if err!=nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(txNew)

	//signtx,err := access.Sign(value, toAddress)
	//if err!=nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(signtx)

}
