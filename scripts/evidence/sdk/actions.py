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
