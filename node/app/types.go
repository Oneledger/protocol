/*
	Copywrite 2017-2018 OneLedger

	Declare all of the types used by the Application
*/
package app

import (
	"github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"
)

type Message []byte // Contents of a transaction
type Key []byte     // Database key

// ApplicationContext keeps all of the upper level global values.
type ApplicationContext struct {
	types.BaseApplication

	log log.Logger // inherited logger
	db  Datastore  // key/value database in memory
}

// NewApplicationContext initializes a new application
func NewApplicationContext(logger log.Logger) *ApplicationContext {
	return &ApplicationContext{
		log: logger,
		db:  *NewDatastore(),
	}
}
