#
# Docker images
#
docker-fullnode: prepare-volume
	docker build -t oneledger/fullnode -f ./DOCKER/Dockerfile --label oneledger --tag="oneledger/fullnode" .

run-fullnode:
	docker run --volume $(CURDIR)/VOLUME:/home/oneledger/go/test --env ID=David -it oneledger/fullnode

local-testnet:
	docker-compose up

local-testnet-down:
	docker-compose down

prepare-volume:
	@./DOCKER/scripts/copy-bin.sh
	@./DOCKER/scripts/testnet.sh

