import os.path as path
import sys

sdkcom_path = path.abspath(path.join(path.dirname(__file__), "../.."))
sys.path.append(sdkcom_path)
from sdkcom import *


class NetWorkDelegate:
    def __init__(self, delegationaddress, amount, keypath):
        self.delegationaddress = delegationaddress
        self.amount = amount
        self.keypath = keypath

    def _network_Delegate(self):
        req = {
            "delegationAddress": self.delegationaddress,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(self.amount),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }

        resp = rpc_call('tx.AddNetworkDelegation', req)
        print resp
        return resp["result"]["rawTx"]

    def query_delegation(self):
        req = {
            "delegationAddress": self.delegationaddress,
        }
        resp = rpc_call('query.ListDelegation', req)
        print resp
        result = resp["result"]
        # print json.dumps(resp, indent=4)
        return result["delegationStats"]

    def _network_undelegate(self, amount):
        req = {
            "delegator": self.delegationaddress,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(amount),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }

        resp = rpc_call('tx.NetworkUndelegate', req)
        print resp
        return resp["result"]["rawTx"]

    def send_network_Delegate(self):
        # create Tx
        raw_txn = self._network_Delegate()

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### delegation added"
                return result["txHash"]

    def send_network_undelegate(self, amount):
        # createTx
        raw_txn = self._network_undelegate(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send undelegate Failed : ", result
                sys.exit(-1)
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### undelegate"

    def send_network_undelegate_shoud_fail(self, amount):
        # createTx
        raw_txn = self._network_undelegate(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if result["ok"]:
                print "Send undelegate should fail, but it doesn't : ", result
                sys.exit(-1)
        print "################### malicious undelegate failed as expected"

    def query_undelegate(self):
        req = {
            "delegator": '0lt' + self.delegationaddress
        }

        resp = rpc_call('query.GetUndelegatedAmount', req)
        print json.dumps(resp, indent=4)
        result = resp["result"]
        return result


class WithdrawRewards:
    def __init__(self, delegator, amount, keypath):
        self.delegator = delegator
        self.amount = amount
        self.keypath = keypath

    def _request(self):
        req = {
            "delegator": self.delegator,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(self.amount),
            },
        }
        resp = rpc_call('tx.WithdrawDelegRewards', req)
        print resp
        return resp["result"]["rawTx"]

    def send(self, expect_succeed=True):
        # create Tx
        raw_txn = self._request()

        # sign Tx
        signed = sign(raw_txn, self.delegator, self.keypath)

        # broadcast Tx
        result = broadcast_sync(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"] and expect_succeed:
                sys.exit(-1)
            else:
                print "################### withdrawal successfully initiated: "


class FinalizeRewards:
    def __init__(self, delegator, keypath):
        self.delegator = delegator
        self.keypath = keypath
    def _request_finalize(self, finalize_amount):
        req = {
            "delegator": self.delegator,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(finalize_amount),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.FinalizeDelegRewards', req)
        print resp
        return resp["result"]["rawTx"]

    def send_finalize(self, finalize_amount, expect_succeed=True):
        # create Tx
        raw_txn = self._request_finalize(finalize_amount)

        # sign Tx
        signed = sign(raw_txn, self.delegator, self.keypath)

        # broadcast Tx
        result = broadcast_sync(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"] and expect_succeed:
                sys.exit(-1)
            else:
                print "################### finalize rewards sent"


def sign(raw_tx, address, keypath):
    resp = rpc_call('owner.SignWithSecureAddress',
                    {"rawTx": raw_tx, "address": address, "password": "1234", "keypath": keypath})
    print resp
    return resp["result"]


def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    if "result" in resp:
        return resp["result"]
    else:
        return resp


def broadcast_sync(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    return resp["result"]


def query_total(only_active):
    req = {
        "onlyActive": only_active
    }
    resp = rpc_call('query.GetTotalNetwkDelegation', req)
    print json.dumps(resp, indent=4)
    result = resp["result"]
    return result


def check_query_undelegated(result, pending_count_expected):
    if not result['height']:
        sys.exit(-1)
    if len(result['pendingAmount']) != pending_count_expected:
        sys.exit(-1)


def check_query_total(result, total_expected):
    if not result['height']:
        sys.exit(-1)
    if result['totalAmount'] != total_expected:
        sys.exit(-1)


def query_rewards(delegator):
    req = {
        "delegator": delegator,
        "inclPending": True,
    }
    resp = rpc_call('query.GetDelegRewards', req)
    print json.dumps(resp, indent=4)
    if "result" in resp:
        result = resp["result"]
    else:
        result = ""
    return result
