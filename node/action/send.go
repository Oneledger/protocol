/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"bytes"
	"encoding/hex"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

// simple send transaction
type Send struct {
	Base
	SendTo SendTo    `json:"SendTo"`
	Fee    data.Coin `json:"fee"`
	//Gas    data.Coin `json:"gas"`
}

func init() {
	serial.Register(Send{})
	serial.Register(SendTo{})
}

func (transaction *Send) Validate() status.Code {
	log.Debug("Validating Send Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	if !transaction.SendTo.IsValid() {
		log.Debug("the send to is not valid", "sendTo", transaction.SendTo)
		return status.MISSING_VALUE
	}

	if !transaction.Fee.IsCurrency("OLT") {
		log.Debug("Wrong Fee token", "fee", transaction.Fee)
		return status.INVALID
	}

	if transaction.Fee.LessThan(global.Current.MinSendFee) {
		log.Debug("Missing Fee", "fee", transaction.Fee)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *Send) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Send Transaction for CheckTx")

	balances := GetBalances(app)

	balance := balances.Get(transaction.Base.Owner)
	if balance == nil {
		log.Debug("Failed to get the balance of the owner", "owner", transaction.Base.Owner)
		return status.MISSING_VALUE
	}

	if !CheckSendTo(*balance, transaction.SendTo, transaction.Fee) {
		return status.INVALID
	}

	return status.SUCCESS
}

func (transaction *Send) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Send) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Send Transaction for DeliverTx")

	balances := GetBalances(app)

	ownerBalance := balances.Get(transaction.Base.Owner)
	if ownerBalance == nil {
		log.Debug("Failed to get the balance of the owner", "owner", transaction.Base.Owner)
		return status.MISSING_VALUE
	}

	if !CheckSendTo(*ownerBalance, transaction.SendTo, transaction.Fee) {
		return status.INVALID
	}

	//change owner balance
	ownerBalance.MinusAmount(transaction.SendTo.Amount)
	ownerBalance.MinusAmount(transaction.Fee)

	//change receiver balance
	receiverBalance := balances.Get(transaction.SendTo.AccountKey)
	if receiverBalance == nil {
		receiverBalance = data.NewBalance()
	}
	receiverBalance.AddAmount(transaction.SendTo.Amount)

	accounts := GetAccounts(app)
	payment, err := accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Error("Failed to get payment account", "status", err)
		return err
	}

	paymentBalance := balances.Get(payment.AccountKey())
	paymentBalance.AddAmount(transaction.Fee)

	balances.Set(transaction.Base.Owner, ownerBalance)
	balances.Set(transaction.SendTo.AccountKey, receiverBalance)
	balances.Set(payment.AccountKey(), paymentBalance)

	return status.SUCCESS
}

func (transaction *Send) Resolve(app interface{}) Commands {
	return []Command{}
}

func (transaction Send) TransactionTags(app interface{}) Tags {
	tags := transaction.Base.TransactionTags(app)

	tagReceiver := transaction.SendTo.AccountKey.String()
	tag1 := common.KVPair{
		Key:   []byte("tx.receiver"),
		Value: []byte(tagReceiver),
	}

	participants := hex.EncodeToString(transaction.Owner) + "," + hex.EncodeToString(transaction.Target)

	tag2 := common.KVPair{
		Key:   []byte("tx.participants"),
		Value: []byte(participants),
	}

	tags = append(tags, tag1)
	tags = append(tags, tag2)

	return tags
}

type SendTo struct {
	AccountKey id.AccountKey `json:"account_key"`
	Amount     data.Coin     `json:"coin"`
}

func (st SendTo) IsValid() bool {
	if st.Amount.LessThan(0) {
		return false
	}
	//todo: should check the transaction.To is a validate accountkey
	if bytes.Compare(st.AccountKey, []byte("")) == 0 {
		return false
	}
	return true
}

func CheckSendTo(balance data.Balance, sendTo SendTo, fee data.Coin) bool {

	//todo: the lessthan check and is valid check should be move into the Coin.IsValid() function together
	if sendTo.Amount.LessThan(0) {
		log.Debug("FAILED: Less Than 0", "out", sendTo)
		return false
	}

	if !sendTo.Amount.IsValid() {
		log.Debug("FAILED: sent amount on Currency is not valid")
		return false
	}

	total := data.NewBalance()
	total.AddAmount(sendTo.Amount)
	total.AddAmount(fee)
	log.Debug("send check", "balance", balance, "total", total)
	if !balance.IsEnoughBalance(*total) {
		log.Debug("FAILED: Balance not enough for the send", "balance", balance, "sendTo", sendTo)
		return false
	}

	return true
}

func TransferVT(app interface{}, validator types.Validator) status.Code {
	log.Debug("Processing Transfer of VT to Zero Account")

	balances := GetBalances(app)
	identities := GetIdentities(app)
	accounts := GetAccounts(app)

	validatorId := identities.FindTendermint(hex.EncodeToString(validator.Address))

	if validatorId.Name == "" {
		log.Error("Missing validator identity argument")
		return status.MISSING_DATA
	}

	validatorBalance := balances.Get(validatorId.AccountKey)
	if validatorBalance == nil {
		log.Debug("Failed to get the balance of the validator", "validatorAccountKey", validatorId.AccountKey)
		return status.MISSING_VALUE
	}

	amount := validatorBalance.GetAmountByName("VT")
	validatorBalance.MinusAmount(amount)

	zeroId, err := accounts.FindName("Zero")

	if err != status.SUCCESS {
		log.Error("Failed to get Zero account", "status", err)
		return err
	}

	zeroBalance := balances.Get(zeroId.AccountKey())
	if zeroBalance == nil {
		zeroBalance = data.NewBalance()
	}
	zeroBalance.AddAmount(amount)

	balances.Set(validatorId.AccountKey, validatorBalance)
	balances.Set(zeroId.AccountKey(), zeroBalance)

	return status.SUCCESS
}
