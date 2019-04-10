package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/status"
)

func CreatePaymentRequest(app Application, quotient data.Coin, height int64) action.Transaction {
	chainId := app.Admin.Get(chainKey)
	identities := app.Validators.Approved

	log.Debug("Paying to", "len", len(identities), "identities", identities)

	sendto := make([]action.SendTo, len(identities))

	for i, identity := range identities {
		if identity.Name == "" {
			log.Error("Missing Party argument")
			return nil
		}

		party, err := app.Identities.FindName(identity.Name)
		if err == status.MISSING_DATA {
			log.Debug("CreatePaymentRequest", "PartyMissingData", err)
			return nil
		}

		//todo : here we use index of the approved identity in the
		// Validators.approved to simplify the verification of validators in payment
		// need to makes this more secure.
		sendto[i] = action.SendTo{
			AccountKey: party.AccountKey,
			Amount:     quotient,
		}
	}

	payment, err := app.Accounts.FindName(global.Current.PaymentAccount)
	if err != status.SUCCESS {
		log.Fatal("Payment Account not found")
	}

	paymentBalance := app.Balances.Get(payment.AccountKey(), true)
	log.Debug("CreatePaymentRequest", "paymentBalance", paymentBalance)

	//numberValidators := data.NewCoinFromInt(int64(len(identities)), "OLT")
	//totalPayment := quotient.Multiply(numberValidators)

	/*
		inputs = append(inputs,
			action.NewSendInput(payment.AccountKey(), paymentBalance.GetCoinByName("OLT")))

		outputs = append(outputs,
			action.NewSendOutput(payment.AccountKey(), paymentBalance.GetCoinByName("OLT").Minus(totalPayment)))
	*/

	// Create base transaction
	send := &action.Payment{
		Base: action.Base{
			Type:     action.PAYMENT,
			ChainId:  string(chainId.([]byte)),
			Owner:    payment.AccountKey(),
			Signers:  GetSigners(payment.AccountKey(), app),
			Sequence: height, //global.Current.Sequence,
		},
		SendTo: sendto,
	}

	return send
}
