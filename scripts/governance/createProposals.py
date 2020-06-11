import sys
import time
from sdk import *

addr_list = addresses()

_pid = "id_20000"
_proposer = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_initial_funding_insufficient = (int("1") * 10 ** 8)
_initial_funding_too_much = (int("100") * 10 ** 9)
_funding_goal = (int("10") * 10 ** 9)

proposals = [Proposal("proposal description A", "codeChange", "10001", "headline A", _proposer, _initial_funding),
             Proposal("proposal description B", "codeChange", "10002", "headline B",  _proposer, _initial_funding),
             Proposal("proposal description C", "configUpdate", "10003", "headline C",  _proposer, _initial_funding),
             Proposal("proposal description E", "general", "10005", "headline E", _proposer, _initial_funding),
             Proposal("proposal description F", "general", "10006", "headline F", _proposer, _initial_funding)
             ]
proposal_littleInitFund = Proposal("proposal description D", "configUpdate", "10004", "headline D", _proposer, _initial_funding_insufficient)
proposal_hugeInitFund = Proposal("proposal description G", "general", "10007", "headline G", _proposer, _initial_funding_too_much)

if __name__ == "__main__":
    # create normal proposals
    for prop in proposals:
        prop.send_create()
        print "proposal id:", prop.pid

    # create proposal with little initial fund
    proposal_littleInitFund.send_create()

    # create proposal with huge initial fund
    proposal_hugeInitFund.send_create()

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

    print bcolors.OKGREEN + "#### Test create proposals succeed" + bcolors.ENDC
    print ""
