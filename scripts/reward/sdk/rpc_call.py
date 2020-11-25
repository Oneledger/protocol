import json
import os
import os.path as path
import sys
sdkcom_p = path.abspath(path.join(path.dirname(__file__), "../.."))
sys.path.append(sdkcom_p)
from sdkcom import *

import requests

url = "http://127.0.0.1:26602/jsonrpc"
devnet = os.path.join(os.environ['OLDATA'], "devnet")
docker_path = get_volume_info()
if is_docker():
    devnet = docker_path
node_0 = os.path.join(devnet, "0-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")
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

def converBigInt(value):
    return str(value)
