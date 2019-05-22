

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
# run unit tests on project packages
#
utest: 
	go test github.com/Oneledger/protocol/data \
		github.com/Oneledger/protocol/serialize \
		github.com/Oneledger/protocol/utils \
		-coverprofile a.out

coverage:
	go tool cover -html=a.out -o cover.html
