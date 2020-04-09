import time
from sdk.actions import *
import sys


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


if __name__ == "__main__":

    name = "expiring2.ol"
    addrs = addresses()
    if len(addrs) < 4:
        # Create New account
        result = new_account('charlie2')
        print result
        addrs = addresses()

    print addrs

    # expiring after one block
    create_price = (int("10000010") * 10 ** 14)
    print "create price:", create_price

    print bcolors.WARNING + "*** Create domain ***" + bcolors.ENDC
    raw_txn = create_domain(name, addrs[0], create_price)
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result

    if not result["ok"]:
        sys.exit(-1)

    print "###################"
    print
