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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	web3types "github.com/Oneledger/protocol/web3/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (svc *Service) GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	contractStore := svc.stateDB.GetContractStore()

	height, err := svc.stateAndHeaderByNumberOrHash(blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	svc.log.Debug("eth_getStorageAt", "address", address, "key", key, "height", height)

	prefixKey := utils.GetStorageByAddressKey(address, common.HexToHash(key).Bytes())
	prefixStore := evm.AddressStoragePrefix(address)

	value := common.Hash{}
	storeKey := contractStore.GetStoreKey(prefixStore, prefixKey.Bytes())
	rawValue := contractStore.State.GetVersioned(height, storeKey)
	if len(rawValue) > 0 {
		value.SetBytes(rawValue)
	}
	return value[:], nil
}

func (svc *Service) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	accountKeeper := svc.stateDB.GetAccountKeeper()
	contractStore := svc.stateDB.GetContractStore()

	height, err := svc.stateAndHeaderByNumberOrHash(blockNrOrHash)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	svc.log.Debug("eth_getCode", "address", address, "height", height)

	ethAcc, err := accountKeeper.GetVersionedAccount(height, address.Bytes())
	if err != nil {
		return hexutil.Bytes{}, nil
	}
	storeKey := contractStore.GetStoreKey(evm.KeyPrefixCode, ethAcc.CodeHash)
	code := contractStore.State.GetVersioned(height, storeKey)
	return code[:], nil
}

func (svc *Service) Call(call web3types.CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := svc.stateAndHeaderByNumberOrHash(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	svc.log.Debug("eth_call", "args", call, "height", height)

	result, err := svc.callContract(call, height)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return hexutil.Bytes{}, newRevertError(result)
	}
	return result.Return(), result.Err
}

func (svc *Service) EstimateGas(call web3types.CallArgs) (hexutil.Uint64, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.log.Debug("eth_estimateGas", "args", call)
	height := svc.getState().Version()
	result, err := svc.callContract(call, height)
	if err != nil {
		return 0, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return 0, newRevertError(result)
	}
	// TODO: change 1000 buffer for more accurate buffer
	gas := result.UsedGas + 1000
	return hexutil.Uint64(gas), result.Err
}

func (svc *Service) callContract(call web3types.CallArgs, height int64) (*action.ExecutionResult, error) {
	accountKeeper := svc.stateDB.GetAccountKeeper()
	contractStore := svc.stateDB.GetContractStore()

	// TODO: Add versioning support
	stateDB := action.NewCommitStateDB(contractStore, accountKeeper, svc.log)
	bhash := stateDB.GetHeightHash(uint64(height))
	stateDB.SetHeightHash(uint64(height), bhash, false)

	block := svc.ext.Block(height).Block
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
	ethAcc, _ := accountKeeper.GetVersionedAccount(height, call.From.Bytes())
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
