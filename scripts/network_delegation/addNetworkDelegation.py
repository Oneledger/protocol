from sdk import *
if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_5)
    update_keystore(node_5, node_4)
    delegationAmount = "1000000"
    delegationAmount_long = delegationAmount + '0' * 18
    no_of_delagations = 5
    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount_long, node_5 + "/keystore/")
    for i in range(no_of_delagations):
        newDelegation.send_network_Delegate()

    query_result = query_delegation([delegationAccount])
    expected_delegation = int(delegationAmount_long) * no_of_delagations
    expected_pending_delegation = 0
    expected_pending_rewards = 0
    check_query_delegation(query_result, 0, expected_delegation, expected_pending_delegation, True, expected_pending_rewards)

    # test query ListDelegation with address that has no delegation at all
    random_address = '0lt3685fd5502ecba760ff5783eb233b8982ed03b6c'
    query_result = query_delegation([random_address])
    expected_delegation = 0
    expected_pending_delegation = 0
    expected_pending_rewards = 0
    check_query_delegation(query_result, 0, expected_delegation, expected_pending_delegation, False, expected_pending_rewards)


print bcolors.OKGREEN + "#### Delegation Succeeded" + bcolors.ENDC
