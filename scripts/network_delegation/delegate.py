import sys, getopt
from sdk import *

def parse_params(argv):
    tx = ""
    delegator = ""
    amount = 0
    try:
      opts, args = getopt.getopt(argv,"duwra:m:",["address=","amount="])
    except getopt.GetoptError:
      print 'delegate.py -d -a <address> -m <amount>'
      print 'delegate.py -u -a <address> -m <amount>'
      print 'delegate.py -w -a <address> -m <amount>'
      print 'delegate.py -r -a <address> -m <amount>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-d", "--delegate"):
         tx = "delegate"
      if opt in ("-u", "--undelegate"):
         tx = "undelegate"
      if opt in ("-w", "--withdraw"):
         tx = "withdraw"
      if opt in ("-r", "--reinvest"):
         tx = "reinvest"

      if opt in ("-a", "--address"):
         delegator = arg
      if opt in ("-m", "--amount"):
         # in big OLT
         amount = int(arg)
    return tx, delegator, amount

if __name__ == "__main__":
    # parse options
    tx, delegator, amount = parse_params(sys.argv[1:])
    keypath = path.join(fullnode, "keystore")
    amt = str(amount) + '0' * 18

    if tx == "delegate":
        print "delegate, delegator=" + delegator + ", amt= " + amt
        tx = NetWorkDelegate(delegator, amt, keypath)
        tx.send_network_Delegate(exit_on_err=True, mode=TxCommit)
    elif tx == "undelegate":
        print "undelegate, delegator=" + delegator + ", amt= " + amt
        tx = NetWorkDelegate(delegator, amt, keypath)
        tx.send_network_undelegate(amt, exit_on_err=True, mode=TxCommit)
    elif tx == "withdraw":
        print "withdraw, delegator=" + delegator + ", amt= " + amt
        tx = WithdrawRewards(delegator, int(amt), keypath)
        tx.send(exit_on_err=True, mode=TxCommit)
    elif tx == "reinvest":
        print "reinvest, delegator=" + delegator + ", amt= " + amt
        tx = ReinvestRewards(delegator, keypath)
        tx.send(int(amt), exit_on_err=True, mode=TxCommit)
    else:
        result = query_delegation([delegator])
        print "----------------------------------"
        print "delegationStats: "
        print result[0]["delegationStats"]
        print "delegationRewardsStats: "
        print result[0]["delegationRewardsStats"]

        total = query_delegation_total(0)
        print "----------------------------------"
        print "activeAmount= " + total["activeAmount"]
        print "pendingAmount= " + total["pendingAmount"]
        print "totalAmount= " + total["totalAmount"]
        print "height= " + str(total["height"])
