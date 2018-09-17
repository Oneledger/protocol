# Local Testnet with Docker Compose

## System Requirements

+ [Docker](https://docs.docker.com/engine/installation/)
+ [Docker-Compose](https://docs.docker.com/compose/install/)

## Running the local testnet

The local testnet starts 4 fullnodes through Docker-Compose.

From the `protocol/` directory, run the following make targets:

*If* you don't have the tendermint and the protocol binaries installed, you can install them with:

```make build```

Prepare the docker image for running a fullnode:

```make docker-fullnode```

Turn on the local testnet with Docker Compose:

```make local-testnet```

Each node exposes their RPC ports to the host machine through 26611 (David), 26711 (Alice), 26811 (Bob), 26911 (Carol).

Turn the testnet off:

```make local-testnet-down```

You can start a single node that joins the local testnet by running:

```make run-singlenode```

