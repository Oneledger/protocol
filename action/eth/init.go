package eth

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
)

const (
	totalETHSupply     = "10000000000000000000" // 10 ETH
	lockBalanceAddress = "13371337"
	totalTTCSupply        = "10000000000000000000" // 10 ETH
	TTClockBalanceAddress = "23372337"
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
	return nil
}

func EnableInternalETH(r action.Router) error {
	err := r.AddHandler(action.ETH_REPORT_FINALITY_MINT, reportFinalityMintTx{})
	if err != nil {
		return errors.Wrap(err, "reportFinaityMintTx")
	}
	return nil
}
