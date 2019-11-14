/*

 */

package jobs

import (
	"errors"
	"sync"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const rootkey = "rootkey"

type JobStore struct {
	storage.SessionedStorage
	ser  serialize.Serializer
	lock sync.Mutex
}

func NewJobStore(config config.Server, dbDir string) *JobStore {

	store := storage.NewStorageDB(storage.KEYVALUE, "validatorJobs", dbDir, config.Node.DB)

	return &JobStore{
		store,
		serialize.GetSerializer(serialize.PERSISTENT),
		sync.Mutex{},
	}
}

func (js *JobStore) SaveJob(job Job) error {

	key := storage.StoreKey("job:" + job.GetJobID())
	typKey := storage.StoreKey("jobtype:" + job.GetJobID())

	dat, err := js.ser.Serialize(job)
	if err != nil {
		return err
	}

	js.lock.Lock()
	session := js.BeginSession()
	err = session.Set(key, dat)
	if err != nil {
		return err
	}
	err = session.Set(typKey, []byte(job.GetType()))
	if err != nil {
		return err
	}
	ok := session.Commit()
	js.lock.Unlock()
	if !ok {
		return errors.New("err commiting to job store")
	}
	return nil
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

	js.lock.Lock()
	session := js.BeginSession()

	_, err := session.Delete(key)
	if err != nil {
		return err
	}

	_, err = session.Delete(typKey)
	if err != nil {
		return err
	}
	ok := session.Commit()
	js.lock.Unlock()
	if !ok {
		return errors.New("error committing to job store")
	}
	return nil
}
