package serialize

import (
	"bytes"
	"github.com/vmihailenco/msgpack"
)

type msgpackStrategy struct {
	metaStrategy
}

func (m *msgpackStrategy) Serialize(obj interface{}) ([]byte, error) {

	b := m.serializeString(obj)
	if len(b) > 0 {
		return b, nil
	}

	if apr, ok := obj.(DataAdapter); ok {
		obj = apr.Data()
	}
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).SortMapKeys(true)
	err := enc.Encode(obj)
	return buf.Bytes(), err
}

func (m *msgpackStrategy) Deserialize(src []byte, dest interface{}) error {
	if dest == nil {
		return ErrIncorrectWrapper
	}

	if apr, ok := dest.(DataAdapter); ok {
		err := m.wrapDataAdapter(src, dest, m.deserialize, apr)
		return err
	}

	return m.deserialize(src, dest)
}

func (msgpackStrategy) deserialize(src []byte, dest interface{}) error {
	return msgpack.Unmarshal(src, dest)
}
