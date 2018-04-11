GOTOOLS = \
	  github.com/mitchellh/gox \
	  github.com/Masterminds/glide

all: test 

tools:
	go get -u -v $(GOTOOLS)
build:
	go build github.com/Oneledger/prototype/node/cmd/fullnode

install:
	go install github.com/Oneledger/prototype/node/cmd/fullnode

test:
	go test -v `glide novendor`
.PHONY: all test build

