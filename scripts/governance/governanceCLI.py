import sys
import time
from sdk import *

addr_list = addresses()

_pid_pass = "id_30010"
_pid_fail = "id_30011"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

def test_pass_proposal_cli():
    _prop = Proposal(_pid_pass, "general", "proposal for vote", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(1)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> 25%
    vote_proposal_cli(encoded_pid, "YES", node_0, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> 25%
    vote_proposal(encoded_pid, "NO", url_1, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 3rd vote --> 50%
    vote_proposal(encoded_pid, "YES", url_2, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 4th vote --> 75%
    vote_proposal(encoded_pid, "YES", url_3, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeCompleted, ProposalStatusCompleted)

    # list proposal under another node
    list_proposal_cli(encoded_pid, node_1)

def test_fail_proposal_cli():
    _prop = Proposal(_pid_fail, "general", "proposal for vote", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> NO--0%
    vote_proposal_cli(encoded_pid, "NO", node_0, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> NO--25%
    vote_proposal(encoded_pid, "YES", url_1, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 3rd vote --> NO--50%
    vote_proposal(encoded_pid, "NO", url_2, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeInsufficientVotes, ProposalStatusCompleted)

    # list proposal
    list_proposal_cli(encoded_pid, node_2)

if __name__ == "__main__":
    # test pass a proposal using cli
    test_pass_proposal_cli()

    # test fail a proposal using cli
    test_fail_proposal_cli()
