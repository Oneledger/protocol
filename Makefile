

export GO111MODULE=on

#
# Updata the dependencies
#
update:
	go mod vendor
#
# Build and install a copy in bin
#
install:
	go install -i github.com/Oneledger/protocol/cmd/...

# Enable the clevelDB
install_c:
	CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" go install -tags "cleveldb" github.com/Oneledger/protocol/cmd/...

#
# test with send transaction in loadtest
#
fulltest: reset
	@./scripts/testsend
	@./scripts/stopNodes

#
# Check out the running status
#
status:
	@./scripts/status

#
# install and restart the network
#
restart: install_c
	@./scripts/stopNodes
	@./scripts/startDev

#
# run unit tests on project packages
#
utest:
	rm -rf /tmp/OneLedger-accounts.db
	go test github.com/Oneledger/protocol/data \
		github.com/Oneledger/protocol/data/accounts \
		github.com/Oneledger/protocol/data/balance \
		github.com/Oneledger/protocol/data/keys \
		github.com/Oneledger/protocol/data/governance \
		github.com/Oneledger/protocol/data/rewards \
		github.com/Oneledger/protocol/data/network_delegation \
		github.com/Oneledger/protocol/data/evidence \
		github.com/Oneledger/protocol/action/transfer \
		github.com/Oneledger/protocol/serialize \
		github.com/Oneledger/protocol/utils \
		github.com/Oneledger/protocol/rpc \
		github.com/Oneledger/protocol/identity \
		github.com/Oneledger/protocol/app \
		github.com/Oneledger/protocol/action/staking \
		github.com/Oneledger/protocol/action/evidence \
		-coverprofile a.out

loadtest: reset
	python scripts/loadtest/run_tests.py

coverage:
	go tool cover -html=a.out -o cover.html

#
# run ons tests
#
onstest: reset
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/renew_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/update_domain.py
	@./scripts/stopNodes

#
# remove it once Jenkins side done
#
withdrawtest: reset
	@./scripts/stopNodes

#
# run governance tests
#

govtest: reset
	@./scripts/testsend
	python scripts/governance/createProposals.py
	python scripts/governance/fundProposals.py
	python scripts/governance/cancelProposals.py
	python scripts/governance/withdrawFunds.py
	python scripts/governance/voteProposals.py
	python scripts/governance/governanceCLI.py
	make reset
	@./scripts/testsend
	python scripts/governance/queryProposals.py
	python scripts/governance/getFundsByFunder.py
	python scripts/governance/queryProposalOptions.py
	@./scripts/stopNodes

#
# run evidence tests
#
evidence: reset
	python scripts/evidence/allegation_loadtest.py
	make reset_no_install
	python scripts/evidence/test_allegation_no.py
	make reset_no_install
	python scripts/evidence/test_allegation_yes.py
	make reset_no_install
	python scripts/evidence/test_release.py
	@./scripts/stopNodes

#
# run staking tests
#
stakingtest: reset
	python scripts/staking/self_staking.py
	@./scripts/stopNodes

transfertest: reset
	python scripts/transfer/testSendPool.py
	@./scripts/stopNodes

# run rewards tests
rewardtest: reset
	@./scripts/testsend
	python scripts/reward/testRewards.py
	make reset
	python scripts/reward/testWithdraw.py
	python scripts/reward/listRewards.py
	@./scripts/stopNodes

delegationtest: reset
	python scripts/network_delegation/networkUndelegate.py
	python scripts/network_delegation/addNetworkDelegation.py
	make reset
	python scripts/network_delegation/withdrawRewards.py
	python scripts/network_delegation/reinvestRewards.py
# 	make reset
# 	python scripts/network_delegation/finalizeRewards.py
# 	make reset
# 	python scripts/network_delegation/withdrawDelegation.py
	@./scripts/stopNodes

alltest: reset
	@./scripts/testsend
	@./scripts/getValidators
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/purchase_expired.py
	python scripts/ons/create_delete_subdomain.py
	python scripts/ons/renew_domain.py
	python scripts/txTypes/listTxTypes.py
	@./scripts/stopNodes

reset: install_c
	make reset_no_install
# 	@./scripts/testapply
# 	@./scripts/testsend

reset_no_install:
	@./scripts/stopNodes
	@./scripts/resetDev
	@./scripts/startDev

resetInvalidValues: install_c
	@./scripts/stopNodes
	@./scripts/resetDev_invalidValues
	@./scripts/startDev
# 	@./scripts/testapply
# 	@./scripts/testsend

resetMain: install_c
	@./scripts/stopNodes
	@./scripts/resetMainnet
	@./scripts/startMainnet

rpcAuthtest: install_c
	@./scripts/stopNodes
	@./scripts/resetDev
	python scripts/rpcAuth/setup.py
	@./scripts/startDev
	python scripts/rpcAuth/rpcTestAuth.py
	@./scripts/stopNodes

stop:
	@./scripts/stopNodes

start:
	@./scripts/startDev


save: reset
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/renew_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/update_domain.py
	@./scripts/stopNodes
	make install_c
	@./scripts/saveState
	@./scripts/startDev

testData: 
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/renew_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/update_domain.py


updatetest: reset
	python scripts/governance/optUpdate.py
	@./scripts/testsend
	python scripts/reward/testWithdraw.py
	python scripts/reward/listRewards.py
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/renew_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/update_domain.py
	python scripts/governance/createProposals.py
	python scripts/governance/fundProposals.py
	python scripts/governance/cancelProposals.py
	python scripts/governance/voteProposals.py
	python scripts/governance/governanceCLI.py
	python scripts/governance/optTestCatchup.py
	python scripts/governance/optValidatorStaking.py
	@./scripts/stopNodes
