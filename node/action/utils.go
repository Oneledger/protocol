package action

import (
	"github.com/Oneledger/protocol/node/err"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/id"
	"github.com/go-errors/errors"
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


// Sign the locker with preImage and nonce for message passed, the message should be the full information of
// Transaction. The nonce is used to preventing the 3rd party from get the message even through he get the preImage,
// where nonce should only be known by the participants of the message sharing
func (bl *BoxLocker) Sign( preImage []byte, nonce []byte, message Message) (bool, error) {
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), append(preImage, nonce...))

	signature, err := privKey.Sign(_hash(message))

	if err != nil {
		return false, errors.New("BoxLocker sign error")
	}
	if bl.Signature != nil && bl.PubKey != nil {
		bl.Signature = signature
		bl.PubKey = pubkey
		return true, nil
	} else if bl.Signature == signature && bl.PubKey == pubkey {
		return true, nil
	}
	return false, errors.New("Sign Error for preImage and nonce not match")
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
	//TODO: add the necessary mechanism fo create the swap
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

//
func (sb *SwapBox) unlock(localPreImage []byte, remotePreImage, nonce []byte, message Message) bool {
	a, err1 := sb.LocalLocker.Sign(remotePreImage,nonce,message)
	b, err2 := sb.RemoteLocker.Sign(localPreImage,nonce,message)

	if err1 != nil {
		log.Error("Remote preImage is wrong")
		return false
	}
	if err2 != nil {
		log.Error("Local preImage is wrong")
		return false
	}
	return a&&b
}

func SharePreImage(remoteAddress id.Address, preImage []byte) (err.Code) {
	log.Info("preImage: "+ string(preImage) + " is shared with: "+remoteAddress.String())
	return SubmitToOl(remoteAddress, preImage)
}

//TODO: this function is called to send message to another address in oneledger node, move it to corret place
func SubmitToOl(dest id.Address, message Message) err.Code{
	log.Info(string(message)+"is send.")
	return err.SUCCESS
}

func (sb *SwapBox)CheckRemoteBox (box *SwapBox) bool {
	//TODO: check the transaction initialed on remote node. need to update Swap struct to get the remote information
	return true
}

func (sb *SwapBox) Execute(box *SwapBox) {
	//TODO: call /node/chains/btcrpc/interfaces or /node/chains/btcrpc/interfaces to finalize the HTLC transaction
	if sb.CheckRemoteBox(box) {
		sb.unlock([]byte{},[]byte{},[]byte{},Message{})
	}
}
