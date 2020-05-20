import sys

from sdk.actions import *
import time

proposals = ["proposal description A", "codeChange"]

def SendCreateProposal():
    addr_list = addresses()
    print addr_list[0]

    initialFunding = (int("10023450")*10**14)

    raw_txn = create_proposal(proposals[1], addr_list[0], proposals[0], initialFunding)

    signed = sign(raw_txn, addr_list[0])

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if not result["ok"]:
        sys.exit(-1)


if __name__ == "__main__":
    SendCreateProposal()
    time.sleep(2)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")
