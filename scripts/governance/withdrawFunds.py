import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_30034"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_funder = addr_list[1]
_funder_never_fund = addr_list[2]
_funds_amount = (int("2") * 10 ** 9)
_withdraw_amount = (int("2") * 10 ** 9)
_withdraw_amount_too_much = (int("5") * 10 ** 9)

_prop = Proposal(_pid, "general", "proposal for funds withdrawing", "proposal headline", _proposer, _initial_funding)
_encoded_pid = _prop.get_encoded_pid()

_wait = 6


def fund_proposal(pid, amount, funder):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(2)


def withdraw_fund(pid, funder, amount, beneficiary):
    # fund the proposal
    fund_withdraw = ProposalFundsWithdraw(pid, funder, amount, beneficiary)
    fund_withdraw.withdraw_fund(funder)
    time.sleep(2)


def withdraw_fund_malicious(wrong_funder, pid, funder, amount, beneficiary):
    # fund the proposal
    fund_withdraw = ProposalFundsWithdraw(pid, funder, amount, beneficiary)
    fund_withdraw.withdraw_fund(wrong_funder)
    time.sleep(2)


if __name__ == "__main__":

    print bcolors.WARNING + "*** Start testing funds withdraw ***" + bcolors.ENDC

    print "#### PROPOSER'S BALANCE BEFORE CREATING PROPOSAL: ####"
    query_balance(_proposer)

    # create proposal
    _prop.send_create()

    print "#### PROPOSER'S BALANCE AFTER CREATING PROPOSAL: ####"
    query_balance(_proposer)

    print "#### FUNDER'S BALANCE BEFORE FUNDING: ####"
    query_balance(_funder)

    print "#### THIS PROPOSAL BEFORE FUNDING: ####"
    query_proposal(_encoded_pid)

    # fund proposal
    fund_proposal(_encoded_pid, _initial_funding, _funder)
    time.sleep(5)

    print "#### THIS PROPOSAL AFTER FUNDING: ####"
    query_proposal(_encoded_pid)

    print "#### FUNDER'S BALANCE AFTER FUNDING: ####"
    query_balance(_funder)

    for x in range(_wait):
        print("wait for 60s, " + str(_wait * 10 - x * 10) + "s left")
        time.sleep(10)

    print bcolors.WARNING + "#### TRY TO WITHDRAW NOT FUNDED PROPOSAL, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund(_encoded_pid, _funder_never_fund, _withdraw_amount_too_much, _funder_never_fund)
    time.sleep(5)

    print bcolors.WARNING + "#### TRY TO WITHDRAW MORE THAN FUNDED AMOUNT, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund(_encoded_pid, _funder, _withdraw_amount_too_much, _funder)
    time.sleep(5)

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

    print "#### FUNDER'S BALANCE BEFORE WITHDRAWING: ####"
    query_balance(_funder)

    # withdraw proposal funds
    withdraw_fund(_encoded_pid, _funder, _withdraw_amount, _funder)
    time.sleep(5)

    print "#### FUNDER'S BALANCE AFTER WITHDRAWING: ####"
    query_balance(_funder)

    print "#### THIS PROPOSAL AFTER FULLY WITHDRAWING: ####"
    query_proposal(_encoded_pid)

    print bcolors.WARNING + "#### TRY TO WITHDRAW WHEN NO FUNDS TO WITHDRAW, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund(_encoded_pid, _funder, _withdraw_amount, _funder)
    time.sleep(5)

    print bcolors.WARNING + "#### TRY TO WITHDRAW OTHER FUNDER'S FUNDS, SHOULD FAIL: ####" + bcolors.ENDC
    withdraw_fund_malicious(_proposer, _encoded_pid, _funder, _withdraw_amount, _funder)
    time.sleep(5)

    # below applied when run withdrawFunds.py only
    # print "#### ACTIVE PROPOSALS: ####"
    # activeList = query_proposals("active")
    # if len(activeList) != 0:
    #     sys.exit(-1)
    #
    # print "#### PASSED PROPOSALS: ####"
    # passedList = query_proposals("passed")
    # if len(passedList) != 0:
    #     sys.exit(-1)

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals(ProposalStateFailed)
    if len(failedList) == 0:
        sys.exit(-1)