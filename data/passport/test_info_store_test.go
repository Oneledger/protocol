package passport

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	memDb     db.DB
	store     *TestInfoStore
	cs        *storage.State
	userId1   UserID
	userId2   UserID
	userId3   UserID
	addresses []keys.Address
	scanner1  UserID
	scanner2  UserID
	scanner3  UserID
	admAddrs  []keys.Address
	tid1      string
	tid2      string
	tid3      string
	info1     *TestInfo
	info2     *TestInfo
	info3     *TestInfo
	final1     *TestInfo
	final2     *TestInfo
	final3     *TestInfo
	updt1     *UpdateTestInfo
	updt2     *UpdateTestInfo
	updt3     *UpdateTestInfo
	log1      *ReadLog
	log2      *ReadLog
	log3      *ReadLog
	log4      *ReadLog
	log5      *ReadLog
)

func init() {
	userId1 = "person1"
	userId2 = "person2"
	userId3 = "person3"
	scanner1 = "scanner1"
	scanner2 = "scanner2"
	scanner3 = "scanner3"

	for i := 0; i < 3; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		addresses = append(addresses, h.Address())
	}

	for i := 0; i < 3; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		admAddrs = append(admAddrs, h.Address())
	}

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)
	nowStamp := now.Format(time.RFC3339)
	yesterdayStamp := yesterday.Format(time.RFC3339)
	tomorrowStamp := tomorrow.Format(time.RFC3339)

	sum1 := sha256.Sum256([]byte("test-id-1"))
	sum2 := sha256.Sum256([]byte("test-id-2"))
	sum3 := sha256.Sum256([]byte("test-id-3"))
	tid1 = hex.EncodeToString(sum1[:])
	tid2 = hex.EncodeToString(sum2[:])
	tid3 = hex.EncodeToString(sum3[:])

	info1 = NewTestInfo(tid1, userId1, TestCOVID19, TestSubAntiBody, "Audacia Bioscience", COVID19Pending,
		"YorkHospital", yesterdayStamp, "Edward", "", "", "", yesterdayStamp, "Tims", "showed similar symptoms")
	updt1 = NewUpdateTestInfo(tid1, userId1, TestCOVID19,   COVID19Negative, "YorkHospital", nowStamp, "Matthew", "analyze finished")
	final1 = NewTestInfo(tid1, userId1, TestCOVID19, TestSubAntiBody, "Audacia Bioscience", COVID19Negative,
		"YorkHospital", yesterdayStamp, "Edward", "YorkHospital", nowStamp, "Matthew",
		yesterdayStamp, "Tims", "showed similar symptoms\nanalyze finished (notes updated at: " + nowStamp + ")\n")


	info2 = NewTestInfo(tid2, userId1, TestCOVID19, TestSubAntigen, "MAG ELISA", COVID19Pending,
		"YorkHospital", nowStamp, "Edward", "", "", "", nowStamp, "Jackson", "want to test again")
	updt2 = NewUpdateTestInfo(tid2, userId1, TestCOVID19, COVID19Negative, "YorkHospital", tomorrowStamp, "Johnson", "analyze done")
	final2 = NewTestInfo(tid2, userId1, TestCOVID19, TestSubAntigen, "MAG ELISA", COVID19Negative,
		"YorkHospital", nowStamp, "Edward", "YorkHospital", tomorrowStamp, "Johnson",
		nowStamp, "Jackson", "want to test again\nanalyze done (notes updated at: " + tomorrowStamp + ")\n")

	info3 = NewTestInfo(tid3, userId2, TestCOVID19, TestSubPCR, "Gnomegen", COVID19Pending,
		"MichaelGarron Hospital", nowStamp, "George", "", "", "", nowStamp, "Obama", "had a fever")
	updt3 = NewUpdateTestInfo(tid3, userId2, TestCOVID19, COVID19Positive, "MichaelGarron Hospital", tomorrowStamp,
		"Kumar", "analyze done")
	final3 = NewTestInfo(tid3, userId2, TestCOVID19, TestSubPCR, "Gnomegen", COVID19Positive,
		"MichaelGarron Hospital", nowStamp, "George", "MichaelGarron Hospital", tomorrowStamp,
		"Kumar", nowStamp, "Obama", "had a fever\nanalyze done (notes updated at: " + tomorrowStamp + ")\n")

	log1 = NewReadLog("BorderService1", scanner1, admAddrs[0], userId1, addresses[0], TestCOVID19, yesterdayStamp)
	log2 = NewReadLog("BorderService1", scanner1, admAddrs[0], userId2, addresses[1], TestCOVID19, nowStamp)
	log3 = NewReadLog("BorderService1", scanner1, admAddrs[0], userId1, addresses[0], TestCOVID19, nowStamp)
	log4 = NewReadLog("BorderService2", scanner2, admAddrs[1], userId2, addresses[1], TestCOVID19, nowStamp)
	log5 = NewReadLog("BorderService2", scanner3, admAddrs[2], userId3, addresses[2], TestCOVID19, tomorrowStamp)
}

func setup() {
	fmt.Println("####### Testing info store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("cs", memDb))
	store = NewTestInfoStore("testinfo", "org", "read", cs)
}

func TestNewTestInfoStore(t *testing.T) {
	setup()
	assert.NotNil(t, store)
}

func TestInfoStore_AddTestInfo(t *testing.T) {
	setup()

	// add test results
	err := store.AddTestInfo(info1)
	assert.Nil(t, err)
	err = store.AddTestInfo(info2)
	assert.Nil(t, err)
}

func TestInfoStore_UpdateTestInfo(t *testing.T) {
	setup()

	// add and update test result
	err := store.AddTestInfo(info1)
	assert.Nil(t, err)
	err = store.UpdateTestInfo(updt1)
	assert.Nil(t, err)

	// check test info list
	infoList1, err := store.GetTestInfoByID(userId1, TestCOVID19)
	assert.Nil(t, err)
	assert.EqualValues(t, infoList1, []*TestInfo{final1})

	// update before create, MUST fail
	err = store.UpdateTestInfo(updt2)
	assert.Nil(t, err)

	// create one more test
	err = store.AddTestInfo(info2)
	assert.Nil(t, err)

	// check test info list again
	infoList1, err = store.GetTestInfoByID(userId1, TestCOVID19)
	assert.Nil(t, err)
	assert.EqualValues(t, infoList1, []*TestInfo{info2, final1})
}

func TestInfoStore_GetTestInfoByID(t *testing.T) {
	setup()

	// add test results
	err := store.AddTestInfo(info1)
	assert.Nil(t, err)
	err = store.AddTestInfo(info2)
	assert.Nil(t, err)
	err = store.UpdateTestInfo(updt2)
	err = store.AddTestInfo(info3)
	assert.Nil(t, err)
	err = store.UpdateTestInfo(updt3)

	// get test results person1
	infoList1, err := store.GetTestInfoByID(userId1, TestCOVID19)
	assert.Nil(t, err)
	assert.EqualValues(t, infoList1, []*TestInfo{final3, info1})

	// get test results person2
	infoList2, err := store.GetTestInfoByID(userId2, TestCOVID19)
	assert.Nil(t, err)
	assert.EqualValues(t, infoList2, []*TestInfo{final3})
}

func TestInfoStore_Iterate(t *testing.T) {
	setup()

	// add test results
	err := store.AddTestInfo(info1)
	assert.Nil(t, err)
	err = store.AddTestInfo(info2)
	assert.Nil(t, err)
	err = store.AddTestInfo(info3)
	assert.Nil(t, err)
	store.State.Commit()

	// iterate thru store
	userCount := 0
	store.Iterate(func(id UserID, infoList []*TestInfo) bool {
		userCount++
		if id == userId1 {
			assert.EqualValues(t, []*TestInfo{info2, info1}, infoList)
		}
		if id == userId2 {
			assert.EqualValues(t, []*TestInfo{info3}, infoList)
		}
		return false
	})
	assert.Equal(t, 2, userCount)
}

func TestInfoStore_IterateOrg(t *testing.T) {
	setup()

	// add test results
	err := store.AddTestInfo(info1)
	assert.Nil(t, err)
	err = store.AddTestInfo(info2)
	assert.Nil(t, err)
	err = store.AddTestInfo(info3)
	assert.Nil(t, err)
	store.State.Commit()

	// iterate thru store
	count := 0
	store.IterateOrgTests(TestCOVID19, TokenTypeID("YorkHospital"), "Tims", "",
		func(test TestType, org TokenTypeID, uploadedBy UserID, id UserID, num int) bool {
			assert.Equal(t, 1, num)
			count++
			return false
		})
	assert.Equal(t, 1, count)
	store.IterateOrgTests(TestCOVID19, TokenTypeID("YorkHospital"), "Jackson", "",
		func(test TestType, org TokenTypeID, uploadedBy UserID, id UserID, num int) bool {
			assert.Equal(t, 1, num)
			count++
			return false
		})
	assert.Equal(t, 2, count)
	store.IterateOrgTests(TestCOVID19, TokenTypeID("MichaelGarron Hospital"), "Obama", "",
		func(test TestType, org TokenTypeID, uploadedBy UserID, id UserID, num int) bool {
			assert.Equal(t, 1, num)
			count++
			return false
		})
	assert.Equal(t, 3, count)
}

func TestInfoStore_LogRead(t *testing.T) {
	setup()

	// scan some persons
	err := store.LogRead(log1)
	assert.Nil(t, err)
	err = store.LogRead(log2)
	assert.Nil(t, err)
	err = store.LogRead(log3)
	assert.Nil(t, err)
	err = store.LogRead(log4)
	assert.Nil(t, err)
	err = store.LogRead(log5)
	assert.Nil(t, err)
	store.State.Commit()

	// read logs
	logs, err := store.GetReadLogs(TestCOVID19, "BorderService1", scanner1, userId1)
	assert.Nil(t, err)
	assert.EqualValues(t, []*ReadLog{log3, log1}, logs)
	logs, err = store.GetReadLogs(TestCOVID19, "BorderService1", scanner1, userId2)
	assert.Nil(t, err)
	assert.EqualValues(t, []*ReadLog{log2}, logs)
	logs, err = store.GetReadLogs(TestCOVID19, "BorderService2", scanner2, userId2)
	assert.Nil(t, err)
	assert.EqualValues(t, []*ReadLog{log4}, logs)
	logs, err = store.GetReadLogs(TestCOVID19, "BorderService2", scanner3, userId3)
	assert.Nil(t, err)
	assert.EqualValues(t, []*ReadLog{log5}, logs)
}

func TestInfoStore_IterateReadLogs(t *testing.T) {
	setup()

	// scan some persons
	err := store.LogRead(log1)
	assert.Nil(t, err)
	err = store.LogRead(log2)
	assert.Nil(t, err)
	err = store.LogRead(log3)
	assert.Nil(t, err)
	err = store.LogRead(log4)
	assert.Nil(t, err)
	err = store.LogRead(log5)
	assert.Nil(t, err)
	store.State.Commit()

	// iterate thru store
	count := 0
	store.IterateReadLogs(TestCOVID19, "BorderService1", scanner1, func(id UserID, logs []*ReadLog) bool {
		if id == userId1 {
			assert.EqualValues(t, []*ReadLog{log3, log1}, logs)
		}
		if id == userId2 {
			assert.EqualValues(t, []*ReadLog{log2}, logs)
		}
		count += len(logs)
		return false
	})
	assert.Equal(t, 3, count)
	store.IterateReadLogs(TestCOVID19, "BorderService2", scanner2, func(id UserID, logs []*ReadLog) bool {
		assert.Equal(t, userId2, id)
		assert.EqualValues(t, []*ReadLog{log4}, logs)
		count += len(logs)
		return false
	})
	assert.Equal(t, 4, count)
	store.IterateReadLogs(TestCOVID19, "BorderService2", scanner3, func(id UserID, logs []*ReadLog) bool {
		assert.Equal(t, userId3, id)
		assert.EqualValues(t, []*ReadLog{log5}, logs)
		count += len(logs)
		return false
	})
	assert.Equal(t, 5, count)
}
