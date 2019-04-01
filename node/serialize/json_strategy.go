package serialize

import (
	"encoding/json"
)

type jsonStrategy struct {
	metaStrategy
}

// Serialize
func (j *jsonStrategy) Serialize(obj interface{}) ([]byte, error) {

	//check if data adapter
	if apr, ok := obj.(DataAdapter); ok {
		obj = apr.Data()
	}

	return json.Marshal(obj)
}

func (j *jsonStrategy) Deserialize(src []byte, dest interface{}) error {

	// check if object satisfies adapter interface and likewise
	// call the decorator while handling it the plain deserialize method
	if apr, ok := dest.(DataAdapter); ok {
		err := j.wrapDataAdapter(src, dest, j.deserialize, apr)
		return err
	}

	return j.deserialize(src, dest)
}

func (jsonStrategy) deserialize(src []byte, dest interface{}) error {
	return json.Unmarshal(src, dest)
}
