package nexus

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(Nexus), "nexus")

}

func EnableNexus(r action.Router) error {

	err := r.AddHandler(action.NEXUS, nexusTx{})
	if err != nil {
		return errors.Wrap(err, "nexusTx")
	}
	return nil
}
