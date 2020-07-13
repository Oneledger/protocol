from time import sleep

from sdk import *

_wait = 50

if __name__ == "__main__":
    staking_test = Staking(node_4)
    staking_test.prepare()
    addr_list = addresses()

    staking_test.stake('3', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 3, 0, False)
    staking_test.checkValidatorSet(4, True, 3)
    time.sleep(5)

    staking_test.unstake('1', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 2, 0, True)
    staking_test.checkValidatorSet(4, True, 2)
    time.sleep(5)

    staking_test.stake('1', addr_list[2], False)
    staking_test.unstake('1', addr_list[2], False)

    staking_test.unstake('2', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 0, 0, True)
    staking_test.checkValidatorSet(4, False, 0)

    print("wait for " + str(_wait) + "s")
    for x in range(_wait):
        print(str(_wait - x) + "s left")
        time.sleep(1)

    staking_test.checkStatus(addr_list[1], 0, 3, False)
    staking_test.withdraw('3', addr_list[1])
    staking_test.checkStatus(addr_list[1], 0, 0, False)

