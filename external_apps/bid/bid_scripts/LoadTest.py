import datetime
import time

from sdk import *

addr_list = addresses()
asset_name = "testd"
offerAmount = int("1")*10**18
counterOfferAmount = int("10")*10**18
counterBidAmount = int("7")*10**18
dt = datetime.datetime.now()

if __name__ == "__main__":
    for i in range(0, 10000):
        # time.sleep(0.1)
        utc_time = dt.utcnow()
        deadline = utc_time + datetime.timedelta(0, 20)
        deadline_stamp = int((deadline - datetime.datetime(1970, 1, 1)).total_seconds())
        bidConv = BidConv(addr_list[1], asset_name + str(i), 0x22, addr_list[0], offerAmount, counterOfferAmount, counterBidAmount, deadline_stamp)
        print i
        print bidConv.send_create_async()


    # query_bidConv("ffe0616a546fd587962d4085155661bb8663abc5e2706d2b46a2c03818d0086c")
    # result = query_bidConvs(0x04, addr_list[1], "", 0x22, addr_list[0])
    # print len(result["bidConvStats"])


