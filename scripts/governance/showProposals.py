from sdk import *

if __name__ == "__main__":
    print "#### ACTIVE PROPOSALS: ####"
    activeList = query_proposals("active")
    print activeList

    print "#### PASSED PROPOSALS: ####"
    passedList = query_proposals("passed")
    print passedList

    print "#### FAILED PROPOSALS: ####"
    failedList = query_proposals("failed")
    print failedList
