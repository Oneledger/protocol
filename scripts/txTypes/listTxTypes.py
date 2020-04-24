import requests
import json

class bcolors:
    WARNING = '\033[93m'
    ENDC = '\033[0m'

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

print bcolors.WARNING + "*** List Tx Types ***" + bcolors.ENDC
resp = rpc_call('query.ListTxTypes', {})
print resp