/*
	Copyright 2017-2018 OneLedger

	Cli to init a node (server)
*/
package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"

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
	genesis       string
	nodeName      string
	numValidators int
	numFullnodes  int
	// Total amount of funds to be shared across each node
	totalFunds           int64
	initialTokenHolders  []string
	chainID              string
	ethUrl               string
	deploySmartcontracts bool
	cloud                bool
}

var initCmdArgs = &InitCmdArguments{}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initCmdArgs.nodeName, "node_name", "Node", "Name of the node")
	initCmd.Flags().StringVar(&initCmdArgs.genesis, "genesis", "", "Genesis file to use to generate new node Key file")
	initCmd.Flags().IntVar(&initCmdArgs.numValidators, "validators", 4, "Number of validators to initialize mainnetnet with")
	initCmd.Flags().IntVar(&initCmdArgs.numFullnodes, "fullnodes", 1, "Number of fullnodes to initialize mainnetnet with")
	initCmd.Flags().Int64Var(&initCmdArgs.totalFunds, "total_funds", 1000000000, "The total amount of tokens in circulation")
	initCmd.Flags().StringSliceVar(&initCmdArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")
	initCmd.Flags().StringVar(&initCmdArgs.chainID, "chain_id", "", "Specify a chain ID, a random one is generated if not given")
	initCmd.Flags().StringVar(&initCmdArgs.ethUrl, "eth_rpc", "HTTP://127.0.0.1:7545", "URL for ethereum network")
	initCmd.Flags().BoolVar(&initCmdArgs.deploySmartcontracts, "deploy_smart_contracts", true, "deploy eth contracts")
	initCmd.Flags().BoolVar(&initCmdArgs.cloud, "cloud_deploy", false, "set true for deploying on cloud")
}

type initContext struct {
	genesis  *consensus.GenesisDoc
	logger   *log.Logger
	rootDir  string
	nodeName string
}

// Given the path of a genesis file and a specified root directory, initNode creates all the configuration files
// needed to run a fullnode inside that specified directory
func runInitNode(cmd *cobra.Command, _ []string) error {
	setEnvVariablesGanache()
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
	cfg := config.DefaultServerConfig()
	cfg.Node.NodeName = initCmdArgs.nodeName
	cfgPath := filepath.Join(rootDir, config.FileName)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create config file at %s", cfgPath))
	}
	csRoot := consensus.RootDirName
	csConfig := consensus.ConfigDirName
	csData := consensus.DataDirName

	configDir := filepath.Join(rootDir, csRoot, csConfig)
	dataDir := filepath.Join(rootDir, csRoot, csData)
	nodeDataDir := filepath.Join(rootDir, "nodedata")

	// Put the genesis file in the right place
	if initCmdArgs.genesis != "" {
		err = cfg.SaveFile(cfgPath)
		dirs := []string{configDir, dataDir, nodeDataDir}
		for _, dir := range dirs {
			err = os.MkdirAll(dir, config.DirPerms)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Dir creation failed at %s", dir))
			}
		}

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
		// Make node Key

	} else {
		//fmt.Println("No genesis file provided, node is not runnable until genesis file is provided at: ", filepath.Join(configDir, consensus.GenesisFilename))
		fmt.Println("Genarating Genesis file  : ")
	}
	nodeList, _, err := generatePVKeys(rootDir)
	if err != nil {
		return errors.Wrap(err, "Failed to Get NodeList")
	}
	cdo := &ethchain.ChainDriverOption{}
	url, err := getEthUrl(initCmdArgs.ethUrl)
	if err != nil {
		return err
	}
	fmt.Println("Deployment Network :", url)
	fmt.Println("Deploy Smart contracts : ", initCmdArgs.deploySmartcontracts)
	if initCmdArgs.deploySmartcontracts {
		if len(initCmdArgs.ethUrl) > 0 {
			cdo, err = getEthOpt(url, nodeList)
			if err != nil {
				return errors.Wrap(err, "failed to deploy the initial eth contract")
			}
		}
	}

	cdoBytes, err := json.Marshal(cdo)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filepath.Join(rootDir, "cdOpts.json"), cdoBytes, os.ModePerm)

	//err = genesisDoc.SaveAs(filepath.Join(rootDir, "genesis.json"))
	//if err != nil {
	//	return err
	//}
	return nil
}

func generatePVKeys(rootDir string) ([]node, []consensus.GenesisValidator, error) {
	totalNodes := initCmdArgs.numValidators + initCmdArgs.numFullnodes
	nodeList := make([]node, totalNodes)
	validatorList := make([]consensus.GenesisValidator, initCmdArgs.numValidators)
	for i := 0; i < totalNodes; i++ {
		// Make node Key
		nodename := initCmdArgs.nodeName + strconv.Itoa(i)
		folder := filepath.Join(rootDir, nodename)
		err := os.MkdirAll(folder, config.DirPerms)
		if err != nil {
			return nil, nil, err
		}

		isValidator := i < initCmdArgs.numValidators
		nodekey, err := p2p.LoadOrGenNodeKey(filepath.Join(folder, consensus.NodeKeyFilename))
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to generate node Key")
		}
		// Make private Validator file
		pvFile := privval.GenFilePV(filepath.Join(folder, consensus.PrivValidatorKeyFilename),
			filepath.Join(folder, consensus.PrivValidatorStateFilename))
		pvFile.Save()

		ecdsaPrivKey := secp256k1.GenPrivKey()
		ecdsaPrivKeyBytes := base64.StdEncoding.EncodeToString([]byte(ecdsaPrivKey[:]))
		ecdsaPk, err := keys.GetPrivateKeyFromBytes([]byte(ecdsaPrivKey[:]), keys.SECP256K1)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error generating secp256k1 private Key")
		}
		fmt.Println("ecdsaPK :", nodekey)
		ecdsaFile := strings.Replace(consensus.PrivValidatorKeyFilename, ".json", "_ecdsa.json", 1)
		f, err := os.Create(filepath.Join(folder, ecdsaFile))

		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to open file to write Validator ecdsa private Key")
		}
		noofbytes, err := f.Write([]byte(ecdsaPrivKeyBytes))
		if err != nil && noofbytes != len(ecdsaPrivKeyBytes) {
			return nil, nil, errors.Wrap(err, "failed to write Validator ecdsa private Key")
		}
		err = f.Close()
		if err != nil && noofbytes != len(ecdsaPrivKeyBytes) {
			return nil, nil, errors.Wrap(err, "failed to save Validator ecdsa private Key")
		}
		n := node{IsValidator: isValidator, Key: nodekey, EsdcaPk: ecdsaPk}
		if isValidator {
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodename,
				Power:   1,
			}
			n.Validator = validator
			validatorList[i] = validator
		}
		nodeList[i] = n

	}
	//jsonData, err := json.Marshal(persistentPeers)
	//if err != nil {
	//	return nil, nil, errors.Wrap(err, "Error in marshalling nodeList to Json")
	//}
	//ioutil.WriteFile("persistantpeers.json", jsonData, 0600)
	return nodeList, validatorList, nil
}

func getEthOpt(conn string, nodeList []node) (*ethchain.ChainDriverOption, error) {

	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return nil, errors.Wrap(err, "Error Reading File")
	}
	b1 := make([]byte, 64)
	pk, err := f.Read(b1)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private Key")
	}
	//fmt.Println("Private Key used to deploy : ", string(b1[:pk]))
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
		privkey := keys.ETHSECP256K1TOECDSA(node.EsdcaPk.Data)
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
		if node.Validator.Address.String() == "" {
			continue
		}

		initialValidatorList = append(initialValidatorList, addr)
		tx := types.NewTransaction(nonce, addr, validatorInitialFund, auth.GasLimit, auth.GasPrice, (nil))
		fmt.Println("Validator Address :", addr.Hex(), ":", validatorInitialFund)
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

	nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))

	address, _, _, err := contract.DeployLockRedeem(auth, cli, initialValidatorList, lock_period)
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
		ContractABI:    contract.LockRedeemABI,
		ERCContractABI: contract.LockRedeemERCABI,
		TokenList: []ethchain.ERC20Token{{
			TokName:        "TTC",
			TokAddr:        tokenAddress,
			TokAbi:         contract.ERC20BasicABI,
			TokTotalSupply: totalTTCSupply,
		}},
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

func getEthUrl(ethUrlArg string) (string, error) {

	u, err := url.Parse(ethUrlArg)
	if err != nil {
		return "", err
	}
	if strings.Contains(u.Host, "infura") && !strings.Contains(u.Path, os.Getenv("API_KEY")) {
		setEnvVariablesInfura()
		u.Path = u.Path + "/" + os.Getenv("API_KEY")
		return u.String(), nil
	}
	return ethUrlArg, nil
}

func getInitialState(args *InitCmdArguments, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
	btcOption bitcoin.ChainDriverOption) consensus.AppState {
	olt := balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"}
	vt := balance.Currency{Id: 1, Name: "VT", Chain: chain.ONELEDGER, Unit: "vt"}
	obtc := balance.Currency{Id: 2, Name: "BTC", Chain: chain.BITCOIN, Decimal: 8, Unit: "satoshi"}
	oeth := balance.Currency{Id: 3, Name: "ETH", Chain: chain.ETHEREUM, Decimal: 18, Unit: "wei"}
	ottc := balance.Currency{Id: 4, Name: "TTC", Chain: chain.TESTTOKEN, Decimal: 18, Unit: "testUnits"} //Tokens count by number ,Unit 1
	currencies := []balance.Currency{olt, vt, obtc, oeth, ottc}
	feeOpt := fees.FeeOption{
		FeeCurrency:   olt,
		MinFeeDecimal: 9,
	}
	balances := make([]consensus.BalanceState, 0, len(nodeList))
	staking := make([]consensus.Stake, 0, len(nodeList))
	domains := make([]consensus.DomainState, 0, len(nodeList))
	fees_db := make([]consensus.BalanceState, 0, len(nodeList))
	total := olt.NewCoinFromInt(args.totalFunds)

	var initialAddrs []keys.Address
	initAddrIndex := 0
	for _, addr := range args.initialTokenHolders {
		tmpAddr := keys.Address{}
		err := tmpAddr.UnmarshalText([]byte(addr))
		if err != nil {
			fmt.Println("Error adding initial address:", addr)
			continue
		}
		initialAddrs = append(initialAddrs, tmpAddr)
	}

	for _, node := range nodeList {
		if !node.IsValidator {
			continue
		}

		h, err := node.EsdcaPk.GetHandler()
		if err != nil {
			fmt.Println("err")
			panic(err)
		}

		var stakeAddr keys.Address
		if len(initialAddrs) > 0 {
			if initAddrIndex > (len(initialAddrs) - 1) {
				initAddrIndex = 0
			}
			stakeAddr = initialAddrs[initAddrIndex]
			initAddrIndex++
		} else {
			stakeAddr = node.Key.PubKey().Address().Bytes()
		}

		pubkey, _ := keys.PubKeyFromTendermint(node.Validator.PubKey.Bytes())
		st := consensus.Stake{
			ValidatorAddress: node.Validator.Address.Bytes(),
			StakeAddress:     stakeAddr,
			Pubkey:           pubkey,
			ECDSAPubKey:      h.PubKey(),
			Name:             node.Validator.Name,
			Amount:           *vt.NewCoinFromInt(node.Validator.Power).Amount,
		}
		staking = append(staking, st)
	}

	if len(args.initialTokenHolders) > 0 {
		for _, acct := range initialAddrs {
			share := total.DivideInt64(int64(len(args.initialTokenHolders)))
			balances = append(balances, consensus.BalanceState{
				Address:  acct,
				Currency: olt.Name,
				Amount:   *olt.NewCoinFromAmount(*share.Amount).Amount,
			})
			balances = append(balances, consensus.BalanceState{
				Address:  acct,
				Currency: vt.Name,
				Amount:   *vt.NewCoinFromInt(100).Amount,
			})
		}
	} else {
		for _, node := range nodeList {
			amt := int64(100)
			if !node.IsValidator {
				amt = 1
			}
			share := total.DivideInt64(int64(len(nodeList)))
			balances = append(balances, consensus.BalanceState{
				Address:  node.Key.PubKey().Address().Bytes(),
				Currency: olt.Name,
				Amount:   *olt.NewCoinFromAmount(*share.Amount).Amount,
			})
			balances = append(balances, consensus.BalanceState{
				Address:  node.Key.PubKey().Address().Bytes(),
				Currency: vt.Name,
				Amount:   *vt.NewCoinFromInt(amt).Amount,
			})
		}
	}

	return consensus.AppState{
		Currencies: currencies,
		Balances:   balances,
		Staking:    staking,
		Domains:    domains,
		Fees:       fees_db,
		Governance: consensus.GovernanceState{
			FeeOption:   feeOpt,
			ETHCDOption: option,
			BTCCDOption: btcOption,
			ONSOptions:  onsOption,
		},
	}
}
