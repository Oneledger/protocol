import sys

from rpc_call import *


def query_rewards(validator):
    req = {
        "validator": validator,
    }

    resp = rpc_call('query.ListRewardsForValidator', req)

    if "result" in resp:
        result = resp["result"]
    else:
        result = ""

    # print json.dumps(resp, indent=4)
    return result


def withdraw_rewards(validator):
    req = {
        "validator": validator,
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.WithdrawRewards', req)
    if "result" in resp:
        result = resp["result"]
    else:
        result = ""
    return result


def list_validators():
    resp = rpc_call('query.ListValidators', {})
    result = resp["result"]

    # print json.dumps(resp, indent=4)
    return result


def sign(raw_tx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": raw_tx, "address": address})
    return resp["result"]


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]


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


class Withdraw:
    def __init__(self, valindatorWalletAddress, amount):
        self.address = valindatorWalletAddress
        self.withdrawAmount = amount

    def _withdraw_reward(self):
        req = {
            "validatorSigningAddress": self.address,
            "withdrawAmount": {
                "currency": "OLT",
                "value": self.withdrawAmount,
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.WithdrawRewards', req)
        result = resp["result"]
        return result["rawTx"]

    def send_withdraw(self):
        # create TX , with Validator address from context
        raw_txn = self._withdraw_reward()
        # sign Tx , with validators wallet account to pay for olt
        signed = sign(raw_txn, self.address)
        # broadcast this signed TX .
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Not OK"
                sys.exit(-1)
            else:
                print "################### Withdraw completed : "
                return result["txHash"]
