from sdk import *

if __name__ == "__main__":
    print "#### ACTIVE PROPOSALS: ####"
    activeList = query_proposals(ProposalStateActive)
    print activeList

    print "#### PASSED PROPOSALS: ####"
    passedList = query_proposals(ProposalStatePassed)
    print passedList

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals(ProposalStateFailed)
    print failedList
