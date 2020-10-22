import sys
import time

from rpc_call import *
from constant import *

def sign(raw_tx, address, keypath):
    resp = rpc_call('owner.SignWithSecureAddress',
                    {"rawTx": raw_tx, "address": address, "password": "1234", "keypath": keypath})
    return resp["result"]

def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    if "result" in resp:
        return resp["result"]
    else:
        return resp

def broadcast_sync(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]

def broadcast_async(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxAsync', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]

def broadcast(raw_tx, signature, pub_key, mode=TxCommit):
    if mode == TxCommit:
        return broadcast_commit(raw_tx, signature, pub_key)
    elif mode == TxSync:
        return broadcast_sync(raw_tx, signature, pub_key)
    elif mode == TxAsync:
        return broadcast_async(raw_tx, signature, pub_key)
    else:
        return broadcast_commit(raw_tx, signature, pub_key)

def query_balance(address):
    req = {
        "currency": "OLT",
        "address": address
    }
    resp = rpc_call('query.CurrencyBalance', req)

    if "result" in resp:
        result = resp["result"]["balance"]
    else:
        result = ""
    return int(float(result))

def wait_for(blocks, url=url_0):
    resp = rpc_call('query.ListValidators', {}, url)
    hstart = resp["result"]["height"]
    hcur = hstart
    while hcur - hstart < blocks:
        time.sleep(0.5)
        resp = rpc_call('query.ListValidators', {}, url)
        hcur = resp["result"]["height"]
