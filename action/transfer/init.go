package transfer

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

func init() {

	serialize.RegisterConcrete(new(Send), "action_send")

}

func EnableSend(r action.Router) error {

	err := r.AddHandler(action.SEND, sendTx{})
	if err != nil {
		return errors.Wrap(err, "sendTx")
	}
	return nil
}
