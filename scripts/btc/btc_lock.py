
import requests
import json
import time
import sys
import struct
import binascii

url = "http://127.0.0.1:26602/jsonrpc"
headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}

def rpc_call(method, params):
    payload = {
        "method": method,
        "params": params,
        "id": 1,
        "jsonrpc": "2.0"
    }
    print(json.dumps(payload))
    response = requests.request("POST", url, data=json.dumps(payload), headers=headers, timeout=200)
    print(response, response.text)
    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp

def converBigInt(value):
    return str(value)


def prepare_lock(txhash, index):
    req = {
        "hash": txhash,
        "index": index,
    }
    resp = rpc_call('btc.PrepareLock', req)

    return resp["result"]


def add_signature(txn, sign, address, tracker_name):

    req = {
        "txn": txn,
        "signature": sign,
        "address": address,
        "tracker_name": tracker_name,
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('btc.AddUserSignatureAndProcessLock', req)

    return resp["result"]


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]

def sign(rawTx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": rawTx, "address": address})
    return resp["result"]


def broadcast_commit(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    return resp["result"]

if __name__ == "__main__":
    result = prepare_lock("860a32ef84ed54df86d207112d1f8d3d5ad28751b25cc7e2107ef55cccbc7586", 1)

    txn = "AQAAAAGGdbzMXPV+EOLHXLJRh9JaPY0fLREH0obfVO2E7zIKhgEAAAAA/////wEsVi4AAAAAABSq5lHld6v+HZUYct6LSCMruHh/cwAAAAA="
    signature = "akcwRAIgNiUpfUIjKkQ8akN8qCSfpKtF6tYuRCwouEJVsG2mdtkCIFkFT8/tp2yOzok9HfNbMtOrB4FYsLdOc/aHH2IA8qdYASECxSWjS07kwg3F6iXysp3oPIhRpQKCnGhaCWUGahIaZhc="


    print()
    print()

    addrs = addresses()

    result = add_signature(txn, signature, addrs[0], result['tracker_name'])

    signed = sign(result['rawTx'], addrs[0])
    print "signed create domain tx:", signed
    print

    result = broadcast_commit(result['rawTx'], signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print