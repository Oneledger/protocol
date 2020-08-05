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
	param := i.(*parameterStruct)
	fmt.Println("INTERNALTX1 ", param.name)
}
func internalTX2(i interface{}) {
	param := i.(*parameterStruct)
	fmt.Println("INTERNALTX2 ", param.name)
}

func init() {
	app = parameterStruct{
		name: "Test App",
	}
	cRouter = NewRouter()
}

func TestRouter_AddBlockBeginner(t *testing.T) {
	err := cRouter.AddBlockBeginner("InternalTX1", internalTX1)
	assert.NoError(t, err)
	err = cRouter.AddBlockBeginner("InternalTX1", internalTX2)
	assert.NoError(t, err)
}

func TestRouter_IterateBlockBeginner(t *testing.T) {
	fmt.Println("Testing Iteration")
	functionlist := cRouter.IterateBlockBeginner()
	for _, function := range functionlist {
		function(app)
	}
}
