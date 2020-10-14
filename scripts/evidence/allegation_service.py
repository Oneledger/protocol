from sdk import *


def create_allegation(reporterAccount, maliciousAccount):
    print reporterAccount
    print maliciousAccount
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


def query_requests():
    requests = ByzantineFault_Requests()
    print("No of requests in queue : " + str(len(requests)))
    return requests


def wait_1_block():
    height = GetBlockHeight()
    check_height = height + 1
    print("Waiting height %d to proceed (current: %s)" % (
        check_height, height,
    ))
    while check_height >= height:
        height = GetBlockHeight()
        time.sleep(1)
    print("Height %s ready" % check_height)

def main():
    reporterAccount, v2, v3, validators, maliciousAccount = setup()
    height = GetBlockHeight()
    check_height = 5
    print("Waiting height %d to proceed (current: %s)" % (
        check_height, height,
    ))
    while check_height >= height:
        height = GetBlockHeight()
        time.sleep(1)
    print("Height %s ready" % check_height)

    query_requests()

    numOfAllegationsPerform = 1
    for i in range(numOfAllegationsPerform):
        create_allegation(reporterAccount, maliciousAccount)

    wait_1_block()

    requests = query_requests()
    request_id = requests[0]['ID']
    for validator, keypath in validators:
        vote_allegation(validator, request_id, keypath)

    wait_1_block()
    release(maliciousAccount, node_1)


if __name__ == '__main__':
    main()

# olclient byzantine_fault allegation --address 963d5d0d501e877c3d2e4d96161cc73330524013 --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test
# olclient byzantine_fault allegation --address f290359332e2c2f31559c05c871730d277f28c23 --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test
# olclient byzantine_fault allegation --address 6ceaac9b49ea6dd2f1075ab4ffc3ab594cb899ec --maliciousAddress 6ea01e67198c80fab1a9b4e65dcc581d7d76fa7a --blockHeight 1 --password 1234 --proofMsg test


# Cannot freeze frozen Validator
# Frozen Validator cannot stake,unstake ,withdraw stake  -> error frozen validator
# Validator is frozen in block X , Stake and Active = false in block X +1
# Frozen validator cannot accumulate rewards  , checked rewards diff between frozen and active
# Frozen
# {
#     "address": "0lt05e5e4ea96d1a6fd36acdd7196fe738f38191c06",
#     "matureBalance": "84383561643835616437",
#     "pendingAmount": "0",
#     "totalAmount": "84383561643835616437",
#     "withdrawnAmount": "0"
# }
# Active
# {
#     "address": "0lt096a4dee464e74c8f0b667835ba04413998fd958",
#     "matureBalance": "128493150684931506835",
#     "pendingAmount": "9589041095890410958",
#     "totalAmount": "138082191780821917793",
#     "withdrawnAmount": "0"
# }
# Frozen Validator can withdraw block rewards
# Reelase tx does nto refun stake
