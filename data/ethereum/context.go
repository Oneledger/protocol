package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

type TrackerCtx struct {
	Tracker      *Tracker
	TrackerStore *TrackerStore
	JobStore     *jobs.JobStore
	CurrNodeAddr keys.Address
	Validators   *identity.ValidatorStore
	Logger       *log.Logger
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore, vs *identity.ValidatorStore, log *log.Logger) *TrackerCtx {
	return &TrackerCtx{
		Tracker:      t,
		CurrNodeAddr: addr,
		JobStore:     js,
		TrackerStore: ts,
		Validators:   vs,
		Logger:       log,
	}
}
