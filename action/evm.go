package action

import (
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

// AbciHeaderToTendermint is a util function to parse a tendermint ABCI Header to
// tendermint types Header.
func AbciHeaderToTendermint(header *abci.Header) tmtypes.Header {
	return tmtypes.Header{
		Version: version.Consensus{
			Block: version.Protocol(header.Version.Block),
			App:   version.Protocol(header.Version.App),
		},
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    header.Time,

		LastBlockID: tmtypes.BlockID{
			Hash: header.LastBlockId.Hash,
			PartsHeader: tmtypes.PartSetHeader{
				Total: int(header.LastBlockId.PartsHeader.Total),
				Hash:  header.LastBlockId.PartsHeader.Hash,
			},
		},
		LastCommitHash: header.LastCommitHash,
		DataHash:       header.DataHash,

		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,

		EvidenceHash:    header.EvidenceHash,
		ProposerAddress: header.ProposerAddress,
	}
}

// GetHashFn implements vm.GetHashFunc for protocol.
func GetHashFn(s *CommitStateDB) ethvm.GetHashFunc {
	return func(height uint64) ethcmn.Hash {
		return s.GetHeightHash(height)
	}
}

func EthereumConfig(chainID string) *ethparams.ChainConfig {
	return &ethparams.ChainConfig{
		ChainID:        utils.HashToBigInt(chainID),
		HomesteadBlock: big.NewInt(0),

		DAOForkBlock:   big.NewInt(0),
		DAOForkSupport: true,

		EIP150Block: big.NewInt(0),
		EIP150Hash:  common.Hash{},

		EIP155Block: big.NewInt(0),
		EIP158Block: big.NewInt(0),

		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),

		EWASMBlock: nil,
	}
}

type EVMConfig struct {
	gasPrice  *big.Int
	gasLimit  uint64
	extraEIPs []int
}

func NewEVMConfig(gasPrice *big.Int, gasLimit uint64, extraEIPs []int) *EVMConfig {
	return &EVMConfig{
		gasPrice:  gasPrice,
		gasLimit:  gasLimit,
		extraEIPs: extraEIPs,
	}
}

type EVMTransaction struct {
	stateDB      *CommitStateDB
	header       *abci.Header
	from         keys.Address
	to           *keys.Address
	nonce        uint64
	value        *big.Int
	data         []byte
	ecfg         *EVMConfig
	isSimulation bool
	state        *storage.State
}

var (
	DefaultGasLimit uint64   = 10_000_000
	DefaultGasPrice *big.Int = big.NewInt(0)
)

func NewEVMTransaction(stateDB *CommitStateDB, header *abci.Header, from keys.Address, to *keys.Address, nonce uint64, value *big.Int, data []byte, isSimulation bool) *EVMTransaction {
	return &EVMTransaction{
		stateDB:      stateDB,
		header:       header,
		from:         from,
		to:           to,
		nonce:        nonce,
		value:        value,
		data:         data,
		isSimulation: isSimulation,
		// NOTE: Decide what to do with the gas price as we have a fee system
		ecfg: NewEVMConfig(DefaultGasPrice, DefaultGasLimit, make([]int, 0)),
	}
}

func (etx *EVMTransaction) NewEVM() *ethvm.EVM {
	blockCtx := ethvm.BlockContext{
		CanTransfer: ethcore.CanTransfer,
		Transfer:    ethcore.Transfer,
		GetHash:     GetHashFn(etx.stateDB),
		Coinbase:    ethcmn.Address{}, // there's no beneficiary since we're not mining
		GasLimit:    etx.ecfg.gasLimit,
		BlockNumber: big.NewInt(etx.header.GetHeight()),
		Time:        big.NewInt(etx.header.Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
	}

	vmConfig := ethvm.Config{
		ExtraEips: etx.ecfg.extraEIPs,
	}

	txCtx := ethvm.TxContext{
		Origin:   etx.Origin(),
		GasPrice: etx.ecfg.gasPrice,
	}

	ethConfig := EthereumConfig(etx.header.ChainID)
	return ethvm.NewEVM(blockCtx, txCtx, etx.stateDB, ethConfig, vmConfig)
}

func (etx *EVMTransaction) Origin() ethcmn.Address {
	return ethcmn.BytesToAddress(etx.from)
}

func (etx *EVMTransaction) From() ethcmn.Address {
	return ethcmn.BytesToAddress(etx.from)
}

func (etx *EVMTransaction) To() *ethcmn.Address {
	if etx.to == nil {
		return nil
	}
	ethTo := ethcmn.BytesToAddress(*etx.to)
	return &ethTo
}

func (etx *EVMTransaction) Apply(vmenv *ethvm.EVM, tx RawTx) (*ExecutionResult, error) {
	msg := ethtypes.NewMessage(etx.From(), etx.To(), etx.nonce, etx.value, uint64(tx.Fee.Gas), etx.ecfg.gasPrice, nil, nil, etx.data, make(ethtypes.AccessList, 0), true)

	txHash := etx.stateDB.GetCurrentTxHash()

	if !etx.isSimulation {
		// Clear cache of accounts to handle changes outside of the EVM
		etx.stateDB.UpdateAccounts()
	}

	executionResult, err := ApplyMessage(vmenv, msg, new(ethcore.GasPool).AddGas(uint64(uint64(tx.Fee.Gas))))

	if !etx.isSimulation {
		if err == nil {
			// calculating bloom
			logs, err := etx.stateDB.GetLogs(txHash)
			if err != nil {
				return nil, err
			}
			bloomInt := big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs))
			etx.stateDB.logger.Debug("tx hash", txHash, "bloom created", bloomInt)
			etx.stateDB.Bloom.Or(etx.stateDB.Bloom, bloomInt)
			etx.stateDB.logger.Debug("block bloom updated", etx.stateDB.Bloom)
		}
		// Ensure any modifications are committed to the state
		if err := etx.stateDB.Finalise(true); err != nil {
			return nil, err
		}
		// Commit state objects to store
		if _, err := etx.stateDB.Commit(true); err != nil {
			return nil, err
		}
		etx.stateDB.logger.Debugf("State finalized\n")
	}
	return executionResult, err
}
