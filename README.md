# OneLedger MVP

Initial version for OneLedger's protocol, focusing on cross-chain transactions and identity management for the two different chains, Ethereum and Bitcoin.

The main part of the system is the blockchain full-node, which is currently an ABCi app based on lowlevel parts of the Tendermint consensus library. As development continues we will replace the Tendermint components with our own multi-level consensus engine, in order to achieve the advanced functionality we have our lined in our Whitepaper.

## Components


### Node

	- Implements a fullnode with Tendermint consensus (PoS)

	- Supports cross-chain transactions (as synchronized swaps on Hashed Timelocks)
		- Etheruem
		- Bitcoin

	- Supports Identity Management for OneLedger, Ethereum and Bitcoin

	- See README.md in node subdirectory for more details

### Consensus

	- Preliminary proof-of stake consensus engine, single-layer
	
	- Post-MVP core development

	
### Lite Chain

	- Simple chain and protocol for internal testing
	
	- Implements a miner with Satosi consenses (PoW)
	
