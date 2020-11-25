import sys, time, subprocess

from common import *
from constant import *
from rpc_call import *

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

def createAccount(node, funds=0, funder="", pswd="1234"):
    args = ['olclient', 'account', 'add', "--password", pswd]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    newaccount = output[1].split(":")[1].strip()[3:]

    if funds > 0:
        sendFunds(funder, newaccount, str(funds), pswd, node)
        balance = query_balance(newaccount)
        if balance != funds:
            sys.exit(-1)
    return newaccount

def height():
    resp = rpc_call('query.ListValidators', {})
    return resp["result"]["height"]

def wait_for(blocks):
    resp = rpc_call('query.ListValidators', {})
    hstart = resp["result"]["height"]
    hcur = hstart
    while hcur - hstart < blocks:
        time.sleep(0.5)
        resp = rpc_call('query.ListValidators', {})
        hcur = resp["result"]["height"]

def wait_until(method, req, predict):
    resp = rpc_call(method, req)
    result = resp["result"]
    while not predict(result):
        time.sleep(1.0)
        resp = rpc_call(method, req)
        result = resp["result"]
    return result
