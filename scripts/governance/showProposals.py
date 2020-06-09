from sdk import *

if __name__ == "__main__":
    print "#### ACTIVE PROPOSALS: ####"
    activeList, funds, votes = query_proposals("active")
    print activeList

    print "#### PASSED PROPOSALS: ####"
    passedList, funds, votes = query_proposals("passed")
    print passedList

    print "#### FAILED PROPOSALS: ####"
    failedList, funds, votes = query_proposals("failed")
    print failedList
