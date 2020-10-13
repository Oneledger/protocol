import json
import os

import requests

url = "http://127.0.0.1:26602/jsonrpc"

devnet = os.path.join(os.environ['OLDATA'], "devnet")
node_0 = os.path.join(devnet, "0-Node")
node_1 = os.path.join(devnet, "1-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")

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
    def __init__(self, reporter, malicious, proofmessage, password, blockheight):
        self.reporter = reporter
        self.malicious = malicious
        self.proofmessage = proofmessage
        self.password = password
        self.blockheight = blockheight

    def _create_allegation(self):
        req = {
            "address": self.reporter,
            "maliciousAddress": self.malicious,
            "blockHeight": self.blockheight,
            "ProofMsg": self.proofmessage
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('tx.Allegation', req)
        return resp["result"]["rawTx"]
