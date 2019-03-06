package states

import (
	"../types"
)

/*
	Basic state representation of the consensus state machine
*/
type BasicState interface {
	Change()
	Process()
}

type Propose struct {
}

type Prevote struct {
}

type Precommit struct {
}

type NewHeight struct {
}

type Commit struct {
}

type DriveChain struct {
}
