import os, shutil
import os.path as path
from sdkcom import oltest, loadtest, fullnode_dev, fullnode_prod, addValidatorWalletAccounts

class TestThreads:
    def __init__(self):
        self.threads = []

    def add_threads(self, threads):
        self.threads.extend(threads)

    def setup_threads(self, interval):
        # setup node account
        # automatically create node account for convenience
        if oltest == "1":
            addValidatorWalletAccounts(fullnode_dev)
        else:
            addValidatorWalletAccounts(fullnode_prod)

        # setup each thread
        thread_interval = interval * len(self.threads)
        for i, t in enumerate(self.threads):
            t.setup(thread_interval)

    def run_threads(self):
        # start all threads
        for i, t in enumerate(self.threads):
            t.start()

    def stop_threads(self):
        # stop all threads
        for i, t in enumerate(self.threads):
            t.stop()
        self.join_threads()

    def join_threads(self):
        for t in self.threads:
            t.join()

    def clean(self):
        # delete dirs and files
        print("cleaning loadtest...")
        shutil.rmtree(loadtest, ignore_errors=True)

threads = TestThreads()
