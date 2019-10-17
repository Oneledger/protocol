/*

 */

package identity

import (
	"strings"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity/internal"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type JobProcess func(job Job)

func ProcessAllJobs(key keys.PrivateKey, jobs []Job) []Job {

	internal.NewService()

	for i := range jobs {

		if !jobs[i].IsMyJobDone(key) {

			jobs[i].DoMyJob()
		}

		if jobs[i].IsSufficient() {

			jobs[i].DoFinalize()
		}
	}

	return jobs
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

func (js *JobStore) RangeJobs(pro JobProcess) {
	start := []byte("job:        ")
	end := []byte("job:~~~~~~~~")
	isAsc := true

	js.IterateRange(start, end, isAsc, func(key, val []byte) bool {

		jobID := strings.TrimPrefix(string(key), "job:")

		typKey := storage.StoreKey("jobtype:" + jobID)
		typ, err := js.Get(typKey)
		if err != nil {
			//
		}

		job := makeJob(val, string(typ))
		pro(job)

		return false
	})
}
