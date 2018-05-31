GOTOOLS = \
	  github.com/mitchellh/gox \
	  github.com/Masterminds/glide

all: test 

#
# Initialize the dependencies
#
setup: tools
	dep init

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
# Build a local copy
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
# Run any of the unit tests
#
utest:
	go test -v `glide novendor`

#
# System Testing, bring up everything
#
test: start
	echo "Running"

#
# Full Testing, bring up everything, send a transaction and then stop
#
full: stop reset install
	./scripts/startChain
	@sleep 9 # TODO: should have a way of knowing the chain is ready
	olclient send --address 127.0.0.1:46601
	@sleep 3
	./scripts/stopChain

#
# Start up a node, make sure it is stopped first
#
start: stop install
	./scripts/startChain

#
# System Testing, bring it all down
#
list: 
	pgrep -lf tendermint

#
# System Testing, bring it all down
#
stop: 
	./scripts/stopChain

#
# Reset all of the chains, back to genesis
#
reset:
	./scripts/resetChain 4

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
# targets that are not associated with files
#
#
.PHONY: all build install utest test tools update stop start full

