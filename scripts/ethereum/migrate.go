/*


transaction cost 	52004 gas
execution cost 	81135 gas


133139-116162

transaction cost 	62027 gas
 execution cost 	54155 gas

*/

package main

import (
	"bufio"
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
	oldValidatorsDir              = "/home/tanmay/Codebase/Test/Testing Migrate/devnetOld/"
	newValidatorsDir              = "/home/tanmay/Codebase/Test/Testing Migrate/devnetNew/"
	lockRedeemContractAddr        = "0x1016b294a58a2d9cc2bac156fe5f77e54da00cd2"
	lockRedeem_KratosContractAddr = ""                      //Generated while migrating
	numofValidatorsOld            = big.NewInt(8)           // Kainos Validators
	lock_period                   = big.NewInt(500)         // Approx 2 hours
	gasPriceM                     = big.NewInt(18000000000) // Not Used currently multiplying suggested gas price
	gasLimitM                     = uint64(1700000)
	validatorInitialFund          = big.NewInt(1000000000000000000) // 1 Ether
	Cfg                           = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")

	lockRedeemABI              = contract.LockRedeemABI
	lockRedeemKratosABI        = contract.LockRedeemKratosABI
	maxSecondsToWaitForTX      = 180
	Log                        = logger.NewDefaultLogger(os.Stdout).WithPrefix("ethMigrate")
	Client                     *ethclient.Client
	ContractAbi                abi.ABI
	KratosContractAbi          abi.ABI
	OldSmartContractAddress    = common.HexToAddress(lockRedeemContractAddr)
	KratosSmartContractAddress = common.HexToAddress(lockRedeem_KratosContractAddr)

	Deployer         = ""
	DeployersAddress = common.Address{}
)

func init() {

	Client, _ = ethclient.Dial(Cfg.Connection)
	ContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemABI))
	KratosContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemKratosABI))
}

// Redeem locked if tracker fails . User redeems more funds than he has .

func main() {
	if lock_period.Int64() == 0 {
		Log.Info("Lock Period is zero")
		return
	}
	gasPrice, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		Log.Fatal(err)
	}
	gasPriceM = big.NewInt(0).Add(gasPrice, big.NewInt(0).Div(gasPrice, big.NewInt(2)))
	fmt.Printf("New Validator PrivateKeys at location : %s \nValidators Fund amount : %s \n", newValidatorsDir, validatorInitialFund)
	fmt.Printf("New Block Period  : %d \nNumber of Validators in Kainos : %d \nKainos Smart Comtract Address :%s  \n", lock_period.Int64(), numofValidatorsOld, OldSmartContractAddress.String())

	fmt.Printf("Gaslimit  : %d \nGasPrice : %d \n", gasLimitM, gasPriceM)
	fmt.Printf("Press the Enter Key to continue")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	address := ethDeployKratoscontract()
	KratosSmartContractAddress = address
	fmt.Printf("Press the Enter Key to continue Migrate From : %s To %s \n", OldSmartContractAddress.String(), KratosSmartContractAddress.String())
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println("New contract isActive : ", checkIsActive())
	fmt.Printf("Press the Enter Key to continue , Old Validator PrivateKeys at location : %s", oldValidatorsDir)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	migrateContract()
	time.Sleep(time.Second * 15)
	fmt.Println("New contract isActive/ Wait for tx to be confirmed on etheruem  if not active : ", checkIsActive())
	fmt.Printf("Press Enter to trasfer fund from old Validator :")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	takeFunds()
}

// Stop Kainos network
// Dump genesis ( Old genesis has 8 Validator address)
// Keep a copy of the PrivateKeys_ecdsa for old validators.
// Edit the Genesis  (Manual)
// Add new validator address  -> 2 ? ( Add public key and address to genesis , use privatekey_ecdsa for  sc deployment)
// Use PrivateKey_ecdsa of these new validators to deploy LockRedeem_Kratos (New Validators become witnesses in Kratos)
// Add smart contract address to genesis
// Witness list does not change with stake / unstake
// Run migrate script  (Separate ,not part of olclient )
// Inputs  ->  old privatekeys-ecdsa , old smart contract address ,new smart contract address
// Outcome -> Balance transferred to new smart contract
// OLD smart contract ACTIVE = False   , NEW smart contract ACTIVE = True

// Test
// Lock some ether on current devnet
// Do the steps mentioned above
// Redeem from new network which has the new smart contract address

func checkIsActive() bool {
	fmt.Println("Kratos Smart contract Address: ", KratosSmartContractAddress.String())
	ctrct, err := contract.NewLockRedeemKratos(KratosSmartContractAddress, Client)
	if err != nil {
		Log.Fatal("Unable to get LockRedeemKratos :", err)
	}
	callopts := &ethereum.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		Context:     context.Background(),
	}

	ok, err := ctrct.ActiveStatus(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract status :", err)
		return false
	}
	return ok
}

func migrateContract() {
	for i := 0; i < int(numofValidatorsOld.Int64()); i++ {

		bytesData, err := ContractAbi.Pack("migrate", KratosSmartContractAddress)
		if err != nil {
			fmt.Println(err)
			return
		}
		folder := oldValidatorsDir + strconv.Itoa(i) + "-Node/consensus/config/"
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

		gasLimit := gasLimitM // in units
		gasPrice := gasPriceM
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
	fmt.Println("Kratos Funding Initial Validators : ")
	err, contractaddress := deployethcdcontract(validatorset)
	if err != nil {
		Log.Fatal("error deploying smart contract", err)
	}
	return contractaddress
}

func deployethcdcontract(initialValidators []common.Address) (error, common.Address) {
	address := common.Address{}
	err := os.Setenv("ETHPKPATH", "/tmp/pkdata")
	if err != nil {
		Log.Error("Unable to set ETHPKPATH")
		return err, address
	}
	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return errors.Wrap(err, "Error Reading File"), address
	}
	b1 := make([]byte, 64)
	pk, err := f.Read(b1)
	if err != nil {
		return errors.Wrap(err, "Error reading private key"), address
	}
	//fmt.Println("Private key used to deploy : ", string(b1[:pk]))
	pkStr := string(b1[:pk])
	privatekey, err := crypto.HexToECDSA(pkStr)

	if err != nil {
		return err, address
	}
	cli, err := ethclient.Dial(Cfg.Connection)
	if err != nil {
		return err, address
	}

	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return err, address
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	DeployersAddress = fromAddress
	fmt.Println("Deployers address ", fromAddress.String())
	auth := bind.NewKeyedTransactor(privatekey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimitM  // in units
	auth.GasPrice = gasPriceM

	initialValidatorList := make([]common.Address, 0, 10)
	validatorInitialFund := validatorInitialFund //300000000000000000

	for _, valAddr := range initialValidators {

		nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return err, address
		}

		addr := valAddr

		initialValidatorList = append(initialValidatorList, addr)
		tx := types.NewTransaction(nonce, addr, validatorInitialFund, gasLimitM, gasPriceM, nil)
		fmt.Println(addr.Hex(), ":", validatorInitialFund, "wei")
		chainId, _ := cli.ChainID(context.Background())
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privatekey)
		if err != nil {
			return errors.Wrap(err, "signing tx"), address
		}
		err = cli.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return errors.Wrap(err, "sending"), address
		}
		time.Sleep(1 * time.Second)
	}

	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err, address
	}

	auth.Nonce = big.NewInt(int64(nonce))

	address, tx, _, err := contract.DeployLockRedeemKratos(auth, cli, initialValidatorList, lock_period, OldSmartContractAddress, numofValidatorsOld)
	if err != nil {
		return errors.Wrap(err, "Deployement Eth LockRedeem Kratos"), address
	}
	Log.Info("Confirming Smart contract Deployment on the Network :", address.String())
	ok = CheckTXstatus(tx.Hash())
	if !ok {
		return errors.New("Transaction Could not be Confirmed on the main-net"), address
	}
	fmt.Printf("LockRedeemContractKratosAddr = \"%v\"\n", address.Hex())

	return nil, address
}

func CheckTXstatus(txHash common.Hash) bool {
	for counter := 0; counter < maxSecondsToWaitForTX; counter++ {
		result, err := Client.TransactionReceipt(context.Background(), txHash)
		if err == nil && result.Status == types.ReceiptStatusSuccessful {
			return true
		}
		time.Sleep(time.Second * 1)
	}
	return false
}

func takeFunds() {
	for i := 0; i < int(numofValidatorsOld.Int64()); i++ {
		folder := oldValidatorsDir + strconv.Itoa(i) + "-Node/consensus/config/"
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

		publicKey := privatekey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			Log.Fatal("error casting public key to ECDSA")
			return
		}
		validatorAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		nonce, err := Client.PendingNonceAt(context.Background(), validatorAddress)
		if err != nil {
			Log.Fatal(err)
		}

		gasLimit := int64(gasLimitM) // in units
		gasPrice, err := Client.SuggestGasPrice(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		gasCost := gasPrice.Mul(gasPrice, big.NewInt(gasLimit))
		////spareWei := big.NewInt(1000000000000000)
		currentBalance, err := Client.BalanceAt(context.Background(), validatorAddress, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		currentBalance.Sub(currentBalance, gasCost)
		g, err := Client.SuggestGasPrice(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		tx := types.NewTransaction(nonce, DeployersAddress, currentBalance, uint64(gasLimit), g, nil)
		chainID, err := Client.ChainID(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privatekey)
		if err != nil {
			Log.Fatal(err)
			return
		}
		err = Client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			Log.Fatal(err, validatorAddress.String())
			return
		}
		fmt.Println("Funds transferred from :", validatorAddress.Hex(), "Amount Transfered :", currentBalance)
	}
}
