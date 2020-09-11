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
    def __init__(self, owner, asset, assetType, bidder, amount, counter_amount, counter_bid_amount, deadline,
                 bidConvId=None):
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

    def _counter_offer_wrong_owner(self, id):
        req = {
            "bidConvId": id,
            "assetOwner": self.bidder,
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

    def send_create_async(self):
        # createTx
        raw_txn = self._create_bid()

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_async(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        return result

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

    def send_counter_offer_malicious(self, id):
        # createTx
        raw_txn = self._counter_offer(id)

        # sign Tx
        signed = sign(raw_txn, self.bidder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if result["ok"]:
                print "Send Malicious Counter Offer Succeed : ", result
            else:
                print "################### Malicious Counter Offer Failed As Expected"

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


def broadcast_async(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxAsync', {
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


def query_balance(address):
    req = {"address": address}
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    return resp["result"]
