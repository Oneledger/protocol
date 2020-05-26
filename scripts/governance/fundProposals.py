import sys
import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_20018"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

_prop = Proposal(_pid, "general", "proposal for fund", _proposer, _initial_funding)

def fund_and_check_state(pid, amount, funder):
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(2)

if __name__ == "__main__":
    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.get_encoded_pid()

    # check proposal state
    prop = query_proposal(encoded_pid)

    # fund proposal
    prop_fund = ProposalFund(encoded_pid, _each_funding, addr_list[1])
    prop_fund.send_fund()
    time.sleep(3)

    # check proposal state
    prop = query_proposal(encoded_pid)
    if prop['state'] != ProposalStateActive:
        sys.exit(-1)
    if prop['proposal']['Status'] != ProposalStatusFunding:
        sys.exit(-1)

    # fund proposal again
    prop_fund = ProposalFund(encoded_pid, _each_funding, addr_list[2])
    prop_fund.send_fund()
    time.sleep(2)

    # check proposal state
    prop = query_proposal(encoded_pid)
    if prop['state'] != ProposalStateActive:
        sys.exit(-1)
    if prop['proposal']['Status'] != ProposalStatusVoting:
        sys.exit(-1)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")
    
