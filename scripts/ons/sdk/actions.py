from rpc_call import rpc_call, converBigInt
import json

def create_domain(name, owner_hex, price):
    req = {
        "name": name,
        "owner": owner_hex,
        "account": owner_hex,
        "buyingPrice": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.ONS_CreateRawCreate', req)
    return resp["result"]["rawTx"]


def send_domain(name, frm, price):
    resp = rpc_call('tx.ONS_CreateRawSend', {
        "name": name,
        "from": frm,
        "amount": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "gasPrice": {
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
        "cancelSale": False,
        "gasPrice": {
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
        "cancelSale": True,
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]


def buy_domain(name, buyer, price):
    resp = rpc_call('tx.ONS_CreateRawBuy', {
        "name": name,
        "buyer": buyer,
        "account": buyer,
        "offering": {
            "currency": "OLT",
            "value": converBigInt(price),
        },
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]


def send(frm, to, amt):
    resp = rpc_call('tx.CreateRawSend', {
        "from": frm,
        "to": to,
        "amount": {
            "currency": "OLT",
            "value": converBigInt(amt),
        },
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    })
    return resp["result"]["rawTx"]


def get_domain_on_sale():
    resp = rpc_call('query.ONS_GetDomainOnSale', {'onSale': True})
    return resp


def new_account(name):
    resp = rpc_call('owner.GenerateNewAccount', {'name': name})
    return resp['result']


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


def get_domains(owner_hex, onsale):
    req = {
        "owner": owner_hex,
        "onSale": onsale,
    }

    resp = rpc_call('query.ONS_GetDomainByOwner', req)
    return resp["result"]["domains"]

def get_parent_domains(owner_hex, onsale):
    req = {
        "owner": owner_hex,
        "onSale": onsale,
    }

    resp = rpc_call('query.ONS_GetParentDomainByOwner', req)
    return resp["result"]["domains"]

def get_sub_domains(name):
    req = {
        "name": name,
    }

    resp = rpc_call('query.ONS_GetSubDomainByName', req)
    return resp["result"]["domains"]

def print_all_domains(owner_addr):
    result = get_domains(owner_addr, False)
    print json.dumps(result, indent=4)
    print

def print_all_parent_domains(owner_addr):
    result = get_parent_domains(owner_addr, False)
    print json.dumps(result, indent=4)
    print

def print_all_sub_domains(owner_addr):
    result = get_sub_domains(owner_addr)
    print json.dumps(result, indent=4)
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

    resp = rpc_call('tx.ONS_CreateRawCreate', req)
    return resp["result"]["rawTx"]


def delete_sub_domain(name, owner_hex):
    req = {
        "owner": owner_hex,
        "name": name,
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 400000,
    }

    resp = rpc_call('tx.ONS_CreateRawDeleteSub', req)
    return resp["result"]["rawTx"]


def renew_domain(name, owner_hex, price):
    req = {
        "owner": owner_hex,
        "name": name,
        "account": owner_hex,
        "buyingprice": {
            "currency": "OLT",
            "value": converBigInt(price)
        },
        "gasprice": {
                    "currency": "OLT",
                    "value": "1000000000",
         },
         "gas": 400000,
    }
    resp = rpc_call('tx.ONS_CreateRawRenew', req)
    return resp["result"]["rawTx"]