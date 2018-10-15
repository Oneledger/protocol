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
	RootType  reflect.Type
	DataType  reflect.Type
	KeyType   *TypeEntry
	ValueType *TypeEntry
}

/*
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
*/

var dataTypes map[string]TypeEntry
var ignoreTypes map[string]TypeEntry

func NewPrimitiveEntry(name string, category Category, root reflect.Type) *TypeEntry {
	return &TypeEntry{
		Name:      name,
		Category:  category,
		RootType:  root,
		DataType:  root,
		KeyType:   nil,
		ValueType: nil,
	}
}

func init() {
	// Load in all of the primitives
	dataTypes = map[string]TypeEntry{
		"bool": *NewPrimitiveEntry("boo", PRIMITIVE, reflect.TypeOf(bool(true))),

		"int":   *NewPrimitiveEntry("int", PRIMITIVE, reflect.TypeOf(int(0))),
		"int8":  *NewPrimitiveEntry("int8", PRIMITIVE, reflect.TypeOf(int8(0))),
		"int16": *NewPrimitiveEntry("int16", PRIMITIVE, reflect.TypeOf(int16(0))),
		"int32": *NewPrimitiveEntry("int32", PRIMITIVE, reflect.TypeOf(int32(0))),
		"int64": *NewPrimitiveEntry("int64", PRIMITIVE, reflect.TypeOf(int64(0))),

		"uint":   *NewPrimitiveEntry("uint", PRIMITIVE, reflect.TypeOf(uint(0))),
		"uint8":  *NewPrimitiveEntry("uint8", PRIMITIVE, reflect.TypeOf(uint8(0))),
		"byte":   *NewPrimitiveEntry("byte", PRIMITIVE, reflect.TypeOf(byte(0))),
		"uint16": *NewPrimitiveEntry("uint16", PRIMITIVE, reflect.TypeOf(uint16(0))),
		"uint32": *NewPrimitiveEntry("uint32", PRIMITIVE, reflect.TypeOf(uint32(0))),
		"uint64": *NewPrimitiveEntry("uint64", PRIMITIVE, reflect.TypeOf(uint64(0))),

		"float32": *NewPrimitiveEntry("float32", PRIMITIVE, reflect.TypeOf(float32(0))),
		"float64": *NewPrimitiveEntry("float64", PRIMITIVE, reflect.TypeOf(float64(0))),

		"complex64":  *NewPrimitiveEntry("complex64", PRIMITIVE, reflect.TypeOf(complex64(0))),
		"complex128": *NewPrimitiveEntry("complex128", PRIMITIVE, reflect.TypeOf(complex128(0))),

		"string": *NewPrimitiveEntry("string", PRIMITIVE, reflect.TypeOf(string(""))),
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
		dataTypes[name] = TypeEntry{name, STRUCT, typeOf, typeOf, nil, nil}
		return
	}
	if IsPrimitiveContainer(base) {
		ubase := UnderlyingType(base)
		underType := GetTypeEntry(ubase.String(), 1)

		underType.Name = name
		underType.RootType = typeOf
		dataTypes[name] = underType
		//log.Dump("Full Entry is", dataTypes[name])
		return
	}

	dataTypes[name] = TypeEntry{name, category, typeOf, typeOf, nil, nil}
}

// Force an entry into the table
func RegisterForce(name string, category Category, dataType reflect.Type, keyType *TypeEntry, valueType *TypeEntry) {
	dataTypes[name] = TypeEntry{name, category, dataType, dataType, keyType, valueType}
}

func RegisterIgnore(base interface{}) {
	if ignoreTypes == nil {
		ignoreTypes = make(map[string]TypeEntry)
	}

	name := GetBaseTypeString(base)

	typeOf := reflect.TypeOf(base)
	var category Category = PRIMITIVE
	if IsStructure(base) {
		ignoreTypes[name] = TypeEntry{name, STRUCT, typeOf, typeOf, nil, nil}
		return
	}
	if IsPrimitiveContainer(base) {
		ubase := UnderlyingType(base)
		ignoreType := GetTypeEntry(ubase.String(), 1)
		ignoreType.Name = name
		ignoreTypes[name] = ignoreType
		return
	}

	ignoreTypes[name] = TypeEntry{name, category, typeOf, typeOf, nil, nil}
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
		return TypeEntry{name, UNKNOWN, reflect.Type(nil), reflect.Type(nil), nil, nil}
	}

	if groups[1] == "map" {
		keyTypeName := groups[2]
		valueTypeName := groups[3]
		keyType := GetTypeEntry(keyTypeName, 1)
		valueType := GetTypeEntry(valueTypeName, 1)
		finalType := reflect.MapOf(keyType.RootType, valueType.RootType)
		return TypeEntry{
			Name:      name,
			Category:  MAP,
			RootType:  finalType,
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
			RootType:  finalType,
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
			RootType:  finalType,
			DataType:  finalType,
			ValueType: &arrayType,
		}
	}
	return TypeEntry{name, UNKNOWN, reflect.Type(nil), reflect.Type(nil), nil, nil}
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
