from time import sleep

from sdk import *

_wait = 50

if __name__ == "__main__":
    staking_test = Staking(node_4)
    staking_test.prepare()
    addr_list = addresses()

    staking_test.stake('6', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 6, 0, False)
    staking_test.checkValidatorSet(4, True, 6)
    time.sleep(5)

    staking_test.unstake('1', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 5, 0, True)
    staking_test.checkValidatorSet(4, True, 5)
    time.sleep(5)

    staking_test.stake('1', addr_list[2], False)
    staking_test.unstake('1', addr_list[2], False)

    staking_test.unstake('5', addr_list[1], True)
    staking_test.checkStatus(addr_list[1], 0, 0, True)
    staking_test.checkValidatorSet(4, False, 0)

    print("wait for " + str(_wait) + "s")
    for x in range(_wait):
        print(str(_wait - x) + "s left")
        time.sleep(1)

    staking_test.checkStatus(addr_list[1], 0, 6, False)
    staking_test.withdraw('6', addr_list[1])
    staking_test.checkStatus(addr_list[1], 0, 0, False)

