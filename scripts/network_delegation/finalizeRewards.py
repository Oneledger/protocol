
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


def check_balance(before, after, expected_diff):
    diff = after - before
    # print diff
    # print expected_diff
    if diff != expected_diff:
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
        amt = i + 2
        withdraw = WithdrawRewards(delegator, amt, node_0 + "/keystore/")
        withdraw.send(True)
        pending.append(str(amt) + '0' * 18)
        total += amt
        wait_for(1)
    total_str = str(total) + '0' * 18

    # query and check pending withdrawal
    res = query_rewards(delegator)
    check_rewards(res, '0', '0', pending)
    print "#### Successfully withdrawn delegator rewards"

    # query and check again after maturity
    wait_for(4)
    res1 = query_rewards(delegator)
    check_rewards(res1, '', total_str, [])
    print "#### Successfully matured delegator rewards"

    # finalize more than withdrawn
    finalize = FinalizeRewards(delegator, node_0 + "/keystore/")
    finalize.send_finalize(int(total_str) * 2, False)
    print bcolors.OKGREEN + "#### finalize rewards more than withdrawn failed as expected" + bcolors.ENDC

    # query balance
    balance_before = query_balance(delegator)
    # finalize withdrawn rewards
    finalize.send_finalize(total_str, True)
    # query and check
    wait_for(2)
    res2 = query_rewards(delegator)
    check_rewards(res2, '', '0', [])
    balance_after = query_balance(delegator)
    check_balance(balance_before, balance_after, total)
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
    balance_final = query_balance(delegator)
    check_balance(balance_after, balance_final, balance)
    print bcolors.OKGREEN + "#### Successfully finalized all rewards" + bcolors.ENDC
