import hashlib
import sys

from rpc_call import *

# Proposal Types
ProposalTypeInvalid = 0xEE
ProposalTypeConfigUpdate = 0x20
ProposalTypeCodeChange = 0x21
ProposalTypeGeneral = 0x22

# Proposal Status
ProposalStatusFunding = 0x23
ProposalStatusVoting = 0x24
ProposalStatusCompleted = 0x25

# Proposal Outcome
ProposalOutcomeInProgress = 0x26
ProposalOutcomeInsufficientFunds = 0x27
ProposalOutcomeInsufficientVotes = 0x28
ProposalOutcomeCancelled = 0x29
ProposalOutcomeCompleted = 0x30

# Proposal States
ProposalStateInvalid = 0xEE
ProposalStateActive = 0x31
ProposalStatePassed = 0x32
ProposalStateFailed = 0x33

# Vote Opinions
OPIN_POSITIVE = 0x1
OPIN_NEGATIVE = 0x2
OPIN_GIVEUP = 0x3
OpinMap = {OPIN_POSITIVE: 'YES', OPIN_NEGATIVE: 'NO', OPIN_GIVEUP: 'GIVEUP'}


class bcolors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'


class Proposal:
    def __init__(self, pid, pType, description, headline, proposer, init_fund):
        self.pid = pid
        self.pty = pType
        self.headline = headline
        self.des = description
        self.proposer = proposer
        self.init_fund = init_fund

    def _create_proposal(self):
        req = {
            "proposalId": self.pid,
            "headline": self.headline,
            "description": self.des,
            "proposer": self.proposer,
            "proposalType": self.pty,
            "initialFunding": {
                "currency": "OLT",
                "value": convertBigInt(self.init_fund),
            },
            "configUpdate": {
                "bitcoinChainDriverOption": {
                    "TotalSupply": "1",
                    "BlockConfirmation": 6,
                    "ChainType": "TESTINGUPDATE",
                    "TotalSupplyAddr": "oneledgerSupplyAddress"
                },
                "ethchaindriverOption": {
                    "ERCContractAddress": "0x0000000000000000000000000000000000000000",
                    "ContractAddress": "0x0000000000000000000000000000000000000000",
                    "BlockConfirmation": 0,
                    "ERCContractABI": "",
                    "TotalSupply": "",
                    "TokenList": [{
                        "TokName": "string",
                        "TokAddr": "0x0000000000000000000000000000000000000000",
                        "TokAbi": "string",
                        "TokTotalSupply": "string",
                    }],
                    "ContractABI": "",
                    "TotalSupplyAddr": ""
                },
                "feeOption": {
                    "feeCurrency": {
                        "decimal": 18,
                        "unit": "nue",
                        "id": 0,
                        "chain": 0,
                        "name": "OLT"
                    },
                    "minFeeDecimal": 9
                },
                "rewardOptions": {
                    "rewardPoolAddress": "rewardpool",
                    "rewardInterval": 150
                },
                "propOptions": {
                    "configUpdate": {
                        "passedFundDistribution": {
                            "burn": 18,
                            "executionCost": 18,
                            "bountyPool": 10,
                            "validators": 18,
                            "proposerReward": 18,
                            "feePool": 18
                        },
                        "failedFundDistribution": {
                            "burn": 10,
                            "executionCost": 20,
                            "bountyPool": 50,
                            "validators": 10,
                            "proposerReward": 0,
                            "feePool": 10
                        },
                        "fundingGoal": "10000000000",
                        "proposalExecutionCost": "executionCostConfig",
                        "votingDeadline": 12,
                        "initialFunding": "1000000000",
                        "fundingDeadline": 12,
                        "passPercentage": 51
                    },
                    "bountyProgramAddr": "oneledgerBountyProgram",
                    "codeChange": {
                        "passedFundDistribution": {
                            "burn": 18,
                            "executionCost": 18,
                            "bountyPool": 10,
                            "validators": 18,
                            "proposerReward": 18,
                            "feePool": 18
                        },
                        "failedFundDistribution": {
                            "burn": 10,
                            "executionCost": 20,
                            "bountyPool": 50,
                            "validators": 10,
                            "proposerReward": 0,
                            "feePool": 10
                        },
                        "fundingGoal": "10000000000",
                        "proposalExecutionCost": "executionCostCodeChange",
                        "votingDeadline": 12,
                        "initialFunding": "1000000000",
                        "fundingDeadline": 12,
                        "passPercentage": 51
                    },
                    "general": {
                        "passedFundDistribution": {
                            "burn": 18,
                            "executionCost": 18,
                            "bountyPool": 10,
                            "validators": 18,
                            "proposerReward": 18,
                            "feePool": 18
                        },
                        "failedFundDistribution": {
                            "burn": 10,
                            "executionCost": 20,
                            "bountyPool": 50,
                            "validators": 10,
                            "proposerReward": 0,
                            "feePool": 10
                        },
                        "fundingGoal": "10000000000",
                        "proposalExecutionCost": "executionCostGeneral",
                        "votingDeadline": 12,
                        "initialFunding": "1000000000",
                        "fundingDeadline": 12,
                        "passPercentage": 51
                    }
                },
                "onsOptions": {
                    "currency": "OLT",
                    "firstLevelDomains": [
                        "ol"
                    ],
                    "baseDomainPrice": "1000000000000000000000",
                    "perBlockFees": "100000000000000"
                },
            },
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 400000,
        }
        resp = rpc_call('tx.CreateProposal', req)
        return resp["result"]["rawTx"]

    def send_create(self):
        # createTx
        raw_txn = self._create_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                self.pid = self.get_encoded_pid()
                self.txHash = "0x" + result["txHash"]
                print "################### proposal created: " + self.pid

    def get_encoded_pid(self):
        hash_handler = hashlib.md5()
        hash_handler.update(self.pid)
        hash_val = hash_handler.digest()
        return hash_val.encode('hex')

    def tx_created(self):
        resp = tx_by_hash(self.txHash)
        return resp["result"]["tx_result"]


class ProposalFund:
    def __init__(self, pid, value, address):
        self.pid = pid
        self.value = value
        self.funder = address

    def _fund_proposal(self):
        req = {
            "proposalId": self.pid,
            "fundValue": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "funderAddress": self.funder,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }

        resp = rpc_call('tx.FundProposal', req)
        return resp["result"]["rawTx"]

    def send_fund(self):
        # create Tx
        raw_txn = self._fund_proposal()

        # sign Tx
        signed = sign(raw_txn, self.funder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal funded: " + self.pid
                return result["txHash"]


class ProposalCancel:
    def __init__(self, pid, proposer, reason):
        self.pid = pid
        self.proposer = proposer
        self.reason = reason

    def _cancel_proposal(self):
        req = {
            "proposalId": self.pid,
            "proposer": self.proposer,
            "reason": self.reason,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }

        resp = rpc_call('tx.CancelProposal', req)
        return resp["result"]["rawTx"]

    def send_cancel(self):
        # create Tx
        raw_txn = self._cancel_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                print "################### failed to cancel proposal: " + self.pid
                return False
            else:
                print "################### proposal canceled: " + self.pid
                return True
        else:
            print "################### failed to cancel proposal: " + self.pid
            return False


class ProposalVote:
    def __init__(self, pid, opinion, url, address):
        self.pid = pid
        self.opin = opinion
        self.voter = url
        self.address = address

    def _vote_proposal(self):
        req = {
            "proposalId": self.pid,
            "opinion": self.opin,
            "address": self.address,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.VoteProposal', req, self.voter)
        result = resp["result"]
        return result["rawTx"], result['signature']['Signed'], result['signature']['Signer']

    def send_vote(self):
        # create and let validator sign Tx
        raw_txn, signed_0, signer_0 = self._vote_proposal()

        # payer sign Tx
        res = sign(raw_txn, self.address)
        signed_1 = res['signature']['Signed']
        signer_1 = res['signature']['Signer']

        # signatures
        sig0 = {"Signer": signer_0, "Signed": signed_0}
        sig1 = {"Signer": signer_1, "Signed": signed_1}
        sigs = [sig1, sig0]

        # broadcast Tx
        result = broadcast_commit_mtsig(raw_txn, sigs)

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal voted:" + self.pid + "opinion: " + OpinMap[self.opin]
                return result["txHash"]


class ProposalFundsWithdraw:
    def __init__(self, pid, funder, value, beneficiary):
        self.pid = pid
        self.funder = funder
        self.value = value
        self.benefi = beneficiary

    def _withdraw_funds(self, funder_address):
        req = {
            "proposalId": self.pid,
            "funderAddress": funder_address,
            "withdrawValue": {
                "currency": "OLT",
                "value": convertBigInt(self.value),
            },
            "beneficiaryAddress": self.benefi,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.WithdrawProposalFunds', req)
        print resp
        return resp["result"]["rawTx"]

    def withdraw_fund(self, contr_address):
        # create Tx
        raw_txn = self._withdraw_funds(contr_address)

        # sign Tx
        signed = sign(raw_txn, self.funder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                print bcolors.FAIL + "################### proposal funds withdraw failed:" + result[
                    "log"] + bcolors.ENDC
                sys.exit(-1)
            else:
                print "################### proposal funds withdrawn:" + self.pid
                return result["txHash"]
        else:
            print bcolors.FAIL + "################### proposal funds withdraw failed:" + result["error"][
                "message"] + bcolors.ENDC
            sys.exit(-1)

    def withdraw_fund_should_fail(self, contr_address):
        # create Tx
        raw_txn = self._withdraw_funds(contr_address)

        # sign Tx
        signed = sign(raw_txn, self.funder)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                print bcolors.FAIL + "################### proposal funds withdraw failed:" + result[
                    "log"] + bcolors.ENDC
                return result["txHash"]
            else:
                sys.exit(-1)
        else:
            print bcolors.FAIL + "################### proposal funds withdraw failed:" + result["error"][
                "message"] + bcolors.ENDC


class ProposalFinalize:
    def __init__(self, pid, address):
        self.pid = pid
        self.proposer = address

    def _finalize_proposal(self):
        req = {
            "proposalId": self.pid,
            "proposer": self.proposer,
            "gasPrice": {
                "currency": "OLT",
                "value": "1000000000",
            },
            "gas": 40000,
        }
        resp = rpc_call('tx.FinalizeProposal', req)
        return resp["result"]["rawTx"]

    def send_finalize(self):
        # create Tx
        raw_txn = self._finalize_proposal()

        # sign Tx
        signed = sign(raw_txn, self.proposer)

        # broadcast Tx
        result = broadcast_commit(raw_txn, signed['signature']['Signed'], signed['signature']['Signer'])

        if "ok" in result:
            if not result["ok"]:
                sys.exit(-1)
            else:
                print "################### proposal finalized: " + self.pid
                return result["txHash"]


def addresses():
    resp = rpc_call('owner.ListAccountAddresses', {})
    return resp["result"]["addresses"]


def sign(raw_tx, address):
    resp = rpc_call('owner.SignWithAddress', {"rawTx": raw_tx, "address": address})
    return resp["result"]


def broadcast_commit(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxCommit', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    print resp
    if "result" in resp:
        return resp["result"]
    else:
        return resp


def broadcast_commit_mtsig(raw_tx, sigs):
    resp = rpc_call('broadcast.TxCommitMtSig', {
        "rawTx": raw_tx,
        "signatures": sigs,
    })
    print resp
    if "result" in resp:
        return resp["result"]
    else:
        return resp


def broadcast_sync(raw_tx, signature, pub_key):
    resp = rpc_call('broadcast.TxSync', {
        "rawTx": raw_tx,
        "signature": signature,
        "publicKey": pub_key,
    })
    return resp["result"]


def query_proposals(prefix, proposer="", proposalType=ProposalTypeInvalid):
    req = {
        "state": prefix,
        "proposer": proposer,
        "proposalType": proposalType,
    }

    resp = rpc_call('query.ListProposals', req)
    print resp
    result = resp["result"]
    print json.dumps(resp, indent=4)
    return result["proposalStats"]


def query_proposal(proposal_id):
    req = {
        "proposalId": proposal_id,
    }
    resp = rpc_call('query.ListProposal', req)
    stat = resp["result"]["proposalStats"][0]
    print json.dumps(resp, indent=4)
    return stat["proposal"], stat["funds"]


def query_governanceState(height):
    req = {
        "height": height,
    }
    resp = rpc_call('query.GetGovernanceOptionsForHeight', req)
    result = resp["result"]

    # print json.dumps(resp, indent=4)
    return result["govOptions"]["bitcoinChainDriverOption"]["ChainType"]


def query_balance(address):
    req = {"address": address}
    resp = rpc_call('query.Balance', req)
    print json.dumps(resp, indent=4)
    return resp["result"]
