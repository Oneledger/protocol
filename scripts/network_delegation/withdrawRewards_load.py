from sdk import *

cfg_dev = TestConfig(fullnode_dev, 1010000, 100)
cfg_prod = TestConfig(fullnode_prod, 1010000, 10000)

class WithdrawRewardsTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(WithdrawRewardsTxLoad, self).__init__(cfg, tid, "WithdrawRewardsTxLoad", free_thread=True)
        self.balance = 0

    def setup(self, interval):
        super(WithdrawRewardsTxLoad, self).setup(interval)
        self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
        self.tx_deleg = NetWorkDelegate(self.test_account, '1000000' + '0' * 18, self.key_path)
        self.tx_deleg.send_network_Delegate(mode=TxCommit)
        self.tx_draw = WithdrawRewards(self.test_account, 1, self.key_path)

    def run_tx(self, i):
        if self.balance == 0:
            self.log("waiting for rewards distribution...")
            self.balance = self.tx_draw.waitfor_rewards('10' + '0' * 18, "balance")
            self.log("rewards distributed: {} OLT".format(self.balance))
        super(WithdrawRewardsTxLoad, self).run_tx(i)
        log = self.tx_draw.send(exit_on_err=False, mode=TxAsync)
        self.balance -= 1
        if len(log) > 0:
            self.log(log)

    def stop(self):
        super(WithdrawRewardsTxLoad, self).stop()

    @classmethod
    def dev(cls, numof_threads):
        return [WithdrawRewardsTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]

    @classmethod
    def prod(cls, numof_threads):
        return [WithdrawRewardsTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
    