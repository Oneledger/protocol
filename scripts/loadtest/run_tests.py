import sys, getopt, signal
from test_threads import *
from sdkcom import *
from network_delegation import *

# for local
def add_threads_dev(threads):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.dev(1))

    # UnDelegate, 2 threads
    threads.add_threads(UnDelegateTxLoad.dev(1))

    # WithdrawRewards, 2 threads
    threads.add_threads(WithdrawRewardsTxLoad.dev(1))

    # ReinvestRewards, 2 threads
    threads.add_threads(ReinvestRewardsTxLoad.dev(1))

# for devnet
def add_threads_prod(threads):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.prod(1))

    # UnDelegate, 2 threads
    threads.add_threads(UnDelegateTxLoad.prod(1))

    # WithdrawRewards, 2 threads
    threads.add_threads(WithdrawRewardsTxLoad.prod(1))

    # ReinvestRewards, 2 threads
    threads.add_threads(ReinvestRewardsTxLoad.prod(1))

def abort_loadtest(signal, frame):
    threads.stop_threads()
    sys.exit(0)

def parse_params(argv):
    clean_run = False
    txs_persec = TXS_PER_SEC_NORMAL
    try:
      opts, args = getopt.getopt(argv,"cs:",["speed="])
    except getopt.GetoptError:
      print 'run_tests.py -s <speed>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-c", "--clean"):
         clean_run = True
      if opt in ("-s", "--speed"):
         txs_persec = arg
    return clean_run, 1000 / int(txs_persec)

if __name__ == "__main__":
    # parse options
    clean_run, interval = parse_params(sys.argv[1:])

    # configuration based on environment
    if oltest == "1":
        add_threads_dev(threads)
    else:
        add_threads_prod(threads)

    # clean up test folder
    if clean_run:
        threads.clean()

    # setup threads before run
    threads.setup_threads(interval)

    # run threads
    signal.signal(signal.SIGINT, abort_loadtest)
    threads.run_threads()

    # join threads
    threads.join_threads()
