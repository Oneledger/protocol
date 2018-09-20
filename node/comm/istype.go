/*
	Copyright 2017-2018 OneLedger
*/

package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// Language primitives
func IsPrimitive(input interface{}) bool {
	switch input.(type) {
	// Booleans
	case bool:
		return true

	// Integers of different sizes
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true

	// Unsigned integers
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true

	// Floating point values
	case float32:
		return true
	case float64:
		return true

	// Complex numbers
	case complex64:
		return true
	case complex128:
		return true

	// Strings
	case string:
		return true
	}

	return false
}

func IsInterface(input interface{}) bool {
	kind := reflect.TypeOf(input).Kind()
	log.Debug("Kind", "kind", kind)
	if kind == reflect.Interface {
		return true
	}
	return false
}

func IsStructure(input interface{}) bool {
	kind := reflect.TypeOf(input).Kind()
	if kind == reflect.Struct {
		return true
	}
	return false
}

func IsStructurePtr(input interface{}) bool {
	kind := reflect.TypeOf(input).Kind()
	if kind == reflect.Ptr {
		return true
	}
	return false
}

// Difficult data types
func IsDifficult(input interface{}) bool {
	if IsContainer(input) || IsSpecial(input) {
		return true
	}
	return false
}

// Container data types
func IsContainer(input interface{}) bool {
	kind := reflect.TypeOf(input).Kind()

	switch kind {
	case reflect.Array:
		return true
	case reflect.Slice:
		return true
	case reflect.Map:
		return true
	}
	return false
}

// Special datatypes
func IsSpecial(input interface{}) bool {
	kind := reflect.TypeOf(input).Kind()

	switch kind {
	case reflect.Chan:
		return true
	case reflect.Func:
		return true
	case reflect.UnsafePointer:
		return true
	}
	return false
}
