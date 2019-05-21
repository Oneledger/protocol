

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
	go install github.com/Oneledger/protocol/cmd/...

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