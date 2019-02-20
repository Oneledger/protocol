package action

import (
	"encoding/hex"
	"errors"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"time"

	"github.com/Oneledger/protocol/node/serial"
	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

//general hash method
func _hashToStringBytes(item interface{}) []byte {
	sum := _hash(item)
	sumStr := hex.EncodeToString(sum[:])
	return []byte(sumStr)
}

func _hash(item interface{}) []byte {

	hasher := ripemd160.New()

	buffer, err := serial.Serialize(item, serial.JSON)
	if err != nil {
		log.Fatal("hash serialize failed", "err", err)
	}
	_, err = hasher.Write(buffer)
	if err != nil {
		log.Fatal("hasher failed", "err", err)
	}
	return hasher.Sum(nil)
}

// General BoxLocker struct to act as locker for any information exchange box of transactions, if verify valid
// then the lock can be release, otherwise box modified somehow
type BoxLocker struct {
	Signature *btcec.Signature
	PubKey    *btcec.PublicKey
}

// Sign the locker with preImage and nonce for message passed, the message should be the full information of
// Transaction. The nonce is used to preventing the 3rd party from get the message even through he get the preImage,
// where nonce should only be known by the participants of the message sharing
func (bl *BoxLocker) Sign(preImage []byte, nonce []byte, message Message) error {
	privKey, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), append(preImage, nonce...))

	signature, err := privKey.Sign(_hash(message))

	if err != nil {
		return errors.New("BoxLocker sign error")
	}
	if bl.Signature != nil && bl.PubKey != nil {
		bl.Signature = signature
		bl.PubKey = pubkey
		return nil
	} else if bl.Signature == signature && bl.PubKey == pubkey {
		return nil
	}
	return errors.New("sign error for preImage and nonce not match")
}

// Verify the pubKey with the signature for the message got.
func (bl *BoxLocker) Verify(message Message) bool {
	return bl.Signature.Verify(_hash(message), bl.PubKey)
}

type MultiSigBox struct {
	Lockers []BoxLocker `json:"lockers"`
	M       int         `json:"m"`
	N       int         `json:"n"`
	message Message     `json:"message"`
}

func NewMultiSigBox(m int, n int, message Message) *MultiSigBox {
	return &MultiSigBox{
		message: message,
		M:       m,
		N:       n,
	}
}

func (msb MultiSigBox) Sign(locker *BoxLocker) error {
	if !locker.Verify(msb.message) {
		return errors.New("signature not match message")
	}

	if len(msb.Lockers) == msb.N {
		return errors.New("box is already locked")
	}

	for _, box := range msb.Lockers {
		if box.Signature.IsEqual(locker.Signature) {
			return errors.New("signature already included")
		}
	}

	msb.Lockers = append(msb.Lockers, *locker)
	return nil
}

func (msb MultiSigBox) Unlock(signs []BoxLocker) *Message {
	cnt := 0
	for _, locker := range msb.Lockers {
		if locker.Verify(msb.message) {
			cnt++
		}
	}
	if cnt >= msb.M {
		return &msb.message
	}
	return nil
}

type HtlContract struct {
	Transaction     Transaction `json:"transaction"`
	Secrethash      [32]byte    `json:"secrethash"`
	LockHeight      int64       `json:"lockheight"`
	StartFromHeight int64       `json:"startfromheight"`
}

func init() {
	serial.Register(HtlContract{})
}

func CreateHtlContract(app interface{}, transaction Transaction, scrhash [32]byte, lock time.Duration) *HtlContract {
	height := GetHeight(app)
	lockHeight := int64(lock.Seconds())

	return &HtlContract{
		Transaction:     transaction,
		Secrethash:      scrhash,
		LockHeight:      lockHeight,
		StartFromHeight: height,
	}
}

func GetNodeAccount(app interface{}) id.Account {

	accounts := GetAccounts(app)
	account, _ := accounts.FindName(global.Current.NodeAccountName)
	if account == nil {
		log.Error("Node does not have account", "name", global.Current.NodeAccountName)
		return nil
	}

	return account
}
