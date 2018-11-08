package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

func CreatePaymentRequest(app Application, identities []id.Identity, quotient data.Coin) []byte {
	var signers []id.PublicKey
	chainId := app.Admin.Get(chainKey)
	inputs := make([]action.SendInput, 0)
	outputs := make([]action.SendOutput, 0)

	for _, identity := range identities {
		if identity.Name == "" {
			log.Error("Missing Party argument")
			return nil
		}

		// TODO: Can't convert identities to accounts, this way!
		log.Debug("CreatePaymentRequest", "IdentityName", identity.Name)

		party, status := app.Identities.FindName(identity.Name)
		if status == err.MISSING_DATA {
			log.Debug("CreatePaymentAccount1", "MissingDataStatus", status)
			return nil
		}

		//if status != err.SUCCESS {
		//	log.Fatal("CreatePaymentRequest", "SuccessStatus", status)
		//}

		partyAccountKey := party.AccountKey
		log.Debug("CreatePaymentAccountKey", "AccountKey", partyAccountKey)
		partyBalance := app.Utxo.Get(party.AccountKey)
		if partyBalance == nil {
			interimBalance := data.NewBalance(0, "OLT")
			partyBalance = &interimBalance
		}
		log.Warn("CreatePaymentRequest", "UTXOPartyBalance", partyBalance)

		//log.Dump("AccountKeys", party, counterParty)

		//if args.Currency == "" || args.Amount == "" {
		//	log.Error("Missing an amount argument")
		//	return nil
		//}

		//log.Dump("Balances", partyBalance, counterPartyBalance)

		if partyBalance == nil {
			log.Error("Missing Balance", "party", partyBalance)
			return nil
		}

		//fee := conv.GetCoin(args.Fee, args.Currency)
		//gas := conv.GetCoin(args.Gas, args.Currency)

		inputs = append(inputs,
			action.NewSendInput(party.AccountKey, partyBalance.Amount))

		outputs = append(outputs,
			action.NewSendOutput(party.AccountKey, partyBalance.Amount.Plus(quotient)))
	}

	payment, status := app.Accounts.FindName("Payment-OneLedger")
	if status != err.SUCCESS {
		log.Fatal("Payment Account not found")
	}
	paymentBalance := app.Utxo.Get(payment.AccountKey())
	log.Debug("CreatePaymentRequest", "paymentBalance", paymentBalance)

	numberValidators := data.NewCoin(int64(len(identities)), "OLT")
	fullPayment := quotient.Multiply(numberValidators)

	inputs = append(inputs,
		action.NewSendInput(payment.AccountKey(), paymentBalance.Amount))

	outputs = append(outputs,
		action.NewSendOutput(payment.AccountKey(), paymentBalance.Amount.Minus(fullPayment)))

	// Create base transaction
	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  string(chainId.([]byte)),
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     data.NewCoin(0, "OLT"),
		Gas:     data.NewCoin(0, "OLT"),
	}

	signed := action.SignTransaction(send)
	packet := action.PackRequest(action.SEND, signed)

	return packet
}

func CreatePaymentRequest2(app Application, identity id.Identity, quotient data.Coin) []byte {
	var signers []id.PublicKey
	chainId := app.Admin.Get(chainKey)
	inputs := make([]action.SendInput, 0)
	outputs := make([]action.SendOutput, 0)

	if identity.Name == "" {
		log.Error("Missing Party argument")
		return nil
	}

	// TODO: Can't convert identities to accounts, this way!
	log.Debug("CreatePaymentRequest", "IdentityName", identity.Name)

	party, status := app.Identities.FindName(identity.Name)
	if status == err.MISSING_DATA {
		log.Debug("CreatePaymentAccount1", "MissingDataStatus", status)
		return nil
	}

	//if status != err.SUCCESS {
	//	log.Fatal("CreatePaymentRequest", "SuccessStatus", status)
	//}

	partyAccountKey := party.AccountKey
	log.Debug("CreatePaymentAccountKey", "AccountKey", partyAccountKey)
	partyBalance := app.Utxo.Get(party.AccountKey)

	log.Debug("CreatePaymentBalance", "partyBalance", partyBalance)

	//log.Dump("AccountKeys", party, counterParty)

	//if args.Currency == "" || args.Amount == "" {
	//	log.Error("Missing an amount argument")
	//	return nil
	//}

	//log.Dump("Balances", partyBalance, counterPartyBalance)

	if partyBalance == nil {
		log.Error("Missing Balance", "party", partyBalance)
		return nil
	}

	//fee := conv.GetCoin(args.Fee, args.Currency)
	//gas := conv.GetCoin(args.Gas, args.Currency)

	inputs = append(inputs,
		action.NewSendInput(party.AccountKey, partyBalance.Amount))

	outputs = append(outputs,
		action.NewSendOutput(party.AccountKey, partyBalance.Amount.Plus(quotient)))

	payment, status := app.Accounts.FindName("Payment-OneLedger")
	if status != err.SUCCESS {
		log.Fatal("Payment Account not found")
	}
	paymentBalance := app.Utxo.Get(payment.AccountKey())
	log.Debug("CreatePaymentRequest", "paymentBalance", paymentBalance)

	//numberValidators := data.NewCoin(int64(len(identities)), "OLT")
	//fullPayment := quotient.Multiply(numberValidators)

	inputs = append(inputs,
		action.NewSendInput(payment.AccountKey(), paymentBalance.Amount))

	outputs = append(outputs,
		action.NewSendOutput(payment.AccountKey(), paymentBalance.Amount.Minus(quotient)))

	// Create base transaction
	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  string(chainId.([]byte)),
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     data.NewCoin(0, "OLT"),
		Gas:     data.NewCoin(0, "OLT"),
	}

	signed := action.SignTransaction(send)
	packet := action.PackRequest(action.SEND, signed)

	return packet
}
