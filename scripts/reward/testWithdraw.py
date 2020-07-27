from time import sleep

from sdk.common import *

tolerance = 1


def testrewardswithdraw(validatorAccounts, result):
    # query balance before
    for i in range(len(validatorAccounts)):
        balance_before = query_balance(validatorAccounts[i])
        node = str(i) + "-Node"
        nodedir = os.path.join(devnet, node)
        # query balance after
        withdrawAmount = 12000
        args = ['olclient', 'rewards', 'withdraw', '--address', validatorAccounts[i], '--amount', str(withdrawAmount),
                '--password', '1234']
        process = subprocess.Popen(args, cwd=nodedir)
        process.wait()
        time.sleep(1)
        balance_after = query_balance(validatorAccounts[i])
        if balance_after - balance_before < withdrawAmount - tolerance or balance_after - balance_before > withdrawAmount + tolerance:
            print "Withdraw amount does not match | Diff :" + str(
                balance_after - balance_before) + " | Current Tolerance " \
                                                  "is :" + str(
                tolerance)
            sys.exit(-1)
        del balance_before
        del balance_after
        print "Withdrawn :" + str(withdrawAmount) + "| for Validator :", str(validatorAccounts[i])


def addValidatorAccounts(numofValidators):
    validatorAcounts = []
    for i in range(numofValidators):
        args = ['olclient', 'show_node_id']
        node = str(i) + "-Node"
        nodedir = os.path.join(devnet, node)
        process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        sleep(1)
        pubKey = output[0].split(",")[0].split(":")[1].strip()
        f = open(os.path.join(nodedir, "consensus", "config", "node_key.json"), "r")
        contents = json.loads(f.read())
        privKey = contents['priv_key']['value']
        args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
        process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        sleep(1)
        validatorAcounts.append(output[1].split(":")[1].strip()[3:])
    return validatorAcounts


def addNewAccount():
    args = ['olclient', 'account', 'add', "--password", '1234']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    sleep(1)
    newaccount = output[1].split(":")[1].strip()[3:]
    balance = query_balance(newaccount)
    if balance != 0:
        sys.exit(-1)
    return newaccount


if __name__ == "__main__":
    testrewardswithdraw([addNewAccount()])
    sys.exit(-1)
    validatorAccounts = addValidatorAccounts(4)
    # send some funds to pool through olclient
    args = ['olclient', 'sendpool', '--root', node_0, '--amount', '1000000', '--party', validatorAccounts[0],
            '--poolName',
            'RewardsPool', '--fee', '0.0001']
    process = subprocess.Popen(args)
    process.wait()

    testrewardswithdraw(validatorAccounts)

print bcolors.OKGREEN + "#### Withdraw block rewards succeed" + bcolors.ENDC
