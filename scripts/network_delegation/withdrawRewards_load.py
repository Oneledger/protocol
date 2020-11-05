from sdk import *

cfg_dev = TestConfig(fullnode_dev, 1010000, 100)
cfg_prod = TestConfig(fullnode_prod, 1010000, 1000)

class WithdrawRewardsTxLoad(TxLoad):
    def __init__(self, cfg, tid):
        super(WithdrawRewardsTxLoad, self).__init__(cfg, tid, "WithdrawRewardsTxLoad", free_thread=True)
        self.balance = 0 # in 0.001 OLT

    def setup(self, interval):
        new_run = super(WithdrawRewardsTxLoad, self).setup(interval, False)
        if new_run:
            self.tx_deleg = NetWorkDelegate(self.test_account, '1000000' + '0' * 18, self.key_path)
            self.tx_deleg.send_network_Delegate(mode=TxAsync)
        self.tx_draw = WithdrawRewards(self.test_account, 10**16, self.key_path)

    def run_tx(self, i):
        if self.balance == 0:
            wait_for(1) # wait 1 block to refresh reduced rewards balance
            self.log("waiting for rewards distribution...")
            self.balance = waitfor_rewards(self.tx_draw.delegator, "1", "balance") * 100
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
    