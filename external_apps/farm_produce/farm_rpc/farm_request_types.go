package farm_rpc

import "github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"

type GetBatchByIDRequest struct {
	BatchID farm_data.BatchID `json:"batchId"`
}

type GetBatchByIDReply struct {
	ProductBatch farm_data.Product `json:"productBatch"`
	Height       int64             `json:"height"`
}
