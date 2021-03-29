package query

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/service/query/filters"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	abci "github.com/tendermint/tendermint/abci/types"
)

// revertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

func newRevertError(result *action.ExecutionResult) *revertError {
	reason, errUnpack := ethabi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

func parseBlockNumber(blockNumber string, defBlock uint64) (*big.Int, error) {
	if blockNumber == "latest" {
		return big.NewInt(int64(defBlock)), nil
	} else {
		b, err := strconv.Atoi(blockNumber)
		if err != nil {
			return big.NewInt(0), fmt.Errorf("Failed to parse block number: %s", err)
		}
		return big.NewInt(int64(b)), nil
	}
}

func parseAddresses(addresses []keys.Address) []ethcmn.Address {
	result := make([]ethcmn.Address, 0, len(addresses))
	for _, addrB := range addresses {
		result = append(result, ethcmn.BytesToAddress(addrB))
	}
	return result
}

func parseTopics(topics [][][]byte) [][]ethcmn.Hash {
	result := make([][]ethcmn.Hash, len(topics))
	for x := range topics {
		for y := range topics[x] {
			result[x][y] = ethcmn.BytesToHash(topics[x][y])
		}
	}
	return result
}

func (svc *Service) EVMTransactionLogs(args client.EVMTransactionLogsRequest, reply *client.EVMLogsReply) error {
	height := big.NewInt(svc.contracts.State.Version())
	addresses := parseAddresses(args.Addresses)
	topics := parseTopics(args.Topics)
	stateDB := action.NewCommitStateDB(svc.contracts, svc.accountKeeper, svc.logger)

	logs, err := stateDB.GetLogs(ethcmn.BytesToHash(args.TransactionHash))
	if err != nil {
		return err
	}

	for _, log := range filters.FilterLogs(logs, height, height, addresses, topics) {
		topicsS := make([]string, 0, len(log.Topics))
		for _, topic := range log.Topics {
			topicsS = append(topicsS, topic.Hex())
		}
		rLog := client.EVMLogReply{
			TransactionHash: log.TxHash.Hex(),
			BlockHeight:     strconv.Itoa(int(log.BlockNumber)),
			BlockHash:       log.BlockHash.Hex(),
			Address:         keys.Address(log.Address.Bytes()),
			Data:            log.Data,
			Topics:          topicsS,
		}
		reply.Logs = append(reply.Logs, rLog)
	}
	return nil
}

func (svc *Service) EVMAccount(args client.EVMAccountRequest, reply *client.EVMAccountReply) error {
	acc, err := svc.accountKeeper.GetAccount(args.Address)
	if err != nil {
		svc.logger.Error("error getting evm account", err)
		return err
	}
	reply.Address = args.Address
	reply.Balance = acc.Balance().String()
	reply.CodeHash = ethcmn.Bytes2Hex(acc.CodeHash)
	reply.Nonce = acc.Sequence
	return nil
}

// EVMCall call smart contract code to get the result
func (svc *Service) EVMCall(args client.SendTxRequest, reply *client.EVMCallReply) error {
	height := svc.contracts.State.Version()
	stateDB := action.NewCommitStateDB(svc.contracts, svc.accountKeeper, svc.logger)
	bhash := stateDB.GetHeightHash(uint64(height))
	stateDB.SetBlockHash(bhash)
	// TODO: Change this
	header := &abci.Header{
		ChainID: "test-1",
		Height:  height,
		Time:    time.Now(),
	}

	var to *keys.Address
	if len(args.To) != 0 {
		to = &args.To
	}
	tx := action.RawTx{
		Type: action.SC_EXECUTE,
		Fee: action.Fee{
			Price: args.GasPrice,
			Gas:   args.Gas,
		},
	}

	evmTx := action.NewEVMTransaction(stateDB, header, args.From, to, args.Amount.Value.BigInt(), args.Data)

	// TODO: Move in some constant
	timeout := 5 * time.Second
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
		return fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}

	if result.Failed() {
		return result.Err
	}

	if err != nil {
		return fmt.Errorf("err: %w (supplied gas %d)", err, args.Gas)
	}

	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return newRevertError(result)
	}

	*reply = client.EVMCallReply{
		Result: result.Return(),
	}
	return nil
}

// EVMEstimateGas call smart contract code to get the gas
func (svc *Service) EVMEstimateGas(args client.SendTxRequest, reply *client.EVMEstimateGasReply) error {
	height := svc.contracts.State.Version()
	stateDB := action.NewCommitStateDB(svc.contracts, svc.accountKeeper, svc.logger)
	bhash := stateDB.GetHeightHash(uint64(height))
	stateDB.SetBlockHash(bhash)
	// TODO: Change this
	header := &abci.Header{
		ChainID: "test-1",
		Height:  height,
		Time:    time.Now(),
	}

	var to *keys.Address
	if len(args.To) != 0 {
		to = &args.To
	}
	tx := action.RawTx{
		Type: action.SC_EXECUTE,
		Fee: action.Fee{
			Price: args.GasPrice,
			Gas:   args.Gas,
		},
	}

	evmTx := action.NewEVMTransaction(stateDB, header, args.From, to, args.Amount.Value.BigInt(), args.Data)

	// TODO: Move in some constant
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	vmenv := evmTx.NewEVM()
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		vmenv.Cancel()
	}()

	gasUsed, err := evmTx.EstimateGas(vmenv, tx)
	if err != nil {
		return err
	}
	reply.GasUsed = gasUsed

	return nil
}
