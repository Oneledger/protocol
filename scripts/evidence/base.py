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
)


validators = ListValidators()

valDict = {data['name']: data for data in validators}

reporter = valDict['0']['address'][3:]
reporter_staking = valDict['0']['stakeAddress'][3:]
reporter_staking_key = GetNodeKey('0')
malicious = valDict['1']['address'][3:]

voter_2 = valDict['2']['address'][3:]
voter_3 = valDict['3']['address'][3:]

reporter_creads = GetNodeCreds('0')
voter_2_creds = GetNodeCreds('2')
voter_3_creds = GetNodeCreds('3')

statuses = {
    1: 'Voting',
    2: 'Innocent',
    3: 'Guilty',
}


def get_status_display(status):
    return statuses[status]


def set_up():
    # adding accouts for validators
    is_added = Account_Add(node_0, reporter_creads['pub'], reporter_creads['priv'], '1234')
    assert is_added is True, 'Failed to add account for %s' % reporter_creads['address']

    is_added = Account_Add(node_2, voter_2_creds['pub'], voter_2_creds['priv'], '1234')
    assert is_added is True, 'Failed to add account for %s' % voter_2_creds['address']

    is_added = Account_Add(node_3, voter_3_creds['pub'], voter_3_creds['priv'], '1234')
    assert is_added is True, 'Failed to add account for %s' % voter_3_creds['address']

    staking_pub_key = NodeID()
    is_added = Account_Add(node_0, staking_pub_key, reporter_staking_key['priv'], '1234')
    assert is_added is True, 'Failed to add staking account for %s' % reporter_creads['address']

    print("Accounts for nodes 0, 2, 3 and staking were created!")

    is_sent = Send(node_0, reporter_staking, reporter, 10, '1234', currency='OLT', fee='0.0001')
    assert is_sent is True, 'Failed to send on %s' % reporter

    is_sent = Send(node_0, reporter_staking, voter_2, 10, '1234', currency='OLT', fee='0.0001')
    assert is_sent is True, 'Failed to send on %s' % voter_2

    is_sent = Send(node_0, reporter_staking, voter_3, 10, '1234', currency='OLT', fee='0.0001')
    assert is_sent is True, 'Failed to send on %s' % voter_3

    print("Validator balances for 0, 2, 3 ready!")
