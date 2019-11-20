package bitcoin

import "github.com/Oneledger/protocol/data/jobs"

type BTCTransitionContext struct {
	Tracker  *Tracker
	JobStore *jobs.JobStore
}
