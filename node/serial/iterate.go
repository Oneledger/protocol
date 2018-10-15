package serial

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/Oneledger/protocol/node/log"
)

// Action is the context for the iteration, each ProcessField function gets an updated pointer
type Action struct {
	// Config Items
	IgnorePrimitives bool

	// Current Values
	Path      Stack
	Name      string
	Index     int
	IsPointer bool

	// Processing Function
	ProcessField func(*Action, interface{}) interface{}

	// Kids that have already been processed
	Processed map[string]Parameters
}

type Parameters struct {
	Children map[string]interface{}
}

// GetBaseType returns the underlying value, even if it is a pointer
func GetBaseType(base interface{}) reflect.Type {
	element := reflect.TypeOf(base)
	if element.Kind() == reflect.Ptr {
		return element.Elem()
	}
	return element
}

// GetBaseValue returns the underlying value, even if it is a pointer
func GetBaseValue(base interface{}) reflect.Value {
	element := reflect.ValueOf(base)
	if element.Kind() == reflect.Ptr {
		return element.Elem()
	}
	return element
}

// Find the type string that matches this variable
func GetBaseTypeString(base interface{}) string {
	valueOf := GetBaseValue(base)
	return valueOf.Type().String()
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
	typeOf := reflect.TypeOf(input)

	// TODO: Shouldn't need to manually ignore recursion
	if typeOf.String() == "big.Int" {
		return []Child{}
	}

	kind := typeOf.Kind()

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

// Get Children from a structure
func GetChildrenStruct(input interface{}) []Child {
	typeOf := reflect.TypeOf(input)
	valueOf := GetBaseValue(input)

	count := valueOf.NumField()

	var children []Child
	children = make([]Child, count)

	for i := 0; i < count; i++ {
		fieldType := typeOf.Field(i)
		field := valueOf.Field(i)

		//if element.IsValid() && element.CanInterface() {
		if field.CanInterface() {
			name := fieldType.Name
			value := field.Interface()
			kind := field.Kind()

			children[i] = Child{Name: name, Number: i, Value: value, Kind: kind}
		} else {
			// Have to recreate the parent, so be able to recreate the child...
			avalueOf := reflect.New(valueOf.Type()).Elem()
			avalueOf.Set(valueOf)
			rfield := avalueOf.Field(i)
			afield := reflect.NewAt(rfield.Type(), unsafe.Pointer(rfield.UnsafeAddr())).Elem()
			name := fieldType.Name
			value := afield.Interface()
			kind := afield.Kind()

			children[i] = Child{Name: name, Number: i, Value: value, Kind: kind}
		}
	}
	return children
}

// Get Children from a Map
func GetChildrenMap(input interface{}) []Child {
	valueOf := GetBaseValue(input)

	var children []Child
	children = make([]Child, 0)

	for _, key := range valueOf.MapKeys() {
		value := valueOf.MapIndex(key)
		kind := value.Kind()

		name := Value2String(key)
		children = append(children, Child{Name: name, Value: value.Interface(), Kind: kind})
	}
	return children
}

func Value2String(key reflect.Value) string {
	switch key.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(key.Int(), 10)
	case reflect.String:
		return key.String()
	}
	return key.String()
}

// Get Children from a slice
func GetChildrenSlice(input interface{}) []Child {
	typeOf := reflect.TypeOf(input)
	valueOf := GetBaseValue(input)

	var children []Child
	children = make([]Child, 0)

	// TODO: Check this, optimization?
	if typeOf.String() == "[]byte" {
		return children
	}

	for i := 0; i < valueOf.Len(); i++ {
		value := valueOf.Index(i).Interface()
		kind := reflect.ValueOf(value).Kind()

		// Use a string index.
		name := fmt.Sprintf("%d", i)
		children = append(children, Child{Name: name, Value: value, Kind: kind})
	}
	return children
}

// Get children from an array
func GetChildrenArray(input interface{}) []Child {
	valueOf := GetBaseValue(input)

	var children []Child
	children = make([]Child, 0)

	for i := 0; i < valueOf.Len(); i++ {
		value := valueOf.Index(i).Interface()
		kind := reflect.ValueOf(value).Kind()

		name := fmt.Sprintf("%d", i)
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
		action.Path = *NewStack()
		action.Path.Push("root")
	}

	// Some types of not implemented
	if IsDifficult(input) {
		log.Fatal("Can't deal with this", "input", input)
	}

	// Short cut if specified
	if action.IgnorePrimitives && IsPrimitive(input) {
		return input
	}

	parent := action.Path.StringPeekN(0)
	action.Path.Push(action.Name)

	if IsPointer(input) {
		if !IsNilValue(input) {
			input = reflect.ValueOf(input).Elem().Interface()
		}
		action.IsPointer = true
	} else {
		action.IsPointer = false
	}

	if action.Processed[parent].Children == nil {
		action.Processed[parent] = Parameters{make(map[string]interface{}, 0)}
	}

	// Walk the children first -- post-order traversal
	if IsContainer(input) {

		// Save the original values
		name := action.Path.StringPeekN(0)
		pointer := action.IsPointer

		children := GetChildren(input)
		for i := 0; i < len(children); i++ {
			action.Name = children[i].Name

			Iterate(children[i].Value, action)

			// Restore the action values, since they were overwritten
			action.Name = name
			action.IsPointer = pointer
		}
	}

	result := action.ProcessField(action, input)
	action.Path.Pop()

	return result
}
