package eth

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
)

type ReportFinalityMint struct {
	TrackerName      ethereum.TrackerName
	Locker           action.Address
	ValidatorAddress action.Address
	VoteIndex         int64
}

var _ action.Msg = &ReportFinalityMint{}

func (m *ReportFinalityMint) Signers() []action.Address {
	return []action.Address{
		m.ValidatorAddress,
	}
}

func (m *ReportFinalityMint) Type() action.Type {
	return action.ETH_REPORT_FINALITY_MINT
}

func (m *ReportFinalityMint) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.ETH_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: m.Locker.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(m.TrackerName.Hex()),
	}
	tag4 := common.KVPair{
		Key:   []byte("tx.validator"),
		Value: m.ValidatorAddress.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (m *ReportFinalityMint) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ReportFinalityMint) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

var _ action.Tx = reportFinalityMintTx{}

type reportFinalityMintTx struct {
}

func (r reportFinalityMintTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	fmt.Println("Starting Validate ")
	f := &ReportFinalityMint{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	if f.VoteIndex < 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (r reportFinalityMintTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	fmt.Println("START CHECK TX CheckFinality")
	//return runCheckFinalityMint(ctx,tx)
	f := &ReportFinalityMint{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		ctx.Logger.Error(err, "err getting tracker")
		//	return false, action.Response{Log: err.Error()}
	}



	//if tracker.State != trackerlib.BusyBroadcasting {
	//	return false, action.Response{Log: errors.New("tracker not available for finalizing").Error()}
	//}
	//
	//if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
	//	return false, action.Response{Log: "tracker process not owned by user"}
	//}
	//
	//if !ctx.Validators.IsValidatorAddress(f.ValidatorAddress) {
	//	return false, action.Response{Log: "transaction sender not a validator"}
	//}
	//
	if tracker.Finalized() {
		return true, action.Response{Log: "tracker already finalized"}
	}


	ctx.Logger.Error("Trying to add vote (CHECK TX)")
	index,ok := tracker.CheckIfVoted(f.ValidatorAddress)
	ctx.Logger.Info("Before voting", ok,"Index :",index ,"F.index" ,f.VoteIndex)
	err = tracker.AddVote(f.ValidatorAddress, f.VoteIndex)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to add vote").Error()}
	}
	fmt.Printf("%b \n", tracker.FinalityVotes)

	index,ok = tracker.CheckIfVoted(f.ValidatorAddress)

	ctx.Logger.Info("After voting (CHECK TX)", ok ,"Index :",index)
	ctx.Logger.Info("Vote Count (CHECK TX): ",tracker.GetVotes())

	//if tracker.State != trackerlib.BusyFinalizing {tracker.State = trackerlib.BusyFinalizing}


	ctx.Logger.Info("IS finalaized (CHECK TX) :"  ,tracker.Finalized())
	ctx.Logger.Info("Tracker Votes (CHECK TX): " ,tracker.GetVotes())

	if tracker.Finalized() {


		err = mintTokens(ctx, tracker, *f)
		if err !=nil {
			return false, action.Response{Log:errors.Wrap(err,"UNABLE TO MINT TOKENS").Error()}
		}


		return true, action.Response{Log: "MINTING SUCCESSFUL"}
	}

	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		ctx.Logger.Info("Unable to save the tracker",err)
		return false, action.Response{Log:errors.Wrap(err,"unable to save the tracker").Error()}
	}
	fmt.Println("TRACKER SAVED AT CHECK FINALITY (VOTES) (CHECK TX): " , tracker.GetVotes())
	ctx.Logger.Info("(CHECK TX) COMPLETED")
	return true, action.Response{Log: "vote success, not ready to mint: "+ strconv.Itoa(tracker.GetVotes())}

}

func (r reportFinalityMintTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	fmt.Println("START DELIVER TX CheckFinality")
	fmt.Println("Starting runCheck Finality Mint Internal Trasaction")
	f := &ReportFinalityMint{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		ctx.Logger.Error(err, "err getting tracker")
		//	return false, action.Response{Log: err.Error()}
	}



	//if tracker.State != trackerlib.BusyBroadcasting {
	//	return false, action.Response{Log: errors.New("tracker not available for finalizing").Error()}
	//}
	//
	//if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
	//	return false, action.Response{Log: "tracker process not owned by user"}
	//}
	//
	//if !ctx.Validators.IsValidatorAddress(f.ValidatorAddress) {
	//	return false, action.Response{Log: "transaction sender not a validator"}
	//}
	//
	if tracker.Finalized() {
		return true, action.Response{Log: "tracker already finalized"}
	}


	ctx.Logger.Error("Trying to add vote (DELIVER TX)")
	index,ok := tracker.CheckIfVoted(f.ValidatorAddress)
	ctx.Logger.Info("Before voting", ok,"Index :",index ,"F.index" ,f.VoteIndex)
	err = tracker.AddVote(f.ValidatorAddress, f.VoteIndex)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to add vote").Error()}
	}
	fmt.Printf("%b \n", tracker.FinalityVotes)

	index,ok = tracker.CheckIfVoted(f.ValidatorAddress)

	ctx.Logger.Info("After voting (DELIVER TX)", ok ,"Index :",index)
	ctx.Logger.Info("Vote Count (DELIVER TX): ",tracker.GetVotes())

	//if tracker.State != trackerlib.BusyFinalizing {tracker.State = trackerlib.BusyFinalizing}


	ctx.Logger.Info("IS finalaized (DELIVER TX) :"  ,tracker.Finalized())
	ctx.Logger.Info("Tracker Votes (DELIVER TX) : " ,tracker.GetVotes())

	if tracker.Finalized() {


		err = mintTokens(ctx, tracker, *f)
		if err !=nil {
			return false, action.Response{Log:errors.Wrap(err,"UNABLE TO MINT TOKENS").Error()}
		}


		return true, action.Response{Log: "MINTING SUCCESSFUL"}
	}

	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		ctx.Logger.Info("Unable to save the tracker",err)
		return false, action.Response{Log:errors.Wrap(err,"unable to save the tracker").Error()}
	}
	fmt.Println("TRACKER SAVED AT CHECK FINALITY (VOTES) (DELIVER TX): " , tracker.GetVotes())
	ctx.Logger.Info("(DELIVER TX) COMPLETED")
	return true, action.Response{Log: "vote success, not ready to mint: "+ strconv.Itoa(tracker.GetVotes())}
    //return runCheckFinalityMint(ctx,tx)

}



func runCheckFinalityMint(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	fmt.Println("Starting runCheck Finality Mint Internal Trasaction")
	f := &ReportFinalityMint{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		ctx.Logger.Error(err, "err getting tracker")
	//	return false, action.Response{Log: err.Error()}
	}



	//if tracker.State != trackerlib.BusyBroadcasting {
	//	return false, action.Response{Log: errors.New("tracker not available for finalizing").Error()}
	//}
	//
	//if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
	//	return false, action.Response{Log: "tracker process not owned by user"}
	//}
	//
	//if !ctx.Validators.IsValidatorAddress(f.ValidatorAddress) {
	//	return false, action.Response{Log: "transaction sender not a validator"}
	//}
	//
	if tracker.Finalized() {
		return true, action.Response{Log: "tracker already finalized"}
	}


	ctx.Logger.Error("Trying to add vote ")
	index,ok := tracker.CheckIfVoted(f.ValidatorAddress)
	ctx.Logger.Info("Before voting", ok,"Index :",index ,"F.index" ,f.VoteIndex)
	err = tracker.AddVote(f.ValidatorAddress, f.VoteIndex)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to add vote").Error()}
	}
	fmt.Printf("%b \n", tracker.FinalityVotes)

	index,ok = tracker.CheckIfVoted(f.ValidatorAddress)

	ctx.Logger.Info("After voting", ok ,"Index :",index)
	ctx.Logger.Info("Vote Count : ",tracker.GetVotes())

	//if tracker.State != trackerlib.BusyFinalizing {tracker.State = trackerlib.BusyFinalizing}


	ctx.Logger.Info("IS finalaized  :"  ,tracker.Finalized())
	ctx.Logger.Info("Tracker Votes  : " ,tracker.GetVotes())

	if tracker.Finalized() {


		err = mintTokens(ctx, tracker, *f)
		if err !=nil {
			return false, action.Response{Log:errors.Wrap(err,"UNABLE TO MINT TOKENS").Error()}
		}


		return true, action.Response{Log: "MINTING SUCCESSFUL"}
	}

	err = ctx.ETHTrackers.Set(tracker)
	if err != nil {
		ctx.Logger.Info("Unable to save the tracker",err)
		return false, action.Response{Log:errors.Wrap(err,"unable to save the tracker").Error()}
	}
    fmt.Println("TRACKER SAVED AT CHECK FINALITY (VOTES): " , tracker.GetVotes())
	ctx.Logger.Info("Voting Done ,unable to mint yet")
	return true, action.Response{Log: "vote success, not ready to mint: "+ strconv.Itoa(tracker.GetVotes())}
}

func (reportFinalityMintTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	// return action.BasicFeeHandling(ctx, signedTx, start, size, 1)

	return true, action.Response{}
}

func mintTokens(ctx *action.Context, tracker *trackerlib.Tracker, oltTx ReportFinalityMint) error {
	curr, ok := ctx.Currencies.GetCurrencyByName("ETH")
	if !ok {
		return errors.New("ETH currency not allowed")
	}
	lockAmount,err := GetAmount(tracker)
	if err != nil {
		return err
	}

	//tracker.State = trackerlib.Minted
	//err = ctx.ETHTrackers.Set(tracker)
	//if err != nil {
	//	return err
	//}

	oEthCoin := curr.NewCoinFromAmount(balance.Amount(*lockAmount))
	err = ctx.Balances.AddToAddress(oltTx.Locker, oEthCoin)
	if err != nil {
		ctx.Logger.Error(err)
		return errors.New("Unable to mint")
	}

	return nil
}
