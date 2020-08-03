package evidence

import (
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

const (
	VOTING   int8 = 0x1
	INNOCENT int8 = 0x2
	GUILTY   int8 = 0x3

	YES int8 = 0x1
	NO  int8 = 0x2
)

func VoteToString(status int8) string {
	switch status {
	case VOTING:
		return "Voting"
	case INNOCENT:
		return "Innocent"
	case GUILTY:
		return "Guilty"
	default:
		return "Not set"
	}
}

func ChoiceStrToInt8(choice string) int8 {
	switch choice {
	case "yes":
		return YES
	case "no":
		return NO
	default:
		return 0
	}
}

type AllegationTracker struct {
	Requests map[int64]bool
}

func (es *EvidenceStore) getAllegationTrackerKey() []byte {
	key := []byte(fmt.Sprintf("_atark"))
	return key
}

func (es *EvidenceStore) GetAllegationTracker() (*AllegationTracker, error) {
	dat, err := es.Get(es.getAllegationTrackerKey())
	if err != nil {
		return nil, err
	}
	at := &AllegationTracker{
		Requests: make(map[int64]bool),
	}
	if len(dat) == 0 {
		return at, nil
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, at)
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (es *EvidenceStore) SetAllegationTracker(at *AllegationTracker) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(at)
	if err != nil {
		return err
	}
	return es.Set(es.getAllegationTrackerKey(), dat)
}

func (es *EvidenceStore) GetRequestID() (*AllegationRequestID, error) {
	prefix := es.getAllegationRequestIDKey()
	dat, err := es.Get(prefix)
	if err != nil {
		return nil, err
	}
	ar := &AllegationRequestID{}
	if len(dat) != 0 {
		err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ar)
		if err != nil {
			return nil, err
		}
	}
	return ar, nil
}

func (es *EvidenceStore) SetRequestID(ar *AllegationRequestID) error {
	prefix := es.getAllegationRequestIDKey()
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ar)
	if err != nil {
		return err
	}
	err = es.Set(prefix, dat)
	if err != nil {
		return err
	}
	return nil
}

func (es *EvidenceStore) GenerateRequestID() (int64, error) {
	es.mux.Lock()
	defer es.mux.Unlock()
	ar, err := es.GetRequestID()
	if err != nil {
		return 0, err
	}

	ar.ID++

	err = es.SetRequestID(ar)
	if err != nil {
		return 0, err
	}
	return ar.ID, nil
}

func (es *EvidenceStore) getAllegationRequestKey(requestID int64) []byte {
	key := []byte(fmt.Sprintf("_ark_%d", requestID))
	return key
}

func (es *EvidenceStore) getAllegationRequestIDKey() []byte {
	key := []byte(fmt.Sprintf("_arid"))
	return key
}

func (es *EvidenceStore) SetAllegationRequest(ar *AllegationRequest) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ar)
	if err != nil {
		return err
	}
	return es.Set(es.getAllegationRequestKey(ar.ID), dat)
}

func (es *EvidenceStore) GetAllegationRequest(ID int64) (*AllegationRequest, error) {
	dat, err := es.Get(es.getAllegationRequestKey(ID))
	if err != nil {
		return nil, err
	}
	if len(dat) == 0 {
		return nil, fmt.Errorf("Request %d not found", ID)
	}
	ar := &AllegationRequest{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

type AllegationRequestID struct {
	ID int64
}

type AllegationVote struct {
	Address keys.Address
	Choice  int8
}

type AllegationRequest struct {
	ID               int64
	ReporterAddress  keys.Address
	MaliciousAddress keys.Address
	BlockHeight      int64
	ProofMsg         string
	Status           int8
	Votes            []*AllegationVote
}

func (ar *AllegationRequest) Bytes() ([]byte, error) {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ar)
	if err != nil {
		return []byte{}, fmt.Errorf("validator not serializable %s", err)
	}
	return value, nil
}

func (ar *AllegationRequest) FromBytes(msg []byte) (*AllegationRequest, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, ar)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize account from bytes %s", err)
	}
	return ar, nil
}

func NewAllegationRequest(ID int64, reporterAddress keys.Address, maliciousAddress keys.Address, blockHeight int64, proofMsg string) *AllegationRequest {
	return &AllegationRequest{
		ID:               ID,
		ReporterAddress:  reporterAddress,
		MaliciousAddress: maliciousAddress,
		BlockHeight:      blockHeight,
		ProofMsg:         proofMsg,
		Status:           VOTING,
		Votes:            make([]*AllegationVote, 0),
	}
}
