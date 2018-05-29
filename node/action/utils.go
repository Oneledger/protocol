package action

import (
	"github.com/Oneledger/protocol/node/err"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec"

	"github.com/Oneledger/protocol/node/log"
)

//general hash method for the actions messages
func _hash(bytes []byte) []byte {

	hasher := ripemd160.New()

	hasher.Write(bytes)

	return hasher.Sum(nil)
}


// General BoxLocker struct to act as locker for any information exchange box of transactions, if verify valid
// then the lock can be release, otherwise box modified somehow
type BoxLocker struct{
	Signature *btcec.Signature
	PubKey *btcec.PublicKey
}


// Sign the locker with preImage and nonce for message passed, the message should be the full information of Transaction
// The nonce is used to preventing the 3rd party from get the message even through he get the preImage, where nonce
// should only be known by the participants of the message sharing
func (bl *BoxLocker) Sign( preImage []byte, nonce []byte, message Message) error {
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), append(preImage, nonce...))

	signature, err := privKey.Sign(_hash(message))

	if err != nil {
		log.Error("BoxLocker sign error")
	}

	bl.Signature = signature
	bl.PubKey = pubkey
	return err
}

// Verify the pubKey with the signature for the message got.
func (bl *BoxLocker) Verify(message Message) bool {
	return bl.Signature.Verify(_hash(message),bl.PubKey)
}


type SwapBox struct{
	Swap Swap
	LocalLocker BoxLocker
	RemoteLocker BoxLocker
}


func CreateSwapBox(swap Swap) *SwapBox {
	return &SwapBox{swap, BoxLocker{}, BoxLocker{}}
}


func (sb *SwapBox) LocalSign(preImage []byte, nonce []byte, message Message) {
	sb.LocalLocker.Sign(preImage, nonce, message)
}

func (sb *SwapBox) CounterSign(preImage []byte, nonce []byte, message Message) {
	sb.RemoteLocker.Sign(preImage, nonce, message)
}

func (sb *SwapBox) Verify(message Message) bool{
	if !sb.LocalLocker.Verify(message) {
		log.Error("LocalLocker verify failed")
		return false
	}
	if !sb.RemoteLocker.Verify(message) {
		log.Error("RemoteLocker verify failed")
		return false
	}
	log.Info("SwapBox verified successfully")
	return true
}

func SharePreImage(remote string, preImage []byte) (err.Code) {
	log.Info("preImage: "+ string(preImage) + " is shared with: "+remote)
	return err.SUCCESS
}

func (sb *SwapBox) CheckRemoteChain (){

}

func (sb *SwapBox) ExecuteSwap() {

}
