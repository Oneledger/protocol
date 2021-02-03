package passport

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type TestInfoStore struct {
	State      *storage.State
	szlr       serialize.Serializer
	prefix     []byte
	prefixOrg  []byte
	prefixRead []byte
}

func NewTestInfoStore(prefix, prefixOrg, prefixRead string, state *storage.State) *TestInfoStore {
	return &TestInfoStore{
		State:      state,
		prefix:     storage.Prefix(prefix),
		prefixOrg:  storage.Prefix(prefixOrg),
		prefixRead: storage.Prefix(prefixRead),
		szlr:       serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (ts *TestInfoStore) WithState(state *storage.State) *TestInfoStore {
	ts.State = state
	return ts
}

func (ts *TestInfoStore) AddTestInfo(info *TestInfo) (err error) {
	key := ts.getTestKey(info.PersonID, info.Test)
	infoList, err := ts.get(key)
	if err != nil {
		return
	}

	// sort by upload time
	infoList = append([]*TestInfo{info}, infoList...)
	sort.Slice(infoList, func(i, j int) bool {
		return infoList[i].TestID < infoList[j].TestID
	})

	err = ts.set(key, infoList)
	if err != nil {
		return
	}

	err = ts.addToOrg(info)
	return
}

func (ts *TestInfoStore) UpdateTestInfo(info *UpdateTestInfo) (updated bool, err error) {
	key := ts.getTestKey(info.PersonID, info.Test)
	infoList, err := ts.get(key)
	if err != nil {
		return
	}
	// update fields if possible
	for _, cur := range infoList {
		if cur.TestID == info.TestID {
			cur.Result = info.Result
			cur.AnalysisOrg = info.AnalysisOrg
			cur.AnalyzedAt = info.AnalyzedAt
			cur.AnalyzedBy = info.AnalyzedBy
			if info.Notes != "" {
				cur.Notes = cur.Notes + "\n" + info.Notes + " (notes updated at: " + info.AnalyzedAt + ")\n"
			}
			updated = true
			break
		}
	}
	// write to database
	err = ts.set(key, infoList)
	return
}

func (ts *TestInfoStore) GetTestInfoByID(id UserID, test TestType) (infoList []*TestInfo, err error) {
	key := ts.getTestKey(id, test)
	infoList, err = ts.get(key)
	return
}

func (ts *TestInfoStore) Iterate(fn func(id UserID, infoList []*TestInfo) bool) (stopped bool) {
	prefix := ts.prefix
	return ts.State.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			// keys in format "userId_testType"
			keys := strings.Split(string(key[len(prefix):]), storage.DB_PREFIX)
			if len(keys) != 2 {
				fmt.Printf("failed to deserialize test info keys")
				return true
			}
			userId := UserID(keys[0])

			infoList := []*TestInfo{}
			err := ts.szlr.Deserialize(value, &infoList)
			if err != nil {
				logger.Error("failed to deserialize test info list")
				return true
			}
			return fn(userId, infoList)
		},
	)
}

func (ts *TestInfoStore) IterateOrgTests(test TestType, org TokenTypeID, uploadedBy UserID, id UserID,
	fn func(test TestType, org TokenTypeID, uploadedBy, person UserID, num int) bool) (stopped bool) {
	// append filters
	prefix := ts.prefixOrg
	prefix = append(prefix, (test.String() + storage.DB_PREFIX)...)
	if len(org.String()) > 0 {
		prefix = append(prefix, (org.String() + storage.DB_PREFIX)...)
		if len(uploadedBy.String()) > 0 {
			prefix = append(prefix, (uploadedBy.String() + storage.DB_PREFIX)...)
			if len(id.String()) > 0 {
				prefix = append(prefix, id.String()...)
			}
		}
	}
	// iterate
	return ts.State.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			// keys in format "testType_orgId_uploadBy_person"
			keys := strings.Split(string(key[len(ts.prefixOrg):]), storage.DB_PREFIX)
			if len(keys) != 4 {
				fmt.Printf("failed to deserialize test info org keys")
				return true
			}
			test := GetTestType(keys[0])
			org = TokenTypeID(keys[1])
			uploadedBy = UserID(keys[2])
			person := UserID(keys[3])

			num := 0
			err := ts.szlr.Deserialize(value, &num)
			if err != nil {
				logger.Error("failed to deserialize num of tests")
				return true
			}
			return fn(test, org, uploadedBy, person, num)
		},
	)
}

func (ts *TestInfoStore) LogRead(log *ReadLog) (err error) {
	key := ts.getReadKey(log.Test, log.Org, log.ReadBy, log.Person)
	logs, err := ts.getLogs(key)
	if err != nil {
		return
	}

	// sort by scan time
	logs = append([]*ReadLog{log}, logs...)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].ReadAt > logs[j].ReadAt
	})

	err = ts.set(key, logs)
	return
}

func (ts *TestInfoStore) GetReadLogs(test TestType, org TokenTypeID, readBy, id UserID) (logs []*ReadLog, err error) {
	key := ts.getReadKey(test, org, readBy, id)
	logs, err = ts.getLogs(key)
	return
}

func (ts *TestInfoStore) IterateReadLogs(test TestType, org TokenTypeID, readBy UserID, fn func(id UserID, logs []*ReadLog) bool) (stopped bool) {
	// append filters
	prefix := ts.prefixRead
	prefix = append(prefix, (test.String() + storage.DB_PREFIX)...)
	if len(org.String()) > 0 {
		prefix = append(prefix, (org.String() + storage.DB_PREFIX)...)
		if len(readBy.String()) > 0 {
			prefix = append(prefix, readBy.String()...)
		}
	}
	// iterate
	return ts.State.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			// keys in format "userId_testType"
			keys := strings.Split(string(key[len(ts.prefixRead):]), storage.DB_PREFIX)
			if len(keys) != 4 {
				fmt.Printf("failed to deserialize read log keys")
				return true
			}
			userId := UserID(keys[3])

			logs := []*ReadLog{}
			err := ts.szlr.Deserialize(value, &logs)
			if err != nil {
				logger.Error("failed to deserialize read logs")
				return true
			}
			return fn(userId, logs)
		},
	)
}

//-----------------------------helper functions
// get COVID key
func (ts *TestInfoStore) getTestKey(id UserID, test TestType) []byte {
	key := fmt.Sprintf("%s%s_%s", string(ts.prefix), id, test)
	return storage.StoreKey(key)
}

func (ts *TestInfoStore) getTestOrgKey(test TestType, org TokenTypeID, uploadedBy UserID, id UserID) []byte {
	key := fmt.Sprintf("%s%s_%s_%s_%s", string(ts.prefixOrg), test, org, uploadedBy, id)
	return storage.StoreKey(key)
}

func (ts *TestInfoStore) getReadKey(test TestType, org TokenTypeID, readBy, id UserID) []byte {
	key := fmt.Sprintf("%s%s_%s_%s_%s", string(ts.prefixRead), test, org, readBy, id)
	return storage.StoreKey(key)
}

// Set test info by key
func (ts *TestInfoStore) set(key storage.StoreKey, obj interface{}) error {
	dat, err := ts.szlr.Serialize(obj)
	if err != nil {
		return err
	}
	err = ts.State.Set(key, dat)
	return err
}

// Get test info list by key
func (ts *TestInfoStore) get(key storage.StoreKey) (infoList []*TestInfo, err error) {
	dat, err := ts.State.Get(key)
	if err != nil {
		return
	}
	infoList = []*TestInfo{}
	if len(dat) == 0 {
		return
	}
	err = ts.szlr.Deserialize(dat, &infoList)
	return
}

// Get number of tests
func (ts *TestInfoStore) getNumofTests(key storage.StoreKey) (num int, err error) {
	dat, err := ts.State.Get(key)
	if err != nil {
		return
	}
	num = 0
	if len(dat) == 0 {
		return
	}
	err = ts.szlr.Deserialize(dat, &num)
	return
}

// Get logs by key
func (ts *TestInfoStore) getLogs(key storage.StoreKey) (logs []*ReadLog, err error) {
	dat, err := ts.State.Get(key)
	if err != nil {
		return
	}
	logs = []*ReadLog{}
	if len(dat) == 0 {
		return
	}
	err = ts.szlr.Deserialize(dat, &logs)
	return
}

// Associate test info with organization
func (ts *TestInfoStore) addToOrg(info *TestInfo) (err error) {
	key := ts.getTestOrgKey(info.Test, info.TestOrg, info.UploadedBy, info.PersonID)
	num, err := ts.getNumofTests(key)
	if err != nil {
		return
	}

	num++
	err = ts.set(key, num)
	return
}
