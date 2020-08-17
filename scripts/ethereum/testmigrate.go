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
	lockRedeemABI                 = contract.LockRedeemABI
	lockRedeemKratosABI           = contract.LockRedeemKratosABI
	lockRedeemContractAddr        = "0x99e709597677ea3fb5160e46e8ea3d4989f8dfc0"
	lockRedeem_KratosContractAddr = "0x6Cd9acb9dabB14D94559943879f7f5C24B699973"
	//cfg = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Cfg                        = config.DefaultEthConfig("rinkeby", "de5e96cbb6284d5ea1341bf6cb7fa401")
	Log                        = logger.NewDefaultLogger(os.Stdout).WithPrefix("testeth")
	Client                     *ethclient.Client
	ContractAbi                abi.ABI
	KratosContractAbi          abi.ABI
	OldSmartContractAddress    = common.HexToAddress(lockRedeemContractAddr)
	KratosSmartContractAddress = common.HexToAddress(lockRedeem_KratosContractAddr)
	oldValidatorsDir           = "/home/tanmay/Codebase/Test/Pk-dev3/"
)

func init() {

	Client, _ = ethclient.Dial(Cfg.Connection)
	ContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemABI))
	KratosContractAbi, _ = abi.JSON(strings.NewReader(lockRedeemKratosABI))
}

// Redeem locked if tracker fails . User redeems more funds than he has .

func main() {
	fmt.Println("NEW contract isActive : ", checkIsActive())
	migrateContract()
	time.Sleep(time.Second * 15)
	fmt.Println("NEW contract isActive : ", checkIsActive())
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
		Log.Fatal("Unable to get LOCKREDEEM FUTURE")
	}
	callopts := &ethereum.CallOpts{
		Pending:     true,
		BlockNumber: nil,
		Context:     context.Background(),
	}

	bool, err := ctrct.ActiveStatus(callopts)
	if err != nil {
		Log.Fatal("Unable to check contract status")
		return false
	}
	return bool
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
