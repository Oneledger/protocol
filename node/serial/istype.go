/*
	Copyright 2017-2018 OneLedger
*/

package serial

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// Language primitives
func IsPrimitive(input interface{}) bool {
	if input == nil {
		return false
	}

	switch input.(type) {

	// Booleans
	case bool:
		return true

	// Integers of different sizes
	case int, int8, int16, int32, int64:
		return true

	// Unsigned integers, includes byte
	case uint, byte, uint16, uint32, uint64:
		return true

	// Floating point values
	case float32, float64:
		return true

	// Complex numbers
	case complex64, complex128:
		return true

	// Strings
	case string:
		return true
	}

	return false
}

func UnderlyingType(input interface{}) reflect.Type {
	//log.Dump("Calling underlyingType", input)

	var proto interface{}
	if input == nil {
		return reflect.TypeOf(proto)
	}

	typeOf := reflect.TypeOf(input)
	if typeOf == nil {
		log.Fatal("Invalid type")
	}

	kind := typeOf.Kind()

	entry := GetTypeEntry(kind.String(), 1)
	if entry.Category != UNKNOWN {
		return entry.DataType
	}

	switch kind {
	case reflect.Ptr:
		return UnderlyingType(reflect.ValueOf(input).Elem().Interface())

	case reflect.Struct:
		return typeOf

	case reflect.Map:
		keyType := typeOf.Key()
		keyValue := reflect.New(keyType).Interface()
		elementType := typeOf.Elem()
		elementValue := reflect.New(elementType).Interface()
		return reflect.MapOf(UnderlyingType(keyValue), UnderlyingType(elementValue))

	case reflect.Slice:
		elementType := typeOf.Elem()
		elementValue := reflect.New(elementType).Interface()
		return reflect.SliceOf(UnderlyingType(elementValue))

	case reflect.Array:
		elementType := typeOf.Elem()
		elementValue := reflect.New(elementType).Interface()
		return reflect.ArrayOf(32, UnderlyingType(elementValue))
	}

	log.Warn("Not sure, so interface{} type")
	return reflect.TypeOf(proto)
}

func IsPrimitiveArray(input interface{}) bool {
	entry := GetTypeEntry(GetBaseTypeString(input), 1)
	if entry.Category == ARRAY {
		if entry.ValueType.Category == PRIMITIVE {
			return true
		}
	}

	underlying := UnderlyingType(input)
	if underlying != nil {
		entry = GetTypeEntry(underlying.String(), 1)
		if entry.Category == ARRAY {
			if entry.ValueType.Category == PRIMITIVE {
				return true
			}
		}
	}
	return false
}

func IsPrimitiveSlice(input interface{}) bool {
	entry := GetTypeEntry(GetBaseTypeString(input), 1)
	if entry.Category == SLICE {
		if entry.ValueType.Category == PRIMITIVE {
			return true
		}
	}

	underlying := UnderlyingType(input)
	if underlying != nil {
		entry = GetTypeEntry(underlying.String(), 1)
		if entry.Category == SLICE {
			if entry.ValueType.Category == PRIMITIVE {
				return true
			}
		}
	}
	return false
}

func IsInterface(input interface{}) bool {
	if input == nil {
		return false
	}

	kind := reflect.TypeOf(input).Kind()
	if kind == reflect.Interface {
		return true
	}
	return false
}

// See if the underlying Value is nil
func IsNilValue(input interface{}) bool {
	if input == nil {
		return true
	}

	valueOf := reflect.ValueOf(input)
	kind := valueOf.Kind()

	switch kind {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		if valueOf.IsNil() {
			return true
		}
	}
	return false
}

func IsPointer(input interface{}) bool {
	if input == nil {
		return false
	}

	kind := reflect.TypeOf(input).Kind()
	if kind == reflect.Ptr {
		return true
	}
	return false
}

func IsStructure(input interface{}) bool {
	if input == nil {
		return false
	}

	kind := reflect.TypeOf(input).Kind()
	if kind == reflect.Struct {
		return true
	}
	return false
}

// An array or slice that is based around primtives, like []byte
func IsPrimitiveContainer(input interface{}) bool {
	if !IsContainer(input) {
		return false
	}

	if IsPrimitiveArray(input) {
		return true
	}

	if IsPrimitiveSlice(input) {
		return true
	}
	return false
}

// Container data types
func IsContainer(input interface{}) bool {
	if input == nil {
		return false
	}

	kind := reflect.TypeOf(input).Kind()

	switch kind {
	case reflect.Struct:
		return true

	case reflect.Array:
		return true

	case reflect.Slice:
		return true

	case reflect.Map:
		return true
	}
	return false
}

// Difficult data types
func IsDifficult(input interface{}) bool {
	if input == nil {
		// TODO: Should be true?
		return false
	}

	if IsSpecial(input) {
		return true
	}
	return false
}

// Special datatypes, not to be handled yet.
func IsSpecial(input interface{}) bool {
	if input == nil {
		return false
	}

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
