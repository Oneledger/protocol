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

	err := cRouter.Add(BlockBeginner, cfunction{
		function:      internalTX1,
		functionParam: parameterStruct{name: "TEST"},
	})
	assert.NoError(t, err)
	err = cRouter.Add(BlockBeginner, cfunction{
		function:      internalTX2,
		functionParam: parameterStruct{name: "TEST2"},
	})
	assert.NoError(t, err)
	err = cRouter.Add(BlockEnder, cfunction{
		function:      internalTX2,
		functionParam: parameterStruct{name: "TEST3"},
	})
}

func TestRouter_IterateBlockBeginner(t *testing.T) {
	functionlist, err := cRouter.Iterate(BlockBeginner)
	assert.Len(t, functionlist, 2)
	assert.NoError(t, err)
	for _, function := range functionlist {
		function.function(function.functionParam)
	}
	functionlist, err = cRouter.Iterate(BlockEnder)
	assert.Len(t, functionlist, 1)
	assert.NoError(t, err)
	for _, function := range functionlist {
		function.function(function.functionParam)
	}
	functionlist, err = cRouter.Iterate(3)
	assert.Error(t, err)
}
