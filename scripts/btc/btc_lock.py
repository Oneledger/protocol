
import requests
import json

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
        "btc_fees": 50000,
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
    result = prepare_lock("4c1f929f875e90908a54a06b49903155998154a64cb930d56625bd713a49e7d6", 1)


    print(result)

    # sys.exit(1)

    txn = "0100000001d6e7493a71bd2566d530b94ca6548199553190496ba0548a90905e879f921f4c0100000000ffffffff0140290f0000000000140678fa5af71cdcfbfa849cb439bae319ba27f0ff00000000".decode('hex').encode('base64')
    signature = "47304402207e22394f66239f342c2a7ac6d377d3a90885672f1554a987615c177c81b9fa5c02200df35360c64a784bb3b19ea7604b3e9fb4fc2ea81bfd8a8751b3e2b8d3828975012102f2ecadea41a08ebbb9aa1ee2b75a311e9096086d4a8f9ca2a3cb7a56fa742ef0".decode('hex').encode('base64')

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