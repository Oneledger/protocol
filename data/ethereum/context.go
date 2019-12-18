package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
)

type TrackerCtx struct {
	Tracker      *Tracker
	TrackerStore *TrackerStore
	JobStore     *jobs.JobStore
	CurrNodeAddr keys.Address
	Validators   *identity.ValidatorStore
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore, vs *identity.ValidatorStore) *TrackerCtx {
	return &TrackerCtx{
		Tracker:      t,
		CurrNodeAddr: addr,
		JobStore:     js,
		TrackerStore: ts,
		Validators:   vs,
	}
}
