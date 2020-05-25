import json
from rpc_call import *

def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]


def sign(raw_tx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": raw_tx, "address": address})
    return resp["result"]


def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    if "result" in resp:
        return resp["result"]
    else:
        return resp


def broadcast_sync(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]


def create_proposal(proposal_id, prop_type, proposer, desc, initial_funding):
    req = {
        "proposal_id": proposal_id,
        "description": desc,
        "proposer": proposer,
        "proposal_type": prop_type,
        "initial_funding": {
            "currency": "OLT",
            "value": convertBigInt(initial_funding),
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

def vote_proposal(proposal_id, opinion, node_url):
    req = {
        "proposal_id": proposal_id,
        "opinion": opinion,
        "gasPrice": {
            "currency": "OLT",
            "value": "1000000000",
        },
        "gas": 40000,
    }
    resp = rpc_call('tx.CreateVote', req, node_url)
    result = resp["result"]
    print resp
    return result["rawTx"], result['signature']['Signed'], result['signature']['Signer']

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

def query_proposal(proposal_id):
    req = {"proposal_id": proposal_id}
    resp = rpc_call('query.GetProposalByID', req)
    print json.dumps(resp, indent=4)
    return resp["result"]
