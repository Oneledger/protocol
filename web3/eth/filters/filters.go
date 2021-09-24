package filters

import (
	"math/big"
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/vm"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

type Filter struct {
	svc     rpctypes.EthService
	stateDB *vm.CommitStateDB
	logger  *log.Logger

	addresses []common.Address
	topics    [][]common.Hash

	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewBlockFilter(svc rpctypes.EthService, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Create a generic filter and convert it into a block filter
	filter := newFilter(svc, addresses, topics)
	filter.block = block
	filter.logger.Debug("NewBlockFilter", "block", block)
	return filter
}

// NewRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func NewRangeFilter(svc rpctypes.EthService, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
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
	filter := newFilter(svc, addresses, topics)
	filter.begin = begin
	filter.end = end
	filter.logger.Debug("NewRangeFilter", "begin", filter.begin, "end", filter.end)
	return filter
}

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(svc rpctypes.EthService, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		svc:       svc,
		logger:    log.NewLoggerWithPrefix(os.Stdout, "filters"),
		addresses: addresses,
		topics:    topics,
		stateDB:   svc.GetStateDB(),
	}
}

// BloomStatus returns the BloomBitsBlocks and the number of processed sections maintained
// by the chain indexer.
func (f *Filter) BloomStatus() (uint64, uint64) {
	return 4096, 0
}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs(_ rpctypes.Web3Context) ([]*ethtypes.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != (common.Hash{}) {
		block, err := f.svc.GetBlockByHash(f.block, false)
		if block == nil {
			return nil, err
		}
		return f.blockLogs(block)
	}
	lastHeight := ethrpc.LatestBlockNumber
	block, _ := f.svc.GetBlockByNumber(ethrpc.BlockNumberOrHash{BlockNumber: &lastHeight}, false)
	if block == nil {
		return nil, nil
	}
	head := block.Number

	if f.begin == -1 {
		f.begin = int64(head)
	}
	if f.end == -1 {
		f.end = int64(head)
	}
	// Gather all indexed logs, and finish with non indexed ones
	logs := []*ethtypes.Log{}
	for i := f.begin; i <= f.end; i++ {
		blockNum := ethrpc.BlockNumber(i)
		block, err := f.svc.GetBlockByNumber(ethrpc.BlockNumberOrHash{BlockNumber: &blockNum}, false)
		if block == nil {
			return logs, err
		}

		if len(block.Transactions) == 0 {
			continue
		}

		logsMatched := f.checkMatches(block.Hash, block.Transactions)
		logs = append(logs, logsMatched...)
	}

	return logs, nil
}

// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(block *rpctypes.Block) (logs []*types.Log, err error) {
	f.logger.Debug("blockLogs", "block height", block.Number)
	if !bloomFilter(block.LogsBloom, f.addresses, f.topics) {
		return []*ethtypes.Log{}, nil
	}
	var logsList = [][]*ethtypes.Log{}

	blockLogs, err := f.stateDB.GetLogs(block.Hash)
	if err != nil {
		return []*ethtypes.Log{}, err
	}

	f.logger.Debug("blockLogs", "iterate txs", len(block.Transactions))
	for _, itx := range block.Transactions {
		txHash, ok := itx.(common.Hash)
		if !ok {
			continue
		}
		if logs, ok := blockLogs.Logs[txHash]; ok {
			logsList = append(logsList, logs)
		}
	}
	f.logger.Debug("blockLogs", "check unfiltered", len(logsList))
	unfiltered := []*ethtypes.Log{}
	for _, logs := range logsList {
		unfiltered = append(unfiltered, logs...)
	}
	f.logger.Debug("blockLogs", "unfiltered", len(unfiltered))
	fLogs := filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
	f.logger.Debug("blockLogs", "filtered", len(fLogs))
	if len(fLogs) == 0 {
		return []*ethtypes.Log{}, nil
	}
	return fLogs, nil
}

// checkMatches checks if the logs from the a list of transactions transaction
// contain any log events that  match the filter criteria. This function is
// called when the bloom filter signals a potential match.
func (f *Filter) checkMatches(bHash common.Hash, transactions []interface{}) []*ethtypes.Log {
	unfiltered := []*ethtypes.Log{}

	blockLogs, err := f.stateDB.GetLogs(bHash)
	if err != nil {
		return []*ethtypes.Log{}
	}
	for _, itx := range transactions {
		txHash, ok := itx.(common.Hash)
		if !ok {
			continue
		}
		logs, ok := blockLogs.Logs[txHash]
		if !ok {
			// ignore error if transaction didn't set any logs (eg: when tx type is not
			// MsgEthereumTx or MsgEthermint)
			continue
		}

		unfiltered = append(unfiltered, logs...)
	}

	return filterLogs(unfiltered, big.NewInt(f.begin), big.NewInt(f.end), f.addresses, f.topics)
}
