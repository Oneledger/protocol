import requests
import json
import time
import sys

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

def create_domain(name, owner_hex, price):
    resp = rpc_call('tx.ONS_CreateRawCreate', {
        "name": name,
        "owner": owner_hex,
        "account": owner_hex,
        "price": {
            "currency": "OLT",
            "value": price,
        },
        "fee": {
            "currency": "OLT",
            "value": "0",
        },
        "gas": 0,
    })
    return resp["result"]["rawTx"]



def send_domain(name, frm, price):
    resp = rpc_call('tx.ONS_CreateRawSend', {
        "name": name,
        "from": frm,
        "amount": {
            "currency": "OLT",
            "value": price,
        },
        "fee": {
            "currency": "OLT",
            "value": "0",
        },
        "gas": 0,
    })
    return resp["result"]["rawTx"]


def sell_domain(name, owner_hex, price):
    resp = rpc_call('tx.ONS_CreateRawSale', {
        "name": name,
        "owner": owner_hex,
        "price": {
            "currency": "OLT",
            "value": price,
        },
        "cancel_sale": False,
        "fee": {
            "currency": "OLT",
            "value": "0",
        },
        "gas": 0,
    })
    return resp["result"]["rawTx"]

def cancel_sell_domain(name, owner_hex, price):
    resp = rpc_call('tx.ONS_CreateRawSale', {
        "name": name,
        "owner": owner_hex,
        "price": {
            "currency": "OLT",
            "value": price,
        },
        "cancel_sale": True,
        "fee": {
            "currency": "OLT",
            "value": "0",
        },
        "gas": 0,
    })
    return resp["result"]["rawTx"]

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

    raw_txn = create_domain("alice1.olt", addrs[0], "100.2345")
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = send_domain("alice1.olt", addrs[0], "10")
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print
    time.sleep(2)
    raw_txn = sell_domain("alice1.olt", addrs[0], "10.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print
    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = send_domain("alice1.olt", addrs[0], "100")
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print
    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = send_domain("alice1.olt", addrs[0], "100")
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print
    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = cancel_sell_domain("alice1.olt", addrs[0], "10.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print
    if result["ok"] != True:
        sys.exit(-1)


