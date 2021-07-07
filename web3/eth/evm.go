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

	ethAcc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(height, address.Bytes())
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
	svc.logger.Debugf("eth_call args %+v height %d", call, height)

	result, err := svc.callContract(call, svc.getStateHeight(height))
	if err != nil {
		return hexutil.Bytes{}, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return hexutil.Bytes{}, rpctypes.NewRevertError(result)
	}
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
	// TODO: Add versioning support
	stateDB := action.NewCommitStateDB(svc.ctx.GetContractStore(), svc.ctx.GetAccountKeeper(), svc.logger)
	bhash := stateDB.GetHeightHash(uint64(height))
	stateDB.SetHeightHash(uint64(height), bhash, false)

	block := svc.ctx.GetAPI().Block(height).Block
	header := &abci.Header{
		ChainID: block.ChainID,
		Height:  block.Height,
		Time:    block.Time,
	}

	var from keys.Address = call.From.Bytes()
	var to keys.Address
	if call.To != nil {
		to = call.To.Bytes()
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

	// TODO: Change on account nonce with involving of pending transactions
	ethAcc, _ := svc.ctx.GetAccountKeeper().GetVersionedAccount(height, call.From.Bytes())
	nonce := ethAcc.Sequence

	price := action.Amount{
		Currency: "OLT",
		Value:    balance.Amount(*big.NewInt(gasPrice)),
	}

	// TODO: Change this when the new sig mechanics will be implemented
	tx := action.RawTx{
		Type: action.SC_EXECUTE,
		Fee: action.Fee{
			Price: price,
			Gas:   gas,
		},
	}
	evmTx := action.NewEVMTransaction(stateDB, header, from, &to, nonce, value, data, true)

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
