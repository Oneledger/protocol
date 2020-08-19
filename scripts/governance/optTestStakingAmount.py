from sdk import *

addr_list = addresses()

_pid_pass_0 = "id_30067"
_pid_pass_1 = "id_30068"
_pid_pass_2 = "id_30069"
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


def stake_new_validators_less_than_min():
    stake(node_4, '2000000')
    stake(node_5, '2000000')


if __name__ == "__main__":
    stake_new_validators_less_than_min()
    if getActiveValidators() != 4:
        sys.exit(-1)
    if getAllValidators() != 6:
        sys.exit(-1)
    print bcolors.OKBLUE + "#### Initial Active Validator count : " + str(
        getActiveValidators()) + "  Total Validator Count :" + str(getAllValidators()) + bcolors.ENDC
    update_param = 2000000
    update = {"stakingOptions.minSelfDelegationAmount": update_param}
    opt = query_governanceState()
    old_height_proposal = opt["lastUpdateHeight"]["proposal"]
    old_height_staking = opt["lastUpdateHeight"]["staking"]
    old_height_fee = opt["lastUpdateHeight"]["fee"]

    test_change_gov_options(update, _pid_pass_1)
    opt = query_governanceState()
    new_height = opt["lastUpdateHeight"]["staking"]
    if opt["govOptions"]["stakingOptions"]["minSelfDelegationAmount"] != str(update_param):
        print "Value not updated"
        sys.exit(-1)
    if opt["lastUpdateHeight"]["proposal"] - old_height_proposal != 0:
        print "Height changed proposal ,Update type staking"
        sys.exit(-1)
    if opt["lastUpdateHeight"]["staking"] - old_height_staking <= 0:
        print "Height not changed staking ,Update type staking"
        sys.exit(-1)
    if opt["lastUpdateHeight"]["fee"] - old_height_fee != 0:
        print "Height changed fee , Update type staking"
        sys.exit(-1)
    print bcolors.OKBLUE + "Option Update Successful : stakingOptions.maturityTime | At Height " + str(
        opt["lastUpdateHeight"]["staking"]) + bcolors.ENDC

    if getActiveValidators() != 6:
        sys.exit(-1)
    if getAllValidators() != 6:
        sys.exit(-1)

    print bcolors.OKBLUE + "#### Active Validator/Total Validator count : " + str(getActiveValidators()) + " / " + str(
        getAllValidators()) + " | New Validators Added | minSelfDelegationAmount = " + str(update_param) + bcolors.ENDC

    update_param = 6000000
    update = {"stakingOptions.minSelfDelegationAmount": update_param}
    test_change_gov_options(update, _pid_pass_2)
    opt = query_governanceState()

    if getActiveValidators() != 3:
        sys.exit(-1)
    if getAllValidators() != 6:
        sys.exit(-1)

    print bcolors.OKBLUE + "#### Active Validator/Total Validator count : " + str(getActiveValidators()) + " / " + str(
        getAllValidators()) + " | Validators Removed | minSelfDelegationAmount = " + str(update_param) + bcolors.ENDC

    clean_and_catchup()
