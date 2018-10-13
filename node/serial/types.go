/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/node/log"
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
	Category Category
	DataType reflect.Type
}

func (entry TypeEntry) String() string {
	switch entry.Category {
	case UNKNOWN:
		return "UNKNOWN: " + entry.DataType.String()
	case PRIMITIVE:
		return "PRIMITIVE: " + entry.DataType.String()
	case STRUCT:
		return "STRUCT: " + entry.DataType.String()
	case MAP:
		return "MAP: " + entry.DataType.String()
	case SLICE:
		return "SLICE: " + entry.DataType.String()
	case ARRAY:
		return "ARRAY: " + entry.DataType.String()
	}
	return "Invalid!"
}

var dataTypes map[string]TypeEntry

func init() {
	dataTypes = map[string]TypeEntry{
		"int":   TypeEntry{PRIMITIVE, reflect.TypeOf(int(0))},
		"int8":  TypeEntry{PRIMITIVE, reflect.TypeOf(int8(0))},
		"int16": TypeEntry{PRIMITIVE, reflect.TypeOf(int16(0))},

		"uint": TypeEntry{PRIMITIVE, reflect.TypeOf(uint(0))},

		"float32": TypeEntry{PRIMITIVE, reflect.TypeOf(int32(0))},
		"string":  TypeEntry{PRIMITIVE, reflect.TypeOf(string(""))},
	}
}

// Register a structure by its name
func Register(base interface{}) {

	// Allocate on the first call
	if dataTypes == nil {
		dataTypes = make(map[string]TypeEntry)
	}
	if IsStructure(base) {
		dataTypes[reflect.TypeOf(base).String()] = TypeEntry{STRUCT, reflect.TypeOf(base)}
	} else {
		dataTypes[reflect.TypeOf(base).String()] = TypeEntry{PRIMITIVE, reflect.TypeOf(base)}
	}
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

func ParseType(name string, count int) TypeEntry {
	automata := regexp.MustCompile(`(.*)\[(.*)\](.*)`)
	groups := automata.FindStringSubmatch(name)

	if groups == nil || len(groups) != 4 {
		log.Dump("Invalid Substring Match for "+name, groups)
		return TypeEntry{UNKNOWN, nil}
	}

	if groups[1] == "map" {
		log.Dump("Allocating a Map", groups)

		keyTypeName := groups[2]
		valueTypeName := groups[3]
		keyType := GetTypeEntry(keyTypeName, 1)
		valueType := GetTypeEntry(valueTypeName, 1)
		finalType := reflect.MapOf(keyType.DataType, valueType.DataType)
		return TypeEntry{
			Category: MAP,
			DataType: finalType,
		}

	} else if groups[1] == "" {
		log.Dump("Allocating a array", groups)
		size := GetInt(groups[2], 0)
		arrayTypeName := groups[3]
		arrayType := GetTypeEntry(arrayTypeName, size)
		finalType := reflect.ArrayOf(count, arrayType.DataType)
		return TypeEntry{
			Category: ARRAY,
			DataType: finalType,
		}

	} else {
		log.Dump("Allocating a slice", groups)
		size := GetInt(groups[2], 0)
		sliceTypeName := groups[3]
		sliceType := GetTypeEntry(sliceTypeName, size)
		finalType := reflect.SliceOf(sliceType.DataType)
		return TypeEntry{
			Category: SLICE,
			DataType: finalType,
		}
	}
	log.Debug("Don't know what this really is?")

	return TypeEntry{UNKNOWN, nil}
}

func GetTypeEntry(name string, size int) TypeEntry {
	name = strings.TrimPrefix(name, "*")

	// Static data types
	log.Debug("Searching Statically for a match", "name", name)
	typeEntry, ok := dataTypes[name]
	if !ok {
		log.Debug("Searching Dynamically for a match", "name", name)
		// dynamic data types -- like maps and slices
		entry := ParseType(name, size)
		return entry
	}

	log.Dump("Found Match", typeEntry)
	return typeEntry
}
