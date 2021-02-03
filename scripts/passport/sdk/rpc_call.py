import os
import requests
import json

oldata_path = os.getenv('OLDATA')
secure_accounts_path = oldata_path + "/devnet/0-Node/"

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


def convert_big_int(value):
    return str(value)
