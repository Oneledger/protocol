from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_5)
    delegationAmount = "1000000"
    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_5 + "/keystore/")
    for i in range(5):
        newDelegation.send_network_Delegate()
    print newDelegation.query_delegation()
    print bcolors.OKGREEN + "#### Delegation Succeeded" + bcolors.ENDC
