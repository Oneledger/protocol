
import time
import sys

from sdk.actions import *


if __name__ == "__main__":
    addrs = addresses()
    print addrs

    create_price = (int("10023450")*10**14)
    print "create price:", create_price
    name = "alice.ol"
    raw_txn = create_domain(name, addrs[0], create_price)
    print "raw create domain tx:", raw_txn

    signed = sign(raw_txn, addrs[0])
    print "signed create domain tx:", signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    raw_txn = send_domain(name, addrs[0], "10")
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

    sell_price = (int("105432")*10**14)
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

    print "Get Domain on Sale"
    resp = get_domain_on_sale()
    print resp

    print "#################" \
          "##"
    print

    raw_txn = cancel_sell_domain(name, addrs[0], sell_price)
    print raw_txn
    print

    signed = sign(raw_txn, addrs[0])
    print signed
    print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print
    if result["ok"] != True:
        sys.exit(-1)

    print "Get Domain on Sale"
    resp = get_domain_on_sale()
    print resp

    print "#################" \
          "##"
    print
