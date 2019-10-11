/*

 */

package btc

import (
	"github.com/Oneledger/protocol/chains/bitcoin"
	btc_data "github.com/Oneledger/protocol/data/bitcoin"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
)

func (s *Service) PrepareLock(hash string, index uint32, net *chaincfg.Params) ([]byte, error) {
	cd := bitcoin.NewChainDriver("")

	btc := gobcy.API{"dd53aae66b83431ca57a1f656af8ed69", "btc", "main"}
	tx, err := btc.GetTX(hash, nil)
	if err != nil {
		return nil, err
	}

	if tx.Confirmations > 10 {
		return nil, errors.New("source transaction doesn't have enough confirmations")
	}

	hashh, _ := chainhash.NewHashFromStr(tx.Hash)
	balance := int64(tx.Outputs[index].Value)
	source := btc_data.NewUTXO(hashh, index, balance)

	utxoInit := btc_data.NewUTXO(&chainhash.Hash{}, 0, 0)

	tracker, err := s.trackerStore.GetTrackerForLock()
	if err != nil {
		return nil, errors.Wrap(err, "error getting tracker for lock")
	}
	txn := cd.PrepareLock(utxoInit, source, tracker.NextLockScriptAddress)

	return txn, nil
}

func (s *Service) AddUserSignatureAndBroadcast(txn []byte, signature []byte) error {
	cd := bitcoin.NewChainDriver("")

	newTx := cd.AddUserLockSignature(txn, signature)

}
