package governance
//
//import (
//	"github.com/Oneledger/protocol/action"
//	"github.com/Oneledger/protocol/data/governance"
//	"github.com/Oneledger/protocol/data/keys"
//	"github.com/pkg/errors"
//	"github.com/tendermint/tendermint/libs/kv"
//)
//
//var _ action.Msg = &WithdrawProposal{}
//
//type WithdrawProposal struct {
//	proposalType   governance.ProposalType
//	description    string
//	contributor    keys.Address
//	withDrawAmount action.Amount
//	beneficiary    keys.Address
//
//}
//
//func (c WithdrawProposal) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
//	withdrawProposal := WithdrawProposal{}
//	err := withdrawProposal.Unmarshal(signedTx.Data)
//	if err != nil {
//		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
//	}
//
//	//validate basic signature
//	err = action.ValidateBasic(signedTx.RawBytes(), withdrawProposal.Signers(), signedTx.Signatures)
//	if err != nil {
//		return false, err
//	}
//
//	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
//	if err != nil {
//		return false, err
//	}
//
//	options := ctx.Proposals.GetOptionsByType(withdrawProposal.proposalType)
//
//
//	//TODO Check if that same person is asking for withdrawal
//	//and A person cannot ask for withdrawal of another person's funds
//
//
//	//TODO Check the proposal, and make sure that the proposal is in CANCELLED state. Else reject the transaction
//
//	//TODO Check if the contributor has sufficient funds to withdraw for that proposal. Check for the corresponding value in Proposal Fund Store
//
//
//	return true, nil
//}
//
//func (c WithdrawProposal) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	ctx.Logger.Debug("Processing CreateProposal Transaction for CheckTx", tx)
//	return runTx(ctx, tx)
//}
//
//func (c WithdrawProposal) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	ctx.Logger.Debug("Processing CreateProposal Transaction for DeliverTx", tx)
//	return runTx(ctx, tx)
//}
//
//func (c WithdrawProposal) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
//	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
//}
//
//func runTx(ctx *action.Context, signedTx action.RawTx) (bool, action.Response) {
//	withDrawProposal := WithdrawProposal{}
//	err := withDrawProposal.Unmarshal(signedTx.Data)
//	if err != nil {
//		return false, action.Response{}
//	}
//
//
//
//
//	return true, action.Response{}
//}
//
//func (c WithdrawProposal) Signers() []action.Address {
//	panic("implement me")
//}
//
//func (c WithdrawProposal) Type() action.Type {
//	panic("implement me")
//}
//
//func (c WithdrawProposal) Tags() kv.Pairs {
//	panic("implement me")
//}
//
//func (c WithdrawProposal) Marshal() ([]byte, error) {
//	panic("implement me")
//}
//
//func (c WithdrawProposal) Unmarshal(bytes []byte) error {
//	panic("implement me")
//}
