import toml
import os



def setup_config(usr, privKey):

    print('Setup Config file for Node 1 ...')
    conf = os.environ['OLDATA'] + "/devnet/1-Node/config.toml"
    Configuration = toml.load(conf)

    node_conf = Configuration['node']
    auth_conf = node_conf['Auth']

    if usr:
        auth_conf['owner_credentials'] = ['username:password']
    else:
        auth_conf['owner_credentials'] = None

    if privKey:
        auth_conf['rpc_private_key'] = \
            'pSKuldvwewRtuUcItfNhCgFy+RTscEnUejF2YRtnqvl98z17rUJLebNRvVwGSO0v3PFGhfng/CvUSru+qYD5Dw=='
    else:
        auth_conf['rpc_private_key'] = ""

    f = open(conf, 'w')
    toml.dump(Configuration, f)



if __name__ == "__main__":

    print('rpcAuth Test Script running')

    print('************ Test Get Token API ************')

    setup_config(True, True)
