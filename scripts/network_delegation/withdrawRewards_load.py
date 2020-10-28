from sdk import *

cfg_dev = TestConfig(fullnode_dev, 1010000, 100)
cfg_prod = TestConfig(fullnode_prod, 1010000, 10000)

class WithdrawRewardsTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(WithdrawRewardsTxLoad, self).__init__(cfg, tid, "WithdrawRewardsTxLoad", free_thread=True)
        self.wait = True

    def setup(self, interval):
        super(WithdrawRewardsTxLoad, self).setup(interval)
        self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
        self.tx_deleg = NetWorkDelegate(self.test_account, "1000000", self.key_path)
        self.tx_deleg.send_network_Delegate(mode=TxCommit)
        self.tx_draw = WithdrawRewards(self.test_account, 1, self.key_path)

    def run_tx(self, i):
        if self.wait:
            self.log("waiting for rewards distribution...")
            self.tx_draw.waitfor_rewards("10")
            self.wait = False
        super(WithdrawRewardsTxLoad, self).run_tx(i)
        log = self.tx_draw.send("1", exit_on_err=True, mode=TxAsync)
        if len(log) > 0:
            self.log(log)
        if i % 10 == 0:
            self.wait = True

    def stop(self):
        super(WithdrawRewardsTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads):
        return [WithdrawRewardsTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads):
        return [WithdrawRewardsTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    