import time
from actions import *

def fund_proposal(pid, amount, funder, secs=1):
    # fund the proposal
    prop_fund = ProposalFund(pid, amount, funder)
    prop_fund.send_fund()
    time.sleep(secs)

def vote_proposal(pid, opinion, url, address, secs=1):
    # fund the proposal
    prop_vote = ProposalVote(pid, opinion, url, address)
    prop_vote.send_vote()
    time.sleep(secs)

def check_proposal_state(pid, state_expected, status_expected):
    # check proposal state
    prop = query_proposal(pid)
    if prop['state'] != state_expected:
        sys.exit(-1)
    if prop['proposal']['Status'] != status_expected:
        sys.exit(-1)
