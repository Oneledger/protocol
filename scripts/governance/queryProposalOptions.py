from sdk import *


class bcolors:
    WARNING = '\033[93m'
    ENDC = '\033[0m'


if __name__ == "__main__":
    print bcolors.WARNING + "*** Proposal Options ***" + bcolors.ENDC
    print query_proposal_options()
