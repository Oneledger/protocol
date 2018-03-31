package core

import (
  "../utils"
  "crypto/sha256"
  "crypto/ecdsa"
  "encoding/hex"
  "crypto/rand"
  "crypto/elliptic"
  "math/big"

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

type TxPreSignature struct {
  TxId          []byte
  TxOutputIndex int
  PubKey        []byte
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
  txCloned := tx.Clone()
  txCloned.Id = []byte{}
  hash := sha256.Sum256(txCloned.Serialize())
  return hash[:]
}

func (tx *Tx) Clone() Tx {
  txCloned := Tx{}
  copier.Copy(&tx, &txCloned)
  return txCloned
}

func NewPreSignature(currentTx *Tx, prevTxMap map[string]Tx, currentInput *TxInput) TxPreSignature {
  prevTx := prevTxMap[hex.EncodeToString(currentInput.Id)]
  return TxPreSignature{currentTx.Id, currentInput.OutputIndex, prevTx.Outputs[currentInput.OutputIndex].PubKeyHash}
}

func GetSignedData(privKey ecdsa.PrivateKey,dataToSign *TxPreSignature) []byte{
  data := fmt.Sprintf("%x\n", dataToSign)
  r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(data))
  if err != nil {
    log.Panic(err)
  }
  signature := append(r.Bytes(), s.Bytes()...)
  return signature
}

func validatePrevTransactions(tx *Tx, prevTxMap map[string]Tx) bool {
  for _, input := range tx.Inputs {
    if prevTxMap[hex.EncodeToString(input.Id)].Id == nil {
      return false
    }
  }
  return true;
}

func (tx *Tx) Sign(privKey ecdsa.PrivateKey, prevTxMap map[string]Tx){
  if tx.IsCoinbase() {
    return
  }
  if validatePrevTransactions(tx, prevTxMap) {
    log.Panic("ERROR: Previous transaction is not correct")
  }

  for inputId, input := range tx.Inputs {
    dataToSign := NewPreSignature(tx, prevTxMap, &input)
    signature := GetSignedData(privKey, &dataToSign)
    tx.Inputs[inputId].Signature = signature
  }
}

func extractBytes(data []byte) (big.Int, big.Int) {
  a := big.Int{}
  b := big.Int{}
  len := len(data)
  a.SetBytes(data[:(len/2)])
  b.SetBytes(data[(len/2):])
  return a, b
}

func verifiedInputTx(input *TxInput,dataToSign *TxPreSignature) bool{
  r, s := extractBytes(input.Signature)
  x, y := extractBytes(input.PubKey)
  data := fmt.Sprintf("%x\n", dataToSign)
  curve := elliptic.P256()
  rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&rawPubKey, []byte(data), &r, &s)
}

func (tx *Tx) Verify(prevTxMap map[string]Tx) bool{
  if tx.IsCoinbase() {
    return true
  }
  if validatePrevTransactions(tx, prevTxMap) {
    log.Panic("ERROR: Previous transaction is not correct")
  }
  for _, input := range tx.Inputs {
    dataToSign := NewPreSignature(tx, prevTxMap, &input)
    if verifiedInputTx(&input, &dataToSign) == false {
      return false
    }
  }
  return true
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
