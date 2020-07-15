import os
import requests
import json

url_tmTx = "http://127.0.0.1:26600/tx"
url_0 = "http://127.0.0.1:26602/jsonrpc"
url_1 = "http://127.0.0.1:26605/jsonrpc"
url_2 = "http://127.0.0.1:26608/jsonrpc"
url_3 = "http://127.0.0.1:26611/jsonrpc"
url_4 = "http://127.0.0.1:26614/jsonrpc"

devnet = os.path.join(os.environ['OLDATA'], "devnet")
node_0 = os.path.join(devnet, "0-Node")
node_1 = os.path.join(devnet, "1-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")

headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}

def rpc_call(method, params, url=url_4):
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
