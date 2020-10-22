import os, errno
import os.path as path
import time, datetime
import threading

from rpc_call import oltest, loadtest
from common import *
from constant import *

class TestConfig:
    def __init__(self, node_root, init_fund, numof_txs, interval=INTERVAL_NORMAL):
        """
        Args:
            node_root (str): e.g., /opt/data/devnet/0-Node or /opt/data/fullnode
            init_fund (int): total amount of OLT needed for this thread's tests
            numof_txs (int): total number of txs this thread will send to node
            interval  (int): interval in milliseconds between 2 Txs.
        """
        self.node_root = node_root
        self.init_fund = init_fund
        self.numof_txs = numof_txs
        self.interval = interval
        self.test_root = loadtest

class TxLoad(threading.Thread):
    def __init__(self, cfg, tid, name):
        super(TxLoad, self).__init__()
        self.cfg = cfg
        self.tid = tid
        self.name = name
        self.stop_event = threading.Event()
        self.test_path = path.join(cfg.test_root, name)
        self.key_path = path.join(cfg.node_root, "keystore")
        self.log_file = path.join(self.test_path, "{}_thread_{}.log".format(self.name, self.tid))

    def setup(self):
        if not path.exists(self.cfg.test_root):
            os.mkdir(self.cfg.test_root)
        if not path.exists(self.test_path):
            os.mkdir(self.test_path)
        self.flog = open(self.log_file, "a+")
        self.log("{}_thread_{} setting up...".format(self.name, self.tid))
        self.node_account = nodeAccount(self.cfg.node_root)

    def run(self):
        self.log("{}_thread_{} started".format(self.name, self.tid))
        for i in range(self.cfg.numof_txs):
            self.run_tx(i + 1)
            time.sleep(self.cfg.interval / 1000)
            if self.stop_event.is_set():
                break
        self.log("{}_thread_{} stopped".format(self.name, self.tid))
        self.flog.close()

    def run_tx(self, i):
        self.log("{}_thread_{} sending {}th transactions".format(self.name, self.tid, i))

    def stop(self):
        self.log("{}_thread_{} stopping...".format(self.name, self.tid))
        if not self.stop_event.is_set():
            self.stop_event.set()

    def log(self, msg, stdout=True):
        now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        msg = now + "---" + msg
        self.flog.write(msg + "\n")
        self.flog.flush()
        if stdout:
            print msg
