package serialize

import (
	"github.com/vmihailenco/msgpack"
)

type msgpackStrategy struct {
	metaStrategy
}

func (m *msgpackStrategy) Serialize(obj interface{}) ([]byte, error) {

	if apr, ok := obj.(DataAdapter); ok {
		obj = apr.Data()
	}
	return msgpack.Marshal(obj)
}

func (m *msgpackStrategy) Deserialize(src []byte, dest interface{}) error {

	if apr, ok := dest.(DataAdapter); ok {
		err := m.wrapDataAdapter(src, dest, m.deserialize, apr)
		return err
	}

	return m.deserialize(src, dest)
}

func (msgpackStrategy) deserialize(src []byte, dest interface{}) error {
	return msgpack.Unmarshal(src, dest)
}
