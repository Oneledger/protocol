package bitcoin

import (
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/identity"
)

type BTCTransitionContext struct {
	Tracker    *Tracker
	JobStore   *jobs.JobStore
	Validators *identity.ValidatorStore
}

func NewTrackerCtx(tracker *Tracker, js *jobs.JobStore, vs *identity.ValidatorStore) *BTCTransitionContext {
	return &BTCTransitionContext{
		Tracker:    tracker,
		JobStore:   js,
		Validators: vs,
	}
}
