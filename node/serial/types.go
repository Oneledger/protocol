/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

type Category int

const (
	UNKNOWN = iota
	INTERFACE
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

func (entry TypeEntry) String() string {
	switch entry.Category {
	case UNKNOWN:
		return fmt.Sprintf("%s UNKNOWN", entry.Name)
	case INTERFACE:
		return fmt.Sprintf("%s INTERFACE", entry.Name)
	case PRIMITIVE:
		if entry.DataType == nil {
			return fmt.Sprintf("%s PRIMITIVE (missing)", entry.Name)
		}
		return fmt.Sprintf("%s PRIMITIVE %s", entry.Name, entry.DataType.String())
	case STRUCT:
		return fmt.Sprintf("%s STRUCT %s", entry.Name, entry.DataType.String())
	case MAP:
		return fmt.Sprintf("%s MAP %s (%s,%s)", entry.Name, entry.DataType.String(),
			entry.KeyType.String(), entry.ValueType.String())
	case SLICE:
		return fmt.Sprintf("%s SLICE %s (%s)", entry.Name, entry.DataType.String(), entry.ValueType.String())
	case ARRAY:
		if entry.ValueType == nil {
			return fmt.Sprintf("%s ARRAY %s (missing)", entry.Name, entry.DataType.String())
		}
		return fmt.Sprintf("%s ARRAY %s (%s)", entry.Name, entry.DataType.String(), entry.ValueType.String())
	}
	return "Invalid!"
}

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

func DumpTypes() {
	log.Dump("Known Data Types are:", dataTypes)
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

func RegisterInterface(base interface{}) {
	typeOf := reflect.TypeOf(base)
	element := typeOf.Elem()
	name := element.String()

	entry := GetTypeEntry(name, 1)
	if entry.Category != UNKNOWN {
		// Already registered
		return
	}

	var category Category = INTERFACE

	dataTypes[name] = TypeEntry{name, category, element, element, nil, nil}

	//DumpTypes()
}

// Register a structure by its name
func Register(base interface{}) {

	name := GetBaseTypeString(base)

	entry := GetTypeEntry(name, 1)
	if entry.Category != UNKNOWN {
		// Already registered
		return
	}

	var category Category = PRIMITIVE

	typeOf := reflect.TypeOf(base)

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
		return
	}

	if IsInterface(base) {
		ubase := UnderlyingType(base)
		underType := GetTypeEntry(ubase.String(), 1)

		underType.Name = name
		underType.RootType = typeOf
		dataTypes[name] = underType
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
	Register(base)

	name := GetBaseTypeString(base)

	typeOf := reflect.TypeOf(base)
	var category Category = PRIMITIVE
	/*
		if IsStructure(base) {
			ignoreTypes[name] = TypeEntry{name, STRUCT, typeOf, typeOf, nil, nil}
			return
		}
		if IsPrimitiveContainer(base) {
			ubase := UnderlyingType(base)
			underType := GetTypeEntry(ubase.String(), 1)
			underType.Name = name
			underType.RootType = typeOf
			ignoreTypes[name] = ignoreType
			return
		}
	*/
	ignoreTypes[name] = TypeEntry{name, category, typeOf, typeOf, nil, nil}
}

func IgnoreType(dataType reflect.Type) bool {
	if _, ok := ignoreTypes[dataType.String()]; ok {
		return true
	}
	return false
}

func IgnoreVariable(value interface{}) bool {
	return IgnoreType(GetBaseType(value))
}

func GetTypeEntry(name string, size int) TypeEntry {
	name = strings.TrimPrefix(name, "*")

	// Static data types
	typeEntry, ok := dataTypes[name]
	if !ok {
		// dynamic data types -- like maps and slices
		entry := ParseType(name, size)
		return entry
	}
	return typeEntry
}

// Given a data type string, break it down into reflect.Type entries
func ParseType(name string, count int) TypeEntry {
	if name == "" {
		return TypeEntry{name, UNKNOWN, reflect.Type(nil), reflect.Type(nil), nil, nil}
	}

	automata := regexp.MustCompile(`(.*?)\[(.*?)\](.*)`)
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

		var finalType reflect.Type
		if sliceType.Category == UNKNOWN {
			return TypeEntry{name, UNKNOWN, reflect.Type(nil), reflect.Type(nil), nil, nil}
			/*
				var prototype interface{}
				finalType = reflect.SliceOf(reflect.TypeOf(prototype))
			*/
		} else {
			finalType = reflect.SliceOf(sliceType.DataType)
		}

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

// TODO: This should be in the convert package, but it shares data with this one
func GetInt(value string, defaultValue int) int {

	// TODO: Should be ParseInt and should specific 64 or 32
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
