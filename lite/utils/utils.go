package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"

	"golang.org/x/crypto/ripemd160"
)

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func ReverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func HashPubKey(pubKey []byte) []byte {
	sha256Sum := sha256.Sum256(pubKey)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Sum[:])
	if err != nil {
		log.Panic(err)
	}
	return ripemd160Hasher.Sum(nil)
}

func Serialize(anything interface{}) []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(anything)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}
func Deserialize(data []byte, anything interface{}) interface{} {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&anything)
	if err != nil {
		log.Panic(err)
	}
	return anything
}
