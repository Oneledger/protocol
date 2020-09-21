package vpart_action

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_data"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type InsertPart struct {
	VIN           vpart_data.Vin         `json:"vin"`
	PartType      string                 `json:"partType"`
	DealerName    string                 `json:"dealerName"`
	DealerAddress string                 `json:"dealerAddress"`
	StockNum      vpart_data.StockNumber `json:"stockNum"`
	Year          int                    `json:"year"`
	Operator      keys.Address           `json:"operator"`
}

type InsertPartTx struct {
}

var _ action.Msg = &InsertPart{}

var _ action.Tx = &InsertPartTx{}

func (i InsertPart) Signers() []action.Address {
	return []action.Address{i.Operator}
}

func (i InsertPart) Type() action.Type {
	return VPART_INSERT
}

func (i InsertPart) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.VIN"),
		Value: []byte(i.VIN),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.partType"),
		Value: []byte(i.PartType),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(i.Type().String()),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (i InsertPart) Marshal() ([]byte, error) {
	return json.Marshal(i)
}

func (i *InsertPart) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, i)
}

func (i InsertPartTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	insertPart := InsertPart{}
	err := insertPart.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(ErrFailedToUnmarshal, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), insertPart.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if VIN is valid
	err = insertPart.VIN.Err()
	if err != nil {
		return false, vpart_data.ErrInvalidVIN.Wrap(err)
	}

	//Check if farm ID is valid
	err = insertPart.StockNum.Err()
	if err != nil {
		return false, vpart_data.ErrInvalidStockNum.Wrap(err)
	}

	//Check if operator address is valid oneLedger address
	err = insertPart.Operator.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (i InsertPartTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (i InsertPartTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck InsertPartTx transaction for CheckTx", tx)
	return runInsertPart(ctx, tx)
}

func (i InsertPartTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver InsertPartTx transaction for DeliverTx", tx)
	return runInsertPart(ctx, tx)
}

func runInsertPart(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	insertPart := InsertPart{}
	err := insertPart.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrFailedToUnmarshal, insertPart.Tags(), err)
	}

	//1. get vPart store
	vPartStore, err := GetVPartStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrGettingVPartStore, insertPart.Tags(), err)
	}

	//2. check if this part is already in the db
	if vPartStore.Exists(insertPart.VIN, insertPart.PartType) {
		return helpers.LogAndReturnFalse(ctx.Logger, vpart_data.ErrVPartAlreadyExists, insertPart.Tags(), err)
	}

	//3. construct new part
	vPart := vpart_data.NewVPart(
		insertPart.VIN,
		insertPart.PartType,
		insertPart.DealerName,
		insertPart.DealerAddress,
		insertPart.StockNum,
		insertPart.Year,
	)

	//4. insert the part
	err = vPartStore.Set(vPart)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, vpart_data.ErrInsertingPart, insertPart.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, insertPart.Tags(), "insert_vpart_success")
}