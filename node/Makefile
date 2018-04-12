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
	go build github.com/Oneledger/prototype/node/cmd/fullnode

#
# Build and install a copy into bin
#
install:
	go install github.com/Oneledger/prototype/node/cmd/fullnode

#
# Run any of the unit tests
#
utest:
	go test -v `glide novendor`

#
# System Testing, bring up everything
#
test:
	./scripts/startNode

#
# Reset all of the chains, back to genesis
#
reset:
	./scripts/resetChain

#
# targets that are not associated with files
#
.PHONY: all utest reset test install build tools update

