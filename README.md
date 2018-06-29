# OneLedger Protocol (MVP)

Initial version for the OneLedger Protocol, focusing on cross-chain transactions and identity management for the two different chains, Ethereum and Bitcoin.

The main part of the system is the blockchain full-node, which is currently an ABCi app based on low-level parts of the Tendermint consensus library. As development continues we will replace the Tendermint components with our own multi-level consensus engine, in order to achieve the advanced functionality we have outlined in our Whitepaper.

## Table of Contents

* [Getting Started](#getting-started)
  * [Quick Start (Docker)](#quick-start-docker)
  * [Manual Compilation](#manual-compilation)
    * [System Requirements](#system-requirements)
    * [Install](#install)
  * [Make Targets](#make-targets)
* [Components](#components)
  * [Node](#node)
  * [Consensus](#consensus)
  * [Lite Chain](#lite-chain)
* [License](#license)

## Getting Started

These instructions will get a copy of the OneLedger Protocol up and running on your local machine.

### Quick Start (Docker)

The OneLedger Protocol can be started as a Docker container. This is the easiest way to set up the OneLedger Protocol on your computer. You'll need to install [Docker CE](https://docs.docker.com/install/) on your computer to use this method.

Build the protocol image:
```
$ git clone https://github.com/Oneledger/protocol.git && cd protocol
$ docker build -t oneledger/protocol .
```

Start the protocol container (runs through your shell interactively):
```
$ docker run -it oneledger/protocol
```

See [Make Targets](#make-targets) to see a list of `make` commands you can run to interact with the OneLedger Protocol.

### Manual Compilation

#### System Requirements

Ensure your system meets the following requirements:

* Operating System must be one of the following:
  * macOS
    * Ensure [Xcode Developer Tools](https://developer.apple.com/xcode/) is installed
  * A Debian-based Linux distribution
    * *Ubuntu* is recommended for the smoothest install, otherwise you will have to set up your distribution for installing PPAs
* [Go](https://golang.org/) version 1.7 or higher
  * Need an explicit [GOPATH](https://github.com/Oneledger/protocol/wiki/Environment-Variables#setting-up-the-gopath) environment variable set
* [git](https://git-scm.com/)

#### Install

Before running any install scripts, ensure your GOPATH is set up explicitly on your user account. Visit the [Setting up the GOPATH](https://github.com/Oneledger/protocol/wiki/Environment-Variables#setting-up-the-gopath) page on the wiki for more help.

First clone the repository for the OneLedger Protocol with `go get`:

```
$ go get github.com/Oneledger/protocol
```

Before running any scripts, you'll need to set up the required [environment variables](https://github.com/Oneledger/protocol/wiki/Environment-Variables) properly.

Install the required dependencies:

```
$ cd "$OLROOT/protocol/node"
$ make setup
```

General scripts for running the OneLedger Protocol are inside [node/scripts](node/scripts).

### Make Targets

If everything is set up properly, you can begin testing the OneLedger Protocol with the `make` targets provided. Run the following scripts from inside the `$OLREPO/node` directory:

| Target | Description |
| --- | --- |
| `make test` | Tests system initialization, brings up everything |
| `make swaptest` | Test swap mechanics between BTC and ETH |
| `make fulltest` | Does a full test, makes use of test scripts in the `/tests` folder |

## Components

### Node

* Implements a fullnode with Tendermint consensus PoS (Proof-of-Stake)
* Supports cross-chain transactions (as synchronized swaps on Hashed Timelocks)
  * Etheruem
  * Bitcoin
* Supports Identity Management for OneLedger, Ethereum and Bitcoin
* See the `node` README at [node/README.md](node/README.md) for more details

### Consensus

* Preliminary PoS consensus engine, single-layer
* Post-MVP core development

### Lite Chain

* Simple chain and protocol for internal testing
* Implements a miner with Satoshi consensus PoW (Proof-of-Work)

## License

The OneLedger Protocol is released under the terms of the Apache 2.0 license. See [LICENSE.md](LICENSE.md) for more details or visit https://www.apache.org/licenses/LICENSE-2.0.html.
