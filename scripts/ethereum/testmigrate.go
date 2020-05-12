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
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/net/context"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	logger "github.com/Oneledger/protocol/log"
)

var (
	lockRedeemABI          = contract.LockRedeemABI
	lockRedeemFutureABI    = contract.LockRedeemFutureABI
	lockRedeemContractAddr = "0x3A15Be9D65cF8f6f966672Ff4606cF7D07A651A8"
	newSmartcontract       = "0x4610C5CfF3bDb7C013e4729787149A127B740B03"
	//cfg = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Cfg                     = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Log                     = logger.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	Client                  *ethclient.Client
	ContractAbi             abi.ABI
	FutureContractAbi       abi.ABI
	OldSmartContractAddress = common.HexToAddress(lockRedeemContractAddr)
	NEWsmartContract        = common.HexToAddress(newSmartcontract)
	readDir                 = "/home/tanmay/Codebase/Test/devnet/"
)

func init() {

	Client, _ = ethclient.Dial(Cfg.Connection)
	ContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemABI))
	FutureContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemFutureABI))
}

// Redeem locked if tracker fails . User redeems more funds than he has .

func main() {
	fmt.Println("NEW contract isActive : ", checkIsActive())
	migrateContract()
	time.Sleep(time.Second * 15)
	fmt.Println("NEW contract isActive : ", checkIsActive())
}

func checkIsActive() bool {
	ctrct, err := contract.NewLockRedeemFuture(NEWsmartContract, Client)
	if err != nil {
		Log.Fatal("Unable to get LOCKREDEEM FUTURE")
	}
	callopts := &ethereum.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		Context:     context.Background(),
	}
	bool, err := ctrct.IsActive(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract status")
		return false
	}
	return bool
}

func migrateContract() {
	for i := 0; i < 4; i++ {

		bytesData, err := ContractAbi.Pack("migrate", NEWsmartContract)
		if err != nil {
			fmt.Println(err)
			return
		}
		folder := readDir + strconv.Itoa(i) + "-Node/consensus/config/"
		ecdspkbytes, err := ioutil.ReadFile(filepath.Join(folder, "priv_validator_key_ecdsa.json"))
		if err != nil {
			Log.Fatal(err)
			return
		}
		ecdsPrivKey, err := base64.StdEncoding.DecodeString(string(ecdspkbytes))
		if err != nil {
			Log.Fatal(err)
			return
		}
		pkey, err := keys.GetPrivateKeyFromBytes(ecdsPrivKey[:], keys.SECP256K1)
		if err != nil {
			fmt.Println("Privatekey from String ", err)
			return
		}
		privatekey := keys.ETHSECP256K1TOECDSA(pkey.Data)

		publicKeyRedeem := privatekey.Public()
		publicKeyECDSARedeem, ok := publicKeyRedeem.(*ecdsa.PublicKey)
		if !ok {
			Log.Fatal("error casting public key to ECDSA")
			return
		}
		validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSARedeem)
		//redeemAddress := redeemRecipientAddress.Bytes()
		nonce, err := Client.PendingNonceAt(context.Background(), validatorAddress)
		if err != nil {
			Log.Fatal(err)
		}

		gasLimit := uint64(700000) // in units
		gasPrice, err := Client.SuggestGasPrice(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		value := big.NewInt(0)
		tx2 := types.NewTransaction(nonce, OldSmartContractAddress, value, gasLimit, gasPrice, bytesData)
		chainID, err := Client.ChainID(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), privatekey)
		if err != nil {
			Log.Fatal(err)
			return
		}
		ts2 := types.Transactions{signedTx2}

		rawTxBytes2 := ts2.GetRlp(0)
		txNew2 := &types.Transaction{}
		err = rlp.DecodeBytes(rawTxBytes2, txNew2)

		err = Client.SendTransaction(context.Background(), signedTx2)
		if err != nil {
			Log.Fatal(err)
			return
		}
		Log.Info("Validator Signed  :", validatorAddress.Hex())
	}
}
