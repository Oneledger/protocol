import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_30010"
_proposer = addr_list[0]
_initial_funding = (int("10023450") * 10 ** 14)
_contributor = addr_list[1]

_prop = Proposal(_pid, "general", "proposal for fund withdraw", _proposer, _initial_funding)
_prop_fund = ProposalFund(_pid, 100, _contributor)
_prop_fund_withdraw = ProposalFundsWithdraw(_pid, _contributor, 100, _contributor)

if __name__ == "__main__":
    # create proposal
    _prop.send_create()
    time.sleep(5)

    # fund proposal
    _prop_fund.send_fund()
    time.sleep(20)

    # withdraw proposal funds
    _prop_fund_withdraw.withdraw_fund()
    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")