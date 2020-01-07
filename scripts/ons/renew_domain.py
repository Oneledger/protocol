import time
from sdk.actions import *

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
    # Create New account
    result = new_account('charlie')
    print result

    name = "expiring.ol"
    addrs = addresses()

    print addrs

    # expiring after one block
    create_price = (int("10000001")*10**14)
    print "create price:", create_price

    raw_txn = create_domain(name, addrs[0], create_price)
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print

    time.sleep(2)

    print_all_domains(addrs[0])

    print bcolors.WARNING + "*** Renewing domain ***" + bcolors.ENDC

    raw_txn = renew_domain(name, addrs[0], (int("200")*10**18))
    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print

    print_all_domains(addrs[0])