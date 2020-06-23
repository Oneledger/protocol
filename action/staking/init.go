package staking

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(Stake), "stake")
	serialize.RegisterConcrete(new(Unstake), "unstake")
	serialize.RegisterConcrete(new(Withdraw), "withdraw")
}

func EnableStaking(r action.Router) error {

	err := r.AddHandler(action.STAKE, stakeTx{})
	if err != nil {
		return errors.Wrap(err, "stakeTx")
	}

	err = r.AddHandler(action.UNSTAKE, unstakeTx{})
	if err != nil {
		return errors.Wrap(err, "unstakeTx")
	}

	err = r.AddHandler(action.WITHDRAW, withdrawTx{})
	if err != nil {
		return errors.Wrap(err, "withdrawTx")
	}
	return nil
}
