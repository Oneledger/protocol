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


def create_domain(name, owner_hex, price):
    req = {
        "name": name,
        "owner": owner_hex,
        "account": owner_hex,
        "price": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.ONS_CreateRawCreate', req)
    print "Response",resp
    return resp["result"]["rawTx"]


def update_domain(name, owner_hex, set_active):
    req = {
        "owner": owner_hex,
        "account": owner_hex,
        "name": name,
        "active": set_active,
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000"
        },
        "gas": 400000,
    }
    resp = rpc_call("tx.ONS_CreateRawUpdate", req)
    return resp["result"]["rawTx"]


def send_domain(name, frm, price):
    resp = rpc_call('tx.ONS_CreateRawSend', {
        "name": name,
        "from": frm,
        "amount": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]


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


def broadcast_sync(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]


if __name__ == "__main__":
    addrs = addresses()
    print addrs

    name = "alice1.olt"
    create_price = (int("1002345") * 10 ** 14)
    print "create price:", create_price

    raw_txn = create_domain("alice1.olt", addrs[0], create_price)
    print "raw create domain tx:", raw_txn

    signed = sign(raw_txn, addrs[0])
    print "signed create domain tx:", signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if not result["ok"]:
        sys.exit(-1)

    active_state = False
    raw_txn = update_domain(name, addrs[0], active_state)
    print "rax update domain tx:", raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print "signed update TX :", signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)
