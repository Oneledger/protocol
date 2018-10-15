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
		ProcessField: PrintNode,
		Name:         "base",
	}
	Iterate(base, action)
}

// PrintNode is called for each node
func PrintNode(action *Action, input interface{}) interface{} {
	log.Dump("Node", input)
	return true
}

// Make a deep copy of the data
func Clone(base interface{}) interface{} {
	action := &Action{
		ProcessField: CloneNode,
		Name:         "base",
	}
	result := Iterate(base, action)

	for _, value := range action.Processed["base"].Children {
		return value
	}
	return result
}

// Clone each element
func CloneNode(action *Action, input interface{}) interface{} {

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

// Clone and add in SerialWrapper
func Extend(base interface{}) interface{} {

	// Don't need to recurse
	if IsPrimitive(base) {
		typeof := reflect.TypeOf(base).Name()
		dict := make(map[string]interface{})
		dict[""] = base
		wrapper := SerialWrapper{Type: typeof, Size: 1, Fields: dict}
		return wrapper
	}

	action := &Action{
		ProcessField: ExtendNode,
		Name:         "base",
	}

	result := Iterate(base, action)

	for _, value := range action.Processed["base"].Children {
		return value
	}

	return result
}

// Extend the input by replacing all structures with a wrapper
func ExtendNode(action *Action, input interface{}) interface{} {

	parent := action.Path.StringPeekN(1)

	if input == nil || IsNilValue(input) {
		return input
	}

	if GetBaseTypeString(input) == "big.Int" {
		return input
	}

	if IsContainer(input) {
		mapping, size := ConvertMap(input)
		log.Dump("Have Mapping", mapping)

		// Attach all of the interface children
		for key, value := range action.Processed[action.Name].Children {
			mapping[key] = value
			delete(action.Processed[action.Name].Children, key)
		}
		log.Dump("Revised Mapping", mapping)

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
			Size:   size,
		}
		action.Processed[parent].Children[action.Name] = wrapper

		return wrapper
	}
	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {

	if IsSerialWrapper(base) {
		wrapper := base.(SerialWrapper)
		typeEntry := GetTypeEntry(wrapper.Type, wrapper.Size)
		if typeEntry.Category == PRIMITIVE {
			return wrapper.Fields[""]
		}
	}

	action := &Action{
		ProcessField: ContractNode,
		Name:         "base",
	}

	result := Iterate(base, action)

	for _, value := range action.Processed["base"].Children {
		return value
	}

	return result
}

// Replace any incoming SerialWrappers with the correct structure
func ContractNode(action *Action, input interface{}) interface{} {
	grandparent := action.Path.StringPeekN(2)
	if grandparent == "" {
		// Top-level, just use the parent
		grandparent = action.Path.StringPeekN(1)
	}

	log.Dump("ContractNode", input)
	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		stype := wrapper.Type
		size := wrapper.Size

		result := Alloc(stype, size)
		/*
			if stype == "ed25519.PubKeyEd25519" {
				log.Dump("New ed25519.PubKeyEd25519", result)
			}
		*/

		// Needs to come from the serialized name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		} else {
			action.IsPointer = false
		}

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		/*
			if stype == "ed25519.PubKeyEd25519" {
				log.Dump("Set ed25519.PubKeyEd25519", result)
			}
		*/

		SetProcessed(action, grandparent, action.Name, result)

		return CleanValue(action, result)
	}

	if IsSerialWrapperMap(input) {
		wrapper := input.(map[string]interface{})
		stype := wrapper["Type"].(string)
		sizeFloat := wrapper["Size"].(float64)
		size := int(sizeFloat)

		result := Alloc(stype, size)
		/*
			if stype == "id.AccountKey" {
				log.Dump("New AccountKey", result)
			}
		*/

		// Needs to come from the serialized name
		if strings.HasPrefix(stype, "*") {
			action.IsPointer = true
		} else {
			action.IsPointer = false
		}

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		/*
			if stype == "id.AccountKey" {
				log.Dump("Set AccountKey", result)
			}
		*/

		SetProcessed(action, grandparent, action.Name, result)
		return CleanValue(action, result)
	}

	SetProcessed(action, grandparent, action.Name, input)
	return CleanValue(action, input)
}

func CleanValue(action *Action, input interface{}) interface{} {
	if input == nil {
		return nil
	}

	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return input
	}

	if action.IsPointer {
		return input
	}

	element := reflect.ValueOf(input).Elem().Interface()
	return element
}

// Set the as a processed result, and handle pointers nicely.
func SetProcessed(action *Action, parent string, name string, input interface{}) {
	if input == nil {
		action.Processed[parent].Children[name] = nil
	} else {
		action.Processed[parent].Children[name] = CleanValue(action, input)
	}
}
