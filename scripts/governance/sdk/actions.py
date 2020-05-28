import json
import sys
import hashlib
from rpc_call import *

#Proposal Status
ProposalStatusFunding    = 0x23
ProposalStatusVoting     = 0x24
ProposalStatusCompleted  = 0x25

#Proposal States
ProposalStateError   = 0xEE
ProposalStateActive  = 0x31
ProposalStatePassed  = 0x32
ProposalStateFailed  = 0x33

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


class Proposal:
    def __init__(self, pid, pType, description, proposer, init_fund):
        self.pid = pid
        self.pty = pType
        self.des = description
        self.proposer = proposer
        self.init_fund = init_fund

    def _create_proposal(self):
        req = {
            "proposal_id": self.pid,
            "description": self.des,
            "proposer": self.proposer,
            "proposal_type": self.pty,
            "initial_funding": {
                "currency": "OLT",
                "value": convertBigInt(self.init_fund),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.CreateProposal', req)
        return resp["result"]["rawTx"]

    def send_create(self):
        # createTx
        raw_txn = self._create_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                self.pid = self.get_encoded_pid()
                self.txHash = "0x" + result["txHash"]
                print "################### proposal created: " + self.pid

    def get_encoded_pid(self):
        hash_handler = hashlib.md5()
        hash_handler.update(self.pid)
        hash_val = hash_handler.digest()
        return hash_val.encode('hex')

    def tx_created(self):
        resp = tx_by_hash(self.txHash)
        return resp["result"]["tx_result"]


class ProposalFund:
    def __init__(self, pid, value, address):
        self.pid = pid
        self.value = value
        self.funder = address

    def _fund_proposal(self):
        req = {
            "proposal_id": self.pid,
            "fund_value": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "funder_address": self.funder,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
    
        resp = rpc_call('tx.FundProposal', req)
        return resp["result"]["rawTx"]

    def send_fund(self):
        # create Tx
        raw_txn = self._fund_proposal()

        # sign Tx
        signed = sign(raw_txn, self.funder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal funded: " + self.pid
                return result["txHash"]


class ProposalVote:
    def __init__(self, pid, opinion, url, address):
        self.pid = pid
        self.opin = opinion
        self.voter = url
        self.address = address

    def _vote_proposal(self):
        req = {
            "proposal_id": self.pid,
            "opinion": self.opin,
            "address": self.address,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.VoteProposal', req, self.voter)
        result = resp["result"]
        return result["rawTx"], result['signature']['Signed'], result['signature']['Signer']

    def send_vote(self):
        # create and let validator sign Tx
        raw_txn, signed_0, signer_0 = self._vote_proposal()

        # payer sign Tx
        res = sign(raw_txn, self.address)
        signed_1 = res['signature']['Signed']
        signer_1 = res['signature']['Signer']

        # signatures
        sig0 = {"Signer": signer_0, "Signed": signed_0}
        sig1 = {"Signer": signer_1, "Signed": signed_1}
        sigs = [sig1, sig0]

        # broadcast Tx
        result = broadcast_commit_mtsig(raw_txn, sigs)
        
        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal voted:" + self.pid + "opinion: " + self.opin
                return result["txHash"]


class ProposalFundsWithdraw:
    def __init__(self, pid, contributor, value, beneficiary):
        self.pid = pid
        self.contr = contributor
        self.value = value
        self.benefi = beneficiary

    def _withdraw_funds(self):
        req = {
            "proposal_id": self.pid,
            "contributor_address": self.contr,
            "withdraw_value": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "beneficiary_address": self.benefi,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.WithdrawProposalFunds', req)
        print resp
        return resp["result"]["rawTx"]

    def withdraw_fund(self):
        # create Tx
        raw_txn = self._withdraw_funds()

        # sign Tx
        signed = sign(raw_txn, self.contr)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal funds withdrawed:" + self.pid
                return result["txHash"]


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]


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

def broadcast_commit_mtsig(raw_tx, sigs):
    resp = rpc_call('broadcast.TxCommitMtSig', {
        "rawTx": raw_tx,
        "signatures": sigs,
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

def query_proposals(prefix):
    req = {
        "prefix": prefix,
        "gasPrice":
        {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }

    resp = rpc_call('query.GetProposals', req)
    print json.dumps(resp, indent=4)
    return resp["result"]["proposals"]

def query_proposal(proposal_id):
    req = {"proposal_id": proposal_id}
    resp = rpc_call('query.GetProposalByID', req)
    return resp["result"]


def query_balance(address):
    req = {"address": address}
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    return resp["result"]
