from sdk import *


def create_allegation(reporterAccount, maliciousAccount):
    newAllegation = Byzantine(reporterAccount, maliciousAccount, "test", "1234", 1, node_0 + "/keystore/")
    newAllegation.send_allegation()


def vote_allegation(validator, requestId, keypath):
    vote = 1
    newVote = Vote(validator, requestId, vote, keypath + "/keystore/")
    newVote.send_vote()


def release(owneraddress, keypath):
    newrelease = Release(owneraddress, keypath + "/keystore/")
    newrelease.send_release()


def setup():
    reporterAccount = addOwnerAccount(node_0)
    v2 = addOwnerAccount(node_2)
    v3 = addOwnerAccount(node_3)
    addValidatorWalletAccounts(node_1)
    validators = [(reporterAccount, node_0), (v2, node_2), (v3, node_3)]
    maliciousAccount = addOwnerAccount(node_1)
    # print "olclient byzantine_fault allegation --address " + reporterAccount + " --maliciousAddress " + maliciousAccount + " --blockHeight 1 --password 1234 --proofMsg test"
    # print "olclient byzantine_fault allegation --address " + v2 + " --maliciousAddress " + maliciousAccount + " --blockHeight 1 --password 1234 --proofMsg test"
    # print "olclient byzantine_fault allegation --address " + v3 + " --maliciousAddress " + maliciousAccount + " --blockHeight 1 --password 1234 --proofMsg test"
    return reporterAccount, v2, v3, validators, maliciousAccount


def query_requests(req=1):
    requests = ByzantineFault_Requests()
    print("No of requests in queue : " + str(len(requests)))
    assert len(requests) <= req, 'Vote requests must not be more than ' + str(req)
    return requests


def allegations(numOfAllegationsPerform, reporterAccount, maliciousAccount):
    for i in range(numOfAllegationsPerform):
        create_allegation(reporterAccount, maliciousAccount)


def voting(validators):
    requests = query_requests()
    request_id = requests[0]['ID']
    for validator, keypath in validators:
        vote_allegation(validator, request_id, keypath)


def main():

    reporterAccount, v2, v3, validators, maliciousAccount = setup()
    numOfAllegationsPerform = 1
    wait_blocks(5)
    query_requests()
    # Test case 1 : Multiple requests against same validator in different blocks
    for i in range(1):
        allegations(numOfAllegationsPerform, reporterAccount, maliciousAccount)
        wait_blocks(80)
        voting(validators)
        wait_blocks(1)
        # release(maliciousAccount, node_1)
        # wait_blocks(1)
    # Test case 1 : Multiple requests different Validators in same and different blocks
    # for i in range(3):
    #     allegations(numOfAllegationsPerform, reporterAccount, maliciousAccount)
    #     allegations(numOfAllegationsPerform, reporterAccount, v2)
    #     allegations(numOfAllegationsPerform, reporterAccount, v3)
    #     wait_blocks(1)
    # wait_blocks(1)
    # query_requests(3)


if __name__ == '__main__':
    main()

# olclient byzantine_fault allegation --address 963d5d0d501e877c3d2e4d96161cc73330524013 --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test
# olclient byzantine_fault allegation --address f290359332e2c2f31559c05c871730d277f28c23 --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test
# olclient byzantine_fault allegation --address 6ceaac9b49ea6dd2f1075ab4ffc3ab594cb899ec --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test
