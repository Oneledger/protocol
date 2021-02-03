from rpc_call import *
from actions import add_account


class Account:
    def __init__(self, name, password, pubkey, privatekey, user_id=""):
        self.name = name
        self.password = password
        self.pubkey = pubkey
        self.privkey = privatekey
        self.user_id = user_id

# super admins
kevin = Account("kevin", "kevin", "RnYwf7S92kZDApLeCiHCksdmqwrbmVXkAJq95puF6m4=",
                "Bqv5vB+Kf/fsnRZ7oV+3HnBbP4BLqYdXDNDCPVzOsRFGdjB/tL3aRkMCkt4KIcKSx2arCtuZVeQAmr3mm4Xqbg==",
                '3d642f9b15ba58ca8f669fe83a0aff3a2bc71af1fcc0021e46916f3f499982a8')
charlie = Account("charlie", "charlie", "Dd7r+MP6OSac3eVzlTc6UQ5eSZm/z03MpYPap2OTUho=",
                  "gqPLq5/rB7GI34pMy9P48eJngAF2xs51RqkldC9HUK8N3uv4w/o5Jpzd5XOVNzpRDl5Jmb/PTcylg9qnY5NSGg==",
                  'c6befab40db15c0ef70f61d032dc50e3458e5b2ff466030db74b37f50cd60adc')
tanmay = Account("tanmay", "tanmay", "jqS9naVA6/O27PlV4rICS/IU3ZV69NDn5jXawcMrNf0=",
                 "5PoC0Vkdb5jzaY4BROoGgl2LXy4g4iatxIC+mPqPGE6OpL2dpUDr87bs+VXisgJL8hTdlXr00OfmNdrBwys1/Q==",
                 '9570519d458e70e6ffd70c82e8b40652b985d0c5cb76332b865f767da68a9107')
hao = Account("hao", "hao", "Yl2fYPB8d5di2is7BhxfdiZGTcyUXnX2mF3vX+jCoLc=",
              "xlAUYiO9C7J0HQaJYpRTOM2LqbRKM/D9viq1o/E+AmdiXZ9g8Hx3l2LaKzsGHF92JkZNzJRedfaYXe9f6MKgtw==",
              '9a47f317a44532e6ecffc54625420cc6318c1647b21e06adc18aa0628bf2fe6e')

# hospital admins
admin1 = Account("admin1", "admin1", "2TZNBEzI4WuxVqMdfbCbBP0s7tECk/GD4FNccdk8lI8=",
              "4Y4/ngP62E85a84w3VQ6I+vEZFtay1RLZpQlpa8yCcPZNk0ETMjha7FWox19sJsE/Szu0QKT8YPgU1xx2TyUjw==")
admin2 = Account("admin2", "admin2", "fSwegoOLWv/05tzgHuXKwNn/WI1rhVNi7Hhdg4MRnCQ=",
              "vuhchLRHnQTeeyH8AyUKTHxU/D7mWcKTrauTHJ2sn2N9LB6Cg4ta//Tm3OAe5crA2f9YjWuFU2LseF2DgxGcJA==")

# persons
person1 = Account("person1", "person1", "pRCAoJZq/Qp83JShrbrnaDJ/6C9gWC7Bs/GOrrIVSvA=",
              "2yPU/GIrL7IBxIGFeq8ytdaemR6ruxGSdFgA+FpU8YqlEICglmr9CnzclKGtuudoMn/oL2BYLsGz8Y6ushVK8A==")
person2 = Account("person2", "person2", "EAivztxUuR+ox0Ol7alKEs8YVKMhClHtVQk8Zl17hCU=",
              "eTG2YOn0SIVDIoV4Ect1tJQlV3BsovNgHgecCHY0g0IQCK/O3FS5H6jHQ6XtqUoSzxhUoyEKUe1VCTxmXXuEJQ==")
person3 = Account("person3", "person3", "D/zti3OPewgDoCKgYq4yvOJBhVVPCpwkCZBjcY84F9Q=",
              "2rAnP54qkuNyWHZUIcsgyPF87HEQZNV/e9j7QdQhM84P/O2Lc497CAOgIqBirjK84kGFVU8KnCQJkGNxjzgX1A==")

accounts = [kevin, charlie, tanmay, hao, admin1, admin2, person1, person2, person3]

super_user_addresses = {'kevin': '0lt5e18ca7b173e672e3cc6aca9790965af82e9c0ac',
                        'charlie': '0lt9e6891c2cc3de5fea4147f8f3ea813e8a19a5bf6',
                        'tanmay': '0lt6c850a34f35318ba0f71a8c48788d9ce0c000af0',
                        'hao': '0lt5fbdc9c5864b5d03e3a70807dc652cc6cbcfff0f'}
admin_addresses = {'admin1': '0lta0ca4068e341272e95be8039db6d8c1d060f573d',
                   'admin2': '0ltaf2ef4ea0dd63885c7bf57a8c8e03777ca5d879d'}
addrs = {'person1': '0lt3adbe60a334a7051aeb3fa36c8db60bd21118136',
         'person2': '0ltf73e1f21dc3fe1084e15f8bdc7624d0ae229cd95',
         'person3': '0lt4a546fdfeb28b8e3c6225093f9b966703e3a9a04'}

node_accounts = ['acct1', 'acct2', 'acct3', 'acct4', 'acct5', 'acct6']

def setup_accounts():
    for acct in accounts:
        os.system('olclient account add' + ' --name ' + acct.name + ' --password ' + acct.password +
                  ' --privkey ' + acct.privkey + ' --pubkey ' + acct.pubkey + ' --root ' + secure_accounts_path)

    for acc in node_accounts:
        add_account(acc)


''' ****************************************  SUPER USER ACCOUNT INFO ***********************************************
0lt5e18ca7b173e672e3cc6aca9790965af82e9c0ac
{
 "type": 0,
 "name": "kevin",
 "publicKey": {
  "keyType": "ed25519",
  "data": "RnYwf7S92kZDApLeCiHCksdmqwrbmVXkAJq95puF6m4="
 },
 "privateKey": {
  "keyType": "ed25519",
  "data": "Bqv5vB+Kf/fsnRZ7oV+3HnBbP4BLqYdXDNDCPVzOsRFGdjB/tL3aRkMCkt4KIcKSx2arCtuZVeQAmr3mm4Xqbg=="
 }
}


0lt9e6891c2cc3de5fea4147f8f3ea813e8a19a5bf6
{
 "type": 0,
 "name": "charlie",
 "publicKey": {
  "keyType": "ed25519",
  "data": "Dd7r+MP6OSac3eVzlTc6UQ5eSZm/z03MpYPap2OTUho="
 },
 "privateKey": {
  "keyType": "ed25519",
  "data": "gqPLq5/rB7GI34pMy9P48eJngAF2xs51RqkldC9HUK8N3uv4w/o5Jpzd5XOVNzpRDl5Jmb/PTcylg9qnY5NSGg=="
 }
}


0lt6c850a34f35318ba0f71a8c48788d9ce0c000af0
{
 "type": 0,
 "name": "tanmay",
 "publicKey": {
  "keyType": "ed25519",
  "data": "jqS9naVA6/O27PlV4rICS/IU3ZV69NDn5jXawcMrNf0="
 },
 "privateKey": {
  "keyType": "ed25519",
  "data": "5PoC0Vkdb5jzaY4BROoGgl2LXy4g4iatxIC+mPqPGE6OpL2dpUDr87bs+VXisgJL8hTdlXr00OfmNdrBwys1/Q=="
 }
}


0lt5fbdc9c5864b5d03e3a70807dc652cc6cbcfff0f
{
 "type": 0,
 "name": "hao",
 "publicKey": {
  "keyType": "ed25519",
  "data": "Yl2fYPB8d5di2is7BhxfdiZGTcyUXnX2mF3vX+jCoLc="
 },
 "privateKey": {
  "keyType": "ed25519",
  "data": "xlAUYiO9C7J0HQaJYpRTOM2LqbRKM/D9viq1o/E+AmdiXZ9g8Hx3l2LaKzsGHF92JkZNzJRedfaYXe9f6MKgtw=="
 }
}
'''