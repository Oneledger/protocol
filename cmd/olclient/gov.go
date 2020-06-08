/*
	Copyright 2017-2018 OneLedger

	Cli to interact with a with the chain.
*/
package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	accounts2 "github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

// Arguments to vote a proposal
type VoteArguments struct {
	ProposalID string `json:"proposalID"`
	Opinion    string `json:"opinion"`
	Address    []byte `json:"address"`
	Password   string `json:"password"`
	GasPrice   string `json:"gasPrice"`
	Gas        int64  `json:"gas"`
}

// Arguments to list proposals
type ListProposalArguments struct {
	ProposalID string `json:"proposalID"`
	State      string `json:"state"`
	Proposer   []byte `json:"proposer"`
	Type       string `json:"type"`
}

var (
	govCmd = &cobra.Command{
		Use:   "gov",
		Short: "governance",
		Long:  "handling governance",
	}

	voteProposalCmd = &cobra.Command{
		Use:   "vote",
		Short: "Vote a proposal",
		RunE:  voteProposal,
	}

	listProposalCmd = &cobra.Command{
		Use:   "list",
		Short: "list proposals",
		RunE:  listProposals,
	}

	voteArgs = &VoteArguments{}
	listArgs = &ListProposalArguments{}
)

func (args *VoteArguments) ClientRequest(currencies *balance.CurrencySet) (client.VoteProposalRequest, error) {
	olt, _ := currencies.GetCurrencyByName("OLT")

	_, err := strconv.ParseFloat(args.GasPrice, 64)
	if err != nil {
		return client.VoteProposalRequest{}, err
	}
	amt := olt.NewCoinFromString(padZero(args.GasPrice)).Amount

	return client.VoteProposalRequest{
		ProposalId: args.ProposalID,
		Address:    args.Address,
		Opinion:    args.Opinion,
		GasPrice:   action.Amount{Currency: "OLT", Value: *amt},
		Gas:        args.Gas,
	}, nil
}

func setVoteArgs() {
	// Transaction Parameters
	voteProposalCmd.Flags().StringVar(&voteArgs.ProposalID, "id", "", "proposal ID")
	voteProposalCmd.Flags().BytesHexVar(&voteArgs.Address, "address", []byte{}, "validator address")
	voteProposalCmd.Flags().StringVar(&voteArgs.Opinion, "opinion", "YES", "vote opinion, YES / NO / GIVEUP")
	voteProposalCmd.Flags().StringVar(&voteArgs.Password, "password", "", "password to access secure wallet.")
	voteProposalCmd.Flags().StringVar(&voteArgs.GasPrice, "gasprice", "0", "include a gas price in OLT")
	voteProposalCmd.Flags().Int64Var(&voteArgs.Gas, "gas", 40000, "gas limit")
}

func setListArgs() {
	// Transaction Parameters
	listProposalCmd.Flags().StringVar(&listArgs.ProposalID, "id", "", "proposal ID")
	listProposalCmd.Flags().StringVar(&listArgs.State, "state", "", "proposal state, active / passed / failed")
	listProposalCmd.Flags().BytesHexVar(&listArgs.Proposer, "proposer", []byte{}, "proposer address")
	listProposalCmd.Flags().StringVar(&listArgs.Type, "type", "", "proposal type, codechange / configupdate / general")
}

func init() {
	RootCmd.AddCommand(govCmd)
	govCmd.AddCommand(voteProposalCmd)
	govCmd.AddCommand(listProposalCmd)
	setVoteArgs()
	setListArgs()
}

func voteProposal(cmd *cobra.Command, args []string) error {
	ctx := NewContext()

	// Prompt for password
	if len(voteArgs.Password) == 0 {
		voteArgs.Password = PromptForPassword()
	}

	// Create new Wallet and User Address
	wallet, err := accounts2.NewWalletKeyStore(keyStorePath)
	if err != nil {
		ctx.logger.Error("failed to create secure wallet", err)
		return err
	}

	// Verify User Password
	usrAddress := keys.Address(voteArgs.Address)
	authenticated, err := wallet.VerifyPassphrase(usrAddress, voteArgs.Password)
	if !authenticated {
		ctx.logger.Error("authentication error", err)
		return err
	}

	// Create request
	fullnode := ctx.clCtx.FullNodeClient()
	currencies, err := fullnode.ListCurrencies()
	if err != nil {
		ctx.logger.Error("failed to get currencies", err)
		return err
	}
	req, err := voteArgs.ClientRequest(currencies.Currencies.GetCurrencySet())
	if err != nil {
		ctx.logger.Error("failed to get request", err)
		return err
	}

	// Validator sign Tx
	reply, err := fullnode.VoteProposal(req)
	if err != nil {
		ctx.logger.Error("failed to create vote request", err)
		return err
	}
	sigV := action.Signature{Signer: reply.Signature.Signer, Signed: reply.Signature.Signed}

	// Deserialize RawTx
	signedTx := &action.SignedTx{}
	err = serialize.GetSerializer(serialize.NETWORK).Deserialize(reply.RawTx, &signedTx.RawTx)
	if err != nil {
		return errors.New("error de-serializing rawTx")
	}

	// Open secure wallet
	if !wallet.Open(usrAddress, voteArgs.Password) {
		ctx.logger.Error("failed to open secure wallet")
		return errors.New("failed to open secure wallet")
	}

	// Sign Tx using secure wallet
	pub, signed, err := wallet.SignWithAddress(reply.RawTx, usrAddress)
	wallet.Close()
	if err != nil {
		ctx.logger.Error("error signing transaction", err)
		return err
	}
	sigW := action.Signature{Signer: pub, Signed: signed}

	// Organize signatures and packet
	signedTx.Signatures = append(signedTx.Signatures, sigW, sigV)
	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(signedTx)
	if packet == nil || err != nil {
		return errors.New("error serializing packet: " + err.Error())
	}

	// Broadcast Tx
	result, err := ctx.clCtx.BroadcastTxCommit(packet)
	if err != nil {
		ctx.logger.Error("error in BroadcastTxCommit", err)
	}

	BroadcastStatus(ctx, result)

	return nil
}

// List proposal information
func listProposals(cmd *cobra.Command, args []string) error {
	Ctx := NewContext()
	fullnode := Ctx.clCtx.FullNodeClient()

	// List single proposal
	if listArgs.ProposalID != "" {
		req := client.ListProposalRequest{
			ProposalId: governance.ProposalID(listArgs.ProposalID),
		}

		reply, err := fullnode.ListProposal(req)
		if err != nil {
			return errors.New("error in getting proposal")
		}

		printProposal(reply.Proposal, reply.Fund, &reply.Vote)
		fmt.Println("Height: ", reply.Height)
	} else { // List multiple proposals
		req := client.ListProposalsRequest{
			State:        listArgs.State,
			Proposer:     listArgs.Proposer,
			ProposalType: listArgs.Type,
		}

		reply, err := fullnode.ListProposals(req)
		if err != nil {
			return errors.New("error in getting proposals")
		}

		if len(reply.Proposals) == 0 {
			return nil
		}

		for i, p := range reply.Proposals {
			printProposal(p, reply.ProposalFunds[i], &reply.ProposalVotes[i])
		}
		fmt.Println("Height: ", reply.Height)
	}
	return nil
}

func printProposal(p governance.Proposal, funds balance.Amount, stat *governance.VoteStatus) {
	fmt.Println("ProposalID : ", p.ProposalID)
	fmt.Println("Type       : ", p.Type.String())
	fmt.Println("Status     : ", p.Status.String())
	fmt.Println("Outcome    : ", p.Outcome.String())
	fmt.Println("Description: ", p.Description)
	fmt.Println("Proposer   : ", p.Proposer.Humanize())
	fmt.Println("Funding Deadline: ", p.FundingDeadline)
	fmt.Println("Funding Goal    : ", p.FundingGoal)
	fmt.Println("Voting Deadline : ", p.VotingDeadline)
	fmt.Println("Pass Percentage : ", p.PassPercentage, "%")
	fmt.Println("Current Funds   : ", funds.String())
	if p.Status == governance.ProposalStatusVoting {
		fmt.Println("Power YES  : ", stat.PowerYes)
		fmt.Println("Power NO   : ", stat.PowerNo)
		fmt.Println("Power All  : ", stat.PowerAll)
		fmt.Println("Vote Result: ", stat.Result.String())
	}
	fmt.Println()
}
