import requests
import json
import time
import sys
import struct
import binascii

from sdk.actions import *

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
        "buyingprice": {
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


def update_domain(name, owner_hex, set_active, uri):
    req = {
        "owner": owner_hex,
        "account": owner_hex,
        "name": name,
        "active": set_active,
        "uri": uri,
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

    name = "alic.ol"
    create_price = (int("10023450")*10**14)
    print "create price:", create_price

    raw_txn = create_domain(name, addrs[0], create_price)
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
    print_domain(name)

    time.sleep(2)

    print bcolors.WARNING + "*** Update Domain With Inactive Flag ***" + bcolors.ENDC
    active_state = False
    uri = ""
    raw_txn = update_domain(name, addrs[0], active_state, uri)
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
    print_domain(name)

    time.sleep(2)

    print bcolors.WARNING + "*** Update Domain With Active Flag ***" + bcolors.ENDC
    active_state = True
    raw_txn = update_domain(name, addrs[0], active_state, uri)
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
    print_domain(name)

    time.sleep(2)

    print bcolors.WARNING + "*** Update Domain With Uri ***" + bcolors.ENDC
    active_state = True
    uri = "http://192.168.0.1"
    raw_txn = update_domain(name, addrs[0], active_state, uri)
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
    print_domain(name)

    time.sleep(2)

    print bcolors.WARNING + "*** Update Domain to Reset Uri ***" + bcolors.ENDC
    active_state = True
    uri = ""
    raw_txn = update_domain(name, addrs[0], active_state, uri)
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
    print_domain(name)