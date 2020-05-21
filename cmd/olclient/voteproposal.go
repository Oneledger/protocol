/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
)

type VoteProposalArguments struct {
	ProposalID string `json:"proposalID"`
	Address    []byte `json:"address"`
	Opinion    string `json:"opinion"`
	Fee        string `json:"fee"`
	Gas        int64  `json:"gas"`
}

func (args *VoteProposalArguments) ClientRequest(currencies *balance.CurrencySet) client.VoteProposalRequest {

	return client.VoteProposalRequest{}
}

var voteProposalCmd = &cobra.Command{
	Use:   "voteproposal",
	Short: "Vote a proposal",
	RunE:  voteProposal,
}

var voteProposalArgs = &VoteProposalArguments{}

func setVoteProposalArgs() {
	// Transaction Parameters
	voteProposalCmd.Flags().StringVar(&voteProposalArgs.ProposalID, "proposalID", "", "proposal ID")
	voteProposalCmd.Flags().BytesHexVar(&voteProposalArgs.Address, "address", []byte{}, "validator address")
	voteProposalCmd.Flags().StringVar(&voteProposalArgs.Opinion, "opinion", "YES", "vote opinion, YES / NO / GIVEUP")
}

func init() {
	RootCmd.AddCommand(voteProposalCmd)
	setVoteProposalArgs()
}

// IssueRequest sends out a sendTx to all of the nodes in the chain
func voteProposal(cmd *cobra.Command, args []string) error {
	ctx := NewContext()

	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}
	// Create message
	req, err := sendfundsargs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}
	fmt.Println(req)
	reply, err := fullnode.SendTx(req)
	if err != nil {
		ctx.logger.Error("failed to create SendTx", err)
		return err
	}
	packet := reply.RawTx

	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
}
