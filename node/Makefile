GOTOOLS = \
	  github.com/mitchellh/gox \
	  github.com/Masterminds/glide

all: test 

#
# Initialize the dependencies
#
setup: tools
	cd setup
	./install.sh
	cd ..
	dep init
	dep ensure

#
# Update the dependencies if they have changed.
#
update: 
	dep ensure

#
# Get the tools if they are missing
#
tools:
	go get -u -v $(GOTOOLS)

#
# Build a local copy, not that useful
#
build:
	go build -o testFullnode github.com/Oneledger/protocol/node/cmd/fullnode
	go build -o testOlclient github.com/Oneledger/protocol/node/cmd/olclient

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
status: 
	./scripts/status

#
# Run any of the unit tests
#
utest:
	go test -v `glide novendor`

#
# System Testing, bring up everything
#
test: reset install start
	@echo "OneLedger Initialize"
	./scripts/stopOneLedger

#
# Full Testing, bring up everything, send a transaction and then stop
#
fulltest: reset install
	./tests/full.sh

#
# Full Testing, bring up everything, send a transaction and then stop
#
swaptest: reset install
	./tests/swap.sh
	./tests/identity.sh
	./tests/accounts.sh
	./tests/utxo.sh
	./scripts/stopOneLedger

#
# Start up a node, make sure it is stopped first
#
restart: reset install
	./scripts/startOneLedger

start: install
	./scripts/startOneLedger

#
# System Testing, bring it all down
#
stop: 
	./scripts/stopOneLedger

#
# Reset all of the chains, back to genesis
#
reset: stop
	./scripts/resetOneLedger 4

#
# Startup a tmux monitor
#
monitor:
	./scripts/monitor

#
# stop a tmux monitor
#
stopmon:
	./scripts/stopmon

#
# targets that are not associated with any files
#
.PHONY: all build install utest test tools update stop start full
