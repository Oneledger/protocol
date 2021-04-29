import datetime
import time

from sdk import *

addr_list = addresses()
domain_name = "testowneraccept5.ol"
offerAmount = int("5")*10**18
counterOfferAmount = int("10")*10**18
counterBidAmount = int("7")*10**18
dt = datetime.datetime.now()
utc_time = dt.utcnow()
deadline = utc_time + datetime.timedelta(0, 60)
deadline_stamp = int((deadline - datetime.datetime(1970, 1, 1)).total_seconds())
bidConv = BidConv(addr_list[0], domain_name, 0x21, addr_list[1], offerAmount, counterOfferAmount, counterBidAmount, deadline_stamp)
if __name__ == "__main__":
    print "################ create domain ################"
    create_price = (int("1001") * 10 ** 18)
    raw_txn = create_domain(domain_name, addr_list[0], create_price)
    signed = sign(raw_txn, addr_list[0])
    time.sleep(1)
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

    time.sleep(5)
    print "bidder balance"
    query_balance(bidConv.bidder)
    print "owner balance"
    query_balance(bidConv.owner)
    bidConv.send_create()
    bidConvs = query_bidConvs(0x01, addr_list[0], domain_name, 0x21, addr_list[1])
    id = bidConvs["bidConvStats"][0]["bidConv"]["bidId"]
    print id

    time.sleep(5)
    print "bidder balance"
    query_balance(bidConv.bidder)
    print "owner balance"
    query_balance(bidConv.owner)
    bidConv.send_counter_offer(id)

    time.sleep(5)
    print "bidder balance"
    query_balance(bidConv.bidder)
    print "owner balance"
    query_balance(bidConv.owner)
    query_bidConvs(0x01, addr_list[0], domain_name, 0x21, addr_list[1])

    bidConv.send_bid_offer(id)

    time.sleep(5)
    print "bidder balance"
    query_balance(bidConv.bidder)
    print "owner balance"
    query_balance(bidConv.owner)
    query_bidConvs(0x01, addr_list[0], domain_name, 0x21, addr_list[1])

    bidConv.send_owner_decision(id, 1)

    time.sleep(5)
    print "bidder balance"
    query_balance(bidConv.bidder)
    print "owner balance"
    query_balance(bidConv.owner)
    result = query_bidConvs(0x02, addr_list[0], domain_name, 0x21, addr_list[1])
    if len(result["bidConvStats"]) != 1:
        sys.exit(-1)