/*

 */

package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/net/context"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	oclient "github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	se "github.com/Oneledger/protocol/service/ethereum"
)

func main() {
	// var initialValidators = []string{"0xdE11f49F87A0eF71A805Bb47Ba0473432AA5E07a"}
	//
	const LockRedeemABI = contract.LockRedeemABI
	contractAddr := "0x4B1BBb85929333569b7c109A815b8e50fEad8096"
	cfg := config.EthereumChainDriverConfig{Connection: "http://localhost:7545"}

	var log = log.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	UserprivKey, err := crypto.HexToECDSA("247eb7922a29ee18b1e2877b5859e0932d03ceae941c3f60047647637580bd17")
	// Account 10 from Ganache

	//
	contractAbi, _ := abi.JSON(strings.NewReader(LockRedeemABI))
	bytesData, err := contractAbi.Pack("lock")
	if err != nil {
		return
	}
	client, err := cfg.Client()
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {

		log.Fatal(err)
	}
	gasLimit := uint64(6721974) // in units

	auth := bind.NewKeyedTransactor(UserprivKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	value := big.NewInt(1000000000000000000) // in wei (1 eth)

	toAddress := common.HexToAddress(contractAddr)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, bytesData)
	fmt.Println("Trasaction Unsigned : ", tx)
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

	rpcclient, err := rpc.NewClient("http://localhost:26602")
	if err != nil {
		fmt.Println(err)
		return
	}
	result := &oclient.ListCurrenciesReply{}
	err = rpcclient.Call("query.ListCurrencies", struct{}{}, result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
	accReply := &oclient.ListAccountsReply{}
	err = rpcclient.Call("owner.ListAccounts", struct{}{}, accReply)
	if err != nil {
		fmt.Println("query account failed", err)
		return
	}

	acc := accReply.Accounts[0]

	olt, ok := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")
	req := se.OLTLockRequest{
		RawTx:   rawTxBytes,
		Address: acc.Address(),
		Fee:     action.Amount{Currency: olt.Name, Value: *balance.NewAmountFromInt(10000000000)},
		Gas:     400000,
	}
	//
	fmt.Println("REQUEST  : ", req)
	reply := &se.OLTLockReply{}
	err = rpcclient.Call("eth.CreateRawExtLock", req, reply)

	fmt.Println("REPLY    : ", reply, err)
	signReply := &oclient.SignRawTxResponse{}
	err = rpcclient.Call("owner.SignWithAddress", oclient.SignRawTxRequest{
		RawTx:   reply.RawTX,
		Address: acc.Address(),
	}, signReply)
	if err != nil {
		fmt.Println(err)
		return
	}

	bresult := &oclient.BroadcastReply{}
	err = rpcclient.Call("broadcast.TxCommit", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("broadcast result: ", bresult.OK)
	//}

}
