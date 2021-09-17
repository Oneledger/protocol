/*
	Copyright 2017-2018 OneLedger

	Cli to init a node (server)
*/
package main

import (
	"context"
	crypto2 "crypto"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	tendermint "github.com/tendermint/tendermint/types"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/data/keys"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/log"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize node (server)",
	RunE:  runInitNode,
}

type InitCmdArguments struct {
	genesis  string
	nodeName string
}

var initCmdArgs = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCmdArgs.nodeName, "node_name", "Node", "Name of the node")
	initCmd.Flags().StringVar(&initCmdArgs.genesis, "genesis", "", "Genesis file to use to generate new node key file")
}

type initContext struct {
	genesis  *config.GenesisDoc
	logger   *log.Logger
	rootDir  string
	nodeName string
}

// Given the path of a genesis file and a specified root directory, initNode creates all the configuration files
// needed to run a fullnode inside that specified directory
func runInitNode(cmd *cobra.Command, _ []string) error {
	setEnvVariables()
	rootDir, err := filepath.Abs(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Invalid root directory specified %s", rootArgs))
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err = os.Mkdir(rootDir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create the specified rootdir at %s", rootDir))
		}
	}

	// Generate new configuration file
	cfg := &config.Server{}
	cfgPath := filepath.Join(rootDir, config.FileName)
	err = cfg.ReadFile(cfgPath)
	if err != nil {
		fmt.Println("failed to read config.toml: ", err)
		fmt.Println("generating default configuration, need manual configure to run the node")
		cfg = config.DefaultServerConfig()
	}

	cfg.Node.NodeName = initCmdArgs.nodeName
	csRoot := consensus.RootDirName
	csConfig := consensus.ConfigDirName
	csData := consensus.DataDirName

	configDir := filepath.Join(rootDir, csRoot, csConfig)
	dataDir := filepath.Join(rootDir, csRoot, csData)
	nodeDataDir := filepath.Join(rootDir, "nodedata")
	err = cfg.SaveFile(cfgPath)
	dirs := []string{configDir, dataDir, nodeDataDir}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, config.DirPerms)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("dir creation failed at %s", dir))
		}
	}

	// Put the genesis file in the right place
	if initCmdArgs.genesis != "" {

		genesisPath, err := filepath.Abs(initCmdArgs.genesis)
		if err != nil {
			return errors.Wrap(err, "invalid genesis file path")
		}
		fmt.Println("verifying genesis file provided")
		genesis, err := tendermint.GenesisDocFromFile(genesisPath)
		if err != nil {
			return err
		}
		err = genesis.SaveAs(filepath.Join(configDir, consensus.GenesisFilename))
		if err != nil {
			return errors.Wrap(err, "Failed to save genesis file")
		}

	} else {
		fmt.Println("No genesis file provided, node is not runnable until genesis file is provided at: ", filepath.Join(configDir, consensus.GenesisFilename))
	}
	err = generatePVKeys(configDir, dataDir)
	if err != nil {
		return errors.Wrap(err, "Failed to Get NodeList")
	}
	return nil
}

type nodeKeys struct {
	Validator_addr   tendermint.Address
	Validator_pubKey crypto2.PublicKey
	Witness_addr     string
	Witness_pubKey   ecdsa.PublicKey
}

func generatePVKeys(configDir string, dataDir string) error {
	_, err := p2p.LoadOrGenNodeKey(filepath.Join(configDir, consensus.NodeKeyFilename))
	if err != nil {
		return errors.Wrap(err, "Failed to generate node key")
	}
	// Make private validator file
	pvFile := privval.LoadOrGenFilePV(filepath.Join(configDir, consensus.PrivValidatorKeyFilename),
		filepath.Join(dataDir, consensus.PrivValidatorStateFilename))
	pvFile.Save()

	//Check if folder already has ECDSA KEY ,if it has use existing else generate new
	ecdsaPk := keys.PrivateKey{}
	ecdsaFile := strings.Replace(consensus.PrivValidatorKeyFilename, ".json", "_ecdsa.json", 1)
	if _, err := os.Stat(filepath.Join(configDir, ecdsaFile)); err == nil {
		ecdspkbytes, err := ioutil.ReadFile(filepath.Join(configDir, ecdsaFile))
		if err != nil {
			return err
		}
		ecdsPrivKey, err := base64.StdEncoding.DecodeString(string(ecdspkbytes))
		if err != nil {
			return err
		}
		ecdsaPk, err = keys.GetPrivateKeyFromBytes(ecdsPrivKey[:], keys.SECP256K1)
		if err != nil {
			return err
		}
	} else {
		ecdsaPrivKey := secp256k1.GenPrivKey()
		ecdsaPrivKeyBytes := base64.StdEncoding.EncodeToString([]byte(ecdsaPrivKey[:]))
		ecdsaPk, err = keys.GetPrivateKeyFromBytes([]byte(ecdsaPrivKey[:]), keys.SECP256K1)
		if err != nil {
			return errors.Wrap(err, "error generating secp256k1 private key")
		}

		f, err := os.Create(filepath.Join(configDir, ecdsaFile))

		if err != nil {
			return errors.Wrap(err, "failed to open file to write validator ecdsa private key")
		}
		noofbytes, err := f.Write([]byte(ecdsaPrivKeyBytes))
		if err != nil && noofbytes != len(ecdsaPrivKeyBytes) {
			return errors.Wrap(err, "failed to write validator ecdsa private key")
		}
		err = f.Close()
		if err != nil && noofbytes != len(ecdsaPrivKeyBytes) {
			return errors.Wrap(err, "failed to save validator ecdsa private key")
		}

	}

	//priv, err := ecdsaPk.GetHandler()
	//if err != nil {
	//	return err
	//}
	//fmt.Println(ecdsaPk.Data)
	//pub, err := priv.PubKey().GetHandler()
	//if err != nil {
	//	return err
	//}
	//witnesspubkey := priv.PubKey()
	//witnessaddr := pub.Address()
	privkey := keys.ETHSECP256K1TOECDSA(ecdsaPk.Data)
	pubkey := privkey.Public()
	witnesspubkey, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("failed to cast pubkey")
	}
	witnessaddr := crypto.PubkeyToAddress(*witnesspubkey)

	node :=
		nodeKeys{
			Validator_addr:   pvFile.GetAddress(),
			Validator_pubKey: pvFile.GetPubKey(),
			Witness_addr:     witnessaddr.Hex(),
			Witness_pubKey:   *witnesspubkey,
		}
	fmt.Printf("%+v\n", node)
	return nil
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func getEthUrl(ethUrlArg string) (string, error) {
	u, err := url.Parse(ethUrlArg)
	if err != nil {
		return "", err
	}
	//&& !strings.Contains(u.Path, os.Getenv("API_KEY"))
	if strings.Contains(u.Host, "infura") {
		u.Path = u.Path + "/" + os.Getenv("API_KEY")
		return u.String(), nil
	}
	return ethUrlArg, nil
}

func getEthOpt(conn string, nodeList []node) (*ethchain.ChainDriverOption, error) {

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

	if err != nil {
		return nil, err
	}
	cli, err := ethclient.Dial(conn)
	if err != nil {
		return nil, err
	}

	publicKey := privatekey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	gasPrice, err := cli.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	gasLimit := uint64(6721974) // in units

	auth := bind.NewKeyedTransactor(privatekey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice

	initialValidatorList := make([]common.Address, 0, 10)
	lock_period := big.NewInt(25)

	tokenSupplyTestToken := new(big.Int)
	validatorInitialFund := big.NewInt(300000000000000000)
	tokenSupplyTestToken, ok = tokenSupplyTestToken.SetString("1000000000000000000000", 10)
	if !ok {
		return nil, errors.New("Unabe to create total supplu for token")
	}
	if !ok {
		return nil, errors.New("Unable to create wallet transfer amount")
	}
	for _, node := range nodeList {
		if node.isValidator {
			privkey := keys.ETHSECP256K1TOECDSA(node.esdcaPk.Data)
			nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
			if err != nil {
				return nil, err
			}
			pubkey := privkey.Public()
			ecdsapubkey, ok := pubkey.(*ecdsa.PublicKey)
			if !ok {
				return nil, errors.New("failed to cast pubkey")
			}
			addr := crypto.PubkeyToAddress(*ecdsapubkey)

			initialValidatorList = append(initialValidatorList, addr)
			tx := types.NewTransaction(nonce, addr, validatorInitialFund, auth.GasLimit, auth.GasPrice, nil)
			fmt.Println("validator Address :", addr.Hex(), ":", validatorInitialFund)
			chainId, _ := cli.ChainID(context.Background())
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privatekey)
			if err != nil {
				return nil, errors.Wrap(err, "signing tx")
			}
			err = cli.SendTransaction(context.Background(), signedTx)
			if err != nil {
				return nil, errors.Wrap(err, "sending")
			}
			time.Sleep(1 * time.Second)
		}
	}

	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	oldaddress := common.Address{}
	num_of_validators := big.NewInt(8)
	address, _, _, err := contract.DeployLockRedeem(auth, cli, initialValidatorList, lock_period, oldaddress, num_of_validators)
	if err != nil {
		return nil, errors.Wrap(err, "Deployement Eth LockRedeem")
	}
	tokenAddress := common.Address{}
	ercAddress := common.Address{}

	/*auth.Nonce = big.NewInt(int64(nonce + 1))
	tokenAddress, _, _, err := ethcontracts.DeployERC20Basic(auth, cli, tokenSupplyTestToken)
	if err != nil {
		return nil, errors.Wrap(err, "Deployement Test Token")
	}
	auth.Nonce = big.NewInt(int64(nonce + 2))
	ercAddress, _, _, err := ethcontracts.DeployLockRedeemERC(auth, cli, initialValidatorList)
	if err != nil {
		return nil, errors.Wrap(err, "Deployement ERC LockRedeem")
	}*/

	fmt.Printf("LockRedeemContractAddr = \"%v\"\n", address.Hex())
	fmt.Printf("TestTokenContractAddr = \"%v\"\n", tokenAddress.Hex())
	fmt.Printf("LockRedeemERC20ContractAddr = \"%v\"\n", ercAddress.Hex())
	//tokenAbiMap := make(map[*common.Address]string)
	//tokenAbiMap[&tokenAddress] = contract.ERC20BasicABI
	return &ethchain.ChainDriverOption{
		ContractABI:        contract.LockRedeemABI,
		ERCContractABI:     "",
		TokenList:          []ethchain.ERC20Token{},
		ContractAddress:    address,
		ERCContractAddress: ercAddress,
		TotalSupply:        totalETHSupply,
		TotalSupplyAddr:    lockBalanceAddress,
		BlockConfirmation:  ethBlockConfirmation,
	}, nil
}

func getBtcOpt() bitcoin.ChainDriverOption {
	return bitcoin.ChainDriverOption{
		"testnet3",
		totalBTCSupply,
		lockBalanceAddress,
		btcBlockConfirmation,
	}
}
