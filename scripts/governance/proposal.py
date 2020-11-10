import sys, getopt
import time, datetime
from sdk import *
from sdkcom import *

_initial_funding_general = "5000" + "0"*18
_funding_goal_general = "15000" + "0" * 18

now = datetime.datetime.now()
propID = "propID_" + now.strftime("%Y-%m-%d %H:%M:%S")

def parse_params(argv):
    pid = ""
    funds = ""
    try:
      opts, args = getopt.getopt(argv,"cp:f:",["propid=","fund="])
    except getopt.GetoptError:
      print 'proposal.py -p <proposer> -f <funds>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-p", "--propid"):
         pid = arg
      if opt in ("-f", "--funds"):
         funds = int(arg) * 10 ** 18
    return pid, funds

def create_proposal(pid, proposer):
    proposal = Proposal(pid, "general", "test proposal general", "headline of general proposal", proposer, _initial_funding_general)
    proposal.keypath = path.join(fullnode, "keystore")
    proposal.send_create_prod()
    print "proposal id:", proposal.pid

def fund_proposal(pid, amount, funder):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.keypath = path.join(fullnode, "keystore")
    prop_fund.send_fund()

if __name__ == "__main__":
    # parse options
    pid, funds = parse_params(sys.argv[1:])
    print "proposal.py started"
    create = True if funds == "" else False

    # node account as proposer
    proposer = nodeAccount(fullnode_dev)

    # create or fund
    if create:
      create_proposal(propID, proposer)
    else:
      fund_proposal(pid, funds, proposer)
