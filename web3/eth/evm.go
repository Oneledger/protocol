package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/utils"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (svc *Service) GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	svc.logger.Debug("eth_getStorageAt", "address", address, "key", key, "height", height)

	prefixKey := utils.GetStorageByAddressKey(address, common.HexToHash(key).Bytes())
	prefixStore := evm.AddressStoragePrefix(address)

	value := common.Hash{}
	storeKey := svc.ctx.GetContractStore().GetStoreKey(prefixStore, prefixKey.Bytes())
	rawValue := svc.ctx.GetContractStore().State.GetVersioned(svc.getStateHeight(height), storeKey)
	if len(rawValue) > 0 {
		value.SetBytes(rawValue)
	}
	return value[:], nil
}

func (svc *Service) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	svc.logger.Debug("eth_getCode", "address", address, "height", height)

	ethAcc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(address.Bytes(), height)
	if err != nil {
		return hexutil.Bytes{}, nil
	}
	storeKey := svc.ctx.GetContractStore().GetStoreKey(evm.KeyPrefixCode, ethAcc.CodeHash)
	code := svc.ctx.GetContractStore().State.GetVersioned(svc.getStateHeight(height), storeKey)
	return code[:], nil
}

func (svc *Service) Call(call rpctypes.CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return nil, err
	}
	svc.logger.Debugf("eth_call args data '%s' with height '%d'", call, height)

	result, err := svc.callContract(call, height)
	if err != nil {
		svc.logger.Debug("eth_call", "err", err)
		return hexutil.Bytes{}, err
	}
	svc.logger.Debug("eth_call", "result", result)
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		rErr := rpctypes.NewRevertError(result)
		svc.logger.Debug("eth_call", "revert err", rErr)
		return hexutil.Bytes{}, rErr
	}
	svc.logger.Debug("eth_call", "return", result.Return(), "err", result.Err)
	return result.Return(), result.Err
}

func (svc *Service) EstimateGas(call rpctypes.CallArgs) (hexutil.Uint64, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.logger.Debug("eth_estimateGas", "args", call)
	height := svc.getState().Version()
	result, err := svc.callContract(call, height)
	if err != nil {
		return 0, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return 0, rpctypes.NewRevertError(result)
	}
	// TODO: change 1000 buffer for more accurate buffer
	gas := result.UsedGas + 1000
	return hexutil.Uint64(gas), result.Err
}

func (svc *Service) callContract(call rpctypes.CallArgs, height int64) (*action.ExecutionResult, error) {
	var (
		blockNum  uint64
		isPending bool
	)
	switch height {
	case rpctypes.EarliestBlockNumber:
		blockNum = 1
	default:
		blockNum = uint64(svc.getState().Version())
		isPending = true
	}
	// TODO: Add versioning support
	stateDB := action.NewCommitStateDB(svc.ctx.GetContractStore(), svc.ctx.GetAccountKeeper(), svc.logger)
	bhash := stateDB.GetHeightHash(blockNum)
	stateDB.SetHeightHash(blockNum, bhash, false)

	block := svc.ctx.GetAPI().Block(int64(blockNum)).Block
	header := &abci.Header{
		ChainID: block.ChainID,
		Height:  block.Height,
		Time:    block.Time,
	}

	var from keys.Address = call.From.Bytes()
	var to *keys.Address
	if call.To != nil {
		to = new(keys.Address)
		*to = call.To.Bytes()
	}

	var gasPrice int64 = action.DefaultGasPrice.Int64()
	if call.GasPrice != nil {
		gasPrice = call.GasPrice.ToInt().Int64()
	}

	var gas int64 = int64(action.DefaultGasLimit)
	if call.Gas != nil {
		gas = int64(*call.Gas)
	}

	value := new(big.Int)
	if call.Value != nil {
		value = call.Value.ToInt()
	}

	// Set Data if provided
	var data []byte
	if call.Data != nil {
		data = []byte(*call.Data)
	}

	nonce := svc.getNonce(common.BytesToAddress(from.Bytes()), blockNum, utils.HashToBigInt(block.ChainID), isPending)

	svc.logger.Debug("eth_callContract", "from", from, "to", to, "gasPrice", gasPrice, "gas", gas, "value", value, "data", data, "nonce", nonce)

	price := action.Amount{
		Currency: "OLT",
		Value:    balance.Amount(*big.NewInt(gasPrice)),
	}

	tx := action.RawTx{
		Type: action.NEXUS,
		Fee: action.Fee{
			Price: price,
			Gas:   gas,
		},
	}
	evmTx := action.NewEVMTransaction(stateDB, header, from, to, nonce, value, data, true)

	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	vmenv := evmTx.NewEVM()
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		vmenv.Cancel()
	}()

	result, err := evmTx.Apply(vmenv, tx)
	if vmenv.Cancelled() {
		return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	return result, err
}

func (svc *Service) getNonce(address common.Address, height uint64, chainID *big.Int, isPending bool) uint64 {
	var (
		nonce        uint64
		pendingNonce uint64
	)
	if isPending {
		pendingNonce = rpcutils.GetPendingTxCountByAddress(svc.getTMClient(), address)
	}
	ethAcc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(address.Bytes(), int64(height))
	if err == nil {
		nonce = ethAcc.Sequence
	}
	return nonce + pendingNonce
}
