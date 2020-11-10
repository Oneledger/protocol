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

    def send_network_Delegate(self, exit_on_err=True, mode=TxCommit):
        # create Tx
        raw_txn = self._network_Delegate()

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "################### delegation failed"
                if exit_on_err:
                    sys.exit(-1)
            else:
                print "################### delegation added"
        return "" if "log" not in result else result["log"]

    def send_network_undelegate(self, amount, exit_on_err=True, mode=TxCommit):
        # createTx
        raw_txn = self._network_undelegate(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "Send undelegate Failed : ", result
                if exit_on_err:
                    sys.exit(-1)
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### undelegate"
        return "" if "log" not in result else result["log"]

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

    def send_network_withdraw(self, amount, exit_on_err=True, mode=TxCommit):
        # createTx
        raw_txn = self._network_withdraw(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegationaddress, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "Send withdraw delegation Failed : ", result
                if exit_on_err:
                    sys.exit(-1)
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### withdraw delegation"
        return result["log"]

    def waitfor_matured(self, amount):
        req = {
            "delegationAddress": '0lt' + self.delegationaddress
        }
        def until(result):
            matured = result["delegationStats"]["matured"]
            matured_olt = matured.split(" ")[0]
            return int(matured_olt) >= int(amount)
        wait_until("query.ListDelegation", req, until)

def waitfor_rewards(delegator, amount, status):
    req = {
        "delegator": delegator,
        "inclPending": False,
    }
    amount += "0"*18
    def until(result):
        actual = result[status]
        return int(actual) >= int(amount)
    result = wait_until("query.GetDelegRewards", req, until)
    actual = (int(result[status]) / 10 ** 18)
    return actual

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
        return resp["result"]["rawTx"]

    def send(self, exit_on_err=True, mode=TxCommit):
        # create Tx
        raw_txn = self._request()

        # sign Tx
        signed = sign(raw_txn, self.delegator, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "################### withdrawal initiation failed"
                if exit_on_err:
                    sys.exit(-1)
            else:
                print "################### withdrawal successfully initiated"
        return "" if "log" not in result else result["log"]

class ReinvestRewards:
    def __init__(self, delegator, keypath):
        self.delegator = delegator
        self.keypath = keypath

    def _request(self, amount):
        req = {
            "delegator": self.delegator,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(amount),
            },
        }
        resp = rpc_call('tx.ReinvestDelegRewards', req)
        return resp["result"]["rawTx"]

    def send(self, amount, exit_on_err=True, mode=TxCommit):
        # create Tx
        raw_txn = self._request(amount)

        # sign Tx
        signed = sign(raw_txn, self.delegator, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"]:
                print "################### failed to reinvest rewards"
                if exit_on_err:
                    sys.exit(-1)
            else:
                print "################### successfully reinvested rewards"
        return "" if "log" not in result else result["log"]

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

    def send_finalize(self, finalize_amount, exit_on_err=True, mode=TxCommit):
        # create Tx
        raw_txn = self._request_finalize(finalize_amount)

        # sign Tx
        signed = sign(raw_txn, self.delegator, self.keypath)

        # broadcast Tx
        result = broadcast(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'], mode)
        if "ok" in result:
            if not result["ok"] and exit_on_err:
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


def query_delegation(delegation_addresses=None):
    if (delegation_addresses == None):
        delegation_addresses = []
    req = {
        "delegationAddresses": delegation_addresses,
    }

    resp = rpc_call('query.ListDelegation', req)
    # print resp
    result = resp["result"]
    return result["allDelegStats"]


def check_query_delegation(query_result, index, expected_delegation, expected_pending_delegation, have_reward, expected_pending_rewards):
    actual_delegation = query_result[index]['delegationStats']['active']
    actual_pending_delegation = query_result[index]['delegationStats']['pending']
    actual_rewards = query_result[index]['delegationRewardsStats']['active']
    actual_pending_rewards = query_result[index]['delegationRewardsStats']['pending']

    if int(actual_delegation) != expected_delegation:
        print "actual_delegation"
        print actual_delegation
        print "expected_delegation"
        print expected_delegation
        sys.exit(-1)

    if int(actual_pending_delegation) != expected_pending_delegation:
        print "actual_pending_delegation"
        print actual_pending_delegation
        print "expected_pending_delegation"
        print expected_pending_delegation
        sys.exit(-1)

    if int(actual_rewards) == 0 and have_reward:
        print "no rewards found!"
        sys.exit(-1)

    if int(actual_pending_rewards) != expected_pending_rewards:
        print "actual_pending_rewards"
        print actual_pending_rewards
        print "expected_pending_rewards"
        print expected_pending_rewards
        sys.exit(-1)
