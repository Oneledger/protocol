from sdk import *

if __name__ == "__main__":
    delegationAccount = addValidatorWalletAccounts(node_5)
    delegationAmount = "1000000"
    no_of_delagations = 5
    newDelegation = NetWorkDelegate(delegationAccount, delegationAmount, node_5 + "/keystore/")
    for i in range(no_of_delagations):
        newDelegation.send_network_Delegate()
    if int(newDelegation.query_delegation()['active'].split(" ")[0]) != int(delegationAmount) * no_of_delagations:
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Delegation Succeeded" + bcolors.ENDC
