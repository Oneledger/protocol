import sys, getopt, signal
from test_threads import *
from sdkcom import *
from network_delegation import *

# for local
def add_threads_dev(threads, interval):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.dev(3, interval))

# for devnet
def add_threads_prod(threads, interval):
    # Delegate, 2 threads
    threads.add_threads(DelegateTxLoad.prod(3, interval))

def abort_loadtest(signal, frame):
    threads.stop_threads()
    sys.exit(0)

def parse_params(argv):
    interval = INTERVAL_NORMAL
    try:
      opts, args = getopt.getopt(argv,"i:",["interval="])
    except getopt.GetoptError:
      print 'run_tests.py -i <interval>'
      sys.exit(-1)
    for opt, arg in opts:
      if opt in ("-i", "--interval"):
         interval = arg
    return int(interval)

if __name__ == "__main__":
    # parse options
    interval = parse_params(sys.argv[1:])

    # configuration based on environment
    if oltest == "1":
        add_threads_dev(threads, interval)
    else:
        add_threads_prod(threads, interval)

    # clean up test folder
    threads.clean()

    # setup threads before run
    threads.setup_threads()

    # run threads
    signal.signal(signal.SIGINT, abort_loadtest)
    threads.run_threads()
