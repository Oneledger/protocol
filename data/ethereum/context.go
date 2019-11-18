package ethereum

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
)

type TrackerCtx struct {
	tracker      *Tracker
	jobStore     *jobs.JobStore
	currNodeAddr keys.Address
}

func NewTrackerCtx(t *Tracker, currAddr keys.Address, jobStore *jobs.JobStore) TrackerCtx {
	return TrackerCtx{
		tracker:      t,
		currNodeAddr: currAddr,
		jobStore:     jobStore,
	}
}
