import sys
from datetime import datetime, timedelta

import pytz

from rpc_call import *

TOKEN_HOSPITAL = 0x402
TOKEN_SCREENER = 0x403

TransactionTypes = {
    TOKEN_HOSPITAL: "tx.CreateHospitalAdminToken",
    TOKEN_SCREENER: "tx.CreateScreenerAdminToken"
}

# RFC3339 time stamps
now = pytz.UTC.localize(datetime.utcnow())
creation_timestamp = (now - timedelta(hours=5)).isoformat('T')
TZero = "0001-01-01T00:00:00Z"


class Colours:
    def __init__(self):
        pass

    HEADER = '\033[95m'
    OK_BLUE = '\033[94m'
    OK_GREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    END_C = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]


def add_account(name):
    resp = rpc_call('owner.GenerateNewAccount', {'name': name})


def get_auth_tokens(org, user_id):
    req = {
        "org": org,
        "id": user_id
    }
    resp = rpc_call('query.GetTokensByOrganization', req)
    return resp['result']


def read_test_info(test, org, admin, person):
    req = {
        "org": org,
        "admin": admin,

        "test": test,
        "uploadedOrg": "",
        "uploadedBy": "",
        "person": person,
        "timeBegin": "",
        "timeEnd": ""
    }
    resp = rpc_call('query.PSPT_QueryTestInfoByID', req)
    return resp['result']['infoList']


def read_test_info_from_user(test, person):
    req = {
        "test": test,
        "person": person,
    }
    resp = rpc_call('query.PSPT_QueryTestInfoForUser', req)
    return resp['result']['infoList']


def filter_test_info(test, org, admin, person, time_begin, time_end, uploaded_org, uploaded_by, analysis_org="",
                     analyzed_by="", analyze_time_begin="", analyze_time_end=""):
    req = {
        "org": org,
        "admin": admin,

        "test": test,
        "uploadedOrg": uploaded_org,
        "uploadedBy": uploaded_by,
        "person": person,
        "testTimeBegin": time_begin,
        "testTimeEnd": time_end,
        "analysisOrg": analysis_org,
        "analyzedBy": analyzed_by,
        "analyzeTimeBegin": analyze_time_begin,
        "analyzeTimeEnd": analyze_time_end
    }
    resp = rpc_call('query.PSPT_FilterTestInfo', req)
    return resp['result']['infoList']


def filter_read_logs(test, org, admin, org_read, read_by, person, time_begin, time_end):
    req = {
        "test": test,
        "org": org,
        "admin": admin,
        "readOrg": org_read,
        "readBy": read_by,
        "person": person,
        "timeBegin": time_begin,
        "timeEnd": time_end,
    }
    resp = rpc_call('query.PSPT_FilterReadLogs', req)
    return resp['result']['logs']


def sign(raw_tx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": raw_tx, "address": address})
    return resp["result"]


def secure_sign(raw_tx, address, password, key_path):
    req = {
        'rawTx': raw_tx,
        'address': address,
        'password': password,
        'keyPath': key_path
    }
    resp = rpc_call('owner.SignWithSecureAddress', req)
    return resp["result"]


def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp


def validate_result(passed, result, description):
    print Colours.HEADER + description + Colours.END_C
    if passed:
        print Colours.OK_BLUE + '########## Expected To PASS ##########' + Colours.END_C
        if 'result' in result:
            print Colours.OK_GREEN + '########## Result Validation Passed ##########' + Colours.END_C
            print result['result']
        else:
            print Colours.FAIL + '########## Result Validation Failed ##########' + Colours.END_C
            print result
            sys.exit(-1)
    else:
        print Colours.OK_BLUE + '########## Expected To FAIL ##########' + Colours.END_C
        if 'error' in result:
            print Colours.OK_GREEN + '########## Result Validation Passed ##########' + Colours.END_C
            print result['error']
        elif 'result' in result and not result['result']['ok']:
            print Colours.OK_GREEN + '########## Test Case Passed ##########' + Colours.END_C
            print result['result']
        else:
            print Colours.FAIL + '########## Result Validation Failed ##########' + Colours.END_C
            print result
            sys.exit(-1)

    print ''


class AdminToken:
    def __init__(self, user, token_type, token_type_id, owner_address, super_user_address, super_user, creation_time):
        self.user = user
        self.tokenTypeID = token_type_id
        self.tokenType = token_type
        self.tokenSubType = token_type + 1
        self.ownerAddress = owner_address
        self.superUserAddress = super_user_address
        self.superUser = super_user
        self.creationTime = creation_timestamp
        self.resp = ""
        self.req = {}

    def build(self):
        self.req = {"userId": self.user,
                    "tokenTypeId": self.tokenTypeID,
                    "tokenType": self.tokenType,
                    "tokenSubType": self.tokenSubType,
                    "ownerAddress": self.ownerAddress,
                    "superUserAddress": self.superUserAddress,
                    "superUser": self.superUser,
                    "creationTime": self.creationTime
                    }

    def send_wrong_signer(self, secure_account):
        self.build()
        self.req['superUserAddress'] = self.ownerAddress
        raw_tx = rpc_call("tx.CreateHospitalAdminToken", self.req)
        signed_tx = secure_sign(raw_tx['result']['rawTx'], self.superUserAddress, secure_account.password,
                                secure_accounts_path + 'keystore/')
        self.resp = broadcast_commit(raw_tx['result']['rawTx'], signed_tx['signature']['Signed'],
                                     signed_tx['signature']['Signer'])

    def send(self, secure_account):
        self.build()
        raw_tx = rpc_call(TransactionTypes[self.tokenType], self.req)
        signed_tx = secure_sign(raw_tx['result']['rawTx'], self.superUserAddress, secure_account.password,
                                secure_accounts_path + 'keystore/')
        self.resp = broadcast_commit(raw_tx['result']['rawTx'], signed_tx['signature']['Signed'],
                                     signed_tx['signature']['Signer'])

    def get_resp(self):
        return self.resp


class TestInfo:
    def __init__(self, test_id="", test_org="", admin="", person="", test=0, test_sub_type=0, manufacturer="", test_result=0,
                 tested_at="", uploaded_at="", tested_by="", notes="", analysis_org="",
                 analyzed_at=TZero, analyzed_by=""):
        self.testId = test_id
        self.personId = person
        self.test = test
        self.subTest = test_sub_type
        self.manufacturer = manufacturer
        self.testResult = test_result
        self.testOrg = test_org
        self.testedAt = tested_at
        self.testedBy = tested_by
        self.analysisOrg = analysis_org
        self.analyzedAt = analyzed_at
        self.analyzedBy = analyzed_by
        self.uploadedBy = admin
        self.uploadedAt = uploaded_at
        self.notes = notes

    def from_dict(self, **entries):
        self.__dict__.update(entries)
        return self

    def __eq__(self, other):
        mineJson = json.dumps(self.__dict__, sort_keys=True)
        otherJson = json.dumps(other.__dict__, sort_keys=True)
        return mineJson == otherJson

    def __ne__(self, other):
        return not self.__eq__(other)

    def build_add_test(self, admin_address):
        return {
            "testId": self.testId,
            "person": self.personId,
            "test": self.test,
            "subTest": self.subTest,
            "manufacturer": self.manufacturer,
            "testResult": self.testResult,
            "testOrg": self.testOrg,
            "testedAt": self.testedAt,
            "testedBy": self.testedBy,
            "admin": self.uploadedBy,
            "adminAddress": admin_address,
            "uploadedAt": self.uploadedAt,
            "notes": self.notes,
        }

    def send_add_test(self, admin_account, admin_address, should_succeed=True):
        req = self.build_add_test(admin_address)
        raw_tx = rpc_call("tx.AddTestInfo", req)
        signed_tx = secure_sign(raw_tx['result']['rawTx'], admin_address, admin_account.password,
                                secure_accounts_path + 'keystore/')
        
        resp = broadcast_commit(raw_tx['result']['rawTx'], signed_tx['signature']['Signed'],
                                signed_tx['signature']['Signer'])
        result = resp["result"]
        if "ok" in result:
            if not result["ok"] and should_succeed:
                sys.exit(-1)
            else:
                print resp

    def build_update_test(self, admin_address):
        return {
            "testId": self.testId,
            "person": self.personId,
            "test": self.test,
            "testResult": self.testResult,
            "analysisOrg": self.analysisOrg,
            "analyzedAt": self.analyzedAt,
            "analyzedBy": self.analyzedBy,
            "admin": self.analyzedBy,
            "adminAddress": admin_address,
            "notes": self.notes,
        }

    def send_update_test(self, admin_account, admin_address, should_succeed=True):
        req = self.build_update_test(admin_address)
        raw_tx = rpc_call("tx.UpdateTestInfo", req)
        signed_tx = secure_sign(raw_tx['result']['rawTx'], admin_address, admin_account.password,
                                secure_accounts_path + 'keystore/')

        resp = broadcast_commit(raw_tx['result']['rawTx'], signed_tx['signature']['Signed'],
                                signed_tx['signature']['Signer'])
        result = resp["result"]
        if "ok" in result:
            if not result["ok"] and should_succeed:
                sys.exit(-1)
            else:
                print resp

    @classmethod
    def read(cls, org, admin, admin_address, person, address, test, read_at, password):
        req = {"org": org,
               "admin": admin,
               "adminAddress": admin_address,
               "person": person,
               "address": address,
               "test": test,
               "readAt": read_at,
               }
        raw_tx = rpc_call("tx.ReadTestInfo", req)
        signed_tx = secure_sign(raw_tx['result']['rawTx'], admin_address, password,
                                secure_accounts_path + 'keystore/')
        resp = broadcast_commit(raw_tx['result']['rawTx'], signed_tx['signature']['Signed'],
                                signed_tx['signature']['Signer'])
        return read_test_info(test, org, admin, person)


class ReadLog:
    def __init__(self, org="", admin="", admin_address="", person="", address="", test=0, read_at=""):
        self.org = org
        self.readBy = admin
        self.adminAddress = admin_address
        self.person = person
        self.address = address
        self.test = test
        self.readAt = read_at

    def from_dict(self, **entries):
        self.__dict__.update(entries)
        return self

    def __eq__(self, other):
        mineJson = json.dumps(self.__dict__, sort_keys=True)
        otherJson = json.dumps(other.__dict__, sort_keys=True)
        return mineJson == otherJson

    def __ne__(self, other):
        return not self.__eq__(other)
