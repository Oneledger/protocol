/*

 */

package btc

import (
	"fmt"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
)

func (s *Service) PrepareLock(hash string, index uint32, net *chaincfg.Params) ([]byte, error) {
	cd := bitcoin.NewChainDriver(s.trackerStore)

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

	utxoInit := bitcoin.NewUTXO(&chainhash.Hash{}, 0, 0)
	txn := cd.PrepareLock(utxoInit, source)

	return txn, nil
}

func (s *Service) AddUserSignatureAndBroadcast(txn []byte, signature []byte) error {
	cd := bitcoin.NewChainDriver(s.trackerStore)

	newTx := cd.AddUserLockSignature(txn, signature)

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18831",
		User:         "oltest01",
		Pass:         "olpass01",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	clt, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	hash, err := cd.BroadcastTx(newTx, clt)
	if err != nil {
		return err
	}

}
