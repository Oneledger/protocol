package states

import (
	"../types"

)


/*
	This is the basic
 */
type BasicState interface {
	Change()
	Process()
}



type Propose struct {

}

type Prevote struct {

}

type Precommit struct{

}

type NewHeight struct{

}

type Commit struct {

}


type DriveChain struct{

}