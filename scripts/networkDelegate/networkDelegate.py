from sdk import *

if __name__ == "__main__":
    userAccount = addValidatorWalletAccounts(node_5)
    delegationAccount = addNewAccount(node_5)
    delegationAmount = "1000000"
    output = sendFunds(userAccount, delegationAccount, delegationAmount, "1234", node_5)
    print output
    newDelegation = NetWorkDelegate(userAccount, delegationAccount, delegationAmount)
    newDelegation.send_network_Delegate()
    print bcolors.OKGREEN + "#### Delagtion Succedded" + bcolors.ENDC
