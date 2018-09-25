/*
	Copyright 2017-2018 OneLedger
*/
package comm

import (
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

func Print(base interface{}) {
	action := &Action{
		ProcessField: PrintIt,
		ParentName:   "root",
		Name:         "base",
	}
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

// Clone each element
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

	copy := input
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

	action := &Action{
		ProcessField: ExtendIt,
		ParentName:   "root",
		Name:         "base",
	}

	result := Iterate(base, action)
	return result
}

func ExtendIt(action *Action, input interface{}) interface{} {
	log.Debug("ExtendIt", "action", action, "input", input)

	if IsContainer(input) {
		if !IsStructure(input) {
			log.Fatal("Can't handle other containers yet", "input", input)
		}

		mapping := ConvertMap(input)
		parent := strings.TrimPrefix(action.ParentName, "*")

		// Attach all of the interface children
		for key, value := range action.Processed[action.Name].Children {
			log.Debug("Fixing Children", "name", action.Name, "key", key, "processed", action.Processed)
			mapping[key] = value
			delete(action.Processed[action.Name].Children, key)
		}

		typestr := reflect.TypeOf(input).String()

		if action.IsPointer {
			typestr = "*" + typestr
		}

		log.Debug("Wrapping", "Name", action.Name, "typestr", typestr, "Parent", parent)
		wrapper := SerialWrapper{
			Type:   typestr,
			Fields: mapping,
		}

		action.Processed[parent].Children[action.Name] = wrapper

		return wrapper
	}

	// In general return a child
	for _, value := range action.Processed[action.ParentName].Children {
		log.Debug("Pushing up a child", "parent", action.ParentName)
		return value
	}

	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{
		VisitPrimitives: true,
		ProcessField:    ContractIt,
		ParentName:      "root",
		Name:            "base",
	}
	result := Iterate(base, action)
	return result
}

// Replace any incoming SerialWrappers with the correct structure
func ContractIt(action *Action, input interface{}) interface{} {
	log.Debug("ContractIt", "input", input, "action", action)

	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		result := NewStruct(wrapper.Type)

		// Fill it with the deserialized values
		/*
			for key, value := range wrapper.Fields {
				log.Debug("Setting Field", key, value)
				Set(result, key, value)
			}
		*/

		// Overwrite with any better children
		for key, value := range action.Processed["Fields"].Children {
			log.Debug("Overwriting Modified Children", key, value)
			Set(result, key, value)
		}

		log.Debug("Pushing up", "parent", action.ParentName, "child", action.Name, "result", result)
		action.Processed[action.ParentName].Children[action.Name] = result
		return result
	}

	if IsSerialWrapperMap(input) {
		wrapper := input.(map[string]interface{})
		stype := wrapper["Type"].(string)
		result := NewStruct(stype)

		// Fill it with the deserialized values
		/*
			for key, value := range wrapper["Fields"].(map[string]interface{}) {
				log.Debug("Setting Field", key, value)
				Set(result, key, value)
			}
		*/

		// Overwrite with any better children
		for key, value := range action.Processed["Fields"].Children {
			log.Debug("Overwriting Modified Children", key, value)
			Set(result, key, value)
		}

		log.Debug("Map Pushing up", "parent", action.ParentName, "child", action.Name, "result", result)
		if strings.HasPrefix(stype, "*") {
			action.Processed[action.ParentName].Children[action.Name] = result
		} else {
			// Get the underlying element
			element := reflect.ValueOf(result).Elem().Interface()
			action.Processed[action.ParentName].Children[action.Name] = element
		}
		return result
	}

	action.Processed[action.ParentName].Children[action.Name] = input
	return input
}
