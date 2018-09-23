package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// Set a structure with a given value, convert as necessary
func Set(parent interface{}, fieldName string, child interface{}) bool {
	kind := reflect.ValueOf(parent).Kind()

	switch kind {

	case reflect.Struct:
		return SetStruct(parent, fieldName, child)

	case reflect.Map:
		return SetMap(parent, fieldName, child)

	case reflect.Ptr:
		return SetStruct(parent, fieldName, child)

	}

	return false
}

// SetStructure takes a pointer to a structure and sets it.
func SetStruct(parent interface{}, fieldName string, child interface{}) bool {
	if child == nil {
		return false
	}

	log.Debug("SetStruct", "parent", parent, "fieldName", fieldName)

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

	/*
		if field.Kind() != reflect.Interface {
			log.Warn("Field is not an interface", "kind", field.Kind(), "field", field)
			return false
		}
	*/

	newValue := ConvertValue(child, field.Type())

	field.Set(newValue)

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
		log.Debug("Converted", "result", result)
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
