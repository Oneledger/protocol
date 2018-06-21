/*
	Copyright 2017 - 2018 Oneledger
*/
package htlc

import "github.com/btcsuite/btcutil"

func NewInitiateCmd(cp2Addr *btcutil.AddressPubKeyHash, amount btcutil.Amount, lockTime int64) *InitiateCmd {
	return &InitiateCmd{
		cp2Addr:  cp2Addr,
		amount:   amount,
		lockTime: lockTime,
	}
}

func NewParticipateCmd() *ParticipateCmd {
	return &ParticipateCmd{
		cp1Addr:    nil,
		amount:     2000,
		secretHash: []byte(nil),
		lockTime:   1000,
	}
}

func NewRedeemCmd() *RedeemCmd {
	return &RedeemCmd{
		contract:   nil,
		contractTx: nil,
		secret:     nil,
	}
}

func NewRefundCmd() *RefundCmd {
	return &RefundCmd{
		contract:   nil,
		contractTx: nil,
	}
}

func NewExtractSecretCmd() *ExtractSecretCmd {
	return &ExtractSecretCmd{
		redemptionTx: nil,
		secretHash:   nil,
	}
}

func NewAuditContractCmd() *AuditContractCmd {
	return &AuditContractCmd{
		contract:   nil,
		contractTx: nil,
	}
}
