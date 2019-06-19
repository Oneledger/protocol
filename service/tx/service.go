package tx

import (
	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func Name() string {
	return "tx"
}

type Service struct {
	balances    *balance.Store
	router      action.Router
	accounts    accounts.Wallet
	logger      *log.Logger
	nodeContext node.Context
}

func NewService(
	balances *balance.Store,
	router action.Router,
	accounts accounts.Wallet,
	nodeCtx node.Context,
	logger *log.Logger,
) *Service {
	return &Service{
		balances: balances,
		router:   router,
		accounts: accounts,
		logger:   logger,
	}
}

// SendTx exists for maintaining backwards compatibility with existing olclient implementations. It returns
// a signed transaction
// TODO: deprecate this
func (svc *Service) SendTx(args client.SendTxRequest, reply *client.SendTxReply) error {
	send := action.Send{
		From:   keys.Address(args.From),
		To:     keys.Address(args.To),
		Amount: args.Amount,
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.BaseTx{
		Data: send,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	if _, err := svc.accounts.GetAccount(args.From); err != nil {
		return errors.New("Account doesn't exist. Send a raw tx instead")
	}

	pubKey, signed, err := svc.accounts.SignWithAddress(tx.Bytes(), send.From)
	if err != nil {
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}

	return nil
}

func (svc *Service) ApplyValidator(args client.ApplyValidatorRequest, reply *client.ApplyValidatorReply) error {
	if len(args.Name) < 1 {
		args.Name = svc.nodeContext.NodeName
	}

	if len(args.Address) < 1 {
		handler, err := svc.nodeContext.PubKey().GetHandler()
		if err != nil {
			return err
		}
		args.Address = handler.Address()
	}

	pubkey := &keys.PublicKey{keys.GetAlgorithmFromTmKeyName(args.TmPubKeyType), args.TmPubKey}
	if len(args.TmPubKey) < 1 {
		*pubkey = svc.nodeContext.ValidatorPubKey()
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

	pubKey, signed, err := svc.accounts.SignWithAccountIndex(tx.Bytes(), 0)
	if err != nil {
		return err
	}
	tx.Signatures = []action.Signature{{pubKey, signed}}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return errors.Wrap(err, "err while network serialization")
	}

	*reply = client.ApplyValidatorReply{RawTx: packet}

	return nil
}
