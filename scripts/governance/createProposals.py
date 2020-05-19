from sdk.actions import *
import time

proposals = ["proposal description A", "codeChange"]

def SendCreateProposal():
    addrList = addresses()
    print addrList[0]

    initialFunding = (int("10023450")*10**14)

    raw_txn = create_proposal(proposals[1], addrList[0], proposals[0], initialFunding)

    signed = sign(raw_txn, addrList[0])

    result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])
    print result
    print "###################"
    print

    if result["ok"] != True:
        sys.exit(-1)

if __name__ == "__main__":
    SendCreateProposal()
    time.sleep(2)
    query_proposals("active")
    query_proposals("passed")
    query_proposals("failed")