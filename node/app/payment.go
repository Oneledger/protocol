package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/status"
)

func CreatePaymentRequest(app Application, identities []id.Identity, quotient data.Coin, height int64) []byte {
	chainId := app.Admin.Get(chainKey)
	inputs := make([]action.SendInput, 0)
	outputs := make([]action.SendOutput, 0)

	for _, identity := range identities {
		if identity.Name == "" {
			log.Error("Missing Party argument")
			return nil
		}

		party, err := app.Identities.FindName(identity.Name)
		if err == status.MISSING_DATA {
			log.Debug("CreatePaymentRequest", "PartyMissingData", err)
			return nil
		}

		partyBalance := app.Balances.Get(party.AccountKey)
		if partyBalance == nil {
			interimBalance := data.NewBalanceFromString(0, "OLT")
			partyBalance = &interimBalance
		}

		//fee := conv.GetCoin(args.Fee, args.Currency)
		//gas := conv.GetCoin(args.Gas, args.Currency)

		inputs = append(inputs,
			action.NewSendInput(party.AccountKey, partyBalance.GetAmountByName("OLT")))

		outputs = append(outputs,
			action.NewSendOutput(party.AccountKey, partyBalance.GetAmountByName("OLT").Plus(quotient)))
	}

	payment, err := app.Accounts.FindName("Payment")
	if err != status.SUCCESS {
		log.Fatal("Payment Account not found")
	}
	paymentBalance := app.Balances.Get(payment.AccountKey())
	log.Debug("CreatePaymentRequest", "paymentBalance", paymentBalance)

	numberValidators := data.NewCoin(int64(len(identities)), "OLT")
	totalPayment := quotient.Multiply(numberValidators)

	inputs = append(inputs,
		action.NewSendInput(payment.AccountKey(), paymentBalance.GetAmountByName("OLT")))

	outputs = append(outputs,
		action.NewSendOutput(payment.AccountKey(), paymentBalance.GetAmountByName("OLT").Minus(totalPayment)))

	// Create base transaction
	send := &action.Payment{
		Base: action.Base{
			Type:     action.PAYMENT,
			ChainId:  string(chainId.([]byte)),
			Owner:    payment.AccountKey(),
			Signers:  GetSigners(payment.AccountKey(), app),
			Sequence: height, //global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     data.NewCoin(0, "OLT"),
		Gas:     data.NewCoin(0, "OLT"),
	}

	return action.PackRequest(SignTransaction(send, app))
}
