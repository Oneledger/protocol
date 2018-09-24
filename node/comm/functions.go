/*
	Copyright 2017-2018 OneLedger
*/
package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

func Print(base interface{}) {
	action := &Action{ProcessField: PrintIt}
	Iterate(base, action)
}

func PrintIt(action *Action, input interface{}) interface{} {
	log.Debug("PrintIt", "input", input, "action", action)
	return true
}

// Make a deep copy of the data
func Clone(base interface{}) interface{} {
	action := &Action{ProcessField: CloneIt}
	result := Iterate(base, action)
	return result
}

func CloneIt(action *Action, input interface{}) interface{} {
	if IsContainer(input) {

		copy := reflect.New(reflect.TypeOf(input)).Interface()

		// Overwrite with children
		for key, value := range action.Processed[action.Name].Children {
			Set(copy, key, value)
		}

		action.Processed[action.ParentName].Children[action.Name] = copy
		return copy

	} else if IsPointer(input) {
		copy := action.Processed[action.Name].Children[action.Name]
		element := copy
		action.Processed[action.ParentName].Children[action.Name] = element
		return element
	}

	copy := action.Value
	action.Processed[action.ParentName].Children[action.Name] = copy
	return copy
}

// ConvertMap takes a structure and return a map of its elements
func ConvertMap(structure interface{}) map[string]interface{} {
	var result map[string]interface{}

	children := GetChildren(structure)
	result = make(map[string]interface{}, len(children))

	for _, child := range children {
		result[child.Name] = child.Value
	}
	return result
}

// Clone and add in SerialWrapper
func Extend(base interface{}) interface{} {
	log.Debug("Extend")

	action := &Action{ProcessField: ExtendIt}
	result := Iterate(base, action)
	return result
}

func ExtendIt(action *Action, input interface{}) interface{} {
	if IsContainer(action.Value) {
		if !IsStructure(action.Value) {
			log.Fatal("Can't handle other containers yet", "value", action.Value)
		}

		mapping := ConvertMap(action.Value)

		// Attach all of the interface children
		for key, value := range action.Processed[action.Name].Children {
			mapping[key] = value
			delete(action.Processed[action.Name].Children, key)
		}

		pre := ""
		if IsPointer(action.Value) {
			pre = "*"
		}
		wrapper := SerialWrapper{
			Type:   pre + reflect.TypeOf(input).String(),
			Fields: mapping,
		}

		action.Processed[action.ParentName].Children[action.Name] = wrapper
		return wrapper
	}
	for _, value := range action.Processed[action.ParentName].Children {
		return value
	}
	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{ProcessField: ContractIt}
	result := Iterate(base, action)
	return result
}

func ContractIt(action *Action, input interface{}) interface{} {
	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		result := NewStruct(wrapper.Type)

		// Fill it with the deserialized values
		for key, value := range wrapper.Fields {
			Set(result, key, value)
		}

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			Set(result, key, value)
		}
		return result
	}

	action.Processed[action.ParentName].Children[action.Name] = input
	return input
}
