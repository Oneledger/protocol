package passport

import (
	"os"

	"github.com/Oneledger/protocol/log"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("passport")
}

const (
	// Test ID length
	TestIDLength int = 64

	// User ID
	UserIvalid UserID = ""

	// Test types
	TestInvalid TestType = 0x100
	TestCOVID19 TestType = 0x101

	// Test sub types
	TestSubInvalid  TestSubType = 0x100
	TestSubAntiBody TestSubType = 0x201
	TestSubAntigen  TestSubType = 0x202
	TestSubPCR      TestSubType = 0x203

	// Test results
	COVID19Positive TestResult = 0x301
	COVID19Negative TestResult = 0x302
	COVID19Pending  TestResult = 0x303

	// Token types
	TokenSuperAdmin TokenType = 0x401
	TokenHospital   TokenType = 0x402
	TokenScreener   TokenType = 0x403

	// Screener types
	ScreenerInvalid       TokenSubType = 0x500
	ScreenerGeneral       TokenSubType = 0x501
	ScreenerBorderService TokenSubType = 0x502
	ScreenerUniversity    TokenSubType = 0x503

	// Token type id (org ID)
	TypeIDInvalid    TokenTypeID = ""
	TypeIDSuperAdmin TokenTypeID = "SuperAdminGroup"

	// Token roles
	RoleSuperAdmin    TokenRole = 0x601
	RoleHospitalAdmin TokenRole = 0x602
	RoleScreenerAdmin TokenRole = 0x603

	// Token permissions
	PermitUpload      TokenPermission = 0x00000001
	PermitScan        TokenPermission = 0x00000002
	PermitQueryTest   TokenPermission = 0x00000004
	PermitQueryRead   TokenPermission = 0x00000008
	PermitQueryTokens TokenPermission = 0x00000010
	PermitSuper       TokenPermission = 0xffffffff
)
