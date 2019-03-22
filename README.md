[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=master)](https://travis-ci.org/Oneledger/protocol)
[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=release)](https://travis-ci.org/Oneledger/protocol)
[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=develop)](https://travis-ci.org/Oneledger/protocol)
# OneLedger Protocol 

This repo is the OneLedger blockchain full-node, which is currently an ABCi app based on low-level parts of the Tendermint consensus library

## Table of Contents

* [Getting Started](#getting-started)
    * [Development](#Development)
    * [Join OneLedger Testnet](#join-oneledger-testnet)
* [License](#license)

## Getting Started

### Development

   These instructions will get a copy of the OneLedger Protocol up and running on your local machine.

#### System Requirements

Ensure your system meets the following requirements:

* Operating System must be one of the following:
  * macOS
    * Ensure [Xcode Developer Tools](https://developer.apple.com/xcode/) is installed
  * A Debian-based Linux distribution
    * *Ubuntu* is recommended for the smoothest install, otherwise you will have to set up your distribution for installing PPAs
* [Go](https://golang.org/) version 1.11 or higher
* [git](https://git-scm.com/)

#### Install

Before running any install scripts, ensure your GOPATH and the necessary environment variables are set up on your account. Visit the [Setting up the GOPATH](https://github.com/Oneledger/protocol/wiki/Environment-Variables#setting-up-the-gopath) page on the wiki for more help.

First clone the repository for the OneLedger Protocol:

```
$ mkdir -p $OLROOT
$ cd $OLROOT
$ git clone github.com/Oneledger/protocol
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
| `make status` | Lists all running nodes |
| `make monitor` | Start tmux session |
| `make stopmon` | Stops tmux session |

See [Make Targets](#make-targets) to see a list of `make` commands you can run to interact with the OneLedger Protocol.

### Join OneLedger Testnet 

* [Setup a node](https://github.com/Oneledger/protocol/wiki/Chronos-Set-Up-Instructions-v0.8.1)

* [Become Validator](https://github.com/Oneledger/protocol/wiki/Chronos-Validator-Instructions-v0.8.1)

* [Check the explorer](https://oneledger.network)


## License

The OneLedger Protocol is released under the terms of the Apache 2.0 license. See [LICENSE.md](LICENSE.md) for more details or visit https://www.apache.org/licenses/LICENSE-2.0.html.
