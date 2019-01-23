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

	balance := balances.Get(transaction.Base.Owner, false)
	if balance == nil {
		log.Debug("Failed to get the balance of the owner", "owner", transaction.Base.Owner)
		return status.MISSING_VALUE
	}

	if !IsEnoughBalance(*balance, transaction.SendTo.Amount, transaction.Fee) {
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

	ownerBalance := balances.Get(transaction.Base.Owner, false)
	if ownerBalance == nil {
		log.Debug("Failed to get the balance of the owner", "owner", transaction.Base.Owner)
		return status.MISSING_VALUE
	}

	if !IsEnoughBalance(*ownerBalance, transaction.SendTo.Amount, transaction.Fee) {
		log.Debug("Owner balance is not enough", "balance", ownerBalance, "amount", transaction.SendTo.Amount, "fee", transaction.Fee)
		return status.INVALID
	}

	//change owner balance
	ownerBalance.MinusAmount(transaction.SendTo.Amount)
	ownerBalance.MinusAmount(transaction.Fee)
	balances.Set(transaction.Base.Owner, ownerBalance)

	//change receiver balance
	receiverBalance := balances.Get(transaction.SendTo.AccountKey, false)
	if receiverBalance == nil {
		receiverBalance = data.NewBalance()
	}
	receiverBalance.AddAmount(transaction.SendTo.Amount)
	balances.Set(transaction.SendTo.AccountKey, receiverBalance)

	accounts := GetAccounts(app)
	payment, err := accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Fatal("Failed to get payment account", "status", err)
		return err
	}
	paymentBalance := balances.Get(payment.AccountKey(), false)
	paymentBalance.AddAmount(transaction.Fee)
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

func IsEnoughBalance(balance data.Balance, value data.Coin, fee data.Coin) bool {

	//todo: the lessthan check and is valid check should be move into the Coin.IsValid() function together
	if value.LessThan(0) {
		log.Debug("FAILED: Less Than 0", "value", value)
		return false
	}

	if !value.IsValid() {
		log.Debug("FAILED: sent amount on Currency is not valid")
		return false
	}

	total := data.NewBalance()
	total.AddAmount(value)
	total.AddAmount(fee)
	log.Debug("send check", "balance", balance, "total", total)
	if !balance.IsEnoughBalance(*total) {
		log.Debug("FAILED: Balance not enough for the send", "balance", balance, "value", value, "fee", fee)
		return false
	}

	return true
}

func TransferVT(app interface{}, applyValidator id.ApplyValidator) status.Code {
	log.Debug("Processing Transfer of VT to Payment Account")

	balances := GetBalances(app)
	identities := GetIdentities(app)
	accounts := GetAccounts(app)

	validatorId := identities.FindTendermint(hex.EncodeToString(applyValidator.Validator.Address))

	if validatorId.Name == "" {
		log.Error("Missing validator identity argument")
		return status.MISSING_DATA
	}

	validatorBalance := balances.Get(validatorId.AccountKey, false)
	if validatorBalance == nil {
		log.Debug("Failed to get the balance of the validator", "validatorAccountKey", validatorId.AccountKey)
		return status.MISSING_VALUE
	}

	amount := applyValidator.Stake
	validatorBalance.MinusAmount(amount)

	paymentId, err := accounts.FindName("Payment")

	if err != status.SUCCESS {
		log.Error("Failed to get Payment account", "status", err)
		return err
	}

	paymentBalance := balances.Get(paymentId.AccountKey(), false)
	if paymentBalance == nil {
		paymentBalance = data.NewBalance()
	}
	paymentBalance.AddAmount(amount)

	balances.Set(validatorId.AccountKey, validatorBalance)
	balances.Set(paymentId.AccountKey(), paymentBalance)

	return status.SUCCESS
}
