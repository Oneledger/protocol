package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
)

type TrackerCtx struct {
	tracker      *Tracker
	trackerStore *TrackerStore
	jobStore     *jobs.JobStore
	currNodeAddr keys.Address
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore) TrackerCtx {
	return TrackerCtx{
		tracker:      t,
		currNodeAddr: addr,
		jobStore:     js,
		trackerStore: ts,
	}
}
