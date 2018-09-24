package comm

import (
	"fmt"
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// Action is the context for the iteration, each ProcessField function gets an updated pointer
type Action struct {
	Path  string
	Name  string
	Value interface{}

	ProcessField func(*Action, interface{}) interface{}

	ParentName  string
	ChildNumber int
	ChildName   string

	Processed map[string]Parameters
}

type Parameters struct {
	Children map[string]interface{}
}

// GetValue returns the underlying value, even if it is a pointer
func GetValue(base interface{}) reflect.Value {
	element := reflect.ValueOf(base)
	if element.Kind() == reflect.Ptr {
		log.Warn("Have an unexpected pointer!")
		return element.Elem()
	}
	return element
}

// Extract this info once, even though it is used in multiple levels of the recursion
type Child struct {
	Kind   reflect.Kind
	Number int
	Name   string
	Value  interface{}
}

// Get the Fields from a structure, and return them into a field array
func GetChildren(input interface{}) []Child {
	kind := reflect.TypeOf(input).Kind()

	switch kind {
	case reflect.Struct:
		return GetChildrenStruct(input)

	case reflect.Map:
		return GetChildrenMap(input)

	case reflect.Array:
		return GetChildrenArray(input)

	case reflect.Slice:
		return GetChildrenSlice(input)
	}

	return []Child{}
}

func GetChildrenStruct(input interface{}) []Child {
	typeOf := reflect.TypeOf(input)
	valueOf := GetValue(input)

	count := valueOf.NumField()

	var children []Child
	children = make([]Child, count)

	for i := 0; i < count; i++ {
		name := typeOf.Field(i).Name
		value := valueOf.Field(i).Interface()
		kind := valueOf.Field(i).Kind()

		children[i] = Child{Name: name, Number: i, Value: value, Kind: kind}
	}

	return children
}

func GetChildrenMap(input interface{}) []Child {
	valueOf := GetValue(input)

	var children []Child
	children = make([]Child, 0)

	for _, key := range valueOf.MapKeys() {
		name := key.String()
		value := valueOf.MapIndex(key).Interface()
		kind := reflect.ValueOf(value).Kind()

		children = append(children, Child{Name: name, Value: value, Kind: kind})
	}

	return children
}

func GetChildrenSlice(input interface{}) []Child {
	valueOf := GetValue(input)

	var children []Child
	children = make([]Child, 0)

	for i := 0; i < valueOf.Len(); i++ {
		name := fmt.Sprintf("%d", i)
		value := valueOf.Index(i)
		kind := reflect.ValueOf(value).Kind()
		children = append(children, Child{Name: name, Value: value, Kind: kind})
	}

	return children
}

func GetChildrenArray(input interface{}) []Child {
	valueOf := GetValue(input)

	var children []Child
	children = make([]Child, 0)

	for i := 0; i < valueOf.Len(); i++ {
		name := fmt.Sprintf("%d", i)
		value := valueOf.Index(i)
		kind := reflect.ValueOf(value).Kind()
		children = append(children, Child{Name: name, Value: value, Kind: kind})
	}

	return children
}

// Iterate the variables in memory, executing functions at each node in the traversal
func Iterate(input interface{}, action *Action) interface{} {
	// TODO: add in cycle detection

	// Initialize on first call
	if action.Processed == nil {
		action.Processed = make(map[string]Parameters, 0)
		action.Processed[""] = Parameters{make(map[string]interface{}, 0)}
	}

	if IsDifficult(input) {
		log.Fatal("Can't deal with this", "input", input)
	}

	//log.Debug("Iterate", "input", input)

	action.Value = input
	action.ParentName = action.Name
	action.Name = action.ChildName

	parent := action.ParentName
	name := action.Name

	if IsPointer(input) {
		element := reflect.ValueOf(input).Elem().Interface()

		action.ChildName = action.Name
		action.ChildNumber = 0

		Iterate(element, action)

		// Restore the action values, since they were overwritten
		action.Value = input
		action.Name = name
		action.ParentName = parent

		return action.ProcessField(action, input)
	}

	// Walk the children first -- post-order traversal
	if IsContainer(input) {
		action.Processed[action.Name] = Parameters{make(map[string]interface{}, 0)}

		action.Path = reflect.TypeOf(input).Name()

		children := GetChildren(input)
		for i := 0; i < len(children); i++ {
			action.ChildName = children[i].Name
			action.ChildNumber = i

			Iterate(children[i].Value, action)

			// Restore the action values, since they were overwritten
			action.Value = input
			action.Name = name
			action.ParentName = parent
		}
	}

	//log.Debug("Iterate after Children", "input", input)
	result := action.ProcessField(action, input)
	return result
}
