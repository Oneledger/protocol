/*
	Copyright 2017-2018 OneLedger
*/

package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

func SetStruct(base interface{}, fieldNum int, value interface{}) bool {
	valueOf := reflect.ValueOf(base)
	//elem := valueOf.Elem()
	elem := valueOf
	if elem.Kind() == reflect.Struct {
		field := elem.Field(fieldNum)
		if field.IsValid() {
			if field.CanSet() {
				if field.Kind() == reflect.Interface {
					//field.SetInt(value)
					field.Set(reflect.ValueOf(value).Elem())
					return true
				}
			}
		}
	}
	log.Warn("Unable to set", "base", base, "fieldNum", fieldNum, "value", value)
	return false
}

func Print(base interface{}) {
	action := &Action{ProcessField: PrintIt}

	Iterate(base, action)
}

func PrintIt(action *Action, input interface{}) bool {
	log.Debug("PrintIt", "action", action, "value", input)
	return true
}

// Clone and add in SerialWrapper
func Extend(base interface{}) interface{} {
	action := &Action{ProcessField: ExtendIt}

	Iterate(base, action)

	var last interface{}
	for key, value := range action.Children {
		log.Debug("Have a Final Child", "key", key)
		last = value
	}
	return last
}
func ExtendIt(action *Action, input interface{}) bool {
	log.Debug("ExtendIt", "struct", action.Struct, "value", input)
	return true
}

// Remove any SerialWrappers
func Contract(base interface{}) interface{} {
	action := &Action{ProcessField: ContractIt}

	Iterate(base, action)

	var last interface{}
	for key, value := range action.Children {
		log.Debug("Have a Final Child", "key", key)
		last = value
	}
	return last
}
func ContractIt(action *Action, input interface{}) bool {
	log.Debug("ContractIt", "struct", action.Struct, "value", input)
	return true
}

// Make a deep copy of the data
func Clone(base interface{}) interface{} {
	action := &Action{ProcessField: CloneIt}

	Iterate(base, action)

	var last interface{}
	for key, value := range action.Children {
		log.Debug("Have a Final Child", "key", key)
		last = value
	}
	return last
}

func CloneIt(action *Action, input interface{}) bool {
	log.Debug("CloneIt", "action", action, "input", input)

	if reflect.TypeOf(action.Value).Kind() == reflect.Struct {
		log.Debug("Have a structure, will copy!!!")
		copy := action.Value // Implicit Copy
		for key, value := range action.Children {
			SetStruct(copy, key, value)
			delete(action.Children, key)
		}
		log.Debug("Copy in Parent", "key", action.Field, "copy", copy)
		action.Children[action.Field] = copy
	}
	return true
}
