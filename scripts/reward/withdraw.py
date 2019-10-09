import requests
import json
import time
import sys
import struct
import binascii

url = "http://127.0.0.1:26606/jsonrpc"
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


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]

def sign(rawTx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": rawTx,"address": address})
    return resp["result"]


def broadcast_commit(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    return resp["result"]

def broadcast_sync(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]

def withdraw(frm, to):
    resp = rpc_call('tx.WithdrawReward', {
        "from": frm,
        "to": to,
        "fee": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]

def balance(addr):
    resp = rpc_call('query.Balance', {
        "address": addr,
    })
    return resp["result"]["balance"]

if __name__ == "__main__":
    addrs = addresses()
    print addrs

    dest_addr = "0xdeadbeafdeadbeafdeadbeafdeadbeafdeadbeaf"
    raw_txn = withdraw(addrs[0], dest_addr)
    print "raw withdraw tx:", raw_txn

    signed = sign(raw_txn, addrs[0])
    print "signed withdraw tx:", signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result

    bal = balance(dest_addr)
    print "addr", dest_addr, "balance", bal
    if len(bal) > 0:
        print "Test withdraw succeed"
    print "###################"
    print
