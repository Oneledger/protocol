package action

import (
	"encoding/hex"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

const (
	SkipFee  = -1
	WrongFee = 0
)

type MsgData []byte

type Msg interface {
	// return the necessary signers for the message, should have consistent order across the network
	Signers() []Address

	Type() Type

	Tags() kv.Pairs

	Marshal() ([]byte, error)

	Unmarshal([]byte) error
}

type Fee struct {
	Price Amount `json:"price"`
	Gas   int64  `json:"gas"`
}

type Signature struct {
	Signer keys.PublicKey
	Signed []byte
}

type TxTypeDescribe struct {
	TxTypeNum    Type
	TxTypeString string
}

func (s Signature) Verify(msg []byte) bool {
	handler, err := s.Signer.GetHandler()
	if err != nil {
		return false
	}
	return handler.VerifyBytes(msg, s.Signed)
}

type RawTx struct {
	Type Type    `json:"type"`
	Data MsgData `json:"data"`
	Fee  Fee     `json:"fee"`
	Memo string  `json:"memo"`
}

func (t *RawTx) RawBytes() []byte {
	value, err := serialize.GetSerializer(serialize.NETWORK).Serialize(t)
	if err != nil {
		logger.Error("failed to serialize tx: ", t)
	}
	return value
}

type SignedTx struct {
	RawTx
	Signatures []Signature `json:"signatures"`
}

func (t *SignedTx) SignedBytes() []byte {
	value, err := serialize.GetSerializer(serialize.NETWORK).Serialize(t)
	if err != nil {
		logger.Error("failed to serialize tx: ", t)
	}
	return value
}

func ValidateBasic(data []byte, signerAddr []Address, signatures []Signature) error {
	if len(signatures) != len(signerAddr) {
		return ErrUnmatchSigner
	}
	for i, s := range signerAddr {
		pkey := signatures[i].Signer
		h, err := pkey.GetHandler()
		if err != nil {
			return ErrInvalidPubkey
		}
		if !h.Address().Equal(s) {
			return errors.Wrap(ErrUnmatchSigner, hex.EncodeToString(h.Address())+","+hex.EncodeToString(s))
		}

		if !h.VerifyBytes(data, signatures[i].Signed) {
			return ErrInvalidSignature
		}
	}

	return nil
}

func ValidateFee(feeOpt *fees.FeeOption, fee Fee) error {
	if fee.Price.Currency != feeOpt.FeeCurrency.Name {
		return ErrInvalidFeeCurrency
	}
	minFee := feeOpt.MinFee()
	if minFee.Amount.BigInt().Cmp(fee.Price.Value.BigInt()) > 0 {
		return ErrInvalidFeePrice
	}
	return nil
}

func StakingPayerFeeHandling(ctx *Context, feePayer keys.Address, signedTx SignedTx, start Gas, size Gas, signatureCnt Gas) (bool, Response) {
	ctx.State.ConsumeVerifySigGas(signatureCnt)
	ctx.State.ConsumeStorageGas(size)
	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	if used > signedTx.Fee.Gas {
		return false, Response{Log: ErrGasOverflow.Error(), GasWanted: signedTx.Fee.Gas, GasUsed: signedTx.Fee.Gas}
	}

	signer := signedTx.Signatures[0].Signer
	h, err := signer.GetHandler()
	if err != nil {
		return false, Response{Log: err.Error()}
	}
	addr := h.Address()

	val, err := ctx.Validators.Get(addr)
	if err != nil {
		return false, Response{Log: err.Error()}
	}

	if !val.Address.Equal(feePayer) {
		return false, Response{Log: errors.Wrap(err, "wrong fee payer").Error()}
	}

	charge := signedTx.Fee.Price.ToCoin(ctx.Currencies).MultiplyInt64(int64(used))
	err = ctx.Balances.MinusFromAddress(val.StakeAddress, charge)
	if err != nil {
		return false, Response{Log: errors.Wrap(err, "charge fee").Error()}
	}
	err = ctx.FeePool.AddToPool(charge)
	if err != nil {
		return false, Response{Log: err.Error()}
	}
	return true, Response{GasWanted: signedTx.Fee.Gas, GasUsed: used}
}

func ContractFeeHandling(ctx *Context, signedTx SignedTx, gasUsed Gas, start Gas) (bool, Response) {
	switch gasUsed {
	case SkipFee:
		return true, Response{}
	case WrongFee:
		return false, Response{Log: ErrInvalidVmExecution.Marshal(), GasWanted: signedTx.Fee.Gas}
	}

	ctx.State.ConsumeContractGas(gasUsed)
	ctx.Logger.Detailf("Gas - used: %d, limit: %d\n", gasUsed, signedTx.Fee.Gas)

	if int64(gasUsed) > signedTx.Fee.Gas {
		return false, Response{Log: ErrGasOverflow.Marshal(), GasWanted: signedTx.Fee.Gas, GasUsed: int64(gasUsed)}
	}

	charge := signedTx.Fee.Price.ToCoin(ctx.Currencies).MultiplyInt(int(gasUsed))

	err := ctx.FeePool.AddToPool(charge)
	if err != nil {
		return false, Response{Log: ErrInvalidVmExecution.Wrap(err).Marshal(), GasWanted: signedTx.Fee.Gas}
	}
	ctx.Logger.Detailf("Add to pool %s\n", charge)
	// NOTE: For the smart contract we need to eat used gas to prevent redundant vm execution of attacker
	return true, Response{GasWanted: signedTx.Fee.Gas, GasUsed: int64(gasUsed)}
}

func BasicFeeHandling(ctx *Context, signedTx SignedTx, start Gas, size Gas, signatureCnt Gas) (bool, Response) {
	ctx.State.ConsumeVerifySigGas(signatureCnt)
	ctx.State.ConsumeStorageGas(size)
	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	if used > signedTx.Fee.Gas {
		return false, Response{Log: ErrGasOverflow.Error(), GasWanted: signedTx.Fee.Gas, GasUsed: signedTx.Fee.Gas}
	}

	// only charge the first signer for now
	signer := signedTx.Signatures[0].Signer
	h, err := signer.GetHandler()
	if err != nil {
		return false, Response{Log: err.Error()}
	}
	addr := h.Address()

	charge := signedTx.Fee.Price.ToCoin(ctx.Currencies).MultiplyInt64(int64(used))
	err = ctx.Balances.MinusFromAddress(addr, charge)
	if err != nil {
		return false, Response{Log: errors.Wrap(err, "charge fee").Error()}
	}
	err = ctx.FeePool.AddToPool(charge)
	if err != nil {
		return false, Response{Log: err.Error()}
	}
	return true, Response{GasWanted: signedTx.Fee.Gas, GasUsed: used}
}

func GetEvent(pairs kv.Pairs, eventType string) []types.Event {
	var eventList []types.Event
	event := types.Event{
		Type:       eventType,
		Attributes: pairs,
	}

	return append(eventList, event)
}
