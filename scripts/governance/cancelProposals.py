import sys
import time
from sdk import *

addr_list = addresses()

_pid = 20036
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("3") * 10 ** 9)
_big_funding = (int("8") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

def gen_prop():
    global _pid
    prop = Proposal(str(_pid), "general", "proposal for fund", "proposal headline", _proposer, _initial_funding)
    _pid += 1
    return prop

def test_normal_cancel():
    # create proposal
    prop = gen_prop()
    prop.send_create()
    time.sleep(1)
    encoded_pid = prop.pid

    # check proposal state
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

    # cancel this proposal
    cancel_proposal(encoded_pid, _proposer, "changed mind")
    check_proposal_state(encoded_pid, ProposalStateFailed, ProposalStatusCompleted)
    return encoded_pid

def test_cancel_noactive_proposal(pid_not_active):
    # cancel this no-active proposal, should fail
    res = cancel_proposal(pid_not_active, _proposer, "try a weird cancel")
    if res:
        sys.exit(-1)
    check_proposal_state(pid_not_active, ProposalStateFailed, ProposalStatusCompleted)

def test_cancel_proposal_in_voting_status():
    # create proposal
    prop = gen_prop()
    prop.send_create()
    time.sleep(1)
    encoded_pid = prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _big_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    # cancel this proposal, should fail
    res = cancel_proposal(encoded_pid, _proposer, "too late to changed mind")
    if res:
        sys.exit(-1)
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

def test_cancel_someone_else_proposal():
    # create proposal
    prop = gen_prop()
    prop.send_create()
    time.sleep(1)
    encoded_pid = prop.pid

    # cancel this proposal, should fail
    res = cancel_proposal(encoded_pid, addr_list[1], "do bad things")
    if res:
        sys.exit(-1)
    check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusFunding)

if __name__ == "__main__":
    pid_canceled = test_normal_cancel()

    test_cancel_noactive_proposal(pid_canceled)

    test_cancel_proposal_in_voting_status()

    test_cancel_someone_else_proposal()

    print "#### Test cancel proposals succeed: ####"
    print ""
    