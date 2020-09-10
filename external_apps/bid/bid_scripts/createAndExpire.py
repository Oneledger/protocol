import datetime
import time

from sdk import *

addr_list = addresses()
domain_name = "testcreate.ol"
dt = datetime.datetime.now()
utc_time = dt.utcnow()
deadline = utc_time + datetime.timedelta(0, 20)
deadline_stamp = int((deadline - datetime.datetime(1970, 1, 1)).total_seconds())
bidConv = BidConv(addr_list[0], domain_name, 0x21, addr_list[1], 5, 10, 7, deadline_stamp)

if __name__ == "__main__":
    print "################ create domain ################"
    create_price = (int("10023450") * 10 ** 14)
    raw_txn = create_domain(domain_name, addr_list[0], create_price)
    signed = sign(raw_txn, addr_list[0])
    time.sleep(1)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    time.sleep(3)

    bidConv.send_create()

    time.sleep(3)
    query_bidConvs(0x01, addr_list[0], domain_name, 0x21, addr_list[1])

    time.sleep(15)
    query_bidConvs(0x04, addr_list[0], domain_name, 0x21, addr_list[1])
