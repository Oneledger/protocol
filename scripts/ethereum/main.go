/*


transaction cost 	52004 gas
execution cost 	81135 gas


133139-116162

transaction cost 	62027 gas
 execution cost 	54155 gas

*/

package main

import (
	"crypto/ecdsa"
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
	LockRedeemABI    = contract.LockRedeemABI
	TestTokenABI     = contract.ERC20BasicABI
	LockRedeemERCABI = contract.LockRedeemERCABI
	// LockRedeemERC20ABI = contract.ContextABI
	LockRedeemContractAddr      = "0x8bA980D33E8325AD815B03a4F9f642A6ae862F0d"
	TestTokenContractAddr       = "0xe2d8814DECcE7DF7c4E78bF61b8FaF5DE706A707"
	LockRedeemERC20ContractAddr = "0x061d8824a0553F34cf1cf48baB54bd585fF300d9"

	cfg               = config.DefautEthConfigRinkeby()
	log               = logger.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	UserprivKey       *ecdsa.PrivateKey
	UserprivKeyRedeem *ecdsa.PrivateKey
	spamKey           *ecdsa.PrivateKey

	client                 *ethclient.Client
	contractAbi            abi.ABI
	valuelock              = createValue("1000") // in wei (1 eth)
	valueredeem            = createValue("10")
	valuelockERC20         = createValue("1000000000000000000")
	valueredeemERC20       = createValue("100000000000000000")
	fromAddress            common.Address
	redeemRecipientAddress common.Address
	spamAddress            common.Address

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
	privKey := "6c24a44424c8182c1e3e995ad3ccfb2797e3f7ca845b99bea8dead7fc9dccd09"
	if strings.Contains(cfg.Connection, "rinkeby") {
		privKey = "02038529C9AB706E9F4136F4A4EB51E866DBFE22D5E102FD3A22C14236E1C2EA"
	}
	UserprivKey, _ = crypto.HexToECDSA(privKey)
	//UserprivKey, _ = crypto.HexToECDSA("02038529C9AB706E9F4136F4A4EB51E866DBFE22D5E102FD3A22C14236E1C2EA")

	UserprivKeyRedeem, _ = crypto.HexToECDSA(privKey)

	spamKey, _ = crypto.HexToECDSA(privKey)

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

	spampub := spamKey.Public()
	spamecdsapub, ok := spampub.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")

	}
	spamAddress = crypto.PubkeyToAddress(*spamecdsapub)
}

func main() {

	lock()
	//time.Sleep(time.Second * 5)
	//send12trasactions()
	time.Sleep(1 * time.Minute)
	redeem()
	//time.Sleep(5 * time.Second)
	//send12trasactions()
	//redeem()

	//erc20lock()
	///time.Sleep(10 * time.Second)
	//erc20Redeem()
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

	rpcclient, err := rpc.NewClient("http://localhost:26602") //104.196.191.206:26604
	//rpcclient, err := rpc.NewClient("https://fullnode-sdk.devnet.oneledger.network/")
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
	err = rpcclient.Call("broadcast.TxSync", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Lock broadcast result: ", bresult.OK)
	if !bresult.OK {
		fmt.Println(bresult.Log)
	}
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

	value := big.NewInt(10000000000000000)
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
	err = client.SendTransaction(context.Background(), signedTx2)
	if err != nil {
		log.Fatal(err)
	}

	//time.Sleep(time.Second * 15)

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
	err = rpcclient.Call("broadcast.TxSync", oclient.BroadcastRequest{
		RawTx:     reply.RawTX,
		Signature: signReply.Signature.Signed,
		PublicKey: signReply.Signature.Signer,
	}, bresult2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Redeem broadcast result: ", bresult2.OK)
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

func erc20Redeem() {
	lockRedeemERCABI, _ := abi.JSON(strings.NewReader(LockRedeemERCABI))
	bytesData, err := lockRedeemERCABI.Pack("redeem", valueredeemERC20, toAddressTestToken)
	if err != nil {
		log.Fatal("unable to pack")
	}
	bytesDataExecuteRedeemd, err := lockRedeemERCABI.Pack("executeredeem", toAddressTestToken)
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
	Executetx := types.NewTransaction(nonce+1, toAdddressLockRedeemERC, value, gasLimit, gasPrice, bytesDataExecuteRedeemd)
	chainID, err := client.ChainID(context.Background())
	if err != nil {

		log.Fatal(err)
	}

	signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), UserprivKeyRedeem)
	if err != nil {

		log.Fatal(err)
	}
	signedExecuteRedeem, err := types.SignTx(Executetx, types.NewEIP155Signer(chainID), UserprivKeyRedeem)
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
	time.Sleep(30 * time.Second)
	err = client.SendTransaction(context.Background(), signedExecuteRedeem)
	if err == nil {
		fmt.Println("Executed Redeem  for user : ", redeemRecipientAddress)
	}

}

func send12trasactions() {
	for i := 0; i <= 26; i++ {
		contractAbi, _ := abi.JSON(strings.NewReader(LockRedeemABI))
		bytesData, err := contractAbi.Pack("lock")
		if err != nil {
			return
		}

		//redeemAddress := redeemRecipientAddress.Bytes()
		nonce, err := client.PendingNonceAt(context.Background(), spamAddress)
		if err != nil {
			log.Fatal(err)
		}

		gasLimit := uint64(6321974) // in units

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		auth2 := bind.NewKeyedTransactor(spamKey)
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

		signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), spamKey)
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
	}
	fmt.Println("12 transactions sent")
}
