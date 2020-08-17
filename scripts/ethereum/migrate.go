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
	oldValidatorsDir              = ""
	newValidatorsDir              = ""
	lockRedeemContractAddr        = ""
	lockRedeem_KratosContractAddr = ""
	numofValidatorsOld            = big.NewInt(8)
	lock_period                   = big.NewInt(0)
	gasPriceM                     = big.NewInt(18000000000) // Not Used currently multiplying suggested gas price
	gasLimitM                     = uint64(700000)
	Cfg                           = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")

	lockRedeemABI       = contract.LockRedeemABI
	lockRedeemKratosABI = contract.LockRedeemKratosABI

	Log                        = logger.NewDefaultLogger(os.Stdout).WithPrefix("testethMigrate")
	Client                     *ethclient.Client
	ContractAbi                abi.ABI
	KratosContractAbi          abi.ABI
	OldSmartContractAddress    = common.HexToAddress(lockRedeemContractAddr)
	KratosSmartContractAddress = common.HexToAddress(lockRedeem_KratosContractAddr)
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
	ethDeployKratoscontract()
	//fmt.Println("NEW contract isActive : ", checkIsActive())
	//migrateContract()
	//time.Sleep(time.Second * 15)
	//fmt.Println("NEW contract isActive : ", checkIsActive())
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
	ctrct, err := contract.NewLockRedeemKratos(KratosSmartContractAddress, Client)
	if err != nil {
		Log.Fatal("Unable to get LOCKREDEEM Kratos")
	}
	callopts := &ethereum.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		Context:     context.Background(),
	}

	ok, err := ctrct.ActiveStatus(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract status")
		return false
	}
	return ok
}

func migrateContract() {
	for i := 0; i < 4; i++ {

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
		gasPrice, err := Client.SuggestGasPrice(context.Background())
		if err != nil {
			Log.Fatal(err)
			return
		}
		gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(0).Div(gasPrice, big.NewInt(2)))
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

func ethDeployKratoscontract() {
	var validatorset []common.Address
	for i := 0; i < 4; i++ {
		folder := newValidatorsDir + strconv.Itoa(i) + "-Node/consensus/config/"
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
		validatorset = append(validatorset, validatorAddress)
	}
	for _, v := range validatorset {
		fmt.Println(v.String())
	}
	err := deployethcdcontract(validatorset)
	if err != nil {
		fmt.Println(err)
	}

}

func deployethcdcontract(initialValidators []common.Address) error {
	err := os.Setenv("ETHPKPATH", "/tmp/pkdata")
	if err != nil {
		Log.Error("Unable to set ETHPKPATH")
		return err
	}
	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return errors.Wrap(err, "Error Reading File")
	}
	b1 := make([]byte, 64)
	pk, err := f.Read(b1)
	if err != nil {
		return errors.Wrap(err, "Error reading private key")
	}
	//fmt.Println("Private key used to deploy : ", string(b1[:pk]))
	pkStr := string(b1[:pk])
	privatekey, err := crypto.HexToECDSA(pkStr)

	if err != nil {
		return err
	}
	cli, err := ethclient.Dial(Cfg.Connection)
	if err != nil {
		return err
	}

	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasLimit := uint64(7000000) // in units

	auth := bind.NewKeyedTransactor(privatekey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	gasPrice, err := Client.SuggestGasPrice(context.Background())
	if err != nil {
		Log.Fatal(err)
	}
	gasPrice = big.NewInt(0).Add(gasPrice, big.NewInt(0).Div(gasPrice, big.NewInt(2)))
	auth.GasPrice = gasPrice
	initialValidatorList := make([]common.Address, 0, 10)

	tokenSupplyTestToken := new(big.Int)
	validatorInitialFund := big.NewInt(30000000000000000) //300000000000000000

	tokenSupplyTestToken, ok = tokenSupplyTestToken.SetString("1000000000000000000000", 10)
	if !ok {
		return errors.New("Unabe to create total supply for token")
	}

	for _, valAddr := range initialValidators {

		nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return err
		}

		addr := valAddr

		initialValidatorList = append(initialValidatorList, addr)
		tx := types.NewTransaction(nonce, addr, validatorInitialFund, auth.GasLimit, auth.GasPrice, nil)
		fmt.Println(addr.Hex(), ":", validatorInitialFund, "wei")
		chainId, _ := cli.ChainID(context.Background())
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privatekey)
		if err != nil {
			return errors.Wrap(err, "signing tx")
		}
		err = cli.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return errors.Wrap(err, "sending")
		}
		time.Sleep(1 * time.Second)
	}

	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))

	address, _, _, err := contract.DeployLockRedeemKratos(auth, cli, initialValidatorList, lock_period, OldSmartContractAddress, numofValidatorsOld)
	if err != nil {
		return errors.Wrap(err, "Deployement Eth LockRedeem Kratos")
	}
	tokenAddress := common.Address{}
	ercAddress := common.Address{}

	fmt.Printf("LockRedeemContractKratosAddr = \"%v\"\n", address.Hex())
	fmt.Printf("TestTokenContractAddr = \"%v\"\n", tokenAddress.Hex())
	fmt.Printf("LockRedeemERC20ContractAddr = \"%v\"\n", ercAddress.Hex())
	return nil
}
