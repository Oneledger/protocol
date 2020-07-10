from sdk import *
addr_list = addresses()
_pid = 22003
_proposer_0 = addr_list[0]
_initial_funding = (int("2") * 10 ** 9)
_funding = (int("1") * 10 ** 9)

def gen_prop(proposer, prop_type):
    global _pid
    prop = Proposal(str(_pid), prop_type, "proposal for fund", "proposal headline", proposer, _initial_funding)
    _pid += 1
    return prop


if __name__ == "__main__":
    prop_0 = gen_prop(_proposer_0, "general")
    prop_0.send_create()
    time.sleep(1)
    print addr_list[0]
    fund_proposal(prop_0.pid, _funding, addr_list[0])
    fund_proposal(prop_0.pid, _funding, addr_list[0])

    result = get_funds_for_proposal_by_funder(prop_0.pid, _proposer_0)
    if "amount" not in result:
        sys.exit(-1)
    amount = result["amount"]
    if amount != "4000000000":
        sys.exit(-1)

    print bcolors.OKGREEN + "#### Test get funds for proposal by funder succeed" + bcolors.ENDC

