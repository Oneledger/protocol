from sdk import *

addr_list = addresses()

_pid_pass = "id_20059"
_pid_fail = "id_20063"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)


def test_pass_finalize_proposal():
    _prop = Proposal(_pid_pass, "general", "headline", "proposal for vote", _proposer, _initial_funding)

    # create proposal
    _prop.send_create()
    time.sleep(3)
    encoded_pid = _prop.pid

    # 1st fund
    fund_proposal(encoded_pid, _funding_goal_general, addr_list[0])

    # 1st vote --> 25%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_0, addr_list[0])

    # 2nd vote --> 25%
    vote_proposal(encoded_pid, OPIN_NEGATIVE, url_1, addr_list[0])

    # 3rd vote --> 50%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_2, addr_list[0])

    # 4th vote --> 75%
    vote_proposal(encoded_pid, OPIN_POSITIVE, url_3, addr_list[0])


    time.sleep(3)
    return encoded_pid


if __name__ == "__main__":
    # test pass a proposal
    pid = test_pass_finalize_proposal()
    prop, funds = query_proposal(pid)
    if prop["outcome"] != 49:
        print "Exiting Outcome is not ProposalOutcomeCompletedYes"
        sys.exit(1)

    print "PassedProposals ###############"
    print query_proposals(ProposalStatePassed)

    print "FinalizedProposlals ################"
    pList = query_proposals(ProposalStateFinalized)
    print pList
    # if len(pList) == 0:
    #     print "Exiting Proposal was not finalized"
    #     sys.exit(1)
