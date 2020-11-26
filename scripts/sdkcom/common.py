import os, sys, json, subprocess, re

devnull = open(os.devnull, 'wb')

def addValidatorWalletAccounts(node):
    args = ['olclient', 'show_node_id']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    pubKey = output[0].split(",")[0].split(":")[1].strip()
    f = open(os.path.join(node, "consensus", "config", "node_key.json"), "r")
    contents = json.loads(f.read())
    privKey = contents['priv_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    update_keystore(node)

    if "exists" in output[0]:
        args = ['olclient', 'list']
        args_in_use = args_wrapper(args, node)
        process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        return output[2].split(" ")[1].strip()[3:]
    return output[1].split(":")[1].strip()[3:]

def addValidatorAccounts(node):
    f = open(os.path.join(node, "consensus", "config", "priv_validator_key.json"), "r")
    contents = json.loads(f.read())
    pubKey = contents['pub_key']['value']
    privKey = contents['priv_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    update_keystore(node)

    if "exists" in output[0]:
        print "account already exists"
        return ""
    return output[1].split(":")[1].strip()[3:]

def nodeAccount(node):
    args = ['olclient', 'show_node_id']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    address = output[0].split(",")[1].split(":")[1].strip()
    return address

def sdkIPAddress(node):
    args = ['olfullnode', 'status']
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    sdkport = output[2].split(":")[1].strip().split(" ")
    if is_docker():
        ip_addr = '127.0.0.1'
    else:
        ip_addr = sdkport[2].strip()
    return ip_addr + ":" + sdkport[0].strip()


def sendFunds(party, counterparty, amount, password, node):
    args = ['olclient', 'send', "--password", password, "--party", party, "--counterparty", counterparty, "--amount",
            amount, "--fee", "0.001"]
    args_in_use = args_wrapper(args, node)
    process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    return output

def check_balance(before, after, expected_diff):
    diff = after - before
    # print diff
    # print expected_diff
    if diff != expected_diff:
        print "actual difference:"
        print after - before
        print "expected difference:"
        print expected_diff
        sys.exit(-1)


def is_docker():
    args = ['pgrep', 'docker']
    process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    if not output:
        return False
    args = ['docker', 'ps']
    process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    process.wait()
    error = process.stderr.readlines()
    if error:
        return False
    docker_output = process.stdout.readlines()
    for line in docker_output:
        if '-Node' in line:
            return True
    return False


def args_wrapper(args, cwd):
    node_name = cwd.split('/')[-1]
    result = [args, cwd]
    if not is_docker():
        return result
    execute_command = 'cd /opt/data/devnet &&'
    count = 0
    for arg in args:
        count = count + 1
        # remove --root flag since we are already in node folder
        if arg == '--root' or args[count-2] == '--root':
            continue
        execute_command += ' ' + arg
    new_args = ['docker', 'exec', node_name, 'bash', '-c', execute_command]
    return [new_args, '/']


def get_volume_info(container_name='0-Node'):
    args = ['docker', 'inspect', container_name]
    process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    process.wait()
    error = process.stderr.readlines()
    output = process.stdout.readlines()
    for line in output:
        volume_path = re.search('"Source": "(.*)\d-Node"', line)
        if volume_path:
            return volume_path.group(1)
    return None

# this is needed because sometimes we use keystore belongs to one node and sign the tx using another node
# and docker instance cannot get what's outside its own volume
def update_keystore(from_node):
    from_keystore = os.path.join(from_node, 'keystore')
    if not os.path.isdir(from_keystore):
        sys.exit(-1)
    from_keystore = os.path.join(from_keystore, '*')
    parent_folder = os.path.dirname(from_node)
    subdirs = [os.path.join(parent_folder, o) for o in os.listdir(parent_folder) if os.path.isdir(os.path.join(parent_folder,o))]
    for dir in subdirs:
        if dir != from_node:
            to_keystore = os.path.join(dir, 'keystore')
            args_copy = 'mkdir -p ' + to_keystore + ' && cp ' + from_keystore + ' ' + to_keystore
            process = subprocess.Popen(args_copy, stderr=subprocess.PIPE, shell=True)
            process.wait()
            err = process.stderr.readlines()
            if err:
                print err
                sys.exit(-1)
