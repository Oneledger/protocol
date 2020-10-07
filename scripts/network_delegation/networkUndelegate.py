import time

from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_3)
    delegationAmount = "1000000"
    undelegate_amount_1 = "100000"
    undelegate_amount_2 = "200000"
    undelegate_amount_3 = "300000"
    malicious_undelegate_amount = "2000000"

    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_3 + "/keystore/")
    newDelegation.send_network_Delegate()

    time.sleep(2)
    newDelegation.send_network_undelegate(undelegate_amount_1)
    newDelegation.send_network_undelegate(undelegate_amount_2)
    time.sleep(5)
    newDelegation.send_network_undelegate(undelegate_amount_3)
    time.sleep(2)
    result = newDelegation.query_undelegate()
    check_query_undelegated(result, 3, "0")
    total_result = query_total()
    check_query_total(total_result, "1000000000000000000000000")

    newDelegation.send_network_undelegate_shoud_fail(malicious_undelegate_amount)
    print bcolors.OKGREEN + "#### Undelegation Succeded" + bcolors.ENDC




