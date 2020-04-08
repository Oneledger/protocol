package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethChain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/tendermint/tendermint/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Oneledger/protocol/app"
	olNode "github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ConsensusParams contains consensus critical parameters that determine the
// validity of blocks. Originally from Tendermint but redefined here to
// customize the JSON output as all values need to be encoded as string.
type ConsensusParams struct {
	Block     BlockParams     `json:"block"`
	Evidence  EvidenceParams  `json:"evidence"`
	Validator ValidatorParams `json:"validator"`
}

// BlockParams define limits on the block size and gas plus minimum time
// between blocks.
type BlockParams struct {
	MaxBytes int64 `json:"max_bytes,string"`
	MaxGas   int64 `json:"max_gas,string"`
	// Minimum time increment between consecutive blocks (in milliseconds)
	// Not exposed to the application.
	TimeIotaMs int64 `json:"time_iota_ms,string"`
}

// EvidenceParams determine how we handle evidence of malfeasance.
type EvidenceParams struct {
	MaxAge int64 `json:"max_age,string"` // only accept new evidence more recent than this
}

// ValidatorParams restrict the public key types validators can use.
// NOTE: uses ABCI pubkey naming, not Amino names.
type ValidatorParams struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type publicKey struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

type GenesisValidator struct {
	Address string    `json:"address"`
	PubKey  publicKey `json:"pub_key"`
	Power   int64     `json:"power,string"`
	Name    string    `json:"name"`
}

var saveStateCmd = &cobra.Command{
	Use:   "save_state",
	Short: "Save Chain State to a file",
	RunE:  SaveState,
}

type saveStateCmdContext struct {
	cfg       *config.Server
	logger    *log.Logger
	outputDir string
	filename  string
	rootDir   string
	chainId   string
	version   int64
}

var (
	saveStateCtx    = &saveStateCmdContext{}
	genesisTimeFile = "/tmp/genTime.tmp"
)

func (ctx *saveStateCmdContext) init(rootDir string) error {
	ctx.logger = log.NewLoggerWithPrefix(os.Stdout, "olfullnode node")

	cfg := &config.Server{}
	rootPath, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}

	ctx.rootDir = rootPath

	err = cfg.ReadFile(cfgPath(rootPath))
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration file at at %s", cfgPath(rootPath))
	}

	ctx.cfg = cfg

	return nil
}

func init() {
	RootCmd.AddCommand(saveStateCmd)
	saveStateCmd.Flags().StringVarP(&saveStateCtx.outputDir, "outDir", "o", "./", "Directory to store Chain State File, default current folder.")
	saveStateCmd.Flags().StringVarP(&saveStateCtx.chainId, "chainId", "c", "OneLedger-DEV", "Chain ID for each node to start with.")
	saveStateCmd.Flags().StringVarP(&saveStateCtx.filename, "filename", "f", "genesis_dump.json", "Name of file that stores the Chain State.")
	saveStateCmd.Flags().Int64Var(&saveStateCtx.version, "version", 0, "the version that need to be dumped, default the latest version")
}

func SaveState(cmd *cobra.Command, args []string) error {
	ctx := saveStateCtx
	err := ctx.init(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	appNodeContext, err := olNode.NewNodeContext(ctx.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create app's node context")
	}

	application, err := app.NewApp(ctx.cfg, appNodeContext)
	if err != nil {
		return errors.Wrap(err, "failed to create new app")
	}

	err = application.Prepare()
	if err != nil {
		return err
	}

	return SaveChainState(application, saveStateCtx.filename, saveStateCtx.outputDir)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func formatConsensusParams(params *types.ConsensusParams) ConsensusParams {
	cParams := ConsensusParams{
		Block:     BlockParams(params.Block),
		Evidence:  EvidenceParams(params.Evidence),
		Validator: ValidatorParams(params.Validator),
	}
	return cParams
}

func setGenesisTime(genesisTime *time.Time) error {
	if fileExists(genesisTimeFile) {
		file, _ := ioutil.ReadFile(genesisTimeFile)
		err := json.Unmarshal(file, genesisTime)
		if err != nil {
			return err
		}
	} else {
		file, _ := json.Marshal(genesisTime)
		err := ioutil.WriteFile(genesisTimeFile, file, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func startBlock(writer io.Writer, name string) bool {
	_, err := writer.Write([]byte(name + ":{"))
	if err != nil {
		return false
	}
	return true
}

func endBlock(writer io.Writer) bool {
	_, err := writer.Write([]byte("}"))
	if err != nil {
		return false
	}
	return true
}

func writeStruct(writer io.Writer, obj interface{}) bool {
	str, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return false
	}
	_, err = writer.Write(str)
	if err != nil {
		return false
	}
	_, err = writer.Write([]byte("\n"))
	if err != nil {
		return false
	}

	return true
}

func writeStructWithTag(writer io.Writer, obj interface{}, tag string) bool {
	str, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return false
	}
	_, err = writer.Write([]byte("\"" + tag + "\"" + ":"))
	_, err = writer.Write(str)
	_, err = writer.Write([]byte(",\n"))
	return true
}

func writeListWithTag(ctx app.StorageCtx, writer io.Writer, tag string) bool {
	delimiter := ","

	_, err := writer.Write([]byte("\"" + tag + "\"" + ":["))
	switch section := tag; section {
	case "validators":
		DumpValidatorsToFile(ctx.Validators, writer, writeStruct)
	case "balances":
		DumpBalanceToFile(ctx.Balances, writer, writeStruct)
	case "staking":
		DumpStakingToFile(ctx.Validators, writer, writeStruct)
	case "domains":
		DumpDomainToFile(ctx.Domains, ctx.Version, writer, writeStruct)
	case "trackers":
		DumpTrackerToFile(ctx.Trackers.WithPrefixType(ethereum.PrefixOngoing), writer, writeStruct)
		DumpTrackerToFile(ctx.Trackers.WithPrefixType(ethereum.PrefixPassed), writer, writeStruct)
		DumpTrackerToFile(ctx.Trackers.WithPrefixType(ethereum.PrefixFailed), writer, writeStruct)
	case "fees":
		DumpFeesToFile(ctx.FeePool, writer, writeStruct)
		delimiter = ""
	}
	_, err = writer.Write([]byte("]"))
	_, err = writer.Write([]byte(delimiter + "\n"))

	if err != nil {
		return false
	}

	return true
}

func SaveChainState(application *app.App, filename string, directory string) error {
	ctx := application.Context.Storage()
	version, err := ctx.Chainstate.LoadVersion(saveStateCtx.version)
	if err != nil {
		return err
	}
	ctx.Version = version
	appState := consensus.AppState{}
	appState.Currencies, err = ctx.Govern.GetCurrencies()
	appState.Chain.Hash = nil  //ctx.Hash
	appState.Chain.Version = 0 //ctx.Version

	chainID := saveStateCtx.chainId
	genesisDoc, err := consensus.NewGenesisDoc(chainID, appState)
	if err != nil {
		err = errors.Wrap(err, "Failed to create Genesis object")
		fmt.Println(err)
		return err
	}

	//Get Genesis Time if temp file created. Otherwise save Genesis time to file.
	err = setGenesisTime(&genesisDoc.GenesisTime)
	if err != nil {
		return err
	}

	genesisDoc.AppHash = []byte{}

	genesis, err := json.Marshal(genesisDoc)
	jsonDecoder := json.NewDecoder(strings.NewReader(string(genesis)))

	//Start writing state to output file
	path := filepath.Join(directory, filename)
	writer, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = writer.Close()
	}()

	token, err := jsonDecoder.Token()
	_, err = fmt.Fprint(writer, token)
	_, err = writer.Write([]byte("\n"))

	for jsonDecoder.More() {
		token, err = jsonDecoder.Token()

		switch value := fmt.Sprintf("%s", token); value {
		case "genesis_time":
			writeStructWithTag(writer, genesisDoc.GenesisTime.UTC(), value)
		case "chain_id":
			writeStructWithTag(writer, genesisDoc.ChainID, value)
		case "consensus_params":
			writeStructWithTag(writer, formatConsensusParams(genesisDoc.ConsensusParams), value)
		}
	}

	writeListWithTag(ctx, writer, "validators")

	writeStructWithTag(writer, genesisDoc.AppHash, "app_hash")

	startBlock(writer, "\"app_state\"")
	writeStructWithTag(writer, appState.Currencies, "currencies")
	writeStructWithTag(writer, GetGovernance(ctx.Govern), "governance")
	writeStructWithTag(writer, appState.Chain, "state")
	writeListWithTag(ctx, writer, "balances")
	writeListWithTag(ctx, writer, "staking")
	writeListWithTag(ctx, writer, "domains")
	writeListWithTag(ctx, writer, "trackers")
	writeListWithTag(ctx, writer, "fees")
	endBlock(writer)

	token, err = jsonDecoder.Token()
	_, err = fmt.Fprint(writer, token)
	if err != nil {
		return err
	}

	return nil
}

func DumpFeesToFile(st *fees.Store, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","
	st.Iterate(func(addr keys.Address, coin balance.Coin) bool {
		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {
				return true
			}
		}
		fee := consensus.BalanceState{}
		fee.Address = addr
		fee.Currency = coin.Currency.Name

		if coin.Amount != nil {
			fee.Amount = *coin.Amount
		}

		fn(writer, fee)
		iterator++
		return false
	})
	return
}

func DumpStakingToFile(vs *identity.ValidatorStore, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","
	vs.Iterate(func(key keys.Address, validator *identity.Validator) bool {
		stake := consensus.Stake{}
		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {
				return true
			}
		}
		stake.Amount = validator.Staking
		stake.ECDSAPubKey = validator.ECDSAPubKey
		stake.Pubkey = validator.PubKey
		stake.ValidatorAddress = validator.Address
		stake.StakeAddress = validator.StakeAddress
		stake.Name = validator.Name
		//stake.Power = validator.Power

		fn(writer, stake)
		iterator++
		return false
	})

	return
}

//Retrieves complete list of Balance records and writes them to an io stream.
func DumpBalanceToFile(bs *balance.Store, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","
	bs.IterateAll(func(addr keys.Address, coin string, amt balance.Amount) bool {
		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {
				return true
			}
		}
		balance := consensus.BalanceState{}
		balance.Address = addr
		balance.Amount = amt
		balance.Currency = coin

		fn(writer, balance)
		iterator++
		return false
	})
	return
}

func DumpDomainToFile(ds *ons.DomainStore, height int64, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","

	ds.Iterate(func(name ons.Name, domain *ons.Domain) bool {

		if domain.ExpireHeight < height {
			return false
		}

		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {

				return true
			}
		}

		domainState := consensus.DomainState{}
		domainState.Name = domain.Name.String()
		domainState.Beneficiary = domain.Beneficiary
		domainState.Owner = domain.Owner
		domainState.CreationHeight = 0
		domainState.LastUpdateHeight = 0
		domainState.ExpireHeight = domain.ExpireHeight - height
		domainState.ActiveFlag = domain.ActiveFlag
		domainState.OnSaleFlag = domain.OnSaleFlag
		domainState.URI = domain.URI
		domainState.SalePrice = domain.SalePrice

		if !fn(writer, domainState) {
			return true
		}
		iterator++
		return false
	})
	return
}

//Save all Current trackers to Genesis file. Currently only supported for Ethereum.
//TODO: Add support for Bitcoin
func DumpTrackerToFile(ts *ethereum.TrackerStore, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","

	ts.Iterate(func(name *ethChain.TrackerName, tracker *ethereum.Tracker) bool {
		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {
				return true
			}
		}

		trackerState := consensus.Tracker{}
		trackerState.Type = tracker.Type
		trackerState.State = tracker.State
		trackerState.FinalityVotes = tracker.FinalityVotes
		trackerState.ProcessOwner = tracker.ProcessOwner
		trackerState.SignedETHTx = tracker.SignedETHTx
		trackerState.TrackerName = tracker.TrackerName
		trackerState.Validators = tracker.Validators
		trackerState.To = tracker.To

		fn(writer, trackerState)
		iterator++

		return false
	})
}

func GetGovernance(gs *governance.Store) *consensus.GovernanceState {
	btcOption, err := gs.GetBTCChainDriverOption()
	if err != nil {
		fmt.Print("Error Reading BTC chain driver options: ", err)
		return nil
	}

	ethOption, err := gs.GetETHChainDriverOption()
	if err != nil {
		fmt.Print("Error Reading ETH chain driver options: ", err)
		return nil
	}
	onsOption, err := gs.GetONSOptions()
	if err != nil {
		fmt.Print("Error Reading ONS Domain options: ", err)
		return nil
	}

	feeOption, err := gs.GetFeeOption()
	if err != nil {
		fmt.Print("Error Reading Fee options: ", err)
		return nil
	}

	return &consensus.GovernanceState{
		FeeOption:   *feeOption,
		ETHCDOption: *ethOption,
		BTCCDOption: *btcOption,
		ONSOptions:  *onsOption,
	}
}

func DumpValidatorsToFile(vs *identity.ValidatorStore, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","

	vs.Iterate(func(key keys.Address, validator *identity.Validator) bool {
		if iterator != 0 {
			_, err := writer.Write([]byte(delimiter))
			if err != nil {
				return true
			}
		}

		stake := GenesisValidator{
			Address: hex.EncodeToString(validator.Address.Bytes()),
			PubKey: publicKey{
				Type:  "tendermint/PubKeyEd25519",
				Value: validator.PubKey.Data,
			},
			Name:  validator.Name,
			Power: validator.Power,
		}

		fn(writer, stake)
		iterator++
		return false
	})

	return
}
