package rewards

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
)

func EnableRewards(r action.Router) error {
	err := r.AddHandler(action.WITHDRAW_REWARD, withdrawTx{})
	if err != nil {
		return errors.Wrap(err, "withdrawRewardsTx")
	}

	return nil
}
