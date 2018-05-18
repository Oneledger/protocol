# Entity/Relationships

## Notes

- All keys are numeric (how do we generate these?)
- All entities have a readable string name (or composite)
- Lists are ordered, sets are not
- Domain tables are persistent enumerations that are changable at runtime

## Model

Identity -- an individual in meatspace

	- User Name (unqiue string)
	-----------------------

	- Contact Info (email address?)
	- Primary Account
	- Set of Secondary Accounts
	- Map of ChainNodes

Account -- information needed to access an account on a chain

	- Address (hash of pubkey)
	-----------------------

	- Chain
	- Public key
	- Private key
	- Owner (Identity or Group)

Balance -- the last known balance of an account on a chain (leafs in UTXO)

	- Address (hash of pubkey)
	-----------------------

	- Coin

Group -- a hiearchtical (or graph) collection, fine-grained

	- Group Name (unqiue string)
	-----------------------

	- Public key?
	- Private key?

	- Set of Roles
	- Set of Data

	- Set of Included Identities 
	- Set of Included Groups (children)

Role -- the ability to view data or perform an an action

	- Role Id
	-----------------------

	- Description
	- Action
	- Data?

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

	- Name
	- Location
	- Permissions?

Transaction -- an action sent to our chain (instruction, statement, expression, transaction)

	-----------------------

	- Transaction Type 
	- Transaction Data

TransactionDomain

	- Send
	- ExternalSend
	- ExternalLock
	- Swap (N party)
	- Prepare
	- Commit
	- Forget

SwapInstruction -- synchronize two (N?) transactions on different chains

	- List of ExternalSend
	- Sequence 
	-----------------------

	- Fee
	- Gas
	
Send -- a transaction on our chain

	- Set of Inputs
	- Set of Outputs
	- Sequence Number (replay protection)
	-----------------------

	- Fee
	- Gas

ExternalSend -- a transction send to bitcoin, ethereum, cosmos, aion, etc. 

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
