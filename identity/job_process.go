/*

 */

package identity

import (
	"strings"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type JobProcess func(job Job) Job

func ProcessAllJobs(ctx *JobsContext, js *JobStore) {

	js.RangeJobs(func(job Job) Job {

		if !job.IsMyJobDone(ctx) {
			job.DoMyJob(ctx)
		}

		if job.IsSufficient() {
			job.DoFinalize()
		}

		return job
	})
}

type JobStore struct {
	storage.ChainState
	ser serialize.Serializer
}

func NewJobStore(store *storage.ChainState) *JobStore {
	return &JobStore{
		*store,
		serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (js *JobStore) SaveJob(job Job) error {

	key := storage.StoreKey("job:" + job.GetJobID())
	typKey := storage.StoreKey("jobtype:" + job.GetJobID())

	dat, err := js.ser.Serialize(job)
	if err != nil {
		return err
	}

	js.Set(key, dat)
	return js.Set(typKey, []byte(job.GetType()))
}

func (js *JobStore) GetJob(jobID string) Job {
	key := storage.StoreKey("job:" + jobID)
	typKey := storage.StoreKey("jobtype:" + jobID)

	dat, err := js.Get(key)
	if err != nil {
		return nil
	}
	typ, err := js.Get(typKey)
	if err != nil {
		return nil
	}

	return makeJob(dat, string(typ))
}

func (js *JobStore) DeleteJob(job Job) error {

	key := storage.StoreKey("job:" + job.GetJobID())
	typKey := storage.StoreKey("jobtype:" + job.GetJobID())

	_, err := js.Delete(key)
	_, err = js.Delete(typKey)

	return err
}

func (js *JobStore) RangeJobs(pro JobProcess) {
	start := []byte("job:        ")
	end := []byte("job:~~~~~~~~")
	isAsc := true

	jobkeys := make([]string, 0, 20)

	js.IterateRange(start, end, isAsc, func(key, val []byte) bool {

		jobkeys = append(jobkeys, string(key))

		return false
	})

	for _, key := range jobkeys {
		jobID := strings.TrimPrefix(string(key), "job:")

		job := js.GetJob(jobID)

		job = pro(job)

		if job.IsDone() {
			js.DeleteJob(job)
		} else {
			js.SaveJob(job)
		}
	}
}
