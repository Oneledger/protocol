from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_3)
    print 'delegationAccount: ', delegationAccount
    query_full_balance(delegationAccount)
    delegationAmount = "1000000"
    undelegate_amount_1 = "100000"
    undelegate_amount_2 = "200000"
    undelegate_amount_3 = "300000"

    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_3 + "/keystore/")
    newDelegation.send_network_Delegate()

    time.sleep(2)
    newDelegation.send_network_undelegate(undelegate_amount_1)
    newDelegation.send_network_undelegate(undelegate_amount_2)
    time.sleep(5)
    newDelegation.send_network_undelegate(undelegate_amount_3)
    time.sleep(2)
    result = newDelegation.query_undelegate()
    check_query_undelegated(result, 3)
    total_result = query_total(0)
    check_query_total(total_result, "1000000000000000000000000")
    total_result_only_active = query_total(1)
    check_query_total(total_result_only_active, "400000000000000000000000")
    newDelegation.send_network_withdraw(undelegate_amount_1)
    time.sleep(1)
    print 'delegationAccount: ', delegationAccount
    query_full_balance(delegationAccount)

    print bcolors.OKGREEN + "#### Withdraw Succeded" + bcolors.ENDC