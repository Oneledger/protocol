import json
import os
import os.path as path
import sys

import requests
sdkcom_p = path.abspath(path.join(path.dirname(__file__), "../.."))
sys.path.append(sdkcom_p)

from sdkcom import *

headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}


def rpc_call(method, params, url=url_0):
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
