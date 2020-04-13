package main

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
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
	"github.com/Oneledger/protocol/data/ons"

	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
)

type genesisCmdArgument struct {
	// Number of validators
	numValidators    int
	numNonValidators int
	outputDir        string
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

var genesisCmdArgs = &genesisCmdArgument{}

var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Initializes a genesis file for OneLedger network",
	RunE:  runGenesis,
}

func init() {
	initCmd.AddCommand(testnetCmd)
	testnetCmd.Flags().IntVar(&testnetArgs.numValidators, "validators", 4, "Number of validators to initialize devnet with")
	testnetCmd.Flags().IntVar(&testnetArgs.numNonValidators, "nonvalidators", 0, "Number of non-validators to initialize the devnet with")
	testnetCmd.Flags().StringVarP(&testnetArgs.outputDir, "dir", "o", "./", "Directory to store initialization files for the devnet, default current folder")
	testnetCmd.Flags().BoolVar(&testnetArgs.allowSwap, "enable_swaps", false, "Allow swaps")
	testnetCmd.Flags().BoolVar(&testnetArgs.createEmptyBlock, "empty_blocks", false, "Allow creating empty blocks")
	testnetCmd.Flags().StringVar(&testnetArgs.chainID, "chain_id", "", "Specify a chain ID, a random one is generated if not given")
	testnetCmd.Flags().StringVar(&testnetArgs.dbType, "db_type", "goleveldb", "Specify the type of DB backend to use: (goleveldb|cleveldb)")
	testnetCmd.Flags().StringVar(&testnetArgs.namesPath, "names", "", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	// 1 billion by default
	testnetCmd.Flags().Int64Var(&testnetArgs.totalFunds, "total_funds", 1000000000, "The total amount of tokens in circulation")
	testnetCmd.Flags().StringSliceVar(&testnetArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")
	testnetCmd.Flags().StringVar(&testnetArgs.ethUrl, "eth_rpc", "", "URL for ethereum network")
	testnetCmd.Flags().BoolVar(&testnetArgs.deploySmartcontracts, "deploy_smart_contracts", false, "deploy eth contracts")
	testnetCmd.Flags().BoolVar(&testnetArgs.cloud, "cloud_deploy", false, "set true for deploying on cloud")
	testnetCmd.Flags().IntVar(&testnetArgs.loglevel, "loglevel", 3, "Specify the log level for olfullnode. 0: Fatal, 1: Error, 2: Warning, 3: Info, 4: Debug, 5: Detail")

}

func runGenesis(_ *cobra.Command, _ []string) error {

	ctx, err := newDevnetContext(testnetArgs)
	if err != nil {
		return errors.Wrap(err, "runDevnet failed")
	}
	args := testnetArgs
	if !args.cloud {
		setEnvVariablesGanache()
	}
	totalNodes := args.numValidators + args.numNonValidators

	if totalNodes > len(ctx.names) {
		return fmt.Errorf("Don't have enough node names, can't specify more than %d nodes", len(ctx.names))
	}

	if args.dbType != "cleveldb" && args.dbType != "goleveldb" {
		ctx.logger.Error("Invalid dbType specified, using goleveldb...", "dbType", args.dbType)
		args.dbType = "goleveldb"
	}

	generatePort := portGenerator(26600)

	validatorList := make([]consensus.GenesisValidator, args.numValidators)
	nodeList := make([]node, totalNodes)
	persistentPeers := make([]string, totalNodes)
	url, err := getEthUrl(args.ethUrl)
	if err != nil {
		return err
	}
	// Create the GenesisValidator list and its key files priv_validator_key.json and node_key.json
	for i := 0; i < totalNodes; i++ {
		isValidator := i < args.numValidators
		nodeName := ctx.names[i]
		nodeDir := filepath.Join(args.outputDir, nodeName+"-Node")
		configDir := filepath.Join(nodeDir, "consensus", "config")
		dataDir := filepath.Join(nodeDir, "consensus", "data")
		nodeDataDir := filepath.Join(nodeDir, "nodedata")

		// Generate new configuration file
		cfg := config.DefaultServerConfig()

		ethConnection := config.EthereumChainDriverConfig{Connection: url}
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

		dirs := []string{configDir, dataDir, nodeDataDir}
		for _, dir := range dirs {
			err := os.MkdirAll(dir, config.DirPerms)
			if err != nil {
				return err
			}
		}

		// Make node key
		nodeKey, err := p2p.LoadOrGenNodeKey(filepath.Join(configDir, "node_key.json"))
		if err != nil {
			ctx.logger.Error("error load or genning node key", "err", err)
			return err
		}

		// Make private validator file
		pvFile := privval.LoadOrGenFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(dataDir, "priv_validator_state.json"))
		pvFile.Save()

		ecdsaPrivKey, _ := btcec.NewPrivateKey(btcec.S256())
		ecdsaPrivKeyBytes := base64.StdEncoding.EncodeToString([]byte(ecdsaPrivKey.Serialize()))

		ecdsaPk, err := keys.GetPrivateKeyFromBytes([]byte(ecdsaPrivKey.Serialize()), keys.BTCECSECP)
		if err != nil {
			return errors.Wrap(err, "error generating BTCECSECP private key")
		}

		f, err := os.Create(filepath.Join(configDir, "priv_validator_key_ecdsa.json"))
		if err != nil {
			return errors.Wrap(err, "failed to open file to write validator ecdsa private key")
		}
		nPrivKeyBytes, err := f.WriteString(ecdsaPrivKeyBytes)
		if err != nil && nPrivKeyBytes != len(ecdsaPrivKeyBytes) {
			return errors.Wrap(err, "failed to write validator ecdsa private key")
		}
		err = f.Close()
		if err != nil {
			return errors.Wrap(err, "failed to save validator ecdsa private key")
		}

		// Save the nodes to a list so we can iterate again and
		n := node{isValidator: isValidator, cfg: cfg, dir: nodeDir, key: nodeKey, esdcaPk: ecdsaPk}
		if isValidator {
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodeName,
				Power:   1,
			}
			validatorList[i] = validator
			n.validator = validator
		}
		nodeList[i] = n
		persistentPeers[i] = n.connectionDetails()

	}

	// Create the non validator nodes

	// Create the genesis file
	chainID := "OneLedger-" + randStr(2)
	if args.chainID != "" {
		chainID = args.chainID
	}

	cdo := &ethchain.ChainDriverOption{}
	fmt.Println("Deployment Network :", url)
	fmt.Println("Deploy Smart contracts : ", args.deploySmartcontracts)
	if args.deploySmartcontracts {
		if len(args.ethUrl) > 0 {
			cdo, err = deployethcdcontract(url, nodeList)
			if err != nil {
				return errors.Wrap(err, "failed to deploy the initial eth contract")
			}
		}
	}

	perblock, _ := big.NewInt(0).SetString("100000000000000", 10)
	baseDomainPrice, _ := big.NewInt(0).SetString("1000000000000000000000", 10)
	onsOp := &ons.Options{
		Currency:          "OLT",
		PerBlockFees:      *balance.NewAmountFromBigInt(perblock),
		FirstLevelDomains: []string{"ol"},
		BaseDomainPrice:   *balance.NewAmountFromBigInt(baseDomainPrice),
	}

	btccdo := bitcoin.ChainDriverOption{
		"testnet3",
		totalBTCSupply,
		lockBalanceAddress,
		btcBlockConfirmation,
	}

	states := initialState(args, nodeList, *cdo, *onsOp, btccdo)

	genesisDoc, err := consensus.NewGenesisDoc(chainID, states)
	if err != nil {
		return errors.Wrap(err, "failed to create new genesis file")
	}
	genesisDoc.Validators = validatorList

	for i := 0; i < totalNodes; i++ {
		nodeName := ctx.names[i]
		nodeDir := filepath.Join(args.outputDir, nodeName+"-Node")
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

	for _, node := range nodeList {
		node.cfg.P2P.PersistentPeers = persistentPeers
		// Modify the btc and eth ports
		if args.allowSwap && isSwapNode(node.cfg.Node.NodeName) {
			node.cfg.Network.BTCAddress = generateAddress(generateBTCPort(), false)
			node.cfg.Network.ETHAddress = generateAddress(generateETHPort(), false)
		}
		//	node.cfg.EthChainDriver.ContractAddress = contractaddr
		err := node.cfg.SaveFile(filepath.Join(node.dir, config.FileName))
		if err != nil {
			return err
		}
	}

	ctx.logger.Info("Created configuration files for", strconv.Itoa(totalNodes), "nodes in", args.outputDir)

	return nil
}

func mainnetInitialState(args *genesisCmdArgument, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
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
		if !node.isValidator {
			continue
		}

		h, err := node.esdcaPk.GetHandler()
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
			stakeAddr = node.key.PubKey().Address().Bytes()
		}

		pubkey, _ := keys.PubKeyFromTendermint(node.validator.PubKey.Bytes())
		st := consensus.Stake{
			ValidatorAddress: node.validator.Address.Bytes(),
			StakeAddress:     stakeAddr,
			Pubkey:           pubkey,
			ECDSAPubKey:      h.PubKey(),
			Name:             node.validator.Name,
			Amount:           *vt.NewCoinFromInt(node.validator.Power).Amount,
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
			if !node.isValidator {
				amt = 1
			}
			share := total.DivideInt64(int64(len(nodeList)))
			balances = append(balances, consensus.BalanceState{
				Address:  node.key.PubKey().Address().Bytes(),
				Currency: olt.Name,
				Amount:   *olt.NewCoinFromAmount(*share.Amount).Amount,
			})
			balances = append(balances, consensus.BalanceState{
				Address:  node.key.PubKey().Address().Bytes(),
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
