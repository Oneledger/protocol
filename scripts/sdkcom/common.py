import subprocess
import time
from actions import *

def addValidatorWalletAccounts(node):
    args = ['olclient', 'show_node_id']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    time.sleep(1)
    pubKey = output[0].split(",")[0].split(":")[1].strip()
    f = open(os.path.join(node, "consensus", "config", "node_key.json"), "r")
    contents = json.loads(f.read())
    privKey = contents['priv_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()

    if "exists" in output[0]:
        args = ['olclient', 'list']
        process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        return output[2].split(" ")[1].strip()[3:]
    time.sleep(1)
    return output[1].split(":")[1].strip()[3:]

def createAccount(node, funds=0, funder="", pswd="1234"):
    args = ['olclient', 'account', 'add', "--password", pswd]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    newaccount = output[1].split(":")[1].strip()[3:]

    if funds > 0:
        sendFunds(funder, newaccount, str(funds), pswd, node)
        balance = query_balance(newaccount)
        if balance != funds:
            sys.exit(-1)
    return newaccount

def sendFunds(party, counterparty, amount, password, node):
    args = ['olclient', 'send', "--password", password, "--party", party, "--counterparty", counterparty, "--amount",
            amount, "--fee", "0.001"]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    return output
