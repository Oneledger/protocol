import json
import os
import subprocess
import time
from rpc_call import *
import psutil
from sdk import *

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
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
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
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
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
    DETACHED_PROCESS = 0x00000008
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdin=None, stdout=None, stderr=None, close_fds=True)
    return True
    # process.wait()
    # output = process.stdout.read()
    # if u'Returned Successfully' in output:
    #     return True
    # return False


def ByzantineFault_Vote(node, request_id, address, choice, password):
    args = [
        'olclient', 'byzantine_fault', 'vote',
        '--address', address,
        '--requestID', str(request_id),
        '--choice', choice,
        '--password', password,
    ]
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    if u'Returned Successfully' in output:
        return True
    return False


def KillNode(node):
    args = [
        "pgrep", "-f", node,
    ]
    args_in_use = [args, node]
    if is_docker():
        args_in_use = args_wrapper(["pgrep", "olfullnode"], node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    pid = output.strip()
    if not pid.isdigit():
        return False
    if is_docker():
        args_in_use = args_wrapper(['kill', pid], node)
        process = subprocess.Popen(args_in_use[0], stdout=subprocess.PIPE, stderr=devnull)
        process.wait()
        output = process.stdout.read()
        if not output:
            return True
        return False
    pr = psutil.Process(int(pid))
    pr.terminate()
    result = pr.wait(timeout=2)
    if not result:
        return True
    return False


def StartNode(node, node_1_log):
    args = [
        'olfullnode', 'node', '--root', node, '>>', node_1_log, "2>&1 &"
    ]
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.read()
    pid = output.strip()
    if not pid.isdigit():
        return False
    if is_docker():
        args_in_use = args_wrapper(['kill', pid], node)
        process = subprocess.Popen(args_in_use[0], stdout=subprocess.PIPE, stderr=devnull)
        process.wait()
        output = process.stdout.read()
        if not output:
            return True
        return False
    pr = psutil.Process(int(pid))
    pr.terminate()
    result = pr.wait(timeout=2)
    if not result:
        return True
    return False


def addValidatorWalletAccounts(node):
    args = ['olclient', 'show_node_id']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.readlines()
    time.sleep(1)
    pubKey = output[0].split(",")[0].split(":")[1].strip()
    f = open(os.path.join(node, "consensus", "config", "node_key.json"), "r")
    contents = json.loads(f.read())
    privKey = contents['priv_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
    process.wait()
    output = process.stdout.readlines()

    if "exists" in output[0]:
        args = ['olclient', 'list']
        args_in_use = args_wrapper(args, node)
        process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE, stderr=devnull)
        process.wait()
        output = process.stdout.readlines()
        return output[2].split(" ")[1].strip()[3:]
    time.sleep(1)
    return output[1].split(":")[1].strip()[3:]


def addOwnerAccount(node):
    f = open(os.path.join(node, "consensus", "config", "priv_validator_key.json"), "r")
    contents = json.loads(f.read())
    privKey = contents['priv_key']['value']
    pubKey = contents['pub_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()

    if "exists" in output[0]:
        args = ['olclient', 'list']
        args_in_use = args_wrapper(args, node)
        process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        return output[2].split(" ")[1].strip()[3:]
    time.sleep(1)
    return output[1].split(":")[1].strip()[3:]


def addNewAccount(node):
    args = ['olclient', 'account', 'add', "--password", '1234']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    return output[1].split(":")[1].strip()[3:]


def sendFunds(party, counterparty, amount, password, node):
    args = ['olclient', 'send', "--password", password, "--party", party, "--counterparty", counterparty, "--amount",
            amount, "--fee", "0.001"]
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    return output
