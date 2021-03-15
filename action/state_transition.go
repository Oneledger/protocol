package action

import (
	"fmt"
	"math"
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *ethcore.GasPool
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      ethvm.StateDB
	evm        *ethvm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() ethcmn.Address
	To() *ethcmn.Address

	GasPrice() *big.Int
	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte
	AccessList() ethtypes.AccessList
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas         uint64 // Total used gas but include the refunded gas
	Err             error  // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData      []byte // Returned data from evm(function result or data supplied with revert opcode)
	ContractAddress keys.Address
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	return result.Err
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Err != nil {
		return nil
	}
	return ethcmn.CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.Err != ethvm.ErrExecutionReverted {
		return nil
	}
	return ethcmn.CopyBytes(result.ReturnData)
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, accessList ethtypes.AccessList, isContractCreation bool, isHomestead, isEIP2028 bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isContractCreation && isHomestead {
		gas = ethparams.TxGasContractCreation
	} else {
		gas = ethparams.TxGas
	}
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := ethparams.TxDataNonZeroGasFrontier
		if isEIP2028 {
			nonZeroGas = ethparams.TxDataNonZeroGasEIP2028
		}
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, ethcore.ErrGasUintOverflow
		}
		gas += nz * nonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/ethparams.TxDataZeroGas < z {
			return 0, ethcore.ErrGasUintOverflow
		}
		gas += z * ethparams.TxDataZeroGas
	}
	if accessList != nil {
		gas += uint64(len(accessList)) * ethparams.TxAccessListAddressGas
		gas += uint64(accessList.StorageKeys()) * ethparams.TxAccessListStorageKeyGas
	}
	return gas, nil
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *ethvm.EVM, msg Message, gp *ethcore.GasPool) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *ethvm.EVM, msg Message, gp *ethcore.GasPool) (*ExecutionResult, error) {
	return NewStateTransition(evm, msg, gp).TransitionDb()
}

// EstimateGas used to estimate gas with the maximum limit of the gas
func EstimateGas(evm *ethvm.EVM, msg Message, gp *ethcore.GasPool) (uint64, error) {
	return NewStateTransition(evm, msg, gp).EstimateGas()
}

// to returns the recipient of the message.
func (st *StateTransition) to() ethcmn.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return ethcmn.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) buyGas() error {
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)
	if have, want := st.state.GetBalance(st.msg.From()), mgval; have.Cmp(want) < 0 {
		return fmt.Errorf("%w: address %v have %v want %v", ethcore.ErrInsufficientFunds, st.msg.From().Hex(), have, want)
	}
	if err := st.gp.SubGas(st.msg.Gas()); err != nil {
		return err
	}
	st.gas += st.msg.Gas()

	st.initialGas = st.msg.Gas()
	st.state.SubBalance(st.msg.From(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	// Make sure this transaction's nonce is correct.
	if st.msg.CheckNonce() {
		stNonce := st.state.GetNonce(st.msg.From())
		if msgNonce := st.msg.Nonce(); stNonce < msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ethcore.ErrNonceTooHigh,
				st.msg.From().Hex(), msgNonce, stNonce)
		} else if stNonce > msgNonce {
			return fmt.Errorf("%w: address %v, tx: %d state: %d", ethcore.ErrNonceTooLow,
				st.msg.From().Hex(), msgNonce, stNonce)
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the evm execution result with following fields.
//
// - used gas:
//      total gas used (including gas being refunded)
// - returndata:
//      the returned data from evm
// - concrete execution error:
//      various **EVM** error which aborts the execution,
//      e.g. ErrOutOfGas, ErrExecutionReverted
//
// However if any consensus issue encountered, return the error directly with
// nil evm execution result.
func (st *StateTransition) TransitionDb() (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas
	// 6. caller has enough balance to cover asset transfer for **topmost** call

	// Check clauses 1-3, buy gas if everything is correct
	fmt.Println("st.gas: before preCheck: ", st.gas)
	if err := st.preCheck(); err != nil {
		return nil, err
	}
	fmt.Println("st.gas: after preCheck: ", st.gas)
	msg := st.msg
	sender := ethvm.AccountRef(msg.From())
	homestead := st.evm.ChainConfig().IsHomestead(st.evm.Context.BlockNumber)
	istanbul := st.evm.ChainConfig().IsIstanbul(st.evm.Context.BlockNumber)
	contractCreation := msg.To() == nil

	// Check clauses 4-5, subtract intrinsic gas if everything is correct
	gas, err := IntrinsicGas(st.data, st.msg.AccessList(), contractCreation, homestead, istanbul)
	fmt.Println("intristic gas: ", gas)
	if err != nil {
		return nil, err
	}
	if st.gas < gas {
		return nil, fmt.Errorf("%w: have %d, want %d", ethcore.ErrIntrinsicGas, st.gas, gas)
	}
	st.gas -= gas
	fmt.Println("st.gas: after minus", st.gas)

	// Check clause 6
	if msg.Value().Sign() > 0 && !st.evm.Context.CanTransfer(st.state, msg.From(), msg.Value()) {
		return nil, fmt.Errorf("%w: address %v", ethcore.ErrInsufficientFundsForTransfer, msg.From().Hex())
	}

	// Set up the initial access list.
	if st.evm.ChainConfig().IsBerlin(st.evm.Context.BlockNumber) {
		st.state.PrepareAccessList(msg.From(), msg.To(), st.evm.ActivePrecompiles(), msg.AccessList())
	}

	var (
		ret   []byte
		vmerr error // vm errors do not effect consensus and are therefore not assigned to err
		ca    ethcmn.Address
	)
	if contractCreation {
		ret, ca, st.gas, vmerr = st.evm.Create(sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(msg.From(), st.state.GetNonce(msg.From())+1)
		ret, st.gas, vmerr = st.evm.Call(sender, st.to(), st.data, st.gas, st.value)
	}
	fmt.Println("gas after evaluation: ", st.gas)
	st.refundGas()
	fmt.Println("gas after refund: ", st.gas)

	result := &ExecutionResult{
		UsedGas:    st.gasUsed(),
		Err:        vmerr,
		ReturnData: ret,
	}
	fmt.Println("used gas: ", result.UsedGas)
	if contractCreation {
		result.ContractAddress = keys.Address(ca.Bytes())
	}
	return result, nil
}

func (st *StateTransition) EstimateGas() (uint64, error) {
	if err := st.preCheck(); err != nil {
		return 0, err
	}
	msg := st.msg
	sender := ethvm.AccountRef(msg.From())
	homestead := st.evm.ChainConfig().IsHomestead(st.evm.Context.BlockNumber)
	istanbul := st.evm.ChainConfig().IsIstanbul(st.evm.Context.BlockNumber)
	contractCreation := msg.To() == nil

	gas, err := IntrinsicGas(st.data, st.msg.AccessList(), contractCreation, homestead, istanbul)
	if err != nil {
		return 0, err
	}
	if st.gas < gas {
		return 0, fmt.Errorf("%w: have %d, want %d", ethcore.ErrIntrinsicGas, st.gas, gas)
	}
	st.gas -= gas

	if msg.Value().Sign() > 0 && !st.evm.Context.CanTransfer(st.state, msg.From(), msg.Value()) {
		return 0, fmt.Errorf("%w: address %v", ethcore.ErrInsufficientFundsForTransfer, msg.From().Hex())
	}

	if st.evm.ChainConfig().IsBerlin(st.evm.Context.BlockNumber) {
		st.state.PrepareAccessList(msg.From(), msg.To(), st.evm.ActivePrecompiles(), msg.AccessList())
	}

	var vmerr error

	if contractCreation {
		_, _, st.gas, vmerr = st.evm.Create(sender, st.data, st.gas, st.value)
	} else {
		_, st.gas, vmerr = st.evm.Call(sender, st.to(), st.data, st.gas, st.value)
	}
	st.refundGas()
	return st.gasUsed(), vmerr
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(st.msg.From(), remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	fmt.Println("initialGas gas: ", st.initialGas)
	return st.initialGas - st.gas
}
