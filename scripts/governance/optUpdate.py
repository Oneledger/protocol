from sdk import *

addr_list = addresses()

_pid_pass = "id_20061"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)


def test_change_gov_options():
    _prop = Proposal(_pid_pass, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding)

    # create proposal
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


if __name__ == "__main__":
    print "#### Governance State Before : ####"
    opt = query_governanceState()
    print "bitcoinChainDriverOption.ChainType :" + opt["bitcoinChainDriverOption"]["ChainType"]
    test_change_gov_options()
    #
    # print "#### FINALIZED PROPOSALS: ####"
    # proposalstats = query_proposals(0X34)
    print "#### Governance State After : ####"
    opt = query_governanceState()
    print "bitcoinChainDriverOption.ChainType :" + opt["bitcoinChainDriverOption"]["ChainType"]

#
# print proposalstats["height"]
#
# print "#### FINALIZEFAILED PROPOSALS: ####"
# query_proposals("finalizeFailed")
