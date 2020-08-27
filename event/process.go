/*

 */

package event

import (
	"fmt"
	"runtime/debug"
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
	BtcInterval time.Duration
	EthInterval time.Duration
	OltInterval time.Duration
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
	tickerOlt := time.NewTicker(j.opt.OltInterval)
	go func() {
		for {
			select {
			case <-tickerBtc.C:
				ProcessAllJobs(j.ctx, j.store.WithChain(chain.BITCOIN))
				DeleteCompletedJobs(j.ctx, j.store.WithChain(chain.BITCOIN))
			case <-tickerEth.C:
				ProcessAllJobs(j.ctx, j.store.WithChain(chain.ETHEREUM))
			case <-tickerOlt.C:
				ProcessAllJobs(j.ctx, j.store.WithChain(chain.ONELEDGER))
				DeleteCompletedJobs(j.ctx, j.store.WithChain(chain.ONELEDGER))
			case <-j.quit:
				tickerBtc.Stop()
				tickerEth.Stop()
				tickerOlt.Stop()
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

		//ctx.Logger.Info("Trying to do job : ", job.GetType())

		func() {
			defer func() {
				if r := recover(); r != nil {
					ctx.Logger.Info("panic in job: ", job.GetJobID())
					ctx.Logger.Info(r)
					debug.PrintStack()
					panic(r)
				}
			}()
			job.DoMyJob(ctx)
		}()

		return job
	})
}

func RangeJobs(js *jobs.JobStore, pro JobProcess) {

	jobkeys := make([]string, 0, 20)
	js.Iterate(func(job jobs.Job) {
		//fmt.Println("Searching Jobstore",job.GetType())
		if !job.IsDone() && !job.IsFailed() {
			jobkeys = append(jobkeys, job.GetJobID())
		}
	})
	for _, key := range jobkeys {

		job, err := js.GetJob(key)
		if err != nil {
			fmt.Println("err get job by key", key)
			continue
		}
		job = pro(job)

		err = js.SaveJob(job)
		if err != nil {
			fmt.Println("range job; err saving job", err)
		}

	}
}

func DeleteCompletedJobs(ctx *JobsContext, js *jobs.JobStore) {

	jobkeys := make([]string, 0, 20)
	js.Iterate(func(job jobs.Job) {
		if job.IsDone() {
			jobkeys = append(jobkeys, job.GetJobID())
		}
	})
	for _, key := range jobkeys {

		job, err := js.GetJob(key)
		if err != nil {
			fmt.Println("err getting job by key", err)
		}

		if job.IsDone() {
			err = js.DeleteJob(job)
			if err != nil {
				fmt.Println("err deleting job", err)
			}
		}
	}
}
