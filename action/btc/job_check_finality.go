/*

 */

package btc

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/bitcoin"
)

type JobBTCCheckFinality struct {
	TrackerName string

	JobID string

	TransactionConfirmed bool

	Done bool
}

func (j *JobBTCCheckFinality) DoMyJob(ctxI interface{}) {

	ctx, _ := ctxI.(action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)
	ok, _ := cd.CheckFinality(tracker.ProcessTxId)
	if ok {
		// todo
	}

}

func (j *JobBTCCheckFinality) IsMyJobDone(ctxI interface{}) bool {
	panic("implement me")
}

func (j *JobBTCCheckFinality) IsSufficient(ctxI interface{}) bool {
	panic("implement me")
}

func (j *JobBTCCheckFinality) DoFinalize() {

}

// simple getters
func (j *JobBTCCheckFinality) GetType() string {
	return JobTypeBTCCheckFinality
}

func (j *JobBTCCheckFinality) GetJobID() string {
	return j.JobID
}

func (j *JobBTCCheckFinality) IsDone() bool {
	return j.Done
}
