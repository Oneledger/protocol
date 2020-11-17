import sys, getopt
from sdkcom import *

def parse_params(argv):
    amount = 0
    try:
      opts, args = getopt.getopt(argv,"f:",["funds="])
    except getopt.GetoptError:
      print 'account.py -f <funds>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-f", "--funds"):
         # in big OLT
         amount = int(arg)
    return amount

if __name__ == "__main__":
    # parse options
    amount = parse_params(sys.argv[1:])

    # create account
    funder = nodeAccount(fullnode)
    account = createAccount(fullnode, amount, funder)
    print account
