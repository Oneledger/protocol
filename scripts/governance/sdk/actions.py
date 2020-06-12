import hashlib
import sys

from rpc_call import *

#Proposal Types
ProposalTypeConfigUpdate = 0x20
ProposalTypeCodeChange   = 0x21
ProposalTypeGeneral      = 0x22

#Proposal Status
ProposalStatusFunding    = 0x23
ProposalStatusVoting     = 0x24
ProposalStatusCompleted  = 0x25

#Proposal Outcome
ProposalOutcomeInProgress         = 0x26
ProposalOutcomeInsufficientFunds  = 0x27
ProposalOutcomeInsufficientVotes  = 0x28
ProposalOutcomeCancelled          = 0x29
ProposalOutcomeCompleted          = 0x30


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
    def __init__(self, pid, pType, description, headline, proposer, init_fund):
        self.pid = pid
        self.pty = pType
        self.headline = headline
        self.des = description
        self.proposer = proposer
        self.init_fund = init_fund

    def _create_proposal(self):
        req = {
            "proposalId": self.pid,
            "headline": self.headline,
            "description": self.des,
            "proposer": self.proposer,
            "proposalType": self.pty,
            "initialFunding": {
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
            "proposalId": self.pid,
            "fundValue": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "funderAddress": self.funder,
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

class ProposalCancel:
    def __init__(self, pid, proposer, reason):
        self.pid = pid
        self.proposer = proposer
        self.reason = reason

    def _cancel_proposal(self):
        req = {
            "proposalId": self.pid,
            "proposer": self.proposer,
            "reason": self.reason,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
    
        resp = rpc_call('tx.CancelProposal', req)
        return resp["result"]["rawTx"]

    def send_cancel(self):
        # create Tx
        raw_txn = self._cancel_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                print "################### failed to cancel proposal: " + self.pid
                return False
            else:
                print "################### proposal canceled: " + self.pid
                return True
        else:
            print "################### failed to cancel proposal: " + self.pid
            return False

class ProposalVote:
    def __init__(self, pid, opinion, url, address):
        self.pid = pid
        self.opin = opinion
        self.voter = url
        self.address = address

    def _vote_proposal(self):
        req = {
            "proposalId": self.pid,
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
    def __init__(self, pid, funder, value, beneficiary):
        self.pid = pid
        self.funder = funder
        self.value = value
        self.benefi = beneficiary

    def _withdraw_funds(self, funder_address):
        req = {
            "proposalId": self.pid,
            "funderAddress": funder_address,
            "withdrawValue": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "beneficiaryAddress": self.benefi,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.WithdrawProposalFunds', req)
        print resp
        return resp["result"]["rawTx"]

    def withdraw_fund(self, contr_address):
        # create Tx
        raw_txn = self._withdraw_funds(contr_address)

        # sign Tx
        signed = sign(raw_txn, self.funder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                print bcolors.FAIL + "################### proposal funds withdraw failed:" + result["log"] + bcolors.ENDC
                return result["txHash"]
            else:
                print "################### proposal funds withdrawn:" + self.pid
                return result["txHash"]
        else:
            print bcolors.FAIL + "################### proposal funds withdraw failed:" + result["error"]["message"] + bcolors.ENDC



class ProposalFinalize:
    def __init__(self, pid, address):
        self.pid = pid
        self.proposer = address

    def _finalize_proposal(self):
        req = {
            "proposalId": self.pid,
            "proposer": self.proposer,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.FinalizeProposal', req)
        return resp["result"]["rawTx"]

    def send_finalize(self):
        # create Tx
        raw_txn = self._finalize_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal finalized: " + self.pid
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

def query_proposals(prefix, proposer="", proposal_type=""):
    req = {
        "state": prefix,
        "proposer": proposer,
        "proposalType": proposal_type,
    }

    resp = rpc_call('query.ListProposals', req)
    result = resp["result"]
    print json.dumps(result, indent=4)
    return result["proposalStats"]

def query_proposal(proposal_id):
    req = {
        "proposalId": proposal_id,
    }
    resp = rpc_call('query.ListProposal', req)

    stat = resp["result"]["proposalStats"][0]
    print json.dumps(stat, indent=4)
    return stat["proposal"], stat["funds"]

def query_balance(address):
    req = {"address": address}
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    return resp["result"]
