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
    fund_withdraw.withdraw_fund()
    time.sleep(2)


if __name__ == "__main__":
    print _funds_amount
    print bcolors.WARNING + "*** Start testing funds withdraw ***" + bcolors.ENDC

    print "#### PROPOSER'S BALANCE BEFORE CREATING PROPOSAL: ####"
    query_balance(_proposer)

    # create proposal
    _prop.send_create()

    print "#### PROPOSER'S BALANCE AFTER CREATING PROPOSAL: ####"
    query_balance(_proposer)

    print "#### CONTRIBUTOR'S BALANCE BEFORE FUNDING PROPOSAL: ####"
    query_balance(_contributor)

    # fund proposal
    fund_proposal(_encoded_pid, _initial_funding, _contributor)
    time.sleep(5)

    print "#### CONTRIBUTOR'S BALANCE AFTER FUNDING PROPOSAL: ####"
    query_balance(_contributor)

    for x in range(_wait):
        print("wait for 60s, " + str(_wait * 10 - x * 10) + "s left")
        time.sleep(10)

    print "#### PROPOSER'S BALANCE BEFORE WITHDRAWING: ####"
    query_balance(_proposer)

    # withdraw proposal funds---withdraw initial fund
    withdraw_fund(_encoded_pid, _proposer, _withdraw_amount, _proposer)
    time.sleep(5)

    print "#### PROPOSER'S BALANCE AFTER WITHDRAWING: ####"
    query_balance(_proposer)

    print "#### CONTRIBUTOR'S BALANCE BEFORE WITHDRAWING: ####"
    query_balance(_contributor)

    # withdraw proposal funds
    withdraw_fund(_encoded_pid, _contributor, _withdraw_amount, _contributor)
    time.sleep(5)

    print "#### CONTRIBUTOR'S BALANCE AFTER WITHDRAWING: ####"
    query_balance(_contributor)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")