/*

 */

package btc

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
)

func EnableBTC(r action.Router) error {
	err := r.AddHandler(action.BTC_LOCK, btcLockTx{})
	if err != nil {
		return errors.Wrap(err, "btcLockTx")
	}

	err = r.AddHandler(action.BTC_ADD_SIGNATURE, btcAddSignatureTx{})
	if err != nil {
		return errors.Wrap(err, "btcAddSignatureTx")
	}

	err = r.AddHandler(action.BTC_REPORT_FINALITY_MINT, reportFinalityMintTx{})
	if err != nil {
		return errors.Wrap(err, "reportFinalityMintTx")
	}

	err = r.AddHandler(action.BTC_EXT_MINT, extMintOBTCTx{})
	if err != nil {
		return errors.Wrap(err, "extMintOBTCTx")
	}

	return nil
}
