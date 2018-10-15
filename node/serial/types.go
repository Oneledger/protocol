/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Category int

const (
	UNKNOWN = iota
	PRIMITIVE
	STRUCT
	MAP
	SLICE
	ARRAY
)

type TypeEntry struct {
	Name      string
	Category  Category
	DataType  reflect.Type
	KeyType   *TypeEntry
	ValueType *TypeEntry
}

func (entry TypeEntry) String() string {
	switch entry.Category {
	case UNKNOWN:
		return entry.Name + " UNKNOWN: " + entry.DataType.String()
	case PRIMITIVE:
		return entry.Name + " PRIMITIVE: " + entry.DataType.String()
	case STRUCT:
		return entry.Name + " STRUCT: " + entry.DataType.String()
	case MAP:
		return entry.Name + " MAP: " + entry.DataType.String()
	case SLICE:
		return entry.Name + " SLICE: " + entry.DataType.String()
	case ARRAY:
		return entry.Name + " ARRAY: " + entry.DataType.String()
	}
	return "Invalid!"
}

var dataTypes map[string]TypeEntry
var ignoreTypes map[string]TypeEntry

func init() {
	// Load in all of the primitives
	dataTypes = map[string]TypeEntry{
		"bool": TypeEntry{"boo", PRIMITIVE, reflect.TypeOf(bool(true)), nil, nil},

		"int":   TypeEntry{"int", PRIMITIVE, reflect.TypeOf(int(0)), nil, nil},
		"int8":  TypeEntry{"int8", PRIMITIVE, reflect.TypeOf(int8(0)), nil, nil},
		"int16": TypeEntry{"int16", PRIMITIVE, reflect.TypeOf(int16(0)), nil, nil},
		"int32": TypeEntry{"int32", PRIMITIVE, reflect.TypeOf(int32(0)), nil, nil},
		"int64": TypeEntry{"int64", PRIMITIVE, reflect.TypeOf(int64(0)), nil, nil},

		"uint":   TypeEntry{"uint", PRIMITIVE, reflect.TypeOf(uint(0)), nil, nil},
		"uint8":  TypeEntry{"uint8", PRIMITIVE, reflect.TypeOf(uint8(0)), nil, nil},
		"byte":   TypeEntry{"byte", PRIMITIVE, reflect.TypeOf(byte(0)), nil, nil},
		"uint16": TypeEntry{"uint16", PRIMITIVE, reflect.TypeOf(uint16(0)), nil, nil},
		"uint32": TypeEntry{"uint32", PRIMITIVE, reflect.TypeOf(uint32(0)), nil, nil},
		"uint64": TypeEntry{"uint64", PRIMITIVE, reflect.TypeOf(uint64(0)), nil, nil},

		"float32": TypeEntry{"float32", PRIMITIVE, reflect.TypeOf(float32(0)), nil, nil},
		"float64": TypeEntry{"float64", PRIMITIVE, reflect.TypeOf(float64(0)), nil, nil},

		"complex64":  TypeEntry{"complex64", PRIMITIVE, reflect.TypeOf(complex64(0)), nil, nil},
		"complex128": TypeEntry{"complex128", PRIMITIVE, reflect.TypeOf(complex128(0)), nil, nil},

		"string": TypeEntry{"string", PRIMITIVE, reflect.TypeOf(string("")), nil, nil},
	}
}

// Register a structure by its name
func Register(base interface{}) {

	// TODO: Not necessary?
	name := GetBaseTypeString(base)
	entry := GetTypeEntry(name, 1)
	if entry.Category != UNKNOWN {
		// Most often caused by byte and uint8 sort of being the same, but not in all cases.

		//log.Warn("Duplicate Entry", "name", name, "orig", entry.Name)
		//log.Dump("Dup is", base)
		//debug.PrintStack()
		//log.Dump("Exists", entry)
		return
	}

	typeOf := reflect.TypeOf(base)

	var category Category = PRIMITIVE
	if IsStructure(base) {
		dataTypes[name] = TypeEntry{name, STRUCT, typeOf, nil, nil}
		return
	}
	if IsPrimitiveContainer(base) {
		ubase := UnderlyingType(base)
		underType := GetTypeEntry(ubase.String(), 1)

		underType.Name = name
		dataTypes[name] = underType
		//dataTypes[name] = TypeEntry{name, underType.Category, typeOf, nil, nil}
		return
	}

	dataTypes[name] = TypeEntry{name, category, typeOf, nil, nil}
}

// Force an entry into the table
func RegisterForce(name string, category Category, dataType reflect.Type, keyType *TypeEntry, valueType *TypeEntry) {
	dataTypes[name] = TypeEntry{name, category, dataType, keyType, valueType}
}

func RegisterIgnore(base interface{}) {
	if ignoreTypes == nil {
		ignoreTypes = make(map[string]TypeEntry)
	}

	name := GetBaseTypeString(base)

	typeOf := reflect.TypeOf(base)
	var category Category = PRIMITIVE
	if IsStructure(base) {
		ignoreTypes[name] = TypeEntry{name, STRUCT, typeOf, nil, nil}
		return
	}
	if IsPrimitiveContainer(base) {
		ubase := UnderlyingType(base)
		ignoreType := GetTypeEntry(ubase.String(), 1)
		ignoreType.Name = name
		ignoreTypes[name] = ignoreType
		return
	}

	ignoreTypes[name] = TypeEntry{name, category, typeOf, nil, nil}
}

func GetTypeEntry(name string, size int) TypeEntry {
	name = strings.TrimPrefix(name, "*")

	// Static data types
	//log.Debug("Searching Statically for a match", "name", name)
	typeEntry, ok := dataTypes[name]
	if !ok {
		//log.Dump("Not Found in ", reflect.ValueOf(dataTypes).MapKeys())
		// dynamic data types -- like maps and slices
		//log.Debug("Searching Dynamically for a match", "name", name)
		entry := ParseType(name, size)
		return entry
	}
	return typeEntry
}

// Given a data type string, break it down into reflect.Type entries
func ParseType(name string, count int) TypeEntry {
	automata := regexp.MustCompile(`(.*)\[(.*)\](.*)`)
	groups := automata.FindStringSubmatch(name)

	if groups == nil || len(groups) != 4 {
		return TypeEntry{name, UNKNOWN, nil, nil, nil}
	}

	if groups[1] == "map" {
		keyTypeName := groups[2]
		valueTypeName := groups[3]
		keyType := GetTypeEntry(keyTypeName, 1)
		valueType := GetTypeEntry(valueTypeName, 1)
		finalType := reflect.MapOf(keyType.DataType, valueType.DataType)
		return TypeEntry{
			Name:      name,
			Category:  MAP,
			DataType:  finalType,
			KeyType:   &keyType,
			ValueType: &valueType,
		}

	} else if groups[1] == "" && groups[2] == "" {
		sliceTypeName := groups[3]
		sliceType := GetTypeEntry(sliceTypeName, count)
		//log.Dump(name+" has "+sliceTypeName, sliceType, groups)
		finalType := reflect.SliceOf(sliceType.DataType)
		return TypeEntry{
			Name:      name,
			Category:  SLICE,
			DataType:  finalType,
			ValueType: &sliceType,
		}

	} else {
		// TODO: What if this is a variable?
		//size := GetInt(groups[2], 0)
		arrayTypeName := groups[3]
		arrayType := GetTypeEntry(arrayTypeName, count)
		//log.Dump(name+" has "+arrayTypeName, arrayType, size, groups)
		finalType := reflect.ArrayOf(count, arrayType.DataType)
		return TypeEntry{
			Name:      name,
			Category:  ARRAY,
			DataType:  finalType,
			ValueType: &arrayType,
		}
	}
	return TypeEntry{name, UNKNOWN, reflect.Type(nil), nil, nil}
}

// TODO: This should be in the convert packahe, but it shares data with this one
func GetInt(value string, defaultValue int) int {

	// TODO: Should be ParseInt and should specific 64 or 32
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
