package ethereum

import (
	"math/big"
	"strings"
	"time"

	"github.com/Oneledger/protocol/node/chains/ethereum/htlc"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/Oneledger/protocol/node/comm"
    "os"
    "github.com/pkg/errors"
)

var client *ethclient.Client
var htlContract = &HTLContract{}

type HTLContract struct {
	Address common.Address     `json:"address"`
	Tx      *types.Transaction `json:"Tx"`
}

func getEthClient() *ethclient.Client {
    address := os.Getenv("OLDATA") + "/" + global.Current.ETHAddress
	if client == nil {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(address)
			if err != nil {
				log.Fatal("failed to get geth ipc ", "err", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	} else if id, _ := client.NetworkID(nil); id == big.NewInt(20180229) {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(address)
			if err != nil {
				log.Fatal("failed to get geth ipc ", "err", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	}
	return client

}

func GetAddress() common.Address {
	auth := GetAuth()
	return auth.From
}

func GetAuth() *bind.TransactOpts {
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
		auth, err := bind.NewTransactor(strings.NewReader(key), "3456")
		if err != nil {
			log.Fatal("Can't get pre-allocate auth for Carol", "err", err)
		}
		return auth
	default:
		log.Info("This node don't have pre-allocate Eth account", "nodeName", nodeName)
		return nil
	}
}

func GetHtlContract() *HTLContract {
	cli := getEthClient()

	if htlContract.Tx == nil {
		auth := GetAuth()
		auth.GasLimit = 2000000
		address, tx, _, err := htlc.DeployHtlc(auth, cli, auth.From)
		if err != nil {
			log.Fatal("Failed to create htlc for the node", "err", err)
		}
		log.Debug("Htlc contract created", "address", address, "tx", tx)
        time.Sleep(6 * time.Second)
		htlContract.Address = address
		htlContract.Tx = tx
		return htlContract

	} else {
		balance, err := htlContract.HTLContractObject().Balance(&bind.CallOpts{Pending: true})
		if err != nil {
			log.Fatal("Previous htlc not callable, re-deploy", "err", err, "address", htlContract.Address)
		}
		log.Warn("htlc already initialed", "address", htlContract.Address, "Tx", htlContract.Tx, "Balance", balance)
		time.Sleep(1 * time.Second)
	}
	return htlContract
}
func (h *HTLContract) HTLContractObject() *htlc.Htlc{
    cli := getEthClient()
    contract, err := htlc.NewHtlc(h.Address, cli)
    if err != nil {
        log.Error("Failed to initial htlc contract from address", "address", h.Address)
        return nil
    }
    return contract
}

func (h *HTLContract) Funds(value *big.Int) error {
	auth := GetAuth()
	auth.Value = EtherToWei(value)
	auth.GasLimit = 200000
	contract := h.HTLContractObject()
	if contract == nil {
	    return errors.New("failed to get htlc contract instance")
    }
	tx, err := contract.Funds(auth)
	if err != nil {
		log.Error("Can't fund the htlc", "err", err, "auth", auth)
		return err
	}
	h.Tx = tx
	balance, err := h.HTLContractObject().Balance(&bind.CallOpts{Pending: true})
	if err != nil {
	    log.Error("Can't get balance after fund", "err", err)
    }
	log.Info("Fund htlc", "address", h.Address, "tx", h.Tx, "value", value, "balance", balance)
    time.Sleep(6 * time.Second)
	return nil
}

func (h *HTLContract) Setup(lockTime *big.Int, receiver common.Address, scrHash [32]byte) error {
	auth := GetAuth()
    auth.GasLimit = 200000
    contract := h.HTLContractObject()
    if contract == nil {
        return errors.New("failed to get htlc contract instance")
    }
    tx, err := contract.Setup(auth, lockTime, receiver, scrHash)
	if err != nil {
		log.Error("Can't setup the htlc", "err", err, "auth", auth)
		return err
	}
	h.Tx = tx
    time.Sleep(6 * time.Second)
	r, _ := contract.Receiver(&bind.CallOpts{Pending: true})
	balance, _ := contract.Balance(&bind.CallOpts{Pending: true})
	log.Info("Setup htlc", "address", h.Address, "tx", h.Tx, "receiver", receiver, "r", r, "balance", balance)
	return nil
}

func (h *HTLContract) Redeem(scr []byte) error {
	auth := GetAuth()

	balance, _ := h.HTLContractObject().Balance(&bind.CallOpts{Pending: false})
	log.Debug("balance before redeem", "balance", balance)

    contract := h.HTLContractObject()
    if contract == nil {
        return errors.New("failed to get htlc contract instance")
    }
    tx, err := contract.Redeem(auth, scr)
	if err != nil {
		log.Error("Can't redeem the htlc", "err", err, "auth", auth)
		return err
	}
	h.Tx = tx
	time.Sleep(6 * time.Second)
    balance, err = h.HTLContractObject().Balance(&bind.CallOpts{Pending: false})
    log.Debug("balance after redeem", "balance", balance, "err", err)
	log.Info("Redeem htlc", "address", h.Address, "tx", tx, "scr", scr, "txaddress", tx.Hash())
	return nil
}

func (h *HTLContract) Refund() error {
	auth := GetAuth()
    contract := h.HTLContractObject()
    if contract == nil {
        return errors.New("failed to get htlc contract instance")
    }
    tx, err := contract.Refund(auth)
	if err != nil {
		log.Error("Can't refund the htlc", "err", err, "auth", auth)
		return err
	}
	h.Tx = tx
	log.Info("Refund htlc", "address", h.Address, "tx", h.Tx, "value", tx.Value())
	return nil
}

func (h *HTLContract) Audit(receiver common.Address, value *big.Int, scrHash [32]byte) error {

    valueWei := EtherToWei(value)
    log.Debug("audit htlc contract", "wei", valueWei)
	result, err := h.HTLContractObject().Audit(&bind.CallOpts{Pending: false}, receiver, valueWei, scrHash)
	if err != nil {
		log.Error("Can't audit the htlc", "err", err)
		return err
	}
	log.Info("Audit htlc", "result", result)
	return nil
}

func (h *HTLContract) Extract() []byte {

    result, err := h.HTLContractObject().ExtractMsg(&bind.CallOpts{Pending: false})
    if err != nil {
        log.Error("Failed to extract secret", "err", err)
    }
    return result
}

func (h *HTLContract) Balance() *big.Int {

	balance, err := h.HTLContractObject().Balance(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("Can't get the balance", "err", err)
	}

	return WeiToEther(balance)
}

func (h *HTLContract) ScrHash() [32]byte {

    sh, err := h.HTLContractObject().ScrHash(&bind.CallOpts{Pending: true})
    if err != nil {
        log.Error("Can't get scrHash", "err", err)
    }
    return sh
}

func EtherToWei(value *big.Int) *big.Int {

    return new(big.Int).Mul(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

func WeiToEther(value *big.Int) *big.Int {

	return new(big.Int).Div(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

func GetHTLCFromMessage(message []byte) *HTLContract {
	log.Debug("Parse message to ETH HTLC")
	register := &HTLContract{}

	result, err := comm.Deserialize(message, register)
	if err != nil {
		log.Error("parse htlc contract failed", "err", err)
		return nil
	}
	return  result.(*HTLContract)
}

func (h *HTLContract) ToMessage() []byte {
    msg, err := comm.Serialize(h)
    if err != nil {
        log.Error("Failed to serialize htlc", "err", err)
    }
    return msg
}

func (h *HTLContract) ToKey() []byte {
    return h.Address.Bytes()
}