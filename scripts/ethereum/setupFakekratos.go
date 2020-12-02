package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	logger "github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	newValidatorsDir = "/home/tanmay/Codebase/Test/Kratos2020/"
	Log              = logger.NewDefaultLogger(os.Stdout).WithPrefix("testethMigrate")
	Cfg              = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
)

func main() {
	address := ethDeployKratoscontract()
	fmt.Println(address.String())
}

func ethDeployKratoscontract() common.Address {
	contractaddress := common.Address{}
	var validatorset []common.Address
	for i := 0; i < 4; i++ {
		folder := newValidatorsDir + strconv.Itoa(i) + "-Node/consensus/config/"
		ecdspkbytes, err := ioutil.ReadFile(filepath.Join(folder, "priv_validator_key_ecdsa.json"))
		if err != nil {
			Log.Fatal(err)
			return contractaddress
		}
		ecdsPrivKey, err := base64.StdEncoding.DecodeString(string(ecdspkbytes))
		if err != nil {
			Log.Fatal(err)
			return contractaddress
		}
		pkey, err := keys.GetPrivateKeyFromBytes(ecdsPrivKey[:], keys.SECP256K1)
		if err != nil {
			fmt.Println("Privatekey from String ", err)
			return contractaddress
		}
		privatekey := keys.ETHSECP256K1TOECDSA(pkey.Data)

		publicKey := privatekey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			Log.Fatal("error casting public key to ECDSA")
			return contractaddress
		}
		validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		validatorset = append(validatorset, validatorAddress)
	}
	for _, v := range validatorset {
		fmt.Println(v.String())
	}

	contractaddress, err := deployethcdcontract(validatorset, Cfg.Connection)
	if err != nil {
		Log.Fatal("error deploying smart contract", err)
	}
	return contractaddress
}

func deployethcdcontract(initialValidatorList []common.Address, conn string) (common.Address, error) {
	contractAddress := common.Address{}
	os.Setenv("ETHPKPATH", "/tmp/pkdata")
	b1 := make([]byte, 64)
	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return contractAddress, errors.Wrap(err, "Error Reading File")
	}
	pk, err := f.Read(b1)
	if err != nil {
		return contractAddress, errors.Wrap(err, "Error reading private key")
	}
	pkStr := string(b1[:pk])

	//fmt.Println("Private key used to deploy : ", pkStr)

	privatekey, err := crypto.HexToECDSA(pkStr)

	if err != nil {
		return contractAddress, err
	}
	cli, err := ethclient.Dial(conn)
	if err != nil {
		return contractAddress, err
	}

	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return contractAddress, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasPrice, err := cli.SuggestGasPrice(context.Background())
	if err != nil {
		return contractAddress, err
	}
	gasLimit := uint64(2000000) // in units

	auth := bind.NewKeyedTransactor(privatekey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice
	lock_period := big.NewInt(2500)
	tokenSupplyTestToken := new(big.Int)
	validatorInitialFund := big.NewInt(30000000000000000) //300000000000000000
	tokenSupplyTestToken, ok = tokenSupplyTestToken.SetString("1000000000000000000000", 10)
	if !ok {
		return contractAddress, errors.New("Unable to create total supply for token")
	}

	for _, validator := range initialValidatorList {
		nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return contractAddress, err
		}

		//fmt.Println("nonce", nonce)
		tx := types.NewTransaction(nonce, validator, validatorInitialFund, auth.GasLimit, auth.GasPrice, nil)
		fmt.Println(validator.Hex(), ":", validatorInitialFund, "wei")
		chainId, _ := cli.ChainID(context.Background())
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privatekey)
		if err != nil {
			return contractAddress, errors.Wrap(err, "signing tx")
		}
		err = cli.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return contractAddress, errors.Wrap(err, "sending")
		}
		time.Sleep(1 * time.Second)
	}

	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return contractAddress, err
	}

	auth.Nonce = big.NewInt(int64(nonce))

	num_of_validators := big.NewInt(1)
	contractAddress, _, _, err = contract.DeployLockRedeem(auth, cli, initialValidatorList, lock_period, fromAddress, num_of_validators)
	if err != nil {
		return contractAddress, errors.Wrap(err, "Deployement Eth LockRedeem")
	}
	fmt.Println("Activating contract :", contractAddress.String())
	time.Sleep(10 * time.Second)
	err = activateContract(fromAddress, contractAddress, privatekey, cli)
	if err != nil {
		return contractAddress, errors.Wrap(err, "Unable to activate new contract")
	}
	tokenAddress := common.Address{}
	ercAddress := common.Address{}
	fmt.Printf("LockRedeemContractAddr = \"%v\"\n", contractAddress.Hex())
	fmt.Printf("TestTokenContractAddr = \"%v\"\n", tokenAddress.Hex())
	fmt.Printf("LockRedeemERC20ContractAddr = \"%v\"\n", ercAddress.Hex())
	return contractAddress, nil
}

func activateContract(validatorAddress common.Address, KratosSmartContractAddress common.Address, privatekey *ecdsa.PrivateKey, client *ethclient.Client) error {
	ContractAbi, _ := abi.JSON(strings.NewReader(contract.LockRedeemABI))
	// Fake migration Vote
	bytesData, err := ContractAbi.Pack("MigrateFromOld")
	if err != nil {
		fmt.Println(err)
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), validatorAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(1700000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	value := big.NewInt(0)
	tx2 := types.NewTransaction(nonce, KratosSmartContractAddress, value, gasLimit, gasPrice, bytesData)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return err
	}
	signedTx2, err := types.SignTx(tx2, types.NewEIP155Signer(chainID), privatekey)
	if err != nil {
		return err
	}
	ts2 := types.Transactions{signedTx2}

	rawTxBytes2 := ts2.GetRlp(0)
	txNew2 := &types.Transaction{}
	err = rlp.DecodeBytes(rawTxBytes2, txNew2)

	err = client.SendTransaction(context.Background(), signedTx2)
	if err != nil {
		return err
	}

	// Calling Payable
	nonce, err = client.PendingNonceAt(context.Background(), validatorAddress)
	if err != nil {
		return err
	}
	validatorInitialFund := big.NewInt(10)
	tx := types.NewTransaction(nonce, KratosSmartContractAddress, validatorInitialFund, gasLimit, gasPrice, nil)
	chainId, _ := client.ChainID(context.Background())
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privatekey)
	if err != nil {
		return errors.Wrap(err, "signing tx")
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return errors.Wrap(err, "sending")
	}
	time.Sleep(1 * time.Second)
	return nil
}
