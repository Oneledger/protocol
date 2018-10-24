package action

import (
	"encoding/hex"
	"testing"

	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/assert"
)

func XTestBoxLocker_Sign(t *testing.T) {
	pkBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2d4f87" +
		"20ee63e502ee2869afab7de234b80c")
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))

	locker := BoxLocker{}
	status := locker.Sign(pkBytes[:len(pkBytes)/2], pkBytes[len(pkBytes)/2:], messageHash)

	assert.Equal(t, nil, status, "Sign with preimage and nonce success")
	if status != nil {
		log.Error(status.Error())
		return
	}

	// Serialize and display the signature.
	log.Info("Serialized Signature: %x\n", locker.Signature.Serialize())

	// Verify the signature for the message using the public key.
	verified := locker.Signature.Verify(messageHash, locker.PubKey)
	assert.Equal(t, true, verified, "Signature verified with pubKey")

	// Test sign again with wrong nonce
	err = locker.Sign(pkBytes[:len(pkBytes)/2], pkBytes[len(pkBytes)/2-1:], messageHash)
	assert.Equal(t, nil, err)

	// Test sign again with wrong preImage
	err = locker.Sign(pkBytes[:len(pkBytes)/2-1], pkBytes[len(pkBytes)/2:], messageHash)
	assert.Equal(t, nil, err)

}

func XTestBoxLocker_Verify(t *testing.T) {
	// Decode hex-encoded serialized public key.
	pubKeyBytes, err := hex.DecodeString("02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
	if err != nil {
		log.Error(err.Error())
		return
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Decode hex-encoded serialized signature.
	sigBytes, err := hex.DecodeString("30450220090ebfb3690a0ff115bb1b38b" +
		"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
		"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479")

	if err != nil {
		log.Error(err.Error())
		return
	}

	signature, err := btcec.ParseSignature(sigBytes, btcec.S256())
	if err != nil {
		log.Error(err.Error())
		return
	}

	locker := BoxLocker{signature, pubKey}

	// Verify the signature for the message using the public key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := locker.Verify(messageHash)
	assert.Equal(t, true, verified, "Signature verified")
}

func TestReflection(t *testing.T) {

}
