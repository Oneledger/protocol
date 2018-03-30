package core

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
