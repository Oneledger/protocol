import time

from sdk import *

def delegate(node, account, amount):
    newDelegation = NetWorkDelegate(account, str(amount), node + "/keystore/")
    newDelegation.send_network_Delegate()

if __name__ == "__main__":
    # create validator account
    funder = addValidatorWalletAccounts(node_0)

    # create delegator account
    delegator = createAccount(node_0, 2000000, funder)

    # delegates some OLT
    delegate(node_0, delegator, 1000000)

    # initiate 2 withdrawal
    for i in range(2):
        withdraw = WithdrawRewards(delegator, i+1, node_0 + "/keystore/")
        withdraw.send()
        wait_for(1)

    # query should return 2 withdrawals
    res = query_rewards(delegator)

    print bcolors.OKGREEN + "#### Withdraw delegator rewards succeed" + bcolors.ENDC
