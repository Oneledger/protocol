package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/log"
)

type genesisArgument struct {
	// Number of validators
	numValidators    int
	numNonValidators int
	outputDir        string
	pvkey_Dir        string
	reserved_domains string
	p2pPort          int
	allowSwap        bool
	chainID          string
	dbType           string
	//namesPath        string
	createEmptyBlock bool
	// Total amount of funds to be shared across each node
	totalFunds          int64
	initialTokenHolders []string

	ethUrl               string
	deploySmartcontracts bool
	cloud                bool
	loglevel             int
}

var genesisCmdArgs = &genesisArgument{}

var genesisCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Initializes a genesis file for OneLedger network",
	RunE:  runGenesis,
}

type reservedDomain struct {
	domainName string
	domainUrl  string
}
type mainetContext struct {
	logger *log.Logger
	names  []string
}

func init() {
	initCmd.AddCommand(genesisCmd)
	genesisCmd.Flags().StringVarP(&genesisCmdArgs.pvkey_Dir, "pv_dir", "p", "$OLDATA/mainnet/", "Directory which contains Genesis File and NodeList")
	genesisCmd.Flags().StringVarP(&genesisCmdArgs.reserved_domains, "reserved_domains", "r", "$OLDATA/mainnet/", "Path to the file which contains the domainlist")
	genesisCmd.Flags().StringVarP(&genesisCmdArgs.outputDir, "dir", "o", "$OLDATA/mainnet/", "Directory to store initialization files for the ")
	//genesisCmd.Flags().BoolVar(&genesisCmdArgs.allowSwap, "enable_swaps", false, "Allow swaps")
	genesisCmd.Flags().IntVar(&genesisCmdArgs.numValidators, "validators", 4, "Number of validators to initialize mainnetnet with")
	genesisCmd.Flags().IntVar(&genesisCmdArgs.numNonValidators, "nonvalidators", 1, "Number of fullnodes to initialize mainnetnet with")
	genesisCmd.Flags().BoolVar(&genesisCmdArgs.createEmptyBlock, "empty_blocks", false, "Allow creating empty blocks")
	genesisCmd.Flags().StringVar(&genesisCmdArgs.dbType, "db_type", "goleveldb", "Specify the type of DB backend to use: (goleveldb|cleveldb)")
	//genesisCmd.Flags().StringVar(&genesisCmdArgs.namesPath, "names", "", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	// 1 billion by default
	genesisCmd.Flags().StringVar(&genesisCmdArgs.ethUrl, "eth_rpc", "HTTP://127.0.0.1:7545", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	genesisCmd.Flags().IntVar(&genesisCmdArgs.loglevel, "loglevel", 3, "Specify the log level for olfullnode. 0: Fatal, 1: Error, 2: Warning, 3: Info, 4: Debug, 5: Detail")
	genesisCmd.Flags().Int64Var(&genesisCmdArgs.totalFunds, "total_funds", 400000000, "The total amount of tokens in circulation")
	genesisCmd.Flags().StringSliceVar(&genesisCmdArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")
	genesisCmd.Flags().BoolVar(&genesisCmdArgs.deploySmartcontracts, "deploy_smart_contracts", false, "deploy eth contracts")
}

func newMainetContext(args *genesisArgument) (*mainetContext, error) {

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
		return nil, errors.New("Not enough key Pairs present the the directory ")
	}
	return &mainetContext{
		names:  names,
		logger: logger,
	}, nil

}

func runGenesis(_ *cobra.Command, _ []string) error {
	var reserveDomains []reservedDomain
	var initialAddrs []keys.Address
	setEnvVariables()
	ctx, err := newMainetContext(genesisCmdArgs)
	if err != nil {
		return err
	}
	if len(genesisCmdArgs.initialTokenHolders) > 0 {
		reserveDomains, err = getReservedDomains(genesisCmdArgs.reserved_domains)
		if err != nil {
			return err
		}
		initialAddrs, err = getInitialAddress(initialAddrs, genesisCmdArgs.initialTokenHolders)
		if err != nil {
			return err
		}
	}
	args := genesisCmdArgs
	totalNodes := len(ctx.names)
	if args.dbType != "cleveldb" && args.dbType != "goleveldb" {
		ctx.logger.Error("Invalid dbType specified, using goleveldb...", "dbType", args.dbType)
		args.dbType = "goleveldb"
	}

	generatePort := portGenerator(26600)

	persistentPeers := make([]string, totalNodes)
	nodeList := make([]node, totalNodes)
	validatorList := make([]consensus.GenesisValidator, args.numValidators)
	url, err := getEthUrl(args.ethUrl)
	if err != nil {
		return err
	}
	// Create the GenesisValidator list and its key files priv_validator_key.json and node_key.json
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
			return errors.Wrap(err, "Failed to generate node key")
		}
		err = copyTofolders(readDir, dataDir, configDir)
		if err != nil {
			return err
		}
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
		// cfg.Network.OLVMAddress = generateAddress(generatePort(), true)

		n := node{isValidator: isValidator, cfg: cfg, key: nodekey, esdcaPk: ecdsaPk}
		if isValidator {
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodeName,
				Power:   1,
			}
			n.validator = validator
			validatorList[i] = validator
		}
		nodeList[i] = n
		persistentPeers[i] = n.connectionDetails()
	}

	onsOp := getOnsOpt()
	btccdo := getBtcOpt()
	//cdoBytes, err := ioutil.ReadFile(filepath.Join(genesisCmdArgs.pvkey_Dir, "cdOpts.json"))
	cdo := &ethchain.ChainDriverOption{}
	url, err = getEthUrl(genesisCmdArgs.ethUrl)
	if err != nil {
		return err
	}
	fmt.Println("Ethereum Deployment Network :", url)
	if genesisCmdArgs.deploySmartcontracts {
		if len(genesisCmdArgs.ethUrl) > 0 {
			cdo, err = getEthOpt(url, nodeList)
			if err != nil {
				return errors.Wrap(err, "failed to deploy the initial eth contract")
			}
		}
	}

	//if err != nil {
	//	return err
	//}
	//cdo := ethchain.ChainDriverOption{}
	//err = json.Unmarshal(cdoBytes, &cdo)
	//if err != nil {
	//	return err
	//}
	//os.Remove(filepath.Join(genesisCmdArgs.pvkey_Dir, "cdOpts.json"))
	states := getInitialState(args, nodeList, *cdo, *onsOp, btccdo, reserveDomains, initialAddrs)

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
	for _, node := range nodeList {
		node.cfg.P2P.PersistentPeers = persistentPeers
		err := node.cfg.SaveFile(filepath.Join(args.outputDir, node.cfg.Node.NodeName, config.FileName))
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
	if genesisCmdArgs.chainID != "" {
		chainID = genesisCmdArgs.chainID
	}
	return chainID
}

func getReservedDomains(domainlistPath string) ([]reservedDomain, error) {
	var reserved []reservedDomain
	reservedDomainsBytes, err := os.Open(domainlistPath)
	if err != nil {
		fmt.Println("Error")
		return reserved, err
	}
	fileScanner := bufio.NewScanner(reservedDomainsBytes)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		dom := strings.Split(fileScanner.Text(), ",")
		domain := reservedDomain{
			domainName: dom[0],
			domainUrl:  dom[1],
		}
		reserved = append(reserved, domain)
	}
	return reserved, err
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

func getInitialState(args *genesisArgument, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
	btcOption bitcoin.ChainDriverOption, reservedDomains []reservedDomain, initialAddrs []keys.Address) consensus.AppState {
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
	domains := make([]consensus.DomainState, 0, len(reservedDomains))
	fees_db := make([]consensus.BalanceState, 0, len(nodeList))
	total := olt.NewCoinFromInt(args.totalFunds)
	initAddrIndex := 0

	// staking
	// staking
	stakingOption := delegation.Options{
		MinSelfDelegationAmount: *balance.NewAmount(3000000),
		MinDelegationAmount:     *balance.NewAmount(3000000),
		TopValidatorCount:       32,
		MaturityTime:            150000,
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
			amt := int64(100)
			share := total.DivideInt64(int64(len(args.initialTokenHolders)))
			balances = append(balances, consensus.BalanceState{
				Address:  acct,
				Currency: olt.Name,
				Amount:   *olt.NewCoinFromAmount(*share.Amount).Amount,
			})
			balances = append(balances, consensus.BalanceState{
				Address:  acct,
				Currency: vt.Name,
				Amount:   *vt.NewCoinFromInt(amt).Amount,
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

	if len(args.initialTokenHolders) > 0 {
		domainsPerHolder := len(reservedDomains) / len(args.initialTokenHolders)
		start := 0
		end := 0
		for _, addr := range initialAddrs {
			start = end
			end = end + domainsPerHolder
			domain := reservedDomains[start:end]
			for _, d := range domain {
				domains = append(domains, consensus.DomainState{
					Owner:            addr,
					Beneficiary:      addr,
					Name:             d.domainName,
					CreationHeight:   0,
					LastUpdateHeight: 0,
					ExpireHeight:     4204800000, // 2000 Years . Taking one block every 15s
					ActiveFlag:       false,
					OnSaleFlag:       false,
					URI:              d.domainUrl,
					SalePrice:        nil,
				})
			}
		}
	}

	return consensus.AppState{
		Currencies: currencies,
		Balances:   balances,
		Staking:    staking,
		Domains:    domains,
		Fees:       fees_db,
		Governance: consensus.GovernanceState{
			FeeOption:      feeOpt,
			ETHCDOption:    option,
			BTCCDOption:    btcOption,
			ONSOptions:     onsOption,
			StakingOptions: stakingOption,
		},
	}
}

func getInitialAddress(initialAddrs []keys.Address, initialTokenHolders []string) ([]keys.Address, error) {
	if len(initialTokenHolders) == 0 {
		return nil, errors.New("No address provided for intital token holders")
	}
	//fmt.Println("Initial token holders :", initialTokenHolders)
	for _, addr := range initialTokenHolders {
		tmpAddr := keys.Address{}
		err := tmpAddr.UnmarshalText([]byte(addr))
		if err != nil {
			fmt.Println("Error adding initial address:", addr)
			return nil, err
		}
		initialAddrs = append(initialAddrs, tmpAddr)
	}
	return initialAddrs, nil
}
