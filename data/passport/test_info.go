package passport

import (
	"fmt"
	"time"

	"github.com/Oneledger/protocol/serialize"
)

type TestInfo struct {
	TestID       string      `json:"testId"`
	PersonID     UserID      `json:"personId"`
	Test         TestType    `json:"test"`
	SubTest      TestSubType `json:"subTest"`
	Manufacturer string      `json:"manufacturer"`
	Result       TestResult  `json:"testResult"`

	TestOrg     TokenTypeID `json:"testOrg"`
	TestedAt    string      `json:"testedAt"`
	TestedBy    string      `json:"testedBy"`
	AnalysisOrg TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string      `json:"analyzedAt"`
	AnalyzedBy  string      `json:"analyzedBy"`

	// admin info
	UploadedAt string `json:"uploadedAt"`
	UploadedBy UserID `json:"uploadedBy"`
	Notes      string `json:"notes"`
}

type UpdateTestInfo struct {
	TestID      string      `json:"testId"`
	PersonID    UserID      `json:"personId"`
	Test        TestType    `json:"test"`
	Result      TestResult  `json:"testResult"`
	AnalysisOrg TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string      `json:"analyzedAt"`
	AnalyzedBy  string      `json:"analyzedBy"`
	Notes       string      `json:"notes"`
}

func NewTestInfo(testId string, id UserID, test TestType, subTest TestSubType, manufacturer string, result TestResult,
	testOrg TokenTypeID, testedAt string, testedBy string, analysisOrg TokenTypeID, analyzedAt string, analyzedBy string,
	uploadedAt string, uploadedBy UserID, notes string) *TestInfo {
	info := &TestInfo{
		TestID:       testId,
		PersonID:     id,
		Test:         test,
		SubTest:      subTest,
		Manufacturer: manufacturer,
		Result:       result,

		TestOrg:     testOrg,
		TestedAt:    testedAt,
		TestedBy:    testedBy,
		AnalysisOrg: analysisOrg,
		AnalyzedAt:  analyzedAt,
		AnalyzedBy:  analyzedBy,

		UploadedAt: uploadedAt,
		UploadedBy: uploadedBy,
		Notes:      notes,
	}
	if info.AnalyzedAt == "" {
		info.AnalyzedAt = time.Time{}.Format(time.RFC3339)
	}
	if info.TestedAt == "" {
		info.TestedAt = time.Time{}.Format(time.RFC3339)
	}
	if info.UploadedAt == "" {
		info.UploadedAt = time.Time{}.Format(time.RFC3339)
	}
	return info
}

func NewUpdateTestInfo(testID string, personID UserID, test TestType, result TestResult, analysisOrg TokenTypeID, analyzedAt string, analyzedBy string, notes string) *UpdateTestInfo {
	updateInfo := &UpdateTestInfo{
		TestID: testID,
		PersonID: personID,
		Test: test,
		Result: result,
		AnalysisOrg: analysisOrg,
		AnalyzedAt: analyzedAt,
		AnalyzedBy: analyzedBy,
		Notes: notes,
	}
	if updateInfo.AnalyzedAt == "" {
		updateInfo.AnalyzedAt = time.Time{}.Format(time.RFC3339)
	}
	return updateInfo
}

func (info *TestInfo) TimeUpload() time.Time {
	tm, _ := time.Parse(time.RFC3339, info.UploadedAt)
	return tm
}

func (info *TestInfo) TimeAnalyze() time.Time {
	tm, _ := time.Parse(time.RFC3339, info.AnalyzedAt)
	return tm
}

func (info *TestInfo) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(info)
	if err != nil {
		logger.Error("test info not serializable", err)
		return []byte{}
	}
	return value
}

func (info *TestInfo) FromBytes(msg []byte) (*TestInfo, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, info)
	if err != nil {
		logger.Error("failed to deserialize test info from bytes", err)
		return nil, err
	}
	return info, nil
}

func (info *TestInfo) String() string {
	return fmt.Sprintf("PersonID= %s, TestType= %s, Result= %s, UploadedBy= %s, UploadedAt= %s",
		info.PersonID, info.Test, info.Result, info.UploadedBy, info.UploadedAt)
}
