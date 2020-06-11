import sys
import time
from sdk import *

addr_list = addresses()

_pid = 21050
_proposer_0 = addr_list[0]
_proposer_1 = addr_list[1]
_proposer_2 = addr_list[2]

_initial_funding = (int("2") * 10 ** 9)
_big_funding = (int("8") * 10 ** 9)
_total_funding = _initial_funding + _big_funding
_funding_goal_general = (int("10") * 10 ** 9)

def gen_prop(proposer, prop_type):
    global _pid
    prop = Proposal(str(_pid), prop_type, "proposal for fund", "proposal headline", proposer, _initial_funding)
    _pid += 1
    return prop

def create_some_proposals():
    result_by_id = {}
    result_by_proposer = {_proposer_0: [], _proposer_1: [], _proposer_2: []}
    result_by_type = {ProposalTypeGeneral:[], ProposalTypeCodeChange: [], ProposalTypeConfigUpdate: []}

    # create a proposal that failed VOTING
    prop_0 = gen_prop(_proposer_0, "general")
    prop_0.send_create()
    time.sleep(1)
    fund_proposal(prop_0.pid, _big_funding, addr_list[1])
    vote_proposal(prop_0.pid, "NO", url_0, addr_list[0])
    vote_proposal(prop_0.pid, "NO", url_1, addr_list[1])
    result_by_id[prop_0.pid] = (ProposalTypeGeneral, ProposalOutcomeInsufficientVotes, ProposalStatusCompleted, _total_funding)
    result_by_proposer[_proposer_0].append(prop_0.pid)
    result_by_type[ProposalTypeGeneral].append(prop_0.pid)

    # create a proposal that passed VOTING
    prop_1 = gen_prop(_proposer_1, "codechange")
    prop_1.send_create()
    time.sleep(1)
    fund_proposal(prop_1.pid, _big_funding, addr_list[0])
    vote_proposal(prop_1.pid, "YES", url_0, addr_list[0])
    vote_proposal(prop_1.pid, "YES", url_1, addr_list[1])
    vote_proposal(prop_1.pid, "YES", url_2, addr_list[2])
    result_by_id[prop_1.pid] = (ProposalTypeCodeChange, ProposalOutcomeCompleted, ProposalStatusCompleted, _total_funding)
    result_by_proposer[_proposer_1].append(prop_1.pid)
    result_by_type[ProposalTypeCodeChange].append(prop_1.pid)

    # create a proposal that canceled
    prop_2 = gen_prop(_proposer_1, "codechange")
    prop_2.send_create()
    time.sleep(1)
    cancel_proposal(prop_2.pid, _proposer_1, "changed mind")
    result_by_id[prop_2.pid] = (ProposalTypeCodeChange, ProposalOutcomeCancelled, ProposalStatusCompleted, _initial_funding)
    result_by_proposer[_proposer_1].append(prop_2.pid)
    result_by_type[ProposalTypeCodeChange].append(prop_2.pid)

    # create a proposal that passed FUNDING
    prop_3 = gen_prop(_proposer_2, "general")
    prop_3.send_create()
    time.sleep(1)
    fund_proposal(prop_3.pid, _big_funding, addr_list[0])
    result_by_id[prop_3.pid] = (ProposalTypeGeneral, ProposalOutcomeInProgress, ProposalStatusVoting, _total_funding)
    result_by_proposer[_proposer_2].append(prop_3.pid)
    result_by_type[ProposalTypeGeneral].append(prop_3.pid)

    # create a proposal that in FUNDING
    prop_4 = gen_prop(_proposer_2, "configupdate")
    prop_4.send_create()
    time.sleep(1)
    result_by_id[prop_4.pid] = (ProposalTypeConfigUpdate, ProposalOutcomeInProgress, ProposalStatusFunding, _initial_funding)
    result_by_proposer[_proposer_2].append(prop_4.pid)
    result_by_type[ProposalTypeConfigUpdate].append(prop_4.pid)

    return result_by_id, result_by_proposer, result_by_type

def check_proposals(props, expected_pids):
    pids = []
    for prop in props:
        pids.append(prop["proposal"]["proposal_id"])
    expected_pids.sort()
    pids.sort()
    if pids != expected_pids:
        sys.exit(-1)

if __name__ == "__main__":
    # create some proposals
    by_id, by_proposer, by_type = create_some_proposals()

    # query each proposal by proposal ID and check
    for pid, expected in by_id.items():
        prop_type, outcome, status, funds = expected
        check_proposal_state(pid, outcome, status, prop_type, funds)
    
    print bcolors.OKGREEN + "#### Test query proposals by ID succeed" + bcolors.ENDC

    # query proposals of proposer_0
    props = query_proposals("", _proposer_0, "")
    check_proposals(props, by_proposer[_proposer_0])

    # query proposals of proposer_1
    props = query_proposals("", _proposer_1, "")
    check_proposals(props, by_proposer[_proposer_1])

    # query proposals of proposer_2
    props = query_proposals("", _proposer_2, "")
    check_proposals(props, by_proposer[_proposer_2])

    print bcolors.OKGREEN + "#### Test query proposals by proposer succeed" + bcolors.ENDC

    # query proposals of type "general"
    props = query_proposals("", "", "general")
    check_proposals(props, by_type[ProposalTypeGeneral])

    # query proposals of type "codechange"
    props = query_proposals("", "", "codechange")
    check_proposals(props, by_type[ProposalTypeCodeChange])

    # query proposals of type "configupdate"
    props = query_proposals("", "", "configupdate")
    check_proposals(props, by_type[ProposalTypeConfigUpdate])

    print bcolors.OKGREEN + "#### Test query proposals by type succeed" + bcolors.ENDC
