/*

 */

package internal

import "github.com/btcsuite/btcd/btcec"

type signBTCTxnData struct {
	txnData          []byte
	validatorPrivKey btcec.PrivateKey
}

func signBTCTxn(data interface{}) {

}
