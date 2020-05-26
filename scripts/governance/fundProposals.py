import sys
import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_20021"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("3") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

_prop = Proposal(_pid, "general", "proposal for fund", _proposer, _initial_funding)

def fund_proposal(pid, amount, funder):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(2)

def check_proposal(pid, state_expected, status_expected):
    # check proposal state
    prop = query_proposal(pid)
    if prop['state'] != state_expected:
        sys.exit(-1)
    if prop['proposal']['Status'] != status_expected:
        sys.exit(-1)

if __name__ == "__main__":
    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # check proposal state
    check_proposal(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])
    check_proposal(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 3rd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[2])
    check_proposal(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")
    
