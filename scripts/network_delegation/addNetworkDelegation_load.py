from sdk import *

cfg_dev = TestConfig(fullnode_dev, 110000, 100)
cfg_prod = TestConfig(fullnode_prod, 110000, 10000)

class DelegateTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(DelegateTxLoad, self).__init__(cfg, tid, "DelegateTxLoad")

    def setup(self):
        super(DelegateTxLoad, self).setup()
        self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
        self.tx = NetWorkDelegate(self.test_account, "1", self.key_path)

    def run_tx(self, i):
        super(DelegateTxLoad, self).run_tx(i)
        log = self.tx.send_network_Delegate(exit_on_err=False, mode=TxAsync)
        if len(log) > 0:
            self.log(log)

    def stop(self):
        super(DelegateTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads, interval):
        cfg_dev.interval = interval
        return [DelegateTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads, interval):
        cfg_prod.interval = interval
        return [DelegateTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    