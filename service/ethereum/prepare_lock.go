package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

// Service Called by Wallet
// Input : Lock Request : SIGNED ETH TRANSACTON + OTHER LOCKING PARAMETERS
// OUTPUT :UNSIGNED OLT TRANSANTION
// This function might create a loophope that node owner might use the sgned eth lock tx for his own benificiary address.
// DONT USE IN PRODUCTION
func (svc *Service) CreateRawExtLock(req OLTLockRequest, out *OLTLockReply) error {
	tx := &types.Transaction{}
	err := rlp.DecodeBytes(req.RawTx, tx)
	if err != nil {
		return errors.Wrap(err, "failed to decode provided transaction bytes")
	}

	packets, err := createRawLock(req.Address, req.RawTx, req.Fee, req.Gas)
	if err != nil {
		return errors.Wrap(err, "createRawLock")
	}
	fmt.Println("CreateRawExtLock:", packets)

	*out = OLTLockReply{
		RawTX: packets,
	}
	return nil
}

// Helper Function to create Lock ,and send back unsigned OLT transaction
// Data Field is Lock struct (Tx.data.ETHTxn)

func createRawLock(locker action.Address, rawTx []byte, userfee action.Amount, gas int64) ([]byte, error) {
	// First accept the rawTx
	//tracker := tracker.NewTracker(common.BytesToHash(rawTx))
	lock := eth.Lock{
		Locker:      locker,
		TrackerName: common.BytesToHash(rawTx),
		ETHTxn:      rawTx,
	}

	data, err := lock.Marshal()
	if err != nil {
		return nil, err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{userfee, gas}
	tx := &action.RawTx{
		Type: action.ETH_LOCK,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return nil, codes.ErrSerialization
	}
	return packet, nil
}

// Expects public key , and creates an unsigned TX to send to wallet .
// Wallet signs and then calls onlinelock
func (svc *Service) GetRawLockTX(req ETHLockRequest, out *ETHLockRawTX) error {
	opt := svc.trackerStore.GetOption()
	cd, err := ethereum.NewEthereumChainDriver(svc.config, svc.logger, opt)
	if err != nil {
		return errors.Wrap(err, "GetRawLockTx")
	}
	rawTx, err := cd.PrepareUnsignedETHLock(req.PublicKey, req.Amount)
	if err != nil {
		return err
	}
	*out = ETHLockRawTX{UnsignedRawTx: rawTx}
	return nil
}
