from rpc_call import rpc_call, convertBigInt
import json


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


def broadcast_sync(rawTx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": rawTx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]


def create_proposal(type, proposer, desc, initialFunding):
    req = {
        "description": desc,
        "proposer": proposer,
        "proposal_type": type,
        "initial_funding": {
            "currency": "OLT",
            "value": convertBigInt(initialFunding),
        },
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.CreateProposal', req)
    print resp
    return resp["result"]["rawTx"]


def query_proposals(prefix):
    req = {
        "prefix": prefix,
        "gasPrice":
        {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }

    resp = rpc_call('query.GetProposals', req)
    print json.dumps(resp, indent=4)
