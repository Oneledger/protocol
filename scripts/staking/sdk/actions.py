from rpc_call import *
import subprocess
import time
import sys
import re
from os.path import dirname
import io


class Staking:
    def __init__(self, node):
        self.node = node

    def prepare(self, secs=1):
        args = ['olclient', 'loadtest', '--root', self.node, '--threads', '2', '--interval', '10', '--max-tx', '20']

        # set protocol root path as current path
        target = dirname(dirname(os.getcwd()))
        process = subprocess.Popen(args, cwd=target)
        process.wait()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient prepare failed"
            sys.exit(-1)
        print "################### olclient prepare succeed"

    def stake(self, amount, stake_address, expect_success, secs=1):
        args = ['olclient', 'delegation', 'stake', '--root', self.node, '--amount', amount, '--address',
                stake_address[3:],
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=target, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient stake failed"
            sys.exit(-1)
        # check for keyword "Failed" in stdout
        target_regex = r"Failed"
        match = re.search(target_regex, output[1])
        if (match is None) and (expect_success is False):
            print "olclient stake succeed, but it should fail!"
            sys.exit(-1)
        elif (match is not None) and (expect_success is True):
            print "olclient stake failed"
            sys.exit(-1)

        print "################### olclient stake succeed or failed as expected"

    def unstake(self, amount, stake_address, expect_success, secs=1):
        args = ['olclient', 'delegation', 'unstake', '--root', self.node, '--amount', amount, '--address',
                stake_address[3:],
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=target, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient unstake failed"
            sys.exit(-1)

        # check for keyword "Failed" in stdout
        target_regex = r"Failed"
        match = re.search(target_regex, output[1])
        if (match is None) and (expect_success is False):
            print "olclient stake succeed, but it should fail!"
            sys.exit(-1)
        elif (match is not None) and (expect_success is True):
            print "olclient stake failed"
            sys.exit(-1)

        print "################### olclient unstake succeed or failed as expected"

    def withdraw(self, amount, stake_address, expect_success, secs=1):
        args = ['olclient', 'delegation', 'withdraw', '--root', self.node, '--amount', amount, '--address',
                stake_address[3:],
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=target, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check for keyword "Failed" in stdout
        target_regex = r"Failed"
        match = re.search(target_regex, output[1])
        if (match is None) and (expect_success is False):
            print "olclient withdraw succeed, but it should fail!"
            sys.exit(-1)
        elif (match is not None) and (expect_success is True):
            print "olclient withdraw failed"
            sys.exit(-1)

    def checkStatus(self, stake_address, delegation_amount, withdrawable_amount, hasFifthLine, secs=1):
        args = ['olclient', 'delegation', 'status', '--root', self.node, '--address', stake_address[3:]]
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=target, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient check status failed"
            sys.exit(-1)
        print output

        # check effective delegation amount in output
        target_regex = r"Effective delegation amount: " + str(delegation_amount) + r"$"
        match = re.search(target_regex, output[1])
        if match is None:
            print "olclient check status failed"
            sys.exit(-1)

        # check withdrawable amount in output
        target_regex = r"Withdrawable amount: " + str(withdrawable_amount) + r"$"
        match = re.search(target_regex, output[2])
        if match is None:
            print "olclient check status failed"
            sys.exit(-1)

        # check pending matured amount in output
        if hasFifthLine:
            target_regex = r"Pending matured amount:$"
        else:
            target_regex = r"Pending matured amount: empty$"
        match = re.search(target_regex, output[3])
        if match is None:
            print "olclient check status failed"
            sys.exit(-1)

        # check if there is fifth line
        if (len(output) <= 4 and hasFifthLine) or (len(output) > 4 and (not hasFifthLine)):
            print "olclient check status failed"
            sys.exit(-1)

        print "################### olclient check status succeed"

    def checkValidatorSet(self, node_number, on_list, target_power, secs=1):
        args = ['olclient', 'validatorset', 'status', '--root', self.node]
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=target, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient check validatorset failed"
            sys.exit(-1)

        # check if on list
        target_regex = r"Name " + str(node_number) + r"$"
        match = re.search(target_regex, output[4])
        if ((match is None) and (on_list is True)) or ((match is not None) and (on_list is False)):
            print "olclient check validatorset failed"
            sys.exit(-1)

        # check target power in output
        if on_list is True:
            target_regex = r"Power " + str(target_power) + r"$"
            match = re.search(target_regex, output[3])
            if match is None:
                print "olclient check validatorset failed"
                sys.exit(-1)

        print "################### olclient check validatorset succeed"



def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]
