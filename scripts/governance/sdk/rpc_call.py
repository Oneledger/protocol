import json
import os

import requests

url_tmTx = "http://127.0.0.1:26600/tx"
url_0 = "http://127.0.0.1:26602/jsonrpc"
url_1 = "http://127.0.0.1:26605/jsonrpc"
url_2 = "http://127.0.0.1:26608/jsonrpc"
url_3 = "http://127.0.0.1:26611/jsonrpc"
url_4 = "http://127.0.0.1:26614/jsonrpc"
url_5 = "http://127.0.0.1:26617/jsonrpc"

devnet = os.path.join(os.environ['OLDATA'], "devnet")
node_0 = os.path.join(devnet, "0-Node")
node_1 = os.path.join(devnet, "1-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")
node_5 = os.path.join(devnet, "5-Node")
node_6 = os.path.join(devnet, "6-Node")
node_7 = os.path.join(devnet, "7-Node")
node_8 = os.path.join(devnet, "8-Node")
node_9 = os.path.join(devnet, "9-Node")

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
