/*

 */

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

type KeyChain struct {
	master *hdkeychain.ExtendedKey
}

func NewKeyChain() *KeyChain {
	// Generate a random seed at the recommended length.
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Generate a new master node using the seed.
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &KeyChain{key}
}

func LoadKeyChainFromFile(filename string) (*KeyChain, error) {
	f, err := os.Open("filename")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ek, err := hdkeychain.NewKeyFromString(string(b))
	if err != nil {
		return nil, err
	}

	return &KeyChain{ek}, nil
}

func (kc *KeyChain) SaveToFile(filename string) error {

	f, err := os.OpenFile(filename, 0, os.ModeExclusive)
	if err != nil {
		return err
	}
	enc := hex.NewEncoder(f)
	_, err = enc.Write([]byte(kc.master.String()))
	return err
}

// m / 44' / 0' / 0' / 0 / 0
func (kc *KeyChain) GetBitcoinExternalAddress(i uint32, net *chaincfg.Params) (btcutil.Address, error) {
	k, err := kc.master.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(i)
	if err != nil {
		return nil, err
	}

	privKey, err := k.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return btcutil.NewAddressPubKey(privKey.PubKey().SerializeCompressed(), net)
}

// m / 44' / 0' / 0' / 1 / 0
func (kc *KeyChain) GetBitcoinRedeemAddress(i uint32, net *chaincfg.Params) (btcutil.Address, error) {
	k, err := kc.master.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(1)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(i)
	if err != nil {
		return nil, err
	}

	privKey, err := k.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return btcutil.NewAddressPubKey(privKey.PubKey().SerializeCompressed(), net)
}

// m / 44' / 0' / 0' / 0 / 0
func (kc *KeyChain) GetBitcoinWIF(i uint32, net *chaincfg.Params) (*btcutil.WIF, error) {
	k, err := kc.master.Child(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(0)
	if err != nil {
		return nil, err
	}

	k, err = k.Child(i)
	if err != nil {
		return nil, err
	}

	privKey, err := k.ECPrivKey()
	if err != nil {
		return nil, err
	}

	return btcutil.NewWIF(privKey, net, true)
}
