package core

import (
	 "bytes"
   "encoding/hex"
	 "errors"
	 "fmt"
	 "os"


  "../utils"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 31/Mar/2018 Is The Bitcoin Price Finally Reaching Its Bottom?"

type Blockchain struct {
  tip []byte
  databaseObject *utils.DatabaseObject
}

func CreateBlockchain(address, nodeId string) *Blockchain {
  dbFile := fmt.Sprintf(dbFile, nodeId)
  requireNotExist(dbFile)

  databaseObject := utils.CreateNewDatabaseObject(dbFile, blocksBucket)

  cbtx := NewCoinbaseTx(address, genesisCoinbaseData)
  genesis := NewGenesisBlock(cbtx)

  databaseObject.Set(genesis.Hash, genesis.Serialize())
  databaseObject.Set([]byte("l"), genesis.Hash)

	bc := Blockchain{genesis.Hash, &databaseObject}
	return &bc
}

func NewBlockchain(nodeID string) *Blockchain {
  dbFile := fmt.Sprintf(dbFile, nodeID)
	requireExist(dbFile)

  databaseObject := utils.CreateNewDatabaseObject(dbFile, blocksBucket)
  tip := databaseObject.Get([]byte("l"))
	bc := Blockchain{tip, &databaseObject}

	return &bc
}

func (bc *Blockchain) AddBlock(block *Block){
  blockBytes := bc.databaseObject.Get(block.Hash)
  if blockBytes!= nil {
    return //the block exists
  }
  bc.databaseObject.Set(block.Hash, block.Serialize())
  lastBlockData := bc.databaseObject.GetLastBlockData()
  lastBlock := DeserializeBlock(lastBlockData)
  if block.Height > lastBlock.Height {
    bc.databaseObject.SetLastHash(block.Hash)
    bc.tip = block.Hash
  }
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Tx, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.Id, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Tx{}, errors.New("Transaction is not found")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.Id)

		Outputs:
			for outIdx, out := range tx.Outputs {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.Id)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutputIndex)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// Iterator returns a BlockchainIterat
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.databaseObject.Db}

	return bci
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func requireExist(dbFile string){
  if dbExists(dbFile) == false {
    fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
  }
}

func requireNotExist(dbFile string){
  if dbExists(dbFile) == true {
    fmt.Println("Blockchain already exists.")
		os.Exit(1)
  }
}
