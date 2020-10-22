
class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'
    
# load test intervals, in milliseconds
INTERVAL_NORMAL = 20
INTERVAL_MAXSPEED = 10

# broadcast mode
TxCommit = 1
TxSync = 2
TxAsync = 3
