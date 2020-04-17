package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"
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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/chains/ethereum/contract"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/ons"

	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
)

type mainnetArgument struct {
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

var mainnetCmdArgs = &mainnetArgument{}

var mainnetCmd = &cobra.Command{
	Use:   "mainnet",
	Short: "Initializes a genesis file for OneLedger network",
	RunE:  runMainnet,
}

func init() {
	initCmd.AddCommand(mainnetCmd)
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.numValidators, "validators", 4, "Number of validators to initialize devnet with")
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.numNonValidators, "nonvalidators", 0, "Number of non-validators to initialize the devnet with")
	mainnetCmd.Flags().StringVarP(&mainnetCmdArgs.outputDir, "dir", "o", "./", "Directory to store initialization files for the devnet, default current folder")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.allowSwap, "enable_swaps", false, "Allow swaps")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.createEmptyBlock, "empty_blocks", false, "Allow creating empty blocks")
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.chainID, "chain_id", "", "Specify a chain ID, a random one is generated if not given")
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.dbType, "db_type", "goleveldb", "Specify the type of DB backend to use: (goleveldb|cleveldb)")
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.namesPath, "names", "", "Specify a path to a file containing a list of names separated by newlines if you want the nodes to be generated with human-readable names")
	// 1 billion by default
	mainnetCmd.Flags().Int64Var(&mainnetCmdArgs.totalFunds, "total_funds", 1000000000, "The total amount of tokens in circulation")
	mainnetCmd.Flags().StringSliceVar(&mainnetCmdArgs.initialTokenHolders, "initial_token_holders", []string{}, "Initial list of addresses that hold an equal share of Total funds")
	mainnetCmd.Flags().StringVar(&mainnetCmdArgs.ethUrl, "eth_rpc", "", "URL for ethereum network")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.deploySmartcontracts, "deploy_smart_contracts", false, "deploy eth contracts")
	mainnetCmd.Flags().BoolVar(&mainnetCmdArgs.cloud, "cloud_deploy", false, "set true for deploying on cloud")
	mainnetCmd.Flags().IntVar(&mainnetCmdArgs.loglevel, "loglevel", 3, "Specify the log level for olfullnode. 0: Fatal, 1: Error, 2: Warning, 3: Info, 4: Debug, 5: Detail")

}

func runMainnet(_ *cobra.Command, _ []string) error {

	ctx, err := newMainetContext(mainnetCmdArgs)
	if err != nil {
		return errors.Wrap(err, "runDevnet failed")
	}
	args := mainnetCmdArgs
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
		nodeKey, err := p2p.LoadOrGenNodeKey(filepath.Join(configDir, consensus.NodeKeyFilename))
		if err != nil {
			return errors.Wrap(err, "Failed to generate node key")
		}

		// Make private validator file
		pvFile := privval.GenFilePV(filepath.Join(configDir, consensus.PrivValidatorKeyFilename),
			filepath.Join(dataDir, consensus.PrivValidatorStateFilename))
		pvFile.Save()

		ecdsaPrivKey := secp256k1.GenPrivKey()
		ecdsaPrivKeyBytes := base64.StdEncoding.EncodeToString([]byte(ecdsaPrivKey[:]))
		ecdsaPk, err := keys.GetPrivateKeyFromBytes([]byte(ecdsaPrivKey[:]), keys.SECP256K1)
		if err != nil {
			return errors.Wrap(err, "error generating secp256k1 private key")
		}
		ecdsaFile := strings.Replace(consensus.PrivValidatorKeyFilename, ".json", "_ecdsa.json", 1)
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
		//fmt.Println("witness_key_address: ", ecdsaPrivKey.PubKey().Address().String())
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
			cdo, err = ethContractMainnet(url, nodeList)
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

	states := mainnetInitialState(args, nodeList, *cdo, *onsOp, btccdo)

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
	//Saving config.toml for each node
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

func mainnetInitialState(args *mainnetArgument, nodeList []node, option ethchain.ChainDriverOption, onsOption ons.Options,
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

func ethContractMainnet(conn string, nodeList []node) (*ethchain.ChainDriverOption, error) {

	f, err := os.Open(os.Getenv("ETHPKPATH"))
	if err != nil {
		return nil, errors.Wrap(err, "Error Reading File")
	}
	b1 := make([]byte, 64)
	pk, err := f.Read(b1)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private key")
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
