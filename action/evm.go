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

// GetHashFn implements vm.GetHashFunc for OneLedger protocol. It handles 3 cases:
//  1. The requested height matches the current height (and thus same epoch number, could take from cache)
//  2. The requested height is from an previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func GetHashFn(s *CommitStateDB, header *abci.Header) ethvm.GetHashFunc {
	return func(height uint64) common.Hash {
		switch {
		case header.GetHeight() == int64(height):
			// Case 1: The requested height matches the one from the CommitStateDB so we can retrieve the block
			// hash directly from the CommitStateDB.
			return s.bhash

		case header.GetHeight() > int64(height):
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			return s.GetHeightHash(height)

		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return common.Hash{}
		}
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
	}
}

func (etx *EVMTransaction) NewEVM() *ethvm.EVM {
	blockCtx := ethvm.BlockContext{
		CanTransfer: ethcore.CanTransfer,
		Transfer:    ethcore.Transfer,
		GetHash:     GetHashFn(etx.stateDB, etx.header),
		Coinbase:    ethcmn.BytesToAddress(etx.header.ProposerAddress),
		GasLimit:    DefaultGasLimit,
		BlockNumber: big.NewInt(etx.header.GetHeight()),
		Time:        big.NewInt(etx.header.Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
	}

	vmConfig := ethvm.Config{
		NoBaseFee: true,
		ExtraEips: make([]int, 0), // skip right now
	}

	txCtx := ethvm.TxContext{
		Origin:   etx.Origin(),
		GasPrice: DefaultGasPrice,
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
	// gas price ignoring here as we have a separate handler for it
	// no gas fee cap and tips, as another reward system in protocol
	// nonce not a fake
	// access list is currently not supported
	msg := ethtypes.NewMessage(etx.From(), etx.To(), etx.nonce, etx.value, uint64(tx.Fee.Gas), vmenv.TxContext.GasPrice, big.NewInt(0), big.NewInt(0), etx.data, make(ethtypes.AccessList, 0), etx.isSimulation)

	if !etx.isSimulation {
		// Clear cache of accounts to handle changes outside of the EVM
		etx.stateDB.UpdateAccounts()
	}

	executionResult, err := ApplyMessage(vmenv, msg, new(ethcore.GasPool).AddGas(uint64(uint64(tx.Fee.Gas))))

	if !etx.isSimulation {
		// only if tx is OK otherwise all are already reverted
		if err == nil {
			// calculating bloom for the block
			logs, err := etx.stateDB.GetLogs(etx.stateDB.thash)
			if err != nil {
				// error must not be here, but in case we will have it, to know what to dig
				return nil, err
			}
			bloomInt := big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs))
			etx.stateDB.Bloom.Or(etx.stateDB.Bloom, bloomInt)
		}
		// Ensure any modifications are committed to the state
		if err := etx.stateDB.Finalise(true); err != nil {
			return nil, err
		}
		// Commit state objects to store
		if _, err := etx.stateDB.Commit(true); err != nil {
			return nil, err
		}
	}
	return executionResult, err
}
