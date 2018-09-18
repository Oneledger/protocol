package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// Action is the context for the iteration, each ProcessField function gets an updated pointer
type Action struct {
	Struct       string
	Name         string
	Field        int
	Value        interface{}
	ProcessField func(*Action, interface{}) bool
	Children     map[int]interface{}
}

// Extract this info once, even though it is used in multiple levels of the recursion
type Field struct {
	Name  string
	Value interface{}
	Kind  reflect.Kind
}

// Get the Fields from a structure, and return them into a field array
func GetFields(input interface{}) []Field {
	typeOf := reflect.TypeOf(input)
	valueOf := reflect.ValueOf(input)

	count := valueOf.NumField()

	var fields []Field
	fields = make([]Field, count)

	for i := 0; i < count; i++ {
		name := typeOf.Field(i).Name
		value := valueOf.Field(i).Interface()
		kind := valueOf.Field(i).Kind()
		fields[i] = Field{Name: name, Value: value, Kind: kind}
	}
	return fields
}

// Iterate the variables in memory, executing functions at each node in the traversal
func Iterate(input interface{}, action *Action) {
	// TODO: add in cycle detection

	if action.Children == nil {
		log.Debug("Allocating memory!")
		action.Children = make(map[int]interface{}, 0)
	}

	if IsDifficult(input) {
		log.Fatal("Can't deal with this", "input", input)
	}

	action.Value = input

	// Visit the childten first
	if IsStructure(input) {
		action.Struct = reflect.TypeOf(input).Name()

		fields := GetFields(input)
		for i := 0; i < len(fields); i++ {
			action.Name = fields[i].Name
			action.Field = i
			Iterate(fields[i].Value, action)
		}
	}

	action.Value = input
	action.ProcessField(action, input)
}
