/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

// Print is a post order traversal (not that useful, should be pre-order)
func Print(base interface{}) {
	action := &Action{
		ProcessField: PrintIt,
		Name:         "base",
	}
	Iterate(base, action)
}

// PrintIt is called for each node
func PrintIt(action *Action, input interface{}) interface{} {
	log.Dump("#### Node", input, action)
	return true
}

// Make a deep copy of the data
func Clone(base interface{}) interface{} {
	action := &Action{
		ProcessField: CloneIt,
		Name:         "base",
	}
	result := Iterate(base, action)

	for _, value := range action.Processed["*base"].Children {
		return value
	}

	return result
}

// Clone each element
func CloneIt(action *Action, input interface{}) interface{} {

	parent := action.Path.StringPeekN(1)

	if IsContainer(input) {
		copy := reflect.New(reflect.TypeOf(input)).Interface()

		// Overwrite with children
		for key, value := range action.Processed[action.Name].Children {
			Set(copy, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		SetProcessed(action, parent, action.Name, copy)
		return copy
	}

	copy := input
	SetProcessed(action, parent, action.Name, copy)

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
	//log.Dump("Extend this", base)

	action := &Action{
		ProcessField: ExtendIt,
		Name:         "base",
	}

	result := Iterate(base, action)
	return result
}

// Extend the input by replacing all structures with a wrapper
func ExtendIt(action *Action, input interface{}) interface{} {
	//log.Debug("ExtendIt", "input", input)

	parent := action.Path.StringPeekN(1)

	if IsStructure(input) {
		mapping := ConvertMap(input)

		// Attach all of the interface children
		for key, value := range action.Processed[action.Name].Children {
			mapping[key] = value
			delete(action.Processed[action.Name].Children, key)
		}

		typestr := reflect.TypeOf(input).String()

		if typestr == "reflect.Value" {
			log.Fatal("Have a reflect.Value, bad call")
		}

		if action.IsPointer {
			typestr = "*" + typestr
		}

		wrapper := SerialWrapper{
			Type:   typestr,
			Fields: mapping,
		}
		action.Processed[parent].Children[action.Name] = wrapper

		return wrapper

	} else if IsContainer(input) {
		return input

	}

	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{
		ProcessField: ContractIt,
		Name:         "base",
	}

	result := Iterate(base, action)

	for _, value := range action.Processed["*base"].Children {
		return value
	}
	return result
}

// Replace any incoming SerialWrappers with the correct structure
func ContractIt(action *Action, input interface{}) interface{} {
	//log.Debug("ContractIt", "input", input)

	grandparent := action.Path.StringPeekN(2)
	if grandparent == "" {
		grandparent = action.Name
	}

	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		stype := wrapper.Type
		result := NewStruct(stype)

		// Needs to come from the serialized name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		}

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		SetProcessed(action, grandparent, action.Name, result)
		return result
	}

	if IsSerialWrapperMap(input) {
		wrapper := input.(map[string]interface{})
		stype := wrapper["Type"].(string)

		result := NewStruct(stype)

		// Needs to come from the serialized name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		}

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		SetProcessed(action, grandparent, action.Name, result)
		return result
	}

	SetProcessed(action, grandparent, action.Name, input)
	return input
}

// Set the as a processed result, and handle pointers nicely.
func SetProcessed(action *Action, parent string, name string, input interface{}) {
	if input == nil {

	} else if reflect.TypeOf(input).Kind() != reflect.Ptr {
		action.Processed[parent].Children[name] = input

	} else if action.IsPointer {
		action.Processed[parent].Children[name] = input

	} else {
		element := reflect.ValueOf(input).Elem().Interface()
		action.Processed[parent].Children[name] = element
	}
}
