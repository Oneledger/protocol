package vm

import (
	"math"
	"math/big"
	"os"
	"sync"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	// refunds
	RefundQuotientFrankenstein uint64 = 3

	// defaults
	DefaultDifficulty       *big.Int = big.NewInt(1)
	SimulationBlockGasLimit uint64   = 100_000_000
	DefaultBlockGasLimit    uint64   = math.MaxInt64
	DefaultGasPrice         *big.Int = big.NewInt(1_000_000_000)

	_ Message = (*EVMTransaction)(nil)
)

// GetHashFn implements vm.GetHashFunc for OneLedger protocol. It handles 3 cases:
//  1. The requested height matches the current height (and thus same epoch number, could take from cache)
//  2. The requested height is from an previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func GetHashFn(s *CommitStateDB, header *abci.Header) ethvm.GetHashFunc {
	cache := make(map[uint64]ethcmn.Hash)
	var rw sync.Mutex

	return func(height uint64) common.Hash {
		s.logger.Detail("GetHashFn current header height", header.GetHeight(), "with requested height", height)
		switch {
		case header.GetHeight() == int64(height):
			// Case 1: The requested height matches the one from the CommitStateDB so we can retrieve the block
			// hash directly from the CommitStateDB.
			// NOTE: Will never occur, left in case of VM changes
			return s.bhash

		case header.GetHeight() > int64(height):
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			// We are sure that we will have it in block store as it is commited

			// NOTE: In case of concurrency
			rw.Lock()
			defer rw.Unlock()

			if hash, ok := cache[height]; ok {
				return hash
			}
			hash := s.GetBlockHash(height)
			cache[height] = hash
			return hash

		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return common.Hash{}
		}
	}
}

func EthereumConfig(chainID string) *ethparams.ChainConfig {
	return &ethparams.ChainConfig{
		ChainID:        utils.HashToBigInt(chainID),
		HomesteadBlock: big.NewInt(1),

		DAOForkBlock:   big.NewInt(1),
		DAOForkSupport: true,

		EIP150Block: big.NewInt(1),
		EIP150Hash:  common.Hash{},

		EIP155Block: big.NewInt(1),
		EIP158Block: big.NewInt(1),

		ByzantiumBlock:      big.NewInt(1),
		ConstantinopleBlock: big.NewInt(1),
		PetersburgBlock:     big.NewInt(1),
		IstanbulBlock:       big.NewInt(1),
		MuirGlacierBlock:    big.NewInt(1),
		BerlinBlock:         big.NewInt(1),
		LondonBlock:         big.NewInt(1),
	}
}

type EVMTransaction struct {
	stateDB      *CommitStateDB
	gaspool      *ethcore.GasPool
	header       *abci.Header
	from         keys.Address
	to           *keys.Address
	nonce        uint64
	value        *big.Int
	data         []byte
	accessList   *ethtypes.AccessList
	gas          uint64
	gasPrice     *big.Int
	isSimulation bool

	// debug olvm
	debug bool
}

func NewEVMTransaction(stateDB *CommitStateDB, gaspool *ethcore.GasPool, header *abci.Header, from keys.Address, to *keys.Address, nonce uint64, value *big.Int, data []byte, accessList *ethtypes.AccessList, gas uint64, gasPrice *big.Int, isSimulation bool) *EVMTransaction {
	return &EVMTransaction{
		stateDB:      stateDB,
		gaspool:      gaspool,
		header:       header,
		from:         from,
		to:           to,
		nonce:        nonce,
		value:        value,
		data:         data,
		isSimulation: isSimulation,
		gas:          gas,
		gasPrice:     gasPrice,
		accessList:   accessList,
	}
}

func (etx *EVMTransaction) SetVMDebug(debug bool) {
	etx.debug = debug
}

func (etx *EVMTransaction) NewEVM() *ethvm.EVM {
	blockCtx := ethvm.BlockContext{
		CanTransfer: ethcore.CanTransfer,
		Transfer:    ethcore.Transfer,
		GetHash:     GetHashFn(etx.stateDB, etx.header),
		Coinbase:    ethcmn.BytesToAddress(etx.header.ProposerAddress),
		GasLimit:    etx.gaspool.Gas(),
		BlockNumber: new(big.Int).SetInt64(etx.header.GetHeight()),
		Time:        new(big.Int).SetInt64(etx.header.Time.Unix()),
		Difficulty:  new(big.Int).Set(DefaultDifficulty), // 0 or 1, does not matter, api show 1, so let say it here as 1
	}

	ethConfig := EthereumConfig(etx.header.ChainID)
	vmConfig := ethvm.Config{
		ExtraEips: make([]int, 0),
	}
	if etx.debug {
		vmConfig.Debug = true
		vmConfig.Tracer = ethvm.NewMarkdownLogger(&ethvm.LogConfig{
			Debug:     true,
			Overrides: ethConfig,
		}, os.Stdout)
	}

	txCtx := ethvm.TxContext{
		Origin:   etx.From(),
		GasPrice: etx.gasPrice,
	}

	return ethvm.NewEVM(blockCtx, txCtx, etx.stateDB, ethConfig, vmConfig)
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

func (etx *EVMTransaction) AccessList() ethtypes.AccessList {
	if etx.accessList == nil {
		return make(ethtypes.AccessList, 0)
	}
	return *etx.accessList
}

func (etx *EVMTransaction) IsFake() bool {
	return etx.isSimulation
}

func (etx *EVMTransaction) Nonce() uint64 {
	return etx.nonce
}

func (etx *EVMTransaction) Data() []byte {
	return etx.data
}

func (etx *EVMTransaction) Value() *big.Int {
	return etx.value
}

func (etx *EVMTransaction) Gas() uint64 {
	return etx.gas
}

func (etx *EVMTransaction) GasPrice() *big.Int {
	return etx.gasPrice
}

func (etx *EVMTransaction) Apply() (*ExecutionResult, error) {
	executionResult, err := ApplyMessage(etx.NewEVM(), etx, etx.gaspool)

	if !etx.IsFake() {
		// Ensure any modifications are committed to the state
		if err := etx.stateDB.Finalise(true); err != nil {
			return nil, err
		}
	}
	return executionResult, err
}
