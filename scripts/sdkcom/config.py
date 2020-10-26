import os

url_tmTx = "http://127.0.0.1:26600/tx"
url_0 = "http://127.0.0.1:26602/jsonrpc"
url_1 = "http://127.0.0.1:26605/jsonrpc"
url_2 = "http://127.0.0.1:26608/jsonrpc"
url_3 = "http://127.0.0.1:26611/jsonrpc"
url_4 = "http://127.0.0.1:26614/jsonrpc"
url_5 = "http://127.0.0.1:26617/jsonrpc"
url_fullnode = ""

oltest = os.getenv('OLTEST')
oldata = os.environ['OLDATA']
devnet = os.path.join(oldata, "devnet")
loadtest = os.path.join(oldata, "loadtest")

node_0 = os.path.join(devnet, "0-Node")
node_1 = os.path.join(devnet, "1-Node")
node_2 = os.path.join(devnet, "2-Node")
node_3 = os.path.join(devnet, "3-Node")
node_4 = os.path.join(devnet, "4-Node")
node_5 = os.path.join(devnet, "5-Node")
node_6 = os.path.join(devnet, "6-Node")
node_7 = os.path.join(devnet, "7-Node")
node_8 = os.path.join(devnet, "8-Node")
node_9 = os.path.join(devnet, "9-Node")

fullnode_dev = os.path.join(devnet, "4-Node")
fullnode_prod = os.path.join(oldata, "fullnode")
