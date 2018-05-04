package core

import (
  "bytes"
  "crypto/sha256"
  "encoding/gob"
  "log"
)

type TxInput struct {
  TxId []byte
  Value int
  Signature string
}

type TxOutput struct {
  Value int
  PublicKey string
}

type Transaction struct {
  ID []byte
  Input []TxInput
  Output []TxOutput
}

const reward = 10

func (tx *Transaction) setID() {
  var encoded bytes.Buffer
  var hash [32]byte

  enc := gob.NewEncoder(&encoded)
  err := enc.Encode(tx)
  if err != nil {
    log.Panic(err)
  }
  hash = sha256.Sum256(encoded.Bytes())
  tx.ID = hash[:]
}

func CoinbaseTx(to string, data string) *Transaction {
  if data == "" {
    data = "This is the beginning of Oneledger network"
  }
  txin := TxInput{[]byte{}, 0, data}
  txout := TxOutput{reward, to}
  tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
  tx.setID()
  return &tx
}
