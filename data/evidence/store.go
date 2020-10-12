package evidence

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type EvidenceStore struct {
	state  *storage.State
	prefix []byte
	mux    sync.Mutex
}

func NewEvidenceStore(prefix string, state *storage.State) *EvidenceStore {
	return &EvidenceStore{
		state:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (es *EvidenceStore) WithState(state *storage.State) *EvidenceStore {
	es.state = state
	return es
}

func (es *EvidenceStore) Get(key []byte) ([]byte, error) {
	prefixKey := append(es.prefix, key...)

	dat, err := es.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (es *EvidenceStore) GetVersioned(key []byte, height int64, diff int64) []byte {
	prefixKey := append(es.prefix, key...)

	return es.state.GetVersioned(height-diff, storage.StoreKey(prefixKey))
}

func (es *EvidenceStore) Set(key []byte, value []byte) error {
	prefixKey := append(es.prefix, key...)
	err := es.state.Set(storage.StoreKey(prefixKey), value)
	return err
}

func (es *EvidenceStore) delete(key storage.StoreKey) (bool, error) {
	prefixed := append(es.prefix, key...)
	ok, err := es.state.Delete(prefixed)
	if !ok || err != nil {
		return ok, err
	}
	return ok, nil
}

func (es *EvidenceStore) getSuspiciousValidatorKey(validatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_ssvk_%s", validatorAddress))
	return key
}

func (es *EvidenceStore) getSuspiciousVL() []byte {
	key := []byte(fmt.Sprintf("_svvl"))
	return key
}

func (es *EvidenceStore) GetSuspiciousValidator(validatorAddress keys.Address, height int64, diff int64) (*LastValidatorHistory, error) {
	var dat []byte
	if height == 0 {
		dat, _ = es.Get(es.getSuspiciousValidatorKey(validatorAddress))
	} else {
		dat = es.GetVersioned(es.getSuspiciousValidatorKey(validatorAddress), height, diff)
	}
	if len(dat) == 0 {
		return nil, fmt.Errorf("Suspisious validator not found")
	}
	lvh := &LastValidatorHistory{}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, lvh)
	if err != nil {
		return nil, err
	}
	return lvh, nil
}

func (es *EvidenceStore) IterateSuspiciousValidators(fn func(lvh *LastValidatorHistory) bool) (stopped bool) {
	key := []byte(fmt.Sprintf("_ssvk_"))
	prefixKey := append(es.prefix, key...)
	return es.state.IterateRange(
		prefixKey,
		storage.Rangefix(string(prefixKey)),
		true,
		func(key, value []byte) bool {
			lvh, err := (&LastValidatorHistory{}).FromBytes(value)
			if err != nil {
				fmt.Println("failed to deserialize susp validator")
				return false
			}
			if !lvh.IsFrozen() {
				return false
			}
			return fn(lvh)
		},
	)
}

func (es *EvidenceStore) GetFrozenMap() map[string]bool {
	fMap := make(map[string]bool)
	es.IterateSuspiciousValidators(func(lvh *LastValidatorHistory) bool {
		fMap[lvh.Address.String()] = true
		return false
	})
	return fMap
}

func (es *EvidenceStore) IterateRequests(fn func(ar *AllegationRequest) bool) (stopped bool) {
	key := []byte("_ark_")
	prefixKey := append(es.prefix, key...)
	return es.state.IterateRange(
		prefixKey,
		storage.Rangefix(string(prefixKey)),
		true,
		func(key, value []byte) bool {
			ar, err := (&AllegationRequest{}).FromBytes(value)
			if err != nil {
				fmt.Println("failed to deserialize allegation request")
				return false
			}
			return fn(ar)
		},
	)
}

func (es *EvidenceStore) CheckRequestExists(new *AllegationRequest) bool {
	requestAlreadyExists := false
	es.IterateRequests(func(ar *AllegationRequest) bool {
		if ar.MaliciousAddress.Equal(new.MaliciousAddress) {
			requestAlreadyExists = true
		}
		return false
	})
	return requestAlreadyExists
}

func (es *EvidenceStore) CreateSuspiciousValidator(validatorAddress keys.Address, status int8, height int64, createdAt *time.Time) (*LastValidatorHistory, error) {
	lvh := NewLastValidatorHistory(validatorAddress, status, height, createdAt)
	err := es.UpdateSuspiciousValidator(lvh)
	return lvh, err
}

func (es *EvidenceStore) UpdateSuspiciousValidator(lvh *LastValidatorHistory) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(lvh)
	if err != nil {
		return err
	}
	return es.Set(es.getSuspiciousValidatorKey(lvh.Address), dat)
}

func (es *EvidenceStore) IsFrozenValidator(validatorAddress keys.Address) bool {
	dat, err := es.Get(es.getSuspiciousValidatorKey(validatorAddress))
	if err != nil {
		return false
	}
	if len(dat) == 0 {
		return false
	}
	lvh := &LastValidatorHistory{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, lvh)
	if err != nil {
		return false
	}

	if err != nil {
		return false
	}
	return lvh.IsFrozen()
}

func (es *EvidenceStore) HandleRelease(options *Options, validatorAddress keys.Address, blockHeight int64, blockCreatedAt time.Time) error {
	lvh, err := es.GetSuspiciousValidator(validatorAddress, 0, 0)
	if err != nil {
		return err
	}
	if !lvh.IsFrozen() {
		return fmt.Errorf("Validator \"%s\" already released", validatorAddress)
	}

	isReady, err := lvh.ReleaseReady(options, blockCreatedAt)
	if err != nil {
		return err
	}

	if !isReady {
		return fmt.Errorf("Validator \"%s\" not ready for release", validatorAddress)
	}

	lvh.ReleaseHeight = blockHeight
	lvh.ReleaseAt = &blockCreatedAt
	err = es.UpdateSuspiciousValidator(lvh)
	if err != nil {
		return err
	}
	return nil
}

func (es *EvidenceStore) PerformAllegation(validatorAddress keys.Address, maliciousAddress keys.Address, ID string, blockHeight int64, proofMsg string) error {
	es.mux.Lock()
	defer es.mux.Unlock()

	isBusy := es.IsRequestIDBusy(ID)
	if isBusy {
		return fmt.Errorf("request ID %s already handled\n", ID)
	}

	ar := NewAllegationRequest(ID, validatorAddress, maliciousAddress, blockHeight, proofMsg)
	if es.CheckRequestExists(ar) {
		return errors.New(fmt.Sprintf("allegation for this validator already exists : %s", ar.String()))
	}
	err := es.SetAllegationRequest(ar)
	if err != nil {
		return err
	}

	at, err := es.GetAllegationTracker()
	if err != nil {
		return err
	}

	at.Requests[ID] = true

	err = es.SetAllegationTracker(at)
	if err != nil {
		return err
	}
	return nil
}

func (es *EvidenceStore) Vote(requestID string, voteAddress keys.Address, choice int8) error {
	ar, err := es.GetAllegationRequest(requestID)
	if err != nil {
		return err
	}
	if choice != YES && choice != NO {
		return fmt.Errorf("Invalid choice, only YES or NO available")
	}
	if ar.Status == GUILTY || ar.Status == INNOCENT {
		return fmt.Errorf("Could not vote on closed requet")
	}
	if ar.Votes == nil {
		ar.Votes = make([]*AllegationVote, 0)
	}
	for i := range ar.Votes {
		vote := ar.Votes[i]
		if vote.Address.Equal(voteAddress) {
			return fmt.Errorf("You have been already voted on this request")
		}
	}
	vote := &AllegationVote{
		Address: voteAddress,
		Choice:  choice,
	}
	ar.Votes = append(ar.Votes, vote)

	err = es.SetAllegationRequest(ar)
	if err != nil {
		return err
	}
	return nil
}
