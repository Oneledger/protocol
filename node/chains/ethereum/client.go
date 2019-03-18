package ethereum

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Oneledger/protocol/node/data"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"context"
	"github.com/Oneledger/protocol/node/chains/ethereum/htlc"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func init() {
	serial.Register(HTLContract{})
	serial.Register(common.Address{})
	serial.Register(atomic.Value{})
	serial.Register(common.Hash{})
}

var client *ethclient.Client

type HTLContract struct {
	Address common.Address `json:"address"`
	TxHash  common.Hash    `json:"txhash"`
}

func getEthClient() *ethclient.Client {
	address := global.Current.ETHAddress
	if client == nil {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(address)
			if err != nil {
				log.Fatal("failed to get geth ", "status", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	} else if id, _ := client.NetworkID(nil); id == big.NewInt(20180229) {
		for i := 0; i < 3; i++ {
			cli, err := ethclient.Dial(address)
			if err != nil {
				log.Fatal("failed to get geth ", "status", err)
				time.Sleep(3 * time.Second)
			}
			return cli
		}
	}
	return client

}

func GetAddress(chainKey interface {}) common.Address {
	auth := GetAuth(chainKey)
	return auth.From
}

func GetAuth(chainKey interface {}) *bind.TransactOpts {

	if (chainKey == nil) {
		return nil
	}

	chainKeyJson, err := base64.StdEncoding.DecodeString(chainKey.(string))

	if err != nil {
		return nil
	}

	var chainKeyMap map[string]interface{}

	json.Unmarshal(chainKeyJson, &chainKeyMap)

	key := chainKeyMap["key"].(string)
	password := chainKeyMap["password"].(string)

	keyDecoded, err2 := base64.StdEncoding.DecodeString(key)

	if err2 != nil {
		return nil
	}

	auth, err := bind.NewTransactor(strings.NewReader(string(keyDecoded)), password)

	if err != nil {
		log.Fatal("Can't get Ethereum auth", "status", err)
	}

	return auth
}

func CreateHtlContract(chainKey interface {}) *HTLContract {
	cli := getEthClient()

	auth := GetAuth(chainKey)
	auth.GasLimit = 2000000
	address, tx, _, err := htlc.DeployHtlc(auth, cli, auth.From)
	if err != nil {
		log.Fatal("Failed to create htlc for the node", "status", err)
	}
	log.Debug("Htlc contract created", "address", address, "tx", tx.Hash())
	htlContract := &HTLContract{Address: address, TxHash: tx.Hash()}

	WaitTxSuccess(cli, tx, 20*time.Second, 1*time.Second)

	return htlContract
}

func WaitTxSuccess(client *ethclient.Client, tx *types.Transaction, maxWait time.Duration, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				result, err := client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					if result.Status == types.ReceiptStatusSuccessful {
						ticker.Stop()
					}
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(maxWait)
	close(stop)
	return
}

func (h *HTLContract) HTLContractObject() *htlc.Htlc {
	cli := getEthClient()
	contract, err := htlc.NewHtlc(h.Address, cli)
	if err != nil {
		log.Error("Failed to initial htlc contract from address", "address", h.Address)
		return nil
	}
	return contract
}

func (h *HTLContract) Funds(chainKey interface {}, value *big.Int, lockTime *big.Int, receiver common.Address, scrHash [32]byte) error {
	auth := GetAuth(chainKey)
	auth.Value = EtherToWei(value)
	auth.GasLimit = 400000
	contract := h.HTLContractObject()
	if contract == nil {
		return errors.New("failed to get htlc contract instance")
	}
	tx, err := contract.Funds(auth, lockTime, receiver, scrHash)
	if err != nil {
		log.Error("Can't fund the htlc", "status", err, "auth", auth)
		return err
	}
	h.TxHash = tx.Hash()
	balance, err := h.HTLContractObject().Balance(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("Can't get balance after fund", "status", err)
	}
	log.Info("Fund htlc", "address", h.Address, "tx", h.TxHash, "value", value, "balance", balance)
	WaitTxSuccess(client, tx, 10*time.Second, 1*time.Second)
	return nil
}

func (h *HTLContract) Redeem(chainKey interface{}, scr []byte) error {
	auth := GetAuth(chainKey)

	balance, _ := h.HTLContractObject().Balance(&bind.CallOpts{Pending: false})
	log.Debug("balance before redeem", "balance", balance)

	contract := h.HTLContractObject()
	if contract == nil {
		return errors.New("failed to get htlc contract instance")
	}
	tx, err := contract.Redeem(auth, scr)
	if err != nil {
		log.Error("Can't redeem the htlc", "status", err, "auth", auth)
		return err
	}
	h.TxHash = tx.Hash()
	WaitTxSuccess(client, tx, 10*time.Second, 1*time.Second)
	balance, err = contract.Balance(&bind.CallOpts{Pending: false})
	log.Debug("balance after redeem", "balance", balance, "status", err)
	log.Info("Redeem htlc", "address", h.Address, "tx", h.TxHash, "scr", scr, "txaddress", tx.Hash())
	return nil
}

func (h *HTLContract) Refund(chainKey interface{}) error {
	auth := GetAuth(chainKey)
	auth.GasLimit = 300000
	contract := h.HTLContractObject()
	if contract == nil {
		return errors.New("failed to get htlc contract instance")
	}
	tx, err := contract.Refund(auth)
	if err != nil {
		log.Error("Can't refund the htlc", "status", err, "auth", auth)
		return err
	}
	WaitTxSuccess(client, tx, 10*time.Second, 1*time.Second)

	h.TxHash = tx.Hash()

	balance, err := contract.Balance(&bind.CallOpts{Pending: false})

	log.Debug("Refund htlc", "address", h.Address, "tx", h.TxHash, "balance", balance)

	return nil
}

func (h *HTLContract) Audit(receiver common.Address, value *big.Int, scrHash [32]byte) error {

	valueWei := EtherToWei(value)
	log.Debug("audit htlc contract", "wei", valueWei)
	result, err := h.HTLContractObject().Audit(&bind.CallOpts{Pending: false}, receiver, valueWei, scrHash)
	if err != nil {
		log.Error("Can't audit the htlc", "status", err)
		return err
	}
	log.Info("Audit htlc", "result", result)
	return nil
}

func (h *HTLContract) Extract() []byte {

	result, err := h.HTLContractObject().ExtractMsg(&bind.CallOpts{Pending: false})
	if err != nil {
		log.Error("Failed to extract secret", "status", err)
	}
	return result
}

func (h *HTLContract) Balance() *big.Int {

	balance, err := h.HTLContractObject().Balance(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("Can't get the balance", "status", err)
	}

	return WeiToEther(balance)
}

func (h *HTLContract) ScrHash() [32]byte {

	sh, err := h.HTLContractObject().ScrHash(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Error("Can't get scrHash", "status", err)
	}
	return sh
}

func EtherToWei(value *big.Int) *big.Int {

	return new(big.Int).Mul(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

func WeiToEther(value *big.Int) *big.Int {

	return new(big.Int).Div(value, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}

func (h *HTLContract) Chain() data.ChainType {
	return data.ETHEREUM
}

func (h *HTLContract) ToBytes() []byte {
	msg, err := serial.Serialize(h, serial.JSON)
	if err != nil {
		log.Error("Failed to serialize htlc", "err", err)
	}
	return msg
}

func (h *HTLContract) ToKey() []byte {
	return h.Address.Bytes()
}

func (h *HTLContract) FromBytes(message []byte) {
	_, err := serial.Deserialize(message, h, serial.JSON)
	if err != nil {
		log.Error("Failed deserialize ETH contract", "contract", message, "err", err)
	}
	return
}
