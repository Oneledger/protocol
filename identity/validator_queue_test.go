package identity

import (
	"os"
	"testing"

	"github.com/Oneledger/protocol/utils"
	"github.com/stretchr/testify/assert"
)

var vq1 = utils.NewQueued([]byte("a"), 200, 1)
var vq2 = utils.NewQueued([]byte("b"), 1000, 2)
var vq3 = utils.NewQueued([]byte("c"), 300, 3)
var vq4 = utils.NewQueued([]byte("d"), 2000, 4)

var Vq *ValidatorQueue

func setup() {
	logger.Debug("This is main setting up")
}

func shutdown() {
	logger.Debug("This is main tearing down")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func subtestsetup(funName string) {
	logger.Debug("This is subtest setting up for", funName)
	Vq = &ValidatorQueue{make(utils.PriorityQueue, 0)}
}

func subtestteardown(funName string) {
	logger.Debug("This is subtest tearing down for", funName)
	Vq = nil
}

func TestValidatorQueue_Len(t *testing.T) {
	subtestsetup(t.Name())
	defer subtestteardown(t.Name())
	t.Run("run validator queue len fun", func(t *testing.T) {
		assert.Equal(t, 0, Vq.Len())
	})
	t.Run("run validator queue len fun with actual item in it", func(t *testing.T) {
		Vq.Push(vq1)
		Vq.Push(vq2)
		Vq.Push(vq3)
		assert.Equal(t, 3, Vq.Len())
	})
}

func TestValidatorQueue_Push(t *testing.T) {
	subtestsetup(t.Name())
	defer subtestteardown(t.Name())
	t.Run("Before push", func(t *testing.T) {
		assert.Equal(t, 0, Vq.Len())
	})
	t.Run("After push", func(t *testing.T) {
		Vq.Push(vq1)
		Vq.Push(vq2)
		Vq.Push(vq3)
		assert.Equal(t, 3, Vq.Len())
	})
}

func TestValidatorQueue_Pop(t *testing.T) {
	subtestsetup(t.Name())
	defer subtestteardown(t.Name())
	t.Run("pop empty validator queue", func(t *testing.T) {
		result := Vq.Pop()
		assert.Empty(t, result)
	})
	t.Run("check if pop order is correct with push, pop, append, update fun", func(t *testing.T) {
		Vq.append(vq1)
		Vq.append(vq2)
		Vq.append(vq3)
		Vq.Init()
		result := Vq.Pop()
		assert.Equal(t, int64(1000), result.Priority())
		assert.Equal(t, []byte("b"), result.Value())
		Vq.Push(vq4)
		newResult := Vq.Pop()
		assert.Equal(t, int64(2000), newResult.Priority())
		assert.Equal(t, []byte("d"), newResult.Value())
		newResult2 := Vq.Pop()
		assert.Equal(t, int64(300), newResult2.Priority())
		assert.Equal(t, []byte("c"), newResult2.Value())
		Vq.update(vq1, []byte("o"), 500)
		newResult3 := Vq.Pop()
		assert.Equal(t, int64(500), newResult3.Priority())
		assert.Equal(t, []byte("o"), newResult3.Value())
		newResult4 := Vq.Pop()
		assert.Empty(t, newResult4)
	})
}
