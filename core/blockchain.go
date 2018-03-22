package core

import(
  "boltdb/bolt"
)

type Blockchain struct {
  tip []byte
  db *bold.DB
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
    retur nil
  })
}


func (bc *Blockchain) Iterator() *BlockchainIterator {
  return &BlockchainIterator{bc.tip, bc.db}
}

func NewGenesisBlock(coinbase *Transaction) *Block {
  return NewBlock([]*Transaction{coinbase}, []byte{})
}

func NewBlockchain(address string) *Blockchain {
  var tip []byte
  db, err := bolt.Open(dbFile, 0600, nil)

  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blocksBucket))
    if b == nil {
      coinbaseTx := CoinbaseTx(address, "")
      genesis := NewGenesisBlock(coinbaseTx)
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
