import requests
import json

url = "http://127.0.0.1:26602/jsonrpc"
headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}

def create_domain(name, owner_hex, price):
    payload = {
        "method": "tx.ONS_CreateRawCreate",
        "params": {
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
        },
        "id": 123,
        "jsonrpc": "2.0"
    }

    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp["result"]["rawTx"]


def sell_domain(name, owner_hex, price):
    payload = {
        "method": "tx.ONS_CreateRawSale",
        "params": {
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
        },
        "id": 123,
        "jsonrpc": "2.0"
    }

    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp["result"]["rawTx"]

def cancel_sell_domain(name, owner_hex, price):
    payload = {
        "method": "tx.ONS_CreateRawSale",
        "params": {
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
        },
        "id": 123,
        "jsonrpc": "2.0"
    }

    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp["result"]["rawTx"]

def addresses():
    payload = {
        "method": "owner.ListAccountAddresses",
        "params": {},
        "id": 123,
        "jsonrpc": "2.0"
    }


    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""


    resp = json.loads(response.text)
    return resp["result"]["addresses"]


def sign(rawTx, address):
    payload = {
        "method": "owner.SignWithAddress",
        "params": {
            "rawTx": rawTx,
            "address": address
            },
        "id": 123,
        "jsonrpc": "2.0"
    }


    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""


    resp = json.loads(response.text)

    return resp["result"]


def broadcast_commit(rawTx, signature, pub_key):
    payload = {
        "method": "broadcast.TxCommit",
        "params": {
            "rawTx": rawTx,
            "signature": signature,
            "publicKey": pub_key,
        },
        "id": 123,
        "jsonrpc": "2.0"
    }


    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""


    resp = json.loads(response.text)

    return resp["result"]



if __name__ == "__main__":
    addrs = addresses()

    print addrs[1]

    raw_txn = create_domain("alice.olt", addrs[1], "1.2345")
    print raw_txn

    signed = sign(raw_txn, addrs[1])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    raw_txn = sell_domain("alice.olt", addrs[1], "10.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[1])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print

    '''
    raw_txn = cancel_sell_domain("alice.olt", addrs[1], "10.2345")
    print raw_txn
    print

    signed = sign(raw_txn, addrs[1])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print
    '''

