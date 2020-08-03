import json
import os
import signal
import subprocess

import psutil

from .rpc_call import devnet, node_0


devnull = open(os.devnull, 'wb')


def GetNodeKey(i):
    node = str(i) + "-Node"
    nodedir = os.path.join(devnet, node)
    with open(os.path.join(nodedir, "consensus", "config", "node_key.json"), "r") as priv_key:
        content = json.loads(priv_key.read())
    return {
        "priv": content['priv_key']['value'],
    }


def GetNodeCreds(i):
    node = str(i) + "-Node"
    nodedir = os.path.join(devnet, node)
    with open(os.path.join(nodedir, "consensus", "config", "priv_validator_key.json"), "r") as priv_key:
        content = json.loads(priv_key.read())
    pubKey = content['pub_key']['value']
    privKey = content['priv_key']['value']
    return {
        "address": content['address'].lower(),
        "pub": pubKey,
        "priv": privKey,
    }


def Send(node, party, counterparty, amount, password, currency='OLT', fee=10):
    args = [
        'olclient', 'send',
        '--party', party,
        '--counterparty', counterparty,
        '--amount', str(amount),
        '--currency', currency,
        '--fee', str(fee),
        '--password', password,
    ]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    if u'Returned Successfully' in output:
        return True
    return False


def Account_Add(node, pubkey, privkey, password):
    args = [
        'olclient', 'account', 'add',
        '--pubkey', pubkey,
        '--privkey', privkey,
        '--password', password,
    ]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    if u'Successfully added' in output:
        return True
    if u'key file already exists' in output:
        return True
    return False


def ByzantineFault_Allegation(node, address, malicious_address, block_height, proof_msg, password):
    args = [
        'olclient', 'byzantine_fault', 'allegation',
        '--address', address,
        '--maliciousAddress', malicious_address,
        '--blockHeight', str(block_height),
        '--proofMsg', proof_msg,
        '--password', password,
    ]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    if u'Returned Successfully' in output:
        return True
    return False


def ByzantineFault_Vote(node, request_id, address, choice, password):
    args = [
        'olclient', 'byzantine_fault', 'vote',
        '--address', address,
        '--requestID', str(request_id),
        '--choice', choice,
        '--password', password,
    ]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    if u'Returned Successfully' in output:
        return True
    return False


def KillNode(node):
    args = [
        "pgrep", "-f", node,
    ]
    process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    pid = output.strip()
    if not pid.isdigit():
        return False

    pr = psutil.Process(int(pid))
    pr.terminate()
    result = pr.wait(timeout=2)
    if not result:
        return True
    return False
