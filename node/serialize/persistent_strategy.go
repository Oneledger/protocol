package serialize

import (
	"encoding/base64"
)

type persistentStrategy struct {
	msgpackStrategy
}

// Serialize
func (j *persistentStrategy) Serialize(obj interface{}) ([]byte, error) {

	b, err := j.msgpackStrategy.Serialize(obj)
	if err != nil {
		return []byte{}, err
	}

	ds := base64.StdEncoding.EncodeToString(b)
	return []byte(ds), nil
}

func (j *persistentStrategy) Deserialize(src []byte, dest interface{}) error {

	b, err := base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return err
	}
	return j.msgpackStrategy.Deserialize(b, dest)
}
