package common

import (
	"github.com/Oneledger/protocol/node/action"

	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/err"
)


type SwapBox struct{
	Swap action.Swap
	LocalLocker BoxLocker
	RemoteLocker BoxLocker
}


// act as transaction box locker.
type BoxLocker struct{
	Signature *btcec.Signature
	PubKey *btcec.PublicKey
}

func (bl *BoxLocker) lock( preImage []byte, nonce []byte) error {

	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), preImage)

	signature, err := privKey.Sign(nonce)

	bl.Signature = signature
	bl.PubKey = pubkey
	return err
}

func (bl *BoxLocker) verify(key *btcec.PublicKey, nonce []byte) bool {
	return bl.Signature.Verify(nonce,key)
}



func CreateSwapBox(swap action.Swap, preImage []byte, nonce []byte) *SwapBox {
	locker := BoxLocker{}
	err:= locker.lock(preImage, nonce)

	if err != nil {
		log.Error("Failed to create SwapBox")
	}
	return &SwapBox{swap, locker, nil}
}

func (sb *SwapBox) CounterSign(preImage []byte, nonce []byte) {
	sb.RemoteLocker.lock(preImage, nonce)
}

func (sb *SwapBox) Verify(localKey *btcec.PublicKey , remoteKey *btcec.PublicKey , nonce []byte) bool{
	if !sb.LocalLocker.verify(localKey, nonce) {
		log.Error("LocalLocker verify failed")
		return false
	}
	if !sb.RemoteLocker.verify(remoteKey, nonce) {
		log.Error("RemoteLocker verify failed")
		return false
	}
	log.Info("SwapBox verified successfully")
	return true
}

func (sb *SwapBox) SharePublicKey(remote string) (err.Code) {
	log.Info("public key of signature:"+ string(sb.LocalLocker.Signature.Serialize())+ " shared with"+remote)
	return err.SUCCESS
}

func (sb *SwapBox) CheckRemoteChain (){

}

func (sb *SwapBox) ExecuteSwap() {

}

func _hash(bytes []byte) []byte {

	hasher := ripemd160.New()

	hasher.Write(bytes)

	return hasher.Sum(nil)
}
