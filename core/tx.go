package core

import (
  "../utils"
  "crypto/sha256"
)

const subsidy = 10 //TODO: this will be removed

type Tx struct {
  Id        []byte
  Inputs    []TxInput
  Outputs   []TxOutput
}

func (tx *Tx) IsCoinbase() bool {
  return  len(tx.Inputs) == 1 &&
          len(tx.Inputs[0].Id) == 0 &&
          tx.Inputs[0].OutputIndex == -1
}

func (tx *Tx) Serialize() []byte {
  return utils.Serialize(tx)
}

func DeserializeTx(data []byte) Tx {
  var tx Tx //in order to tell the deserialize the type
  return utils.Deserialize(data, tx)
}

func (tx *Tx) Hash () []byte {
  return sha256.Sum256(tx.Serialize())[:]
}
