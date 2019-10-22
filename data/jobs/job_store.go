/*

 */

package jobs


import (
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

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

func (js *JobStore) GetJob(jobID string) ([]byte, string) {
	key := storage.StoreKey("job:" + jobID)
	typKey := storage.StoreKey("jobtype:" + jobID)

	dat, err := js.Get(key)
	if err != nil {
		return nil, ""
	}
	typ, err := js.Get(typKey)
	if err != nil {
		return nil, ""
	}

	return dat, string(typ)
}

func (js *JobStore) DeleteJob(job Job) error {

	key := storage.StoreKey("job:" + job.GetJobID())
	typKey := storage.StoreKey("jobtype:" + job.GetJobID())

	_, err := js.Delete(key)
	_, err = js.Delete(typKey)

	return err
}

