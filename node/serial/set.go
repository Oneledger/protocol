package serial

import (
	"math/big"
	"reflect"
	"runtime/debug"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
)

// Alloc a new variable of any type
func Alloc(dataType string, size int) interface{} {
	if dataType == "" {
		return nil
	}

	if size == -1 {
		return nil
	}

	entry := GetTypeEntry(dataType, size)

	var value reflect.Value

	switch entry.Category {
	case UNKNOWN:
		log.Fatal("Unknown datatype", "dataType", dataType)

	case PRIMITIVE:
		// Don't need to alloc, only containers
		return nil

	case STRUCT:
		value = reflect.New(entry.DataType)
		if !value.IsValid() {
			log.Warn("Retrying as slice?")
			value = reflect.MakeSlice(entry.DataType, size, size)
			if !value.IsValid() {
				log.Warn("Retrying as byte array")
				value = reflect.ValueOf(make([]byte, size))
			}
		}
		return value.Interface()

	case MAP:
		smap := reflect.MakeMapWithSize(entry.RootType, size)
		value = reflect.New(smap.Type())
		value.Elem().Set(smap)
		return value.Interface()

	case SLICE:
		slice := reflect.MakeSlice(entry.DataType, size, size)
		value = reflect.New(slice.Type())
		value.Elem().Set(slice)
		return value.Interface()

	case ARRAY:
		array := reflect.ArrayOf(size, entry.ValueType.DataType)
		value = reflect.New(array)
		return value.Interface()
	}

	log.Warn("Unknown Category", "dataType", dataType, "entry", entry)
	return nil
}

// Set a structure with a given value, convert as necessary
func Set(parent interface{}, fieldName string, child interface{}) (status bool) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Set Panic", "r", r)
			log.Dump("Parameters", parent, fieldName, child)
			debug.PrintStack()
		}
	}()

	kind := GetBaseValue(parent).Kind()

	switch kind {

	case reflect.Struct:
		//log.Dump("The parent", parent, kind)
		return SetStruct(parent, fieldName, child)

	case reflect.Map:
		return SetMap(parent, fieldName, child)

	case reflect.Slice:
		index, err := strconv.Atoi(fieldName)
		if err != nil {
			log.Fatal("Invalid Index", "fieldName", fieldName)
		}
		return SetSlice(parent, index, child)

	case reflect.Array:
		index, err := strconv.Atoi(fieldName)
		if err != nil {
			log.Fatal("Invalid Index", "fieldName", fieldName)
		}
		return SetArray(parent, index, child)
	}
	return false
}

// SetStructure takes a pointer to a structure and sets it.
func SetStruct(parent interface{}, fieldName string, child interface{}) bool {
	if child == nil {
		return false
	}

	// Convert the interfaces to structures
	//element := reflect.ValueOf(parent).Elem()
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Struct {
		log.Fatal("Not a structure", "element", element, "kind", element.Kind())
	}

	if !CheckValue(element) {
		return false
	}

	field := element.FieldByName(fieldName)

	if !CheckValue(field) {
		return false
	}

	if field.Type().Kind() == reflect.Interface {
		// When setting to a generic interface{}

		value := reflect.ValueOf(child)
		field.Set(value)

	} else {
		newValue := ConvertValue(child, field.Type())
		if newValue.Kind() == reflect.Int {
			field.SetInt(newValue.Int())
		} else {
			field.Set(newValue)
		}
	}
	return true
}

// Convert any value to an arbitrary type
func ConvertValue(value interface{}, fieldType reflect.Type) reflect.Value {
	if value == nil {
		return reflect.ValueOf(nil)
	}

	typeOf := reflect.TypeOf(value)
	valueOf := reflect.ValueOf(value)

	// Remove any pointers, if they exist
	if typeOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
		typeOf = reflect.TypeOf(valueOf)
	}

	switch typeOf.Kind() {

	case reflect.Float64:
		// JSON returns floats for everything :-(
		result := ConvertNumber(fieldType, valueOf)
		return result

	case reflect.String:
		if fieldType.String() == "data.ChainType" {
			entry := GetTypeEntry(fieldType.String(), 1)
			interim, err := strconv.ParseInt(GetString(valueOf), 10, 0)
			if err != nil {
				log.Fatal("Failed to convert int")
			}
			valueof := reflect.New(entry.RootType)
			element := valueof.Elem()
			element.SetInt(interim)
			return element
		}
		if fieldType.Kind() == reflect.Int {
			var result int
			interim, err := strconv.ParseInt(GetString(valueOf), 10, 0)
			if err != nil {
				log.Fatal("Failed to convert int")
			}
			result = int(interim)
			return reflect.ValueOf(result)
		}
	}
	return valueOf
}

func GetString(value reflect.Value) string {
	if value.Kind() == reflect.String {
		return value.String()
	}
	return value.Elem().String()
}

// ConvertNumber handles JSON numbers as they are float64 from the parser
func ConvertNumber(fieldType reflect.Type, value reflect.Value) reflect.Value {

	// TODO: find a better way of handling types that are not structures
	/*
		if fieldType.String() == "data.ChainType" {
			entry := GetTypeEntry(fieldType.String(), 1)
			valueof := reflect.New(entry.DataType)
			var result int64
			result = int64(value.Float())
			element := valueof.Elem()
			element.SetInt(result)
			return element
		}
	*/
	/*
		if fieldType.String() == "action.Type" {
			entry := GetTypeEntry(fieldType.String(), 1)
			valueof := reflect.New(entry.DataType)

			var result int64
			result = int64(value.Float())

			element := valueof.Elem()
			element.SetInt(result)
			return element
		}
	*/

	// TODO: shouldn't be manaually creating big ints
	if fieldType.String() == "*big.Int" {
		converted := big.NewInt(int64(value.Float()))
		return reflect.ValueOf(converted)
	}

	switch fieldType.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(value.Float()))

	case reflect.Int8:
		return reflect.ValueOf(int8(value.Float()))

	case reflect.Int16:
		return reflect.ValueOf(int16(value.Float()))

	case reflect.Int32:
		return reflect.ValueOf(int32(value.Float()))

	case reflect.Int64:
		return reflect.ValueOf(int64(value.Float()))

	case reflect.Uint:
		return reflect.ValueOf(uint(value.Float()))

	case reflect.Uint8:
		return reflect.ValueOf(uint8(value.Float()))

	case reflect.Uint16:
		return reflect.ValueOf(uint16(value.Float()))

	case reflect.Uint32:
		return reflect.ValueOf(uint32(value.Float()))

	case reflect.Uint64:
		return reflect.ValueOf(uint64(value.Float()))

	case reflect.Float32:
		return reflect.ValueOf(float32(value.Float()))
	}
	return value
}

// SetMap takes a pointer to a structure and sets it.
func SetMap(parent interface{}, fieldName string, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Map {
		log.Fatal("Not a map", "element", element)
	}

	if !CheckValue(element) {
		return false
	}

	//key := reflect.ValueOf(fieldName)
	keyType := GetBaseType(parent).Key()
	newKey := ConvertValue(fieldName, keyType)

	fieldType := GetBaseType(parent).Elem()
	newValue := ConvertValue(child, fieldType)

	/*
		if element.Len() < 1 {
			// TODO: Need to figure out a reasonable size here...
			log.Warn("Reallocating Map", "type", element.Type().String())
			entry := GetTypeEntry(element.Type().String(), 100)
			revised := reflect.MakeMapWithSize(entry.RootType, 100)
			value := reflect.New(revised.Type())
			value.Elem().Set(element)
			element = value.Elem()
		}
	*/

	element.SetMapIndex(newKey, newValue)

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetSlice(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Slice {
		log.Fatal("Not a slice", "element", element)
	}

	if !CheckValue(element) {
		return false
	}

	if element.Len() < index {
		log.Warn("Reallocating Slice")
		element = reflect.MakeSlice(element.Index(0).Type(), index+1, index+1)
	}
	newValue := ConvertValue(child, element.Index(index).Type())
	element.Index(index).Set(newValue)

	return true
}

func CheckValue(element reflect.Value) bool {
	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Element not Settable", "element", element)
		return false
	}
	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetArray(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Array {
		log.Fatal("Not a structure", "element", element)
	}

	if !CheckValue(element) {
		return false
	}

	if element.Len() < index {
		log.Warn("Reallocating Array")
	}

	cell := element.Index(index)
	newValue := ConvertValue(child, cell.Type())

	element.Index(index).Set(newValue)

	return true
}
