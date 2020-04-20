package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/log"
)

type mainnetArgument struct {
	// Number of validators
	numValidators    int
	numNonValidators int
	outputDir        string
	pvkey_Dir        string
	p2pPort          int
	allowSwap        bool
	chainID          string
	dbType           string
	namesPath        string
	createEmptyBlock bool
	// Total amount of funds to be shared across each node
	totalFunds          int64
	initialTokenHolders []string

	ethUrl               string
	deploySmartcontracts bool
	cloud                bool
	loglevel             int
}

var mainnetCmdArgs = &mainnetArgument{}

var mainnetCmd = &cobra.Command{
	Use:   "mainnet",
	Short: "Initializes a genesis file for OneLedger network",
	RunE:  runMainnet,
}

type mainetContext struct {
	logger *log.Logger
	names  []string
}

func init() {
	initCmd.AddCommand(mainnetCmd)
	mainnetCmd.Flags().StringVarP(&mainnetCmdArgs.pvkey_Dir, "pv_dir", "p", "/home/tanmay/Codebase/Test/mainnet", "Directory which contains Genesis File and NodeList")
	mainnetCmd.Flags().StringVarP(&mainnetCmdArgs.outputDir, "output_dir", "o", "./", "Directory to store initialization files for the mainet, default current folder")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.allowSwap, "enable_swaps", false, "Allow swaps")
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.numValidators, "validators", 4, "Number of validators to initialize mainnetnet with")
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.numNonValidators, "nonvalidators", 1, "Number of fullnodes to initialize mainnetnet with")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.createEmptyBlock, "empty_blocks", false, "Allow creating empty blocks")
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.dbType, "db_type", "goleveldb", "Specify the type of DB backend to use: (goleveldb|cleveldb)")
	//mainnetCmd.Flags().StringVar(&mainnetCmdArgs.namesPath, "names", "", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	// 1 billion by default
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.ethUrl, "eth_rpc", "HTTP://127.0.0.1:7545", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.loglevel, "loglevel", 3, "Specify the log level for olfullnode. 0: Fatal, 1: Error, 2: Warning, 3: Info, 4: Debug, 5: Detail")
	mainnetCmd.Flags().Int64Var(&mainnetCmdArgs.totalFunds, "total_funds", 1000000000, "The total amount of tokens in circulation")
	mainnetCmd.Flags().StringSliceVar(&mainnetCmdArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")

}

func newMainetContext(args *mainnetArgument) (*mainetContext, error) {
	logger := log.NewLoggerWithPrefix(os.Stdout, "olfullnode mainnet")
	var names []string
	files, err := ioutil.ReadDir(args.pvkey_Dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			names = append(names, file.Name())
		}
	}
	if len(names) != args.numValidators+args.numNonValidators {
		return nil, errors.New("Not enough Key Pairs present the the directory ")
	}
	return &mainetContext{
		names:  names,
		logger: logger,
	}, nil

}

func runMainnet(_ *cobra.Command, _ []string) error {
	ctx, err := newMainetContext(mainnetCmdArgs)
	if err != nil {
		return err
	}
	if err != nil {
		return errors.Wrap(err, "runMainet failed")
	}
	args := mainnetCmdArgs

	totalNodes := len(ctx.names)

	if args.dbType != "cleveldb" && args.dbType != "goleveldb" {
		ctx.logger.Error("Invalid dbType specified, using goleveldb...", "dbType", args.dbType)
		args.dbType = "goleveldb"
	}

	generatePort := portGenerator(26600)

	persistentPeers := make([]string, totalNodes)
	nodeList := make([]node, totalNodes)
	validatorList := make([]consensus.GenesisValidator, args.numValidators)

	// Create the GenesisValidator list and its Key files priv_validator_key.json and node_key.json
	for i, nodeName := range ctx.names {
		isValidator := i < args.numValidators
		readDir := filepath.Join(args.pvkey_Dir, nodeName)
		nodeDir := filepath.Join(args.outputDir, nodeName)
		configDir := filepath.Join(nodeDir, "consensus", "config")
		dataDir := filepath.Join(nodeDir, "consensus", "data")
		nodeDataDir := filepath.Join(nodeDir, "nodedata")
		createDirectories(configDir, dataDir, nodeDataDir)
		ecdspkbytes, err := ioutil.ReadFile(filepath.Join(readDir, "priv_validator_key_ecdsa.json"))
		if err != nil {
			return err
		}
		ecdsPrivKey, err := base64.StdEncoding.DecodeString(string(ecdspkbytes))
		if err != nil {
			return err
		}
		ecdsaPk, err := keys.GetPrivateKeyFromBytes(ecdsPrivKey[:], keys.SECP256K1)
		if err != nil {
			return err
		}
		nodekey, err := p2p.LoadOrGenNodeKey(filepath.Join(readDir, consensus.NodeKeyFilename))
		pvFile := privval.LoadOrGenFilePV(filepath.Join(readDir, "priv_validator_key.json"), filepath.Join(readDir, "priv_validator_state.json"))
		if err != nil {
			return errors.Wrap(err, "Failed to generate node Key")
		}
		err = copyTofolders(readDir, dataDir, configDir)
		if err != nil {
			return err
		}
		// Generate new configuration file
		cfg := config.DefaultServerConfig()

		ethConnection := config.EthereumChainDriverConfig{Connection: args.ethUrl}
		cfg.EthChainDriver = &ethConnection
		cfg.Node.NodeName = nodeName
		cfg.Node.LogLevel = args.loglevel
		cfg.Node.DB = args.dbType
		if args.createEmptyBlock {
			cfg.Consensus.CreateEmptyBlocks = true
			cfg.Consensus.CreateEmptyBlocksInterval = 3000
		} else {
			cfg.Consensus.CreateEmptyBlocks = false
		}

		cfg.Network.RPCAddress = generateAddress(generatePort(), true)
		cfg.Network.P2PAddress = generateAddress(generatePort(), true)
		cfg.Network.SDKAddress = generateAddress(generatePort(), true, true)
		cfg.Network.OLVMAddress = generateAddress(generatePort(), true)

		n := node{IsValidator: isValidator, Cfg: cfg, Key: nodekey, EsdcaPk: ecdsaPk}
		if isValidator {
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodeName,
				Power:   1,
			}
			n.Validator = validator
			validatorList[i] = validator
		}
		nodeList[i] = n
		persistentPeers[i] = n.connectionDetails()
	}

	onsOp := getOnsOpt()
	btccdo := getBtcOpt()
	cdoBytes, err := ioutil.ReadFile(filepath.Join(mainnetCmdArgs.pvkey_Dir, "cdOpts.json"))
	if err != nil {
		return err
	}
	cdo := ethchain.ChainDriverOption{}
	err = json.Unmarshal(cdoBytes, &cdo)
	if err != nil {
		return err
	}
	states := getInitialState(args, nodeList, cdo, *onsOp, btccdo)

	genesisDoc, err := consensus.NewGenesisDoc(getChainID(), states)
	if err != nil {
		return errors.Wrap(err, "failed to create new genesis file")
	}
	genesisDoc.Validators = validatorList

	for _, nodeName := range ctx.names {
		nodeDir := filepath.Join(args.outputDir, nodeName)
		configDir := filepath.Join(nodeDir, "consensus", "config")
		err := genesisDoc.SaveAs(filepath.Join(configDir, "genesis.json"))
		if err != nil {
			return err
		}
	}
	// Save the files to the node's relevant directory
	generateBTCPort := portGenerator(18831)
	generateETHPort := portGenerator(28101)

	var swapNodes []string
	if args.allowSwap {
		swapNodes = ctx.names[1:4]
	}
	isSwapNode := func(name string) bool {
		for _, nodeName := range swapNodes {
			if nodeName == name {
				return true
			}
		}
		return false
	}

	//deploy contract and get contract addr
	//Saving config.toml for each node
	for _, node := range nodeList {
		node.Cfg.P2P.PersistentPeers = persistentPeers
		// Modify the btc and eth ports
		if args.allowSwap && isSwapNode(node.Cfg.Node.NodeName) {
			node.Cfg.Network.BTCAddress = generateAddress(generateBTCPort(), false)
			node.Cfg.Network.ETHAddress = generateAddress(generateETHPort(), false)
		}
		//	node.Cfg.EthChainDriver.ContractAddress = contractaddr
		err := node.Cfg.SaveFile(filepath.Join(args.outputDir, node.Cfg.Node.NodeName, config.FileName))
		if err != nil {
			return err
		}
	}

	ctx.logger.Info("Created configuration files for", strconv.Itoa(totalNodes), "nodes in", args.outputDir)
	return nil
}
func createDirectories(configDir string, dataDir string, nodeDataDir string) error {
	dirs := []string{configDir, dataDir, nodeDataDir}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, config.DirPerms)
		if err != nil {
			return err
		}
	}
	return nil
}
func copyTofolders(readDir string, dataDir string, configDir string) error {
	err := move(filepath.Join(readDir, "priv_validator_state.json"), filepath.Join(dataDir, "priv_validator_state.json"))
	if err != nil {
		return err
	}
	err = move(filepath.Join(readDir, "priv_validator_key.json"), filepath.Join(configDir, "priv_validator_key.json"))
	if err != nil {
		return err
	}
	err = move(filepath.Join(readDir, "priv_validator_key_ecdsa.json"), filepath.Join(configDir, "priv_validator_key_ecdsa.json"))
	if err != nil {
		return err
	}
	err = move(filepath.Join(readDir, "node_key.json"), filepath.Join(configDir, "node_key.json"))
	if err != nil {
		return err
	}
	return nil
}
func move(source string, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		return err
	}
	os.Remove(source)
	return nil
}

func getChainID() string {
	chainID := "OneLedger-" + randStr(2)
	if initCmdArgs.chainID != "" {
		chainID = initCmdArgs.chainID
	}
	return chainID
}

func getOnsOpt() *ons.Options {

	perblock, _ := big.NewInt(0).SetString("100000000000000", 10)
	baseDomainPrice, _ := big.NewInt(0).SetString("1000000000000000000000", 10)
	return &ons.Options{
		Currency:          "OLT",
		PerBlockFees:      *balance.NewAmountFromBigInt(perblock),
		FirstLevelDomains: []string{"ol"},
		BaseDomainPrice:   *balance.NewAmountFromBigInt(baseDomainPrice),
	}
}

func getInitialState(args *mainnetArgument, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
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
