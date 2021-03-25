package action

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Oneledger/protocol/data/keys"
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

// HashFromContext returns the Ethereum Header hash from the context's protocol
// block header.
func HashFromHeader(header *abci.Header) ethcmn.Hash {
	// cast the ABCI header to tendermint Header type
	tmHeader := AbciHeaderToTendermint(header)

	// get the Tendermint block hash from the current header
	tmBlockHash := tmHeader.Hash()

	// NOTE: if the validator set hash is missing the hash will be returned as nil,
	// so we need to check for this case to prevent a panic when calling Bytes()
	if tmBlockHash == nil {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(tmBlockHash.Bytes())
}

// GetHashFn implements vm.GetHashFunc for protocol. It handles 3 cases:
//  1. The requested height matches the current height from context (and thus same epoch number)
//  2. The requested height is from an previous height from the same chain epoch
//  3. The requested height is from a height greater than the latest one
func GetHashFn(header *abci.Header, s *CommitStateDB) ethvm.GetHashFunc {
	return func(height uint64) ethcmn.Hash {
		switch {
		case header.GetHeight() == int64(height):
			// Case 1: The requested height matches the one from the context so we can retrieve the header
			// hash directly from the context.
			return HashFromHeader(header)

		case header.GetHeight() > int64(height):
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			return s.GetHeightHash(height)

		default:
			// Case 3: heights greater than the current one returns an empty hash.
			return ethcmn.Hash{}
		}
	}
}

func EthereumConfig(chainID string) (*ethparams.ChainConfig, error) {
	chainData := strings.Split(chainID, "-")
	id := new(big.Int)
	id, ok := id.SetString(chainData[1], 10)
	if !ok {
		return nil, errors.New("chainId is wrong")
	}

	return &ethparams.ChainConfig{
		ChainID:        id,
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

		YoloV3Block: nil,
		EWASMBlock:  nil,
	}, nil
}

type EVMConfig struct {
	addr      keys.Address
	gasPrice  *big.Int
	gasLimit  uint64
	extraEIPs []int
}

func NewEVMConfig(addr keys.Address, gasPrice *big.Int, gasLimit uint64, extraEIPs []int) *EVMConfig {
	return &EVMConfig{
		addr:      addr,
		gasPrice:  gasPrice,
		gasLimit:  gasLimit,
		extraEIPs: extraEIPs,
	}
}

type EVMTransaction struct {
	stateDB *CommitStateDB
	header  *abci.Header
	from    keys.Address
	to      *keys.Address
	value   *big.Int
	data    []byte
	ecfg    *EVMConfig
}

var (
	DefaultGasLimit uint64   = 10_000_000
	DefaultGasPrice *big.Int = big.NewInt(0)
)

func NewEVMTransaction(stateDB *CommitStateDB, header *abci.Header, from keys.Address, to *keys.Address, value *big.Int, data []byte) *EVMTransaction {
	return &EVMTransaction{
		stateDB: stateDB,
		header:  header,
		from:    from,
		to:      to,
		value:   value,
		data:    data,
		// NOTE: Decide what to do with the gas price as we have a fee system
		ecfg: NewEVMConfig(from, DefaultGasPrice, DefaultGasLimit, make([]int, 0)),
	}
}

func (etx *EVMTransaction) NewEVM() *ethvm.EVM {
	blockCtx := ethvm.BlockContext{
		CanTransfer: ethcore.CanTransfer,
		Transfer:    ethcore.Transfer,
		GetHash:     GetHashFn(etx.header, etx.stateDB),
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

	ethConfig, _ := EthereumConfig(etx.header.ChainID)
	return ethvm.NewEVM(blockCtx, txCtx, etx.stateDB, ethConfig, vmConfig)
}

func (etx *EVMTransaction) Origin() ethcmn.Address {
	return ethcmn.BytesToAddress(etx.ecfg.addr)
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

func (etx *EVMTransaction) GetLastNonce() uint64 {
	return etx.stateDB.GetNonce(etx.From())
}

func (etx *EVMTransaction) Apply(vmenv *ethvm.EVM, tx RawTx) (*ethcore.ExecutionResult, error) {
	fmt.Printf("etx.From(): %s\n", etx.From())
	fmt.Printf("etx.To(): %s\n", etx.To())
	nonce := etx.GetLastNonce()
	msg := ethtypes.NewMessage(etx.From(), etx.To(), nonce, etx.value, uint64(tx.Fee.Gas), etx.ecfg.gasPrice, etx.data, make(ethtypes.AccessList, 0), true)

	statedb := etx.stateDB
	statedb.Prepare(ethcmn.Hash{}, ethcmn.Hash{}, 0)
	msgResult, err := ethcore.ApplyMessage(vmenv, msg, new(ethcore.GasPool).AddGas(uint64(uint64(tx.Fee.Gas))))
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}
	fmt.Printf("msgResult: %v\n", msgResult)
	// Ensure any modifications are committed to the state
	// Only delete empty objects if EIP158/161 (a.k.a Spurious Dragon) is in effect
	if err := statedb.Finalise(vmenv.ChainConfig().IsEIP158(big.NewInt(etx.header.Height))); err != nil {
		return nil, fmt.Errorf("failed to finalize: %v", err)
	}
	return msgResult, nil
}
