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
	err = r.AddHandler(action.BTC_REDEEM, btcRedeemTx{})
	if err != nil {
		return errors.Wrap(err, "btcRedeemTx")
	}

	err = r.AddHandler(action.BTC_ADD_SIGNATURE, &btcAddSignatureTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_BROADCAST_SUCCESS, &btcBroadcastSuccessTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_REPORT_FINALITY_MINT, &reportFinalityMintTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_FAILED_BROADCAST_RESET, &btcBroadcastFailureReset{})
	if err != nil {
		return err
	}

	return nil
}

func EnableBTCInternalTx(r action.Router) error {
	err := r.AddHandler(action.BTC_ADD_SIGNATURE, &btcAddSignatureTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_BROADCAST_SUCCESS, &btcBroadcastSuccessTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_REPORT_FINALITY_MINT, &reportFinalityMintTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_FAILED_BROADCAST_RESET, &btcBroadcastFailureReset{})
	if err != nil {
		return err
	}

	return nil
}
