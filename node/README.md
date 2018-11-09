# OneLedger Fullnode

A node in the OneLedger chain, which includes a process for Tendermint consensus and an ABCi application to handle processing all of the transactions.

Make targets are used for all development activities. 

$GOPATH/test is used for the working directory and logs (Tendermint files should also be redirected here).

## Set Environment

    source node/scripts/setEnvironment

## Setup

	make setup
	make update
	
## Building

	make build (builds a local copy)
	make install
	
## Testing
	make test (stops, installs and runs...)


