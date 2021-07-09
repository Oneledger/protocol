package eth

import (
	rpcfilters "github.com/Oneledger/protocol/web3/eth/filters"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethfilters "github.com/ethereum/go-ethereum/eth/filters"
)

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*ethtypes.Log) []*ethtypes.Log {
	if logs == nil {
		return []*ethtypes.Log{}
	}
	return logs
}

// GetLogs returns logs matching the given argument that are stored within the state.
func (svc *Service) GetLogs(crit ethfilters.FilterCriteria) ([]*ethtypes.Log, error) {
	var filter *rpcfilters.Filter

	svc.logger.Debug("eth_getLogs", "crit", crit)

	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = rpcfilters.NewBlockFilter(svc, *crit.BlockHash, crit.Addresses, crit.Topics)
	} else {
		result, err := svc.getTMClient().Block(nil)
		if err != nil {
			return nil, err
		}
		if result.Block == nil {
			return nil, nil
		}
		// Convert the RPC block numbers into internal representations
		begin := result.Block.Height
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := result.Block.Height
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = rpcfilters.NewRangeFilter(svc, begin, end, crit.Addresses, crit.Topics)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(svc.ctx)
	if err != nil {
		return nil, err
	}
	svc.logger.Debug("eth_getLogs", "count", len(logs))
	return returnLogs(logs), err
}
