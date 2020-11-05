from sdk import *

cfg_dev = TestConfig(fullnode_dev, 1010000, 100)
cfg_prod = TestConfig(fullnode_prod, 1010000, 10000)

class ReinvestRewardsTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(ReinvestRewardsTxLoad, self).__init__(cfg, tid, "ReinvestRewardsTxLoad", free_thread=True)
        self.balance = 0 # in 0.001 OLT

    def setup(self, interval):
        new_run = super(ReinvestRewardsTxLoad, self).setup(interval, False)
        if new_run:
            self.tx_deleg = NetWorkDelegate(self.test_account, '1000000' + '0' * 18, self.key_path)
            self.tx_deleg.send_network_Delegate(mode=TxAsync)
        self.tx_invest = ReinvestRewards(self.test_account, self.key_path)

    def run_tx(self, i):
        if self.balance == 0:
            self.log("waiting for rewards distribution...")
            self.balance = waitfor_rewards(self.test_account, '1', "balance") * 100
            self.log("rewards distributed: {} OLT".format(self.balance))
        super(ReinvestRewardsTxLoad, self).run_tx(i)
        log = self.tx_invest.send(10**16, exit_on_err=False, mode=TxAsync)
        self.balance -= 1
        if len(log) > 0:
            self.log(log)

    def stop(self):
        super(ReinvestRewardsTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads):
        return [ReinvestRewardsTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads):
        return [ReinvestRewardsTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    