package serial

import (
	"reflect"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
)

// Set a structure with a given value, convert as necessary
func Set(parent interface{}, fieldName string, child interface{}) (status bool) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Set Panic", "r", r)
			status = false
		}
	}()

	kind := reflect.ValueOf(parent).Kind()

	if kind == reflect.Ptr {
		kind = reflect.ValueOf(parent).Elem().Kind()
	}

	switch kind {

	case reflect.Struct:
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
	element := reflect.ValueOf(parent).Elem()

	if element.Kind() != reflect.Struct {
		log.Warn("Not a structure", "element", element, "kind", element.Kind())
		return false
	}

	if !element.CanSet() {
		log.Warn("Structure not settable", "element", element)
		return false
	}

	field := element.FieldByName(fieldName)
	if !field.IsValid() {
		log.Warn("Field not found", "fieldName", fieldName, "field", field,
			"element", element, "fieldName", fieldName)
		return false
	}

	if !field.CanSet() {
		log.Warn("Not Settable", "field", field, "fieldName", fieldName)
		return false
	}

	if field.Type().Kind() == reflect.Interface {
		value := reflect.ValueOf(child)
		log.Debug("Trying to set Interface", "child", child, "type", reflect.TypeOf(child), "value", value)

		field.Set(value)

	} else {
		newValue := ConvertValue(child, field.Type())
		field.Set(newValue)
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
		return ConvertNumber(fieldType.Kind(), valueOf)

	case reflect.String:
		if fieldType.Kind() == reflect.Int {
			log.Debug("CONVERTING STR TO INT")
			result, err := strconv.ParseInt(valueOf.Elem().String(), 10, 0)
			if err != nil {
				log.Fatal("Failed to convert int")
			}
			return reflect.ValueOf(result)
		}
	}
	return valueOf
}

// ConvertNumber handles JSON numbers as they are float64 from the parser
func ConvertNumber(kind reflect.Kind, value reflect.Value) reflect.Value {
	switch kind {
	case reflect.Int8:
		return reflect.ValueOf(int8(value.Float()))

	case reflect.Int16:
		return reflect.ValueOf(int16(value.Float()))

	case reflect.Int32:
		return reflect.ValueOf(int32(value.Float()))

	case reflect.Int64:
		return reflect.ValueOf(int64(value.Float()))

	case reflect.Int:
		return reflect.ValueOf(int(value.Float()))

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

func AllocMap(parent interface{}, size int) bool {
	element := GetValue(parent)

	if element.Kind() != reflect.Map {
		log.Warn("Not a structure", "element", element)
		return false
	}

	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Not Settable", "element", element)
		return false
	}

	element.Set(reflect.MakeMapWithSize(element.Type(), size))

	return true
}

// SetMap takes a pointer to a structure and sets it.
func SetMap(parent interface{}, fieldName string, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Map {
		log.Warn("Not a structure", "element", element)
		return false
	}

	key := reflect.ValueOf(fieldName)

	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Not Settable", "element", element)
		return false
	}

	newValue := reflect.ValueOf(child)

	log.Dump("key", key, "newValue", newValue)

	element.SetMapIndex(key, newValue)

	return true
}

func AllocSlice(parent interface{}, size int) bool {
	element := GetValue(parent)

	if element.Kind() != reflect.Map {
		log.Warn("Not a structure", "element", element)
		return false
	}

	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Not Settable", "element", element)
		return false
	}

	element.Set(reflect.MakeSlice(element.Type(), 0, size))

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetSlice(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Slice {
		log.Warn("Not a structure", "element", element)
		return false
	}

	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Not Settable", "element", element)
		return false
	}

	newValue := reflect.ValueOf(child)
	element.Index(index).Set(newValue)

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetArray(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Array {
		log.Warn("Not a structure", "element", element)
		return false
	}

	if !element.IsValid() {
		log.Warn("Map is invalid", "element", element)
		return false
	}

	if !element.CanSet() {
		log.Warn("Not Settable", "element", element)
		return false
	}

	newValue := reflect.ValueOf(child)
	element.Index(index).Set(newValue)

	return true
}
