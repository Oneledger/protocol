package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cRouter Router
	app     *App
)

func internalTX1(app *App) {
	fmt.Println("INTERNALTX1 ", app.name)
}
func internalTX2(app *App) {
	fmt.Println("INTERNALTX2 ", app.name)
}

func init() {
	app = &App{
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
