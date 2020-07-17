import time
import subprocess
from actions import *


def WithdrawRewards(Walletaddress, amount, secs=1):
    # fund the proposal
    withdraw = Withdraw(Walletaddress, amount)
    withdraw.send_withdraw()
    time.sleep(secs)
