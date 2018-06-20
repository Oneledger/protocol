package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/global"
	"time"
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/Oneledger/protocol/node/chains/ethereum/htlc"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

var client *ethclient.Client
var contract *htlc.Htlc

func GetEthClient() (*ethclient.Client) {
	if client == nil {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(global.Current.ETHAddress)
			if err != nil{
				log.Fatal("failed to get geth ipc ", "err", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	} else if id, _ := client.NetworkID(nil); id == big.NewInt(20180229) {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(global.Current.ETHAddress)
			if err != nil{
				log.Fatal("failed to get geth ipc ", "err", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	}
	return client

}

func GethAuth() (*bind.TransactOpts) {
	//todo: generate auth when register without pre-allocate
	nodeName := global.Current.NodeName
	switch nodeName {

	case "Alice-Node":
		key := `{"address":"d7858005867c3449f6673a91f6e4f719f10e12e5","crypto":{"cipher":"aes-128-ctr","ciphertext":"72f2b4ce4af1ff68f3b8c3d0c5b6b3c55e153441633df9ba5b615134f697c108","cipherparams":{"iv":"7839596b4f84b7b2c4cdf01e841a8cec"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"00750dcae280f52db4950517334e02348488e638f2958dffb625ddead8d626f3"},"mac":"341af04c7105565890aac985d13e46b3362312fc4d488944d17d19f03c100386"},"id":"00888165-5663-4e71-beea-1b6055f91cbe","version":3}`
		auth, err := bind.NewTransactor(strings.NewReader(key), "1234")
		if err != nil {
			log.Fatal("Can't get pre-allocate auth for Alice", "err", err)
		}
		return auth
	case "Bob-Node":
		key := `{"address":"aafa2d8980a730b02195f9c8dfeafeb3e69a69ca","crypto":{"cipher":"aes-128-ctr","ciphertext":"cfc36b7deb503116482371b7d2596aa936758b8247279efce461cf0344ae4b31","cipherparams":{"iv":"fc200b937116258856dd0e5a085e011d"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7354c4523dfc70372c8c34616c15dc21448ac40617ffc3a7b3a9af7ee32c37e6"},"mac":"1b322fa3c5789cede87144783f2bd8c4588e5094e68f7880640ca9a5458b8aab"},"id":"87363a39-0171-4640-ba12-b5aacad7aed2","version":3}`
		auth, err := bind.NewTransactor(strings.NewReader(key), "2345")
		if err != nil {
			log.Fatal("Can't get pre-allocate auth for Bob", "err", err)
		}
		return auth
	case "Carol-Node":
		key := `{"address":"8a309f95de0e47edb61de8fa0cf8bdd722271789","crypto":{"cipher":"aes-128-ctr","ciphertext":"81becb7ca37be737af147aa0552b1639b770d76ba98fa82069325fe1ce6e1aa1","cipherparams":{"iv":"5be20f263a46d6cca53cb0ae490245fd"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"12456c9a74778a06449596676cc90f2f046e306b5db74688600c04577529b9c2"},"mac":"6737c9dd93f0abc8e102590984790214c2b9dfc36ea6e2b769e80c19eb22e4e8"},"id":"fbaef12b-a667-4c4e-b4c7-7234ef37cbe9","version":3}`
		auth, err := bind.NewTransactor(strings.NewReader(key),"3456")
		if err != nil {
			log.Fatal( "Can't get pre-allocate auth for Carol", "err", err)
		}
		return auth
	default:
		log.Info("This node don't have pre-allocate Eth account")
		return nil
	}
}

