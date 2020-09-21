package farm_action

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type InsertProduce struct {
	BatchId         farm_data.BatchID `json:"batchId"`
	ItemType        string            `json:"itemType"`
	FarmID          farm_data.FarmID  `json:"farmId"`
	FarmName        string            `json:"farmName"`
	HarvestLocation string            `json:"harvestLocation"`
	HarvestDate     int64             `json:"harvestDate"`
	Classification  string            `json:"classification"`
	Quantity        int               `json:"quantity"`
	Description     string            `json:"description"`
	Operator        keys.Address      `json:"operator"`
}

type InsertProduceTx struct {
}

var _ action.Msg = &InsertProduce{}

var _ action.Tx = &InsertProduceTx{}

func (i InsertProduce) Signers() []action.Address {
	return []action.Address{i.Operator}
}

func (i InsertProduce) Type() action.Type {
	return FARM_INSERT_PRODUCE
}

func (i InsertProduce) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.batchId"),
		Value: []byte(i.BatchId),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(i.Type().String()),
	}

	tags = append(tags, tag, tag1)
	return tags
}

func (i InsertProduce) Marshal() ([]byte, error) {
	return json.Marshal(i)
}

func (i *InsertProduce) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, i)
}

func (i InsertProduceTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	insertProduce := InsertProduce{}
	err := insertProduce.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(ErrFailedToUnmarshal, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), insertProduce.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if batch ID is valid
	err = insertProduce.BatchId.Err()
	if err != nil {
		return false, farm_data.ErrInvalidBatchID.Wrap(err)
	}

	//Check if farm ID is valid
	err = insertProduce.FarmID.Err()
	if err != nil {
		return false, farm_data.ErrInvalidFarmID.Wrap(err)
	}

	//Check if operator address is valid oneLedger address
	err = insertProduce.Operator.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (i InsertProduceTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (i InsertProduceTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck InsertProduceTx transaction for CheckTx", tx)
	return runInsertProduce(ctx, tx)
}

func (i InsertProduceTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver InsertProduceTx transaction for DeliverTx", tx)
}

func runInsertProduce(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	insertProduce := InsertProduce{}
	err := insertProduce.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrFailedToUnmarshal, insertProduce.Tags(), err)
	}

	//1. get produce store
	produceStore, err := GetProduceStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrGettingProduceStore, insertProduce.Tags(), err)
	}

	//2. check if there is produce batch with same batch ID
	if produceStore.Exists(insertProduce.BatchId) {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrBatchIDAlreadyExists, insertProduce.Tags(), err)
	}

	//3. construct new produce batch
	produceBatch := farm_data.NewProduce(
		insertProduce.BatchId,
		insertProduce.ItemType,
		insertProduce.FarmID,
		insertProduce.FarmName,
		insertProduce.HarvestLocation,
		insertProduce.HarvestDate,
		insertProduce.Classification,
		insertProduce.Quantity,
		insertProduce.Description,
	)

	//4. insert the produce batch
	err = produceStore.Set(produceBatch)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrInsertingProduce, insertProduce.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, insertProduce.Tags(), "insert_produce_success")
}
