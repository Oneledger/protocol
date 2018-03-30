package core
import (
    "bytes"
    "crypto/sha256"
    "strconv"
    "time"
    "encoding/gob"
    "log"
)

type Block struct {
  Timestamp     int64
  Transactions  []*Tx
  PrevBlockHash []byte
  Hash          []byte
  Nonce         int
  Height        int
}

func NewBlock (transactions []*Tx, prevBlockHash []byte, height int) *Block {
  block := Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
  pow := NewProofOfWork(block) //TODO: Yeah, we are going to remove this
  nonce, hash := pow.Run()

  block.Hash = hash[:]
  block.Nonce = nonce
  return &block
}

func NewGenesisBlock(coinbase *Tx) *Block {
  return NewBlock([]*Transactions{coinbase}, []byte{}, 0)
}
