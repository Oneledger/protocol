package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	actionGov "github.com/Oneledger/protocol/action/governance"
	"github.com/Oneledger/protocol/app"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	db "github.com/tendermint/tm-db"
	"strconv"
	"testing"
)

const (
	numOfProposals = 5
)

var (
	txList           []*app.RequestDeliverTx
	transactionStore *TransactionStore
)

func init() {
	for i := 0; i < numOfProposals; i++ {
		id := governance.ProposalID("proposal" + strconv.Itoa(i))
		tx := &actionGov.ExpireVotes{
			ProposalID:       id,
			ValidatorAddress: keys.Address("address1"),
		}
		txData, _ := serialize.GetSerializer(serialize.PERSISTENT).Serialize(tx)
		deliverTx := &app.RequestDeliverTx{
			Tx: txData,
		}
		txList = append(txList, deliverTx)
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	transactionStore = NewTransactionStore("tran", cs)
}

func TestTransactionStore_Set(t *testing.T) {
	for _, tx := range txList {
		hash := sha256.New()
		hash.Write(tx.Tx)
		hashData := hash.Sum(nil)
		hashStr := hex.EncodeToString(hashData)
		_ = transactionStore.Set(tx, hashStr)

		newTx, _ := transactionStore.Get(hashStr)
		result := bytes.Compare(tx.Tx, newTx.Tx)

		assert.Equal(t, result, 0)
	}
}

func TestTransactionStore_Delete(t *testing.T) {
	hash := sha256.New()
	hash.Write(txList[0].Tx)
	hashData := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashData)

	res, _ := transactionStore.Delete(hashStr)
	assert.Equal(t, res, true)

	transactionStore.State.Commit()
	res = transactionStore.Exists(hashStr)
	assert.Equal(t, res, false)
}
