import sys
import time
from sdk import *

addr_list = addresses()

_pid = "id_20020"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("3") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

_prop = Proposal(_pid, "general", "proposal for fund", _proposer, _initial_funding)

if __name__ == "__main__":
    # create proposal
    _prop.send_create()
    time.sleep(1)
    encoded_pid = _prop.pid

    # check proposal state
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 3rd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[2])
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")
    
