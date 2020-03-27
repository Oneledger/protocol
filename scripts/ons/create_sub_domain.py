"""
Create Sub Domain

1. Create Domain using existing script.
2. Create Sub Domain based on initial Parent domain.
3. Send Currency to the sub domain.

"""

import sys
import time
import  calendar
import uuid
from sdk.actions import *

if __name__ == "__main__":
    addrs = addresses()

    """
        ****** Create Initial Domain ******
    """

    # Prepare new domain for creation
    name = "alice2" + str(calendar.timegm(time.gmtime())) + uuid.uuid1().hex +".ol"
   # print name
   # name = "alice2.ol"
    create_price = (int("10002345") * 10 ** 14)
   # print "create price:", create_price

    # Get raw transaction for domain creation
    raw_txn = create_domain(name, addrs[0], create_price)
    #print "raw create domain tx:", raw_txn

    # Sign raw Transaction
    signed = sign(raw_txn, addrs[0])
    #print "signed create domain tx:", signed
    #print

    # Broadcast transaction
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

    #print "###################"
    #print

    if not result["ok"]:
     #   print result
        sys.exit(-1)

    """
        ****** Create Sub domain based on initial domain above ******
    """
    #print "---Creating Sub Domain---"

    # Prepare sub domain
    sub_name = "bob." + name
    # Use same create price as above

    # Get raw transaction
    raw_txn = create_sub_domain(sub_name, addrs[0], create_price, "http://myuri.com")
    #print "raw create sub domain transaction: ", raw_txn

    # Sign Transaction
    signed = sign(raw_txn, addrs[0])
    #print "signed create sub domain transaction: ", signed
    #print

    # Broadcast Transaction
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

    #print "###################"
    #print

    if not result["ok"]:
        print result
        sys.exit(-1)

    """
       ****** Send to sub domain on initial domain above ******
   """
   # print "---Send to Domain---"
    raw_txn = send_domain(name, addrs[0], "10")
   # print raw_txn

    signed = sign(raw_txn, addrs[0])
    #print signed
    #print

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

    #print "###################"
    #print

    if not result["ok"]:
        print result
        sys.exit(-1)



    #print_all_domains(addrs[0])
