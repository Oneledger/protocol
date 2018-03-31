package core

import (
	// "bytes"
	// "crypto/ecdsa"
	// "encoding/hex"
	// "errors"
	 "fmt"
	 "log"
	 "os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 31/Mar/2018 Is The Bitcoin Price Finally Reaching Its Bottom?"

type Blockchain struct {
  tip []byte
  db *bolt.DB
}

func CreateBlockchain(address, nodeId string) *Blockchain {
  dbFile := fmt.Sprintf(dbFile, nodeId)
  if dbExists(dbFile) {
    fmt.Println("Blockchain already exists.")
    os.Exit(1)
  }
  var tip []byte
  cbtx := NewCoinbaseTx(address, genesisCoinbaseData)
  genesis := NewGenesisBlock(cbtx)
  db, err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    log.Panic(err)
  }
  err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

func NewBlockchain(nodeID string) *Blockchain {
  dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
