from sdk import *

cfg_dev = TestConfig(fullnode_dev, 110000, 100)
cfg_prod = TestConfig(fullnode_prod, 110000, 1000)

class UnDelegateTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(UnDelegateTxLoad, self).__init__(cfg, tid, "UnDelegateTxLoad")

    def setup(self, interval):
        super(UnDelegateTxLoad, self).setup(interval)
        self.tx = NetWorkDelegate(self.test_account, '100000' + '0' * 18, self.key_path)
        self.tx.send_network_Delegate(mode=TxCommit)

    def run_tx(self, i):
        super(UnDelegateTxLoad, self).run_tx(i)
        log = self.tx.send_network_undelegate('1' + '0' * 18, exit_on_err=False, mode=TxAsync)
        if len(log) > 0:
            self.log(log)

    def stop(self):
        super(UnDelegateTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads):
        return [UnDelegateTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads):
        return [UnDelegateTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    