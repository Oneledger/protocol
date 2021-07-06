package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/rewards"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/ons"

	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"

	"github.com/Oneledger/protocol/log"
)

var (
	//Lock Limits
	totalETHSupply     = "2000000000000000000" // 20 ETH
	totalTTCSupply     = "2000000000000000000" // 2 Token
	totalBTCSupply     = "1000000000"          // 10 BTC
	lockBalanceAddress = "oneledgerSupplyAddress"

	ethBlockConfirmation = int64(12)
	btcBlockConfirmation = int64(6)

	proposalInitialFunding, _ = balance.NewAmountFromString("1000000000", 10)
	proposalFundingGoal, _    = balance.NewAmountFromString("10000000000", 10)
	proposalFundingDeadline   = int64(75001)
	proposalVotingDeadline    = int64(150000)
	//proposalFundingDeadline     = int64(10)
	//proposalVotingDeadline      = int64(12)
	proposalPassPercentage      = 51
	bountyProgramAddr           = "oneledgerBountyProgram"
	executionCostAddrConfig     = "executionCostConfig"
	executionCostAddrCodeChange = "executionCostCodeChange"
	executionCostAddrGeneral    = "executionCostGeneral"
	passedProposalDistribution  = governance.ProposalFundDistribution{
		Validators:     18.00,
		FeePool:        18.00,
		Burn:           18.00,
		ExecutionCost:  18.00,
		BountyPool:     10.00,
		ProposerReward: 18.00,
	}
	failedProposalDistribution = governance.ProposalFundDistribution{
		Validators:     10.00,
		FeePool:        10.00,
		Burn:           10.00,
		ExecutionCost:  20.00,
		BountyPool:     50.00,
		ProposerReward: 00.00,
	}

	estimatedSecondsPerCycle  = int64(1728)
	blockSpeedCalculateCycle  = int64(100)
	burnoutRate, _            = balance.NewAmountFromString("5000000000000000000", 10)
	yearCloseWindow           = int64(3600 * 24)
	yearBlockRewardShare_1, _ = balance.NewAmountFromString("70000000000000000000000000", 10)
	yearBlockRewardShare_2, _ = balance.NewAmountFromString("70000000000000000000000000", 10)
	yearBlockRewardShare_3, _ = balance.NewAmountFromString("40000000000000000000000000", 10)
	yearBlockRewardShare_4, _ = balance.NewAmountFromString("40000000000000000000000000", 10)
	yearBlockRewardShare_5, _ = balance.NewAmountFromString("30000000000000000000000000", 10)
	yearBlockRewardShares     = []balance.Amount{
		*yearBlockRewardShare_1,
		*yearBlockRewardShare_2,
		*yearBlockRewardShare_3,
		*yearBlockRewardShare_4,
		*yearBlockRewardShare_5,
	}

	testnetArgs = &testnetConfig{}

	testnetCmd = &cobra.Command{
		Use:   "devnet",
		Short: "Initializes files for a devnet",
		RunE:  runDevnet,
	}
)

type testnetConfig struct {
	// Number of validators
	numValidators    int
	numNonValidators int
	outputDir        string
	p2pPort          int
	chainID          string
	dbType           string
	namesPath        string
	createEmptyBlock bool
	// Total amount of funds to be shared across each node
	totalFunds          int64
	initialTokenHolders []string

	ethUrl                   string
	deploySmartcontracts     bool
	ethpk                    string
	cloud                    bool
	loglevel                 int
	rewardsInterval          int64
	reserved_domains         string
	maturityTime             int64
	delegRewardsMaturityTime int64
	votingDeadline           int64
	fundingDeadline          int64
	timeoutcommit            int64
	docker                   bool
	subnet                   string
	minVotesRequired         int64
	blockVotesDiff           int64
	topValidators            int64
}

func init() {
	initCmd.AddCommand(testnetCmd)
	testnetCmd.Flags().IntVar(&testnetArgs.numValidators, "validators", 4, "Number of validators to initialize devnet with")
	testnetCmd.Flags().IntVar(&testnetArgs.numNonValidators, "nonvalidators", 0, "Number of non-validators to initialize the devnet with")
	testnetCmd.Flags().StringVarP(&testnetArgs.outputDir, "dir", "o", "./", "Directory to store initialization files for the devnet, default current folder")
	testnetCmd.Flags().BoolVar(&testnetArgs.createEmptyBlock, "empty_blocks", false, "Allow creating empty blocks")
	testnetCmd.Flags().StringVar(&testnetArgs.chainID, "chain_id", "", "Specify a chain ID, a random one is generated if not given")
	testnetCmd.Flags().StringVar(&testnetArgs.dbType, "db_type", "goleveldb", "Specify the type of DB backend to use: (goleveldb|cleveldb)")
	testnetCmd.Flags().StringVar(&testnetArgs.namesPath, "names", "", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	// 1 billion by default
	testnetCmd.Flags().Int64Var(&testnetArgs.totalFunds, "total_funds", 1000000000, "The total amount of tokens in circulation")
	testnetCmd.Flags().StringSliceVar(&testnetArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")
	testnetCmd.Flags().StringVar(&testnetArgs.ethUrl, "eth_rpc", "", "URL for ethereum network")
	testnetCmd.Flags().BoolVar(&testnetArgs.deploySmartcontracts, "deploy_smart_contracts", false, "deploy eth contracts")
	testnetCmd.Flags().StringVar(&testnetArgs.ethpk, "eth_pk", "", "ethereum test private key")
	testnetCmd.Flags().BoolVar(&testnetArgs.cloud, "cloud_deploy", false, "set true for deploying on cloud")
	testnetCmd.Flags().IntVar(&testnetArgs.loglevel, "loglevel", 3, "Specify the log level for olfullnode. 0: Fatal, 1: Error, 2: Warning, 3: Info, 4: Debug, 5: Detail")
	testnetCmd.Flags().Int64Var(&testnetArgs.rewardsInterval, "rewards_interval", 1, "Block rewards interval")
	testnetCmd.Flags().StringVar(&testnetArgs.reserved_domains, "reserved_domains", "", "Directory which contains Reserved domains list")
	testnetCmd.Flags().Int64Var(&testnetArgs.maturityTime, "maturity_time", 109200, "Set Maturity time for staking")
	testnetCmd.Flags().Int64Var(&testnetArgs.delegRewardsMaturityTime, "deleg_rewards_maturity_time", 109200, "Set Maturity time for delegation rewards")
	testnetCmd.Flags().Int64Var(&testnetArgs.fundingDeadline, "funding_deadline", 75001, "Set Maturity time for staking")
	testnetCmd.Flags().Int64Var(&testnetArgs.votingDeadline, "voting_deadline", 150000, "Set Maturity time for staking")
	testnetCmd.Flags().Int64Var(&testnetArgs.timeoutcommit, "timeout_commit", 1000, "Set timecommit for blocks")
	testnetCmd.Flags().Int64Var(&testnetArgs.minVotesRequired, "min_votes", 2, "Set Min votes for evidence")
	testnetCmd.Flags().Int64Var(&testnetArgs.blockVotesDiff, "block_diff", 4, "Set block Votes diff evidence")
	testnetCmd.Flags().Int64Var(&testnetArgs.topValidators, "top_validators", 4, "Set top valdiators for staking")
	testnetCmd.Flags().BoolVar(&testnetArgs.docker, "docker", false, "set true for deploying on docker")
	testnetCmd.Flags().StringVar(&testnetArgs.subnet, "subnet", "10.5.0.0/16", "subnet where all docker containers will run")
}

func randStr(size int) string {
	bz := make([]byte, size)
	_, err := rand.Read(bz)
	if err != nil {
		return "deadbeef"
	}
	return hex.EncodeToString(bz)
}

// Need to maintain a list of nodes and be able to:
// (1) Keep track of all of their P2P addresses including their addresses
// (2) Modify their configurations to have each one have its persistent peer set
type node struct {
	isValidator bool
	cfg         *config.Server
	dir         string
	key         *p2p.NodeKey
	esdcaPk     keys.PrivateKey
	validator   consensus.GenesisValidator
}

func (n node) connectionDetails() string {
	var addr string
	if n.cfg.Network.ExternalP2PAddress == "" {
		addr = n.cfg.Network.P2PAddress
	} else {
		addr = n.cfg.Network.ExternalP2PAddress
	}

	u, _ := url.Parse(addr)
	return fmt.Sprintf("%s@%s", n.key.ID(), u.Host)
}

// This function maintains a running counter of ports
func portGenerator(startingPort int) func() int {
	count := startingPort
	return func() int {
		port := count
		count++
		return port
	}
}

func generateIP(nodeID int) (result string) {
	result = "127.0.0.1"
	if testnetArgs.docker {
		//parse subnet IP
		ip, _, err := net.ParseCIDR(testnetArgs.subnet)
		if err != nil {
			return
		}
		ipv4 := ip.To4()
		if len(ipv4) != 4 {
			return
		}
		ip3 := byte(nodeID) + 10
		result = fmt.Sprintf("%d.%d.%d.%d", ipv4[0], ipv4[1], ipv4[2], ip3)
	}
	return
}

func generateAddress(nodeID int, port int, flags ...bool) string {
	// flags
	var hasProtocol, isRPC bool
	switch len(flags) {
	case 2:
		hasProtocol, isRPC = flags[0], flags[1]
	case 1:
		hasProtocol = flags[0]
	default:
	}

	var prefix string
	ip := generateIP(nodeID)
	protocol := "tcp://"
	if isRPC {
		protocol = "http://"
	}
	if hasProtocol {
		prefix = protocol + ip
	} else {
		prefix = ip
	}

	return fmt.Sprintf("%s:%d", prefix, port)
}

// Just a basic context for the devnet cmd
type devnetContext struct {
	names  []string
	logger *log.Logger
}

func newDevnetContext(args *testnetConfig) (*devnetContext, error) {
	logger := log.NewLoggerWithPrefix(os.Stdout, "olfullnode devnet")

	names := nodeNamesWithZeros("", args.numNonValidators+args.numValidators)
	// TODO: Reading from a file is actually unimplemented right now
	if args.namesPath != "" {
		logger.Warn("--names parameter is unimplemented")
	}

	return &devnetContext{
		names:  names,
		logger: logger,
	}, nil
}

// Returns a list of names with the given prefix and a number after the prefix afterwards
func nodeNamesWithZeros(prefix string, total int) []string {
	names := make([]string, total)
	//maxZeroes := len(strconv.Itoa(total))

	generateName := func(i int) string {
		name := prefix
		num := strconv.Itoa(i)
		// Unpad nums
		return name + num
	}

	for i := 0; i < total; i++ {
		names[i] = generateName(i)
	}
	return names
}

func runDevnet(_ *cobra.Command, _ []string) error {

	ctx, err := newDevnetContext(testnetArgs)
	if err != nil {
		return errors.Wrap(err, "runDevnet failed")
	}
	args := testnetArgs
	if !args.cloud {
		setEnvVariables()
	}
	totalNodes := args.numValidators + args.numNonValidators

	if totalNodes > len(ctx.names) {
		return fmt.Errorf("Don't have enough node names, can't specify more than %d nodes", len(ctx.names))
	}
	var reserveDomains []reservedDomain
	var initialAddrs []keys.Address
	if len(testnetArgs.initialTokenHolders) > 0 {
		initialAddrs, err = getInitialAddress(initialAddrs, testnetArgs.initialTokenHolders)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(testnetArgs.reserved_domains); err == nil {
		reserveDomains, err = getReservedDomains(testnetArgs.reserved_domains)
		if err != nil {
			return err
		}
	}

	if args.dbType != "cleveldb" && args.dbType != "goleveldb" {
		ctx.logger.Error("Invalid dbType specified, using goleveldb...", "dbType", args.dbType)
		args.dbType = "goleveldb"
	}

	generatePort := portGenerator(26600)
	generateWeb3Port := portGenerator(10545)

	validatorList := make([]consensus.GenesisValidator, args.numValidators)
	nodeList := make([]node, totalNodes)
	persistentPeers := make([]string, totalNodes)
	url, err := getEthUrl(args.ethUrl)
	if err != nil {
		return err
	}
	// Create the GenesisValidator list and its key files priv_validator_key.json and node_key.json
	minSelfStaking := int64(3000000)
	validatorCount := int64(0)
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

		cfg.Consensus.TimeoutCommit = config.Duration(testnetArgs.timeoutcommit)

		cfg.Network.RPCAddress = generateAddress(i, generatePort(), true)
		cfg.Network.P2PAddress = generateAddress(i, generatePort(), true)
		cfg.Network.SDKAddress = generateAddress(i, generatePort(), true)
		cfg.Network.Web3Address = generateAddress(i, generateWeb3Port(), true, true)

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
		//ethecdsaPrivKey, err := crypto.GenerateKey()
		//if err != nil {
		//	return errors.Wrap(err, "error generating ETHSECP key")
		//}
		//ethecdsaPk, err := keys.GetPrivateKeyFromBytes(ethecdsaPrivKey.D.Bytes(), keys.ETHSECP)
		//if err != nil {
		//	return errors.Wrap(err, "error generating BTCECSECP private key")
		//}

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
			// each validator takes a different power
			validatorCount++
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodeName,
				Power:   minSelfStaking * validatorCount,
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
			cdo, err = deployethcdcontract(url, nodeList, args.ethpk)
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
	proposalFundingDeadline = args.fundingDeadline
	proposalVotingDeadline = args.votingDeadline
	propOpt := governance.ProposalOptionSet{
		ConfigUpdate: governance.ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrConfig,
		},
		CodeChange: governance.ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrCodeChange,
		},
		General: governance.ProposalOption{
			InitialFunding:         proposalInitialFunding,
			FundingGoal:            proposalFundingGoal,
			FundingDeadline:        proposalFundingDeadline,
			VotingDeadline:         proposalVotingDeadline,
			PassPercentage:         proposalPassPercentage,
			PassedFundDistribution: passedProposalDistribution,
			FailedFundDistribution: failedProposalDistribution,
			ProposalExecutionCost:  executionCostAddrGeneral,
		},
		BountyProgramAddr: bountyProgramAddr,
	}

	rewzOpt := rewards.Options{
		RewardInterval:           args.rewardsInterval,
		RewardPoolAddress:        "rewardpool",
		RewardCurrency:           "OLT",
		EstimatedSecondsPerCycle: estimatedSecondsPerCycle,
		BlockSpeedCalculateCycle: blockSpeedCalculateCycle,
		YearCloseWindow:          yearCloseWindow,
		YearBlockRewardShares:    yearBlockRewardShares,
		BurnoutRate:              *burnoutRate,
	}
	states := initialState(args, nodeList, *cdo, *onsOp, btccdo, propOpt, reserveDomains, initialAddrs, rewzOpt)

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

	//deploy contract and get contract addr

	for _, node := range nodeList {
		node.cfg.P2P.PersistentPeers = persistentPeers
		err := node.cfg.SaveFile(filepath.Join(node.dir, config.FileName))
		if err != nil {
			return err
		}
	}

	ctx.logger.Info("Created configuration files for", strconv.Itoa(totalNodes), "nodes in", args.outputDir)

	return nil
}

func initialState(args *testnetConfig, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
	btcOption bitcoin.ChainDriverOption, propOpt governance.ProposalOptionSet, reservedDomains []reservedDomain, initialAddrs []keys.Address, rewardOpt rewards.Options) consensus.AppState {

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
	rewards := rewards.RewardMasterState{
		RewardState: rewards.NewRewardState(),
		CumuState:   rewards.NewRewardCumuState(),
	}
	domains := make([]consensus.DomainState, 0, len(nodeList))
	fees_db := make([]consensus.BalanceState, 0, len(nodeList))
	total := olt.NewCoinFromInt(args.totalFunds)

	// staking
	stakingOption := delegation.Options{
		MinSelfDelegationAmount: *balance.NewAmount(3000000),
		MinDelegationAmount:     *balance.NewAmount(1),
		TopValidatorCount:       args.topValidators,
		MaturityTime:            args.maturityTime,
	}

	// network delegation
	delegOption := network_delegation.Options{
		RewardsMaturityTime: args.delegRewardsMaturityTime,
	}

	// evidence
	evidenceOption := evidence.Options{
		MinVotesRequired: args.minVotesRequired,
		BlockVotesDiff:   args.blockVotesDiff,

		PenaltyBasePercentage: 30,
		PenaltyBaseDecimals:   100,

		PenaltyBountyPercentage: 50,
		PenaltyBountyDecimals:   100,

		PenaltyBurnPercentage: 50,
		PenaltyBurnDecimals:   100,

		ValidatorReleaseTime:    0,
		ValidatorVotePercentage: 50,
		ValidatorVoteDecimals:   100,

		AllegationPercentage: 50,
		AllegationDecimals:   100,
	}

	// Cutting from current stake or cutting from begining stake

	//var initialAddrs []keys.Address
	initAddrIndex := 0
	//for _, addr := range args.initialTokenHolders {
	//	tmpAddr := keys.Address{}
	//	err := tmpAddr.UnmarshalText([]byte(addr))
	//	if err != nil {
	//		fmt.Println("Error adding initial address:", addr)
	//		continue
	//	}
	//	initialAddrs = append(initialAddrs, tmpAddr)
	//}

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
			Amount:           *balance.NewAmountFromInt(node.validator.Power),
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
				Amount:   *vt.NewCoinFromInt(1000).Amount,
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
					ExpireHeight:     42048000, // 20 Years . Taking one block every 15s
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
		Rewards:    rewards,
		Domains:    domains,
		Fees:       fees_db,
		Governance: governance.GovernanceState{
			FeeOption:       feeOpt,
			ETHCDOption:     option,
			BTCCDOption:     btcOption,
			ONSOptions:      onsOption,
			PropOptions:     propOpt,
			StakingOptions:  stakingOption,
			DelegOptions:    delegOption,
			EvidenceOptions: evidenceOption,
			RewardOptions:   rewardOpt,
		},
	}
}

func deployethcdcontract(conn string, nodeList []node, ethprivKey string) (*ethchain.ChainDriverOption, error) {
	pkStr := ""
	if len(ethprivKey) == 0 {
		b1 := make([]byte, 64)
		f, err := os.Open(os.Getenv("ETHPKPATH"))
		if err != nil {
			return nil, errors.Wrap(err, "Error Reading File")
		}
		pk, err := f.Read(b1)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading private key")
		}
		pkStr = string(b1[:pk])
	} else {
		pkStr = ethprivKey
	}
	//fmt.Println("Private key used to deploy : ", pkStr)

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

	//gasPrice, err := cli.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	gasLimit := uint64(6721974) // in units

	auth := bind.NewKeyedTransactor(privatekey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = big.NewInt(18000000000)

	initialValidatorList := make([]common.Address, 0, 10)
	lock_period := big.NewInt(25)

	tokenSupplyTestToken := new(big.Int)
	validatorInitialFund := big.NewInt(30000000000000000) //300000000000000000
	tokenSupplyTestToken, ok = tokenSupplyTestToken.SetString("1000000000000000000000", 10)
	if !ok {
		return nil, errors.New("Unabe to create total supplu for token")
	}
	if !ok {
		return nil, errors.New("Unable to create wallet transfer amount")
	}
	for _, node := range nodeList {
		privkey := keys.ETHSECP256K1TOECDSA(node.esdcaPk.Data)
		nonce, err := cli.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return nil, err
		}

		//fmt.Println("nonce", nonce)
		pubkey := privkey.Public()
		ecdsapubkey, ok := pubkey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("failed to cast pubkey")
		}
		addr := crypto.PubkeyToAddress(*ecdsapubkey)
		if node.validator.Address.String() == "" {
			continue
		}

		initialValidatorList = append(initialValidatorList, addr)
		tx := types.NewTransaction(nonce, addr, validatorInitialFund, auth.GasLimit, auth.GasPrice, nil)
		fmt.Println(addr.Hex(), ":", validatorInitialFund, "wei")
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

	num_of_validators := big.NewInt(1)
	address, _, _, err := contract.DeployLockRedeem(auth, cli, initialValidatorList, lock_period, fromAddress, num_of_validators)
	if err != nil {
		return nil, errors.Wrap(err, "Deployement Eth LockRedeem")
	}
	time.Sleep(10 * time.Second)
	err = activateContract(fromAddress, address, privatekey, cli)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to activate new contract")
	}
	tokenAddress := common.Address{}
	ercAddress := common.Address{}

	//auth.Nonce = big.NewInt(int64(nonce + 1))
	//tokenAddress, _, _, err := ethcontracts.DeployERC20Basic(auth, cli, tokenSupplyTestToken)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Deployement Test Token")
	//}
	//auth.Nonce = big.NewInt(int64(nonce + 2))
	//ercAddress, _, _, err := ethcontracts.DeployLockRedeemERC(auth, cli, initialValidatorList)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Deployement ERC LockRedeem")
	//}

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

	rawTxBytes2, _ := rlp.EncodeToBytes(ts2[0])
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
