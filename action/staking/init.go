package staking

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(ApplyValidator), "action_av")
}

func EnableApplyValidator(r action.Router) error {

	err := r.AddHandler(action.APPLYVALIDATOR, applyTx{})
	if err != nil {
		return errors.Wrap(err, "applyTx")
	}

	err = r.AddHandler(action.WITHDRAW, withdrawTx{})
	if err != nil {
		return errors.Wrap(err, "withdrawTx")
	}
	return nil
}

func EnablePurgeValidator(r action.Router) error {
	err := r.AddHandler(action.PURGE, purgeTx{})
	if err != nil {
		return errors.Wrap(err, "purgeTx")
	}
	return nil
}
