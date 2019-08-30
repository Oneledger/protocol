/*

 */

package btc

import (
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
)

func (s *Service) PrepareLock(hash string, index uint32) ([]byte, error) {
	cd := bitcoin.NewChainDriver()

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
	source := bitcoin.NewUTXO(hashh, index, balance)

	lockScript, lockScriptAddress, err := bitcoin.CreateMultiSigAddress(1, validatorPubKeys)

	utxoInit := bitcoin.NewUTXO(&chainhash.Hash{}, 0, 0)
	txn := cd.PrepareLock(utxoInit, source, nil, lockScriptAddress)

	return txn, nil
}

func (s *Service) AddUserSignatureAndBroadcast(txn []byte, signature []byte) {
	cd := bitcoin.NewChainDriver()

	newTx := cd.AddUserLockSignature(txn, signature)

}
