from actions import *
from rpc_call import *
from setupAccounts import *

# User ID
UserIvalid = ""

# Test types
COVID19 = 0x101

# Test sub types
TestSubAntiBody = 0x201
TestSubAntigen = 0x202
TestSubPCR = 0x203

# Test results
COVID19Positive = 0x301
COVID19Negative = 0x302
COVID19Pending = 0x303

# Token types
TokenSuperAdmin = 0x401
TokenHospital = 0x402
TokenScreener = 0x403

# Screener types
ScreenerInvalid = 0x500
ScreenerGeneral = 0x501
ScreenerBorderService = 0x502
ScreenerUniversity = 0x503

# Token type id (org ID)
TypeIDInvalid = ""
TypeIDSuperAdmin = "SuperAdminGroup"

# Token roles
RoleSuperAdmin = 0x601
RoleHospitalAdmin = 0x602
RoleScreenerAdmin = 0x603

# Token permissions
PermitUpload = 0x00000001
PermitScan = 0x00000002
PermitQueryTest = 0x00000004
PermitQueryRead = 0x00000008
PermitSuper = 0xffffffff
