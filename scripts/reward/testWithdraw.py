from time import sleep

from sdk.common import *

tolerance = 1
stake_error = "stake address does not match"
success = "Returned Successfully"


def testrewardswithdraw(validatorAccounts, result, error_message):
    # query balance before
    for i in range(len(validatorAccounts)):
        balance_before = query_balance(validatorAccounts[i])
        node = str(i) + "-Node"
        nodedir = os.path.join(devnet, node)
        # query balance after
        withdrawAmount = 12000
        args = ['olclient', 'rewards', 'withdraw', '--address', validatorAccounts[i], '--amount', str(withdrawAmount),
                '--password', '1234']
        process = subprocess.Popen(args, cwd=nodedir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(1)
        if not result:
            if error_message in output[1]:
                return
            print "Output : ", output
            sys.exit(-1)

        balance_after = query_balance(validatorAccounts[i])
        if balance_after - balance_before < withdrawAmount - tolerance or balance_after - balance_before > withdrawAmount + tolerance:
            print "Withdraw amount does not match | Diff :" + str(
                balance_after - balance_before) + " | Current Tolerance " \
                                                  "is :" + str(
                tolerance)
            sys.exit(-1)
        del balance_before
        del balance_after
        print "Withdrawn :" + str(withdrawAmount) + "| To address :", str(validatorAccounts[i])


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
    # Expected to fail with error message for stake address
    newAccount = addNewAccount()
    testrewardswithdraw([newAccount], False, stake_error)
    print bcolors.OKGREEN + "#### Should Fail for withdraw to new address" + bcolors.ENDC
    # Creating accounts for Validators
    validatorAccounts = addValidatorAccounts(4)

    # send some funds to pool through olclient
    args = ['olclient', 'sendpool', '--amount', '1000000', '--party', validatorAccounts[0],
            '--poolName',
            'RewardsPool', '--fee', '0.0001']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    if not success in output[1]:
        print "Send to pool was not successful"
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Success for send to pool" + bcolors.ENDC

    # Withdraw from all four validators
    testrewardswithdraw(validatorAccounts, True, "No Error")
    print bcolors.OKGREEN + "#### Success for withdraw to Staking address" + bcolors.ENDC
    # Unstaking from 0-Node
    args = ['olclient', 'delegation', 'unstake', '--address', validatorAccounts[0], '--amount', '1', '--password',
            '1234']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    if not success in output[1]:
        print "Unstake was not successful"
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Success for unstake" + bcolors.ENDC
    time.sleep(5)

    # Trying to Withdraw now with new address
    testrewardswithdraw([newAccount], True, "No Error")
    print bcolors.OKGREEN + "#### Success for withdraw to new address" + bcolors.ENDC

print bcolors.OKGREEN + "#### Withdraw block rewards succeed" + bcolors.ENDC
