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
	Witnesses    *identity.EthWitnessStore
	Logger       *log.Logger
}

func NewTrackerCtx(t *Tracker, addr keys.Address, js *jobs.JobStore, ts *TrackerStore, ws *identity.EthWitnessStore, log *log.Logger) *TrackerCtx {
	return &TrackerCtx{
		Tracker:      t,
		CurrNodeAddr: addr,
		JobStore:     js,
		TrackerStore: ts,
		Witnesses:    ws,
		Logger:       log,
	}
}
