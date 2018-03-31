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
  "strings"

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
  return utils.Deserialize(data, tx).(Tx)
}

func (tx *Tx) Hash () []byte {
  txCloned := tx.TrimmedClone()
  txCloned.Id = []byte{}
  hash := sha256.Sum256(txCloned.Serialize())
  return hash[:]
}

func (tx *Tx) Clone() Tx {
  txCloned := Tx{}
  copier.Copy(&tx, &txCloned)
  return txCloned
}

func (tx *Tx) TrimmedClone() Tx {
  var inputs []TxInput
  var outputs []TxOutput
  for _, input := range tx.Inputs {
    inputs = append(inputs, TxInput{input.Id, input.OutputIndex, nil, nil})
  }
  for _, output := range tx.Outputs {
    outputs = append(outputs, TxOutput{output.Value, output.PubKeyHash})
  }
  return Tx{tx.Id,inputs, outputs}
}

func (tx *Tx) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Tx) Tx{
  if tx.IsCoinbase() {
    return Tx{}
  }
  for _, input := range tx.Inputs {
    if prevTxs[hex.EncodeToString(input.Id)].Id == nil {
      log.Panic("ERROR: Previous transaction is not correct")
    }
  }
  txCloned := tx.TrimmedClone()

  for inputId, input := range txCloned.Inputs {
    prevTx := prevTxs[hex.EncodeToString(input.Id)]
    txCloned.Inputs[inputId].Signature = nil
    txCloned.Inputs[inputId].PubKey = prevTx.Outputs[input.OutputIndex].PubKeyHash
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

func (tx *Tx) String() string {
  var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.Id))

	for i, input := range tx.Inputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Id))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.OutputIndex))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}












func placeHolder() {}
