package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

type testnetConfig struct {
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
	totalFunds        int64
	singleOriginFunds bool
}

var testnetArgs = &testnetConfig{}

var testnetCmd = &cobra.Command{
	Use:   "devnet",
	Short: "Initializes files for a devnet",
	RunE:  runDevnet,
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
	testnetCmd.Flags().BoolVar(&testnetArgs.singleOriginFunds, "single_origin_funds", false, "If specified, allocates all possible tokens to a single account")
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

func generateAddress(port int, hasProtocol bool) string {
	var prefix string
	ip := "127.0.0.1"
	protocol := "tcp://"
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

// padZeroes takes the maximum number of zeroes allowed and pa
func padZeroes(str string, total int) string {
	prefix := strings.Repeat("0", total-len(str))
	return prefix + str
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

func runDevnet(cmd *cobra.Command, _ []string) error {
	ctx, err := newDevnetContext(testnetArgs)
	if err != nil {
		return errors.Wrap(err, "runDevnet failed")
	}
	args := testnetArgs

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
		cfg.Node.NodeName = nodeName
		cfg.Node.DB = args.dbType
		if args.createEmptyBlock {
			cfg.Consensus.CreateEmptyBlocks = true
		} else {
			cfg.Consensus.CreateEmptyBlocks = false
		}
		cfg.Network.RPCAddress = generateAddress(generatePort(), true)
		cfg.Network.P2PAddress = generateAddress(generatePort(), true)
		cfg.Network.SDKAddress = generateAddress(generatePort(), true)
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
		pvFile := privval.GenFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(dataDir, "priv_validator_state.json"))
		pvFile.Save()

		if isValidator {
			validator := consensus.GenesisValidator{
				Address: pvFile.GetAddress(),
				PubKey:  pvFile.GetPubKey(),
				Name:    nodeName,
				Power:   1,
			}
			validatorList[i] = validator
		}

		// Save the nodes to a list so we can iterate again and
		n := node{isValidator, cfg, nodeDir, nodeKey}
		nodeList[i] = n
		persistentPeers[i] = n.connectionDetails()
	}

	// Create the non validator nodes

	// Create the genesis file
	chainID := "OneLedger-" + randStr(2)
	if args.chainID != "" {
		chainID = args.chainID
	}

	currencies, states := initialState(args, nodeList)

	genesisDoc, err := consensus.NewGenesisDoc(chainID, currencies, states)
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

	for _, node := range nodeList {
		node.cfg.P2P.PersistentPeers = persistentPeers
		// Modify the btc and eth ports
		if args.allowSwap && isSwapNode(node.cfg.Node.NodeName) {
			node.cfg.Network.BTCAddress = generateAddress(generateBTCPort(), false)
			node.cfg.Network.ETHAddress = generateAddress(generateETHPort(), false)
		}

		err := node.cfg.SaveFile(filepath.Join(node.dir, config.FileName))
		if err != nil {
			return err
		}
	}

	ctx.logger.Info("Created configuration files for", strconv.Itoa(totalNodes), "nodes in", args.outputDir)

	return nil
}

func initialState(args *testnetConfig, nodeList []node) ([]balance.Currency, []consensus.StateInput) {
	olt := balance.Currency{"OLT", chain.Type(0), 18}
	vt := balance.Currency{"VT", chain.Type(0), 0}
	currencies := []balance.Currency{olt, vt}

	var out []consensus.StateInput
	// If single origin is active, then the first node in the list should hold all the funds
	if args.singleOriginFunds {
		b0 := balance.NewBalance()
		b0.AddCoin(olt.NewCoinFromInt(args.totalFunds))
		b0.AddCoin(vt.NewCoinFromInt(1))

		out = []consensus.StateInput{
			{
				Address: nodeList[0].key.PubKey().Address().String(),
				Balance: *b0,
			},
		}
		for _, node := range nodeList {
			if !node.isValidator {
				continue
			}
			b := balance.NewBalance()
			b.AddCoin(vt.NewCoinFromInt(1))
			out = append(out, consensus.StateInput{
				Address: node.key.PubKey().Address().String(),
				Balance: *b,
			})
		}
		return currencies, out
	}

	out = make([]consensus.StateInput, len(nodeList))
	for i, node := range nodeList {
		share := args.totalFunds / int64(len(nodeList))
		b := balance.NewBalance()
		b.AddCoin(olt.NewCoinFromInt(share))
		if node.isValidator {
			b.AddCoin(vt.NewCoinFromInt(1))
		}
		out[i] = consensus.StateInput{
			Address: node.key.PubKey().Address().String(),
			Balance: *b,
		}

		if node.isValidator {
			out[i].Balance.AddCoin(vt.NewCoinFromInt(1))
		}
	}
	return currencies, out
}
