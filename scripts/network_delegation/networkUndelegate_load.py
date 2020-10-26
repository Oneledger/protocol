from sdk import *

cfg_dev = TestConfig(fullnode_dev, 110000, 100, INTERVAL_NORMAL)
cfg_prod = TestConfig(fullnode_prod, 110000, 10000, INTERVAL_NORMAL)

class UnDelegateTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(UnDelegateTxLoad, self).__init__(cfg, tid, "UnDelegateTxLoad")

    def setup(self):
        super(UnDelegateTxLoad, self).setup()
        self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
        self.tx = NetWorkDelegate(self.test_account, "100000", self.key_path)
        self.tx.send_network_Delegate(mode=TxCommit)

    def run_tx(self, i):
        super(UnDelegateTxLoad, self).run_tx(i)
        log = self.tx.send_network_undelegate("1", exit_on_err=False, mode=TxAsync)
        if len(log) > 0:
            self.log(log)

    def stop(self):
        super(UnDelegateTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads, interval):
        cfg_dev.interval = interval
        return [UnDelegateTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads, interval):
        cfg_prod.interval = interval
        return [UnDelegateTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    