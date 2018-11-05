#
# Docker targets
#
#
build:
	./DOCKER/scripts/build.sh

# Prepare the base image for docker
docker-fullnode: prepare-volume
	docker build -t oneledger/fullnode -f ./DOCKER/Dockerfile --label oneledger --tag="oneledger/fullnode" .

# Prepare volume for persistence
prepare-volume: reset-volume
	./DOCKER/scripts/copy-bin.sh
	@./DOCKER/scripts/testnet.sh

# Reset the persistent data
reset-volume:
	@./DOCKER/scripts/reset.sh

# start up a local testnet
local-testnet:
	docker-compose up

# bring down a local testnet
local-testnet-down:
	docker-compose down

# Run a node to attach to the testnet
run-singlenode:
	docker run --volume $(CURDIR)/VOLUME:/home/oneledger/go/test --network="protocol_localnet" --env ID=Edwin --env OL_PEERS=192.167.11.1:26611,192.167.11.2:26611,192.167.11.3:26611,192.167.11.4:26611 --rm oneledger/fullnode

.PHONY: docker-fullnode run-singlenode local-testnet local-testnet-down prepare-volume reset-volume build
