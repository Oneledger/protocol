[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=master)](https://travis-ci.org/Oneledger/protocol)
[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=release)](https://travis-ci.org/Oneledger/protocol)
[![Build Status](https://travis-ci.org/Oneledger/protocol.svg?branch=develop)](https://travis-ci.org/Oneledger/protocol)
# OneLedger Protocol 

This repo is the OneLedger blockchain full-node, which is currently an ABCi app based on low-level parts of the Tendermint consensus library

## Table of Contents

* [Getting Started](#getting-started)
    * [Development](#Development)
    * [Make Targets](#make-targets)
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

First clone the repository for the OneLedger Protocol (pick any folder and clone it here):

```
$ git clone github.com/Oneledger/protocol
```

Install the required dependencies:

```
$ cd ./protocol
$ make install
```

General scripts for running the OneLedger Protocol are inside `./protocol/scripts`

### Make Targets

If everything is set up properly, you can begin testing the OneLedger Protocol with the `make` targets provided. Run the following scripts inside the `./protocol/` directory:

| Target | Description |
| --- | --- |
| `make install`| Build and install a copy of Oneledger protocol in bin |
| `make fulltest` | Test with send transaction in loadtest, makes use of test scripts in the `/tests` folder |
| `make status` | Lists all running nodes, Check out the running status|
| `make update` | Updata the dependencies |
| `make install_c` | Enable the clevelDB |


See [Make Targets](#make-targets) to see a list of `make` commands you can run to interact with the OneLedger Protocol.


## License

The OneLedger Protocol is released under the terms of the Apache 2.0 license. See [LICENSE.md](LICENSE.md) for more details or visit https://www.apache.org/licenses/LICENSE-2.0.html.
