
import requests
import json
import sys


from ecdsa import VerifyingKey, SECP256k1

url = "http://127.0.0.1:26602/jsonrpc"
headers = {
    "Content-Type": "application/json",
    "Accept": "application/json",
}

def rpc_call(method, params):
    payload = {
        "method": method,
        "params": params,
        "id": 1,
        "jsonrpc": "2.0"
    }
    print(json.dumps(payload))
    response = requests.request("POST", url, data=json.dumps(payload), headers=headers, timeout=200)
    print(response, response.text)
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

    return resp["result"]


def add_signature(txn, sign, address, tracker_name):

    req = {
        "txn": txn,
        "signature": sign,
        "address": address,
        "tracker_name": tracker_name,
        "gasprice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 400000,
    }
    resp = rpc_call('btc.AddUserSignatureAndProcessLock', req)

    return resp["result"]


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]

def sign(rawTx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": rawTx, "address": address})
    return resp["result"]


def broadcast_commit(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    return resp["result"]

if __name__ == "__main__":
    result = prepare_lock("1b9c234aa1add73f96cc2756de8adef9a0e78a113f660fccdfe9dc3dd6beb8aa", 1)


    print(result)

    sys.exit(1)

    txn = "01000000018675bccc5cf57e10e2c75cb25187d25a3d8d1f2d1107d286df54ed84ef320a860100000000ffffffff016c752e000000000014aae651e577abfe1d951872de8b48232bb8787f7300000000".decode('hex').encode('base64')
    signature = "47304402200ed0cb4f0b29a069e1ef730efd7920cd59c1cfc6f49cbe28c5512d1298d33aa4022016f96575652b19c356188911ad8731975357f7f8fbe60623bbae93864257b5b8012102c525a34b4ee4c20dc5ea25f2b29de83c8851a502829c685a0965066a121a6617".decode('hex').encode('base64')

#
    # comp_str = 'ef817664f54410936e81dce0b93996d7bda2e4c16747aeadc9cf6278c4cd427aae012a1a0a8e'
    # vk = VerifyingKey.from_string(bytearray.fromhex(comp_str), curve=SECP256k1)
    # print(vk.to_string("uncompressed").hex())

    print()
    print()

    addrs = addresses()

    result = add_signature(txn, signature, addrs[0], result['tracker_name'])

    signed = sign(result['rawTx'], addrs[0])
    print "signed create domain tx:", signed
    print

    result = broadcast_commit(result['rawTx'], signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print