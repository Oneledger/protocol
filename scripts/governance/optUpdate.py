from random import randint

from sdk import *

addr_list = addresses()

_pid_pass_0 = "id_20067"
_pid_pass_1 = "id_20068"
_pid_pass_2 = "id_20069"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_each_funding = (int("5") * 10 ** 9)
_funding_goal_general = (int("10") * 10 ** 9)




if __name__ == "__main__":
    update_param = {
        "burn": 10,
        "executionCost": 20,
        "bountyPool": 50,
        "validators": 10,
        "proposerReward": 0,
        "feePool": 10
    }
    update = {"propOptions.configUpdate.passedFundDistribution": update_param}
    updateGov(update, "proposal", "id_" + str(randint(10000, 99999)), True)

    update_param = 109201
    update = {"stakingOptions.maturityTime": update_param}
    updateGov(update, "staking", "id_" + str(randint(10000, 99999)), True)

    update_param = 10
    update = {"feeOption.minFeeDecimal": update_param}
    updateGov(update, "fee", "id_" + str(randint(10000, 99999)), True)

    update_param = 800
    update = {"evidenceOptions.minVotesRequired": update_param}
    updateGov(update, "evidence", "id_" + str(randint(10000, 99999)), True)

    update_param = 100000000000
    update = {"onsOptions.perBlockFees": update_param}
    updateGov(update, "ons", "id_" + str(randint(10000, 99999)), True, False)

    update_param = 10000
    update = {"evidenceOptions.blockVotesDiff": update_param}
    updateGov(update, "evidence", "id_" + str(randint(10000, 99999)), False)

    update_param = 1000
    update = {"evidenceOptions.blockVotesDiff": update_param}
    updateGov(update, "evidence", "id_" + str(randint(10000, 99999)), True)

    update_param = 150002
    update = {"propOptions.configUpdate.votingDeadline": update_param}
    updateGov(update, "proposal", "id_" + str(randint(10000, 99999)), True, False)

    update_param = {
        "penaltyBountyPercentage": 20,
        "penaltyBurnPercentage": 80,
    }
    update = {"evidenceOptions.penaltyPercentage": update_param}
    updateGov(update, "evidence", "id_" + str(randint(10000, 99999)), True, False)

    update_param = {
        "penaltyBountyPercentage": 10,
        "penaltyBurnPercentage": 80,
    }
    update = {"evidenceOptions.penaltyPercentage": update_param}
    updateGov(update, "evidence", "id_" + str(randint(10000, 99999)), False, False)

    clean_and_catchup()
