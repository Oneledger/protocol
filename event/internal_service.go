/*

 */

package event

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	gov_action "github.com/Oneledger/protocol/action/governance"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/consensus"
	ethereum2 "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/bytes"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Service struct {
	nodeCtx node.Context

	logger *log.Logger
	router action.Router

	//only support local client for broadcasting internal txs
	tmrpc *tmclient.Local
}

func NewService(ctx node.Context, logger *log.Logger, router action.Router, tmnode *consensus.Node) *Service {
	return &Service{
		nodeCtx: ctx,
		logger:  logger,
		router:  router,
		tmrpc:   tmclient.NewLocal(tmnode),
	}
}

func (svc Service) allowedTx(tx action.RawTx) error {
	h := svc.router.Handler(tx.Type)
	if h == nil {
		return errors.New("transaction type not allowed")
	}
	return nil
}

type InternalBroadcastRequest struct {
	RawTx action.RawTx `json:"rawTx"`
}

type BroadcastReply struct {
	TxHash bytes.HexBytes `json:"txHash"`
	// OK indicates whether this broadcast was a request.
	// For TxSync, it indicates success of CheckTx. Does not guarantee inclusion of a block
	// For TxAsync, it always returns true
	// For TxCommit, it indicates the success of both CheckTx and DeliverTx. If the broadcast fails is false.
	OK     bool   `json:"ok"`
	Height *int64 `json:"height,omitempty"`
	Log    string `json:"log"`
}

func (reply *BroadcastReply) FromResultBroadcastTx(result *ctypes.ResultBroadcastTx) {
	reply.TxHash = result.Hash
	reply.OK = result.Code == 0
	reply.Height = nil
	reply.Log = result.Log
}

func (svc Service) InternalBroadcast(request InternalBroadcastRequest, reply *BroadcastReply) error {
	if err := svc.allowedTx(request.RawTx); err != nil {
		return err
	}

	priKey := svc.nodeCtx.PrivVal()
	h, err := priKey.GetHandler()
	if err != nil {
		return errors.Wrap(err, "wrong node private validator key")
	}
	signed, err := h.Sign(request.RawTx.RawBytes())
	if err != nil {
		return errors.Wrap(err, "signing failed")
	}
	rawSignedTx := action.SignedTx{
		RawTx: request.RawTx,
		Signatures: []action.Signature{{
			Signer: h.PubKey(),
			Signed: signed,
		}},
	}
	result, err := svc.tmrpc.BroadcastTxSync(rawSignedTx.SignedBytes())
	if err != nil {
		return errors.Wrap(err, "error broadcast sync")
	}

	reply.FromResultBroadcastTx(result)
	return nil

}

//^TODO Replace error with InternalBroadcastStatus
func BroadcastReportFinalityETHTx(ethCtx *JobsContext, trackerName ethereum.TrackerName, jobID string, success bool) error {

	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.QueryAllStores(trackerName)
	if err != nil {
		return err
	}
	if tracker.State == ethereum2.Released || tracker.State == ethereum2.Failed {
		return nil
	}
	index, voted := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
	if voted {
		//Validator has already Voted
		return nil
	}
	reportFailed := &eth.ReportFinality{
		TrackerName:      trackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
		VoteIndex:        index,
		Success:          success,
	}

	txData, err := reportFailed.Marshal()
	if err != nil {
		ethCtx.Logger.Error("Error while preparing mint txn ", jobID, err)
		return err
	}
	uuidNew, _ := uuid.NewUUID()
	internalFailedTx := action.RawTx{
		Type: action.ETH_REPORT_FINALITY_MINT,
		Data: txData,
		Fee:  action.Fee{},
		Memo: jobID + uuidNew.String(),
	}

	req := InternalBroadcastRequest{
		RawTx: internalFailedTx,
	}
	rep := BroadcastReply{}
	err = ethCtx.Service.InternalBroadcast(req, &rep)

	if err != nil || !rep.OK {
		ethCtx.Logger.Error("Error while broadcasting vote to Fail transaction ", jobID, err, rep.Log)
		return err
	}
	return nil
}

func BroadcastGovExpireVotesTx(jobCtx *JobsContext, proposalID governance.ProposalID, jobID string) error {

	expireVotes := &gov_action.ExpireVotes{
		ProposalID:       proposalID,
		ValidatorAddress: jobCtx.ValidatorAddress,
	}

	txData, err := expireVotes.Marshal()
	if err != nil {
		jobCtx.Logger.Error("Error while preparing expire votes txn ", jobID, err)
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	internalExpireVotesTx := action.RawTx{
		Type: action.EXPIRE_VOTES,
		Data: txData,
		Fee:  action.Fee{},
		Memo: jobID + uuidNew.String(),
	}

	req := InternalBroadcastRequest{
		RawTx: internalExpireVotesTx,
	}
	rep := BroadcastReply{}
	err = jobCtx.Service.InternalBroadcast(req, &rep)

	if err != nil || !rep.OK {
		jobCtx.Logger.Error("Error while broadcasting expire votes transaction ", jobID, err, rep.Log)
		return err
	}
	return nil
}

func BroadcastGovFinalizeVotesTx(jobCtx *JobsContext, proposalID governance.ProposalID, jobID string) error {

	finalizeProposal := &gov_action.FinalizeProposal{
		ProposalID:       proposalID,
		ValidatorAddress: jobCtx.ValidatorAddress,
	}

	txData, err := finalizeProposal.Marshal()
	if err != nil {
		jobCtx.Logger.Error("Error while preparing finalizing votes txn ", jobID, err)
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	internalFinalizeTx := action.RawTx{
		Type: action.PROPOSAL_FINALIZE,
		Data: txData,
		Fee:  action.Fee{},
		Memo: jobID + uuidNew.String(),
	}
	req := InternalBroadcastRequest{
		RawTx: internalFinalizeTx,
	}
	rep := BroadcastReply{}
	err = jobCtx.Service.InternalBroadcast(req, &rep)

	if err != nil || !rep.OK {
		jobCtx.Logger.Error("Error while broadcasting finalizing votes transaction ", jobID, err, rep.Log)
		return err
	}
	return nil
}
