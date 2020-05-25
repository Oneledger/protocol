import sys
import time

from sdk.actions import *

addr_list = addresses()

_pid = "id_20000"
_proposer = addr_list[0]
_initial_funding = (int("10023450") * 10 ** 14)
_initial_funding_insufficient = (int("1000"))

proposals = [Proposal("proposal description A", "codeChange", "10001", _initial_funding),
             Proposal("proposal description B", "codeChange", "10002", _initial_funding),
             Proposal("proposal description C", "configUpdate", "10003", _initial_funding),
             Proposal("proposal description D", "configUpdate", "10004", _initial_funding_insufficient),
             Proposal("proposal description E", "general", "10005", _initial_funding),
             Proposal("proposal description F", "general", "10006", _initial_funding)
             ]

if __name__ == "__main__":
    # create proposals
    for prop in proposals:
        prop.send_create()

    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")
