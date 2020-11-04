from sdk import *
# below is removed since withdraw logic is moved to block beginner, OLP-1267

# cfg_dev = TestConfig(fullnode_dev, 110000, 100)
# cfg_prod = TestConfig(fullnode_prod, 110000, 10000)
#
# class WithdrawDelegationTxLoad(TxLoad):
#     def __init__(self, cfg, tid):
#         super(WithdrawDelegationTxLoad, self).__init__(cfg, tid, "WithdrawDelegationTxLoad")
#         self.wait = True
#
#     def setup(self, interval):
#         super(WithdrawDelegationTxLoad, self).setup(interval)
#         self.test_account = createAccount(node=self.cfg.node_root, funds=self.cfg.init_fund, funder=self.node_account)
#         self.tx = NetWorkDelegate(self.test_account, "100000", self.key_path)
#         self.tx.send_network_Delegate(mode=TxCommit)
#         self.tx.send_network_undelegate("100000", True, mode=TxCommit)
#
#     def run_tx(self, i):
#         if self.wait:
#             self.log("waiting for undelegated amount to mature...")
#             self.tx.waitfor_matured("100000")
#             self.wait = False
#         super(WithdrawDelegationTxLoad, self).run_tx(i)
#         log = self.tx.send_network_withdraw("1", exit_on_err=False, mode=TxAsync)
#         if len(log) > 0:
#             self.log(log)
#
#     def stop(self):
#         super(WithdrawDelegationTxLoad, self).stop()
#
#     @classmethod
#     def dev(cls, numof_threads):
#         return [WithdrawDelegationTxLoad(cfg_dev, tid+1) for tid in range(numof_threads)]
#
#     @classmethod
#     def prod(cls, numof_threads):
#         return [WithdrawDelegationTxLoad(cfg_prod, tid+1) for tid in range(numof_threads)]
