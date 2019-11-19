package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
)

type TrackerCtx struct {
	Tracker      *Tracker
	TrackerStore *TrackerStore
	JobStore     *jobs.JobStore
	CurrNodeAddr keys.Address
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore) TrackerCtx {
	return TrackerCtx{
		Tracker:      t,
		CurrNodeAddr: addr,
		JobStore:     js,
		TrackerStore: ts,
	}
}
