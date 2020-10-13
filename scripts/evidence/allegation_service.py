from base import *
from sdk import *


def create_allegation(reporterAccount, maliciousAccount):
    newAllegation = Byzantine(reporterAccount, maliciousAccount, "test", "1234", 1, node_0 + "/keystore/")
    newAllegation.send_allegation()


def vote_allegation(validator, requestId, keypath):
    vote = 1
    newVote = Vote(validator, requestId, vote, keypath + "/keystore/")
    newVote.send_vote()


def main():
    numOfAllegationsPerform = 1
    # reporterAccount = addOwnerAccount(node_0)
    # v2 = addOwnerAccount(node_2)
    # v3 = addOwnerAccount(node_3)
    # validators = [(reporterAccount, node_0), (v2, node_2), (v3, node_3)]
    # maliciousAccount = addOwnerAccount(node_1)
    # print reporterAccount +"   "+maliciousAccount
    # print v2 +"   "+maliciousAccount
    # sys.exit(-1)
    # for i in range(numOfAllegationsPerform):
    #     create_allegation(reporterAccount, maliciousAccount)
    requests = ByzantineFault_Requests()
    assert len(requests) >= 1, 'Vote requests must not be empty'
    # assert len(requests) == 1, 'Only 1 vote request must exist'
    print("No of requests in queue : " + str(len(requests)))
    request_id = requests[0]['ID']
    for validator, keypath in validators:
        vote_allegation(validator, request_id, keypath)
    request_id = requests[1]['ID']
    for validator, keypath in validators:
        vote_allegation(validator, request_id, keypath)


if __name__ == '__main__':
    main()

# olclient byzantine_fault allegation --address e704c38bbc5fdff57e5b4c703b9d5da08133415a --maliciousAddress 131bbb3a7eb7a311eea050c6c73cf8da57cfa617 --blockHeight 1 --password 1234 --proofMsg test
# olclient byzantine_fault allegation --address e02667569dc4a8cae0c05f117da807f595979904 --maliciousAddress 131bbb3a7eb7a311eea050c6c73cf8da57cfa617 --blockHeight 1 --password 1234 --proofMsg test
