/*

 */

package jobs

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const rootkey = "rootkey"

type JobStore struct {
	storage.SessionedStorage
	chain chain.Type
	ser   serialize.Serializer
	lock  sync.Mutex
}

func NewJobStore(config config.Server, dbDir string) *JobStore {

	store := storage.NewStorageDB(storage.KEYVALUE, "validatorJobs", dbDir, config.Node.DB)

	return &JobStore{
		SessionedStorage: store,
		chain:            chain.Type(-1),
		ser:              serialize.GetSerializer(serialize.LOCAL),
		lock:             sync.Mutex{},
	}
}
func (js *JobStore) WithChain(chain chain.Type) *JobStore {
	js.chain = chain
	return js
}

func (js *JobStore) SaveJob(job Job) error {

	key := storage.StoreKey("job:" + js.chain.String() + ":" + job.GetJobID())
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

	ok := session.Commit()
	js.lock.Unlock()
	if !ok {
		return errors.New("err commiting to job store")
	}
	return nil
}

func (js *JobStore) GetJob(jobID string) (Job, error) {
	key := storage.StoreKey("job:" + js.chain.String() + ":" + jobID)
	dat, err := js.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, key.String())
	}
	var job Job
	err = js.ser.Deserialize(dat, &job)
	if err != nil {
		return nil, err
	}
	return job, err
}

func (js *JobStore) DeleteJob(job Job) error {

	key := storage.StoreKey("job:" + js.chain.String() + ":" + job.GetJobID())

	js.lock.Lock()
	session := js.BeginSession()

	_, err := session.Delete(key)
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

func (js *JobStore) Iterate(fn func(job Job)) {
	start := []byte("job:" + js.chain.String() + ":")
	end := []byte("job:" + js.chain.String() + storage.DB_RANGEFIX)
	isAsc := true

	session := js.BeginSession()
	iter := session.GetIterator()

	iter.IterateRange(start, end, isAsc, func(key, val []byte) bool {
		var job Job
		err := js.ser.Deserialize(val, &job)
		if err != nil {
			return false
		}
		fn(job)
		return false
	})

}
