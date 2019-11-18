package ethereum

import "github.com/Oneledger/protocol/data/keys"

type TrackerCtx struct {
	tracker *Tracker
	me      keys.Address
}

func NewTrackerCtx(t *Tracker, me keys.Address) TrackerCtx {
	return TrackerCtx{
		tracker: t,
		me:      me,
	}
}
