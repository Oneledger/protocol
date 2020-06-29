import time

from actions import *


def WithdrawRewards(Walletaddress, secs=1):
    # fund the proposal
    withdraw = Withdraw(Walletaddress)
    withdraw.send_withdraw()
    time.sleep(secs)
