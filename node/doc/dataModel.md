# Entity/Relationships

## Notes:

	- All keys are numeric? (how do we generate these?)
	- Decentralized keys are hard to make unique
	- All entities have a readable string name (or composite), that isn't the key
	- Lists are ordered, sets are not
	- Domain tables are persistent enumerations that are changable at runtime
	- Permissions should be on both code and data, and they should be separated
	- Single keyed entities are stored in a key/value database.
	- Multi-keyed entries are similar, but have indicies which map the external key to the main one.
	- No key entries are not persistent
	- Domain entries are stored as a list under a constant

# Data Model

## Authentication: Users, Identities and Accounts

**Relationship**: user <<->> group <-> identity <->> wallet <->> account


**Group/User** -- a hiearchtical (or graph?) collection, fine-grained, a group can own a public/private key, and any members of that group can use it.

	- Group Id
	-----------------------
	- Group Name (naming collisions?)
	- Public key
	- Private key

	- Set of Accessible Roles
	- Set of Accessible Data

	- Set of Identities and Groups (parents or children?)

**Identity** -- an identity that is known across the chains, tied to one public key

	- User Id / Chain
	-------------------------
	- User Name (naming collisions?)
	- Public Key
	- Contact Info
	- Default Account (Wallet?)
	- Set of Accounts (Set of Wallets?)
	- Set of Trusted ChainNodes
	- Set of Trusted Channels


## Authorization

**Relationship**: identity <->> role, identity <->> data

**Role** -- the ability to view or perform an an action

	- Role Id
	-----------------------
	- Description
	- Actions
	- Data Access

## Low-level accounts and addresses

**Address** -- Polymorphic, chain-independent

	- Bitcoin temporary random address
	- Bitcoin script address
	- Tendermint hash of public key
	- Ethereum hash of public key

**AccountKey**

	- Can be an identity (wrong?)
	- Can be an account on a chain
	- Can be an temporary account on a chain

**Account** -- information needed to access an account on a chain

	- AccountKey 
	-----------------------
	- Name
	- Chain 
	- Public key
	- Private key
	- Balance@height
	- Set of Active HTLC states

**Wallet** -- a group of accounts, can be one-offs, includes the file structure and is protected by a passphrase. Needs to be located in a secure place.

	- Wallet Id
	-----------------------
	- Description
	- Set of Accounts

**Balance** -- the last known balance of an account on a chain (leafs in UTXO)

	- Account Key
	-----------------------
	- Coin
	- Date (?)

## External Chains

**Chain** -- a decentralized system that can be driver by the protocol. There can be many instances of the same code-base, like testnets and a mainnet. Forks are possible too.

	- Chain Id (Name?)
	-----------------------
	- Unique access point (how to get to the chain)
	-----------------------
	- Chain type
	- Chain Driver
	- Set of accessible nodes (not all nodes)
	- Schema
		- Basic Features
		- Extended Features

**ChainNode** -- a specific trusted node in a chain

	- Node Location
	-----------------------
	- NodeName
	- Location
	- Permissions?
	- Trust?

## Transactions (relative, active)

**Transaction** -- an action sent to our chain (instruction, statement, expression, transaction)

	- Transaction Type 
	- Transaction Data
	-----------------------

**Payment** -- the cost of doing business

	- Fee
	- Gas

**Input** -- accounts that go into the transaction

	- Address
	-----------------------
	- Amount

**Output** -- accounts after the transaction

	- Address
	-----------------------
	- Amount

**Coin** -- a fungible unit of account (so no key)

	- Currency
	- Count

**TransactionDomain** -- Types of available transactions

	- Send
	- ExternalSend
	- ExternalLock
	- Swap (N party)
	- Prepare
	- Commit
	- Forget

## Transactions

SwapInstruction -- synchronize two (N?) transactions on different chains

	- List of Send/ExternalSend
	- Sequence 
	-----------------------
	- Payment
	
Send -- a transaction on our chain (we don't need this anymore...)

	- Set of Inputs
	- Set of Outputs
	- Sequence Number (replay protection)
	-----------------------
	- Payment

ExternalSend -- a transction send to oneledger, bitcoin, ethereum, cosmos, aion, etc. 

	- Set of Inputs
	- Set of Outputs
	- Sequence Number (replay protection)
	-----------------------
	- Payment

ExternalLock -- setup a locked asset on bitcoin, ethereum, cosmos, aion, etc. 

	- Set of Addresses
	- Sequence Number (replay protection)
	-----------------------
	- Payment

CurrencyDomain

	OLT
	BTC
	ETH
	USD
	CAD

### Algorithms

*Algorithm*: **ExternalSwap** Steps:

	- test to see if this is me?
	- see if both sides of swap are transactions
	- create lockbox on primary chain
	- wait for other box to be ready
		- check
		- wait for other box to be ready
	- sign other box on secondary chain 
	- wait -- initialized...
		- check
		- wait -- initialized...
	- verify box on primary chain
	- verify box on secondary chain
	- wait -- verified...
		- check
		- wait -- verified...
	- send pre-image to counterparty via transaction
	- wait
	- read the opposite chain
	- open lockbox

