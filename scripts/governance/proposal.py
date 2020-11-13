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
    amount = ""
    tx = "create"
    try:
      opts, args = getopt.getopt(argv,"fcwp:m:",["propid=","amount="])
    except getopt.GetoptError:
      print 'proposal.py -p <proposal> -m <amount>'
      print 'proposal.py -p <proposal> -c'
      print 'proposal.py -p <proposal> -m <amount> -w'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-p", "--propid"):
         pid = arg
      if opt in ("-m", "--amount"):
         amount = int(arg) * 10 ** 18
      if opt in ("-w", "--withdraw"):
         tx = "withdraw"
      if opt in ("-c", "--cancel"):
         tx = "cancel"
      if opt in ("-f", "--fund"):
         tx = "fund"
    return pid, amount, tx

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

def withdraw_fund(pid, funder, amount, beneficiary):
    fund_withdraw = ProposalFundsWithdraw(pid, funder, amount, beneficiary)
    fund_withdraw.withdraw_fund(funder)
    time.sleep(2)

def cancel_proposal(pid, proposer, reason, secs=1):
    prop_cancel = ProposalCancel(pid, proposer, reason)
    res = prop_cancel.send_cancel()
    time.sleep(secs)
    return res

if __name__ == "__main__":
    # parse options
    pid, amount, tx = parse_params(sys.argv[1:])

    # node account as proposer
    proposer = nodeAccount(fullnode)

    # create or fund or withdraw or cancel
    if tx == "create":
      create_proposal(propID, proposer)
    elif tx == "fund":
      print pid, amount, proposer
      fund_proposal(pid, amount, proposer)
    elif tx == "withdraw":
      withdraw_fund(pid, proposer, amount, proposer)
    elif tx == "cancel":
      cancel_proposal(pid, proposer, "no reason")
