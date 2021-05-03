from sdk import *

addr_list = addresses()

_pid = "id_20000"
_pid1 = "id_20001"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_funding_goal = (int("10") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)


def expired_votes():
    _prop = Proposal(_pid1, "general", "proposal for vote expired", "proposal headline", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> 25%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> 25%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)


def completed_votes():
    _prop = Proposal(_pid, "general", "proposal for vote", "proposal headline", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _each_funding, addr_list[0])

    # 2nd fund
    fund_proposal(encoded_pid, _each_funding, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 1st vote --> 25%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 2nd vote --> 25%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[1])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 3rd vote --> 50%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_2, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeInProgress, ProposalStatusVoting)

    # 4th vote --> 75%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_3, addr_list[2])
    check_proposal_state(encoded_pid, ProposalOutcomeCompletedYes, ProposalStatusCompleted)


def show_proposals():
    print "#### ACTIVE PROPOSALS: ####"
    activeList = query_proposals(ProposalStateActive)
    print activeList

    print "#### PASSED PROPOSALS: ####"
    passedList = query_proposals(ProposalStatePassed)
    print passedList

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals(ProposalStateFailed)
    print failedList


if __name__ == "__main__":

    completed_votes()
    show_proposals()

    expired_votes()
    time.sleep(31)
    show_proposals()
