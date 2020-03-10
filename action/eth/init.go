//Package for transactions related to Etheruem
package eth

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
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

	err = r.AddHandler(action.ETH_REDEEM, ethRedeemTx{})
	if err != nil {
		return errors.Wrap(err, "ethRedeemTx")
	}

	err = r.AddHandler(action.ERC20_LOCK, ethERC20LockTx{})
	if err != nil {
		return errors.Wrap(err, "ERC20LockTx")
	}

	err = r.AddHandler(action.ERC20_REDEEM, ethERC20RedeemTx{})
	if err != nil {
		return errors.Wrap(err, "ERC20Redeem)")
	}
	return nil
}

func EnableInternalETH(r action.Router) error {
	err := r.AddHandler(action.ETH_REPORT_FINALITY_MINT, reportFinalityMintTx{})
	if err != nil {
		return errors.Wrap(err, "reportFinaityMintTx")
	}
	return nil
}
