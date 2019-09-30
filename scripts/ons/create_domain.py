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
        "fee": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.ONS_CreateRawCreate',req)
    return resp["result"]["rawTx"]



def send_domain(name, frm, price):
    resp = rpc_call('tx.ONS_CreateRawSend', {
        "name": name,
        "from": frm,
        "amount": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "fee": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]


def sell_domain(name, owner_hex, price):
    resp = rpc_call('tx.ONS_CreateRawSale', {
        "name": name,
        "owner": owner_hex,
        "price": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "cancel_sale": False,
        "fee": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]

def cancel_sell_domain(name, owner_hex, price):
    resp = rpc_call('tx.ONS_CreateRawSale', {
        "name": name,
        "owner": owner_hex,
        "price": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "cancel_sale": True,
        "fee": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]

def get_domain_on_sale():
    resp = rpc_call('query.ONS_GetDomainOnSale', {'onSale': True})
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



if __name__ == "__main__":
    addrs = addresses()
    print addrs

    create_price = (int("1002345")*10**14)
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

    sell_price = (int("105432")*10**14)
    raw_txn = sell_domain("alice1.olt", addrs[0], sell_price)
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

    raw_txn = send_domain("alice1.olt", addrs[0], (int("100")*10**18))
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

    raw_txn = send_domain("alice1.olt", addrs[0], (int("100")*10**18))
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

    print "Get Domain on Sale"
    resp = get_domain_on_sale()
    print resp

    print "#################" \
          "##"
    print

    raw_txn = cancel_sell_domain("alice1.olt", addrs[0], sell_price)
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

    print "Get Domain on Sale"
    resp = get_domain_on_sale()
    print resp

    print "#################" \
          "##"
    print
