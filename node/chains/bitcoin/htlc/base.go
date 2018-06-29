/*
	Copyright 2017 - 2018 Oneledger

	Functions to wrap the atomic swap mechanics for Bitcoin
*/
package htlc

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// Initiate a new swpa between parties
func NewInitiateCmd(cp2Addr *btcutil.AddressPubKeyHash, amount btcutil.Amount, lockTime int64) *InitiateCmd {
	return &InitiateCmd{
		cp2Addr:  cp2Addr,
		amount:   amount,
		lockTime: lockTime,
	}
}

// Audit an existing contract
func NewAuditContractCmd(contract []byte, contractTx *wire.MsgTx) *AuditContractCmd {
	return &AuditContractCmd{
		contract:   contract,
		contractTx: contractTx,
	}
}

// Participate in the swap
func NewParticipateCmd(cp1Addr *btcutil.AddressPubKeyHash, amount btcutil.Amount,
	secretHash [32]byte, lockTime int64) *ParticipateCmd {
	return &ParticipateCmd{
		cp1Addr:    cp1Addr,
		amount:     amount,
		secretHash: secretHash[:],
		lockTime:   lockTime,
	}
}

// Redeem the asset within the contact, the swap is approved
func NewRedeemCmd(contract []byte, contractTx *wire.MsgTx, secret []byte) *RedeemCmd {
	return &RedeemCmd{
		contract:   contract,
		contractTx: contractTx,
		secret:     secret,
	}
}

// Collect the refund of the contract, it has failed somehow
func NewRefundCmd(contract []byte, contractTx *wire.MsgTx) *RefundCmd {
	return &RefundCmd{
		contract:   contract,
		contractTx: contractTx,
	}
}

// Extract the secret from the transaction as issued by the counterparty
func NewExtractSecretCmd(redemptionTx *wire.MsgTx, secretHash [32]byte) *ExtractSecretCmd {
	return &ExtractSecretCmd{
		redemptionTx: redemptionTx,
		secretHash:   secretHash[:],
	}
}
