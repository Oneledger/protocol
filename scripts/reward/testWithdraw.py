from os.path import dirname
from time import sleep

from sdk.common import *

tolerance = 10


def testrewardswithdraw(validatorAccounts):
    # query balance before
    for i in range(len(validatorAccounts)):
        print validatorAccounts[i]
        balance_before = query_balance(validatorAccounts[i])
        target = dirname(dirname(os.getcwd()))
        node = target + "/Codebase/Test/devnet/" + str(i) + "-Node"
        # query balance after
        withdrawAmount = 12000
        args = ['olclient', 'rewards', 'withdraw', '--address', validatorAccounts[i], '--amount', str(withdrawAmount),
                '--password', '1234']
        process = subprocess.Popen(args, cwd=node)
        process.wait()
        balance_after = query_balance(validatorAccounts[i])
        if balance_after - balance_before < withdrawAmount - tolerance or balance_after - balance_before > withdrawAmount + tolerance:
            print "Withdraw amount does not match"
            sys.exit(-1)
        del balance_before
        del balance_after
        print "Withdrawn :" + str(withdrawAmount) + "| for Validator :", str(validatorAccounts[i])


def addValidatorAccounts(numofValidators):
    validatorAcounts = []
    for i in range(numofValidators):
        args = ['olclient', 'show_node_id']
        target = dirname(dirname(os.getcwd()))
        node = str(i) + "-Node"
        dir = target + "/Codebase/Test/devnet/" + node
        process = subprocess.Popen(args, cwd=dir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        sleep(1)
        pubKey = output[0].split(",")[0].split(":")[1].strip()
        f = open(dir + "/consensus/config/node_key.json", "r")
        contents = json.loads(f.read())
        privKey = contents['priv_key']['value']
        args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
        process = subprocess.Popen(args, cwd=dir, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        sleep(1)
        validatorAcounts.append(output[1].split(":")[1].strip()[3:])
    return validatorAcounts


if __name__ == "__main__":
    validatorAccounts = addValidatorAccounts(4)
    # send some funds to pool through olclient
    args = ['olclient', 'sendpool', '--root', node_0, '--amount', '1000000', '--party', validatorAccounts[0],
            '--poolName',
            'RewardsPool', '--fee', '0.0001']
    process = subprocess.Popen(args)
    process.wait()

    # # test rewards distribution
    # validators = testRewardsDistribution()
    #
    # test rewards withdraw
    testrewardswithdraw(validatorAccounts)

print bcolors.OKGREEN + "#### Withdraw block rewards succeed" + bcolors.ENDC
