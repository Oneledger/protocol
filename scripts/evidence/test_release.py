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
    KillNode,
)

from base import valDict, malicious


def test_release():
    print("Starting release test")
    # checking if validators is not frozen
    fMap = GetFrozenMap()
    assert len(fMap) == 0, 'Frozen validator found'

    is_killed = KillNode(node_1)
    assert is_killed is True, 'Failed to kill node %s' % node_1

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
    assert valDict['1']['address'] in fMap, 'Validator %s not frozen' % malicious

    print("Test for release successfully passed!")


if __name__ == "__main__":
    test_release()
