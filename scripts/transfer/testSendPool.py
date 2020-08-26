import json
import os
import re
import subprocess
import sys
import time

success = "Returned Successfully"
sendAmount = "10000"

url = "http://127.0.0.1:26602/jsonrpc"
devnet = os.path.join(os.environ['OLDATA'], "devnet")
node_0 = os.path.join(devnet, "0-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


def addValidatorAccounts(numofValidators):
    validatorAcounts = []
    for i in range(numofValidators):
        args = ['olclient', 'show_node_id']
        node = str(i) + "-Node"
        nodedir = os.path.join(devnet, node)
        process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        pubKey = output[0].split(",")[0].split(":")[1].strip()
        f = open(os.path.join(nodedir, "consensus", "config", "node_key.json"), "r")
        contents = json.loads(f.read())
        privKey = contents['priv_key']['value']
        args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
        process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        if "exists" in output[0]:
            args = ['olclient', 'list']
            process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
            process.wait()
            output = process.stdout.readlines()
            validatorAcounts.append(output[2].split(" ")[1].strip()[3:])
        validatorAcounts.append(output[1].split(":")[1].strip()[3:])
    return validatorAcounts


if __name__ == "__main__":
    # Creating accounts for with funds
    validatorAccounts = addValidatorAccounts(1)
    # send some funds to pool through olclient
    args = ['olclient', 'sendpool', '--amount', sendAmount, '--party', validatorAccounts[0],
            '--poolname', 'RewardsPool', '--fee', '0.0001', '--password', '1234']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.read()
    if not success in output:
        print "Send to pool was not successful"
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Success for send to pool" + bcolors.ENDC
    time.sleep(1)
    # Check balance
    args = ['olclient', 'balance', '--poolname', 'RewardsPool']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.read()
    m = re.search(r'(?:[, ])Balance: (\d+)', output)
    if m.group(1) != sendAmount:
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Success for reading balance" + bcolors.ENDC

    # Trying to Withdraw now with new address

print bcolors.OKGREEN + "#### Send Pool Test succeed" + bcolors.ENDC
