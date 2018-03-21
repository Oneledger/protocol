package core

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
  Input TxInput
  Output TxOutput
}
