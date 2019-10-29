
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
    print(json.dumps(payload))
    response = requests.request("POST", url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp

def converBigInt(value):
    return str(value)


def prepare_lock(txhash, index):
    req = {
        "hash": txhash,
        "index": index,
    }
    resp = rpc_call('btc.PrepareLock', req)
    print resp
    return resp["result"]


if __name__ == "__main__":
    prepare_lock("860a32ef84ed54df86d207112d1f8d3d5ad28751b25cc7e2107ef55cccbc7586", 1)