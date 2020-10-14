import sys
import time

from sdk import *

def delegate(node, account, amount):
    newDelegation = NetWorkDelegate(account, amount, node + "/keystore/")
    newDelegation.send_network_Delegate()

def check_rewards(result, balance, matured, pending):
    if balance != '':
        balance = str(balance) + '0' * 18
        if result['balance'] < balance:
            sys.exit(-1)
    if matured != '' and result['matured'] != matured:
        sys.exit(-1)
    if pending != None:
        if len(result['pending']) != len(pending):
            sys.exit(-1)
        for i, amt in enumerate(pending):
            if amt != result['pending'][i]['amount']:
                sys.exit(-1)

if __name__ == "__main__":
    # create validator account
    funder = addValidatorWalletAccounts(node_0)

    # create delegator account
    delegator = createAccount(node_0, 2500000, funder)

    # delegates some OLT and wait for rewards distribution
    delegate(node_0, delegator, '2000000')
    wait_for(4)

    # query and check balance
    res = query_rewards(delegator)
    check_rewards(res, '6', '0', [])

    # initiate 2 withdrawals
    pending = []
    total = 0
    for i in range(2):
        amt = i+2
        withdraw = WithdrawRewards(delegator, amt, node_0 + "/keystore/")
        withdraw.send(True)
        pending.append(str(amt) + '0' * 18)
        total += amt
        wait_for(1)
    total = str(total) + '0' * 18

    # query and check pending withdrawal
    res = query_rewards(delegator)
    check_rewards(res, '0', '0', pending)
    print "#### Successfully withdrawn delegator rewards"

    # query and check again after maturity
    wait_for(4)
    res1 = query_rewards(delegator)
    check_rewards(res1, '', total, [])
    print "#### Successfully matured delegator rewards"

    # finalize more than matured
    finalize = FinalizeRewards(delegator, node_0 + "/keystore/")
    finalize.send_finalize(total*2, False)
    print bcolors.OKGREEN + "#### Overdraw finalize rewards failed as expected" + bcolors.ENDC

    # finalize all matured rewards
    finalize.send_finalize(total, True)
    # query and check
    wait_for(2)
    res2 = query_rewards(delegator)
    check_rewards(res2, '', '0', [])
    print bcolors.OKGREEN + "#### Successfully finalized delegator rewards" + bcolors.ENDC

    # withdraw all balance
    balance = int(res['balance']) / 1000000000000000000
    withdraw = WithdrawRewards(delegator, balance, node_0 + "/keystore/")
    withdraw.send(True)
    print "#### Successfully withdrawn all rewards"

    # finalize withdraw rewards again
    wait_for(6)
    finalize.send_finalize(str(balance) + '0' * 18, True)
    # query and check
    wait_for(3)
    res3 = query_rewards(delegator)
    check_rewards(res3, '', '0', [])
    print bcolors.OKGREEN + "#### Successfully finalized all rewards" + bcolors.ENDC