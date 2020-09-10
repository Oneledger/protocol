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


def create_domain(name, owner_hex, price):
    req = {
        "name": name,
        "owner": owner_hex,
        "account": owner_hex,
        "buyingPrice": {
            "currency": "OLT",
            "value": convertBigInt(price),
        },
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.ONS_CreateRawCreate', req)
    return resp["result"]["rawTx"]

class BidConv:
    def __init__(self, owner, asset, assetType, bidder, amount, counter_amount, counter_bid_amount, deadline, bidConvId=None):
        self.bidConvId = bidConvId
        self.owner = owner
        self.asset = asset
        self.assetType = assetType
        self.bidder = bidder
        self.amount = amount
        self.counter_amount = counter_amount
        self.counter_bid_amount = counter_bid_amount
        self.deadline = deadline

    def _create_bid(self):
        req = {
            "bidConvId": self.bidConvId,
            "assetOwner": self.owner,
            "assetName": self.asset,
            "assetType": self.assetType,
            "bidder": self.bidder,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(self.amount),
            },
            "deadline": self.deadline,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.CreateBid', req)
        return resp["result"]["rawTx"]

    def _cancel_bid(self, id):
        req = {
            "bidConvId": id,
            "bidder": self.bidder,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.CancelBid', req)
        return resp["result"]["rawTx"]

    def _counter_offer(self, id):
        req = {
            "bidConvId": id,
            "assetOwner": self.owner,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(self.counter_amount),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.CounterOffer', req)
        return resp["result"]["rawTx"]

    def _bidder_decision(self, id, decision):
        req = {
            "bidConvId": id,
            "bidder": self.bidder,
            "decision": decision,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.BidderDecision', req)
        return resp["result"]["rawTx"]

    def _add_bid_offer(self, id):
        req = {
            "bidConvId": id,
            "bidder": self.bidder,
            "amount": {
                "currency": "OLT",
                "value": convertBigInt(self.counter_bid_amount),
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.CreateBid', req)
        return resp["result"]["rawTx"]

    def _owner_decision(self, id, decision):
        req = {
            "bidConvId": id,
            "owner": self.owner,
            "decision": decision,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('bid_tx.OwnerDecision', req)
        return resp["result"]["rawTx"]

#     def _create_proposal_invalid_info(self, invalid_field):
#         _proposal_info = self._calculate_proposal_info(12)
#         req = {
#             "proposalId": self.get_encoded_pid(),
#             "headline": self.headline,
#             "description": self.des,
#             "proposer": self.proposer,
#             "proposalType": self.pty,
#             "initialFunding": {
#                 "currency": "OLT",
#                 "value": convertBigInt(self.init_fund),
#             },
#             "fundingGoal": _proposal_info.funding_goal,
#             "fundingDeadline": _proposal_info.funding_deadline,
#             "votingDeadline": _proposal_info.voting_deadline,
#             "passPercentage": _proposal_info.pass_percentage,
#             "gasPrice": {
#                 "currency": "OLT",
#                 "value": "1000000000",
#             },
#             "gas": 40000,
#         }
#         if invalid_field == 0:
#             req["fundingGoal"] = "123"
#         elif invalid_field == 1:
#             req["fundingDeadline"] = 0
#         elif invalid_field == 2:
#             req["votingDeadline"] = 0
#         elif invalid_field == 3:
#             req["passPercentage"] = 1
#
#         resp = rpc_call('tx.CreateProposal', req)
#         return resp["result"]["rawTx"]
#
    def send_create(self):
        # createTx
        raw_txn = self._create_bid()

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            print result
            if not result["ok"]:
                print "Send Create Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### BidConv Created"

    def send_cancel(self, id):
        # createTx
        raw_txn = self._cancel_bid(id)

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send Cancel Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### BidConv Cancelled"

    def send_counter_offer(self, id):
        # createTx
        raw_txn = self._counter_offer(id)

        # sign Tx
        signed = sign(raw_txn, self.owner)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send Counter Offer Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### Counter Offer Sent"

    def send_bidder_decision(self, id, decision):
        # createTx
        raw_txn = self._bidder_decision(id, decision)

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send Bidder Decision Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### Bidder Decision Sent"

    def send_bid_offer(self, id):
        # createTx
        raw_txn = self._add_bid_offer(id)

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            print result
            if not result["ok"]:
                print "Add Bid Offer Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### Bid Offer Added"

    def send_owner_decision(self, id, decision):
        # createTx
        raw_txn = self._owner_decision(id, decision)

        # sign Tx
        signed = sign(raw_txn, self.owner)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print "Send Owner Decision Failed : ", result
            else:
                self.txHash = "0x" + result["txHash"]
                print "################### Owner Decision Sent"

#
#     def send_create_invalid_info(self, invalid_field):
#         # createTx
#         raw_txn = self._create_proposal_invalid_info(invalid_field)
#
#         # sign Tx
#         signed = sign(raw_txn, self.proposer)
#
#         # broadcast Tx
#         result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
#
#         if "ok" in result:
#             if result["ok"]:
#                  sys.exit(-1)
#
#     def get_encoded_pid(self):
#         hash_handler = hashlib.sha256()
#         hash_handler.update(self.pid)
#         hash_val = hash_handler.digest()
#         return hash_val.encode('hex')
#
#     def tx_created(self):
#         resp = tx_by_hash(self.txHash)
#         return resp["result"]["tx_result"]
#
#
# class ProposalInfo:
#     def __init__(self, funding_goal, funding_deadline, voting_deadline, pass_percentage):
#         self.funding_goal = funding_goal
#         self.funding_deadline = funding_deadline
#         self.voting_deadline = voting_deadline
#         self.pass_percentage = pass_percentage
#
#
# class ProposalFund:
#     def __init__(self, pid, value, address):
#         self.pid = pid
#         self.value = value
#         self.funder = address
#
#     def _fund_proposal(self):
#         req = {
#             "proposalId": self.pid,
#             "fundValue": {
#                 "currency": "OLT",
#                 "value": convertBigInt(self.value),
#             },
#             "funderAddress": self.funder,
#             "gasPrice": {
#                 "currency": "OLT",
#                 "value": "1000000000",
#             },
#             "gas": 40000,
#         }
#
#         resp = rpc_call('tx.FundProposal', req)
#         return resp["result"]["rawTx"]
#
#     def send_fund(self):
#         # create Tx
#         raw_txn = self._fund_proposal()
#
#         # sign Tx
#         signed = sign(raw_txn, self.funder)
#
#         # broadcast Tx
#         result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
#
#         if "ok" in result:
#             if not result["ok"]:
#                 sys.exit(-1)
#             else:
#                 print "################### proposal funded: " + self.pid
#                 return result["txHash"]
#
#
# class ProposalCancel:
#     def __init__(self, pid, proposer, reason):
#         self.pid = pid
#         self.proposer = proposer
#         self.reason = reason
#
#     def _cancel_proposal(self):
#         req = {
#             "proposalId": self.pid,
#             "proposer": self.proposer,
#             "reason": self.reason,
#             "gasPrice": {
#                 "currency": "OLT",
#                 "value": "1000000000",
#             },
#             "gas": 40000,
#         }
#
#         resp = rpc_call('tx.CancelProposal', req)
#         return resp["result"]["rawTx"]
#
#     def send_cancel(self):
#         # create Tx
#         raw_txn = self._cancel_proposal()
#
#         # sign Tx
#         signed = sign(raw_txn, self.proposer)
#
#         # broadcast Tx
#         result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
#
#         if "ok" in result:
#             if not result["ok"]:
#                 print "################### failed to cancel proposal: " + self.pid
#                 return False
#             else:
#                 print "################### proposal canceled: " + self.pid
#                 return True
#         else:
#             print "################### failed to cancel proposal: " + self.pid
#             return False
#
#
# class ProposalFinalize:
#     def __init__(self, pid, address):
#         self.pid = pid
#         self.proposer = address
#
#     def _finalize_proposal(self):
#         req = {
#             "proposalId": self.pid,
#             "proposer": self.proposer,
#             "gasPrice": {
#                 "currency": "OLT",
#                 "value": "1000000000",
#             },
#             "gas": 40000,
#         }
#         resp = rpc_call('tx.FinalizeProposal', req)
#         return resp["result"]["rawTx"]
#
#     def send_finalize(self):
#         # create Tx
#         raw_txn = self._finalize_proposal()
#
#         # sign Tx
#         signed = sign(raw_txn, self.proposer)
#
#         # broadcast Tx
#         result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
#
#         if "ok" in result:
#             if not result["ok"]:
#                 sys.exit(-1)
#             else:
#                 print "################### proposal finalized: " + self.pid
#                 return result["txHash"]


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


def query_bidConvs(state, owner, assetName, assetType, bidder):
    req = {
        "state": state,
        "owner": owner,
        "assetName": assetName,
        "assetType": assetType,
        "bidder": bidder,
    }

    resp = rpc_call('bid_query.ListBidConvs', req)
    # print resp
    result = resp["result"]
    print json.dumps(resp, indent=4)
    return result


def query_bidConv(bidConv_Id):
    req = {
        "bidConvId": bidConv_Id,
    }
    resp = rpc_call('bid_query.ShowBidConv', req)
    result = resp["result"]
    print json.dumps(resp, indent=4)
    return result


# def query_balance(address):
#     req = {"address": address}
#     resp = rpc_call('query.Balance', req)
#     print json.dumps(resp, indent=4)
#     return resp["result"]
#
#
#
# def get_funds_for_proposal_by_funder(proposalId, funder):
#     req = {
#         "proposalId": proposalId,
#         "funderAddress": funder
#     }
#     resp = rpc_call('query.GetFundsForProposalByFunder', req)
#     if "result" not in resp:
#         sys.exit(-1)
#
#     return resp["result"]
