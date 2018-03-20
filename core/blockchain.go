package core

import(
  "github.com/boltdb/bolt"
)

type Blockchain struct {
  tip []byte
  db *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
  var lastHash []byte
  err := bc.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blockBucket))
    lastHash = b.Get([]byte("l"))
    return nil
  })
  newBlock := NewBlock(data, lastHash)
  err = bc.db.Update(func(tx *bolt.Tx) error{
    b := tx.Bucket([]byte(blocksBucket))
    err := b.Put(newBlock.Hash, newBlock.Serialize())
    err = b.Put([]byte("l"), newBlock.Hash)
    bc.tip = newBlock.Hash
    return nil
  })
}


func (bc *Blockchain) Iterator() *BlockchainIterator {
  return &BlockchainIterator{bc.tip, bc.db}
}

func NewGenesisBlock() *Block {
  return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
  var tip []byte
  db, err := bolt.Open(dbFile, 0600, nil)

  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    if b == nil {
      genesis := NewGenesisBlock()
      b, err := tx.CreateBucket([]byte(blocksBucket))
      err = b.Put(genesis.Hash, genesis.Serialize())
      err = b.Put([]byte("l"), genesis.Hash)//l will store the hash of last block
    }else{
      tip = b.Get([]byte("l"))
    }
    return nil
  })
  bc := Blockchain{tip, db}
  return &bc
}
