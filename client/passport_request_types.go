package client

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/passport"
)

//****** Query ******

type PSPTFilterTestInfoRequest struct {
	// For auth
	Org   passport.TokenTypeID `json:"org"`
	Admin passport.UserID      `json:"admin"`

	// For query
	Test             passport.TestType    `json:"test"`
	UploadedOrg      passport.TokenTypeID `json:"uploadedOrg"`
	UploadedBy       passport.UserID      `json:"uploadedBy"`
	AnalysisOrg      passport.TokenTypeID `json:"analysisOrg"`
	AnalyzedBy       string               `json:"analyzedBy"`
	Person           passport.UserID      `json:"person"`
	TestTimeBegin    string               `json:"testTimeBegin"`
	TestTimeEnd      string               `json:"testTimeEnd"`
	AnalyzeTimeBegin string               `json:"analyzeTimeBegin"`
	AnalyzeTimeEnd   string               `json:"analyzeTimeEnd"`
}

// For query from user, does not require permission
type PSPTQueryTestInfoForUserReq struct {
	Test             passport.TestType    `json:"test"`
	Person           passport.UserID      `json:"person"`
}

type PSPTQueryTestInfoForUserReply struct {
	InfoList []*passport.TestInfo `json:"infoList"`
	Height   int64                `json:"height"`
}

type AuthTokenRequest struct {
	Org passport.TokenTypeID `json:"org"`
	ID  passport.UserID      `json:"id"`
}

type PSPTFilterTestInfoReply struct {
	InfoList []*passport.TestInfo `json:"infoList"`
	Height   int64                `json:"height"`
}

type PSPTFilterReadLogsRequest struct {
	// For auth
	Org   passport.TokenTypeID `json:"org"`
	Admin passport.UserID      `json:"admin"`

	// For query
	Test      passport.TestType    `json:"test"`
	ReadOrg   passport.TokenTypeID `json:"readOrg"`
	ReadBy    passport.UserID      `json:"readBy"`
	Person    passport.UserID      `json:"person"`
	TimeBegin string               `json:"timeBegin"`
	TimeEnd   string               `json:"timeEnd"`
}

type PSPTFilterReadLogsReply struct {
	Logs   []*passport.ReadLog `json:"logs"`
	Height int64               `json:"height"`
}

type AuthTokenReply struct {
	Token *passport.AuthToken `json:"token"`
}

type OrgTokenReply struct {
	Tokens []*passport.AuthToken `json:"tokens"`
}

//****** Transaction ******

type CreateTokenRequest struct {
	User             passport.UserID       `json:"userId"`
	TokenTypeID      passport.TokenTypeID  `json:"tokenTypeId"`
	TokenType        passport.TokenType    `json:"tokenType"`
	TokenSubType     passport.TokenSubType `json:"tokenSubType"`
	OwnerAddress     keys.Address          `json:"ownerAddress"`
	SuperUserAddress keys.Address          `json:"superUserAddress"`
	SuperUser        passport.UserID       `json:"superUser"`
	CreationTime     string                `json:"creationTime"`
}

type CreateTokenReply struct {
	RawTx []byte `json:"rawTx"`
}

type AddTestInfoRequest struct {
	TestID       string               `json:"testId"`
	Person       passport.UserID      `json:"person"`
	Test         passport.TestType    `json:"test"`
	SubTest      passport.TestSubType `json:"subTest"`
	Manufacturer string               `json:"manufacturer"`
	Result       passport.TestResult  `json:"testResult"`

	TestOrg     passport.TokenTypeID `json:"testOrg"`
	TestedAt    string               `json:"testedAt"`
	TestedBy    string               `json:"testedBy"`
	AnalysisOrg passport.TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string               `json:"analyzedAt"`
	AnalyzedBy  string               `json:"analyzedBy"`

	Admin        passport.UserID `json:"admin"`
	AdminAddress keys.Address    `json:"adminAddress"`
	UploadedAt   string          `json:"uploadedAt"`
	Notes        string          `json:"notes"`
}

type UpdateTestInfoRequest struct {
	TestID       string               `json:"testId"`
	Person       passport.UserID      `json:"person"`
	Test         passport.TestType    `json:"test"`
	Result       passport.TestResult  `json:"testResult"`

	AnalysisOrg passport.TokenTypeID `json:"analysisOrg"`
	AnalyzedAt  string               `json:"analyzedAt"`
	AnalyzedBy  string               `json:"analyzedBy"`

	Admin        passport.UserID `json:"admin"`
	AdminAddress keys.Address    `json:"adminAddress"`
	Notes        string          `json:"notes"`
}

type ReadTestInfoRequest struct {
	Org          passport.TokenTypeID `json:"org"`
	Admin        passport.UserID      `json:"admin"`
	AdminAddress keys.Address         `json:"adminAddress"`
	Person       passport.UserID      `json:"person"`
	Address      keys.Address         `json:"address"`
	Test         passport.TestType    `json:"test"`
	ReadAt       string               `json:"readAt"`
}
