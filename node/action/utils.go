package action

import (
	"github.com/Oneledger/protocol/node/err"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec"

	"github.com/Oneledger/protocol/node/log"
)

func _hash(bytes []byte) []byte {

	hasher := ripemd160.New()

	hasher.Write(bytes)

	return hasher.Sum(nil)
}


type BoxLocker struct{
	Signature *btcec.Signature
	PubKey *btcec.PublicKey
}

func (bl *BoxLocker) sign( preImage []byte, nonce []byte, message Message) error {
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), append(preImage, nonce...))

	signature, err := privKey.Sign(_hash(message))

	if err != nil {
		log.Error("BoxLocker sign error")
	}

	bl.Signature = signature
	bl.PubKey = pubkey
	return err
}

func (bl *BoxLocker) verify(key *btcec.PublicKey, message Message) bool {
	return bl.Signature.Verify(_hash(message),key)
}


type SwapBox struct{
	Swap Swap
	LocalLocker BoxLocker
	RemoteLocker BoxLocker
}


func CreateSwapBox(swap Swap, preImage []byte, nonce []byte) *SwapBox {
	return &SwapBox{swap, nil, nil}
}


func (sb *SwapBox) LocalSign(preImage []byte, nonce []byte, message Message) {
	sb.LocalLocker.sign(preImage, nonce, message)
}

func (sb *SwapBox) CounterSign(preImage []byte, nonce []byte, message Message) {
	sb.RemoteLocker.sign(preImage, nonce, message)
}

func (sb *SwapBox) Verify(localKey *btcec.PublicKey , remoteKey *btcec.PublicKey , nonce []byte, message Message) bool{
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

func SharePreImage(remote string, preImage []byte) (err.Code) {
	log.Info("preImage: "+ string(preImage) + " is shared with: "+remote)
	return err.SUCCESS
}

func (sb *SwapBox) CheckRemoteChain (){

}

func (sb *SwapBox) ExecuteSwap() {

}
