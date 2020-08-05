package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cRouter Router
	app     parameterStruct
)

type parameterStruct struct {
	name string
}

func internalTX1(i interface{}) {
	param := i.(parameterStruct)
	fmt.Println("INTERNALTX1 ", param.name)
}
func internalTX2(i interface{}) {
	param := i.(parameterStruct)
	fmt.Println("INTERNALTX2 ", param.name)
}

func init() {
	app = parameterStruct{
		name: "Test App",
	}
	cRouter = NewRouter()
}

func TestRouter_AddBlockBeginner(t *testing.T) {
	err := cRouter.Add(BlockBeginner, internalTX1)
	assert.NoError(t, err)
	err = cRouter.Add(BlockBeginner, internalTX2)
	assert.NoError(t, err)
	err = cRouter.Add(BlockEnder, internalTX2)
	assert.NoError(t, err)
}

func TestRouter_IterateBlockBeginner(t *testing.T) {
	fmt.Println("Iterating Block Beginner")
	functionlist := cRouter.Iterate(BlockBeginner)
	assert.Len(t, functionlist, 2)
	for _, function := range functionlist {
		function(app)
	}
	fmt.Println("Iterating Block Ender")
	functionlist = cRouter.Iterate(BlockEnder)
	assert.Len(t, functionlist, 1)
	for _, function := range functionlist {
		function(app)
	}
}
