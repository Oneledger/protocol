from random import randint

from sdk import *

addr_list = addresses()

_pid_pass_0 = "id_30067"
_pid_pass_1 = "id_30068"
_pid_pass_2 = "id_30069"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)


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
    update = "stakingOptions.minSelfDelegationAmount:" + str(update_param)
    updateGov(update, "staking", "id_" + str(randint(10000, 99999)), True, False)

    if getActiveValidators() != 6:
        sys.exit(-1)
    if getAllValidators() != 6:
        sys.exit(-1)
    #
    print bcolors.OKBLUE + "#### Active Validator/Total Validator count : " + str(getActiveValidators()) + " / " + str(
        getAllValidators()) + " | New Validators Added | minSelfDelegationAmount = " + str(update_param) + bcolors.ENDC

    update_param = 6000000
    update = "stakingOptions.minSelfDelegationAmount:" + str(update_param)
    updateGov(update, "staking", "id_" + str(randint(10000, 99999)), True, False)
    opt = query_governanceState()

    if getActiveValidators() != 3:
        sys.exit(-1)
    if getAllValidators() != 6:
        sys.exit(-1)

    print bcolors.OKBLUE + "#### Active Validator/Total Validator count : " + str(getActiveValidators()) + " / " + str(
        getAllValidators()) + " | Validators Removed | minSelfDelegationAmount = " + str(update_param) + bcolors.ENDC

    clean_and_catchup()
