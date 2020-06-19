import requests
import json

rpc_url = "http://127.0.0.1:26605/jsonrpc"
token_url = "http://127.0.0.1:26605/token"

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

    response = requests.request("POST", rpc_url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        print response.reason
        return ""

    resp = json.loads(response.text)
    return resp


def get_token(username, password):
    payload = {
        "clientId": "123456",
        "username": username,
        "password": password
    }

    response = requests.request("POST", token_url, data=json.dumps(payload), headers=headers)

    if response.status_code != 200:
        return ""

    resp = json.loads(response.text)
    return resp


# Test RPC Call
def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp


if __name__ == "__main__":

    print('rpcAuth Test Script running')

    print('************ Test Get Token API ************')

    print('Send Get Token Request [GOOD Credentials]')
    resp = get_token("username", "password")
    print('Get Token Response:')
    print(resp)

    print('Send Get Token Request [BAD Credentials]')
    bad_token = get_token("username", "wrongPassword")
    print('Get Token Response:')
    print(bad_token)

    print('\n')

    print('Call owner.ListAccountAddresses from RPC Services [GOOD Signature]')
    # Test RPC Calls with Authentication
    headers["Authorization"] = resp['token']
    print(addresses())

    print('Call owner.ListAccountAddresses from RPC Services [BAD Signature]')
    headers["Authorization"] = 'HEfrz3HitJLLJ9NHTbFtGTjFFhs8WHA28odLaWgUfKvfNF6SmXg47KX9T2vVNcYvzRrChPusCafqr8w6UgUqm5oNpqjnR6DdWKBzpShzy5zFhR6D7P5'
    print(addresses())

    print

    print('rpcAuth Test Script Finished')


