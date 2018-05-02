# OneLedger MVP

Initial version for OneLedger's protocol, focusing on cross-chain transactions and identity management for the different chains.

## Components


### Node

	- Implements a fullnode with Tendermint consensus (PoS)

	- Supports cross-chain transactions (as synchronized swaps on Hashed Timelocks)
		- Etheruem
		- Bitcoin

	- Supports Identity Management for OneLedger, Ethereum and Bitcoin

	- See README.md in node subdirectory for more details

### Cli
	- Implements a miner with Satosi consenses (PoW)
