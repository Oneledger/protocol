package governance

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &CreateProposal{}

type CreateProposal struct {
	ProposalID     governance.ProposalID
	ProposalType   governance.ProposalType
	Description    string
	Proposer       keys.Address
	InitialFunding action.Amount
}

func (c CreateProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createProposal := CreateProposal{}
	err := createProposal.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), createProposal.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	options := ctx.ProposalMasterStore.Proposal.GetOptionsByType(createProposal.ProposalType)
	if options == nil {
		return false, errors.New("failed to get proposal options")
	}

	// the currency should be OLT
	currency, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if currency.Name != createProposal.InitialFunding.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, createProposal.InitialFunding.String())
	}

	//Check if Proposal ID is valid
	if len(createProposal.ProposalID) <= 0 {
		return false, errors.New("invalid proposal id")
	}

	//Check if initial funding is greater than minimum amount based on type.
	coin := createProposal.InitialFunding.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(*options.InitialFunding)) {
		return false, action.ErrInvalidAmount
	}

	//Check if Proposal Type is valid
	switch createProposal.ProposalType {
	case governance.ProposalTypeGeneral:
	case governance.ProposalTypeCodeChange:
	case governance.ProposalTypeConfigUpdate:
	default:
		return false, errors.New("invalid proposal type")
	}

	//Check if proposer address is valid oneLedger address
	err = createProposal.Proposer.Err()
	if err != nil {
		return false, errors.Wrap(err, "invalid proposer address")
	}

	if len(createProposal.Description) == 0 {
		return false, errors.New("invalid description of proposal")
	}

	return true, nil
}

func (c CreateProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for CheckTx", tx)
	return runTx(ctx, tx)
}

func (c CreateProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing CreateProposal Transaction for DeliverTx", tx)
	return runTx(ctx, tx)
}

func (c CreateProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runTx(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createProposal := CreateProposal{}
	err := createProposal.Unmarshal(tx.Data)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_failed_deserialize"),
		}
		return false, result
	}

	//Get Proposal options based on type.
	options := ctx.ProposalMasterStore.Proposal.GetOptionsByType(createProposal.ProposalType)

	//Check if initial funding is greater than minimum amount based on type.
	coin := createProposal.InitialFunding.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(*options.InitialFunding)) {
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_insufficient_funds"),
		}
		return false, result
	}

	//Create Proposal and save to Proposal Store
	proposal := governance.NewProposal(
		createProposal.ProposalID,
		createProposal.ProposalType,
		createProposal.Description,
		createProposal.Proposer,
		options.FundingDeadline,
		options.FundingGoal,
		options.VotingDeadline,
		options.PassPercentage)

	//Check if Proposal already exists
	if ctx.ProposalMasterStore.Proposal.Exists(proposal.ProposalID) {
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_failed"),
		}
		return false, result
	}

	err = ctx.ProposalMasterStore.Proposal.Set(proposal)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_failed"),
		}
		return false, result
	}

	//Set generated Proposal ID in transaction response
	createProposal.ProposalID = proposal.ProposalID

	//Deduct initial funding from proposer's address
	funds := createProposal.InitialFunding.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(createProposal.Proposer.Bytes(), funds)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_deduction_failed"),
		}
		return false, result
	}

	//Add initial funds to the Proposal Fund store
	initialFunding := balance.NewAmountFromBigInt(createProposal.InitialFunding.Value.BigInt())
	err = ctx.ProposalMasterStore.ProposalFund.AddFunds(proposal.ProposalID, proposal.Proposer, initialFunding)
	if err != nil {
		//return Funds back to proposer.
		err = ctx.Balances.AddToAddress(createProposal.Proposer, funds)
		if err != nil {
			panic("error returning funds to balance store")
		}
		result := action.Response{
			Events: action.GetEvent(createProposal.Tags(), "create_proposal_funding_failed"),
		}
		return false, result
	}

	result := action.Response{
		Events: action.GetEvent(createProposal.Tags(), "create_proposal_success"),
	}

	return true, result
}

func (c CreateProposal) Signers() []action.Address {
	return []action.Address{c.Proposer}
}

func (c CreateProposal) Type() action.Type {
	return action.PROPOSAL_CREATE
}

func (c CreateProposal) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(c.ProposalID),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(c.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.proposer"),
		Value: c.Proposer.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.proposalType"),
		Value: []byte(string(c.ProposalType)),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.initialFunding"),
		Value: []byte(c.InitialFunding.String()),
	}

	tags = append(tags, tag, tag1, tag2, tag3, tag4)
	return tags
}

func (c CreateProposal) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CreateProposal) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, c)
}
