/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package app

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RPCServerContext struct {
	nodeName     string
	balances     *balance.Store
	accounts     accounts.Wallet
	currencies   *balance.CurrencyList
	cfg          config.Server
	nodeContext  node.Context
	validatorSet *identity.ValidatorStore

	services *client.ExtServiceContext
	logger   *log.Logger
}

func (h *RPCServerContext) ApplyValidator(args client.ApplyValidatorRequest, reply *client.ApplyValidatorReply) error {
	if len(args.Name) < 1 {
		args.Name = h.nodeName
	}

	if len(args.Address) < 1 {
		handler, err := h.nodeContext.PubKey().GetHandler()
		if err != nil {
			return err
		}
		args.Address = handler.Address()
	}

	pubkey := &keys.PublicKey{keys.GetAlgorithmFromTmKeyName(args.TmPubKeyType), args.TmPubKey}
	if len(args.TmPubKey) < 1 {
		*pubkey = h.nodeContext.ValidatorPubKey()
	}

	handler, err := pubkey.GetHandler()
	if err != nil {

		return err
	}

	addr := handler.Address()
	apply := action.ApplyValidator{
		Address:          keys.Address(args.Address),
		Stake:            action.Amount{Currency: "VT", Value: args.Amount},
		NodeName:         args.Name,
		ValidatorAddress: addr,
		ValidatorPubKey:  *pubkey,
	}

	uuidNew, _ := uuid.NewUUID()
	tx := action.BaseTx{
		Data: apply,
		Fee:  action.Fee{action.Amount{Currency: "OLT", Value: "0.1"}, 1},
		Memo: uuidNew.String(),
	}

	pubKey, signed, err := h.accounts.SignWithAccountIndex(tx.Bytes(), 0)
	if err != nil {
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	*reply = client.ApplyValidatorReply{packet}

	return nil
}

// ListValidator returns a list of all validator
func (h *RPCServerContext) ListValidators(_ client.ListValidatorsRequest, reply *client.ListValidatorsReply) error {
	validators, err := h.validatorSet.GetValidatorSet()
	if err != nil {
		return errors.Wrap(err, "err while retrieving validators info")
	}

	*reply = client.ListValidatorsReply{
		Validators: validators,
		Height:     h.balances.Version,
	}

	return nil
}
