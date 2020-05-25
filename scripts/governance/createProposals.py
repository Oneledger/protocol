import sys
import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_20000"
_proposer = addr_list[0]
_initial_funding = (int("10023450") * 10 ** 14)
_initial_funding_insufficient = (int("1000"))

proposals = [Proposal("proposal description A", "codeChange", "10001", _proposer, _initial_funding),
             Proposal("proposal description B", "codeChange", "10002", _proposer, _initial_funding),
             Proposal("proposal description C", "configUpdate", "10003", _proposer, _initial_funding),
             Proposal("proposal description E", "general", "10005", _proposer, _initial_funding),
             Proposal("proposal description F", "general", "10006", _proposer, _initial_funding)
             ]
proposal_littleInitFund = Proposal("proposal description D", "configUpdate", "10004", _proposer, _initial_funding_insufficient)

if __name__ == "__main__":
    # create normal proposals
    for prop in proposals:
        prop.send_create()

    #create proposal with little initial fund
    proposal_littleInitFund.send_create()

    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    activeList = query_proposals("active")
    if len(activeList) != len(proposals):
        sys.exit(-1)

    print "#### PASSED PROPOSALS: ####"
    passedList = query_proposals("passed")
    if len(passedList) != 0:
        sys.exit(-1)

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals("failed")
    if len(failedList) != 0:
        sys.exit(-1)
