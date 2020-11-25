from sdk import *


def delegate(node, account, amount):
    newDelegation = NetWorkDelegate(account, amount, node + "/keystore/")
    newDelegation.send_network_Delegate()

def check_rewards(result, balance, matured, pending):
    if balance != '':
        balance = str(balance) + '0' * 18
        if result['balance'] < balance:
            sys.exit(-1)
    if matured != '' and result['matured'] != matured:
        sys.exit(-1)
    if pending != None:
        if len(result['pending']) != len(pending):
            sys.exit(-1)
        for i, amt in enumerate(pending):
            if amt != result['pending'][i]['amount']:
                sys.exit(-1)


if __name__ == "__main__":
    # create validator account
    funder = addValidatorWalletAccounts(node_0)

    # create delegator account
    delegator = createAccount(node_0, 2500000, funder)
    update_keystore(node_0, node_4)

    # delegates some OLT and wait for rewards distribution
    delegation = '2000000'
    delegate(node_0, delegator, delegation + '0'*18)

    # wait for 10 OLT of rewards
    amount = waitfor_rewards(delegator, "10", "balance")

    # reinvest
    invest = ReinvestRewards(delegator, node_0 + "/keystore/")
    balance_before = int(query_rewards(delegator)['balance']) / 10**18
    invest.send(amount * 20 ** 18, exit_on_err=False, mode=TxCommit)
    print bcolors.OKGREEN + "#### Investing more than actual rewards failed as expected" + bcolors.ENDC
    invest.send(amount * 10 ** 18, exit_on_err=True, mode=TxCommit)
    
    # check rewards left
    balance_after = int(query_rewards(delegator)['balance']) / 10 ** 18
    print bcolors.OKGREEN + "#### Balance before: " + str(balance_before) + " Banlance after: " + str(balance_after) + bcolors.ENDC
    if balance_after >= balance_before:
        sys.exit(-1)

    # check delegation amount
    amount_actual = int(query_delegation([delegator])[0]['delegationStats']['active'])
    amount_expect = int(delegation) + int(amount)
    if amount_actual != amount_expect * 10 ** 18:
        print bcolors.FAIL + "#### Reinvested rewards failed" + bcolors.ENDC
        sys.exit(-1)
    print bcolors.OKGREEN + "#### Successfully reinvested rewards" + bcolors.ENDC
