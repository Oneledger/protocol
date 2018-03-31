package core
import (
    "time"
    "../utils"

    "bytes"
    "encoding/gob"
    "log"
)

type Block struct {
  Timestamp     int64
  Transactions  []*Tx
  PrevBlockHash []byte
  Hash          []byte
  Nonce         int64
  Height        int
}

func NewBlock (transactions []*Tx, prevBlockHash []byte, height int) *Block {
  block := Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
  pow := NewProofOfWork(&block) //TODO: Yeah, we are going to remove this
  nonce, hash := pow.Run()

  block.Hash = hash[:]
  block.Nonce = nonce
  return &block
}

func NewGenesisBlock(coinbase *Tx) *Block {
  return NewBlock([]*Tx{coinbase}, []byte{}, 0)
}

func (block *Block) HashTransactions () []byte{
  var transactions [][]byte
  for _, tx := range block.Transactions {
    transactions = append(transactions, tx.Serialize())
  }
  mTree := utils.NewMerkleTree(transactions)
  return mTree.RootNode.Data
}

func (block *Block) Serialize() []byte {
  var result bytes.Buffer
  encoder := gob.NewEncoder(&result)
  err := encoder.Encode(block)
  if err != nil {
    log.Panic(err)
  }
  return result.Bytes()
}

func DeserializeBlock(d []byte) * Block {
  var block Block
  decoder := gob.NewDecoder(bytes.NewReader(d))
  err := decoder.Decode(&block)
  if err != nil {
    log.Panic(err)
  }
  return &block
}
