import sys, getopt, signal
from test_threads import *
from sdkcom import *
from network_delegation import *

# for local
def add_threads_dev(threads):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.dev(2))

    # UnDelegate, 2 threads
    threads.add_threads(UnDelegateTxLoad.dev(2))

    # WithdrawRewards, 2 threads
    threads.add_threads(WithdrawRewardsTxLoad.dev(2))

    # ReinvestRewards, 2 threads
    threads.add_threads(ReinvestRewardsTxLoad.dev(2))

# for devnet
def add_threads_prod(threads, interval):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.prod(2))

    # UnDelegate, 2 threads
    threads.add_threads(UnDelegateTxLoad.prod(2))

    # UnDelegate, 2 threads
    threads.add_threads(WithdrawRewardsTxLoad.prod(2))

def abort_loadtest(signal, frame):
    threads.stop_threads()
    sys.exit(0)

def parse_params(argv):
    txs_persec = TXS_PER_SEC_NORMAL
    try:
      opts, args = getopt.getopt(argv,"s:",["speed="])
    except getopt.GetoptError:
      print 'run_tests.py -s <speed>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-s", "--speed"):
         txs_persec = arg
    return 1000 / int(txs_persec)

if __name__ == "__main__":
    # parse options
    interval = parse_params(sys.argv[1:])

    # configuration based on environment
    if oltest == "1":
        add_threads_dev(threads)
    else:
        add_threads_prod(threads)

    # clean up test folder
    threads.clean()

    # setup threads before run
    threads.setup_threads(interval)

    # run threads
    signal.signal(signal.SIGINT, abort_loadtest)
    threads.run_threads()

    # join threads
    threads.join_threads()
