import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_30034"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_contributor = addr_list[1]
_funds_amount = (int("2") * 10 ** 9)
_withdraw_amount = (int("2") * 10 ** 9)

_prop = Proposal(_pid, "general", "proposal for funds withdrawing", _proposer, _initial_funding)
_encoded_pid = _prop.get_encoded_pid()

_wait = 6


def fund_proposal(pid, amount, funder):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(2)


def withdraw_fund(pid, contributor, amount, beneficiary):
    # fund the proposal
    fund_withdraw = ProposalFundsWithdraw(pid, contributor, amount, beneficiary)
    fund_withdraw.withdraw_fund(contributor)
    time.sleep(2)


def withdraw_fund_malicious(wrong_contributor, pid, contributor, amount, beneficiary):
    # fund the proposal
    fund_withdraw = ProposalFundsWithdraw(pid, contributor, amount, beneficiary)
    fund_withdraw.withdraw_fund(wrong_contributor)
    time.sleep(2)


if __name__ == "__main__":

    print bcolors.WARNING + "*** Start testing funds withdraw ***" + bcolors.ENDC

    print "#### PROPOSER'S BALANCE BEFORE CREATING PROPOSAL: ####"
    query_balance(_proposer)

    # create proposal
    _prop.send_create()

    print "#### PROPOSER'S BALANCE AFTER CREATING PROPOSAL: ####"
    query_balance(_proposer)

    print "#### CONTRIBUTOR'S BALANCE BEFORE FUNDING: ####"
    query_balance(_contributor)

    print "#### THIS PROPOSAL BEFORE FUNDING: ####"
    query_proposal(_encoded_pid)

    # fund proposal
    fund_proposal(_encoded_pid, _initial_funding, _contributor)
    time.sleep(5)

    print "#### THIS PROPOSAL AFTER FUNDING: ####"
    query_proposal(_encoded_pid)

    print "#### CONTRIBUTOR'S BALANCE AFTER FUNDING: ####"
    query_balance(_contributor)

    for x in range(_wait):
        print("wait for 60s, " + str(_wait * 10 - x * 10) + "s left")
        time.sleep(10)

    print "#### PROPOSER'S BALANCE BEFORE WITHDRAWING: ####"
    query_balance(_proposer)

    print "#### THIS PROPOSAL BEFORE WITHDRAWING: ####"
    query_proposal(_encoded_pid)

    # withdraw proposal funds---withdraw initial fund
    withdraw_fund(_encoded_pid, _proposer, _withdraw_amount, _proposer)
    time.sleep(5)

    print "#### THIS PROPOSAL AFTER WITHDRAWING: ####"
    query_proposal(_encoded_pid)

    print "#### PROPOSER'S BALANCE AFTER WITHDRAWING: ####"
    query_balance(_proposer)

    print "#### CONTRIBUTOR'S BALANCE BEFORE WITHDRAWING: ####"
    query_balance(_contributor)

    # withdraw proposal funds
    withdraw_fund(_encoded_pid, _contributor, _withdraw_amount, _contributor)
    time.sleep(5)

    print "#### CONTRIBUTOR'S BALANCE AFTER WITHDRAWING: ####"
    query_balance(_contributor)

    print "#### THIS PROPOSAL AFTER FULLY WITHDRAWING: ####"
    query_proposal(_encoded_pid)

    print bcolors.WARNING + "#### TRY TO WITHDRAW WHEN NO FUNDS TO WITHDRAW, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund(_encoded_pid, _contributor, _withdraw_amount, _contributor)
    time.sleep(5)

    print bcolors.WARNING + "#### TRY TO WITHDRAW OTHER CONTRIBUTOR'S FUNDS, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund_malicious(_proposer, _encoded_pid, _contributor, _withdraw_amount, _contributor)
    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    activeList = query_proposals("active")
    if len(activeList) != 0:
        sys.exit(-1)

    print "#### PASSED PROPOSALS: ####"
    passedList = query_proposals("passed")
    if len(passedList) != 0:
        sys.exit(-1)

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals("failed")
    if len(failedList) != 1:
        sys.exit(-1)