GOTOOLS = \
	  github.com/mitchellh/gox \
	  github.com/Masterminds/glide

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
	@echo "==== Running dep ensure, this may take awhile... ===="
	dep ensure
	@echo "==== Done setting up dependencies! ===="

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
list:
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
# Full Testing, bring up everything, send a transaction and then stop
#
swaptest: reset install
	@./tests/swap.sh
	@./tests/identity.sh
	@./tests/accounts.sh
	@./tests/utxo.sh
	@./scripts/stopOneLedger

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
	@./scripts/stopOneLedger

#
# Reset all of the chains, back to genesis
#
reset: check stop
	@./scripts/resetOneLedger 4

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
# targets that are not associated with any files
#
.PHONY: all build install utest test tools update stop start full
