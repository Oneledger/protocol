package eth

import (
	"github.com/Oneledger/protocol/action"
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

