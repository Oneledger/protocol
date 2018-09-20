/*
	Copyright 2017-2018 OneLedger
*/

package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

type Test struct {
	Name string
}

// GetValue returns the underlying value, even if it is a pointer
func GetValue(base interface{}) reflect.Value {
	element := reflect.ValueOf(base)
	if element.Kind() == reflect.Ptr {
		return element.Elem()
	}
	return element
}

// SetStructure takes a pointer to a structure and sets it.
func SetStruct(parent interface{}, fieldNum int, child interface{}) bool {
	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Struct {
		log.Warn("Not a structure", "element", element)
		return false
	}

	field := element.Field(fieldNum)
	if !field.IsValid() {
		log.Warn("Field is invalid", "field", field)
		return false
	}

	if !field.CanSet() {
		log.Warn("Not Settable", "field", field)
		return false
	}

	if field.Kind() != reflect.Interface {
		log.Warn("Field is not an interface", "kind", field.Kind(), "field", field)
		return false
	}

	newValue := GetValue(child)

	field.Set(newValue)
	return true
}

func Print(base interface{}) {
	action := &Action{ProcessField: PrintIt}

	Iterate(base, action)
}

func PrintIt(action *Action, input interface{}) bool {
	log.Debug("PrintIt", "action", action, "value", input)
	return true
}

// Clone and add in SerialWrapper
func Extend(base interface{}) interface{} {
	action := &Action{ProcessField: ExtendIt}

	Iterate(base, action)

	var last interface{}
	for _, value := range action.Children {
		last = value
	}
	return last
}
func ExtendIt(action *Action, input interface{}) bool {
	if reflect.TypeOf(action.Value).Kind() == reflect.Ptr {
		var copy interface{}

		if reflect.TypeOf(action.Value).Kind() == reflect.Ptr {
			copy = action.Value
		} else {
			copy = action.Value // Implicit Copy
		}

		for key, value := range action.Children {
			SetStruct(copy, key, value)
			delete(action.Children, key)
		}
		wrapper := SerialWrapper{
			Type:      reflect.TypeOf(copy).Name(),
			IsPointer: false,
			Value:     copy,
		}
		action.Children[action.Field] = wrapper
	}
	return true
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{ProcessField: ContractIt}

	Iterate(base, action)

	var last interface{}
	for key, value := range action.Children {
		log.Debug("Have a Final Child", "key", key)
		last = value
	}
	return last
}
func ContractIt(action *Action, input interface{}) bool {
	log.Debug("ContractIt", "struct", action.Struct, "value", input)
	return true
}

// Make a deep copy of the data
func Clone(base interface{}) interface{} {
	action := &Action{ProcessField: CloneIt}

	Iterate(base, action)

	var last interface{}
	for _, value := range action.Children {
		last = value
	}
	return last
}

func CloneIt(action *Action, input interface{}) bool {

	if reflect.TypeOf(action.Value).Kind() == reflect.Ptr {
		var copy interface{}

		if reflect.TypeOf(action.Value).Kind() == reflect.Ptr {
			copy = action.Value
		} else {
			copy = action.Value // Implicit Copy
		}

		for key, value := range action.Children {
			SetStruct(copy, key, value)
			delete(action.Children, key)
		}
		action.Children[action.Field] = copy
	}
	return true
}
