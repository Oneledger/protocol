from sdk import *

addr_list = addresses()

_pid_fail = "id_20061"
_pid_pass = "id_20063"
_pid_pass2 = "id_20064"
_proposer = addr_list[0]
_initial_funding = 1000000000
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)
_initial_funding_too_less = 100000000


def test_catchup():
    # Create Proposal should Fail becuase funding is too high
    _prop = Proposal(_pid_fail, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding_too_less)
    _prop.send_create()
    time.sleep(3)

    # Update Proposal to decrese initial funding
    _prop = Proposal(_pid_pass, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding)
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid
    # 1st fund
    fund_proposal(encoded_pid, _funding_goal_general, addr_list[0])

    # 1st vote --> 25%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])
    # check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    # 2nd vote --> 25%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[0])
    # check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    # 3rd vote --> 50%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_2, addr_list[0])
    # check_proposal_state(encoded_pid, ProposalStateActive, ProposalStatusVoting)

    # 4th vote --> 75%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_3, addr_list[0])
    # check_proposal_state(encoded_pid, ProposalStatePassed, ProposalStatusCompleted)

    time.sleep(3)

    # Create propsal with lesser funding amount should now pass
    _prop = Proposal(_pid_fail, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding_too_less)
    _prop.send_create()
    time.sleep(3)


if __name__ == "__main__":
    print "#### Governance State Before : ####"
    opt = query_governanceState()
    print "propOptions.configUpdate.initialFunding :" + opt["propOptions"]["configUpdate"]["initialFunding"]
    test_catchup()
    #
    # print "#### FINALIZED PROPOSALS: ####"
    # proposalstats = query_proposals(0X34)
    print "#### Governance State After : ####"
    opt = query_governanceState()
    print "propOptions.configUpdate.initialFunding :" + opt["propOptions"]["configUpdate"]["initialFunding"]

#
# print proposalstats["height"]
#
# print "#### FINALIZEFAILED PROPOSALS: ####"
# query_proposals("finalizeFailed")
