package filters

import (
	"math/big"
	"os"

	"github.com/Oneledger/protocol/log"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
	tmcoretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/store"
)

type Filter struct {
	logger     *log.Logger
	blockStore *store.BlockStore

	addresses []common.Address
	topics    [][]common.Hash

	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewBlockFilter(blockStore *store.BlockStore, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Create a generic filter and convert it into a block filter
	filter := newFilter(blockStore, addresses, topics)
	filter.block = block
	filter.logger.Debug("NewBlockFilter", "block", block)
	return filter
}

// NewRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func NewRangeFilter(blockStore *store.BlockStore, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Flatten the address and topic filter clauses into a single bloombits filter
	// system. Since the bloombits are not positional, nil topics are permitted,
	// which get flattened into a nil byte slice.
	var filters [][][]byte
	if len(addresses) > 0 {
		filter := make([][]byte, len(addresses))
		for i, address := range addresses {
			filter[i] = address.Bytes()
		}
		filters = append(filters, filter)
	}
	for _, topicList := range topics {
		filter := make([][]byte, len(topicList))
		for i, topic := range topicList {
			filter[i] = topic.Bytes()
		}
		filters = append(filters, filter)
	}
	// Create a generic filter and convert it into a range filter
	filter := newFilter(blockStore, addresses, topics)
	filter.begin = begin
	filter.end = end
	filter.logger.Debug("NewRangeFilter", "begin", filter.begin, "end", filter.end)
	return filter
}

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(blockStore *store.BlockStore, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		logger:     log.NewLoggerWithPrefix(os.Stdout, "filters"),
		blockStore: blockStore,
		addresses:  addresses,
		topics:     topics,
	}
}

// BloomStatus returns the BloomBitsBlocks and the number of processed sections maintained
// by the chain indexer.
func (f *Filter) BloomStatus() (uint64, uint64) {
	return 4096, 0
}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs() ([]*ethtypes.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != (common.Hash{}) {
		block := f.blockStore.LoadBlockByHash(f.block.Bytes())
		if block == nil {
			return returnLogs(nil), nil
		}
		blockResults, err := tmrpccore.BlockResults(nil, &block.Height)
		if err != nil {
			return returnLogs(nil), nil
		}
		return f.blockLogs(blockResults)
	}
	lastHeight := f.blockStore.Height()
	block := f.blockStore.LoadBlock(lastHeight)
	if block == nil {
		return returnLogs(nil), nil
	}
	head := block.Height

	if f.begin == -1 {
		f.begin = int64(head)
	}
	if f.end == -1 {
		f.end = int64(head)
	}
	// Gather all indexed logs, and finish with non indexed ones
	logs := []*ethtypes.Log{}
	for i := f.begin; i <= f.end; i++ {
		blockResults, err := tmrpccore.BlockResults(nil, &i)
		if err != nil {
			return logs, nil
		}

		if len(blockResults.TxsResults) == 0 {
			continue
		}

		logsMatched := f.checkMatches(blockResults.TxsResults)
		logs = append(logs, logsMatched...)
	}

	return logs, nil
}

// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(blockResults *tmcoretypes.ResultBlockResults) (logs []*types.Log, err error) {
	f.logger.Debug("blockLogs", "block height", blockResults.Height)

	bloom := rpctypes.GetBlockBloom(blockResults.EndBlockEvents)
	if !bloomFilter(bloom, f.addresses, f.topics) {
		return []*ethtypes.Log{}, nil
	}
	var logsList = [][]*ethtypes.Log{}

	f.logger.Debug("blockLogs", "iterate txs", len(blockResults.TxsResults))
	for index, tx := range blockResults.TxsResults {
		logReceipt := rpctypes.GetTxEthLogs(tx, uint32(index))
		logsList = append(logsList, logReceipt.Logs)
	}
	f.logger.Debug("blockLogs", "check unfiltered", len(logsList))
	unfiltered := []*ethtypes.Log{}
	for _, logs := range logsList {
		unfiltered = append(unfiltered, logs...)
	}
	f.logger.Debug("blockLogs", "unfiltered", len(unfiltered))
	fLogs := FilterLogs(unfiltered, nil, nil, f.addresses, f.topics)
	f.logger.Debug("blockLogs", "filtered", len(fLogs))
	if len(fLogs) == 0 {
		return []*ethtypes.Log{}, nil
	}
	return fLogs, nil
}

// checkMatches checks if the logs from the a list of transactions transaction
// contain any log events that  match the filter criteria. This function is
// called when the bloom filter signals a potential match.
func (f *Filter) checkMatches(transactions []*abci.ResponseDeliverTx) []*ethtypes.Log {
	unfiltered := []*ethtypes.Log{}

	for index, tx := range transactions {
		logReceipt := rpctypes.GetTxEthLogs(tx, uint32(index))
		unfiltered = append(unfiltered, logReceipt.Logs...)
	}

	return FilterLogs(unfiltered, big.NewInt(f.begin), big.NewInt(f.end), f.addresses, f.topics)
}
