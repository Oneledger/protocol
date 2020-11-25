from sdk import *


def delegate(node, account, amount):
    newDelegation = NetWorkDelegate(account, amount, node + "/keystore/")
    newDelegation.send_network_Delegate()


def check_rewards(result, balance, pending):
    if balance != '':
        balance = str(balance) + '0' * 18
        if int(result['balance']) < int(balance):
            sys.exit(-1)
    if pending != None:
        if len(result['pending']) != len(pending):
            sys.exit(-1)
        for i, amt in enumerate(pending):
            if amt != result['pending'][i]['amount']:
                sys.exit(-1)

def check_total_rewards(result, expected_exclude_withdrawn):
    if result < expected_exclude_withdrawn:
        sys.exit(-1)

if __name__ == "__main__":
    # create validator account
    funder = addValidatorWalletAccounts(node_0)

    # create delegator account
    delegator = createAccount(node_0, 2500000, funder)
    update_keystore(node_0, node_4)

    # delegates some OLT and wait for rewards distribution
    delegation_amt = 2000000
    delegation_amt_long = str(delegation_amt) + '0' * 18
    delegate(node_0, delegator, delegation_amt_long)
    wait_for(4)

    # query and check balance
    res = query_rewards(delegator)
    check_rewards(res, '6', [])

    # overdraw MUST fail
    overdraw_amount = '100' + '0' * 18
    withdraw = WithdrawRewards(delegator, overdraw_amount, node_0 + "/keystore/")
    withdraw.send(exit_on_err=False, mode=TxSync)
    print bcolors.OKGREEN + "#### Overdraw rewards failed as expected" + bcolors.ENDC

    # initiate 2 withdrawals
    pending = []
    total = 0
    for i in range(2):
        amt = i + 2
        amt_long = str(amt) + '0' * 18
        withdraw = WithdrawRewards(delegator, amt_long, node_0 + "/keystore/")
        withdraw.send(exit_on_err=True, mode=TxSync)
        pending.append(amt_long)
        total += amt
        wait_for(2)

    # query account balance before mature
    balance_before = query_balance(delegator)

    # query and check pending withdrawal
    res = query_rewards(delegator)
    check_rewards(res, '0', pending)
    print bcolors.OKGREEN + "#### Successfully withdrawn delegator rewards" + bcolors.ENDC

    # query and check again after maturity
    wait_for(4)
    res1 = query_rewards(delegator)
    check_rewards(res1, '', [])
    print bcolors.OKGREEN + "#### Successfully matured delegator rewards" + bcolors.ENDC

    # fully undelegate
    newDelegation = NetWorkDelegate(delegator, delegation_amt_long, node_0 + "/keystore/")
    newDelegation.send_network_undelegate(delegation_amt_long)
    wait_for(4)

    # withdraw all balance
    res2 = query_rewards(delegator)
    withdraw = WithdrawRewards(delegator, res2['balance'], node_0 + "/keystore/")
    withdraw.send(exit_on_err=True, mode=TxSync)

    # test query ListDelegation when there is no delegation and no rewards balance, only pending rewards
    wait_for(2)
    query_result = query_delegation()
    expected_delegation = 0
    expected_pending_delegation = 0
    expected_pending_rewards = int(res2['balance'])
    check_query_delegation(query_result, 0, expected_delegation, expected_pending_delegation, False, expected_pending_rewards)
    print bcolors.OKGREEN + "#### Successfully tested query ListDelegation with pending rewards" + bcolors.ENDC
    print bcolors.OKGREEN + "#### Successfully withdrawn all rewards" + bcolors.ENDC

    # query and check account balance
    balance_after = query_balance(delegator)
    check_balance(balance_before, balance_after, total + delegation_amt)

    # below is to test total rewards query
    # create another delegator account
    funder1 = addValidatorWalletAccounts(node_1)
    delegator1 = createAccount(node_1, 8000000, funder1)
    update_keystore(node_1, node_4)

    # delegates some OLT and wait for rewards distribution
    delegation_amt = '5000000' + '0' * 18
    delegate(node_1, delegator1, delegation_amt_long)
    wait_for(4)

    # initiate 1 withdrawal
    amt1 = 3
    amt1_long = str(amt1) + '0' * 18
    withdraw1 = WithdrawRewards(delegator1, amt1_long, node_1 + "/keystore/")
    withdraw1.send(True)

    wait_for(7)

    # query total rewards and check
    res = query_rewards(delegator)
    res1 = query_rewards(delegator1)
    total = query_total_rewards()
    check_total_rewards(total['totalRewards'], int(res['balance']) * pow(10, 18) + int(res1['balance']) * pow(10, 18))
    print bcolors.OKGREEN + "#### Successfully tested query total rewards" + bcolors.ENDC

