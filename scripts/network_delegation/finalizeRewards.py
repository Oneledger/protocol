# below is removed since finalize withdraw rewards logic is moved to block beginner, OLP-1266
# from sdk import *
#
#
# def delegate(node, account, amount):
#     newDelegation = NetWorkDelegate(account, amount, node + "/keystore/")
#     newDelegation.send_network_Delegate()
#
#
# def check_rewards(result, balance, matured, pending):
#     if balance != '':
#         balance = int(balance) * pow(10, 18)
#         if int(result['balance']) < balance:
#             sys.exit(-1)
#     if matured != '' and result['matured'] != matured:
#         sys.exit(-1)
#     if pending != None:
#         if len(result['pending']) != len(pending):
#             sys.exit(-1)
#         for i, amt in enumerate(pending):
#             if amt != result['pending'][i]['amount']:
#                 sys.exit(-1)
#
#
# def check_total_rewards(result, expected_exclude_withdrawn):
#     if result < expected_exclude_withdrawn:
#         sys.exit(-1)
#
#
# def check_balance(before, after, expected_diff):
#     diff = after - before
#     if not (diff == expected_diff or diff == expected_diff + 1 or diff == expected_diff - 1):
#         print diff
#         print expected_diff
#         sys.exit(-1)
#
#
# if __name__ == "__main__":
#     # create validator account
#     funder = addValidatorWalletAccounts(node_0)
#
#     # create delegator account
#     delegator = createAccount(node_0, 2500000, funder)
#
#     # delegates some OLT and wait for rewards distribution
#     delegate(node_0, delegator, '2000000' + '0' * 18)
#     wait_for(4)
#
#     # query and check balance
#     res = query_rewards(delegator)
#     check_rewards(res, '6', '0', [])
#
#     # initiate 2 withdrawals
#     pending = []
#     total = 0
#     for i in range(2):
#         amt = i + 2
#         amt_long = str(amt) + '0' * 18
#         withdraw = WithdrawRewards(delegator, amt_long, node_0 + "/keystore/")
#         withdraw.send(True)
#         pending.append(str(amt) + '0' * 18)
#         total += amt
#         wait_for(1)
#     total_str = str(total) + '0' * 18
#
#     # query and check pending withdrawal
#     res = query_rewards(delegator)
#     check_rewards(res, '0', '0', pending)
#     print "#### Successfully withdrawn delegator rewards"
#
#     # query and check again after maturity
#     wait_for(4)
#     res = query_rewards(delegator)
#     check_rewards(res, '', total_str, [])
#     print "#### Successfully matured delegator rewards"
#
#     # finalize more than withdrawn
#     finalize = FinalizeRewards(delegator, node_0 + "/keystore/")
#     finalize.send_finalize(int(total_str) * 2, False)
#     print bcolors.OKGREEN + "#### finalize rewards more than withdrawn failed as expected" + bcolors.ENDC
#
#     # query balance
#     balance_before = query_balance(delegator)
#     # finalize withdrawn rewards
#     finalize.send_finalize(total_str, True)
#     # query and check
#     wait_for(2)
#     res = query_rewards(delegator)
#     check_rewards(res, '', '0', [])
#     balance_after = query_balance(delegator)
#     check_balance(balance_before, balance_after, total)
#     print bcolors.OKGREEN + "#### Successfully finalized delegator rewards" + bcolors.ENDC
#
#     # withdraw all balance
#     withdraw_all_balance = int(res['balance'])
#     withdraw = WithdrawRewards(delegator, withdraw_all_balance, node_0 + "/keystore/")
#     withdraw.send(True)
#     print "#### Successfully withdrawn all rewards"
#
#     # finalize withdraw rewards again
#     wait_for(6)
#     finalize.send_finalize(str(withdraw_all_balance), True)
#     # query and check
#     wait_for(3)
#     res = query_rewards(delegator)
#     check_rewards(res, '', '0', [])
#     balance_final = query_balance(delegator)
#     check_balance(balance_after, balance_final, withdraw_all_balance/1000000000000000000)
#     print bcolors.OKGREEN + "#### Successfully finalized all rewards" + bcolors.ENDC
#
#     # below is to test total rewards query
#     # create another delegator account
#     funder1 = addValidatorWalletAccounts(node_1)
#     delegator1 = createAccount(node_1, 8000000, funder1)
#
#     # delegates some OLT and wait for rewards distribution
#     delegate(node_1, delegator1, '5000000' + '0' * 18)
#     wait_for(4)
#
#     # initiate 1 withdrawal
#     amt1 = 3
#     amt1_long = str(amt1) + '0' * 18
#     withdraw1 = WithdrawRewards(delegator1, amt1_long, node_1 + "/keystore/")
#     withdraw1.send(True)
#     wait_for(1)
#
#     # finalize withdraw rewards
#     finalize1 = FinalizeRewards(delegator1, node_1 + "/keystore/")
#     wait_for(6)
#     finalize1.send_finalize(str(amt1) + '0' * 18, True)
#
#     # query total rewards and check
#     res = query_rewards(delegator)
#     res1 = query_rewards(delegator1)
#     total = query_total_rewards()
#     check_total_rewards(total['totalRewards'], int(res['balance']) * pow(10, 18) + int(res1['balance']) * pow(10, 18))
#     print bcolors.OKGREEN + "#### Successfully tested query total rewards" + bcolors.ENDC
