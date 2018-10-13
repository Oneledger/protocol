package serial

import (
	"reflect"

	"github.com/Oneledger/protocol/node/log"
)

// ConvertMap takes a structure and return a map of its elements
func ConvertMap(container interface{}) (map[string]interface{}, int) {
	var result map[string]interface{}

	valueOf := reflect.ValueOf(container)
	log.Dump("ConvertMap", container, valueOf.Kind().String())

	kind := valueOf.Kind()
	switch kind {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		if valueOf.IsNil() {
			return nil, -1
		}
	}

	children := GetChildren(container)

	result = make(map[string]interface{}, len(children))

	for _, child := range children {
		result[child.Name] = child.Value
	}
	return result, len(children)
}
