import requests
import json

url_0 = "http://127.0.0.1:26602/jsonrpc"
url_1 = "http://127.0.0.1:26605/jsonrpc"
url_2 = "http://127.0.0.1:26608/jsonrpc"
url_3 = "http://127.0.0.1:26611/jsonrpc"
url_4 = "http://127.0.0.1:26614/jsonrpc"

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


def convertBigInt(value):
    return str(value)
