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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
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
	lockRedeemContractAddr = "0xdCF4cA39890B73b243850d87317fb5d3D31985cd"
	newSmartcontract       = "0x7fef1061Ea2789E04b125A52f1DEa2fF9Eb9e198"
	//cfg = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Cfg                     = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Log                     = logger.NewDefaultLogger(os.Stdout).WithPrefix("Migrate")
	Client                  *ethclient.Client
	ContractAbi             abi.ABI
	FutureContractAbi       abi.ABI
	OldSmartContractAddress = common.HexToAddress(lockRedeemContractAddr)
	NEWsmartContract        = common.HexToAddress(newSmartcontract)
	readDir                 = os.Getenv("OLDATA") + "/devnet/"
	lock_period             = new(big.Int)
	auth                    *bind.TransactOpts
)

//auth                    *bind.TransactOpts
//)

func init() {
	gasLimit := uint64(2000000) // in units
	lock_period = big.NewInt(25)
	gasPrice := big.NewInt(1000000000)

	os.Setenv("ETHPKPATH", "/tmp/pkdata")
	Client, _ = ethclient.Dial(Cfg.Connection)
	ContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemABI))
	FutureContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemFutureABI))

	privatekey, err := readPkFile()
	if err != nil {
		Log.Fatal("Unable to read private key from ", os.Getenv("ETHPKPATH"))
		return
	}

	auth = bind.NewKeyedTransactor(privatekey)
	auth.GasPrice = gasPrice
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units

	auth.Context = context.Background()

}

func main() {
	//_, err := deployNewContract()
	//if err != nil {
	//	return
	//}
	verifyValidators()
	//NEWsmartContract = newcontract
	checkContract()
	migrateContract()
	time.Sleep(time.Second * 15)
	checkContract()
	verifyValidators()

}

// Redeem locked if tracker fails . User redeems more funds than he has .
func deployNewContract() (common.Address, error) {

	addressfuture, _, _, err := contract.DeployLockRedeemV2(auth, Client, lock_period, OldSmartContractAddress, big.NewInt(4))
	if err != nil {
		Log.Fatal("Unable to Deploy LockRedeemV2", err)
		return common.Address{}, err
	}
	fmt.Println("LockRedeemV2 : ", addressfuture.Hex())
	return addressfuture, nil
}
func readPkFile() (*ecdsa.PrivateKey, error) {
	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return nil, errors.Wrap(err, "Error Reading File")
	}
	b1 := make([]byte, 64)
	pk, err := f.Read(b1)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private key")
	}
	//fmt.Println("Private key used to deploy : ", string(b1[:pk]))
	pkStr := string(b1[:pk])
	privatekey, err := crypto.HexToECDSA(pkStr)
	return privatekey, nil
}

func checkContract() {
	ctrct, err := contract.NewLockRedeemV2(NEWsmartContract, Client)
	if err != nil {
		Log.Fatal("Unable to get LOCKREDEEM FUTURE")
	}
	callopts := &ethereum.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		Context:     context.Background(),
	}
	balance, err := ctrct.GetTotalEthBalance(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract status :", err)
		return
	}
	sign, err := ctrct.GetMigrationCount(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract signatures :", err)
		return
	}
	fmt.Println("Migration Signatures : ", sign, " | ", "Contract balance : ", balance)
	return
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
		publicKeyECDSA, ok := publicKeyRedeem.(*ecdsa.PublicKey)
		if !ok {
			Log.Fatal("error casting public key to ECDSA")
			return
		}
		validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
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

func verifyValidators() {
	for i := 0; i < 4; i++ {
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

		publicKeyValidator := privatekey.Public()
		publicKeyECDSA, ok := publicKeyValidator.(*ecdsa.PublicKey)
		if !ok {
			Log.Fatal("error casting public key to ECDSA")
			return
		}
		validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		//redeemAddress := redeemRecipientAddress.Bytes()
		callopts := &ethereum.CallOpts{
			Pending:     false,
			From:        validatorAddress,
			BlockNumber: nil,
			Context:     context.Background(),
		}
		ctrct, err := contract.NewLockRedeemV2(NEWsmartContract, Client)
		if err != nil {
			Log.Fatal("Unable to get LOCKREDEEM FUTURE")
		}
		bool, err := ctrct.VerifyValidator(callopts)
		Log.Info("Verifying Validator : ", validatorAddress.Hex(), " : ", bool)
	}
}
