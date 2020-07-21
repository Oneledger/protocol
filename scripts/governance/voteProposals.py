from sdk import *

addr_list = addresses()

_pid_pass = "id_20033"
_pid_fail = "id_20043"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)

def test_pass_proposal():
    _prop = Proposal(_pid_pass, "general", "proposal for vote", "proposal headline", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> YES--10%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> YES--10% NO--20%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 3rd vote --> YES--40% NO--20%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_2, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 4th vote --> YES--80% NO--20%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_3, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeCompletedYes, ProposalStatusCompleted)

def test_fail_proposal():
    _prop = Proposal(_pid_fail, "general", "proposal for vote", "proposal headline", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(1)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> NO--0%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> NO--20%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 3rd vote --> NO--50%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_2, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeCompletedNo, ProposalStatusCompleted)

if __name__ == "__main__":
    # test pass a proposal
    test_pass_proposal()

    # test fail a proposal
    test_fail_proposal()

    print bcolors.OKGREEN + "#### Test vote proposals succeed" + bcolors.ENDC
    print ""
