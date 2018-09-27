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
	action := &Action{
		ProcessField: CloneIt,
		Name:         "base",
	}
	result := Iterate(base, action)
	return result
}

// Clone each element
func CloneIt(action *Action, input interface{}) interface{} {

	parent := action.Path.StringPeekN(1)

	if IsPointer(input) {
		copy := action.Processed[action.Name].Children[action.Name]
		element := copy
		action.Processed[parent].Children[action.Name] = element
		return element
	}

	if IsContainer(input) {
		copy := reflect.New(reflect.TypeOf(input)).Interface()

		// Overwrite with children
		for key, value := range action.Processed[action.Name].Children {
			Set(copy, key, value)
		}

		action.Processed[parent].Children[action.Name] = copy
		return copy
	}

	copy := input
	action.Processed[parent].Children[action.Name] = copy

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
		Name:         "base",
	}

	result := Iterate(base, action)
	return result
}

func ExtendIt(action *Action, input interface{}) interface{} {
	log.Debug("ExtendIt", "action", action, "input", input)

	parent := action.Path.StringPeekN(1)

	if IsContainer(input) {
		if !IsStructure(input) {
			log.Fatal("Can't handle other containers yet", "input", input)
		}

		mapping := ConvertMap(input)
		//parent = strings.TrimPrefix(parent, "*")

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

		wrapper := SerialWrapper{
			Type:   typestr,
			Fields: mapping,
		}

		log.Debug("Assigning to", "parent", parent, "name", action.Name, "Processed", action.Processed)
		action.Processed[parent].Children[action.Name] = wrapper

		return wrapper
	}

	// In general return a child
	/*
		for _, value := range action.Processed[parent].Children {
			log.Debug("Pushing up a child", "parent", parent)
			return value
		}
	*/
	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{
		VisitPrimitives: true,
		ProcessField:    ContractIt,
		Name:            "base",
	}

	result := Iterate(base, action)
	for _, value := range action.Processed["*base"].Children {
		return value
	}
	return result
}

// Replace any incoming SerialWrappers with the correct structure
func ContractIt(action *Action, input interface{}) interface{} {
	log.Debug("ContractIt", "input", input, "action", action)

	grandparent := action.Path.StringPeekN(2)
	if grandparent == "" {
		grandparent = action.Name
	}

	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		stype := wrapper.Type
		result := NewStruct(stype)

		// Needs to come from the name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		}

		// Fill it with the deserialized values
		/*
			for key, value := range wrapper.Fields {
				log.Debug("Setting Field", key, value)
				Set(result, key, value)
			}
		*/

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			log.Debug("Overwriting Modified Children", key, value,
				"grandparent", grandparent, "name", action.Name, "Processed", action.Processed)
			Set(result, key, value)
			//delete(action.Processed["Fields"].Children, key)
		}

		if action.IsPointer {
			action.Processed[grandparent].Children[action.Name] = result
		} else {
			element := reflect.ValueOf(result).Elem().Interface()
			action.Processed[grandparent].Children[action.Name] = element
		}
		return result
	}

	if IsSerialWrapperMap(input) {
		wrapper := input.(map[string]interface{})
		stype := wrapper["Type"].(string)
		result := NewStruct(stype)

		// Needs to come from the name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		}

		// Fill it with the deserialized values
		/*
			for key, value := range wrapper["Fields"].(map[string]interface{}) {
				log.Debug("Setting Field", key, value)
				Set(result, key, value)
			}
		*/

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			log.Debug("Map Overwriting Modified Children", key, value,
				"grandparent", grandparent, "name", action.Name)

			Set(result, key, value)
			//delete(action.Processed["Fields"].Children, key)
		}

		log.Debug("Map Pushing up", "isptr", action.IsPointer, "name", action.Name, "grandparent", grandparent, "result", result)
		if action.IsPointer {
			action.Processed[grandparent].Children[action.Name] = result
		} else {
			// Get the underlying element
			element := reflect.ValueOf(result).Elem().Interface()
			action.Processed[grandparent].Children[action.Name] = element
		}

		return result
	}

	log.Debug("Final Pushing up", "grandparent", grandparent, "child", action.Name, "input", input)
	action.Processed[grandparent].Children[action.Name] = input
	return input
}
