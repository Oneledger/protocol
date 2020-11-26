import json
import os
import os.path as path
import sys

import requests
sdkcom_p = path.abspath(path.join(path.dirname(__file__), "../.."))
sys.path.append(sdkcom_p)

from sdkcom import *

class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


url = "http://127.0.0.1:26602/jsonrpc"

devnet = os.path.join(os.environ['OLDATA'], "devnet")
if is_docker():
    devnet = get_volume_info()
node_0 = os.path.join(devnet, "0-Node")
node_1 = os.path.join(devnet, "1-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")
node_1_log = os.path.join(node_1, "olfullnode.log")

headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}


def rpc_call(method, params):
    payload = {
        "method": method,
        "params": params,
        "id": 123,
        "jsonrpc": "2.0"
    }

    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp


def converBigInt(value):
    return str(value)


class Byzantine:
    def __init__(self, reporter, malicious, proofmessage, password, blockheight, keypath):
        self.reporter = reporter
        self.malicious = malicious
        self.proofmessage = proofmessage
        self.password = password
        self.blockheight = blockheight
        self.keypath = keypath

    def create_allegation(self):
        req = {
            "address": self.reporter,
            "maliciousAddress": self.malicious,
            "blockHeight": self.blockheight,
            "ProofMsg": self.proofmessage,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('tx.Allegation', req)
        print resp
        return resp["result"]["rawTx"]

    def send_allegation(self):
        """

        :rtype: object
        """
        # create Tx
        raw_txn = self.create_allegation()

        # sign Tx
        signed = sign(raw_txn, self.reporter, self.keypath)

        # broadcast Tx
        result = broadcast_sync(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print result
            else:
                print bcolors.OKBLUE + "allegation added" + bcolors.ENDC
                return result["txHash"]


class Vote:
    def __init__(self, validator, reqID, choise, keypath):
        self.validator = validator
        self.reqId = reqID
        self.choise = choise
        self.keypath = keypath

    def create_vote(self):
        req = {
            "address": self.validator,
            "requestID": self.reqId,
            "choice": self.choise,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('tx.Vote', req)
        print resp
        return resp["result"]["rawTx"]

    def send_vote(self):
        """

        :rtype: object
        """
        # create Tx
        raw_txn = self.create_vote()

        # sign Tx
        signed = sign(raw_txn, self.validator, self.keypath)

        # broadcast Tx
        result = broadcast_sync(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print result
            else:
                print bcolors.OKGREEN + "Vote added" + bcolors.ENDC
                return result["txHash"]


class Release:
    def __init__(self, ownerAddress, keypath):
        self.ownerAddress = ownerAddress
        self.keypath = keypath

    def create_release(self):
        req = {
            "address": self.ownerAddress,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('tx.Release', req)
        print resp
        return resp["result"]["rawTx"]

    def send_release(self):
        """

        :rtype: object
        """
        # create Tx
        raw_txn = self.create_release()

        # sign Tx
        signed = sign(raw_txn, self.ownerAddress, self.keypath)

        # broadcast Tx
        result = broadcast_sync(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
        if "ok" in result:
            if not result["ok"]:
                print result
            else:
                print bcolors.OKGREEN + "Release TX sent" + bcolors.ENDC
                return result["txHash"]


def sign(raw_tx, address, keypath):
    if is_docker():
        keypath = './keystore/'
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
    return resp["result"]
