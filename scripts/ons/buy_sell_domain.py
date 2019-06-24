

import sys
import requests
import json
import time

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

def buy_domain(name, buyer, price):
    resp = rpc_call('tx.ONS_CreateRawBuy', {
        "name": name,
        "buyer": buyer,
        "account": buyer,
        "offering": {
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

def send(frm, to, amt):
    resp = rpc_call('tx.CreateRawSend', {
        "from": frm,
        "to": to,
        "amount": {
            "currency": "OLT",
            "value": amt,
        },
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

def new_account(name):
    resp = rpc_call('owner.GenerateNewAccount', {'name': name})
    return resp['result']


def sign(rawTx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": rawTx,"address": address})
    return resp["result"]


def broadcast_commit(rawTx, signature='', pub_key=''):
    params = {
        "rawTx": rawTx,
    }
    if signature != '':
        params["signature"] = signature
    if pub_key != '':
        params["publicKey"] = pub_key

    resp = rpc_call('broadcast.TxCommit', params)
    return resp["result"]

def broadcast_sync(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]


def get_domain_on_sale():
    resp = rpc_call('query.ONS_GetDomainOnSale', {'onSale': True})
    return resp


if __name__ == "__main__":
    # Create New account
    result = new_account('charlie')
    print result


    addrs = addresses()

    print addrs

    raw_txn = create_domain("bob2.olt", addrs[0], "100.2345")
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print

    raw_txn = send_domain("bob2.olt", addrs[0], "100")
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
    raw_txn = sell_domain("bob2.olt", addrs[0], "10.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print


    raw_txn = send(addrs[0], addrs[2], "20")
    print result
    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print

    print bcolors.WARNING + "*** Buying domain ***" + bcolors.ENDC

    raw_txn = buy_domain('bob2.olt', addrs[2], '20.0')
    signed = sign(raw_txn, addrs[2])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print


    print bcolors.WARNING + "*** Putting Domain on sale ***" + bcolors.ENDC
    raw_txn = sell_domain("bob2.olt", addrs[2], "99.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[2])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print
    if result["ok"] != True:
        sys.exit(-1)


    print bcolors.WARNING + "*** Get Domains on sale ***" + bcolors.ENDC
    resp = get_domain_on_sale()
    print resp
