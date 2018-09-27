package comm

import (
	"reflect"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
)

// Set a structure with a given value, convert as necessary
func Set(parent interface{}, fieldName string, child interface{}) bool {
	kind := reflect.ValueOf(parent).Kind()

	if kind == reflect.Ptr {
		kind = reflect.ValueOf(parent).Elem().Kind()
	}

	// TODO: Set should be able to handle points properly
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
	if typeOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
		typeOf = reflect.TypeOf(valueOf)
	}

	if typeOf.Kind() == reflect.Float64 && fieldType.Kind() == reflect.Int {
		result := int(valueOf.Float())
		return reflect.ValueOf(result)
	}
	return valueOf
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
	element.SetMapIndex(key, newValue)

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
