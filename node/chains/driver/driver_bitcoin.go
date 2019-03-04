package chaindriver

import (
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/chains/bitcoin/htlc"
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

func (driver BitcoinDriver) CreateSwapContract(receiver []byte, account id.Account, value big.Int, timeout int64, hash [32]byte) Contract {
	address, err := btcutil.DecodeAddress(string(receiver), &chaincfg.RegressionNetParams)

	if err != nil {
		log.Error("failed to get addressPubKeyHash")
		return nil
	}

	cli := bitcoin.GetBtcClient(global.Current.BTCAddress, account.GetChainKey())

	amount := bitcoin.GetAmount(value.String())

	initCmd := htlc.NewInitiateCmd(address.(*btcutil.AddressPubKeyHash), amount, timeout, hash)

	_, err = initCmd.RunCommand(cli)
	if err != nil {
		log.Error("Bitcoin Initiate", "status", err)
		return nil
	}

	contract := &bitcoin.HTLContract{}
	contract.FromMsgTx(initCmd.Contract, initCmd.ContractTx)

	return contract
}

func (driver BitcoinDriver) CreateRedeemContract(contract Contract, account id.Account, hash [32]byte) Contract {
	contractBTC := contract.(*bitcoin.HTLContract)

	cmd := htlc.NewRedeemCmd(contractBTC.Contract, contractBTC.GetMsgTx(), hash[:])

	cli := bitcoin.GetBtcClient(global.Current.BTCAddress, account.GetChainKey())

	_, e := cmd.RunCommand(cli)
	if e != nil {
		log.Error("Bitcoin redeem htlc", "status", e)
		return nil
	}

	redeemcontract := &bitcoin.HTLContract{}
	redeemcontract.FromMsgTx(contractBTC.Contract, cmd.RedeemContractTx)

	return redeemcontract
}

func (driver BitcoinDriver) CreateSwapContractFromMessage(message []byte) Contract{
	contract := &bitcoin.HTLContract{}

	contract.FromBytes(message)

	return contract
}
