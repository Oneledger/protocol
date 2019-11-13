package main

import (
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/app"
	olnode "github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
}

var saveStateCtx = &saveStateCmdContext{}

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
	testnetCmd.Flags().StringVarP(&saveStateCtx.outputDir, "outDir", "c", "./", "Directory to store Chain State File, default current folder.")
	testnetCmd.Flags().StringVarP(&saveStateCtx.filename, "filename", "f", "genesis_dump.json", "Name of file that stores the Chain State.")
}

func SaveState(cmd *cobra.Command, args []string) error {
	ctx := saveStateCtx
	err := ctx.init(rootArgs.rootDir)
	if err != nil {
		return errors.Wrap(err, "failed to initialize config")
	}

	appNodeContext, err := olnode.NewNodeContext(ctx.cfg)
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

	SaveChainState(application, saveStateCtx.outputDir, saveStateCtx.filename)

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
	_, err = writer.Write([]byte("\n"))
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
	case "balances":
		DumpBalanceToFile(ctx.Balances, writer, writeStruct)
	case "staking":
		DumpStakingToFile(ctx.Validators, writer, writeStruct)
	case "domains":
		DumpDomainToFile(ctx.Domains, writer, writeStruct)
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

func SaveChainState(application *app.App, filename string, directory string) {
	ctx := application.Context.Storage()
	appState := consensus.AppState{}
	var err error
	appState.Currencies, err = ctx.Govern.GetCurrencies()
	appState.FeeOption = fees.FeeOption{
		FeeCurrency:   ctx.FeeOption.FeeCurrency,
		MinFeeDecimal: ctx.FeeOption.MinFeeDecimal,
	}
	appState.Chain.Hash = ctx.Hash
	appState.Chain.Version = ctx.Version

	chainID := "OneLedger-" + randStr(2)
	genesisDoc, err := consensus.NewGenesisDoc(chainID, appState)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to create Genesis object"))
	}

	genesisDoc.AppHash = ctx.Hash

	genesis, err := json.Marshal(genesisDoc)
	jsonDecoder := json.NewDecoder(strings.NewReader(string(genesis)))

	//Start writing state to output file
	path := filepath.Join(directory, filename)
	writer, err := os.Create(path)

	token, err := jsonDecoder.Token()
	_, err = fmt.Fprint(writer, token)
	_, err = writer.Write([]byte("\n"))

	for jsonDecoder.More() {
		token, err = jsonDecoder.Token()

		switch value := fmt.Sprintf("%s", token); value {
		case "genesis_time":
			writeStructWithTag(writer, genesisDoc.GenesisTime, value)
		case "chain_id":
			writeStructWithTag(writer, genesisDoc.ChainID, value)
		case "consensus_params":
			writeStructWithTag(writer, genesisDoc.ConsensusParams, value)
		case "validators":
			writeStructWithTag(writer, genesisDoc.Validators, value)
		}
	}

	writeStructWithTag(writer, genesisDoc.AppHash, "app_hash")

	startBlock(writer, "\"app_state\"")
	writeStructWithTag(writer, appState.Currencies, "currencies")
	writeStructWithTag(writer, appState.FeeOption, "fee_option")
	writeStructWithTag(writer, appState.Chain, "chain")
	writeListWithTag(ctx, writer, "balances")
	writeListWithTag(ctx, writer, "staking")
	writeListWithTag(ctx, writer, "domains")
	writeListWithTag(ctx, writer, "fees")
	endBlock(writer)

	token, err = jsonDecoder.Token()
	_, err = fmt.Fprint(writer, token)

	err = writer.Close()
}

type BalanceState struct {
	Address  keys.Address   `json:"address"`
	Currency string         `json:"currency"`
	Amount   balance.Amount `json:"amount"`
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
		fee := BalanceState{}
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

type Stake struct {
	ValidatorAddress keys.Address
	StakeAddress     keys.Address
	Pubkey           keys.PublicKey
	ECDSAPubKey      keys.PublicKey
	Name             string
	Amount           balance.Amount
}

func DumpStakingToFile(vs *identity.ValidatorStore, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","
	vs.Iterate(func(key keys.Address, validator *identity.Validator) bool {
		stake := Stake{}
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
		balance := BalanceState{}
		balance.Address = addr
		balance.Amount = amt
		balance.Currency = coin

		fn(writer, balance)
		iterator++
		return false
	})
	return
}

type DomainState struct {
	OwnerAddress   keys.Address `json:"ownerAddress"`
	AccountAddress keys.Address `json:"accountAddress"`
	Name           string       `json:"name"`
}

func DumpDomainToFile(ds *ons.DomainStore, writer io.Writer, fn func(writer io.Writer, obj interface{}) bool) {
	iterator := 0
	delimiter := ","
	ds.State.IterateRange(
		ds.Prefix,
		storage.Rangefix(string(ds.Prefix)),
		true,
		func(key, value []byte) bool {
			domainState := DomainState{}
			domain := &ons.Domain{}
			err := ds.Szlr.Deserialize(value, domain)
			if err != nil {
				return true
			}
			if iterator != 0 {
				_, err := writer.Write([]byte(delimiter))
				if err != nil {
					return true
				}
			}
			domainState.Name = domain.Name
			domainState.AccountAddress = domain.AccountAddress
			domainState.OwnerAddress = domain.OwnerAddress

			fn(writer, domainState)
			iterator++
			return false
		},
	)
	return
}
