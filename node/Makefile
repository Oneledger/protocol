GOTOOLS = \
	  github.com/mitchellh/gox \
	  github.com/Masterminds/glide \
	  github.com/golang/dep/cmd/dep \
	  github.com/golang/protobuf/protoc-gen-go

all: test

#
# Make sure the environment variables are setup correctly
#
check:
	@./checkEnv || exit 1

#
# Initialize the dependencies
#
setup: check tools
	./setup/install.sh
	#dep init
	echo "==== Running dep ensure, this may take awhile... ===="
	dep ensure
	echo "==== Done setting up dependencies! ===="

#
# Update the dependencies if they have changed.
#
update: check
	dep ensure

#
# Get the tools if they are missing
#
tools:
	go get -u -v $(GOTOOLS)

#
# Build and install a copy into bin
#
install:
	go install github.com/Oneledger/protocol/node/cmd/...

#
# List out all of the tendermint processes
#
show:
	pgrep -lf tendermint

#
# Check out the running status
#
status: check
	@./scripts/status

#
# Run any of the unit tests
#
utest:
	go test -v `glide novendor`

#
# System Testing, bring up everything
#
test: check reset install start
	@echo "OneLedger Initialize"
	@./scripts/stopOneLedger

#
# Full Testing, bring up everything, send a transaction and then stop
#
fulltest: reset install
	@./tests/full.sh

#
# Test validator, bring up everything, send a validator token and then stop
#
validatortest: reset install
	@./tests/validator.sh

#
# Swap Testing, swap BTC with ETH between nodes, do a full reset and rebuild each time
#
swaptest: resetDevnet install
	@./scripts/startBitcoin
	@./scripts/startEthereum
	@./tests/register.sh
	@./tests/swap.sh
	@./scripts/stopDevnet

#
# Smart Contract Testing, install and execute smart contracts
#
contracttest: stop reset start register
	@./tests/contract.sh


#
# Start up a node, make sure it is stopped first
#
restart: reset install
	@./scripts/startOneLedger

start: check install
	@./scripts/startOneLedger

#
# System Testing, bring it all down
#
stop: check
	@./scripts/stopDevnet

#
# Register the devnet users
#
register: check
	@./tests/register.sh

#
# Send money to the devnet users
#
testmint: check
	@./tests/testmint.sh

#
# List out the data for David
#
list: check
	@olclient list -c David

#
# Reset all of the chains, back to genesis
#
reset: check stop
	@./scripts/resetOneLedger

#
# Reset all of the chains, back to genesis
#
resetDevnet: check stop
	@./scripts/resetDevnet

#
# Startup a tmux monitor
#
monitor: check
	@./scripts/monitor

#
# stop a tmux monitor
#
stopmon: check
	@./scripts/stopmon

#
# Compile protocol buffers
#
protobuf:
	protoc ./sdk/pb/base.proto --go_out=plugins=grpc:./

#
# targets that are not associated with any files
#
.PHONY: all build install utest test tools update stop start full

