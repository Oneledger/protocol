import os.path as path
import sys

sdkcom_p = path.abspath(path.join(path.dirname(__file__), "../.."))
sys.path.append(sdkcom_p)

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
        return resp["result"]["rawTx"]

    def query_delegation(self):
        req = {
            "delegationAddress": self.delegationaddress,
        }
        resp = rpc_call('query.ListDelegation', req)
        print resp
        result = resp["result"]
        # prmodeps(resp, indent=4)
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
        print (resp)
        return resp["result"]["rawTx"]

    def _network_withdraw(self, amount):
        req = {
            "delegationAddress": self.delegationaddress,
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
        resp = rpc_call('tx.WithDrawNetworkDelegation', req)
        print resp
        return resp["result"]["rawTx"]

    def send_network_Delegate(self, mode=TxCommit):
        # create Tx
        raw_txn = self._network_Delegate()

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "################### delegation failed"
            else:
                print "################### delegation added"
        return result["log"]

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

    def send_network_withdraw(self, amount):
        # createTx
        raw_txn = self._network_withdraw(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send withdraw delegation Failed : ", result
                sys.exit(-1)
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### withdraw delegation"


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


def query_total_rewards():
    req = {}
    resp = rpc_call('query.GetTotalDelegRewards', req)
    print json.dumps(resp, indent=4)
    if "result" in resp:
        result = resp["result"]
    else:
        result = ""
    return result


def query_full_balance(address):
    req = {
        "address": address
    }
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    if "result" in resp:
        result = resp["result"]
    else:
        result = ""
    return result
