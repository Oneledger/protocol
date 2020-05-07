

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
		github.com/Oneledger/protocol/serialize \
		github.com/Oneledger/protocol/utils \
		github.com/Oneledger/protocol/rpc \
		github.com/Oneledger/protocol/identity \
		github.com/Oneledger/protocol/app \
		github.com/Oneledger/protocol/action/staking \
		-coverprofile a.out

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
# run ons tests
#
withdrawtest: reset
	@./scripts/testsend
	@./scripts/testsend
	python scripts/reward/withdraw.py
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
	@./scripts/stopNodes

alltest: reset
	@./scripts/testsend
	@./scripts/getValidators
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/purchase_expired.py
	python scripts/ons/create_delete_subdomain.py
	python scripts/ons/renew_domain.py
	python scripts/reward/withdraw.py
	python scripts/txTypes/listTxTypes.py
	@./scripts/stopNodes

reset: install_c
	@./scripts/stopNodes
	@./scripts/resetDev
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

save:
	@./scripts/stopNodes
	make install
	@./scripts/saveState
	@./scripts/startDev
