from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_3)
    delegationAmount = "1500000" + '0' * 18
    undelegate_amount_1 = "100000"
    undelegate_amount_2 = "200000"
    undelegate_amount_3 = "400000"
    undelegate_amount_4 = "800000"
    undelegate_amount_1_long = undelegate_amount_1 + '0' * 18
    undelegate_amount_2_long = undelegate_amount_2 + '0' * 18
    undelegate_amount_3_long = undelegate_amount_3 + '0' * 18
    undelegate_amount_4_long = undelegate_amount_4 + '0' * 18
    malicious_undelegate_amount = "2000000" + '0' * 18

    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_3 + "/keystore/")
    newDelegation.send_network_Delegate()

    balance_before = query_balance(delegationAccount)

    time.sleep(2)
    newDelegation.send_network_undelegate(undelegate_amount_1_long)
    newDelegation.send_network_undelegate(undelegate_amount_2_long)

    # test query ListDelegation
    query_result = query_delegation()
    expected_delegation = int(delegationAmount) - int(undelegate_amount_1_long) - int(undelegate_amount_2_long)
    expected_pending_delegation = int(undelegate_amount_1_long) + int(undelegate_amount_2_long)
    expected_pending_rewards = 0
    check_query_delegation(query_result, 0, expected_delegation, expected_pending_delegation, True, expected_pending_rewards)

    time.sleep(5)
    newDelegation.send_network_undelegate(undelegate_amount_3_long)
    time.sleep(2)
    result = newDelegation.query_undelegate()
    check_query_undelegated(result, 3)
    total_result = query_total(0)
    expected_total = int(delegationAmount) - int(undelegate_amount_1_long) - int(undelegate_amount_2_long)
    check_query_total(total_result, str(expected_total))
    total_result_only_active = query_total(1)
    expected_active = expected_total - int(undelegate_amount_3_long)
    check_query_total(total_result_only_active, str(expected_active))

    balance_after = query_balance(delegationAccount)
    check_balance(balance_before, balance_after, int(undelegate_amount_1) + int(undelegate_amount_2))

    newDelegation.send_network_undelegate_shoud_fail(malicious_undelegate_amount)

    # fully undelegate
    newDelegation.send_network_undelegate(undelegate_amount_4_long)
    time.sleep(8)

    # test query ListDelegation when there is no delegation
    query_result = query_delegation()
    expected_delegation = 0
    expected_pending_delegation = 0
    expected_pending_rewards = 0
    check_query_delegation(query_result, 0, expected_delegation, expected_pending_delegation, True, expected_pending_rewards)

    print bcolors.OKGREEN + "#### Undelegation Succeded" + bcolors.ENDC
