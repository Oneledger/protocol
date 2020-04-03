package ethereum

import (
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

func (svc *Service) CreateRawExtLock(req OLTLockRequest, out *OLTReply) error {

	packets, err := createRawLock(req.Address, req.RawTx, req.Fee, req.Gas)
	if err != nil {
		svc.logger.Error(err, codes.ErrPreparingOLTLock.ErrorMsg())
		return codes.ErrPreparingOLTLock
	}

	*out = OLTReply{
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
		Locker: locker,
		ETHTxn: rawTx,
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

// Expects users ethereum address , and creates an unsigned TX to send to wallet .
// Wallet signs and then calls onlinelock
func (svc *Service) GetRawLockTX(req ETHLockRequest, out *ETHLockRawTX) error {
	opt := svc.trackers.GetOption()
	cd, err := ethereum.NewChainDriver(svc.config, svc.logger, opt.ContractAddress, opt.ContractABI, ethereum.ERC)
	if err != nil {
		return errors.Wrap(err, "GetRawLockTx")
	}
	// TODO:Change to address
	rawTx, err := cd.PrepareUnsignedETHLock(req.UserAddress, req.Amount)
	if err != nil {
		svc.logger.Error(codes.ErrPreparingETHLock.Msg)
		return codes.ErrPreparingETHLock
	}
	*out = ETHLockRawTX{UnsignedRawTx: rawTx}
	return nil
}
