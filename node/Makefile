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
	go build -o testFullnode github.com/Oneledger/prototype/node/cmd/...

#
# Build and install a copy into bin
#
install:
	go install github.com/Oneledger/prototype/node/cmd/...

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
	./scripts/startNode
	sleep 9 # TODO: should have a wait of knowing the chain is ready
	olclient send
	sleep 3
	./scripts/stopNode

#
# Start up a node, make sure it is stopped first
#
start: stop install
	./scripts/startNode

#
# System Testing, bring it all down
#
stop: 
	./scripts/stopNode

#
# Reset all of the chains, back to genesis
#
reset:
	./scripts/resetChain

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

