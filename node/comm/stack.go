/*
	Copyright 2017-2018 OneLedger
*/

package comm

import "github.com/Oneledger/protocol/node/log"

type Stack struct {
	Base []interface{}
}

func NewStack() *Stack {
	return &Stack{make([]interface{}, 0)}
}

func (stack *Stack) Push(element interface{}) {
	base := []interface{}{element}
	stack.Base = append(base, stack.Base...)
}

func (stack Stack) Len() int {
	return len(stack.Base)
}

func (stack Stack) Peek() interface{} {
	return stack.Base[0]
}

func (stack Stack) StringPeekN(index int) string {
	if index >= stack.Len() {
		return ""
	}
	return stack.PeekN(index).(string)
}

func (stack Stack) PeekN(index int) interface{} {
	return stack.Base[index]
}

func (stack *Stack) Shift() {
	stack.Base = stack.Base[1:]
}

func (stack *Stack) Pop() interface{} {
	element := stack.Peek()
	stack.Shift()
	return element
}

func (stack Stack) Print() {
	for i := 0; i < len(stack.Base); i++ {
		log.Debug("Stack", "index", i, "value", stack.Base[i])
	}
}
