package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

type TrackerCtx struct {
	Tracker         *Tracker
	TrackersOngoing *TrackerStore
	TrackersFailed  *TrackerStore
	JobStore        *jobs.JobStore
	CurrNodeAddr    keys.Address
	Validators      *identity.ValidatorStore
	Logger          *log.Logger
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore, tsf *TrackerStore, vs *identity.ValidatorStore, log *log.Logger) *TrackerCtx {
	return &TrackerCtx{
		Tracker:         t,
		CurrNodeAddr:    addr,
		JobStore:        js,
		TrackersOngoing: ts,
		TrackersFailed:  tsf,
		Validators:      vs,
		Logger:          log,
	}
}
