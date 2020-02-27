

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
	CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" go install -tags "gcc" github.com/Oneledger/protocol/cmd/...

#
# test with send transaction in loadtest
#
fulltest: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
	@./scripts/testsend
	@./scripts/stopDev

#
# Check out the running status
#
status:
	@./scripts/status



#
# install and restart the network
#
restart: install
	@./scripts/stopDev
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
		github.com/Oneledger/protocol/serialize \
		github.com/Oneledger/protocol/utils \
		github.com/Oneledger/protocol/rpc \
		github.com/Oneledger/protocol/identity \
		github.com/Oneledger/protocol/app \
		-coverprofile a.out

coverage:
	go tool cover -html=a.out -o cover.html


#
# run apply validator tests
#
applytest: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
	@./scripts/testapply
	@./scripts/getValidators
	@./scripts/stopDev

#
# run ons tests
#
onstest: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/renew_domain.py
	python scripts/ons/buy_sell_domain.py
	@./scripts/stopDev

#
# run ons tests
#
withdrawtest: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
	@./scripts/testsend
	@./scripts/testsend
	python scripts/reward/withdraw.py
	@./scripts/stopDev


alltest: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
	@./scripts/testsend
	@./scripts/testapply
	@./scripts/getValidators
	@./scripts/testsend
	python scripts/ons/create_domain.py
	python scripts/ons/create_sub_domain.py
	python scripts/ons/buy_sell_domain.py
	python scripts/ons/purchase_expired.py
	python scripts/ons/create_delete_subdomain.py
	python scripts/reward/withdraw.py
	@./scripts/stopDev


reset: install
	@./scripts/stopDev
	@./scripts/resetDev
	@./scripts/startDev
# 	@./scripts/testapply
# 	@./scripts/testsend


rpcAuthtest: install
	@./scripts/stopDev
	@./scripts/resetDev
	python scripts/rpcAuth/setup.py
	@./scripts/startDev
	python scripts/rpcAuth/rpcTestAuth.py
	@./scripts/stopDev


stop:
	@./scripts/stopDev


start:
	@./scripts/startDev
