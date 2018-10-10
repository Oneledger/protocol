#
# Docker targets
#

# Prepare the base image for docker
docker-fullnode: prepare-volume
	docker build -t oneledger/fullnode -f ./DOCKER/Dockerfile --label oneledger --tag="oneledger/fullnode" .

run-singlenode:
	docker run --volume $(CURDIR)/VOLUME:/home/oneledger/go/test --network="protocol_localnet" --env ID=Edwin --env OL_PEERS=192.167.11.1:26611,192.167.11.2:26611,192.167.11.3:26611,192.167.11.4:26611 --rm oneledger/fullnode

local-testnet:
	docker-compose up

local-testnet-down:
	docker-compose down

prepare-volume: reset-volume
	./DOCKER/scripts/copy-bin.sh
	@./DOCKER/scripts/testnet.sh

reset-volume:
	@./DOCKER/scripts/reset.sh

build:
	./DOCKER/scripts/build.sh

.PHONY: docker-fullnode run-singlenode local-testnet local-testnet-down prepare-volume reset-volume build
