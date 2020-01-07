"""
Create Sub Domain

1. Create Domain using existing script.
2. Create Sub Domain based on initial Parent domain.
3. Send Currency to the sub domain.

"""

import requests
import json
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


def converBigInt(value):
    return str(value)


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
        "gas": 400000,
    }
    resp = rpc_call('tx.ONS_CreateRawCreate', req)
    return resp["result"]["rawTx"]


def get_domains(owner_hex, onsale):
    req = {
        "owner": owner_hex,
        "onSale": onsale,
    }

    resp = rpc_call('query.ONS_GetDomainByOwner', req)
    return resp["result"]["rawTx"]


def print_all_domains(owner_addr):
    raw_tx = get_domains(owner_addr, False)
    signedTx = sign(raw_tx, owner_addr)
    result = broadcast_commit(raw_tx, signedTx['signature']['Signed'], signedTx['signature']['Signer'])
    print result
    print


def create_sub_domain(name, owner_hex, price, uri):
    req = {
        "owner": owner_hex,
        "account": owner_hex,
        "name": name,
        "buyingprice": {
            "currency": "OLT",
            "value": converBigInt(price)
        },
        "uri": uri,
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 400000,
    }

    resp = rpc_call('tx.ONS_CreateRawCreateSub', req)
    return resp["result"]["rawTx"]


if __name__ == "__main__":
    addrs = addresses()

    """
        ****** Create Initial Domain ******
    """

    # Prepare new domain for creation
    name = "alice2.ol"
    create_price = (int("10002345") * 10 ** 14)
    print "create price:", create_price

    # Get raw transaction for domain creation
    raw_txn = create_domain(name, addrs[0], create_price)
    print "raw create domain tx:", raw_txn

    # Sign raw Transaction
    signed = sign(raw_txn, addrs[0])
    print "signed create domain tx:", signed
    print

    # Broadcast transaction
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if not result["ok"]:
        sys.exit(-1)

    """
        ****** Create Sub domain based on initial domain above ******
    """
    print "---Creating Sub Domain---"

    # Prepare sub domain
    sub_name = "bob.alice2.ol"
    # Use same create price as above

    # Get raw transaction
    raw_txn = create_sub_domain(sub_name, addrs[0], create_price, "http://myuri.com")
    print "raw create sub domain transaction: ", raw_txn

    # Sign Transaction
    signed = sign(raw_txn, addrs[0])
    print "signed create sub domain transaction: ", signed
    print

    # Broadcast Transaction
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print
