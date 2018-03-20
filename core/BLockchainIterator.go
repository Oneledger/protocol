package core

import (
  "github.com/boltdb/bolt"
)

type BlockchainIterator struct {
  currentHash []byte
  db *bolt.DB
}


func (blockchainIterator *BlockchainIterator) prev() *Block {
  var block *Block
  err := blockchainIterator.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    encodedBlock := b.Get(blockchainIterator.currentHash)
    block = DeserializeBlock(encodedBlock)
    return nil
  })
  blockchainIterator.currentHash = block.PrevBlockHash
  return block
}
