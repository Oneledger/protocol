import datetime
import time

from sdk import *

addr_list = addresses()
asset_name = "testb"
offerAmount = int("1")*10**18
counterOfferAmount = int("10")*10**18
counterBidAmount = int("7")*10**18
dt = datetime.datetime.now()

if __name__ == "__main__":
    for i in range(0, 10000):
        utc_time = dt.utcnow()
        deadline = utc_time + datetime.timedelta(0, 20)
        deadline_stamp = int((deadline - datetime.datetime(1970, 1, 1)).total_seconds())
        bidConv = BidConv(addr_list[1], asset_name + str(i), 0x22, addr_list[0], offerAmount, counterOfferAmount, counterBidAmount, deadline_stamp)
        print i
        print bidConv.send_create_async()

    # query_bidConv("ff0d46b86480769354dad41552269a5240a08bf451ea93d0952ec1a180fb11ef")

