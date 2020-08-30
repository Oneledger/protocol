import json
import sys
import time

from sdk.actions import *

if __name__ == "__main__":
    addrs = addresses()
    print addrs

    create_price = (int("10002345") * 10 ** 14)
    print "create price:", create_price

    time.sleep(3)
    name = "xyzzz.ol"
    raw_txn = create_domain(name, addrs[0], create_price)
    print "raw create domain tx:", raw_txn

    signed = sign(raw_txn, addrs[0])
    print "signed create domain tx:", signed
    print
    time.sleep(3)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)


    print "Create subdomain abc.xyzzz.ol"
    name = "abc.xyzzz.ol"

    raw_txn = create_sub_domain(name, addrs[0], create_price, '')
    signed = sign(raw_txn, addrs[0])
    time.sleep(3)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print json.dumps(result)
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    time.sleep(3)
    print "send to subdomain abc.xyzzz.ol"
    raw_txn = send_domain(name, addrs[0], "10")
    signed = sign(raw_txn, addrs[0])
    time.sleep(3)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print
    if result["ok"] != True:
        sys.exit(-1)
    time.sleep(2)


    print "delete subdomain abc.xyzzz.ol"
    raw_txn = delete_sub_domain(name, addrs[0])
    signed = sign(raw_txn, addrs[0])
    time.sleep(3)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    if result["ok"] != True:
        sys.exit(-1)
    time.sleep(2)


    print " send to a deleted subdomain (this should fail) "
    print name
    raw_txn = send_domain(name, addrs[0], "10")

    signed = sign(raw_txn, addrs[0])
    time.sleep(3)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

