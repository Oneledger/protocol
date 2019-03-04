package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
	"github.com/Oneledger/protocol/node/chains/common"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"math/big"
)

type BitcoinDriver struct {
	Driver ChainDriverBase
}

func init() {
	serial.Register(BitcoinDriver{})
}

func (driver BitcoinDriver) GetURL() string {
	return global.Current.BTCAddress
}

func (driver BitcoinDriver) GetChainAddress(chainKey interface{}) []byte{
	cli := bitcoin.GetBtcClient(global.Current.BTCAddress, chainKey)
	return []byte(bitcoin.GetRawAddress(cli).String())
}

func (driver BitcoinDriver) GetMethodsList() []string {
	list := []string{"getinfo"}

	return list
}

func (driver BitcoinDriver) ExecuteMethod(method string, params []byte) status.Code {
	return status.NOT_IMPLEMENTED
}

func (driver BitcoinDriver) GetAddressFromByteArray(address []byte) interface{} {
	result, err := btcutil.DecodeAddress(string(address), &chaincfg.RegressionNetParams)

	if err != nil {
		log.Error("failed to get addressPubKeyHash")
		return nil
	}

	return result.(*btcutil.AddressPubKeyHash)
}

func (driver BitcoinDriver) CreateSwapContract(receiver interface{}, account id.Account, value big.Int, timeout int64, hash [32]byte) common.Contract {

	cli := bitcoin.GetBtcClient(global.Current.BTCAddress, account.GetChainKey())

	amount := bitcoin.GetAmount(value.String())

	initCmd := htlc.NewInitiateCmd(receiver, amount, timeout, hash)

	_, err := initCmd.RunCommand(cli)
	if err != nil {
		log.Error("Bitcoin Initiate", "status", err)
		return nil
	}

	contract := &bitcoin.HTLContract{}
	contract.FromMsgTx(initCmd.Contract, initCmd.ContractTx)

	return contract
}

func (driver BitcoinDriver) CreateSwapContractFromMessage(message []byte) common.Contract{
	contract := &bitcoin.HTLContract{}

	contract.FromBytes(message)

	return contract
}
