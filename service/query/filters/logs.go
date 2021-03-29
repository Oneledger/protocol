package filters

import (
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// filterLogs creates a slice of logs matching the given criteria.
// [] -> anything
// [A] -> A in first position of log topics, anything after
// [null, B] -> anything in first position, B in second position
// [A, B] -> A in first position and B in second position
// [[A, B], [A, B]] -> A or B in first position, A or B in second position
func FilterLogs(logs []*ethtypes.Log, fromBlock, toBlock *big.Int, addresses []ethcmn.Address, topics [][]ethcmn.Hash) []*ethtypes.Log {
	var ret []*ethtypes.Log
Logs:
	for _, log := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			continue
		}
		if len(addresses) > 0 && !includes(addresses, log.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(log.Topics) {
			continue
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if log.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}

func includes(addresses []ethcmn.Address, a ethcmn.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}

func bloomFilter(bloom ethtypes.Bloom, addresses []ethcmn.Address, topics [][]ethcmn.Hash) bool {
	var included bool
	if len(addresses) > 0 {
		for _, addr := range addresses {
			if ethtypes.BloomLookup(bloom, addr) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	for _, sub := range topics {
		included = len(sub) == 0 // empty rule set == wildcard
		for _, topic := range sub {
			if ethtypes.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
	}
	return included
}

// returnHashes is a helper that will return an empty hash array case the given hash array is nil,
// otherwise the given hashes array is returned.
func returnHashes(hashes []ethcmn.Hash) []ethcmn.Hash {
	if hashes == nil {
		return []ethcmn.Hash{}
	}
	return hashes
}

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*ethtypes.Log) []*ethtypes.Log {
	if logs == nil {
		return []*ethtypes.Log{}
	}
	return logs
}
