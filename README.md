# OneLedger Protocol 

We launched our testnet Chronos based on v0.8.1. To set up a node please refer to https://github.com/Oneledger/protocol/wiki/Chronos-Set-Up-Instructions-(v0.8.1). 
To check testnet status, refer to https://oneledger.network

This repo is the OneLedger blockchain full-node, which is currently an ABCi app based on low-level parts of the Tendermint consensus library

## Table of Contents

* [Components](#components)
  * [Node](#node)
* [License](#license)

## Components

### Node

* Implements a fullnode with Tendermint consensus PBFT
* Supports cross-chain transactions (as synchronized swaps on Hashed Timelocks)
  * Etheruem
  * Bitcoin
* Supports Identity Management for OneLedger, Ethereum and Bitcoin
* See the `node` README at [node/README.md](node/README.md) for more details


## License

The OneLedger Protocol is released under the terms of the Apache 2.0 license. See [LICENSE.md](LICENSE.md) for more details or visit https://www.apache.org/licenses/LICENSE-2.0.html.
