import os
import time
import subprocess
from actions import *

def fund_proposal(pid, amount, funder, secs=1):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(secs)

def cancel_proposal(pid, proposer, reason, secs=1):
    # fund the proposal
    prop_cancel = ProposalCancel(pid, proposer, reason)
    res = prop_cancel.send_cancel()
    time.sleep(secs)
    return res

def vote_proposal(pid, opinion, url, address, secs=1):
    # vote the proposal
    prop_vote = ProposalVote(pid, opinion, url, address)
    prop_vote.send_vote()
    time.sleep(secs)

def vote_proposal_cli(pid, opinion, node, address, secs=1):
    # vote the proposal through CLI
    args = ['olclient', 'gov', 'vote', '--root', node, '--id', pid, '--address', address[3:], '--opinion', opinion, '--password', 'pass', '--gasprice', '0.00001', '--gas', '40000']

    # set cwd for the purpose of wallet path
    process = subprocess.Popen(args, cwd=os.getcwd())
    process.wait()
    time.sleep(secs)

    # check return code
    if process.returncode != 0:
        print "olclient vote failed"
        sys.exit(-1)
    print "################### proposal voted:" + pid + "opinion: " + opinion

def list_proposal_cli(pid, node):
    # vote the proposal through CLI
    args = ['olclient', 'gov', 'list', '--root', node, '--id', pid]
    process = subprocess.Popen(args)
    process.wait()

    # check return code
    if process.returncode != 0:
        print "olclient list proposal failed"
        sys.exit(-1)

def check_proposal_state(pid, outcome_expected, status_expected, type_expected=ProposalTypeGeneral, funds=-1):
    # check proposal status, outcome, status, fund
    prop, cur_fund = query_proposal(pid)
    if prop['status'] != status_expected:
        sys.exit(-1)
    if prop['outcome'] != outcome_expected:
        sys.exit(-1)
    if prop['proposal_type'] != type_expected:
        sys.exit(-1)
    cur_fund = int(cur_fund)
    if funds != -1 and funds != cur_fund:
        sys.exit(-1)
    