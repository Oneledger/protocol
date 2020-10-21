package evidence

import (
	"fmt"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
	"sort"
	"strings"
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
	Requests map[string]bool
}

func (es *EvidenceStore) CleanTracker() {
	at, err := es.GetAllegationTracker()
	if err != nil {
		return
	}
	requestIdList := make([]string, len(at.Requests))
	for k := range at.Requests {
		requestIdList = append(requestIdList, k)
	}
	sort.Strings(requestIdList)
	countMap := make(map[string]bool)
	for _, r := range requestIdList {
		ar, err := es.GetAllegationRequest(r)
		if err != nil {
			delete(at.Requests, r)
			continue
		}
		if countMap[ar.MaliciousAddress.String()] {
			fmt.Println("Deleting Duplicate Requests")
			es.DeleteAllegationRequest(r)
			delete(at.Requests, r)
		}
		countMap[ar.MaliciousAddress.String()] = true
	}
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
		Requests: make(map[string]bool),
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

func (es *EvidenceStore) IsRequestIDBusy(ID string) bool {
	_, err := es.GetAllegationRequest(ID)
	if err != nil {
		return false
	}
	return true

}

func (es *EvidenceStore) GenerateRequestID() (string, error) {
	uuidNew, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuidNew.String(), nil
}

func (es *EvidenceStore) getAllegationRequestKey(requestID string) []byte {
	key := []byte(fmt.Sprintf("_ark_%s", requestID))
	return key
}

func (es *EvidenceStore) SetAllegationRequest(ar *AllegationRequest) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ar)
	if err != nil {
		return err
	}
	return es.Set(es.getAllegationRequestKey(ar.ID), dat)
}

func (es *EvidenceStore) GetAllegationRequest(ID string) (*AllegationRequest, error) {
	dat, err := es.Get(es.getAllegationRequestKey(ID))
	if err != nil {
		return nil, err
	}
	if len(dat) == 0 {
		return nil, fmt.Errorf("Request %s not found", ID)
	}
	ar := &AllegationRequest{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (es *EvidenceStore) DeleteAllegationRequest(ID string) (bool, error) {
	ok, err := es.delete(es.getAllegationRequestKey(ID))
	if !ok || err != nil {
		return ok, err
	}
	return ok, nil
}

type AllegationVote struct {
	Address keys.Address
	Choice  int8
}

type AllegationRequest struct {
	ID               string
	ReporterAddress  keys.Address
	MaliciousAddress keys.Address
	BlockHeight      int64
	ProofMsg         string
	Status           int8
	Votes            []*AllegationVote
}

func (ar *AllegationRequest) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s
ReporterAddress: %s
MaliciousAddress: %s
Votes %v
`, ar.ID, ar.ReporterAddress, ar.MaliciousAddress, ar.Votes))
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

func NewAllegationRequest(ID string, reporterAddress keys.Address, maliciousAddress keys.Address, blockHeight int64, proofMsg string) *AllegationRequest {
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
