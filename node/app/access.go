/*
	Copyright 2017-2018 OneLedger
*/

package app

// Access to the local persistent databases
func (app Application) GetAdmin() interface{} {
	return app.Admin
}

// Access to the local persistent databases
func (app Application) GetStatus() interface{} {
	return app.Status
}

// Access to the local persistent databases
func (app Application) GetIdentities() interface{} {
	return app.Identities
}

// Access to the local persistent databases
func (app Application) GetAccounts() interface{} {
	return app.Accounts
}

// Access to the local persistent databases
func (app Application) GetBalances() interface{} {
	return app.Balances
}

func (app Application) GetChainID() interface{} {
	return ChainId
}

func (app Application) GetEvent() interface{} {
	return app.Event
}

func (app Application) GetContract() interface{} {
	return app.Contract
}
