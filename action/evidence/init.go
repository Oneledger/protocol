package evidence

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(Release), "release")
	serialize.RegisterConcrete(new(Allegation), "allegation")
	serialize.RegisterConcrete(new(AllegationVote), "allegation_vote")
}

func EnablePenalization(r action.Router) error {

	err := r.AddHandler(action.RELEASE, releaseTx{})
	if err != nil {
		return errors.Wrap(err, "releaseTx")
	}
	err = r.AddHandler(action.ALLEGATION, allegationTx{})
	if err != nil {
		return errors.Wrap(err, "allegationTx")
	}
	err = r.AddHandler(action.ALLEGATION_VOTE, allegationVoteTx{})
	if err != nil {
		return errors.Wrap(err, "allegationVoteTx")
	}
	return nil
}
