import subprocess
import time

from actions import *

stake_error = "stake address does not match"
success = "Returned Successfully"


def WithdrawRewards(Walletaddress, amount, secs=1):
    # fund the proposal
    withdraw = Withdraw(Walletaddress, amount)
    withdraw.send_withdraw()
    time.sleep(secs)


def addValidatorAccounts(numofValidators):
    validatorAcounts = []
    for i in range(numofValidators):
        args = ['olclient', 'show_node_id']
        node = str(i) + "-Node"
        nodedir = os.path.join(devnet, node)
        args_in_use = args_wrapper(args, nodedir)
        process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        pubKey = output[0].split(",")[0].split(":")[1].strip()
        f = open(os.path.join(nodedir, "consensus", "config", "node_key.json"), "r")
        contents = json.loads(f.read())
        privKey = contents['priv_key']['value']
        args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
        args_in_use = args_wrapper(args, nodedir)
        process = subprocess.Popen(args_in_use[0], cwd=args_in_use[1], stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        validatorAcounts.append(output[1].split(":")[1].strip()[3:])
    return validatorAcounts
