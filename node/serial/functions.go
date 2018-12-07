/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"fmt"
	"math/big"
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
	Iterate(base, action, 1)
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
	result := Iterate(base, action, 1)

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

func IsBigint(base interface{}) bool {
	typeOf := GetBaseType(base)

	if typeOf.String() == "big.Int" {
		return true
	}
	return false
}

// Special case to handle big ints, forces it to always be a pointer
func ExtendBigint(base interface{}) interface{} {
	typeOf := "*big.Int"

	dict := make(map[string]interface{})
	convert := base.(big.Int)
	dict["string"] = convert.String()
	wrapper := SerialWrapper{Type: typeOf, Size: 1, Fields: dict}

	return wrapper
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

	if IsBigint(base) {
		return ExtendBigint(base)
	}

	action := &Action{
		ProcessField: ExtendNode,
		Name:         "base",
	}

	result := Iterate(base, action, 1)

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

	// Special case to handle big.Int's inability to marshal properly
	if IsBigint(input) {
		wrapper := ExtendBigint(input)
		action.Processed[parent].Children[action.Name] = wrapper
		return wrapper
	}

	// Ignore this variable because of its type.
	if IgnoreVariable(input) {
		return input
	}

	if IsContainer(input) {
		mapping, size := ConvertMap(input)

		// Override all of the underlying container items.
		for key, value := range action.Processed[action.Name].Children {
			// key has the depth included, needs to be removed.
			mapping[GetFieldName(key)] = value
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
			Size:   size,
		}
		action.Processed[parent].Children[action.Name] = wrapper

		return wrapper
	}
	return input
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {

	var typeEntry TypeEntry

	underlying := GetBaseValue(base).Interface()
	if IsSerialWrapper(underlying) || IsSerialWrapperMap(underlying) {
		wrapper := underlying.(SerialWrapper)
		typeEntry = GetTypeEntry(wrapper.Type, wrapper.Size)
		if typeEntry.Category == PRIMITIVE {
			value := wrapper.Fields[""]
			return ConvertValue(value, typeEntry.DataType).Interface()
		}
	}

	action := &Action{
		ProcessField: ContractNode,
		Name:         "base",
	}

	result := Iterate(base, action, 1)

	for _, value := range action.Processed["base"].Children {
		result = value
		break
		return value
	}

	if typeEntry.Category == SLICE || typeEntry.Category == ARRAY {
		interim := reflect.New(typeEntry.RootType)
		interim.Elem().Set(reflect.ValueOf(result))
		result = interim.Interface()
	}

	return CleanValue(action, result)
}

func ContractBigint(value string, size int) interface{} {
	number := new(big.Int)
	_, err := fmt.Sscan(value, number)
	if err != nil {
		log.Fatal("Invalid integer string", "err", err, "value", value)
	}
	return &number
}

// Replace any incoming SerialWrappers with the correct structure
func ContractNode(action *Action, input interface{}) interface{} {
	grandparent := action.Path.StringPeekN(2)
	if grandparent == "" {
		// Top-level, just use the parent
		grandparent = action.Path.StringPeekN(1)
	}

	if IsSerialWrapper(input) {
		wrapper := input.(SerialWrapper)
		stype := wrapper.Type
		size := wrapper.Size

		if stype == "*big.Int" {
			result := ContractBigint(wrapper.Fields["string"].(string), size)
			action.IsPointer = true // TODO: Should be handled correctly
			SetProcessed(action, grandparent, action.Name, result)
			return CleanValue(action, result)
		}

		result := Alloc(stype, size)

		// Needs to come from the serialized name, not the internal variable
		action.IsPointer = strings.HasPrefix(stype, "*")

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			//log.Dump(action.Name+" Child is "+key, value)
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		SetProcessed(action, grandparent, action.Name, result)
		return CleanValue(action, result)
	}

	if IsSerialWrapperMap(input) {
		wrapper := input.(map[string]interface{})
		stype := wrapper["Type"].(string)
		sizeFloat := wrapper["Size"].(float64)
		size := int(sizeFloat)

		if stype == "*big.Int" {
			fields := wrapper["Fields"].(map[string]interface{})
			result := ContractBigint(fields["string"].(string), size)
			action.IsPointer = true // TODO: Should be handled correctly
			SetProcessed(action, grandparent, action.Name, result)
			return CleanValue(action, result)
		}

		result := Alloc(stype, size)

		// Needs to come from the serialized name
		action.IsPointer = strings.HasPrefix(stype, "*")

		// Overwrite with any better children
		for key, value := range action.Processed[action.Name].Children {
			//log.Dump(action.Name+" Child is "+key, value)
			Set(result, key, value)
			delete(action.Processed[action.Name].Children, key)
		}

		SetProcessed(action, grandparent, action.Name, result)
		return CleanValue(action, result)
	}

	SetProcessed(action, grandparent, action.Name, input)
	return CleanValue(action, input)
}

// Convert from pointers back to their values, if necessary
func CleanValue(action *Action, input interface{}) interface{} {
	if input == nil {
		return nil
	}

	// Return it directly if is already a value
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		return input
	}

	// Return it as a pointer, that's what we want
	if action.IsPointer {
		return input
	}

	// Remove the pointer
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
