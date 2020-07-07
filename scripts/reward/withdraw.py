from sdk.actions import *
from sdk.common import *

addr_list = addresses()


def withdraw_rewards():
    # using address for signing
    WithdrawRewards(addr_list[0], "100")


if __name__ == "__main__":
    withdraw_rewards()
