import sys

from sdk.actions import *
import time

initial_funding = (int("10023450") * 10 ** 14)
initial_funding_insufficient = (int("1000"))
proposals = [["proposal description A", "codeChange", "10001", initial_funding],
             ["proposal description B", "codeChange", "10002", initial_funding],
             ["proposal description C", "configUpdate", "10003", initial_funding],
             ["proposal description D", "configUpdate", "10004", initial_funding_insufficient],
             ["proposal description E", "general", "10005", initial_funding],
             ["proposal description F", "general", "10006", initial_funding]
             ]


def send_create_proposal(proposal):
    addr_list = addresses()
    print addr_list[0]

    raw_txn = create_proposal(proposal[2], proposal[1], addr_list[0], proposal[0], proposal[3])

    signed = sign(raw_txn, addr_list[0])

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print "###################"
    print

    if "ok" in result:
        if not result["ok"]:
            sys.exit(-1)


if __name__ == "__main__":

    for prop in proposals:
        send_create_proposal(prop)

    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")
