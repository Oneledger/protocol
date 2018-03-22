package core

import (
  "github.com/boltdb/bolt"
  "log"
)

type BlockchainIterator struct {
  currentHash []byte
  db *bolt.DB
}


func (blockchainIterator *BlockchainIterator) Prev() *Block {
  var block *Block
  err := blockchainIterator.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blockBucket))
    encodedBlock := b.Get(blockchainIterator.currentHash)
    block = DeserializeBlock(encodedBlock)
    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  blockchainIterator.currentHash = block.PrevBlockHash
  return block
}
