/*

 */

package event

import (
	"time"

	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/jobs"
)

type JobBus struct {
	store *jobs.JobStore
	ctx   *JobsContext
	opt   Option
	quit  chan struct{}
}

type Option struct {
	EthInterval time.Duration
	BtcInterval time.Duration
}

func NewJobBus(opt Option, store *jobs.JobStore) *JobBus {
	return &JobBus{
		store: store,
		opt:   opt,
		quit:  make(chan struct{}),
	}
}

func (j *JobBus) Start(ctx *JobsContext) error {
	j.ctx = ctx
	tickerBtc := time.NewTicker(j.opt.BtcInterval)
	tickerEth := time.NewTicker(j.opt.EthInterval)
	go func() {
		for {
			select {
			case <-tickerBtc.C:
				ProcessAllJobs(j.ctx, j.store.WithChain(chain.BITCOIN))
			case <-tickerEth.C:
				ProcessAllJobs(j.ctx, j.store.WithChain(chain.ETHEREUM))
			case <-j.quit:
				tickerBtc.Stop()
				tickerEth.Stop()
				return
			}
		}
	}()

	return nil
}

func (j *JobBus) Close() error {
	close(j.quit)
	return nil
}

type JobProcess func(job jobs.Job) jobs.Job

func ProcessAllJobs(ctx *JobsContext, js *jobs.JobStore) {

	RangeJobs(js, func(job jobs.Job) jobs.Job {
		if !job.IsDone() {
			job.DoMyJob(ctx)
		}
		return job
	})
}

func RangeJobs(js *jobs.JobStore, pro JobProcess) {

	jobkeys := make([]string, 0, 20)
	js.Iterate(func(job jobs.Job) {
		jobkeys = append(jobkeys, job.GetJobID())
	})
	for _, key := range jobkeys {

		job, err := js.GetJob(key)
		if err != nil {
			continue
		}

		job = pro(job)

		js.DeleteJob(job)

	}
}
