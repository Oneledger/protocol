# Entity/Relationships

## Notes

- All keys are numeric (how do we generate these?)
- Decentralized keys are hard to make unique
- All entities have a readable string name (or composite), that isn't the key
- Lists are ordered, sets are not
- Domain tables are persistent enumerations that are changable at runtime
- Permissions should be on both code and data, and they should be separated

## Model

Identity -- an individual in meatspace

	- User Id / External Node
	-------------------------

	- User Id
	-----------------------

	- User Name (naming collisions?)
	- Contact Info
	- default Account
	- Set of Accounts (Set of Wallets?)
	- Set of Trusted ChainNodes


Group -- a hiearchtical (or graph?) collection, fine-grained

	- Group Id
	-----------------------

	- Group Name (naming collisions?)
	- Public key?
	- Private key?

	- Set of Accessible Roles
	- Set of Accessible Data

	- Set of Identities and Groups (parents or children?)

Role -- the ability to view or perform an an action

	- Role Id
	-----------------------

	- Description
	- Action
	- Data?

Address -- Bitcoin temporary random address
Address -- Tendermint hash of public
Address -- Ethereum hash of public

Account -- information needed to access an account on a chain

	- Account Key - Different (OneLedger, Ethereum or Bitcoin)
	-----------------------

	- Name
	- Chain 
	- Public key
	- Private key
	- Balance@height
	- Set of Active HTLC states

Wallet -- a group of accounts

	-----------------------
	- Description
	- Set of Accounts

Balance -- the last known balance of an account on a chain (leafs in UTXO)

	- Account Key
	-----------------------

	- Coin

Chain -- a decentralized system 

	- Chain Id (Name?)
	-----------------------

	- Unique access point (how to get to the chain)
	-----------------------

	- Set of accessible nodes (not all nodes)
	- schema?
	- features?

ChainNode -- a specific access point in a chain

	-----------------------

	- NodeName
	- Location
	- Permissions?

AccountKey
	- Can be an identity
	- Can be an account on a chain
	- Can be an temporary account on a chain

Inputs -- account and existing balance

	-----------------------


Outputs -- account and new balance

	-----------------------


Transaction -- an action sent to our chain (instruction, statement, expression, transaction)

	-----------------------

	- Transaction Type 
	- Transaction Data

TransactionDomain -- Types of available transactions

	- Send
	- ExternalSend
	- ExternalLock
	- Swap (N party)
	- Prepare
	- Commit
	- Forget

SwapInstruction -- synchronize two (N?) transactions on different chains

	- List of Send/ExternalSend
	- Sequence 
	-----------------------

	- Fee
	- Gas
	
Send -- a transaction on our chain (we don't need this anymore...)

	- Set of Inputs
	- Set of Outputs
	- Sequence Number (replay protection)
	-----------------------

	- Fee
	- Gas

ExternalSend -- a transction send to oneledger, bitcoin, ethereum, cosmos, aion, etc. 

	- Set of Inputs
	- Set of Outputs
	- Sequence Number (replay protection)
	-----------------------

	- Fee
	- Gas (Ethereum only?)

ExternalLock -- setup a locked asset on bitcoin, ethereum, cosmos, aion, etc. 

	- Set of Addresses
	- Sequence Number (replay protection)
	-----------------------

	- Fee
	- Gas (Ethereum only?)

Payment -- the cost of doing business

	-----------------------

	- Fee
	- Gas

Input -- accounts that go into the transaction

	-----------------------

	- Address
	- Amount

Output -- accounts after the transaction

	-----------------------

	- Address
	- Amount

Coin -- a fungible unit of account (so no key)

	-----------------------

	- Currency
	- Count

CurrencyDomain

	OLT
	BTC
	ETH
	USD
	CAD

ExternalSwap Steps:

	- create lockbox on opposite chain
	- sign other box on current chain 
	- verify both boxes
	- send pre-image to counterparty
	- read the opposite chain
	- open lockbox

