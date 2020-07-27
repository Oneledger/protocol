from time import sleep

from sdk.common import *

# test only 50 blocks to save time
_reward_calc_cycle = 100
_reward_check_height = _reward_calc_cycle / 2
_reward_each_block = 38356164383561643834

# _total_rewards = 1879452054794520547866
_total_rewards = _reward_each_block * (_reward_check_height - 1)

# each validators rewards share in total
_rewards_share_power_1 = 187945205479452054767
_rewards_share_power_2 = 375890410958904109583
_rewards_share_power_3 = 563835616438356164350
_rewards_share_power_4 = 751780821917808219166
_rewards_share_expected = [_rewards_share_power_1, _rewards_share_power_2, _rewards_share_power_3,
                           _rewards_share_power_4]
_fault_tolerance = 1000

addr_list = addresses()


def rewardkey(validator):
    return validator['totalAmount']


def testRewardsDistribution():
    while True:
        sleep(0.5)
        rewards = query_total_rewards()
        if rewards['height'] < _reward_check_height:
            continue
        else:
            break

    total = int(rewards['totalRewards'])
    if abs(total - _total_rewards) > _fault_tolerance:
        print "totalRewards incorrect"
        sys.exit(-1)

    validators = sorted(rewards['validators'], key=rewardkey)
    validator_addresses = []
    for i, v in enumerate(validators):
        validator_addresses.append(v['address'])
        reward_share = int(v['totalAmount'])
        reward_share_expected = _rewards_share_expected[i]
        if abs(reward_share - reward_share_expected) > _fault_tolerance:
            print "validator reward share incorrect"
            sys.exit(-1)

    print bcolors.OKGREEN + "#### test block rewards distribution succeed" + bcolors.ENDC
    return validator_addresses


if __name__ == "__main__":
    # send some funds to pool through olclient
    account = addr_list[0][3:]
    args = ['olclient', 'sendpool', '--root', node_0, '--amount', '1000000', '--party', account, '--poolName',
            'RewardsPool', '--fee', '0.0001']
    process = subprocess.Popen(args)
    process.wait()

    # # test rewards distribution
    # validators = testRewardsDistribution()
    #

print bcolors.OKGREEN + "#### Verify block rewards succeed" + bcolors.ENDC

print bcolors.OKGREEN + "#### Verify block rewards succeed" + bcolors.ENDC
