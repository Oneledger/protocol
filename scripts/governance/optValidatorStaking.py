from sdk import *

addr_list = addresses()

_pid_fail = "id_40061"
_pid_pass = "id_40063"
_pid_pass2 = "id_40064"
_proposer = addr_list[0]
_initial_funding = 1000000000
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)
_initial_funding_too_less = 5000


def update_options(update):
    # Update Proposal to increse top validator count
    _prop = Proposal(_pid_pass, "configUpdate", "proposal for vote", "Headline", _proposer, _initial_funding, update)
    # state = _prop.default_gov_state()
    # state['stakingOptions']['topValidatorCount'] = 8
    # _prop.configupdate = state
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


def stake_genesis_initialValidators():
    stake(node_0)
    stake(node_1)
    stake(node_2)
    stake(node_3)


def stake_new_validators():
    stake(node_4)
    stake(node_5)
    stake(node_6)
    stake(node_7)
    stake(node_8)


if __name__ == "__main__":
    if getActiveValidators() != 4:
        sys.exit(-1)
    print bcolors.OKBLUE + "#### Initial Active Validator count : " + str(getActiveValidators()) + bcolors.ENDC
    #  Increasing the Staking of genesis validators so that they stay on top
    stake_genesis_initialValidators()
    #  Staking of genesis validators so that they stay on top
    stake_new_validators()
    time.sleep(1)
    if getActiveValidators() != 8:
        sys.exit(-1)
    print bcolors.OKBLUE + "#### Active Validator count : " + str(getActiveValidators()) + bcolors.ENDC
    update_param = 10
    update = "stakingOptions.topValidatorCount:" + str(update_param)
    update_options(update)
    time.sleep(1)
    opt = query_governanceState()
    print "Last Update Height:" + str(opt["lastUpdateHeight"]["staking"])
    if getActiveValidators() != 9:
        sys.exit(-1)
    print bcolors.OKBLUE + "#### Top Validator count updated to : " + str(update_param) + bcolors.ENDC
    time.sleep(1)
    print bcolors.OKBLUE + "#### Active Validator count : " + str(getActiveValidators()) + bcolors.ENDC

    clean_and_catchup()
