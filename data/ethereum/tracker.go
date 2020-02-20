package ethereum

import (
	//"errors"

	"strconv"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils/transition"
)

type TrackerState int

// Tracker
type Tracker struct {
	// State tracks the current state of the tracker, Also used for locking distributed access
	Type          ProcessType
	State         TrackerState
	TrackerName   ethereum.TrackerName
	SignedETHTx   []byte
	Validators    []keys.Address
	ProcessOwner  keys.Address
	FinalityVotes []Vote
	To            keys.Address
}

//number of validator should be smaller than 64
func NewTracker(typ ProcessType, owner keys.Address, signedEthTx []byte, name ethereum.TrackerName, validators []keys.Address) *Tracker {

	return &Tracker{
		Type:          typ,
		State:         New,
		TrackerName:   name,
		ProcessOwner:  owner,
		SignedETHTx:   signedEthTx,
		Validators:    validators,
		FinalityVotes: make([]Vote, len(validators)),
	}
}

func (t *Tracker) AddVote(addr keys.Address, index int64, vote bool) error {

	if len(t.Validators) <= int(index) {
		return errTrackerInvalidVote
	}

	_, voted := t.CheckIfVoted(addr)
	if voted {
		return errTrackerInvalidVote
	}
	if t.Validators[index].Equal(addr) {
		//		t.FinalityVotes = (t.FinalityVotes | (1 << index))
		//		return nil

		// vote:
		//      0: not vote yet
		//      1: vote for yes
		//      2: vote for no
		v := int8(0)
		if !vote {
			v = 1
		}
		t.FinalityVotes[index] = Vote(v + 1)

	}

	return nil
}



func (t *Tracker) GetJobID(state TrackerState) string {
	return t.TrackerName.String() + storage.DB_PREFIX + strconv.Itoa(int(state))
}

func (t *Tracker) GetVotes() (yes, no int) {
	ycnt := 0
	ncnt := 0
	for _, item := range t.FinalityVotes {
		if item == 1 {
			ycnt++
		}
		if item == 2 {
			ncnt++
		}
	}

	return ycnt, ncnt
}

func (t *Tracker) CheckIfVoted(node keys.Address) (index int64, voted bool) {
	index = int64(-1)
	for i, addr := range t.Validators {
		if addr.Equal(node) {
			index = int64(i)
			break
		}
	}
	//if index < int64(len(t.Validators)) && index >= 0 {
	//	var mybit uint64 = 1 << index
	//	and := t.FinalityVotes & mybit
	//	v = and == mybit
	//}
	if index == -1 {
		return -1, false
	}
	if t.FinalityVotes[index] > 0 {
		return index, true
	}

	return index, false
}


func (t *Tracker) Finalized() bool {
	l := len(t.Validators)
	num := (l * 2 / 3) + 1
	//v := t.FinalityVotes
	//cnt := 0
	//for v >= 1 {
	//	if v%2 == 1 {
	//		cnt++
	//	}
	//	v = v >> 1
	//}
	y, _ := t.GetVotes()
	return y >= num
}


func (t *Tracker) Failed() bool {
	l := len(t.Validators)
	num := (l * 2 / 3) + 1
	//v := t.FinalityVotes
	//cnt := 0
	//for v >= 1 {
	//	if v%2 == 1 {
	//		cnt++
	//	}
	//	v = v >> 1
	//}
	_, n := t.GetVotes()
	return n >= num
}

func (t Tracker) NextStep() string {
	if t.Type == ProcessTypeLock {
		switch t.State {
		case New:
			return BROADCASTING
		case BusyBroadcasting:
			return FINALIZING
		case BusyFinalizing:
			return FINALIZE
		case Finalized:
			return MINTING
		case Released:
			return CLEANUP
		case Failed:
			return CLEANUPFAILED
		}
		return transition.NOOP
	}
	if t.Type == ProcessTypeLockERC {
		switch t.State {
		case New:
			return BROADCASTING
		case BusyBroadcasting:
			return FINALIZING
		case BusyFinalizing:
			return FINALIZE
		case Finalized:
			return MINTING
		case Released:
			return CLEANUP
		case Failed:
			return CLEANUPFAILED
		}
		return transition.NOOP
	}
	if t.Type == ProcessTypeRedeem {
		switch t.State {
		case New:
			return SIGNING
		case BusyBroadcasting:
			return VERIFYREDEEM
		case BusyFinalizing:
			return REDEEMCONFIRM
		case Finalized:
			return BURN
		case Released:
			return CLEANUP
		case Failed:
			return CLEANUPFAILED
		}
	}
	if t.Type == ProcessTypeRedeemERC {
		switch t.State {
		case New:
			return SIGNING
		case BusyBroadcasting:
			return VERIFYREDEEM
		case BusyFinalizing:
			return REDEEMCONFIRM
		case Finalized:
			return BURN
		case Released:
			return CLEANUP
		case Failed:
			return CLEANUPFAILED
		}
	}
	return transition.NOOP
}
