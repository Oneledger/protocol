import os, json, subprocess


def addValidatorWalletAccounts(node):
    args = ['olclient', 'show_node_id']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
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
    return output[1].split(":")[1].strip()[3:]

def nodeAccount(node):
    args = ['olclient', 'show_node_id']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    address = output[0].split(",")[1].split(":")[1].strip()
    return address

def sdkIPAddress(node):
    args = ['olfullnode', 'status']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    sdkport = output[2].split(":")[1].strip().split(" ")
    ip_addr = sdkport[2].strip() + ":" + sdkport[0].strip()
    return ip_addr


def sendFunds(party, counterparty, amount, password, node):
    args = ['olclient', 'send', "--password", password, "--party", party, "--counterparty", counterparty, "--amount",
            amount, "--fee", "0.001"]
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    return output
