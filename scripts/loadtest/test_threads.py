import os, shutil
import os.path as path
from sdkcom import oltest, loadtest, fullnode_dev, addValidatorWalletAccounts

class TestThreads:
    def __init__(self):
        self.threads = []

    def add_threads(self, threads):
        self.threads.extend(threads)

    def setup_threads(self):
        # setup node account
        if oltest == "1":
            # in local,  automatically create node account for convenience
            addValidatorWalletAccounts(fullnode_dev)
        else:
            # in devnet, please manually create node account and send funds
            pass

        # setup each thread
        for i, t in enumerate(self.threads):
            t.setup()

    def run_threads(self):
        # start all threads
        for i, t in enumerate(self.threads):
            t.start()
        # join all threads
        for i, t in enumerate(self.threads):
            t.join()

    def stop_threads(self):
        # stop all threads
        for i, t in enumerate(self.threads):
            t.stop()

    def clean(self):
        # delete dirs and files
        print("cleaning loadtest...")
        shutil.rmtree(loadtest, ignore_errors=True)

threads = TestThreads()
