from __future__ import print_function

import time

from base import *
from sdk.actions import (
    ByzantineFault_Requests,
    GetBlockHeight,
    GetFrozenMap,
    GetActiveMap,
    ValidatorStatus, wait_blocks,
)
from sdk.cmd_call import (
    ByzantineFault_Allegation,
    ByzantineFault_Vote,
)


def test_no_votes():
    print("Starting allegation no test")
    # check requests, so they will be empty
    requests = ByzantineFault_Requests()
    assert len(requests) == 0, 'Vote requests must be empty'

    # perform allegation
    is_allegation_pass = ByzantineFault_Allegation(node_0, reporter, malicious, 1, "test", "1234")
    assert is_allegation_pass is True, 'Failed to perform allegation'

    # check requests, so we will have 1 voting
    wait_blocks(1)
    requests = ByzantineFault_Requests()
    assert len(requests) == 1, 'Vote requests must not be empty'
    request_id = requests[0]['ID']
    status = requests[0]['Status']
    assert status == 1, 'Vote must be pending'

    # storing init delegation amount
    wait_blocks(1)
    validator_status = ValidatorStatus(malicious)
    power = validator_status['power']
    total_amount = int(validator_status['totalDelegationAmount'])

    # 0 percent will be cutted down (from genesis)
    expected_power = int(power)
    expected_total_amount = total_amount

    vote = 'no'

    # vote for no (50% so 2 votes)
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

    # freezeing height and waiting 1 blocks
    height = GetBlockHeight()
    check_height = height + 1
    print("Waiting height %d to proceed (current: %s)" % (
        check_height, height,
    ))
    while check_height >= height:
        height = GetBlockHeight()
        time.sleep(1)
    print("Height %s ready" % check_height)

    # check requests, so we will have 1 innocent
    print("Checking request status...")
    requests = ByzantineFault_Requests()
    assert len(requests) == 0, 'Vote requests must be empty ,request has been cleared'
    # status = requests[0]['Status']
    # assert status == 2, 'Vote must be %s (got %s)' % (
    #     get_status_display(2),
    #     get_status_display(status)
    # )
    print("Checking request done.")

    # freezeing height and waiting 4 blocks
    height = GetBlockHeight()
    check_height = height + 4
    print("Waiting height %d to proceed (current: %s)" % (
        check_height, height,
    ))
    while check_height >= height:
        height = GetBlockHeight()
        time.sleep(1)
    print("Height %s ready" % check_height)

    # checking if validator is frozen
    fMap = GetFrozenMap()
    vMap = GetActiveMap()
    assert len(fMap) == 0, 'Frozen validator found'
    assert len(vMap) == len(validators), 'Not all validators are active!'

    print("Checking power and amounts...")
    validator_status = ValidatorStatus(malicious)
    assert validator_status['power'] == expected_power, 'Power not match (expected: %d, got: %d)' % (
        expected_power,
        validator_status['power'],
    )
    assert int(validator_status['totalDelegationAmount']) == expected_total_amount, 'Amount not match (expected: %d, got: %d)' % (
        expected_total_amount,
        int(validator_status['totalDelegationAmount']),
    )
    print("Test for allegation no successfully passed!")


if __name__ == "__main__":
    set_up()
    test_no_votes()
