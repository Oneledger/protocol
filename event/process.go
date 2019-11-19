/*

 */

package event

import (
	"strings"

	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/serialize"
)

type JobProcess func(job jobs.Job) jobs.Job

func ProcessAllJobs(ctx *JobsContext, js *jobs.JobStore) {

	RangeJobs(js, func(job jobs.Job) jobs.Job {

		job.DoMyJob(ctx)
		return job
	})
}

func RangeJobs(js *jobs.JobStore, pro JobProcess) {

	start := []byte("job:     ")
	end := []byte("job:~~~~~~~~")
	isAsc := true

	jobkeys := make([]string, 0, 20)

	session := js.BeginSession()
	iter := session.GetIterator()

	iter.IterateRange(start, end, isAsc, func(key, val []byte) bool {

		jobkeys = append(jobkeys, string(key))

		return false
	})

	for _, key := range jobkeys {

		jobID := strings.TrimPrefix(string(key), "job:")

		dat, typ := js.GetJob(jobID)
		job := MakeJob(dat, typ)
		if job == nil {
			continue
		}

		job = pro(job)

		js.DeleteJob(job)

	}
}

func MakeJob(data []byte, typ string) jobs.Job {

	ser := serialize.GetSerializer(serialize.PERSISTENT)

	switch typ {
	case JobTypeAddSignature:
		as := JobAddSignature{}
		ser.Deserialize(data, &as)
		return &as
	case JobTypeBTCBroadcast:
		as := JobBTCBroadcast{}
		ser.Deserialize(data, &as)
		return &as
	}

	return nil
}
