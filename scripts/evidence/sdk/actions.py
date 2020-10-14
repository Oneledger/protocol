import time

from rpc_call import rpc_call


def ListValidators():
    resp = rpc_call('query.ListValidators', {})
    return resp['result']['validators']


def GetBlockHeight():
    resp = rpc_call('query.ListValidators', {})
    return resp['result']['height']


def GetFrozenMap():
    resp = rpc_call('query.ListValidators', {})
    return resp['result'].get('fmap', {})


def GetActiveMap():
    resp = rpc_call('query.ListValidators', {})
    return resp['result'].get('vmap', {})


def ByzantineFault_Requests():
    resp = rpc_call('query.VoteRequests', {})
    return resp['result']['Requests']


def NodeID():
    resp = rpc_call('node.ID', {})
    return resp['result']['publicKey']


def ValidatorStatus(address):
    resp = rpc_call('query.ValidatorStatus', {
        'address': address,
    })
    return resp['result']


def wait_blocks(n):
    height = GetBlockHeight()
    check_height = height + n
    print("Waiting height %d to proceed (current: %s)" % (
        check_height, height,
    ))
    while check_height >= height:
        height = GetBlockHeight()
        time.sleep(1)
    print("Height %s ready" % check_height)
