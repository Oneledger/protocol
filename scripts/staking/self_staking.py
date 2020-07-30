from sdk import *

_wait = 50

if __name__ == "__main__":
    staking_test = Staking(node_4)
    new_staking_address = staking_test.prepare()
    addr_list = addresses()
    original_staking_address = staking_test.staking_address

    # stake
    staking_test.stake('4000000', True)
    staking_test.checkStatus(4000000, 0, False)
    staking_test.checkValidatorSet(4, True, 4000000)
    time.sleep(2)

    # unstake
    staking_test.unstake('500000', True)
    staking_test.checkStatus(3500000, 0, True)
    staking_test.checkValidatorSet(4, True, 3500000)
    time.sleep(2)

    # stake from different address, should fail due to stake address mismatch

    staking_test.staking_address = new_staking_address
    staking_test.stake('1', False)

    # unstake from different address, should fail due to stake address mismatch
    staking_test.unstake('1', False)

    # unstake all
    staking_test.staking_address = original_staking_address
    staking_test.unstake('3500000', True)
    staking_test.checkStatus(0, 0, True)
    staking_test.checkValidatorSet(4, False, 0)

    # stake again using new stake address.
    # We allow this Tx, but you won't be able to withdraw old money due to stake address mismatch
    staking_test.staking_address = new_staking_address
    staking_test.stake('4000000', True)
    staking_test.checkStatus(4000000, 0, False)
    staking_test.checkValidatorSet(4, True, 4000000)

    # waits for above unstaked old amount to mature
    print("wait for " + str(_wait) + "s")
    for x in range(_wait):
        print(str(_wait - x) + "s left")
        time.sleep(1)

    # withdraw old money, should fail due to stake address mismatch
    staking_test.staking_address = original_staking_address
    staking_test.checkStatus(0, 4000000, False)
    staking_test.withdraw('6', False)

    # unstake all new money
    staking_test.staking_address = new_staking_address
    staking_test.unstake('4000000', True)
    staking_test.checkStatus(0, 0, True)
    staking_test.checkValidatorSet(4, False, 0)

    # withdraw old money again, should success
    staking_test.staking_address = original_staking_address
    staking_test.checkStatus(0, 4000000, False)
    staking_test.withdraw('4000000', True)
    staking_test.checkStatus(0, 0, False)
