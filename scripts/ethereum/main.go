/*

 */

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/net/context"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	oclient "github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	logger "github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	se "github.com/Oneledger/protocol/service/ethereum"
)

var (
	LockRedeemABI = contract.LockRedeemABI
	TestTokenABI  = contract.ERC20BasicABI
	LockRedeemERCABI =contract.LockRedeemERCABI
	// LockRedeemERC20ABI = contract.ContextABI
	LockRedeemContractAddr = "0x1B56deD528F91B62d4b49eD24d9A6ca5c2b577B8"
	TestTokenContractAddr = "0x76d7e36086de8a8FF194F1F80E1106142D1aB172"
	LockRedeemERC20ContractAddr = "0x62edfD8a5B56aAF449CA6ae7Ab05B39CF58cD903"

	cfg                         = config.DefaultEthConfig()
	log                         = logger.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	UserprivKey                 *ecdsa.PrivateKey
	UserprivKeyRedeem           *ecdsa.PrivateKey

	client                 *ethclient.Client
	contractAbi            abi.ABI
	valuelock              = createValue("10000") // in wei (1 eth)
	valueredeem            = createValue("100")
	valuelockERC20         = createValue("1000000000000000000")
	valueredeemERC20       = createValue("100000000000000000")
	fromAddress            common.Address
	redeemRecipientAddress common.Address

	toAddress               = common.HexToAddress(LockRedeemContractAddr)
	toAddressTestToken      = common.HexToAddress(TestTokenContractAddr)
	toAdddressLockRedeemERC = common.HexToAddress(LockRedeemERC20ContractAddr)
)

func createValue(str string) *big.Int {
	n := new(big.Int)
	n, ok := n.SetString(str, 10)
	if !ok {
		fmt.Println("SetString: error")
		return big.NewInt(0)
	}
	return n
}

func init() {

	UserprivKey, _ = crypto.HexToECDSA("6c24a44424c8182c1e3e995ad3ccfb2797e3f7ca845b99bea8dead7fc9dccd09")

	UserprivKeyRedeem, _ = crypto.HexToECDSA("bdb082c7e42a946c477fa3efee4fb5bdece508b47592d8cb57f5e811cd840a40")

	client, _ = cfg.Client()
	contractAbi, _ = abi.JSON(strings.NewReader(LockRedeemABI))

	publicKey := UserprivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")

	}
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

	publicKeyRedeem := UserprivKeyRedeem.Public()
	publicKeyECDSARedeem, ok := publicKeyRedeem.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")

	}
	redeemRecipientAddress = crypto.PubkeyToAddress(*publicKeyECDSARedeem)
}

func main() {

	//lock()

	//time.Sleep(10 * time.Second)

	//redeem()

	erc20lock()
	time.Sleep(10 * time.Second)
	erc20Redeem()
}

func lock() {
	contractAbi, _ := abi.JSON(strings.NewReader(LockRedeemABI))
	bytesData, err := contractAbi.Pack("lock")
	if err != nil {
		return
	}

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

	tx := types.NewTransaction(nonce, toAddress, valuelock, gasLimit, gasPrice, bytesData)

	//fmt.Println("Transaction Unsigned : ", tx)
	chainID, err := client.ChainID(context.Background())
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
	olt, _ := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")

	accReply := &oclient.ListAccountsReply{}
	err = rpcclient.Call("owner.ListAccounts", struct{}{}, accReply)
	if err != nil {
		fmt.Println("query account failed", err)
		return
	}

	acc := accReply.Accounts[0]

	req := se.OLTLockRequest{
		RawTx:   rawTxBytes,
		Address: acc.Address(),
		Fee:     action.Amount{Currency: olt.Name, Value: *balance.NewAmountFromInt(10000000000)},
		Gas:     400000,
	}
	//
	reply := &se.OLTReply{}
	err = rpcclient.Call("eth.CreateRawExtLock", req, reply)
	signReply := &oclient.SignRawTxResponse{}
	err = rpcclient.Call("owner.SignWithAddress", oclient.SignRawTxRequest{
		RawTx:   reply.RawTX,
		Address: acc.Address(),
	}, signReply)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("after sign call",reply.RawTX)

	bresult := &oclient.BroadcastReply{}
	err = rpcclient.Call("broadcast.TxCommit", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult)

	fmt.Println(hex.EncodeToString(reply.RawTX))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("broadcast result: ", bresult.OK)
	fmt.Println(bresult.Log)
}

func redeem() {

	bytesData, err := contractAbi.Pack("redeem", valueredeem)
	if err != nil {
		fmt.Println(err)
		return
	}

	redeemAddress := redeemRecipientAddress.Bytes()
	nonce, err := client.PendingNonceAt(context.Background(), redeemRecipientAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(6321974) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	auth2 := bind.NewKeyedTransactor(UserprivKeyRedeem)
	auth2.Nonce = big.NewInt(int64(nonce))
	auth2.Value = big.NewInt(0) // in wei
	auth2.GasLimit = gasLimit   // in units
	auth2.GasPrice = gasPrice

	value := big.NewInt(0)
	tx2 := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, bytesData)

	chainID, err := client.ChainID(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), UserprivKeyRedeem)
	if err != nil {

		log.Fatal(err)
	}

	ts2 := types.Transactions{signedTx2}

	rawTxBytes2 := ts2.GetRlp(0)
	txNew2 := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes2, txNew2)

	if err != nil {
		fmt.Println(err)
		return
	}

	_ = client.SendTransaction(context.Background(), signedTx2)
	rpcclient, err := rpc.NewClient("http://localhost:26602")
	if err != nil {
		fmt.Println(err)
		return
	}

	accReply := &oclient.ListAccountsReply{}
	err = rpcclient.Call("owner.ListAccounts", struct{}{}, accReply)
	if err != nil {
		fmt.Println("query account failed", err)
		return
	}

	acc := accReply.Accounts[0]

	result := &oclient.ListCurrenciesReply{}
	err = rpcclient.Call("query.ListCurrencies", struct{}{}, result)
	if err != nil {
		fmt.Println(err)
		return
	}
	olt, _ := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")

	rr := se.RedeemRequest{
		acc.Address(),
		redeemAddress,
		rawTxBytes2,
		action.Amount{Currency: olt.Name, Value: *balance.NewAmountFromInt(10000000000)},
		400000,
	}

	reply := &se.OLTReply{}
	err = rpcclient.Call("eth.CreateRawExtRedeem", rr, reply)

	signReply := &oclient.SignRawTxResponse{}
	err = rpcclient.Call("owner.SignWithAddress", oclient.SignRawTxRequest{
		RawTx:   reply.RawTX,
		Address: acc.Address(),
	}, signReply)
	if err != nil {
		fmt.Println(err)
		return
	}

	bresult2 := &oclient.BroadcastReply{}
	err = rpcclient.Call("broadcast.TxCommit", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("broadcast result: ", bresult2.OK)
}

func erc20lock() {
	tokenAbi, _ := abi.JSON(strings.NewReader(TestTokenABI))
	bytesData, err := tokenAbi.Pack("transfer", toAdddressLockRedeemERC, valuelockERC20)
	if err != nil {
		log.Fatal("unable to pack")
	}

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

	tx := types.NewTransaction(nonce, toAddressTestToken, big.NewInt(0), gasLimit, gasPrice, bytesData)

	chainID, err := client.ChainID(context.Background())
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
	olt, ok := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")
	if !ok {
		fmt.Println(result.Currencies)
		return
	}
	accReply := &oclient.ListAccountsReply{}
	err = rpcclient.Call("owner.ListAccounts", struct{}{}, accReply)
	if err != nil {
		fmt.Println("query account failed", err)
		return
	}

	acc := accReply.Accounts[0]

	req := se.OLTERC20LockRequest{
		RawTx:   rawTxBytes,
		Address: acc.Address(),
		Fee:     action.Amount{Currency: olt.Name, Value: *balance.NewAmountFromInt(10000000000)},
		Gas:     400000,
	}
	reply := &se.OLTReply{}
	err = rpcclient.Call("eth.PrepareOLTERC20Lock", req, reply)
	if err != nil {
		fmt.Println("Error in prepare OLTERCLock ", err)
		return
	}
	signReply := &oclient.SignRawTxResponse{}
	err = rpcclient.Call("owner.SignWithAddress", oclient.SignRawTxRequest{
		RawTx:   reply.RawTX,
		Address: acc.Address(),
	}, signReply)
	if err != nil {
		fmt.Println("Error in signing erc lock ", err)
		return
	}

	bresult := &oclient.BroadcastReply{}
	err = rpcclient.Call("broadcast.TxCommit", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult)

	//fmt.Println(hex.EncodeToString(reply.RawTX))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("broadcast result: ", bresult.OK)
	fmt.Println(bresult.Log)
}

func erc20Redeem()  {
	lockRedeemERCABI, _ := abi.JSON(strings.NewReader(LockRedeemERCABI))
	bytesData, err := lockRedeemERCABI.Pack("redeem",valueredeemERC20,toAddressTestToken)
	if err != nil {
		log.Fatal("unable to pack")
	}
	redeemAddress := redeemRecipientAddress.Bytes()
	nonce, err := client.PendingNonceAt(context.Background(), redeemRecipientAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(6321974) // in units

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	auth2 := bind.NewKeyedTransactor(UserprivKeyRedeem)
	auth2.Nonce = big.NewInt(int64(nonce))
	auth2.Value = big.NewInt(0) // in wei
	auth2.GasLimit = gasLimit   // in units
	auth2.GasPrice = gasPrice

	value := big.NewInt(0)
	tx2 := types.NewTransaction(nonce, toAdddressLockRedeemERC, value, gasLimit, gasPrice, bytesData)

	chainID, err := client.ChainID(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), UserprivKeyRedeem)
	if err != nil {

		log.Fatal(err)
	}

	ts2 := types.Transactions{signedTx2}

	rawTxBytes2 := ts2.GetRlp(0)
	txNew2 := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes2, txNew2)

	if err != nil {
		fmt.Println(err)
		return
	}

	_ = client.SendTransaction(context.Background(), signedTx2)
	rpcclient, err := rpc.NewClient("http://localhost:26602")
	if err != nil {
		fmt.Println(err)
		return
	}

	accReply := &oclient.ListAccountsReply{}
	err = rpcclient.Call("owner.ListAccounts", struct{}{}, accReply)
	if err != nil {
		fmt.Println("query account failed", err)
		return
	}

	acc := accReply.Accounts[0]

	result := &oclient.ListCurrenciesReply{}
	err = rpcclient.Call("query.ListCurrencies", struct{}{}, result)
	if err != nil {
		fmt.Println(err)
		return
	}
	olt, _ := result.Currencies.GetCurrencySet().GetCurrencyByName("OLT")

	rr := se.OLTERC20RedeemRequest{
		acc.Address(),
		redeemAddress,
		rawTxBytes2,
		action.Amount{Currency: olt.Name, Value: *balance.NewAmountFromInt(10000000000)},
		400000,
	}

	reply := &se.OLTReply{}
	err = rpcclient.Call("eth.CreateRawExtERC20Redeem", rr, reply)

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
}




//func mapkey(m map[string]string, value string) (key string, ok bool) {
//	for k, v := range m {
//		if v == value {
//			key = k
//			ok = true
//			return
//		}
//	}
//	return
//}

//methodName := "transfer"
//method,exists:=tokenAbi.Methods[methodName]
//if !exists {
//	fmt.Println("Method '%v' not found in the ABI", methodName)
//	return
//}
//_,ok := mapkey(contract.ERC20BasicFuncSigs,method.Sig())
//if !ok {
//	log.Fatal("Method Signature does not exist")
//}
//
//ss := strings.Split(hex.EncodeToString(rawTxBytes), "a9059cbb")

//tokenAmount,err := hex.DecodeString(ss[1][64:128])
//if err != nil {
//	return nil,err
//}
//receiver := ss[1][24:64]
//fmt.Println(common.HexToAddress(receiver))
//amt := big.NewInt(0).SetBytes(tokenAmount)

//transactionHash := client.SendTransaction(context.Background(), signedTx)

//
//if err != nil {
//	fmt.Println(err)
//	return
//}
//
