/*
	Copyright 2017-2018 OneLedger
*/
package serial

import (
	"reflect"
	"regexp"
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
	dataTypes[reflect.TypeOf(base).String()] = TypeEntry{STRUCT, reflect.TypeOf(base)}
}

func ParseType(name string, count int) TypeEntry {
	if strings.HasPrefix(name, "map[") {
		automata := regexp.MustCompile(regexp.QuoteMeta("map[(.*)](.*)"))
		groups := automata.FindStringSubmatch(name)
		keyTypeName := groups[1]
		valueTypeName := groups[2]
		keyType := GetTypeEntry(keyTypeName)
		valueType := GetTypeEntry(valueTypeName)
		finalType := reflect.MapOf(keyType.DataType, valueType.DataType)
		return TypeEntry{
			Category: MAP,
			DataType: finalType,
		}

	} else if strings.HasPrefix(name, "[]") {
		automata := regexp.MustCompile(regexp.QuoteMeta("[](.*)"))
		groups := automata.FindStringSubmatch(name)
		arrayTypeName := groups[1]
		arrayType := GetTypeEntry(arrayTypeName)
		finalType := reflect.ArrayOf(count, arrayType.DataType)
		return TypeEntry{
			Category: ARRAY,
			DataType: finalType,
		}
	}
	return TypeEntry{UNKNOWN, nil}
}

func GetTypeEntry(name string) TypeEntry {
	name = strings.TrimPrefix(name, "*")

	typeEntry, ok := dataTypes[name]
	if !ok {
		log.Dump("structures", dataTypes)
		log.Fatal("Missing structure type", "name", name)
		return TypeEntry{UNKNOWN, nil}
	}

	return typeEntry
}
