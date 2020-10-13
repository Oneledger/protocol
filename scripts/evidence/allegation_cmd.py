from __future__ import print_function

from base import *
from sdk.actions import (
    ByzantineFault_Requests,
)
from sdk.cmd_call import (
    ByzantineFault_Allegation,
    ByzantineFault_Vote,
)


def test_allegation_requests():
    print("Starting allegation yes test")
    # check requests, so they will be empty
    requests = ByzantineFault_Requests()

    # perform allegation
    is_allegation_pass = ByzantineFault_Allegation(node_0, reporter, malicious, 1, "test", "1234")
    assert is_allegation_pass is True, 'Failed to perform allegation'

    # check requests, so we will have 1 voting


def test_voting(index):
    requests = ByzantineFault_Requests()
    print(len(requests))
    assert len(requests) >= 1, 'Vote requests must not be empty'
    request_id = requests[index]['ID']
    status = requests[index]['Status']
    assert status == 1, 'Vote must be pending'
    vote = 'yes'

    # vote for yes (50% so 2 votes)
    has_voted = ByzantineFault_Vote(node_2, request_id, voter_2, vote, '1234')
    assert has_voted is True, 'Failed to perform vote on request: %s of address: %s' % (
        request_id,
        voter_2,
    )
    print("Vote for address %s on request %s with %s done!" % (voter_2, request_id, vote))

    has_voted = ByzantineFault_Vote(node_3, request_id, voter_3, vote, '1234')
    assert has_voted is True, 'Failed to perform vote on request: %s of address: %s' % (
        request_id,
        voter_3,
    )
    print("Vote for address %s on request %s with %s done!" % (voter_3, request_id, vote))


if __name__ == "__main__":
    numOfAllegationsPerform = 2
    # try:
    set_up()
    for i in range(numOfAllegationsPerform):
        test_allegation_requests()
    # time.sleep(10)
    # for i in range(numOfAllegationsPerform):
    #     test_voting(i)
    # except AssertionError as e:
    #     import ipdb;
    #
    #     ipdb.set_trace()
    #     raise e
