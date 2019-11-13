/*

 */

package btc

import "github.com/Oneledger/protocol/action"

func EnableBTCInternalTx(r action.Router) error {
	err := r.AddHandler(action.BTC_ADD_SIGNATURE, &btcAddSignatureTx{})
	if err != nil {
		return err
	}

	err = r.AddHandler(action.BTC_REPORT_FINALITY_MINT, &reportFinalityMintTx{})
	return err
}
