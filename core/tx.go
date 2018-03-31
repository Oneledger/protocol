package core

import (
  "../utils"
  "crypto/sha256"
  "crypto/ecdsa"
  "encoding/hex"
  "crypto/rand"

  "github.com/jinzhu/copier"
  "log"
  "fmt"

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
  txCloned := tx.clone()
  txCloned.ID = []byte{}
  return sha256.Sum256(txCloned.Serialize())[:]
}

func (tx *Tx) Clone() Tx {
  txCloned := Tx{}
  copier.Copy(&tx, &txCloned)
  return txCloned
}

func (tx *Tx) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Tx) Tx{
  if tx.IsCoinbase() {
    return
  }
  for _, input := range tx.Inputs {
    if prevTxs[hex.EncodeToString(input.Id)].Id == nil {
      log.Panic("ERROR: Previous transaction is not correct")
    }
  }
  txCloned := tx.Clone()

  for inputId, input := range txCloned.Inputs {
    prevTx := prevTxs[hex.EncodeToString(input.Id)]
    txCloned.inputs[inputId].Signature = nil
    txCloned.inputs[inputId].PubKey = prevTx.outpus[input.OutputIndex].PubKeyHash
    dataToSign := fmt.Sprintf("%x\n", txCloned)
    r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
    if err != nil {
      log.Panic(err)
    }
    signature := append(r.Bytes(), s.Bytes()...)
    tx.Inputs[inputId].Signature = signature
    txCloned.Inputs[inputId].PubKey = nil
  }
  return txCloned
}












func placeHolder() {}
