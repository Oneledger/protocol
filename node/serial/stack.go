/*
	Copyright 2017-2018 OneLedger

	Basic stack implementation

	TODO: Currently reversed, which is probably not good for performance.
*/
package serial

import "github.com/Oneledger/protocol/node/log"

type Stack struct {
	base []interface{}
}

func NewStack() *Stack {
	return &Stack{make([]interface{}, 0)}
}

func (stack *Stack) Push(element interface{}) {
	base := []interface{}{element}
	stack.base = append(base, stack.base...)
}

func (stack Stack) Len() int {
	return len(stack.base)
}

func (stack Stack) Peek() interface{} {
	return stack.base[0]
}

func (stack Stack) StringPeekN(index int) string {
	if index >= stack.Len() {
		return ""
	}
	return stack.PeekN(index).(string)
}

func (stack Stack) PeekN(index int) interface{} {
	return stack.base[index]
}

func (stack *Stack) Shift() {
	stack.base = stack.base[1:]
}

func (stack *Stack) Pop() interface{} {
	element := stack.Peek()
	stack.Shift()
	return element
}

func (stack Stack) Print() {
	for i := 0; i < len(stack.base); i++ {
		log.Debug("Stack", "index", i, "value", stack.base[i])
	}
}
