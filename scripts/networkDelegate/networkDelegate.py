from sdk import *

if __name__ == "__main__":
    # create proposal
    delegationAccount = addValidatorWalletAccounts(node_5)
    userAccount = addNewAccount(node_5)

    print bcolors.OKGREEN + "#### Delagtion Succedded" + bcolors.ENDC
    print ""
