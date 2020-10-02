import hashlib
import sys

from rpc_call import *

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


class Undelegate:
    def __init__(self, delegator, amount):
        self.delegator = delegator
        self.amount = amount

    def _create_undelegate(self):
        req = {
            "delegator": self.delegator,
            "amount": self.amount,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }

        resp = rpc_call('tx.NetUndelegate', req)
        return resp["result"]["rawTx"]

    def send_create(self):
        # createTx
        raw_txn = self._create_undelegate()

        # sign Tx
        signed = sign(raw_txn, self.delegator)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send undelegate Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### undelegate success"


def sign(raw_tx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": raw_tx, "address": address})
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
    return resp["result"]


def query_undelegate():
    req = {
        "state": prefix,
        "proposer": proposer,
        "proposalType": proposalType,
    }

    resp = rpc_call('query.GetUndelegatedAmount', req)
    # print resp
    result = resp["result"]
    # print json.dumps(resp, indent=4)
    return result["proposalStats"]


def query_proposal(proposal_id):
    req = {
        "proposalId": proposal_id,
    }
    resp = rpc_call('query.ListProposal', req)
    stat = resp["result"]["proposalStats"][0]
    # print json.dumps(resp, indent=4)
    return stat["proposal"], stat["funds"]


def query_governanceState():
    req = {}
    resp = rpc_call('query.GetGovernanceOptionsForHeight', req)
    result = resp["result"]

    # print json.dumps(resp, indent=4)
    return result


def query_balance(address):
    req = {"address": address}
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    return resp["result"]


def query_proposal_options():
    req = {}
    resp = rpc_call('query.GetProposalOptions', req)

    if "result" not in resp:
        sys.exit(-1)
    if "proposalOptions" not in resp["result"]:
        sys.exit(-1)
    if "height" not in resp["result"]:
        sys.exit(-1)
    # print json.dumps(resp, indent=4)
    return resp["result"]


def get_funds_for_proposal_by_funder(proposalId, funder):
    req = {
        "proposalId": proposalId,
        "funderAddress": funder
    }
    resp = rpc_call('query.GetFundsForProposalByFunder', req)
    if "result" not in resp:
        sys.exit(-1)

    return resp["result"]
