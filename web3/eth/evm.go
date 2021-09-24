package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/vm"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	ethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/rpc"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (svc *Service) GetStorageAt(address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		svc.logger.Debug("eth_getStorageAt", "block err", err)
		return hexutil.Bytes{}, nil
	}
	height = svc.getStateHeight(height)

	svc.logger.Debug("eth_getStorageAt", "address", address, "key", key, "height", height)

	prefixKey := utils.GetStorageByAddressKey(address, common.HexToHash(key).Bytes())
	prefixStore := evm.AddressStoragePrefix(address)

	cs := svc.ctx.GetContractStore()

	value := common.Hash{}
	storeKey := cs.GetStoreKey(prefixStore, prefixKey.Bytes())

	rawValue, err := cs.State.GetAtHeight(height, storeKey)
	svc.logger.Debug("eth_getStorageAt", "address", address, "err", err)
	if len(rawValue) > 0 {
		value.SetBytes(rawValue)
	}
	return value[:], nil
}

func (svc *Service) GetCode(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		svc.logger.Debug("eth_getCode", "block err", err)
		return hexutil.Bytes{}, nil
	}

	height = svc.getStateHeight(height)

	svc.logger.Debug("eth_getCode", "address", address, "height", height)

	stateDB := svc.GetStateDB()
	keeper := stateDB.GetAccountKeeper()
	cs := stateDB.GetContractStore()

	ethAcc, err := keeper.GetVersionedAccount(address.Bytes(), height)
	if err != nil {
		return hexutil.Bytes{}, nil
	}
	svc.logger.Debug("eth_getCode", "address", address, "code hash", ethAcc.CodeHash)
	storeKey := cs.GetStoreKey(evm.KeyPrefixCode, ethAcc.CodeHash)
	code, err := cs.State.GetAtHeight(height, storeKey)
	svc.logger.Debug("eth_getCode", "address", address, "err", err)
	return code[:], nil
}

func (svc *Service) Call(call rpctypes.CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		svc.logger.Debug("eth_call", "block err", err)
		return nil, nil
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
	svc.logger.Debug("eth_estimateGas", "gas", result.UsedGas)
	return hexutil.Uint64(result.UsedGas), err
}

func (svc *Service) callContract(call rpctypes.CallArgs, height int64) (*vm.ExecutionResult, error) {
	var (
		blockNum  uint64
		isPending bool
	)
	switch height {
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
	default:
		blockNum = uint64(svc.getState().Version())
		isPending = true
	}

	stateDB := svc.GetStateDB()
	bhash := stateDB.GetHeightHash(blockNum)
	stateDB.SetHeightHash(blockNum, bhash)

	resBlock := svc.ctx.GetAPI().Block(int64(blockNum))
	if resBlock == nil {
		return nil, errors.New("failed to get block")
	}
	block := resBlock.Block
	header := &abci.Header{
		ChainID: block.ChainID,
		Height:  block.Height,
		Time:    block.Time,
	}

	var from keys.Address
	if call.From != nil {
		from = call.From.Bytes()
	}

	var to *keys.Address
	if call.To != nil {
		to = new(keys.Address)
		*to = call.To.Bytes()
	}

	var gasPrice *big.Int = vm.DefaultGasPrice
	if call.GasPrice != nil {
		gasPrice = call.GasPrice.ToInt()
	}

	var gas uint64 = vm.SimulationBlockGasLimit
	if call.Gas != nil {
		gas = uint64(*call.Gas)
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

	// Set infinite balance to the fake caller account.
	fromAcc := stateDB.GetOrNewStateObject(common.BytesToAddress(from))
	fromAcc.SetBalance(math.MaxBig256)

	evmTx := vm.NewEVMTransaction(stateDB, new(ethcore.GasPool).AddGas(math.MaxUint64), header, from, to, nonce, value, data, nil, gas, gasPrice, true)

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

	result, err := evmTx.Apply()
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
