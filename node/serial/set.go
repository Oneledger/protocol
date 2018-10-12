package serial

import (
	"reflect"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
)

func Alloc(dataType string, size int) interface{} {

	if dataType == "" {
		return nil
	}

	entry := GetTypeEntry(dataType)

	switch entry.Category {
	case PRIMITIVE:
		// Don't need to alloc, only containers
	case STRUCT:
		return reflect.New(entry.DataType).Interface()
	case MAP:
		return reflect.MakeMapWithSize(entry.DataType, size).Interface()
	case SLICE:
		return reflect.MakeSlice(entry.DataType, 0, size).Interface()
	}
	return nil
}

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
		log.Fatal("Not a structure", "element", element, "kind", element.Kind())
	}

	if !CheckValue(element) {
		return false
	}

	field := element.FieldByName(fieldName)

	if !CheckValue(field) {
		return false
	}

	/*
		if field.Type().Kind() == reflect.Interface {
			value := reflect.ValueOf(child)
			log.Debug("Trying to set Interface", "child", child, "type", reflect.TypeOf(child), "value", value)

			field.Set(value)

		} else {
	*/
	newValue := ConvertValue(child, field.Type())
	field.Set(newValue)

	/*
		}
	*/
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
		log.Debug("##### CONVERTING JSON NUMBER TO GO BYTE")
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
	element := GetValue(parent)

	if element.Kind() != reflect.Map {
		log.Fatal("Not a map", "element", element)
	}

	if !CheckValue(element) {
		return false
	}

	key := reflect.ValueOf(fieldName)
	fieldType := GetType(parent).Elem()

	newValue := ConvertValue(child, fieldType)

	log.Dump("key", key, "newValue", newValue)

	element.SetMapIndex(key, newValue)

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetSlice(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Slice {
		log.Fatal("Not a structure", "element", element)
	}

	if !CheckValue(element) {
		return false
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
		log.Warn("Not Settable", "element", element)
		return false
	}
	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetArray(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetValue(parent)

	if element.Kind() != reflect.Array {
		log.Fatal("Not a structure", "element", element)
	}

	if !CheckValue(element) {
		return false
	}

	newValue := ConvertValue(child, element.Index(index).Type())
	element.Index(index).Set(newValue)

	return true
}
