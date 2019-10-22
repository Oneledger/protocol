/*

 */

package app

import (
	"strings"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/serialize"
)

type JobProcess func(job jobs.Job) jobs.Job

func ProcessAllJobs(ctx *action.JobsContext, js *jobs.JobStore) {

	RangeJobs(js, func(job jobs.Job) jobs.Job {

		if !job.IsMyJobDone(ctx) {
			job.DoMyJob(ctx)
		}

		if job.IsSufficient(ctx) {
			job.DoFinalize()
		}

		return job
	})
}

func RangeJobs(js *jobs.JobStore, pro JobProcess) {
	start := []byte("job:     ")
	end := []byte("job:~~~~~~~~")
	isAsc := true

	jobkeys := make([]string, 0, 20)

	js.IterateRange(start, end, isAsc, func(key, val []byte) bool {

		jobkeys = append(jobkeys, string(key))

		return false
	})

	for _, key := range jobkeys {
		jobID := strings.TrimPrefix(string(key), "job:")

		dat, typ := js.GetJob(jobID)
		job := makeJob(dat, typ)

		job = pro(job)

		if job.IsDone() {
			js.DeleteJob(job)
		} else {
			js.SaveJob(job)
		}
	}
}

func makeJob(data []byte, typ string) jobs.Job {

	ser := serialize.GetSerializer(serialize.PERSISTENT)

	switch typ {
	case btc.JobTypeAddSignature:
		as := btc.JobAddSignature{}
		ser.Deserialize(data, &as)
		return &as
	case btc.JobTypeBTCBroadcast:
		as := btc.JobBTCBroadcast{}
		ser.Deserialize(data, &as)
		return &as
	}
}
