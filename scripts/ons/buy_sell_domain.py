

import sys
import requests
import json
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

    name = "bob2.ol"

    addrs = addresses()

    print addrs

    create_price = (int("10002345")*10**14)
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

    if result["ok"] != True:
        sys.exit(-1)

    sell_price = (int("105432")*10**14)

    raw_txn = send_domain(name, addrs[0], (int("100")*10**18))
    print raw_txn

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "#################" \
          "##"
    print

    if result["ok"] != True:
        sys.exit(-1)

    time.sleep(2)
    raw_txn = sell_domain(name, addrs[0], sell_price)
    print raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = send(addrs[0], addrs[3], (int("2000")*10**18))
    print result
    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    print bcolors.WARNING + "*** Buying domain ***" + bcolors.ENDC

    raw_txn = buy_domain(name, addrs[3], (int("20")*10**18))
    print addrs[3]
    print addrs
    signed = sign(raw_txn, addrs[3])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print
    if result["ok"] != True:
        sys.exit(-1)

    print bcolors.WARNING + "*** Putting Domain on sale ***" + bcolors.ENDC
    raw_txn = sell_domain(name, addrs[3], (int("993242")*10**14))
    print raw_txn
    print

    signed = sign(raw_txn, addrs[3])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "############################################"
    print
    if result["ok"] != True:
        sys.exit(-1)

    print bcolors.WARNING + "*** Get Domains on sale ***" + bcolors.ENDC
    resp = get_domain_on_sale()
    print resp
