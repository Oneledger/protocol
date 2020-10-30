from sdk import *
# below is removed since withdraw logic is moved to block beginner, OLP-1267

# cfg_dev = TestConfig(fullnode_dev, 110000, 100)
# cfg_prod = TestConfig(fullnode_prod, 110000, 10000)
#
# class WithdrawDelegationTxLoad(TxLoad):
#     def __init__(self, cfg, tid):
#         super(WithdrawDelegationTxLoad, self).__init__(cfg, tid, "WithdrawDelegationTxLoad")
#
#     def setup(self):
#         super(WithdrawDelegationTxLoad, self).setup()
#         self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
#         self.tx = NetWorkDelegate(self.test_account, "1", self.key_path)
#
#     def run_tx(self, i):
#         super(WithdrawDelegationTxLoad, self).run_tx(i)
#         log = self.tx.send_network_Delegate(exit_on_err=False, mode=TxAsync)
#         if len(log) > 0:
#             self.log(log)
#
#     def stop(self):
#         super(WithdrawDelegationTxLoad, self).stop()
#
#     @classmethod
#     def dev(cls, numof_threads, interval):
#         cfg_dev.interval = interval
#         return [WithdrawDelegationTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]
#
#     @classmethod
#     def prod(cls, numof_threads, interval):
#         cfg_prod.interval = interval
#         return [WithdrawDelegationTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]