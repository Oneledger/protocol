package transfer

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(Send), "action_send")

}

func EnableSend(r action.Router) error {

	err := r.AddHandler(action.SEND, sendTx{})
	if err != nil {
		return errors.Wrap(err, "sendTx")
	}
	err = r.AddHandler(action.SENDPOOL, sendPoolTx{})
	if err != nil {
		return errors.Wrap(err, "sendPoolTx")
	}
	return nil
}
