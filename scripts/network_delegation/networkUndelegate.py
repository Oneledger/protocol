import time

from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_5)
    delegationAmount = "1000000"
    undelegate_amount_1 = "100000"
    undelegate_amount_2 = "200000"
    undelegate_amount_3 = "300000"
    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_5 + "/keystore/")

    for i in range(5):
        newDelegation.send_network_Delegate()
        print bcolors.OKGREEN + "#### Delegation Succeded" + bcolors.ENDC

    time.sleep(2)
    newDelegation.send_network_undelegate(undelegate_amount_1)
    newDelegation.send_network_undelegate(undelegate_amount_2)
    time.sleep(5)
    newDelegation.send_network_undelegate(undelegate_amount_3)
    time.sleep(2)
    result = newDelegation.query_undelegate()
    check_query_undelegated(result, 3, 0)
    total_result = query_total()
    check_query_total(total_result, 500000)

    print bcolors.OKGREEN + "#### Undelegation Test Succeded" + bcolors.ENDC


