import time

from sdk.actions import *

addr_list = addresses()
addr_list_1 = addresses_from_second_node()

_pid = "id_30025"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_contributor = addr_list_1[0]
_funds_amount = (int("2") * 10 ** 5)

_prop = Proposal(_pid, "general", "proposal for funds withdrawing", _proposer, _initial_funding)
_encoded_pid = _prop.get_encoded_pid()
_prop_fund = ProposalFund(_pid, _funds_amount, _contributor)
_prop_fund_withdraw_initial_funding = ProposalFundsWithdraw(_encoded_pid, _proposer, _funds_amount, _proposer)
_prop_fund_withdraw = ProposalFundsWithdraw(_encoded_pid, _contributor, _funds_amount, _contributor)

_wait = 6

if __name__ == "__main__":
    print bcolors.WARNING + "*** Start testing funds withdraw ***" + bcolors.ENDC
    # create proposal
    _prop.send_create()

    # fund proposal
    _prop_fund.send_fund()
    time.sleep(5)

    for x in range(_wait):
        print("wait for 60s, " + str(_wait * 10 - x * 10) + "s left")
        time.sleep(10)

    # withdraw proposal funds---withdraw initial fund
    _prop_fund_withdraw_initial_funding.withdraw_fund()
    time.sleep(5)

    # withdraw proposal funds
    _prop_fund_withdraw_initial_funding.withdraw_fund()
    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")