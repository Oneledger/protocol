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

type InsertProduct struct {
	BatchId         farm_data.BatchID  `json:"batchId"`
	ItemType        farm_data.ItemType `json:"itemType"`
	FarmID          farm_data.FarmID   `json:"farmId"`
	FarmName        string             `json:"farmName"`
	HarvestLocation string             `json:"harvestLocation"`
	HarvestDate     int64              `json:"harvestDate"`
	Classification  string             `json:"classification"`
	Quantity        int                `json:"quantity"`
	Description     string             `json:"description"`
	Operator        keys.Address       `json:"operator"`
}

type InsertProductTx struct {
}

var _ action.Msg = &InsertProduct{}

var _ action.Tx = &InsertProductTx{}

func (i InsertProduct) Signers() []action.Address {
	return []action.Address{i.Operator}
}

func (i InsertProduct) Type() action.Type {
	return FARM_INSERT_PRODUCT
}

func (i InsertProduct) Tags() kv.Pairs {
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

func (i InsertProduct) Marshal() ([]byte, error) {
	return json.Marshal(i)
}

func (i InsertProduct) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, i)
}

func (i InsertProductTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	insertProduct := InsertProduct{}
	err := insertProduct.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), insertProduct.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if batch ID is valid
	if insertProduct.BatchId.Err() != nil {
		return false, farm_data.ErrInvalidBatchID
	}

	//Check if farm ID is valid
	if insertProduct.FarmID.Err() != nil {
		return false, farm_data.ErrInvalidFarmID
	}

	//Check if operator address is valid oneLedger address
	err = insertProduct.Operator.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (i InsertProductTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (i InsertProductTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck CancelProposalTx transaction for CheckTx", tx)
	return runInsertProduct(ctx, tx)
}

func (i InsertProductTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver CancelProposalTx transaction for DeliverTx", tx)
	return runInsertProduct(ctx, tx)
}

func runInsertProduct(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	insertProduct := InsertProduct{}
	err := insertProduct.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, insertProduct.Tags(), err)
	}

	//1. get product store
	productStore, err := GetProductStore(ctx)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrGettingProductStore, insertProduct.Tags(), err)
	}

	//2. check if there is product batch with same batch ID
	if productStore.Exists(insertProduct.BatchId) {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrBatchIDAlreadyExists, insertProduct.Tags(), err)
	}

	//3. construct new product batch
	productBatch := farm_data.NewProduct(
		insertProduct.BatchId,
		insertProduct.ItemType,
		insertProduct.FarmID,
		insertProduct.FarmName,
		insertProduct.HarvestLocation,
		insertProduct.HarvestDate,
		insertProduct.Classification,
		insertProduct.Quantity,
		insertProduct.Description,
	)

	//4. insert the product batch
	err = productStore.Set(productBatch)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, farm_data.ErrInsertingProduct, insertProduct.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, insertProduct.Tags(), "insert_product_success")
}
