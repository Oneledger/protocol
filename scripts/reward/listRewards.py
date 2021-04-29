from time import sleep
from sdk import *


def query_block_rewards(count=1):
    try:
        while count > 0:
            response = list_validators()
            print '-------------------------------- Height: ', response["height"], ' --------------------------------'
            for v in response["validators"]:
                print query_rewards(v["address"])
                print "Validator Reward Stats: ", query_all_rewards(v["address"])
            sleep(5)
            count -= 1
    except KeyboardInterrupt:
        print 'Exiting Test'

if __name__ == "__main__":
    query_block_rewards()
