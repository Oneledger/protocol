package serialize

import (
	"bytes"
	"encoding/gob"
)

type gobStrategy struct {
	metaStrategy
}

// Serialize
func (j *gobStrategy) Serialize(obj interface{}) ([]byte, error) {

	//check if data adapter
	if apr, ok := obj.(DataAdapter); ok {
		obj = apr.Data()
	}

	var bb bytes.Buffer
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(obj)
	return bb.Bytes(), err
}

func (j *gobStrategy) Deserialize(src []byte, dest interface{}) error {

	// check if object satisfies adapter interface and likewise
	// call the decorator while handling it the plain deserialize method
	if apr, ok := dest.(DataAdapter); ok {
		err := j.wrapDataAdapter(src, dest, j.deserialize, apr)
		return err
	}

	return j.deserialize(src, dest)
}

func (gobStrategy) deserialize(src []byte, dest interface{}) error {

	bb := bytes.NewReader(src)
	dec := gob.NewDecoder(bb)
	return dec.Decode(dest)
}
