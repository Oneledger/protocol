package core

import(
  "github.com/boltdb/bolt"
  "log"
)

const blockBucket = "blockBucket" //it is like a table name
const dbFile = "fancyBlock.db" //Mmmm, isn't it fancy?

type Blockchain struct {
  tip []byte
  db *bolt.DB
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
  var lastHash []byte
  err := bc.db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(blockBucket))
    lastHash = b.Get([]byte("l"))
    return nil
  })
  if err != nil {
    log.Panic(err)
  }
  newBlock := NewBlock(transactions, lastHash)
  err = bc.db.Update(func(tx *bolt.Tx) error{
    b := tx.Bucket([]byte(blockBucket))
    err := b.Put(newBlock.Hash, newBlock.Serialize())
    err = b.Put([]byte("l"), newBlock.Hash)
    if err != nil {
      log.Panic(err)
    }
    bc.tip = newBlock.Hash
    return nil
  })

}

func (bc *Blockchain) CloseDB() {
  bc.db.Close()
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
    b := tx.Bucket([]byte(blockBucket))
    if b == nil {
      coinbaseTx := CoinbaseTx(address, "")
      genesis := NewGenesisBlock(coinbaseTx)
      b, err := tx.CreateBucket([]byte(blockBucket))
      err = b.Put(genesis.Hash, genesis.Serialize())
      err = b.Put([]byte("l"), genesis.Hash)//l will store the hash of last block
      return err
    }else{
      tip = b.Get([]byte("l"))
    }
    return nil
  })

  if err != nil {
    log.Panic(err)
  }

  bc := Blockchain{tip, db}
  return &bc
}
