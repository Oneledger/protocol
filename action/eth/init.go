package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ethereum"

	"github.com/pkg/errors"
)

func EnableETH(r action.Router) error {
	err := r.AddHandler(action.ETH_LOCK, ethLockTx{})
	if err != nil {
		return errors.Wrap(err, "ethLockTx")
	}


	err = r.AddHandler(action.ETH_REPORT_FINALITY_MINT, reportFinalityMintTx{})
	if err != nil {
		return errors.Wrap(err, "reportFinalityMintTx")
	}

	err = r.AddHandler(action.ETH_MINT, ethExtMintTx{})
	if err != nil {
		return errors.Wrap(err, "extMintETHTx")
	}

	return nil
}

func EnableInternalETH(r action.Router) error {
	err := r.AddHandler(action.ETH_REPORT_FINALITY_MINT,reportFinalityMintTx{})
	if err != nil {
		return errors.Wrap(err,"reportFinaityMintTx")
	}
	return nil
}

func GetAmount(tracker *ethereum.Tracker) (*big.Int,error) {
	ethTx := &types.Transaction{}
	err := rlp.DecodeBytes(tracker.SignedETHTx, ethTx)
	if err != nil {
		return nil, errors.Wrap(err, "eth txn decode failed")
	}
	return ethTx.Value(),nil

}