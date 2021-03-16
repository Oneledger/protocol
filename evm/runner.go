package evm

import (
	"errors"
	"math/big"
	"strings"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	ethparams "github.com/ethereum/go-ethereum/params"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

// AbciHeaderToTendermint is a util function to parse a tendermint ABCI Header to
// tendermint types Header.
func AbciHeaderToTendermint(header abci.Header) tmtypes.Header {
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
func HashFromContext(ctx *action.Context) ethcmn.Hash {
	// cast the ABCI header to tendermint Header type
	tmHeader := AbciHeaderToTendermint(*ctx.Header)

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
func GetHashFn(ctx *action.Context, csdb *CommitStateDB) ethvm.GetHashFunc {
	return func(height uint64) ethcmn.Hash {
		switch {
		case ctx.Header.GetHeight() == int64(height):
			// Case 1: The requested height matches the one from the context so we can retrieve the header
			// hash directly from the context.
			return HashFromContext(ctx)

		case ctx.Header.Time.Unix() > int64(height):
			// Case 2: if the chain is not the current height we need to retrieve the hash from the store for the
			// current chain epoch. This only applies if the current height is greater than the requested height.
			return csdb.WithContext(ctx).GetHeightHash(height)

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

	return &params.ChainConfig{
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

func NewEVM(ctx *action.Context, ecfg *EVMConfig) *ethvm.EVM {
	s := NewCommitStateDB(ctx)
	blockCtx := ethvm.BlockContext{
		CanTransfer: ethcore.CanTransfer,
		Transfer:    ethcore.Transfer,
		GetHash:     GetHashFn(ctx, s),
		Coinbase:    ethcmn.Address{}, // there's no beneficiary since we're not mining
		GasLimit:    ecfg.gasLimit,    // TODO: Set gas limit from the settings
		BlockNumber: big.NewInt(ctx.Header.GetHeight()),
		Time:        big.NewInt(ctx.Header.Time.Unix()),
		Difficulty:  big.NewInt(0), // unused. Only required in PoW context
	}

	vmConfig := ethvm.Config{
		ExtraEips: ecfg.extraEIPs,
	}

	txCtx := ethvm.TxContext{
		Origin:   ethcmn.BytesToAddress(ecfg.addr),
		GasPrice: ecfg.gasPrice,
	}

	ethConfig, _ := EthereumConfig(ctx.Header.ChainID)
	return ethvm.NewEVM(blockCtx, txCtx, s, ethConfig, vmConfig)
}
