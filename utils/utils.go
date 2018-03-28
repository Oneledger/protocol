package utils

import (
  "bytes"
  "encoding/binary"
  "log"
  "crypto/sha256"

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

func ReverseBytes(data []byte) []byte{
  for i, j := 0, len(data) - 1; i < j ; i, j = i+1, j-1 {
    data[i], data[j] =  data[j], data[i]
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
  return ripemd160Hasher.Sum(nil);
}
