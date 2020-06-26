from time import sleep
from sdk import *


def query_block_rewards():
    try:
        while True:
            response = list_validators()
            print '-------------------------------- Height: ', response["height"], ' --------------------------------'
            for v in response["validators"]:
                print query_rewards(v["address"])
            sleep(5)
    except KeyboardInterrupt:
        print 'Exiting Test'


if __name__ == "__main__":
    query_block_rewards()
