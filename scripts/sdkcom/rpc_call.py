import json
import os
import requests

from config import *
from common import sdkIPAddress

headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}

if oltest == "1":
    url_fullnode = "http://{}/jsonrpc".format(sdkIPAddress(fullnode_dev))
else:
    url_fullnode = "http://{}/jsonrpc".format(sdkIPAddress(fullnode_prod))

def rpc_call(method, params, url=url_fullnode):
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

def tx_by_hash(hash):
    params = {"hash": hash}

    response = requests.get(url_tmTx, params=params)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp

def convertBigInt(value):
    return str(value)
