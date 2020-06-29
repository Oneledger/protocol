from sdk.common import *
from sdk.actions import *



addr_list = addresses()

def withdraw_rewards():
    # using address for signing
    WithdrawRewards(addr_list[0])


if __name__ == "__main__":
    withdraw_rewards()
