package utils

import (
	"bytes"
	"fmt"
	"math/big"

	rpcclient "github.com/Oneledger/protocol/client"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// GetPendingTxCountByAddress is used to get pending tx count (nonce) for user address
// NOTE: Working right now only with legacy tx
func GetPendingTxCountByAddress(tmClient rpcclient.Client, address common.Address) (total uint64) {
	unconfirmed, err := tmClient.UnconfirmedTxs(1000)
	if err != nil {
		return 0
	}
	for _, tx := range unconfirmed.Txs {
		lTx, err := rpctypes.ParseLegacyTx(tx)
		if err != nil {
			// means tx is not legacy and we need to check is tx is ethereum
			// TODO: Add ethereum tx check when it will be released
			continue
		}
		// This is only for legacy tx
		for _, sig := range lTx.Signatures {
			pubKeyHandler, err := sig.Signer.GetHandler()
			if err != nil {
				continue
			}
			// match if signer is a user
			if pubKeyHandler.Address().Equal(address.Bytes()) {
				total++
			}
		}
	}
	return
}

// GetPendingTx search for tx in pool
func GetPendingTx(tmClient rpcclient.Client, hash common.Hash, chainID *big.Int) (*rpctypes.Transaction, error) {
	unconfirmed, err := tmClient.UnconfirmedTxs(1000)
	if err != nil {
		return nil, err
	}

	for _, uTx := range unconfirmed.Txs {
		if bytes.Equal(uTx.Hash(), hash.Bytes()) {
			return rpctypes.LegacyRawBlockAndTxToEthTx(nil, &uTx, chainID, nil)
		}
	}
	return nil, err
}

// GetPendingTransactions search for txs in pool
func GetPendingTxs(tmClient rpcclient.Client, chainID *big.Int) ([]*rpctypes.Transaction, error) {
	unconfirmed, err := tmClient.UnconfirmedTxs(1000)
	if err != nil {
		return nil, err
	}

	transactions := make([]*rpctypes.Transaction, 0)

	for _, uTx := range unconfirmed.Txs {
		tx, err := rpctypes.LegacyRawBlockAndTxToEthTx(nil, &uTx, chainID, nil)
		if err != nil {
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// GetPendingTxsWithCallback search for txs in pool and return in callback form
func GetPendingTxsWithCallback(tmClient rpcclient.Client, chainID *big.Int, callback func(tx *rpctypes.Transaction) bool) error {
	unconfirmed, err := tmClient.UnconfirmedTxs(1000)
	if err != nil {
		return err
	}

	for _, uTx := range unconfirmed.Txs {
		tx, err := rpctypes.LegacyRawBlockAndTxToEthTx(nil, &uTx, chainID, nil)
		if err != nil {
			continue
		}
		stopped := callback(tx)
		if stopped {
			break
		}
	}
	return nil
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func CheckTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.Ether)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f OLT) exceeds the configured cap (%.2f OLT)", feeFloat, cap)
	}
	return nil
}
