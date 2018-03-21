package core

import "fmt"

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
  Id []byte
  Input []TxInput
  Output []TxOutput
}

const reward = 1024

func CoinbaseTx(to string, data string) *Transaction {
  if data == "" {
    data = "This is the beginning of Oneledger network"
  }
  txin := TxInput{[]byte{}, 0, data}
  txout := TxOutput{reward, to}
  tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
  return &tx
}
