from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_3)
    delegationAmount = "1500000"
    undelegate_amount_1 = "100000"
    undelegate_amount_2 = "200000"
    undelegate_amount_3 = "400000"
    malicious_undelegate_amount = "2000000"

    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_3 + "/keystore/")
    newDelegation.send_network_Delegate()

    balance_before = query_balance(delegationAccount)

    time.sleep(2)
    newDelegation.send_network_undelegate(undelegate_amount_1)
    newDelegation.send_network_undelegate(undelegate_amount_2)
    time.sleep(5)
    newDelegation.send_network_undelegate(undelegate_amount_3)
    time.sleep(2)
    result = newDelegation.query_undelegate()
    check_query_undelegated(result, 3)
    total_result = query_total(0)
    expected_total = int(delegationAmount) - int(undelegate_amount_1) - int(undelegate_amount_2)
    check_query_total(total_result, str(expected_total) + '0' * 18)
    total_result_only_active = query_total(1)
    expected_active = expected_total - int(undelegate_amount_3)
    check_query_total(total_result_only_active, str(expected_active) + '0' * 18)

    balance_after = query_balance(delegationAccount)
    check_balance(balance_before, balance_after, int(undelegate_amount_1) + int(undelegate_amount_2))

    newDelegation.send_network_undelegate_shoud_fail(malicious_undelegate_amount)
    print bcolors.OKGREEN + "#### Undelegation Succeded" + bcolors.ENDC
