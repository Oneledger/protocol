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

	//if reflect.TypeOf(action.Value).Kind() == reflect.Ptr {
	if IsContainer(input) {

		copy := reflect.New(reflect.TypeOf(input)).Interface()
		log.Debug("Copied", "copy", copy, "input", input)

		// Overwrite with children
		for key, value := range action.Processed[action.Name].Children {
			log.Debug("Child", key, value)
			Set(copy, key, value)
		}

		action.Processed[action.ParentName].Children[action.Name] = copy
		return copy

	} else if IsPointer(input) {
		return input
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
		}

		wrapper := SerialWrapper{
			Type:   reflect.TypeOf(input).String(),
			Fields: mapping,
		}

		action.Processed[action.ParentName].Children[action.Name] = wrapper
		return wrapper
	}
	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	log.Debug("########## Contract")
	action := &Action{ProcessField: ContractIt}
	result := Iterate(base, action)
	return result
}

func ContractIt(action *Action, input interface{}) interface{} {
	log.Debug("ContractIt", "action", action, "input", input)

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

		//action.Processed[action.ParentName].Children[action.Name] = result
		return result
	}

	action.Processed[action.ParentName].Children[action.Name] = input
	return input
}
