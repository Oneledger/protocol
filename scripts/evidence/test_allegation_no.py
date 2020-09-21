from __future__ import print_function

import subprocess
import os
import json
import time

from sdk.actions import (
    ListValidators,
    ByzantineFault_Requests,
    NodeID,
    GetBlockHeight,
    GetFrozenMap,
    GetActiveMap,
    ValidatorStatus,
)

from sdk.rpc_call import (
    node_0,
    node_1,
    node_2,
    node_3,
)

from sdk.cmd_call import (
    GetNodeCreds,
    Account_Add,
    ByzantineFault_Allegation,
    Send,
    GetNodeKey,
    ByzantineFault_Vote,
)

from base import *


def test_no_votes():
    print("Starting allegation no test")
    # check requests, so they will be empty
    requests = ByzantineFault_Requests()
    assert len(requests) == 0, 'Vote requests must be empty'

    # perform allegation
    is_allegation_pass = ByzantineFault_Allegation(node_0, reporter, malicious, 1, "test", "1234")
    assert is_allegation_pass is True, 'Failed to perform allegation'

    # check requests, so we will have 1 voting
    requests = ByzantineFault_Requests()
    assert len(requests) == 1, 'Vote requests must not be empty'
    request_id = requests[0]['ID']
    status = requests[0]['Status']
    assert status == 1, 'Vote must be pending'

    # storing init delegation amount
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
    assert len(requests) == 1, 'Vote requests must not be empty'
    status = requests[0]['Status']
    assert status == 2, 'Vote must be %s (got %s)' % (
        get_status_display(2),
        get_status_display(status)
    )
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
    try:
        set_up()
        test_no_votes()
    except AssertionError as e:
        raise e
