import sys

from sdk.actions import *
import time

initial_funding = (int("10023450") * 10 ** 14)
proposal = ["id_20000", "general", "proposal description", initial_funding]
addr_list = addresses()

def send_create_proposal(proposal):
    print addr_list[0]

    # create Tx
    raw_txn = create_proposal(proposal[0], proposal[1], addr_list[0], proposal[2], proposal[3])

    # sign Tx
    signed = sign(raw_txn, addr_list[0])

    # broadcast Tx
    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

    if "ok" not in result or not result["ok"]:
        sys.exit(-1)
    print "################### proposal created:" + proposal

def send_vote_proposal(proposal, opinion, url):
    # create Tx
    raw_txn, signed, signer = vote_proposal(proposal, opinion, url)

    # broadcast Tx
    result = broadcast_commit(signed["rawTx"], signed['signature']['Signed'], signed['signature']['Signer'])

    if "ok" not in result or not result["ok"]:
        sys.exit(-1)
    print "################### proposal voted:" + proposal + "opinion: " + opinion

if __name__ == "__main__":
    send_create_proposal(proposal)

    time.sleep(5)

    print "#### ACTIVE PROPOSALS: ####"
    query_proposals("active")

    print "#### PASSED PROPOSALS: ####"
    query_proposals("passed")

    print "#### FAILED PROPOSALS: ####"
    query_proposals("failed")

    
