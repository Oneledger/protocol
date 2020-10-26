import time

from rpc_call import *


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


def sign(raw_tx, address, keypath):
    resp = rpc_call('owner.SignWithSecureAddress',
                    {"rawTx": raw_tx, "address": address, "password": "1234", "keypath": keypath})
    print resp
    return resp["result"]


def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
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
