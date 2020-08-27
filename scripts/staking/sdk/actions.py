import re
import subprocess
from os.path import dirname

import sys
import time

from rpc_call import *


class Staking:
    def __init__(self, node):
        self.node = node
        args = ['olclient', 'show_node_id']
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(1)
        pubKey = output[0].split(",")[0].split(":")[1].strip()
        f = open(os.path.join(self.node, "consensus", "config", "node_key.json"), "r")
        contents = json.loads(f.read())
        privKey = contents['priv_key']['value']
        args = ['olclient', 'account', 'add', '--privkey', privKey, '--pubkey', pubKey, "--password", 'pass']
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(1)
        self.staking_address = output[1].split(":")[1].strip()[3:]

    def prepare(self, secs=1):
        # add new account to node wallet
        args = ['olclient', 'account', 'add', "--password", 'pass']
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(1)
        new_staking_address = output[1].split(":")[1].strip()[3:]

        # Fund this account for future staking
        args = ['olclient', 'send', '--party', self.staking_address, '--counterparty', new_staking_address, '--amount',
                '5000000', '--fee', '0.0001', "--password", 'pass']
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        return new_staking_address

    def stake(self, amount, expect_success, secs=1):

        args = ['olclient', 'delegation', 'stake', '--amount', amount, '--address',
                self.staking_address,
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
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

    def unstake(self, amount, expect_success, secs=1):
        args = ['olclient', 'delegation', 'unstake', '--amount', amount, '--address',
                self.staking_address,
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
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

    def withdraw(self, amount, expect_success, secs=1):
        args = ['olclient', 'delegation', 'withdraw', '--amount', amount, '--address',
                self.staking_address,
                '--password', 'pass']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
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

    def checkStatus(self, delegation_amount, withdrawable_amount, hasFifthLine, secs=1):
        args = ['olclient', 'delegation', 'status', '--address', self.staking_address]
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)
        # check return code
        if process.returncode != 0:
            print "olclient check status failed returncode"
            sys.exit(-1)

        # check effective delegation amount in output
        target_regex = r"Effective delegation amount: " + str(delegation_amount) + r"$"
        match = re.search(target_regex, output[1])
        if match is None:
            print "olclient check status failed Effective delegation amount"
            sys.exit(-1)

        # check withdrawable amount in output
        target_regex = r"Withdrawable amount: " + str(withdrawable_amount) + r"$"
        match = re.search(target_regex, output[2])
        if match is None:
            print "olclient check status failed Withdrawable amount"
            sys.exit(-1)

        # check pending matured amount in output
        if hasFifthLine:
            target_regex = r"Pending matured amount:$"
        else:
            target_regex = r"Pending matured amount: empty$"
        match = re.search(target_regex, output[3])
        if match is None:
            print "olclient check status failed pending matured amount"
            sys.exit(-1)

        # check if there is fifth line
        if (len(output) <= 4 and hasFifthLine) or (len(output) > 4 and (not hasFifthLine)):
            print "olclient check status failed"
            sys.exit(-1)

        print "################### olclient check status succeed"

    def checkValidatorSet(self, node_number, on_list, target_power, secs=1):
        args = ['olclient', 'validatorset', 'status']
        target = dirname(dirname(os.getcwd()))
        # set protocol root path as current path
        process = subprocess.Popen(args, cwd=self.node, stdout=subprocess.PIPE)
        process.wait()
        output = process.stdout.readlines()
        time.sleep(secs)

        # check return code
        if process.returncode != 0:
            print "olclient check validatorset failed returncode"
            sys.exit(-1)

        # check if on list

        target_regex = r"Name " + str(node_number) + r"$"
        match = re.search(target_regex, output[4])
        if ((match is None) and (on_list is True)) or ((match is not None) and (on_list is False)):
            print "olclient check validatorset failed validatorname"
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
