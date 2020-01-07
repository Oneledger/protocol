from sdk.actions import *
import json
import sys
import time

if __name__ == "__main__":
    addrs = addresses()
    print addrs

    create_price = (int("10002345")*10**14)
    print "create price:", create_price


    name = "xyzzz.ol"
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

    name = "abc.xyzzz.ol"

    raw_txn = create_sub_domain(name, addrs[0], create_price, '')
    signed = sign(raw_txn, addrs[0])
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print json.dumps(result)
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)


    raw_txn = send_domain(name, addrs[0], "10")
    signed = sign(raw_txn, addrs[0])
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print
    if result["ok"] != True:
        sys.exit(-1)
    time.sleep(2)


    raw_txn = delete_sub_domain(name, addrs[0])
    signed = sign(raw_txn, addrs[0])

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    if result["ok"] != True:
        sys.exit(-1)
    time.sleep(2)


    print name
    raw_txn = send_domain(name, addrs[0], "10")

    signed = sign(raw_txn, addrs[0])
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

