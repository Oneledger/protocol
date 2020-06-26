from rpc_call import *


def query_rewards(validator):
    req = {
        "validator": validator,
    }

    resp = rpc_call('query.ListRewardsForValidator', req)

    if "result" in resp:
        result = resp["result"]
    else:
        result = ""

    # print json.dumps(resp, indent=4)
    return result


def list_validators():
    resp = rpc_call('query.ListValidators', {})
    result = resp["result"]
    # print json.dumps(resp, indent=4)
    return result
