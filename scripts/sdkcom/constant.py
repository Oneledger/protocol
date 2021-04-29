
class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

# default load test intervals, in milliseconds
INTERVAL_DEFAULT = 200

# load test speed, txs per second
TXS_PER_SEC_NORMAL = 100
TXS_PER_SEC_SLOW = 10

# broadcast mode
TxCommit = 1
TxSync = 2
TxAsync = 3
