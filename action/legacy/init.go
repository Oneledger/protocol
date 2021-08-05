package legacy

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(LegacyConnect), "legacy_connect")
	serialize.RegisterConcrete(new(LegacySend), "legacy_send")

}

func EnableLegacy(r action.Router) error {

	err := r.AddHandler(action.LEGACY_CONNECT, legacyConnectTx{})
	if err != nil {
		return errors.Wrap(err, "legacyConnectTx")
	}
	err = r.AddHandler(action.LEGACY_SEND, legacySendTx{})
	if err != nil {
		return errors.Wrap(err, "legacySendTx")
	}
	return nil
}
