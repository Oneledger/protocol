import subprocess
import time

from actions import *

success = "Returned Successfully"

option_types = ["proposal", "staking", "fee", "evidence", "currency", "eth", "btc", "ons", "rewards", "currency"]


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


def finalize_proposal(pid, address, secs=1):
    # fund the proposal
    prop_finalize = ProposalFinalize(pid, address)
    prop_finalize.send_finalize()
    time.sleep(secs)


def vote_proposal_cli(pid, opinion, node, address, secs=1):
    # vote the proposal through CLI
    args = ['olclient', 'gov', 'vote', '--root', node, '--id', pid, '--address', address[3:], '--opinion', opinion,
            '--password', 'pass', '--gasprice', '0.00001', '--gas', '40000']

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
    if prop['proposalType'] != type_expected:
        sys.exit(-1)
    cur_fund = int(cur_fund)
    if funds != -1 and funds != cur_fund:
        sys.exit(-1)


def getActiveValidators():
    args = ['olclient', 'validatorset']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    active_count = 0
    for out in output:
        if "Active true" in out:
            active_count = active_count + 1
    return active_count


def getAllValidators():
    args = ['olclient', 'validatorset']
    process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    active_count = 0
    for out in output:
        if "Active" in out:
            active_count = active_count + 1
    return active_count


def addValidatorAccounts(node):
    args = ['olclient', 'show_node_id']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    time.sleep(1)
    pubKey = output[0].split(",")[0].split(":")[1].strip()
    f = open(os.path.join(node, "consensus", "config", "node_key.json"), "r")
    contents = json.loads(f.read())
    privKey = contents['priv_key']['value']
    args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", '1234']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.readlines()
    time.sleep(1)
    return output[1].split(":")[1].strip()[3:]


def stake(node, amount='3000000'):
    validatorAccount = addValidatorAccounts(node)
    # trasfer funds from node 0 to staking validator
    # args = ['olclient', 'send', '--party', parentnodeaddre, "--counterparty", validatorAccount, '--amount', '100',
    #         '--password',
    #         '1234', '--fee', '0.001']
    # process = subprocess.Popen(args, cwd=node_0, stdout=subprocess.PIPE)
    # process.wait()
    # output = process.stdout.read()
    # print output
    args = ['olclient', 'delegation', 'stake', '--address', validatorAccount, '--amount', amount, '--password',
            '1234']
    process = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process.wait()
    output = process.stdout.read()
    if not success in output:
        print "Stake was not successful"
    print bcolors.OKBLUE + "#### Stake Successfull for :" + node + bcolors.ENDC


def validateCatchup(node, node2):
    args = ['olclient', 'validatorset', 'status']
    # set protocol root path as current path
    process1 = subprocess.Popen(args, cwd=node, stdout=subprocess.PIPE)
    process2 = subprocess.Popen(args, cwd=node2, stdout=subprocess.PIPE)
    process1.wait()
    process2.wait()
    if process1.returncode != 0:
        print "olclient check validatorset failed returncode"
        sys.exit(-1)
    if process2.returncode != 0:
        print "olclient check validatorset failed returncode"
        sys.exit(-1)
    output1 = process1.stdout.readlines()
    height1 = output1[len(output1) - 1].split(" ")[1]
    output2 = process2.stdout.readlines()
    height2 = output2[len(output2) - 1].split(" ")[1]
    if height2 != height1:
        print "Node Failed to catchup"
        sys.exit(-1)
    return height1
    # check return code

    # check if on list


def clean_and_catchup():
    opt = query_governanceState()

    old_heights = [None] * len(option_types)
    for idx, val in enumerate(option_types):
        old_heights[idx] = opt["lastUpdateHeight"][val]

    call_with_args = "./scripts/governance/clean '%s'" % (str(0))
    os.system(call_with_args)
    # Wait for it to catchup
    time.sleep(30)
    # Validate Catchup
    height = validateCatchup(node_0, node_1)
    opt = query_governanceState()
    for idx, val in enumerate(option_types):
        if opt["lastUpdateHeight"][val] != old_heights[idx]:
            print "Last update height for " + val + "is :" + opt["lastUpdateHeight"][val] + " , it was : " + \
                  old_heights[idx] + " Before Catchup"
            sys.exit(-1)
    print bcolors.OKBLUE + "Last Update Height After Catchup : " + str(opt["lastUpdateHeight"]) + bcolors.ENDC
    print bcolors.OKBLUE + "Block Height After Catchup : " + str(height) + bcolors.ENDC
