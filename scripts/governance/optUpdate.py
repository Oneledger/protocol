from benedict import benedict

from sdk import *

addr_list = addresses()

_pid_pass_0 = "id_20067"
_pid_pass_1 = "id_20068"
_pid_pass_2 = "id_20069"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)


def test_change_gov_options(update, pid):
    _prop = Proposal(pid, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding,
                     update)
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


def updateGov(update, updatetype, p, test_type, checkValue=True):
    opt = query_governanceState()
    old_heights = [None] * len(option_types)
    for idx, val in enumerate(option_types):
        old_heights[idx] = opt["lastUpdateHeight"][val]
    test_change_gov_options(update, p)
    opt = query_governanceState()
    opt = benedict(opt)
    if checkValue and opt['govOptions.' + update.keys()[0]] != update.values()[0] and test_type:
        print "Value not updated"
        sys.exit(-1)
    for idx, val in enumerate(option_types):
        new_height_type = opt["lastUpdateHeight"][val]
        if val == updatetype and new_height_type - old_heights[idx] <= 0 and test_type:
            print "Height not changed" + str(val) + "Test Type :" + str(test_type)
            sys.exit(-1)
        if val == updatetype and new_height_type - old_heights[idx] > 0 and not test_type:
            print "Height changed for negative update " + str(val) + " Test Type :" + str(test_type)
            sys.exit(-1)
        if val != updatetype and new_height_type - old_heights[idx] > 0 and test_type:
            print "Height changed" + str(val) + "Test Type :" + str(test_type)
            sys.exit(-1)
    if test_type:
        print bcolors.OKBLUE + "Option Update Successful : " + str(update.keys()[0]) + "| At Height " + str(
            opt["lastUpdateHeight"][updatetype]) + bcolors.ENDC
    if not test_type:
        print bcolors.OKGREEN + "Option Update NOT Successful (Validation Failed) : " + str(
            update.keys()[0]) + " | For Value " + str(update.values()[0]) + bcolors.ENDC


if __name__ == "__main__":
    # update_param = {
    #     "burn": 10,
    #     "executionCost": 20,
    #     "bountyPool": 50,
    #     "validators": 10,
    #     "proposerReward": 0,
    #     "feePool": 10
    # }
    # update = {"propOptions.configUpdate.passedFundDistribution": update_param}
    # updateGov(update, "proposal", "id_20067", True)
    #
    # update_param = 109201
    # update = {"stakingOptions.maturityTime": update_param}
    # updateGov(update, "staking", "id_20068", True)
    #
    # update_param = 10
    # update = {"feeOption.minFeeDecimal": update_param}
    # updateGov(update, "fee", "id_20069", True)
    #
    # update_param = 800
    # update = {"evidenceOptions.minVotesRequired": update_param}
    # updateGov(update, "evidence", "id_20070", True)
    #
    update_param = 1000
    update = {"evidenceOptions.blockVotesDiff": update_param}
    updateGov(update, "evidence", "id_20071", False)

    update_param = {
        "penaltyBountyPercentage": 20,
        "penaltyBurnPercentage": 80,
    }
    update = {"evidenceOptions.penaltyPercentage": update_param}
    updateGov(update, "evidence", "id_20072", True, False)

    update_param = {
        "penaltyBountyPercentage": 10,
        "penaltyBurnPercentage": 80,
    }
    update = {"evidenceOptions.penaltyPercentage": update_param}
    updateGov(update, "evidence", "id_20072", False, False)

    clean_and_catchup()
